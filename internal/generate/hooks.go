package generate

import (
	"fmt"
	"strings"

	"github.com/HeytalePazguato/plc-gitignore/internal/rules"
)

// PreCommitHook renders a POSIX-sh pre-commit script that refuses to
// commit files matching any of the vendor's known-bad patterns. The
// script is meant to be saved as `.plc-gitignore-hooks/pre-commit` and
// symlinked into `.git/hooks/`.
func PreCommitHook(r rules.RuleSet, opts Options) string {
	var b strings.Builder
	b.WriteString("#!/bin/sh\n")
	b.WriteString(fmt.Sprintf("# pre-commit hook for %s (plc-gitignore v%s)\n", r.DisplayName, defaultVersion(opts.Version)))
	b.WriteString("# Auto-generated. Symlink this into .git/hooks/pre-commit:\n")
	b.WriteString("#   ln -sf ../../.plc-gitignore-hooks/pre-commit .git/hooks/pre-commit\n\n")
	b.WriteString("set -e\n\n")
	b.WriteString("staged=\"$(git diff --cached --name-only)\"\n")
	b.WriteString("violations=\"\"\n\n")
	for _, w := range r.HookWarnings {
		// Translate the gitignore-style pattern into a grep regex.
		regex := patternToRegex(w.Pattern)
		b.WriteString(fmt.Sprintf("# %s\n", w.Reason))
		b.WriteString(fmt.Sprintf("hit=\"$(printf '%%s\\n' \"$staged\" | grep -E '%s' || true)\"\n", regex))
		b.WriteString("if [ -n \"$hit\" ]; then\n")
		b.WriteString(fmt.Sprintf("  violations=\"$violations\\n--- %s (%s):\\n$hit\"\n", w.Pattern, escape(w.Reason)))
		b.WriteString("fi\n\n")
	}
	b.WriteString("if [ -n \"$violations\" ]; then\n")
	b.WriteString("  printf 'plc-gitignore pre-commit: refusing to commit known-bad files:%s\\n' \"$violations\" >&2\n")
	b.WriteString("  printf '\\nTo override (NOT RECOMMENDED): git commit --no-verify\\n' >&2\n")
	b.WriteString("  exit 1\n")
	b.WriteString("fi\n")
	return b.String()
}

func defaultVersion(v string) string {
	if v == "" {
		return "dev"
	}
	return v
}

func escape(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}

// patternToRegex converts a small subset of gitignore-style patterns
// into an extended regex usable by `grep -E`. Mirrors check.matchPattern
// in spirit but produces regex rather than matching directly.
func patternToRegex(p string) string {
	// "_Boot/"     -> "(^|/)_Boot(/|$)"
	// "*.TcLIDs"   -> "\.TcLIDs$"
	// "*.tmc"      -> "\.tmc$"
	if strings.HasSuffix(p, "/") {
		name := strings.TrimSuffix(p, "/")
		return "(^|/)" + regexEscape(name) + "(/|$)"
	}
	if strings.HasPrefix(p, "*.") {
		return regexEscape(p[1:]) + "$"
	}
	return regexEscape(p)
}

func regexEscape(s string) string {
	specials := `\.+?^$()[]{}|`
	var b strings.Builder
	for _, c := range s {
		if strings.ContainsRune(specials, c) {
			b.WriteByte('\\')
		}
		b.WriteRune(c)
	}
	return b.String()
}
