package cmdide

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdantigravity"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcursor"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

func runIDESync(args []string) error {
	opts, remaining := parseIDEOptions(args)
	repoPaths, err := resolveSyncRepoPaths(opts, remaining)
	if err != nil {
		return err
	}
	summary := SyncReposAcrossIDEsDirect(repoPaths, opts)
	printSyncSummary(summary, opts.IsJSON, opts.IsQuiet)
	return nil
}

func resolveSyncRepoPaths(opts IDEOptions, remaining []string) ([]string, error) {
	if opts.TargetDirectory != "" {
		return scanDirectoryRepos(opts.TargetDirectory), nil
	}
	if len(remaining) > 0 && isDirTarget(remaining[0]) {
		return scanDirectoryRepos(remaining[0]), nil
	}
	return fetchStoreRepos()
}

func isDirTarget(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func scanDirectoryRepos(root string) []string {
	var paths []string
	entries, _ := os.ReadDir(root)
	for _, e := range entries {
		sub := filepath.Join(root, e.Name())
		if e.IsDir() && isDirTarget(filepath.Join(sub, ".git")) {
			paths = append(paths, sub)
		}
	}
	return paths
}

func SyncReposAcrossIDEsDirect(repoPaths []string, opts IDEOptions) IDESyncSummary {
	s := IDESyncSummary{TotalRepos: len(repoPaths)}
	syncEditorTargets(repoPaths, opts, &s)
	syncOtherTargets(repoPaths, opts, &s)
	return s
}

func syncEditorTargets(repoPaths []string, opts IDEOptions, s *IDESyncSummary) {
	if opts.IsVSCodeTargeted && !opts.IsDryRun {
		pairs := make([]vscodepm.Pair, len(repoPaths))
		for i, p := range repoPaths {
			pairs[i] = vscodepm.Pair{RootPath: p, Name: filepath.Base(p), Tags: []string{"gitmap"}}
		}
		sum, _ := vscodepm.SyncMode(pairs, vscodepm.MergeModeUnion)
		s.VSCodeAdded, s.VSCodeRetained = sum.Added, sum.Total-sum.Added
	}
	if opts.IsCursorTargeted && !opts.IsDryRun {
		added, total, _ := cmdcursor.SyncCursorProjects(repoPaths)
		s.CursorAdded, s.CursorRetained = added, total-added
	}
}

func syncOtherTargets(repoPaths []string, opts IDEOptions, s *IDESyncSummary) {
	if opts.IsAntigravityTargeted {
		syncAntigravityPaths(repoPaths, opts.IsDryRun, s)
	}
	cli := desktop.ResolveCLI()
	if opts.IsDesktopTargeted && cli != "" && !opts.IsDryRun {
		syncDesktopPaths(repoPaths, cli, s)
	}
}

func syncAntigravityPaths(repoPaths []string, isDryRun bool, s *IDESyncSummary) {
	for _, p := range repoPaths {
		if cmdantigravity.IsRepoRegisteredInAgy(p) {
			s.AntigravityRetained++
		} else if !isDryRun && workspacesync.SyncAntigravity(p, filepath.Base(p)) {
			s.AntigravityAdded++
		}
	}
}

func syncDesktopPaths(repoPaths []string, cli string, s *IDESyncSummary) {
	for _, p := range repoPaths {
		if exec.Command(cli, p).Run() == nil {
			s.DesktopAdded++
		} else {
			s.DesktopFailed++
		}
	}
}

func printSyncSummary(s IDESyncSummary, isJSON, isQuiet bool) {
	if isJSON {
		b, _ := json.MarshalIndent(s, "", "  ")
		fmt.Println(string(b))
		return
	}
	if !isQuiet {
		fmt.Printf("%s✔ IDE Sync Completed:%s %d total repositories\n", constants.ColorGreen, constants.ColorReset, s.TotalRepos)
		fmt.Printf("  • VS Code:     %d added, %d retained\n", s.VSCodeAdded, s.VSCodeRetained)
		fmt.Printf("  • Cursor:      %d added, %d retained\n", s.CursorAdded, s.CursorRetained)
		fmt.Printf("  • Antigravity: %d added, %d retained\n", s.AntigravityAdded, s.AntigravityRetained)
		fmt.Printf("  • Desktop:     %d added, %d failed\n", s.DesktopAdded, s.DesktopFailed)
	}
}
