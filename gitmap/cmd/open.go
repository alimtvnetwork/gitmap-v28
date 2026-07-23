// Package cmd — open.go implements `gitmap open` (alias `op`).
//
// Detects the repo for the current working directory (preferring
// `git rev-parse --show-toplevel` when available, falling back to
// cwd) and launches BOTH GitHub Desktop and VS Code on that path.
//
// Behavior is intentionally idempotent-by-side-effect: GitHub
// Desktop silently skips if the repo is already registered, and
// VS Code happily re-opens an already-open window. The DB upsert
// step mirrors `inject`: best-effort, only when a remote origin
// exists, never aborts the user-visible side effects.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runOpen is the entrypoint for `gitmap open` / `op`.
func runOpen(args []string) {
	checkHelp(constants.CmdOpen, args)

	force := parseInjectForceFlag(constants.CmdOpen, args)

	target, err := resolveOpenTarget()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrOpenResolveCwd, err)
		os.Exit(1)
	}

	repoName := filepath.Base(target)
	fmt.Printf(constants.MsgOpenStart, repoName, target)

	if force {
		fmt.Printf(constants.MsgInjectForceNotice, repoName)
	}

	// Best-effort DB upsert (skipped silently if no origin remote).
	upsertInjectIfRemote(target, repoName)

	stamps := loadInjectStamps(target)

	// Always re-inject when forced; otherwise gate on per-tool stamp.
	if shouldRunDesktop(target, stamps, force) {
		registerSingleDesktop(repoName, target)
		markInjected(target, constants.InjectKindDesktop)
	}

	if shouldRunVSCode(target, stamps, force) {
		openInVSCode(target)
		markInjected(target, constants.InjectKindVSCode)
	}

	fmt.Printf(constants.MsgOpenDone, repoName)
}

// resolveOpenTarget picks the directory to open. Prefers the git
// toplevel (so running `open` from a subfolder still opens the repo
// root), and falls back to plain cwd when git isn't available or
// the folder isn't a repo.
func resolveOpenTarget() (string, error) {
	if root, err := gitTopLevel(); err == nil && len(root) > 0 {
		return root, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(cwd)
	if err != nil {
		return "", err
	}

	return abs, nil
}
