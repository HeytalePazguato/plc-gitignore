package generate

import (
	"strings"
	"testing"

	"github.com/HeytalePazguato/plc-gitignore/internal/rules"
)

func TestPreCommitHook_includesAllWarnings(t *testing.T) {
	out := PreCommitHook(rules.TwinCAT(), Options{Version: "0.0.2"})
	if !strings.HasPrefix(out, "#!/bin/sh") {
		t.Errorf("hook must start with shebang, got: %q", firstLine(out))
	}
	for _, w := range rules.TwinCAT().HookWarnings {
		if !strings.Contains(out, w.Reason) {
			t.Errorf("hook missing reason %q", w.Reason)
		}
	}
}

func TestPreCommitHook_patternToRegexShapes(t *testing.T) {
	cases := map[string]string{
		"_Boot/":    "(^|/)_Boot(/|$)",
		"*.TcLIDs":  "\\.TcLIDs$",
		"*.tmc":     "\\.tmc$",
	}
	for in, want := range cases {
		got := patternToRegex(in)
		if got != want {
			t.Errorf("patternToRegex(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAllVendors_haveCoreRuleShape(t *testing.T) {
	for _, r := range rules.All() {
		if r.Vendor == "" {
			t.Errorf("rule set has empty Vendor: %+v", r)
		}
		if r.DisplayName == "" {
			t.Errorf("rule set %q missing DisplayName", r.Vendor)
		}
		if len(r.Sections) == 0 {
			t.Errorf("rule set %q has no Sections", r.Vendor)
		}
		if len(r.Attributes) == 0 {
			t.Errorf("rule set %q has no Attributes", r.Vendor)
		}
	}
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
