// Package doctor audits an existing .gitignore against a vendor's
// recommended rules and reports which patterns are present, missing,
// or partially configured.
package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HeytalePazguato/plc-gitignore/internal/rules"
)

// CheckStatus is the outcome of a single rule check.
type CheckStatus int

const (
	StatusOK CheckStatus = iota
	StatusMissing
	StatusWarn
)

// CheckResult describes a single audit finding.
type CheckResult struct {
	Pattern string
	Reason  string
	Status  CheckStatus
}

// Report bundles every check plus a score.
type Report struct {
	Vendor       rules.RuleSet
	Results      []CheckResult
	PresentCount int
	TotalCount   int
	HasAttrFile  bool
}

// Audit reads root/.gitignore (and root/.gitattributes for attribute
// checks) and returns a Report scored against r's recommended rules.
func Audit(root string, r rules.RuleSet) (Report, error) {
	ignoreBytes, _ := os.ReadFile(filepath.Join(root, ".gitignore"))
	ignoreLines := splitLines(string(ignoreBytes))
	ignorePresent := map[string]bool{}
	for _, line := range ignoreLines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ignorePresent[line] = true
	}

	attrPath := filepath.Join(root, ".gitattributes")
	_, attrErr := os.Stat(attrPath)
	hasAttr := attrErr == nil

	var results []CheckResult
	present := 0
	for _, s := range r.Sections {
		for _, p := range s.Patterns {
			if p.Negate {
				continue
			}
			pat := p.Glob
			cr := CheckResult{Pattern: pat, Reason: reasonFromSection(s)}
			if ignorePresent[pat] {
				cr.Status = StatusOK
				present++
			} else {
				cr.Status = StatusMissing
			}
			results = append(results, cr)
		}
	}

	if !hasAttr && len(r.Attributes) > 0 {
		results = append(results, CheckResult{
			Pattern: ".gitattributes",
			Reason:  fmt.Sprintf("missing — no merge strategy defined for %s", r.DisplayName),
			Status:  StatusMissing,
		})
	}

	// Cross-section warning: TwinCAT's *.suo without .vs/ is a smell.
	if ignorePresent["*.suo"] && !ignorePresent[".vs/"] {
		results = append(results, CheckResult{
			Pattern: ".vs/",
			Reason:  "*.suo is ignored but .vs/ is not — Visual Studio cache may be committed",
			Status:  StatusWarn,
		})
	}

	total := 0
	for _, r := range results {
		if r.Status != StatusWarn {
			total++
		}
	}

	return Report{
		Vendor:       r,
		Results:      results,
		PresentCount: present,
		TotalCount:   total,
		HasAttrFile:  hasAttr,
	}, nil
}

func reasonFromSection(s rules.Section) string {
	if s.Comment != "" {
		// First line only.
		if i := strings.IndexByte(s.Comment, '\n'); i > 0 {
			return s.Comment[:i]
		}
		return s.Comment
	}
	return strings.ToLower(s.Title)
}

func splitLines(s string) []string {
	return strings.Split(s, "\n")
}
