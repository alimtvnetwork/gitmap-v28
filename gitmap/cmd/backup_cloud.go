package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runBackupCloudPush(args []string) error {
	prof, err := resolveDefaultCloudProfile(args)
	if err != nil {
		return err
	}
	repoSlug := resolveBackupRepoSlug(prof, args)
	cloudDir := filepath.Join(store.BinaryDataDir(), "cloud-backup")
	if prepErr := ensureCloudRepoPrepared(cloudDir, repoSlug); prepErr != nil {
		return prepErr
	}
	snapID, copyErr := createSnapshotFolder(cloudDir, args)
	if copyErr != nil {
		return copyErr
	}
	if pushErr := commitAndPushSnapshot(cloudDir, snapID); pushErr != nil {
		return pushErr
	}
	printBackupSuccess(snapID, repoSlug)
	return nil
}

func resolveDefaultCloudProfile(args []string) (model.GitProfile, error) {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return model.GitProfile{}, apperror.WrapSimple(err, "load profiles:")
	}
	req := extractFlagVal(args, "--profile")
	if req != "" {
		_, p, findErr := pickProfileBySequenceOrName(cfg.Profiles, req)
		return p, findErr
	}
	for _, p := range cfg.Profiles {
		if p.IsDefault || p.Name == cfg.Default {
			return p, nil
		}
	}
	if len(cfg.Profiles) > 0 {
		return cfg.Profiles[0], nil
	}
	return model.GitProfile{Name: "default", Provider: "github"}, nil
}

func resolveBackupRepoSlug(p model.GitProfile, args []string) string {
	customRepo := extractFlagVal(args, "--repo")
	if customRepo != "" {
		return customRepo
	}
	repoName := "gitmap-cloud-backup"
	if p.Name != "" && p.Name != "default" {
		return p.Name + "/" + repoName
	}
	return repoName
}

func ensureCloudRepoPrepared(cloudDir, repoSlug string) error {
	if mkErr := os.MkdirAll(cloudDir, 0755); mkErr != nil {
		return apperror.WrapSimple(mkErr, "create cloud backup dir:")
	}
	ensureRemoteRepoExists(repoSlug)
	gitDir := filepath.Join(cloudDir, ".git")
	if isDirExists(gitDir) {
		cmdPull := exec.Command("git", "pull", "--rebase", "origin", "main")
		cmdPull.Dir = cloudDir
		_ = cmdPull.Run()
		return nil
	}
	cmdClone := exec.Command("gh", "repo", "clone", repoSlug, cloudDir)
	out, err := cmdClone.CombinedOutput()
	if err != nil {
		return initFallbackRepo(cloudDir, repoSlug, string(out))
	}
	return nil
}

func isDirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func ensureRemoteRepoExists(repoSlug string) {
	cmdView := exec.Command("gh", "repo", "view", repoSlug)
	if viewErr := cmdView.Run(); viewErr == nil {
		return
	}
	cmdCreate := exec.Command("gh", "repo", "create", repoSlug, "--private")
	_ = cmdCreate.Run()
}

func initFallbackRepo(cloudDir, repoSlug, logMsg string) error {
	cmdInit := exec.Command("git", "init", "-b", "main")
	cmdInit.Dir = cloudDir
	if initErr := cmdInit.Run(); initErr != nil {
		return apperror.WrapSimple(initErr, "fallback git init: "+logMsg)
	}
	remoteURL := fmt.Sprintf("https://github.com/%s.git", repoSlug)
	cmdRemote := exec.Command("git", "remote", "add", "origin", remoteURL)
	cmdRemote.Dir = cloudDir
	_ = cmdRemote.Run()
	return nil
}

func createSnapshotFolder(cloudDir string, args []string) (string, error) {
	snapID := "snapshot-" + time.Now().Format("2006-01-02-150405")
	destDir := filepath.Join(cloudDir, "snapshots", snapID)
	if mkErr := os.MkdirAll(destDir, 0755); mkErr != nil {
		return "", apperror.WrapSimple(mkErr, "create snapshot dir:")
	}
	copyBackupArtifacts(destDir)
	note := extractFlagVal(args, "--note")
	writeSnapshotManifest(destDir, snapID, note)
	return snapID, nil
}

func commitAndPushSnapshot(cloudDir, snapID string) error {
	cmdAdd := exec.Command("git", "add", ".")
	cmdAdd.Dir = cloudDir
	_ = cmdAdd.Run()

	commitMsg := fmt.Sprintf("backup: capture %s", snapID)
	cmdCommit := exec.Command("git", "commit", "-m", commitMsg)
	cmdCommit.Dir = cloudDir
	_ = cmdCommit.Run()

	cmdPush := exec.Command("git", "push", "-u", "origin", "main")
	cmdPush.Dir = cloudDir
	out, pushErr := cmdPush.CombinedOutput()
	if pushErr != nil {
		return apperror.NewSimple("git push backup failed: "+string(out), "E1079")
	}
	return nil
}

func printBackupSuccess(snapID, repoSlug string) {
	fmt.Printf("\n  %s✓ Cloud backup snapshot created and pushed successfully!%s\n",
		constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  ● Snapshot ID: %s\n", snapID)
	fmt.Printf("  ● Remote Repo: https://github.com/%s\n", repoSlug)
	fmt.Printf("  ● Timestamp:   %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
}
