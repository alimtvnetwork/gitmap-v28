package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

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
	if err != nil {
		return "", apperror.WrapSimple(err, fmt.Sprintf("gh repo create failed: %s", string(out)))
	}

	return fmt.Sprintf("https://github.com/%s", slug), nil
}
