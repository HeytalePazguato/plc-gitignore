package rules

import "testing"

func TestByVendor_known(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"twincat", "twincat"},
		{"TwinCAT", "twincat"},
		{"TWINCAT", "twincat"},
		{"codesys", "codesys"},
		{"Codesys", "codesys"},
	}
	for _, c := range cases {
		r, err := ByVendor(c.in)
		if err != nil {
			t.Errorf("ByVendor(%q) returned error: %v", c.in, err)
			continue
		}
		if r.Vendor != c.want {
			t.Errorf("ByVendor(%q).Vendor = %q, want %q", c.in, r.Vendor, c.want)
		}
	}
}

func TestByVendor_unknown(t *testing.T) {
	if _, err := ByVendor("acme9000"); err == nil {
		t.Fatal("expected error for unknown vendor, got nil")
	}
}

func TestTwinCAT_essentials(t *testing.T) {
	r := TwinCAT()
	mustContainPattern(t, r, "_Boot/")
	mustContainPattern(t, r, "*.TcLIDs")
	mustContainPattern(t, r, "*.tmc")
	mustContainPattern(t, r, "*.compiled")
	mustContainPattern(t, r, ".vs/")
}

func TestCodesys_essentials(t *testing.T) {
	r := Codesys()
	mustContainPattern(t, r, "*.compileinfo")
	mustContainPattern(t, r, "*.object")
	mustContainPattern(t, r, "__pycache__/")
}

func TestTwinCAT_attributes(t *testing.T) {
	r := TwinCAT()
	mustContainAttr(t, r, "*.tsproj")
	mustContainAttr(t, r, "*.plcproj")
}

func mustContainPattern(t *testing.T, r RuleSet, want string) {
	t.Helper()
	for _, s := range r.Sections {
		for _, p := range s.Patterns {
			if p.Glob == want {
				return
			}
		}
	}
	t.Errorf("%s rule set missing pattern %q", r.DisplayName, want)
}

func mustContainAttr(t *testing.T, r RuleSet, want string) {
	t.Helper()
	for _, s := range r.Attributes {
		for _, rule := range s.Rules {
			if rule.Pattern == want {
				return
			}
		}
	}
	t.Errorf("%s rule set missing .gitattributes pattern %q", r.DisplayName, want)
}
