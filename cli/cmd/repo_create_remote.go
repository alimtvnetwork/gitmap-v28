package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func buildRemoteSlug(p createRepoParams) string {
	slug := p.Slug
	if slug == "" {
		slug = SlugifyRepoName(p.Name)
	}

	hasProfile := p.Profile.Name != "" && p.Profile.Name != "default"
	if hasProfile {
		return p.Profile.Name + "/" + slug
	}

	return slug
}

func pushRemoteRepo(p createRepoParams) (string, error) {
	if p.IsSkipRemote {
		return "", nil
	}

	absDir, _ := filepath.Abs(p.LocalDir)
	visibilityFlag := "--private"
	if p.IsPublic {
		visibilityFlag = "--public"
	}

	slug := buildRemoteSlug(p)
	cmd := exec.Command("gh", "repo", "create", slug, visibilityFlag, "--source=.", "--remote=origin", "--push")
	cmd.Dir = absDir
	out, err := cmd.CombinedOutput()
	if err != nil && strings.Contains(strings.ToLower(string(out)), "already exists") {
		return handleExistingRemoteRepo(absDir, slug)
	}
	if err != nil {
		return "", apperror.WrapSimple(err, fmt.Sprintf("gh repo create failed: %s", string(out)))
	}

	return fmt.Sprintf("https://github.com/%s", slug), nil
}

func handleExistingRemoteRepo(absDir, slug string) (string, error) {
	remoteURL := fmt.Sprintf("https://github.com/%s.git", slug)
	cmdSet := exec.Command("git", "remote", "set-url", "origin", remoteURL)
	cmdSet.Dir = absDir
	if err := cmdSet.Run(); err != nil {
		cmdAdd := exec.Command("git", "remote", "add", "origin", remoteURL)
		cmdAdd.Dir = absDir
		_ = cmdAdd.Run()
	}
	cmdPush := exec.Command("git", "push", "-u", "origin", "main")
	cmdPush.Dir = absDir
	_ = cmdPush.Run()

	return fmt.Sprintf("https://github.com/%s", slug), nil
}
