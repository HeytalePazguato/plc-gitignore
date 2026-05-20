package doctor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/HeytalePazguato/plc-gitignore/internal/rules"
)

func TestAudit_emptyRepoFlagsEverything(t *testing.T) {
	dir := t.TempDir()
	rep, err := Audit(dir, rules.TwinCAT())
	if err != nil {
		t.Fatal(err)
	}
	if rep.PresentCount != 0 {
		t.Errorf("expected 0 present, got %d", rep.PresentCount)
	}
	if rep.TotalCount == 0 {
		t.Errorf("expected non-zero total count")
	}
	if rep.HasAttrFile {
		t.Errorf("expected HasAttrFile=false")
	}
}

func TestAudit_completeRepoScoresFull(t *testing.T) {
	dir := t.TempDir()
	// Write a .gitignore containing every TwinCAT pattern.
	var lines []string
	r := rules.TwinCAT()
	for _, s := range r.Sections {
		for _, p := range s.Patterns {
			if !p.Negate {
				lines = append(lines, p.Glob)
			}
		}
	}
	body := ""
	for _, l := range lines {
		body += l + "\n"
	}
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".gitattributes"), []byte("*.tsproj merge=union\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := Audit(dir, r)
	if err != nil {
		t.Fatal(err)
	}
	if rep.PresentCount != rep.TotalCount {
		t.Errorf("expected full score, got %d/%d", rep.PresentCount, rep.TotalCount)
	}
}

func TestAudit_suoWithoutVsTriggersWarn(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.suo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rep, err := Audit(dir, rules.TwinCAT())
	if err != nil {
		t.Fatal(err)
	}
	foundWarn := false
	for _, r := range rep.Results {
		if r.Status == StatusWarn && r.Pattern == ".vs/" {
			foundWarn = true
		}
	}
	if !foundWarn {
		t.Error("expected a StatusWarn for .vs/")
	}
}
