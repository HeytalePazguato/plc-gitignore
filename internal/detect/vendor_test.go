package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFromDir_twincat(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "MyProject.tsproj"), "x")
	mustWrite(t, filepath.Join(dir, "PLC1", "PLC1.plcproj"), "x")
	r, score, err := FromDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if r.Vendor != "twincat" {
		t.Errorf("expected vendor=twincat, got %q", r.Vendor)
	}
	if score < 2 {
		t.Errorf("expected score >= 2, got %d", score)
	}
}

func TestFromDir_codesys(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "App.project"), "x")
	r, _, err := FromDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if r.Vendor != "codesys" {
		t.Errorf("expected vendor=codesys, got %q", r.Vendor)
	}
}

func TestFromDir_rockwell(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "Machine.ACD"), "x")
	mustWrite(t, filepath.Join(dir, "Export.L5X"), "x")
	r, _, err := FromDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if r.Vendor != "rockwell" {
		t.Errorf("expected vendor=rockwell, got %q", r.Vendor)
	}
}

func TestFromDir_noMatch(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "README.md"), "x")
	r, _, err := FromDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if r.Vendor != "" {
		t.Errorf("expected no match, got %q", r.Vendor)
	}
}

func mustWrite(t *testing.T, p, c string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(c), 0o644); err != nil {
		t.Fatal(err)
	}
}
