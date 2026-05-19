// plc-gitignore generates opinionated .gitignore and .gitattributes
// for PLC projects (TwinCAT, Codesys, B&R, Siemens TIA, Rockwell).
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/HeytalePazguato/plc-gitignore/internal/check"
	"github.com/HeytalePazguato/plc-gitignore/internal/detect"
	"github.com/HeytalePazguato/plc-gitignore/internal/doctor"
	"github.com/HeytalePazguato/plc-gitignore/internal/generate"
	"github.com/HeytalePazguato/plc-gitignore/internal/rules"
)

// Stamped at build time via -ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "init":
		os.Exit(runInit(os.Args[2:]))
	case "check":
		os.Exit(runCheck(os.Args[2:]))
	case "doctor":
		os.Exit(runDoctor(os.Args[2:]))
	case "version", "--version", "-v":
		fmt.Printf("plc-gitignore %s (%s, %s)\n", version, commit, date)
	case "help", "--help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `plc-gitignore — opinionated .gitignore for PLC projects

USAGE
  plc-gitignore <command> [flags]

COMMANDS
  init      Generate .gitignore and .gitattributes for a vendor
  check     Scan the current repo for files that should be ignored
  doctor    Audit existing .gitignore against recommended rules
  version   Print version
  help      Show this message

SUPPORTED VENDORS
  %s

EXAMPLES
  plc-gitignore init --vendor twincat
  plc-gitignore init --vendor twincat --with-hmi --with-lfs --with-hooks
  plc-gitignore check
  plc-gitignore doctor

Run "plc-gitignore <command> --help" for command-specific flags.
`, strings.Join(rules.VendorIDs(), ", "))
}

func runInit(args []string) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	vendor := fs.String("vendor", "", "vendor (required): "+strings.Join(rules.VendorIDs(), ", "))
	dir := fs.String("dir", ".", "output directory")
	force := fs.Bool("force", false, "overwrite existing .gitignore/.gitattributes")
	withHMI := fs.Bool("with-hmi", false, "include HMI-specific rules")
	withLFS := fs.Bool("with-lfs", false, "include Git LFS configuration in .gitattributes")
	withHooks := fs.Bool("with-hooks", false, "generate a pre-commit hook in .plc-gitignore-hooks/")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *vendor == "" {
		fmt.Fprintln(os.Stderr, "error: --vendor is required")
		fs.Usage()
		return 2
	}
	r, err := rules.ByVendor(*vendor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	opts := generate.Options{Version: version, WithHMI: *withHMI, WithLFS: *withLFS}
	gi := generate.Gitignore(r, opts)
	ga := generate.Gitattributes(r, opts)

	if err := writeFile(filepath.Join(*dir, ".gitignore"), gi, *force); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if err := writeFile(filepath.Join(*dir, ".gitattributes"), ga, *force); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	fmt.Printf("Generated .gitignore and .gitattributes for %s in %s\n", r.DisplayName, *dir)

	if *withHooks {
		hookDir := filepath.Join(*dir, ".plc-gitignore-hooks")
		if err := os.MkdirAll(hookDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		hookPath := filepath.Join(hookDir, "pre-commit")
		if err := os.WriteFile(hookPath, []byte(generate.PreCommitHook(r, opts)), 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		fmt.Printf("Generated pre-commit hook at %s\n", hookPath)
		fmt.Println("  Activate with: ln -sf ../../.plc-gitignore-hooks/pre-commit .git/hooks/pre-commit")
	}
	return 0
}

func writeFile(path, contents string, force bool) error {
	if _, err := os.Stat(path); err == nil && !force {
		return fmt.Errorf("%s exists (pass --force to overwrite)", path)
	}
	return os.WriteFile(path, []byte(contents), 0o644)
}

func resolveVendor(dir, override string) (rules.RuleSet, error) {
	if override != "" {
		return rules.ByVendor(override)
	}
	r, score, err := detect.FromDir(dir)
	if err != nil {
		return rules.RuleSet{}, err
	}
	if r.Vendor == "" {
		return rules.RuleSet{}, fmt.Errorf("could not auto-detect vendor; pass --vendor explicitly (supported: %s)", strings.Join(rules.VendorIDs(), ", "))
	}
	fmt.Fprintf(os.Stderr, "Detected vendor: %s (matched %d file(s))\n", r.DisplayName, score)
	return r, nil
}

func runCheck(args []string) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	vendor := fs.String("vendor", "", "vendor (auto-detected if omitted)")
	dir := fs.String("dir", ".", "repository root to scan")
	fix := fs.Bool("fix", false, "add missing patterns to .gitignore and untrack offending files")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	r, err := resolveVendor(*dir, *vendor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	findings, err := check.Scan(*dir, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if len(findings) == 0 {
		fmt.Println("OK — no offending files found.")
		return 0
	}
	fmt.Printf("Found %d file(s) that should be ignored:\n\n", len(findings))
	for _, f := range findings {
		fmt.Printf("  %-50s (%s)\n", f.Path, f.Reason)
	}
	if *fix {
		if err := check.Fix(*dir, r, findings); err != nil {
			fmt.Fprintf(os.Stderr, "fix error: %v\n", err)
			return 1
		}
		for _, f := range findings {
			cmd := exec.Command("git", "-C", *dir, "rm", "--cached", "-f", "--ignore-unmatch", f.Path)
			_ = cmd.Run()
		}
		fmt.Printf("\nApplied fixes: patterns appended to .gitignore, files untracked.\n")
		return 0
	}
	fmt.Printf("\nRun: plc-gitignore check --vendor %s --fix\n", r.Vendor)
	fmt.Println("  to add these patterns to .gitignore and remove the files from tracking.")
	return 1
}

func runDoctor(args []string) int {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	vendor := fs.String("vendor", "", "vendor (auto-detected if omitted)")
	dir := fs.String("dir", ".", "repository root to audit")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	r, err := resolveVendor(*dir, *vendor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	report, err := doctor.Audit(*dir, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	fmt.Printf("Detected vendor: %s\n\n", report.Vendor.DisplayName)
	for _, res := range report.Results {
		marker := "?"
		switch res.Status {
		case doctor.StatusOK:
			marker = "[ok] "
		case doctor.StatusMissing:
			marker = "[!!] "
		case doctor.StatusWarn:
			marker = "[~~] "
		}
		fmt.Printf("%s %-40s %s\n", marker, res.Pattern, res.Reason)
	}
	fmt.Printf("\nScore: %d/%d rules present. ", report.PresentCount, report.TotalCount)
	if report.PresentCount < report.TotalCount {
		fmt.Printf("Run 'plc-gitignore init --vendor %s' to generate a complete config.\n", report.Vendor.Vendor)
		return 1
	}
	fmt.Println("All recommended rules are present.")
	return 0
}
