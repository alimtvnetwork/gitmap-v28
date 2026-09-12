package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func parsePurgeArgs(args []string) (pat string, isRestore, isAutoConfirm, isLovable bool) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--restore" {
			isRestore = true
		} else if a == "--lovable-untracked" {
			isLovable = true
		} else if a == "-y" || a == "--confirm" {
			isAutoConfirm = true
		} else if a == "--path" && i+1 < len(args) {
			i++
			pat = args[i]
		} else if !strings.HasPrefix(a, "-") && pat == "" {
			pat = a
		}
	}

	return pat, isRestore, isAutoConfirm, isLovable
}

func runPurgeCmd(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("%s failed: %w\n%s", name, err, string(out))
	}

	return string(out), nil
}

func copyPurgeFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}

func runPurge(args []string) error {
	pat, isRestore, isAutoConfirm, isLovable := parsePurgeArgs(args)
	repo, err := os.Getwd()
	if err != nil {
		return apperror.Wrap(err, "failed to get current directory", nil)
	}

	db, err := store.OpenDefault()
	if err != nil {
		return apperror.Wrap(err, "failed to open database", nil)
	}

	defer db.Close()
	if isLovable {
		return doPurgeLovable(repo)
	}

	if isRestore {
		return doRestore(db, repo)
	}

	if pat == "" {
		return apperror.NewSimple("EXECUTION", "pattern argument required")
	}

	return doPurge(db, repo, pat, isAutoConfirm)
}
