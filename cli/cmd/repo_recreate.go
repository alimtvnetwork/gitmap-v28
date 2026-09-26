package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

// runRecreateRepo implements gitmap recreate-repo <name>.
func runRecreateRepo(args []string) error {
	checkHelp("recreate-repo", args)
	cleanArgs := extractNonFlagTokens(args)
	if len(cleanArgs) == 0 {
		return apperror.NewValidationError("usage: gitmap recreate-repo (recreate) <new-repo-name> [--public] [--private]")
	}

	rawName := cleanArgs[0]
	slug := SlugifyRepoName(rawName)
	isPublic := hasArgFlag(args, "--public") && !hasArgFlag(args, "--private")

	// 1. Verify inside git repo
	cmdCheck := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmdCheck.Run(); err != nil {
		return apperror.NewValidationError("current directory is not a git repository")
	}

	// 2. Create backup branch
	timestamp := time.Now().Format("20060102-150405")
	backupBranch := fmt.Sprintf("backup/recreate-%s", timestamp)
	cmdBackup := exec.Command("git", "branch", backupBranch)
	if err := cmdBackup.Run(); err != nil {
		return apperror.WrapSimple(err, "create backup branch before recreate:")
	}
	fmt.Printf("  %s✓ Created backup branch: %s%s\n", constants.ColorGreen, backupBranch, constants.ColorReset)

	// 3. Handle remote origin
	visibilityFlag := "--private"
	if isPublic {
		visibilityFlag = "--public"
	}

	// Probe current origin
	cmdOrigin := exec.Command("git", "remote", "get-url", "origin")
	if out, err := cmdOrigin.Output(); err == nil && len(strings.TrimSpace(string(out))) > 0 {
		// Rename existing origin to origin-backup
		_ = exec.Command("git", "remote", "rename", "origin", "origin-backup").Run()
	}

	// 4. Create new repository on GitHub via gh CLI
	cmdCreate := exec.Command("gh", "repo", "create", slug, visibilityFlag, "--source=.", "--remote=origin", "--push")
	outCreate, errCreate := cmdCreate.CombinedOutput()
	if errCreate != nil {
		// Restore origin if failed
		_ = exec.Command("git", "remote", "rename", "origin-backup", "origin").Run()

		return apperror.WrapSimple(errCreate, fmt.Sprintf("gh repo create failed: %s", string(outCreate)))
	}

	absDir, _ := filepath.Abs(".")
	workspacesync.SyncAll(absDir, slug)

	fmt.Printf("\n%s● Successfully recreated repository on GitHub%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  • Local Repo:    %s\n", absDir)
	fmt.Printf("  • Remote Repo:   https://github.com/%s\n", slug)
	fmt.Printf("  • Backup Branch: %s\n\n", backupBranch)

	return nil
}
