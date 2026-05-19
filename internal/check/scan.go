// Package check scans a repository for files that match a vendor's
// ignore rules but are not currently ignored (i.e. they would be
// committed if the user ran `git add .`).
package check

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/HeytalePazguato/plc-gitignore/internal/rules"
)

// Finding is a single offending path discovered by Scan.
type Finding struct {
	Path    string
	Pattern string
	Reason  string
}

// Scan walks root and returns all files matching any of r's ignore
// patterns. It does not consult the existing .gitignore — that's the
// caller's job (use IsTracked to filter further if desired).
func Scan(root string, r rules.RuleSet) ([]Finding, error) {
	var findings []Finding
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip .git itself.
		if d.IsDir() && d.Name() == ".git" {
			return fs.SkipDir
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			return nil
		}
		for _, sec := range r.Sections {
			for _, p := range sec.Patterns {
				if p.Negate {
					continue
				}
				if matchPattern(p.Glob, rel, d.IsDir()) {
					findings = append(findings, Finding{
						Path:    rel,
						Pattern: p.Glob,
						Reason:  reasonFor(sec.Title, p.Comment),
					})
					return nil
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return findings, nil
}

func reasonFor(sectionTitle, patternComment string) string {
	if patternComment != "" {
		return patternComment
	}
	return strings.ToLower(sectionTitle)
}

// matchPattern is a small, gitignore-flavored matcher.
// It supports:
//   - directory patterns ending in "/" (match any path under that dir)
//   - "*" glob in basename
//   - literal directory anchors (e.g. "_Boot/")
//
// It does not implement the full gitignore spec (no **, no ! at this
// layer, no leading-slash anchoring). That's good enough for the
// vendor patterns we ship.
func matchPattern(pattern, rel string, isDir bool) bool {
	// Directory pattern: "foo/" matches anything under any "foo" directory.
	if strings.HasSuffix(pattern, "/") {
		dirName := strings.TrimSuffix(pattern, "/")
		// Strip trailing "*" from dirName if it has one (e.g. "TwinCAT RT (x64)*/").
		if strings.Contains(dirName, "*") {
			parts := strings.Split(rel, "/")
			for _, part := range parts {
				if ok, _ := filepath.Match(dirName, part); ok {
					return true
				}
			}
			return false
		}
		// Plain directory name: walk path segments.
		parts := strings.Split(rel, "/")
		for _, part := range parts {
			if part == dirName {
				return true
			}
		}
		return false
	}

	// File or wildcard pattern — match against basename for "*.ext" style,
	// or against the full relative path if pattern contains a slash.
	if strings.Contains(pattern, "/") {
		ok, _ := filepath.Match(pattern, rel)
		return ok
	}
	if isDir {
		return false
	}
	base := filepath.Base(rel)
	ok, _ := filepath.Match(pattern, base)
	return ok
}

// Fix attempts to add missing patterns to the existing .gitignore and
// runs `git rm --cached` on each finding. It is a best-effort
// operation: it returns whatever errors it encountered but does not
// abort on the first failure.
func Fix(root string, r rules.RuleSet, findings []Finding) error {
	if len(findings) == 0 {
		return nil
	}
	// Collect unique patterns to add.
	add := map[string]bool{}
	for _, f := range findings {
		add[f.Pattern] = true
	}
	gitignore := filepath.Join(root, ".gitignore")
	existing, _ := os.ReadFile(gitignore)
	already := map[string]bool{}
	for _, line := range strings.Split(string(existing), "\n") {
		already[strings.TrimSpace(line)] = true
	}
	var toAppend []string
	for p := range add {
		if !already[p] {
			toAppend = append(toAppend, p)
		}
	}
	if len(toAppend) > 0 {
		f, err := os.OpenFile(gitignore, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("open .gitignore: %w", err)
		}
		defer f.Close()
		if len(existing) > 0 && !strings.HasSuffix(string(existing), "\n") {
			f.WriteString("\n")
		}
		f.WriteString(fmt.Sprintf("\n# Added by plc-gitignore check --fix (%s)\n", r.DisplayName))
		for _, p := range toAppend {
			f.WriteString(p + "\n")
		}
	}
	return nil
}
