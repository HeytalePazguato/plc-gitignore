// Package rules defines vendor-specific PLC project rule sets.
//
// Each vendor has a RuleSet describing the .gitignore patterns,
// .gitattributes entries, and optional extensions (HMI, LFS) that
// plc-gitignore generates.
package rules

import "fmt"

// Section groups related ignore rules under a single header.
// The Comment is rendered above the patterns in the generated
// .gitignore so a human reviewing the file can understand what
// each block of rules is for and why.
type Section struct {
	Title    string
	Comment  string
	Patterns []Pattern
}

// Pattern is a single ignore rule with an optional inline comment.
// Negated patterns (those starting with "!") un-ignore a previously
// ignored path — used sparingly.
type Pattern struct {
	Glob    string
	Comment string
	Negate  bool
}

// AttrRule is a single .gitattributes entry.
type AttrRule struct {
	Pattern string
	Attrs   string
	Comment string
}

// AttrSection groups related .gitattributes entries.
type AttrSection struct {
	Title   string
	Comment string
	Rules   []AttrRule
}

// RuleSet bundles everything a single vendor needs.
type RuleSet struct {
	Vendor       string   // canonical lowercase id: "twincat", "codesys", "br", "siemens", "rockwell"
	DisplayName  string   // human-readable: "TwinCAT 3", "Codesys", etc.
	DetectGlobs  []string // file extensions/paths that indicate this vendor (e.g. "*.tsproj")
	Sections     []Section
	Attributes   []AttrSection
	HMISections  []Section     // appended when --with-hmi is passed
	LFSAttrs     []AttrSection // appended when --with-lfs is passed
	HookWarnings []HookWarning // patterns the pre-commit hook should flag
}

// HookWarning is a pattern the pre-commit hook should refuse to allow.
type HookWarning struct {
	Pattern string
	Reason  string
}

// All returns every supported vendor RuleSet.
func All() []RuleSet {
	return []RuleSet{
		TwinCAT(),
		Codesys(),
		BR(),
		Siemens(),
		Rockwell(),
	}
}

// ByVendor returns the RuleSet for the given vendor id (case-insensitive).
// It returns an error if the vendor is not recognized.
func ByVendor(id string) (RuleSet, error) {
	for _, r := range All() {
		if equalFold(r.Vendor, id) {
			return r, nil
		}
	}
	return RuleSet{}, fmt.Errorf("unknown vendor %q (supported: twincat, codesys, br, siemens, rockwell)", id)
}

// VendorIDs returns the lowercase id for every supported vendor.
func VendorIDs() []string {
	ids := make([]string, 0, len(All()))
	for _, r := range All() {
		ids = append(ids, r.Vendor)
	}
	return ids
}

// equalFold is a tiny case-insensitive compare without pulling in strings
// just for one call site.
func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
