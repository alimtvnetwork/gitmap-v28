package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"gopkg.in/yaml.v3"
)

// RepoCreateReport represents structured repository creation status.
type RepoCreateReport struct {
	Name       string `json:"name" yaml:"name"`
	Slug       string `json:"slug" yaml:"slug"`
	Path       string `json:"path" yaml:"path"`
	Visibility string `json:"visibility" yaml:"visibility"`
	RemoteURL  string `json:"remoteUrl" yaml:"remoteUrl"`
	Profile    string `json:"profile" yaml:"profile"`
	Provider   string `json:"provider" yaml:"provider"`
}

func reportCreatedRepo(p createRepoParams, remoteURL string) error {
	absDir, _ := filepath.Abs(p.LocalDir)
	visibility := getRepoVisibility(p)
	if p.IsJSON {
		return renderJSONReport(p, absDir, visibility, remoteURL)
	}
	if p.IsYAML {
		return renderYAMLReport(p, absDir, visibility, remoteURL)
	}

	return renderTextReport(p, absDir, visibility, remoteURL)
}

func getRepoVisibility(p createRepoParams) string {
	if p.IsPublic {
		return "public"
	}

	return "private"
}

func buildRepoCreateReport(p createRepoParams, absDir, visibility, remoteURL string) RepoCreateReport {
	return RepoCreateReport{
		Name:       p.Name,
		Slug:       p.Slug,
		Path:       absDir,
		Visibility: visibility,
		RemoteURL:  remoteURL,
		Profile:    p.Profile.Name,
		Provider:   p.Profile.Provider,
	}
}

func renderJSONReport(p createRepoParams, absDir, visibility, remoteURL string) error {
	rep := buildRepoCreateReport(p, absDir, visibility, remoteURL)
	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal json report:")
	}
	fmt.Println(string(data))

	return nil
}

func renderYAMLReport(p createRepoParams, absDir, visibility, remoteURL string) error {
	rep := buildRepoCreateReport(p, absDir, visibility, remoteURL)
	data, err := yaml.Marshal(rep)
	if err != nil {
		return apperror.WrapSimple(err, "marshal yaml report:")
	}
	fmt.Print(string(data))

	return nil
}

func renderTextReport(p createRepoParams, absDir, visibility, remoteURL string) error {
	fmt.Printf("\n  %s✓ Repository created successfully!%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  ● Name:       %s\n", p.Name)
	fmt.Printf("  ● Slug:       %s\n", p.Slug)
	fmt.Printf("  ● Path:       %s\n", absDir)
	fmt.Printf("  ● Visibility: %s\n", visibility)
	fmt.Printf("  ● Profile:    %s (%s)\n", p.Profile.Name, p.Profile.Provider)
	if remoteURL != "" {
		fmt.Printf("  ● Remote:     %s\n", remoteURL)
	}
	fmt.Println()

	return nil
}
