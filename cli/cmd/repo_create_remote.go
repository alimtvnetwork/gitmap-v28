package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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

func reportCreatedRepo(p createRepoParams, remoteURL string) error {
	absDir, _ := filepath.Abs(p.LocalDir)
	if p.IsJSON {
		res := map[string]string{
			"name": p.Name, "slug": p.Slug, "path": absDir, "remoteUrl": remoteURL,
			"profile": p.Profile.Name, "provider": p.Profile.Provider,
		}

		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(data))

		return nil
	}

	fmt.Printf("\n  %s✓ Repository created successfully!%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  ● Name:      %s\n", p.Name)
	fmt.Printf("  ● Slug:      %s\n", p.Slug)
	fmt.Printf("  ● Path:      %s\n", absDir)
	fmt.Printf("  ● Profile:   %s (%s)\n", p.Profile.Name, p.Profile.Provider)
	if remoteURL != "" {
		fmt.Printf("  ● Remote:    %s\n", remoteURL)
	}

	fmt.Println()

	return nil
}
