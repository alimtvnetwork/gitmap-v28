package cmdsync

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/constants"
)

// Curated defaults sourced from the user-approved baseline (see helptext/sync.md).
// Kept inline (not embed) because these are short, hand-curated and rarely change.
const DefaultGitignoreBaseline = `# =======================================================================
# System, IDEs, and Text Editor Configurations
# =======================================================================
.DS_Store
Thumbs.db
*.bak
*.sw?
*.swp

# JetBrains / Rider / VS Code / Visual Studio
.idea/
.vs/
.vscode/*
!.vscode/extensions.json

# User-specific developer overrides
*.rsuser
*.suo
*.user
*.userosscache
*.sln.docstates
*.ntvs*
*.njsproj
*.sln

# =======================================================================
# Logging Frameworks & Taggings
# =======================================================================
logs/
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*
pnpm-debug.log*
lerna-debug.log*

# =======================================================================
# GitMap & Tooling Internal Workspace Cache
# =======================================================================
.gitmap/

# =======================================================================
# Node.js, Modern Meta-Frameworks, and Tooling
# =======================================================================
node_modules/
.pnpm-store/

# Production builds and deployment artifact targets
dist/
dist-ssr/
build/
out/
.output/
.next/

# Meta-framework runtimes (Nitro, Vinxi, TanStack, etc.)
.nitro/
.vinxi/
.tanstack/

# =======================================================================
# Cloudflare / Serverless Deployment Runtimes
# =======================================================================
.wrangler/
.dev.vars

# =======================================================================
# Go (Golang)
# =======================================================================
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.prof
bin/
pkg/

# =======================================================================
# Python Core & Packaging Ecosystem
# =======================================================================
__pycache__/
*.py[cod]
*$py.class
venv/
.venv/
env/
ENV/
pyenv/
doc/_build/
*.egg-info/

# =======================================================================
# .NET Ecosystem (C# / F# / VB)
# =======================================================================
[Bb]in/
[Oo]bj/

# =======================================================================
# Environment Secrets & Cryptographic Keys
# =======================================================================
.env
.env.local
.env.*.local
*.local
*.pem
*.key
`

const DefaultGitattributesBaseline = `# Set default behavior to automatically normalize line endings
* text=auto

# Explicitly declare text files you want normalized to LF on check-in
*.go text eol=lf
*.py text eol=lf
*.js text eol=lf
*.ts text eol=lf
*.md text eol=lf
*.json text eol=lf

# Explicitly declare dotnet C# source files
*.cs text eol=lf
*.sln text eol=lf
*.csproj text eol=lf

# Denote all files that are truly binary and should not be modified
*.png filter=lfs diff=lfs merge=lfs -text
*.jpg filter=lfs diff=lfs merge=lfs -text
*.gif filter=lfs diff=lfs merge=lfs -text
*.ico filter=lfs diff=lfs merge=lfs -text
`

const DefaultPrettierignoreBaseline = `node_modules
dist
.output
.vinxi
pnpm-lock.yaml
package-lock.json
bun.lock
routeTree.gen.ts
`

const commonUsage = `Usage: gitmap common <target> [flags]

Targets:
  ignore            Union-merge curated defaults into ./.gitignore
  attributes       Union-merge curated defaults into ./.gitattributes
  lfs-install      Run 'git lfs install --local' and merge lfs/common .gitattributes block
                   (delegates to 'gitmap add lfs-install'; supports --dry-run)
  prettier-ignore  Union-merge curated defaults into ./.prettierignore
  prettier-rc      Merge curated JSON defaults into ./.prettierrc (existing keys win)
  all              Run every target above in sequence

Flags:
  --dry-run, -n    Print planned additions without touching disk
  --force,   -f    Overwrite conflicting JSON values in .prettierrc

Behavior:
  - Line-based targets (ignore/attributes/prettier-ignore) append MISSING lines
    only; existing entries are preserved verbatim. Safe to re-run.
  - lfs-install delegates to the same code path as 'gitmap add lfs-install'
    (marker-block managed, idempotent). Requires git-lfs on PATH.
  - .prettierrc is JSON key-union: missing keys are added; existing keys stay
    unless --force is passed.
  - Alias: gitmap co <target>

Examples:
  gitmap common ignore
  gitmap common attributes --dry-run
  gitmap common lfs-install
  gitmap common lfs-install --dry-run
  gitmap common all --dry-run
  gitmap common prettier-rc --force
`

// DispatchCommon routes `gitmap common <target>` subcommands.
func DispatchCommon(command string) (bool, error) {
	if command != constants.CmdCommon && command != constants.CmdCommonAlias {
		return false, nil
	}

	if len(os.Args) < 3 {
		runAllCommonTargets(false, false)

		return true, nil
	}

	if isCommonHelp(os.Args[2]) {
		RenderSyncHelp()

		return true, nil
	}

	return executeCommonTarget(os.Args[2], os.Args[3:])
}

func runAllCommonTargets(dry, force bool) {
	_ = RunCommonLines(constants.FileGitignore, DefaultGitignoreBaseline, dry)
	_ = RunCommonLines(constants.FileGitattributes, DefaultGitattributesBaseline, dry)
	_ = RunCommonLFSInstall(dry)
	_ = RunCommonLines(constants.FilePrettierignore, DefaultPrettierignoreBaseline, dry)
	_ = RunCommonPrettierRC(dry, force)
}

