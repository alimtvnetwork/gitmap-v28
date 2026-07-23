// Package cmd — vscodeworkspace.go: implements `gitmap vscode-workspace`
// (alias `vsws`).
//
// Emits a single multi-root `.code-workspace` file from the same Repo
// table that drives the Project Manager sync, so one click in VS Code
// (File → Open Workspace from File…) opens every cloned repo as a
// folder in one window — no manual "Add Folder to Workspace…" loop.
//
// Distinct from the projects.json sync: the PM sync gives you a flat
// sidebar list of projects (each opens in its own VS Code window);
// the workspace file gives you ONE window with N folders, ideal for
// cross-repo search / refactor. Both surfaces stay in lockstep
// because both read the same DB.
//
// Spec: spec/01-vscode-project-manager-sync/03-workspace-export.md
package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/vscodepm"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/vscodeworkspace"
)

// vscodeWorkspaceFlags carries the parsed CLI flags for one run.
type vscodeWorkspaceFlags struct {
	out        string
	isRelative bool
	tag        string
	rootSubdir string
}

// runVSCodeWorkspace is the dispatcher entry point.
func runVSCodeWorkspace(args []string) {
	checkHelp("vscode-workspace", args)

	flags := parseVSCodeWorkspaceFlags(args)
	records, err := loadReposForWorkspace()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	folders := buildFoldersFromRecords(records, flags.tag, flags.rootSubdir)
	if len(folders) == 0 {
		fmt.Print(constants.MsgVSCodeWorkspaceEmpty)

		return
	}

	writeWorkspaceFile(flags, folders)
}

// parseVSCodeWorkspaceFlags parses the three supported flags.
// Unknown flags trigger the standard flag.PrintDefaults + exit(2).
func parseVSCodeWorkspaceFlags(args []string) vscodeWorkspaceFlags {
	fs := flag.NewFlagSet(constants.CmdVSCodeWorkspace, flag.ExitOnError)
	cfg := vscodeWorkspaceFlags{out: constants.VSCodeWorkspaceDefaultFilename}

	fs.StringVar(&cfg.out, constants.FlagVSCodeWorkspaceOut, cfg.out,
		constants.FlagDescVSCodeWorkspaceOut)
	fs.BoolVar(&cfg.isRelative, constants.FlagVSCodeWorkspaceRelative, false,
		constants.FlagDescVSCodeWorkspaceRelative)
	fs.StringVar(&cfg.tag, constants.FlagVSCodeWorkspaceTag, "",
		constants.FlagDescVSCodeWorkspaceTag)
	fs.StringVar(&cfg.rootSubdir, constants.FlagVSCodeWorkspaceRootSubdir, "",
		constants.FlagDescVSCodeWorkspaceRootSubdir)

	_ = fs.Parse(args)

	return cfg
}

// loadReposForWorkspace reads every tracked repo from the DB.
// Same code path the PM sync ultimately uses (store.ListRepos), so
// the two surfaces never drift.
func loadReposForWorkspace() ([]model.ScanRecord, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrVSCodeWorkspaceDBOpen, err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		return nil, fmt.Errorf(constants.ErrVSCodeWorkspaceDBOpen, err)
	}

	records, err := db.ListRepos()
	if err != nil {
		return nil, fmt.Errorf(constants.ErrVSCodeWorkspaceDBList, err)
	}

	return records, nil
}

// buildFoldersFromRecords converts ScanRecords into Folder tuples,
// applying the optional --tag filter via vscodepm.DetectTags so the
// filter semantics match what the PM sync writes. When rootSubdir is
// non-empty, each folder's path is the repo root joined with that
// subdir; repos missing that subdir are skipped with a notice.
func buildFoldersFromRecords(records []model.ScanRecord, tag, rootSubdir string) []vscodeworkspace.Folder {
	out := make([]vscodeworkspace.Folder, 0, len(records))
	for _, r := range records {
		if !matchesWorkspaceTag(r.AbsolutePath, tag) {
			continue
		}
		path, ok := resolveFolderPath(r.AbsolutePath, rootSubdir)
		if !ok {
			fmt.Fprintf(os.Stderr, constants.MsgVSCodeWorkspaceSubdirSkip, r.RepoName, rootSubdir)
			continue
		}
		out = append(out, vscodeworkspace.Folder{Name: r.RepoName, Path: path})
	}

	return out
}

// resolveFolderPath returns the workspace folder path for one repo.
// Empty rootSubdir => repo root unchanged. Non-empty rootSubdir =>
// joined path, only if it exists on disk as a directory.
func resolveFolderPath(repoRoot, rootSubdir string) (string, bool) {
	if rootSubdir == "" {
		return repoRoot, true
	}
	candidate := filepath.Join(repoRoot, rootSubdir)
	info, err := os.Stat(candidate)
	if err != nil || !info.IsDir() {
		return "", false
	}
	return candidate, true
}

// matchesWorkspaceTag returns true when the tag filter is empty OR
// the rootPath's auto-detected tag set contains the requested tag.
func matchesWorkspaceTag(rootPath, tag string) bool {
	if tag == "" {
		return true
	}
	for _, t := range vscodepm.DetectTagsCustom(rootPath) {
		if strings.EqualFold(t, tag) {
			return true
		}
	}

	return false
}

// writeWorkspaceFile assembles the Workspace, optionally relativizes
// paths, then commits the file atomically.
func writeWorkspaceFile(flags vscodeWorkspaceFlags, folders []vscodeworkspace.Folder) {
	outPath, err := filepath.Abs(flags.out)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	finalFolders := folders
	if flags.isRelative {
		finalFolders, err = vscodeworkspace.Relativize(folders, filepath.Dir(outPath))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	ws := vscodeworkspace.Build(finalFolders)
	if err := vscodeworkspace.WriteAtomic(outPath, ws); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf(constants.MsgVSCodeWorkspaceWritten, outPath, len(ws.Folders))
}
