package cmdprompts

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"path/filepath"
	"sort"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const promptsDriftHelp = `prompts drift — canonical vs copy drift report for muse prompts.

Description:
  Compares the canonical prompt home
  (coding-guidelines-v24 01-prompts/27-muse-prompts/) against every
  registered local copy by content hash. Read-only: reports drift,
  never syncs.

Usage:
  gitmap prompts drift [--canonical <dir>] [--copies <dir,dir>] [--workdir <dir>]

Flags:
  --canonical <dir>   canonical prompts dir
                      (default: ~/workspace/repos/coding-guidelines-v24/01-prompts/27-muse-prompts)
  --copies <dirs>     comma-separated extra copy dirs to compare
  --workdir <dir>     auto-discover copies under this dir (default: ~/workspace/repos)
  --help, -h          show this help

Example:
  gitmap prompts drift
`

// promptCopy is one compared copy location.
type promptCopy struct {
	label string
	dir   string
}

// RunPromptsDrift implements `gitmap prompts drift`.
func RunPromptsDrift(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(promptsDriftHelp)
			return nil
		}
	}
	home, _ := os.UserHomeDir()
	defaultCanonical := filepath.Join(home, "workspace", "repos", "coding-guidelines-v24", "01-prompts", "27-muse-prompts")
	fs := flag.NewFlagSet("prompts drift", flag.ContinueOnError)
	canonicalFlag := fs.String("canonical", defaultCanonical, "canonical prompts dir")
	copiesFlag := fs.String("copies", "", "comma-separated extra copy dirs")
	workdirFlag := fs.String("workdir", filepath.Join(home, "workspace", "repos"), "auto-discovery root")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(promptsDriftHelp)
		return apperror.WrapSimple(parseErr, "prompts.drift.flags")
	}
	canonical := *canonicalFlag

	canonicalFiles, listErr := listPromptFiles(canonical)
	if listErr != nil {
		fmt.Printf("prompts drift: canonical prompts not found at %s\n", canonical)
		fmt.Println("hint: pass --canonical <dir> to point at the canonical prompt home.")
		return nil
	}
	if len(canonicalFiles) == 0 {
		fmt.Printf("prompts drift: no prompt files in %s\n", canonical)
		return nil
	}

	copies := discoverPromptCopies(canonical, *workdirFlag, *copiesFlag)
	if len(copies) == 0 {
		fmt.Println("prompts drift: no local copies found to compare against.")
		return nil
	}

	driftCount := 0
	for _, name := range canonicalFiles {
		canonicalHash, canonicalLines := hashPromptFile(filepath.Join(canonical, name))
		fmt.Printf("%s:\n", name)
		fmt.Printf("  canonical %-8s %s\n", canonicalHash, shortPath(canonical))
		for _, c := range copies {
			copyPath := filepath.Join(c.dir, name)
			copyHash, copyLines := hashPromptFile(copyPath)
			switch {
			case copyHash == "":
				fmt.Printf("  %-10s MISSING   %s\n", c.label, shortPath(c.dir))
				driftCount++
			case copyHash == canonicalHash:
				fmt.Printf("  %-10s ok        %s\n", c.label, copyHash)
			default:
				lineDelta := copyLines - canonicalLines
				fmt.Printf("  %-10s DRIFT     %s (%+d lines vs canonical)\n", c.label, copyHash, lineDelta)
				driftCount++
			}
		}
	}
	fmt.Println()
	if driftCount == 0 {
		fmt.Println("prompts drift: all copies match canonical.")
	} else {
		fmt.Printf("prompts drift: %d drifted/missing file(s) across %d copie(s).\n", driftCount, len(copies))
	}

	return nil
}

// listPromptFiles returns sorted .md filenames in dir.
func listPromptFiles(dir string) ([]string, error) {
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return nil, readErr
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)

	return names, nil
}

// hashPromptFile returns the short sha256 and line count, or "" when missing.
func hashPromptFile(path string) (string, int) {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return "", 0
	}
	sum := sha256.Sum256(data)
	lines := strings.Count(string(data), "\n")

	return hex.EncodeToString(sum[:])[:8], lines
}

// discoverPromptCopies finds <root>/*/01-prompts/27-muse-prompts plus --copies extras.
func discoverPromptCopies(canonical, workdir, copiesFlag string) []promptCopy {
	var copies []promptCopy
	seen := map[string]bool{}
	addCopy := func(label, dir string) {
		clean := filepath.Clean(dir)
		if seen[clean] || clean == filepath.Clean(canonical) {
			return
		}
		if _, statErr := os.Stat(clean); statErr != nil {
			return
		}
		seen[clean] = true
		copies = append(copies, promptCopy{label: label, dir: clean})
	}
	entries, readErr := os.ReadDir(workdir)
	if readErr == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			candidate := filepath.Join(workdir, e.Name(), "01-prompts", "27-muse-prompts")
			if _, statErr := os.Stat(candidate); statErr == nil {
				addCopy(e.Name(), candidate)
			}
		}
	}
	for _, raw := range strings.Split(copiesFlag, ",") {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			continue
		}
		addCopy(filepath.Base(trimmed), trimmed)
	}
	sort.Slice(copies, func(i, j int) bool { return copies[i].label < copies[j].label })

	return copies
}

// shortPath trims the home prefix for compact display.
func shortPath(path string) string {
	home, _ := os.UserHomeDir()
	if home != "" && strings.HasPrefix(path, home) {
		return "~" + strings.TrimPrefix(path, home)
	}

	return path
}
