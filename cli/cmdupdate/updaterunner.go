package cmdupdate

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/release"
	"github.com/alimtvnetwork/gitmap-v28/cli/verbose"
)

// runUpdateRunner is a hidden command that performs the real update work.
// After the binary update completes, it triggers a best-effort schema
// migration so the next CLI invocation never has to repair the database.
//
// Phase 3 of the handoff chain runs at the end via
// scheduleDeployedCleanupHandoff: the freshly-deployed gitmap binary
// (a different file than this handoff copy) is invoked detached to
// run `update-cleanup`. Only the deployed binary can safely remove
// the still-locked handoff copy and the just-renamed *.exe.old.
// See 02-spec/08-generic-update/06-cleanup.md and
// 02-spec/03-general/02f-self-update-orchestration.md.
func runUpdateRunner() error {
	repoPath, err := resolveRepoPath()
	if err != nil {
		return err
	}

	report := resolveReportErrors()

	currentVersion := constants.Version
	targetVersion := readTargetVersion(repoPath)

	initRunnerVerbose()
	fmt.Printf(constants.MsgUpdateStarting)
	fmt.Printf(constants.MsgUpdateRepoPath, repoPath)
	fmt.Printf(constants.MsgUpdateVersionCompare, currentVersion, targetVersion)
	executeUpdate(repoPath, report)
	RunPostUpdateMigrate()
	report.summarize()
	printUpdateSummary(currentVersion, targetVersion, repoPath)
	scheduleDeployedCleanupHandoff()

	return nil
}

// getFlagValue returns the value following a flag like --repo-path <value>.
func getFlagValue(name string) string {
	args := os.Args[2:]
	for i, arg := range args {
		if arg == name && i+1 < len(args) {
			return args[i+1]
		}
	}

	return ""
}

// initRunnerVerbose initializes verbose logging if --verbose flag is present.
func initRunnerVerbose() {
	if !hasFlag(constants.FlagVerbose) {
		return
	}

	log, err := verbose.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)

		return
	}

	defer log.Close()
	log.Log(constants.UpdateRunnerLogStart, constants.RepoPath)
}

// hasFlag checks if a flag is present in os.Args[2:].
func hasFlag(name string) bool {
	for _, arg := range os.Args[2:] {
		if arg == name {
			return true
		}
	}

	return false
}

// copyFile copies src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, in)

	return err
}

// readTargetVersion reads the target version from the repository.
func readTargetVersion(repoPath string) string {
	versionPath := filepath.Join(repoPath, constants.DefaultVersionFile)
	data, err := os.ReadFile(versionPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateVersionRead, versionPath, err)

		return "unknown"
	}

	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateVersionRead, versionPath, err)

		return "unknown"
	}

	if v, ok := rawMap["Version"].(string); ok && len(v) > 0 {
		return v
	}

	if v, ok := rawMap["version"].(string); ok && len(v) > 0 {
		return v
	}

	return "unknown"
}

// printUpdateSummary outputs the final update summary.
func printUpdateSummary(oldVer, newVer, repoPath string) {
	fmt.Printf(constants.MsgUpdateSummaryDetail, oldVer, newVer, repoPath)
	fmt.Printf("  → Update: %s -> %s\n", oldVer, newVer)
	printUpdateChangelog()
	fmt.Printf("\n\n")
}

// printUpdateChangelog prints the last two notes from the changelog.
func printUpdateChangelog() {
	entries, err := release.ReadChangelog()
	hasEntries := err == nil && len(entries) > 0
	if !hasEntries {
		return
	}

	for _, note := range getLastTwoNotes(entries[0].Notes) {
		fmt.Printf("    - %s\n", note)
	}
}

func getLastTwoNotes(notes []string) []string {
	count := len(notes)
	if count > 2 {
		return notes[count-2:]
	}

	return notes
}
