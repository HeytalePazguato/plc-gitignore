# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] — next: 0.1.0

## [0.0.2] - 2026-05-18

### Added

- Three new vendor rule sets: **B&R Automation Studio**, **Siemens TIA
  Portal**, **Rockwell Studio 5000**. Each ships with its own
  `.gitignore` rules, `.gitattributes` merge/binary classifications,
  HMI-specific extensions, and pre-commit hook warnings.
- `plc-gitignore doctor` command — audits an existing `.gitignore`
  against the vendor's recommended rules, reports OK / missing /
  warnings, and scores the result.
- `internal/detect` — auto-detects vendor from project file
  extensions (e.g. `*.tsproj` → TwinCAT, `*.ACD` → Rockwell). `check`
  and `doctor` use it when `--vendor` is omitted.
- `--with-hmi` and `--with-lfs` flags on `init` — opt-in extensions
  for HMI build output and Git LFS configuration.
- `--with-hooks` flag — generates `.plc-gitignore-hooks/pre-commit`
  that refuses to commit known-bad files (build artifacts, license
  blobs, user-specific caches).
- Test coverage for vendor detection, doctor audit, and pre-commit
  hook rendering.

## [0.0.1] - 2026-05-18

### Added

- Initial CLI scaffold: `plc-gitignore init --vendor <name>` generates
  `.gitignore` and `.gitattributes` for a PLC vendor.
- TwinCAT 3 rule set: `_Boot/`, `_CompileInfo/`, `*.TcLIDs`, `*.tmc`,
  `*.compiled`, `.vs/`, user-settings, hardware-scan cache, HMI build
  output (opt-in), license blobs, plus `.gitattributes` merge=union for
  `*.tsproj`/`*.plcproj` and LF normalization for `*.st`.
- Codesys rule set: `*.compileinfo`, `*.object`, `*.app`, Eclipse-style
  workspace metadata, Python scripting caches, device build cache,
  auto-save backups.
- `plc-gitignore check --vendor <name>` scans a repo for files matching
  vendor ignore patterns and reports each finding with a reason.
  `--fix` appends the missing patterns to `.gitignore` and untracks the
  files via `git rm --cached`.
- Build flag stamping (`-X main.version`, `main.commit`, `main.date`)
  so the binary responds to `version` with the spec'd format.
- Test coverage for rule essentials, both generators, scan, and fix.

<!--
  Add entries here while working on `develop` or `release/*`. The heading
  carries the next target version (`next: X.Y.Z`) so it's obvious what
  these entries will ship as. On release, rename to `[X.Y.Z] - YYYY-MM-DD`
  and start a new `[Unreleased] — next: X.Y.Z+1` block.

  Use these subsections, omitting any that don't apply:
    ### Added       — new features
    ### Changed     — changes in existing functionality
    ### Deprecated  — soon-to-be removed features
    ### Removed     — removed features
    ### Fixed       — bug fixes
    ### Security    — vulnerability fixes
-->
