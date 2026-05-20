package check

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HeytalePazguato/plc-gitignore/internal/rules"
)

func TestScan_findsTwinCATOffenders(t *testing.T) {
	dir := t.TempDir()
	mustMkdir(t, filepath.Join(dir, "src", "_CompileInfo"))
	mustMkdir(t, filepath.Join(dir, "src", "PLC1", "_Boot", "Plc"))
	mustMkdir(t, filepath.Join(dir, "src", "PLC1", "POUs"))
	mustWrite(t, filepath.Join(dir, "src", "_CompileInfo", "Route.xml"), "x")
	mustWrite(t, filepath.Join(dir, "src", "PLC1", "PLC1.tmc"), "x")
	mustWrite(t, filepath.Join(dir, "src", "PLC1", "_Boot", "Plc", "Port_851.bootdata"), "x")
	mustWrite(t, filepath.Join(dir, "src", "PLC1", "POUs", "FB_Main.TcLIDs"), "x")
	mustWrite(t, filepath.Join(dir, "src", "PLC1", "POUs", "FB_Main.TcPOU"), "x") // not ignored

	findings, err := Scan(dir, rules.TwinCAT())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) < 4 {
		t.Fatalf("expected at least 4 findings, got %d: %+v", len(findings), findings)
	}
	gotPatterns := map[string]bool{}
	for _, f := range findings {
		gotPatterns[f.Pattern] = true
	}
	for _, want := range []string{"_CompileInfo/", "*.tmc", "_Boot/", "*.TcLIDs"} {
		if !gotPatterns[want] {
			t.Errorf("expected to find pattern %q in findings", want)
		}
	}
}

func TestScan_cleanRepoYieldsNoFindings(t *testing.T) {
	dir := t.TempDir()
	mustMkdir(t, filepath.Join(dir, "src", "POUs"))
	mustWrite(t, filepath.Join(dir, "src", "POUs", "FB_Pump.TcPOU"), "x")
	mustWrite(t, filepath.Join(dir, "README.md"), "x")
	findings, err := Scan(dir, rules.TwinCAT())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %+v", findings)
	}
}

func TestFix_appendsPatterns(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, ".gitignore"), "# existing\n*.bak\n")
	mustMkdir(t, filepath.Join(dir, "_CompileInfo"))
	mustWrite(t, filepath.Join(dir, "_CompileInfo", "x.xml"), "x")
	findings, err := Scan(dir, rules.TwinCAT())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) == 0 {
		t.Fatal("expected findings")
	}
	if err := Fix(dir, rules.TwinCAT(), findings); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	if !contains(s, "_CompileInfo/") {
		t.Errorf("expected appended pattern in .gitignore, got:\n%s", s)
	}
	if !contains(s, "# existing") {
		t.Errorf("expected existing content preserved, got:\n%s", s)
	}
}

func mustMkdir(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
}

func mustWrite(t *testing.T, p, content string) {
	t.Helper()
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}
