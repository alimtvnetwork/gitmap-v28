package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// Curated defaults sourced from the user-approved baseline (see helptext/sync.md).
// Kept inline (not embed) because these are short, hand-curated and rarely change.
const defaultGitignoreBaseline = `# =======================================================================
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

const defaultGitattributesBaseline = `# Set default behavior to automatically normalize line endings
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

const defaultPrettierignoreBaseline = `node_modules
dist
.output
.vinxi
pnpm-lock.yaml
package-lock.json
bun.lock
routeTree.gen.ts
`

var defaultPrettierrcBaseline = map[string]any{
	"printWidth":    100.0,
	"semi":          true,
	"singleQuote":   false,
	"trailingComma": "all",
}

const syncUsage = `Usage: gitmap sync <target> [flags]

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
  - Alias: gitmap sy <target>

Examples:
  gitmap sync ignore
  gitmap sync attributes --dry-run
  gitmap sync lfs-install
  gitmap sync lfs-install --dry-run
  gitmap sync all --dry-run
  gitmap sync prettier-rc --force
`

// dispatchSync routes `gitmap sync <target>` subcommands.
func dispatchSync(command string) bool {
	if command != constants.CmdSync && command != constants.CmdSyncAlias {
		return false
	}
	if len(os.Args) < 3 {
		fmt.Fprint(os.Stderr, syncUsage)
		os.Exit(1)
	}
	sub, rest := os.Args[2], os.Args[3:]
	dry, force := parseSyncFlags(rest)

	switch sub {
	case "ignore":
		runSyncLines(".gitignore", defaultGitignoreBaseline, dry)
	case "attributes":
		runSyncLines(".gitattributes", defaultGitattributesBaseline, dry)
	case "lfs-install":
		runSyncLFSInstall(dry)
	case "prettier-ignore":
		runSyncLines(".prettierignore", defaultPrettierignoreBaseline, dry)
	case "prettier-rc":
		runSyncPrettierRC(dry, force)
	case "all":
		runSyncLines(".gitignore", defaultGitignoreBaseline, dry)
		runSyncLines(".gitattributes", defaultGitattributesBaseline, dry)
		runSyncLFSInstall(dry)
		runSyncLines(".prettierignore", defaultPrettierignoreBaseline, dry)
		runSyncPrettierRC(dry, force)
	default:
		fmt.Fprintf(os.Stderr, "unknown sync target: %s\n\n", sub)
		fmt.Fprint(os.Stderr, syncUsage)
		os.Exit(1)
	}
	return true
}

// parseSyncFlags scans args for --dry-run and --force (position agnostic).
func parseSyncFlags(args []string) (dry, force bool) {
	for _, a := range args {
		switch a {
		case "--dry-run", "-n":
			dry = true
		case "--force", "-f":
			force = true
		}
	}
	return
}

// runSyncLines appends every line from baseline that is not already
// present (verbatim, trimmed compare) in the target file.
func runSyncLines(path, baseline string, dry bool) {
	existing, _ := os.ReadFile(path)
	present := map[string]bool{}
	for _, l := range strings.Split(string(existing), "\n") {
		present[strings.TrimSpace(l)] = true
	}

	var toAdd []string
	for _, l := range strings.Split(baseline, "\n") {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			// Skip comment/blank noise from baseline when file already exists.
			if len(existing) > 0 {
				continue
			}
		}
		if present[trimmed] {
			continue
		}
		present[trimmed] = true
		toAdd = append(toAdd, l)
	}

	if len(toAdd) == 0 {
		fmt.Printf("  ok  %s already has all curated entries\n", path)
		return
	}

	if dry {
		fmt.Printf("  +   %s would gain %d line(s):\n", path, len(toAdd))
		for _, l := range toAdd {
			fmt.Printf("      %s\n", l)
		}
		return
	}

	buf := string(existing)
	if len(buf) > 0 && !strings.HasSuffix(buf, "\n") {
		buf += "\n"
	}
	if len(existing) > 0 {
		buf += "\n# added by gitmap sync\n"
	}
	buf += strings.Join(toAdd, "\n")
	if !strings.HasSuffix(buf, "\n") {
		buf += "\n"
	}
	if err := os.WriteFile(path, []byte(buf), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "  x   %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Printf("  +   %s: added %d line(s)\n", path, len(toAdd))
}

// runSyncPrettierRC does a JSON key-union. Existing keys are kept unless
// --force is passed; missing keys are inserted from the baseline.
func runSyncPrettierRC(dry, force bool) {
	const path = ".prettierrc"
	current := map[string]any{}

	if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &current); err != nil {
			fmt.Fprintf(os.Stderr, "  x   %s: not valid JSON (%v). Fix or delete first.\n", path, err)
			os.Exit(1)
		}
	}

	var added, overwritten []string
	for k, v := range defaultPrettierrcBaseline {
		existing, has := current[k]
		if !has {
			current[k] = v
			added = append(added, k)
			continue
		}
		if force && !syncJSONEqual(existing, v) {
			current[k] = v
			overwritten = append(overwritten, k)
		}
	}

	if len(added) == 0 && len(overwritten) == 0 {
		fmt.Printf("  ok  %s already has all curated keys\n", path)
		return
	}

	sort.Strings(added)
	sort.Strings(overwritten)

	if dry {
		if len(added) > 0 {
			fmt.Printf("  +   %s would add: %s\n", path, strings.Join(added, ", "))
		}
		if len(overwritten) > 0 {
			fmt.Printf("  ~   %s would overwrite (--force): %s\n", path, strings.Join(overwritten, ", "))
		}
		return
	}

	out, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "  x   marshal %s: %v\n", path, err)
		os.Exit(1)
	}
	out = append(out, '\n')
	if err := os.WriteFile(path, out, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "  x   %s: %v\n", path, err)
		os.Exit(1)
	}
	if len(added) > 0 {
		fmt.Printf("  +   %s: added keys %s\n", path, strings.Join(added, ", "))
	}
	if len(overwritten) > 0 {
		fmt.Printf("  ~   %s: overwrote keys %s\n", path, strings.Join(overwritten, ", "))
	}
}

// syncJSONEqual compares two JSON-decoded values by re-marshaling. Small
// payloads only, fine for .prettierrc scalars/arrays. Named uniquely to
// avoid colliding with the identically-shaped helper in chromeprofile_merge.go.
func syncJSONEqual(a, b any) bool {
	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)
	return string(ab) == string(bb)
}

// runSyncLFSInstall delegates to `gitmap add lfs-install`. Kept as a
// thin wrapper so the sync surface stays a single dispatch table while
// the LFS logic (marker block, templates.Merge, git-lfs probe) lives
// once in addlfsinstall.go.
func runSyncLFSInstall(dry bool) {
	args := []string{}
	if dry {
		args = append(args, "--dry-run")
	}
	runAddLFSInstall(args)
}
