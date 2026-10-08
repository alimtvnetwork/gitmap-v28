package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func resolveBaseSlug(slug, name string) string {
	if slug != "" {
		return slug
	}

	return SlugifyRepoName(name)
}

func isProfileSlugPrefix(prof model.GitProfile) bool {
	if prof.Name == "" || prof.Name == "default" {
		return false
	}

	if strings.Contains(prof.Name, " ") {
		return false
	}

	return true
}

func buildRemoteSlug(p createRepoParams) string {
	slug := resolveBaseSlug(p.Slug, p.Name)
	if strings.Contains(slug, "/") {
		return slug
	}

	if isProfileSlugPrefix(p.Profile) {
		return p.Profile.Name + "/" + slug
	}

	return slug
}

func resolveVisibilityFlag(isPublic bool) string {
	if isPublic {
		return "--public"
	}

	return "--private"
}

func executeCreateRemoteRepo(absDir, slug, visibilityFlag string) (string, error) {
	cmd := exec.Command("gh", "repo", "create", slug, visibilityFlag, "--source=.", "--remote=origin", "--push")
	cmd.Dir = absDir
	out, err := cmd.CombinedOutput()
	if err != nil && isRepoCollisionOutput(string(out)) {
		return handleExistingRemoteRepo(absDir, slug)
	}

	if err != nil {
		return "", apperror.WrapSimple(err, fmt.Sprintf("gh repo create failed: %s", string(out)))
	}

	sshURL := fmt.Sprintf("git@github.com:%s.git", slug)
	_ = exec.Command("git", "-C", absDir, "remote", "set-url", "origin", sshURL).Run()

	return fmt.Sprintf("https://github.com/%s", slug), nil
}

func pushRemoteRepo(p createRepoParams) (string, error) {
	if p.IsSkipRemote {
		return "", nil
	}

	absDir, _ := filepath.Abs(p.LocalDir)
	visibilityFlag := resolveVisibilityFlag(p.IsPublic)
	slug := buildRemoteSlug(p)

	hasRemote, probeErr := probeRemoteRepoExists(slug)
	if probeErr == nil && hasRemote {
		return handleExistingRemoteRepo(absDir, slug)
	}

	return executeCreateRemoteRepo(absDir, slug, visibilityFlag)
}

func handleExistingRemoteRepo(absDir, slug string) (string, error) {
	return handleRemoteCollision(absDir, slug)
}
