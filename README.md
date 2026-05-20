# plc-gitignore

Generate opinionated `.gitignore` and `.gitattributes` for PLC projects.
Supports **TwinCAT**, **Codesys**, **B&R Automation Studio**, **Siemens
TIA Portal**, and **Rockwell Studio 5000**. Includes an audit command
that finds files you forgot to ignore.

## Why this exists

GitHub's official [gitignore templates repo](https://github.com/github/gitignore)
covers hundreds of languages and tools. It has nothing for PLC projects.
Every shop that starts using git for PLC code hand-rolls a `.gitignore`
and gets it wrong — committing `_Boot/` directories, `.TcLIDs`
breakpoint metadata, regenerated `.tmc` type caches, and per-user
Visual Studio state.

`plc-gitignore` ships the opinionated, comment-annotated ignore rules
every vendor needs but nobody publishes.

## Install

### Homebrew (macOS / Linux)

```sh
brew install HeytalePazguato/tap/plc-gitignore
```

### Scoop (Windows)

```sh
scoop bucket add HeytalePazguato https://github.com/HeytalePazguato/scoop-bucket
scoop install plc-gitignore
```

### Go

```sh
go install github.com/HeytalePazguato/plc-gitignore/cmd/plc-gitignore@latest
```

### Pre-built binaries

Grab the latest release for your OS/arch from
[Releases](https://github.com/HeytalePazguato/plc-gitignore/releases).

## Usage

```sh
plc-gitignore init --vendor twincat              # generate for TwinCAT
plc-gitignore init --vendor codesys              # generate for Codesys
plc-gitignore init --vendor br                   # generate for B&R
plc-gitignore init --vendor siemens              # generate for Siemens TIA
plc-gitignore init --vendor rockwell             # generate for Rockwell

plc-gitignore init --vendor twincat --with-hmi   # include HMI rules
plc-gitignore init --vendor twincat --with-lfs   # include Git LFS config
plc-gitignore init --vendor twincat --with-hooks # generate pre-commit hook

plc-gitignore check                              # scan repo for missed files
plc-gitignore doctor                             # audit existing .gitignore
```

`check` and `doctor` auto-detect the vendor from file extensions; pass
`--vendor` to override.

## Output

`plc-gitignore init` writes two files (three with `--with-hooks`):

| File | Purpose |
|---|---|
| `.gitignore`     | Vendor-specific ignore rules with inline comments explaining each section. |
| `.gitattributes` | Merge strategy (`merge=union` for project XML), binary classifications, LF normalization. |
| `.plc-gitignore-hooks/pre-commit` | Optional. Refuses commits matching known-bad patterns. Symlink into `.git/hooks/`. |

Every rule has an inline comment explaining *why* it's there.

## Vendors

| Vendor | Detected by | Headline rules |
|---|---|---|
| **TwinCAT 3**   | `*.tsproj`, `*.plcproj`, `*.TcPOU` | `_Boot/`, `_CompileInfo/`, `*.TcLIDs`, `*.tmc`, `.vs/`, `*.suo`, hardware-scan cache |
| **Codesys**     | `*.project`, `*.library`           | `*.compileinfo`, `*.object`, `*.app`, `.metadata/`, Python script caches |
| **B&R**         | `*.apj`, `Physical.pkg`            | `Binaries/`, `Temp/`, `*.isopen`, `LastUser.set`, CPU build cache |
| **Siemens TIA** | `*.ap17`, `*.ap18`, `*.ap19`, `*.ap20` | `.Siemens.Automation.Objectframe.FileStorage/`, `UserFiles/`, archive files |
| **Rockwell**    | `*.ACD`, `*.L5X`, `*.L5K`          | `*.BAK`, `*.CRC`, `*.RSS`, `*.WRK`, `*.SEM` |

## Commands

### `init`

Generates `.gitignore` + `.gitattributes` for the chosen vendor. Refuses
to overwrite existing files unless `--force` is passed.

### `check`

Walks the working tree and reports files matching the vendor's ignore
patterns. With `--fix`, appends missing patterns to `.gitignore` and
runs `git rm --cached` on each finding.

### `doctor`

Reads the existing `.gitignore` (and `.gitattributes`) and reports
which recommended rules are present, missing, or partially configured.
Emits a score: `4/8 rules present`. Exit code 1 when score is incomplete.

## Branch flow

This project follows the blueprint in `BLUEPRINT.md`:

```
develop  →  release/<version>  →  main          (planned releases)
                hotfix/<version> →  main         (emergency patches)
```

Daily work lands on `develop`. Releases ship via `release/*` branches.

## License

MIT — see [LICENSE](LICENSE).
