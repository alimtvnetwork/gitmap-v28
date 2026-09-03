package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runProfilesSetDefault(args []string) error {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return apperror.WrapSimple(err, "load profiles:")
	}
	target, resolveErr := resolveTargetProfileArg(cfg.Profiles, args, "Select default profile")
	if resolveErr != nil {
		return resolveErr
	}
	return applyDefaultProfile(&cfg, target)
}

func applyDefaultProfile(cfg *model.GitProfileConfig, target model.GitProfile) error {
	for i := range cfg.Profiles {
		cfg.Profiles[i].IsDefault = (cfg.Profiles[i].Name == target.Name)
	}
	cfg.Default = target.Name
	cfg.Active = target.Name
	saveErr := store.SaveGitProfiles(*cfg)
	if saveErr != nil {
		return apperror.WrapSimple(saveErr, "save profiles:")
	}
	fmt.Printf("  %s✓ Default Git profile set to: %s (%s)%s\n",
		constants.ColorGreen, target.Name, target.Provider, constants.ColorReset)
	return nil
}

func runProfilesSwitch(args []string) error {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return apperror.WrapSimple(err, "load profiles:")
	}
	target, resolveErr := resolveTargetProfileArg(cfg.Profiles, args, "Switch active profile")
	if resolveErr != nil {
		return resolveErr
	}
	cfg.Active = target.Name
	saveErr := store.SaveGitProfiles(cfg)
	if saveErr != nil {
		return apperror.WrapSimple(saveErr, "save profiles:")
	}
	fmt.Printf("  %s✓ Active profile switched to: %s%s\n",
		constants.ColorGreen, target.Name, constants.ColorReset)
	return nil
}

func runProfilesAdd(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("usage: gitmap profiles add <name> [--provider github|gitlab] [--org]", "E1072")
	}
	name := args[0]
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return apperror.WrapSimple(err, "load profiles:")
	}
	newProfile := buildNewProfile(name, args)
	cfg.Profiles = append(cfg.Profiles, newProfile)
	saveErr := store.SaveGitProfiles(cfg)
	if saveErr != nil {
		return apperror.WrapSimple(saveErr, "save profiles:")
	}
	fmt.Printf("  %s✓ Registered profile: %s (%s, %s)%s\n",
		constants.ColorGreen, newProfile.Name, newProfile.Provider, newProfile.Type, constants.ColorReset)
	return nil
}

func buildNewProfile(name string, args []string) model.GitProfile {
	provider := "github"
	if hasArgFlag(args, "--gitlab") {
		provider = "gitlab"
	}
	profileType := "user"
	if hasArgFlag(args, "--org") {
		profileType = "organization"
	}
	email := extractFlagVal(args, "--email")
	return model.GitProfile{
		ID:         name,
		Name:       name,
		Provider:   provider,
		Type:       profileType,
		Email:      email,
		AuthMethod: "gh-cli",
		LastUsedAt: time.Now(),
	}
}

func runProfilesRemove(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("usage: gitmap profiles rm <name|1-N>", "E1073")
	}
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return apperror.WrapSimple(err, "load profiles:")
	}
	idx, target, findErr := pickProfileBySequenceOrName(cfg.Profiles, args[0])
	if findErr != nil {
		return findErr
	}
	confirmMsg := fmt.Sprintf("Remove profile '%s'?", target.Name)
	if !confirmOrSkip(confirmMsg, args) {
		fmt.Println("  Aborted.")
		return nil
	}
	cfg.Profiles = append(cfg.Profiles[:idx], cfg.Profiles[idx+1:]...)
	return store.SaveGitProfiles(cfg)
}

func runProfilesStatus(_ []string) error {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return apperror.WrapSimple(err, "load profiles:")
	}
	fmt.Printf("\n  %s● Active Git Profile:%s  %s\n", constants.ColorCyan, constants.ColorReset, cfg.Active)
	fmt.Printf("  %s● Default Git Profile:%s %s\n", constants.ColorCyan, constants.ColorReset, cfg.Default)
	fmt.Printf("  %s● Profiles Configured:%s %d\n\n", constants.ColorCyan, constants.ColorReset, len(cfg.Profiles))
	return nil
}

func resolveTargetProfileArg(profiles []model.GitProfile, args []string, promptTitle string) (model.GitProfile, error) {
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		_, prof, err := pickProfileBySequenceOrName(profiles, args[0])
		return prof, err
	}
	if !isInteractiveStdin() {
		return model.GitProfile{}, apperror.NewSimple("missing profile argument in non-interactive mode", "E1074")
	}
	return promptSequenceSelection(profiles, promptTitle)
}

func promptSequenceSelection(profiles []model.GitProfile, promptTitle string) (model.GitProfile, error) {
	fmt.Printf("\n%s (Enter number 1-%d):\n", promptTitle, len(profiles))
	for i, p := range profiles {
		fmt.Printf("  [%d] %s (%s, %s)\n", i+1, p.Name, p.Provider, p.Type)
	}
	fmt.Print("Choice: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return model.GitProfile{}, apperror.WrapSimple(err, "read choice:")
	}
	val := strings.TrimSpace(line)
	_, prof, pickErr := pickProfileBySequenceOrName(profiles, val)
	return prof, pickErr
}

func pickProfileBySequenceOrName(profiles []model.GitProfile, val string) (int, model.GitProfile, error) {
	num, err := strconv.Atoi(val)
	if err == nil && num >= 1 && num <= len(profiles) {
		return num - 1, profiles[num-1], nil
	}
	lower := strings.ToLower(val)
	for i, p := range profiles {
		if strings.ToLower(p.Name) == lower || strings.ToLower(p.ID) == lower {
			return i, p, nil
		}
	}
	return -1, model.GitProfile{}, apperror.NewSimple(fmt.Sprintf("profile not found: %s", val), "E1075")
}