func executeCommonTarget(sub string, rest []string) (bool, error) {
	if strings.HasPrefix(sub, "-") {
		rest = append([]string{sub}, rest...)
		sub = constants.TargetAll
	}
	dry, force := ParseCommonFlags(rest)

	return dispatchTargetAction(sub, dry, force)
}

func dispatchTargetAction(sub string, dry, force bool) (bool, error) {
	if sub == constants.TargetAll {
		runAllCommonTargets(dry, force)

		return true, nil
	}

	return dispatchSingleTarget(sub, dry, force)
}

func dispatchSingleTarget(sub string, dry, force bool) (bool, error) {
	switch sub {
	case constants.TargetIgnore:
		return true, RunCommonLines(constants.FileGitignore, DefaultGitignoreBaseline, dry)
	case constants.TargetAttributes:
		return true, RunCommonLines(constants.FileGitattributes, DefaultGitattributesBaseline, dry)
	case constants.TargetLFSInstall:
		return true, RunCommonLFSInstall(dry)
	case constants.TargetPrettierIgnore:
		return true, RunCommonLines(constants.FilePrettierignore, DefaultPrettierignoreBaseline, dry)
	case constants.TargetPrettierRC:
		return true, RunCommonPrettierRC(dry, force)
	default:
		return true, unknownTargetError(sub)
	}
}

func unknownTargetError(sub string) error {
	return apperror.NewWithDetails(
		"cmd.common.dispatch.unknown",
		constants.ErrCodeUnknownTarget,
		fmt.Sprintf("unknown sync target: %s\n\n%s", sub, commonUsage),
		"cmd.common",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		map[string]any{"subcommand": sub},
	)
}

func isCommonHelp(arg string) bool {
	return arg == constants.FlagHelp || arg == constants.FlagShortHelp || arg == constants.FlagLongHelp
}

func isSyncHelp(arg string) bool {
	return isCommonHelp(arg)
}

// ParseCommonFlags scans args for --dry-run and --force (position agnostic).
func ParseCommonFlags(args []string) (dry, force bool) {
	for _, a := range args {
		switch a {
		case constants.FlagLongDryRun, constants.FlagShortDryRun:
			dry = true
		case constants.FlagLongForce, constants.FlagShortForce:
			force = true
		}
	}

	return
}

// RunCommonLines appends every line from baseline that is not already
// present (verbatim, trimmed compare) in the target file.
func RunCommonLines(path, baseline string, dry bool) error {
	existing, _ := os.ReadFile(path)
	toAdd := computeMissingBaselineLines(existing, baseline)

	if len(toAdd) == 0 {
		fmt.Printf("  ok  %s already has all curated entries\n", path)

		return nil
	}

	if dry {
		printDryRunAdditions(path, toAdd)

		return nil
	}

	return appendLinesToFile(path, existing, toAdd)
}

func computeMissingBaselineLines(existing []byte, baseline string) []string {
	present := map[string]bool{}
	for _, l := range strings.Split(string(existing), "\n") {
		present[strings.TrimSpace(l)] = true
	}

	var toAdd []string
	hasExisting := len(existing) > 0
	for _, l := range strings.Split(baseline, "\n") {
		trimmed := strings.TrimSpace(l)
		if (trimmed == "" || strings.HasPrefix(trimmed, "#")) && hasExisting {
			continue
		}
		if present[trimmed] {
			continue
		}
		present[trimmed] = true
		toAdd = append(toAdd, l)
	}

	return toAdd
}

func printDryRunAdditions(path string, toAdd []string) {
	fmt.Printf("  +   %s would gain %d line(s):\n", path, len(toAdd))
	for _, l := range toAdd {
		fmt.Printf("      %s\n", l)
	}
}

func appendLinesToFile(path string, existing []byte, toAdd []string) error {
	buf := buildAppendedContent(existing, toAdd)
	if err := os.WriteFile(path, []byte(buf), 0o644); err != nil {
		return apperror.WrapWithDetails(
			err,
			"cmd.common.writeLines",
			constants.ErrCodeWriteLines,
			"failed to write synced lines to file",
			"cmd.common",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"path": path},
		)
	}

	fmt.Printf("  +   %s: added %d line(s)\n", path, len(toAdd))

	return nil
}

func buildAppendedContent(existing []byte, toAdd []string) string {
	buf := string(existing)
	if len(buf) > 0 && !strings.HasSuffix(buf, "\n") {
		buf += "\n"
	}
	if len(existing) > 0 {
		buf += constants.CommentAddedByCommon
	}
	buf += strings.Join(toAdd, "\n")
	if !strings.HasSuffix(buf, "\n") {
		buf += "\n"
	}

	return buf
}

// RunCommonLFSInstall delegates to `gitmap add lfs-install`.
func RunCommonLFSInstall(dry bool) error {
	var args []string
	if dry {
		args = append(args, constants.FlagLongDryRun)
	}

	runAddLFSInstall(args)

	return nil
}
