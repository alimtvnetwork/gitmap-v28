package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/mapper"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/scanner"
)

// runAs implements `gitmap as [alias-name] [--force]`.
//
// It must be invoked from inside a Git repository. It:
//  1. Resolves the repo top-level via `git rev-parse --show-toplevel`.
//  2. Builds a ScanRecord for the repo and upserts it into SQLite.
//  3. Creates (or updates with --force) an alias mapping name -> repo.
//
// When alias-name is omitted, the repo folder basename is used.
func runAs(args []string) *apperror.AppError {
	checkHelp(constants.CmdAs, args)
	aliasName, force := parseAsArgs(args)

	root, err := gitTopLevel()
	if err != nil {
		cwd, _ := os.Getwd()
		fmt.Fprintf(os.Stderr, constants.ErrAsNotInRepoFmt, cwd)
		return apperror.New("fatal error", "E9000", nil)
	}

	if aliasName == "" {
		aliasName = filepath.Base(root)
	}

	rec, buildErr := buildSingleRepoRecord(root); if buildErr != nil { return buildErr }
	if upsertErr := upsertSingleRepo(rec); upsertErr != nil { return upsertErr }
	if regErr := registerAlias(aliasName, rec, force); regErr != nil { return regErr }

	// Shell handoff: cd the parent shell to the alias root if invoked
	// via the wrapper function (e.g. `gitmap as foo` from elsewhere).
	WriteShellHandoff(root)
	return nil
}

// parseAsArgs extracts the optional alias-name positional and --force flag.
func parseAsArgs(args []string) (string, bool) {
	fs := flag.NewFlagSet(constants.CmdAs, flag.ExitOnError)
	force := fs.Bool(constants.FlagAsForce, false, "overwrite an existing alias")
	fs.BoolVar(force, constants.FlagAsForceS, false, "overwrite an existing alias (short)")

	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		os.Exit(2)
	}

	rest := fs.Args()
	if len(rest) > 1 {
		fmt.Fprintln(os.Stderr, constants.ErrAsUsage)
		os.Exit(2)
	}

	if len(rest) == 1 {
		return rest[0], *force
	}

	return "", *force
}

// gitTopLevel returns the absolute path of the current repo's top-level dir.
func gitTopLevel() (string, error) {
	cmd := exec.Command(constants.GitBin, constants.GitRevParse, "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("empty top-level")
	}

	return filepath.Clean(root), nil
}

// buildSingleRepoRecord constructs a ScanRecord for one already-known repo.
func buildSingleRepoRecord(absPath string) (model.ScanRecord, *apperror.AppError) {
	repos := []scanner.RepoInfo{{
		AbsolutePath: absPath,
		RelativePath: filepath.Base(absPath),
	}}
	records := mapper.BuildRecords(repos, constants.ModeHTTPS, "")
	if len(records) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrAsResolveFmt, absPath, "no record built")
		return model.ScanRecord{}, apperror.New("fatal error", "E9000", nil)
	}

	return records[0], nil
}
