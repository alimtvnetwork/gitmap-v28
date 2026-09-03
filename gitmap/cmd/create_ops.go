package cmd

import (
	"encoding/json"
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

type createRepoParams struct {
	Name        string
	LocalDir    string
	Description string
	IsPublic    bool
	NoRemote    bool
	IsJSON      bool
	Profile     model.GitProfile
}

func executeCreateRepo(args []string) error {
	params, parseErr := parseCreateParams(args)
	if parseErr != nil {
		return parseErr
	}
	if initErr := initLocalRepo(params); initErr != nil {
		return initErr
	}
	remoteURL, pushErr := pushRemoteRepo(params)
	if pushErr != nil {
		return pushErr
	}
	recordProfileUsage(params.Profile)
	return reportCreatedRepo(params, remoteURL)
}

func parseCreateParams(args []string) (createRepoParams, error) {
	name := args[0]
	isPublic := hasArgFlag(args, "--public")
	noRemote := hasArgFlag(args, "--no-remote")
	isJSON := hasArgFlag(args, "--json")
	desc := extractFlagVal(args, "--description")
	if desc == "" {
		desc = extractFlagVal(args, "-d")
	}
	dir := extractFlagVal(args, "--dir")
	if dir == "" {
		dir = filepath.Join(".", name)
	}
	prof, profErr := resolveCreationProfile(args)
	if profErr != nil {
		return createRepoParams{}, profErr
	}
	return createRepoParams{
		Name: name, LocalDir: dir, Description: desc,
		IsPublic: isPublic, NoRemote: noRemote, IsJSON: isJSON, Profile: prof,
	}, nil
}

func resolveCreationProfile(args []string) (model.GitProfile, error) {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return model.GitProfile{}, apperror.WrapSimple(err, "load profiles:")
	}
	req := extractFlagVal(args, "--profile")
	if req == "" {
		req = extractFlagVal(args, "--org")
	}
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
	return model.GitProfile{Name: "default", Provider: "github", Type: "user"}, nil
}

func initLocalRepo(p createRepoParams) error {
	absDir, absErr := filepath.Abs(p.LocalDir)
	if absErr != nil {
		return apperror.WrapSimple(absErr, "resolve absolute dir:")
	}
	if mkErr := os.MkdirAll(absDir, 0755); mkErr != nil {
		return apperror.WrapSimple(mkErr, "create directory:")
	}
	cmdInit := exec.Command("git", "init", "-b", "main")
	cmdInit.Dir = absDir
	if initErr := cmdInit.Run(); initErr != nil {
		return apperror.WrapSimple(initErr, "git init:")
	}
	writeInitialFiles(absDir, p)
	return commitInitialFiles(absDir)
}

func writeInitialFiles(absDir string, p createRepoParams) {
	readmePath := filepath.Join(absDir, "README.md")
	readmeContent := fmt.Sprintf("# %s\n\n%s\n", p.Name, p.Description)
	_ = os.WriteFile(readmePath, []byte(readmeContent), 0644)

	gitignorePath := filepath.Join(absDir, ".gitignore")
	gitignoreContent := ".DS_Store\nThumbs.db\nnode_modules/\nbin/\n*.log\n"
	_ = os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
}

func commitInitialFiles(absDir string) error {
	cmdAdd := exec.Command("git", "add", ".")
	cmdAdd.Dir = absDir
	_ = cmdAdd.Run()

	cmdCommit := exec.Command("git", "commit", "-m", "feat: initial commit")
	cmdCommit.Dir = absDir
	_ = cmdCommit.Run()
	return nil
}

func pushRemoteRepo(p createRepoParams) (string, error) {
	if p.NoRemote {
		return "", nil
	}
	absDir, _ := filepath.Abs(p.LocalDir)
	visibilityFlag := "--private"
	if p.IsPublic {
		visibilityFlag = "--public"
	}
	slug := p.Name
	if p.Profile.Name != "" && p.Profile.Name != "default" {
		slug = p.Profile.Name + "/" + p.Name
	}
	cmd := exec.Command("gh", "repo", "create", slug, visibilityFlag, "--source=.", "--remote=origin", "--push")
	cmd.Dir = absDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", apperror.NewSimple(fmt.Sprintf("gh repo create failed: %s", string(out)), "E1078")
	}
	return fmt.Sprintf("https://github.com/%s", slug), nil
}

func recordProfileUsage(prof model.GitProfile) {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return
	}
	for i := range cfg.Profiles {
		if cfg.Profiles[i].Name == prof.Name {
			cfg.Profiles[i].UsageCount++
			cfg.Profiles[i].LastUsedAt = time.Now()
			break
		}
	}
	_ = store.SaveGitProfiles(cfg)
}

func reportCreatedRepo(p createRepoParams, remoteURL string) error {
	absDir, _ := filepath.Abs(p.LocalDir)
	if p.IsJSON {
		res := map[string]string{
			"name": p.Name, "path": absDir, "remoteUrl": remoteURL,
			"profile": p.Profile.Name, "provider": p.Profile.Provider,
		}
		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(data))
		return nil
	}
	fmt.Printf("\n  %s✓ Repository created successfully!%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  ● Name:      %s\n", p.Name)
	fmt.Printf("  ● Path:      %s\n", absDir)
	fmt.Printf("  ● Profile:   %s (%s)\n", p.Profile.Name, p.Profile.Provider)
	if remoteURL != "" {
		fmt.Printf("  ● Remote:    %s\n", remoteURL)
	}
	fmt.Println()
	return nil
}
