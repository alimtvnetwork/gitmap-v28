// Package cmd — inject.go implements `gitmap inject` (`inj`).
//
// Purpose: take an existing on-disk folder and "inject" it into the
// user's tooling — register with GitHub Desktop, open in VS Code, and
// (when a `git remote get-url origin` succeeds) upsert into the
// gitmap SQLite database so it shows up in `cd`, `list`, etc.
//
// Forms:
//
//	gitmap inject              # operate on cwd
//	gitmap inject <folder>     # operate on the given folder
//	gitmap inj   ...           # short alias
//
// `<folder>` accepts absolute, relative, or `~`-prefixed paths.
// Reuses `resolveCloneNextFolder` from clonenextfolderdispatch.go for
// path resolution + dir validation, so error messages stay consistent
// with `cn <folder>`.
//
// Per the user's spec answer: any folder is accepted (no `.git/`
// required) — Desktop will silently skip non-repos and VS Code is
// happy to open anything. The DB upsert is conditional: if the folder
// has no `origin` remote, we skip the database write but still do
// Desktop + VS Code, so local-only sandboxes can still be injected
// into the editor without polluting the repo index.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// runInject is the entrypoint for `gitmap inject` / `inj`.
func runInject(args []string) {
	checkHelp("inject", args)

	force := parseInjectForceFlag("inject", args)

	target, err := resolveInjectTarget(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrInjectResolve, err)
		os.Exit(1)
	}

	repoName := filepath.Base(target)
	fmt.Printf(constants.MsgInjectStart, repoName, target)

	if force {
		fmt.Printf(constants.MsgInjectForceNotice, repoName)
	}

	// 1. DB upsert — only when an origin remote is configured. A
	//    local-only repo (or a non-repo folder) silently skips this
	//    step; we still proceed to Desktop + VS Code below.
	upsertInjectIfRemote(target, repoName)

	// Load post-upsert stamps so a freshly-inserted row is correctly
	// seen as "never injected" (zero strings) on first run.
	stamps := loadInjectStamps(target)

	// 2. GitHub Desktop registration — gated by per-tool stamp.
	if shouldRunDesktop(target, stamps, force) {
		registerSingleDesktop(repoName, target)
		markInjected(target, constants.InjectKindDesktop)
	}

	// 3. VS Code open — same gate.
	if shouldRunVSCode(target, stamps, force) {
		openInVSCode(target)
		markInjected(target, constants.InjectKindVSCode)
	}

	// 4. Shell handoff so the parent shell cds into the injected
	//    folder (mirrors clone / cn / cd UX). Skipped silently when
	//    the wrapper isn't installed.
	WriteShellHandoff(target)

	fmt.Printf(constants.MsgInjectDone, repoName)
}

// resolveInjectTarget returns the absolute path of the folder to
// inject. With no positional args it returns cwd; with one positional
// it resolves through the same `cn <folder>` pipeline. Flags are
// ignored at the dispatcher level — `inject` accepts none today.
func resolveInjectTarget(args []string) (string, error) {
	positional := extractPositionalArgs(args)

	if len(positional) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		return cwd, nil
	}

	if len(positional) > 1 {
		return "", fmt.Errorf("expected 0 or 1 folder argument, got %d", len(positional))
	}

	resolved, err := resolveCloneNextFolder(positional[0])
	if err != nil {
		return "", fmt.Errorf("folder not found or not a directory: %s", positional[0])
	}

	return resolved, nil
}

// upsertInjectIfRemote writes the repo to the SQLite DB only when a
// remote origin is configured. Per the user's spec answer, missing
// remotes do NOT abort — we still proceed to Desktop + VS Code so
// local-only sandboxes can be injected into the editor.
//
// The print policy mirrors `upsertDirectClone`: warnings on stderr,
// no exit on failure. A best-effort persistence step shouldn't bring
// down the more useful Desktop/VS Code side effects.
func upsertInjectIfRemote(absPath, repoName string) {
	remoteURL, err := gitutil.RemoteURL(absPath)
	if err != nil || len(remoteURL) == 0 {
		fmt.Printf(constants.MsgInjectNoRemote, repoName)

		return
	}

	rec := model.ScanRecord{
		Slug:         strings.ToLower(repoName),
		RepoName:     repoName,
		RelativePath: repoName,
		AbsolutePath: absPath,
	}
	if strings.HasPrefix(remoteURL, constants.PrefixSSH) || strings.HasPrefix(remoteURL, "git@") {
		rec.SSHUrl = remoteURL
	} else {
		rec.HTTPSUrl = remoteURL
	}

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnInjectDBOpen, err)

		return
	}
	defer db.Close()

	if upsertErr := db.UpsertRepos([]model.ScanRecord{rec}); upsertErr != nil {
		fmt.Fprintf(os.Stderr, constants.WarnInjectDBUpsert, upsertErr)

		return
	}

	fmt.Printf(constants.MsgInjectDBOK, repoName, remoteURL)
}
