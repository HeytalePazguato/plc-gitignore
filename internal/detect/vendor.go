// Package detect inspects a project directory and guesses which PLC
// vendor's rules should apply, based on file extensions present.
package detect

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/HeytalePazguato/plc-gitignore/internal/rules"
)

// FromDir walks root and returns the best-matching RuleSet plus a
// confidence count (how many distinguishing files were found). If no
// vendor matched, the returned RuleSet has an empty Vendor.
func FromDir(root string) (rules.RuleSet, int, error) {
	scores := map[string]int{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "node_modules" {
				return fs.SkipDir
			}
			return nil
		}
		name := d.Name()
		for _, r := range rules.All() {
			for _, glob := range r.DetectGlobs {
				if matched(glob, name) {
					scores[r.Vendor]++
				}
			}
		}
		return nil
	})
	if err != nil {
		return rules.RuleSet{}, 0, err
	}
	bestVendor, bestScore := "", 0
	for v, s := range scores {
		if s > bestScore {
			bestVendor, bestScore = v, s
		}
	}
	if bestVendor == "" {
		return rules.RuleSet{}, 0, nil
	}
	r, err := rules.ByVendor(bestVendor)
	return r, bestScore, err
}

func matched(pattern, name string) bool {
	// strip leading "*" + extension matching: filepath.Match handles "*.ext".
	if strings.Contains(pattern, "/") {
		return false
	}
	ok, _ := filepath.Match(pattern, name)
	return ok
}
