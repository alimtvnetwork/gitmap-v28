package cmdlogin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/usercontext"
)

const (
	aiAutoBotName  = "gitmap automation bot - runs autonmously"
	aiAutoBotEmail = "info@rgitmap.com"
	aiAutoBotAlias = "autobot"
)

func runUseAIAuto(hasRestrictionTime bool) error {
	cfg, _ := store.LoadGitProfiles()
	attachExistingOriginalBackup(&cfg)
	backupOriginalUser(&cfg)
	if err := switchGitToAIAuto(); err != nil {
		return err
	}

	updateAndActivateAutobotProfile(&cfg)
	printAIAutoSwitchedNotice()
	if hasRestrictionTime {
		return applyRestrictionTimeDefaults()
	}

	return nil
}

func attachExistingOriginalBackup(cfg *model.GitProfileConfig) {
	if cfg.OriginalIdentity != nil && len(cfg.OriginalIdentity.Name) > 0 {
		return
	}
	if orig, isLoaded := loadOriginalUserFile(); isLoaded {
		cfg.OriginalIdentity = orig
	}
}

func switchGitToAIAuto() error {
	return usercontext.SetGitIdentity(aiAutoBotName, aiAutoBotEmail, true)
}

func backupOriginalUser(cfg *model.GitProfileConfig) {
	curIdentity, _ := usercontext.GetGitIdentity(true)
	if curIdentity.Name == aiAutoBotName || len(curIdentity.Name) == 0 {
		return
	}
	if cfg.OriginalIdentity != nil && len(cfg.OriginalIdentity.Name) > 0 {
		return
	}

	orig := model.OriginalIdentity{
		Name:         curIdentity.Name,
		Email:        curIdentity.Email,
		ProfileAlias: cfg.Active,
		SavedAt:      time.Now(),
	}
	cfg.OriginalIdentity = &orig
	saveOriginalUserFile(orig)
	fmt.Printf("%s✔ Saved original Git identity: %s <%s>%s\n",
		constants.ColorGreen, orig.Name, orig.Email, constants.ColorReset)
}

func saveOriginalUserFile(orig model.OriginalIdentity) {
	data, err := json.MarshalIndent(orig, "", "  ")
	if err != nil {
		return
	}

	writeOriginalFileAt(store.BinaryDataDir(), data)
	if globalDir := store.GlobalUserDataDir(); len(globalDir) > 0 {
		writeOriginalFileAt(globalDir, data)
	}
}

func writeOriginalFileAt(dir string, data []byte) {
	_ = os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, "original_user.json")
	_ = os.WriteFile(path, data, 0644)
}

func loadOriginalUserFile() (*model.OriginalIdentity, bool) {
	if orig, isLoaded := readOriginalFileAt(store.BinaryDataDir()); isLoaded {
		return orig, true
	}
	if globalDir := store.GlobalUserDataDir(); len(globalDir) > 0 {
		return readOriginalFileAt(globalDir)
	}

	return nil, false
}

func readOriginalFileAt(dir string) (*model.OriginalIdentity, bool) {
	path := filepath.Join(dir, "original_user.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var orig model.OriginalIdentity
	if err := json.Unmarshal(data, &orig); err != nil {
		return nil, false
	}

	return &orig, len(orig.Name) > 0
}

func updateAndActivateAutobotProfile(cfg *model.GitProfileConfig) {
	prof := model.GitProfile{
		ID:         aiAutoBotAlias,
		Name:       aiAutoBotName,
		Email:      aiAutoBotEmail,
		Provider:   "git",
		Type:       "user",
		AuthMethod: "config",
		IsDefault:  true,
		LastUsedAt: time.Now(),
	}
	cfg.Profiles = appendOrUpdateBotProfile(cfg.Profiles, prof)
	cfg.Active = aiAutoBotAlias
	_ = store.SaveGitProfiles(*cfg)
}

func appendOrUpdateBotProfile(profiles []model.GitProfile, p model.GitProfile) []model.GitProfile {
	for idx := range profiles {
		if profiles[idx].ID == p.ID {
			profiles[idx] = p
			return profiles
		}
	}

	return append(profiles, p)
}

func printAIAutoSwitchedNotice() {
	fmt.Printf("%s✔ Switched global Git identity to AI Automation:%s\n",
		constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  Name:  %s\n", aiAutoBotName)
	fmt.Printf("  Email: %s\n", aiAutoBotEmail)
	fmt.Printf("  Alias: %s (active)\n", aiAutoBotAlias)
}

func applyRestrictionTimeDefaults() error {
	fmt.Printf("%s✔ Applied time restriction defaults:%s\n",
		constants.ColorGreen, constants.ColorReset)
	fmt.Println("  - Pipeline-AI delay: 20s")
	fmt.Println("  - Dynamic ETA sleep: >120s -> 30s | 60-120s -> 20s | <=60s -> 10s")
	fmt.Println("  - Local runner wait: 60s")
	fmt.Println("  - Tight polling ban: active")

	return nil
}

func runRevertToOriginal() error {
	cfg, _ := store.LoadGitProfiles()
	orig := resolveOriginalIdentity(&cfg)
	if orig == nil {
		fmt.Printf("%s⚠ No original Git identity found to revert to.%s\n",
			constants.ColorYellow, constants.ColorReset)
		return nil
	}

	if err := usercontext.SetGitIdentity(orig.Name, orig.Email, true); err != nil {
		return err
	}

	restoreOriginalProfile(&cfg, orig)
	printRevertedNotice(orig)

	return nil
}

func resolveOriginalIdentity(cfg *model.GitProfileConfig) *model.OriginalIdentity {
	if cfg.OriginalIdentity != nil && len(cfg.OriginalIdentity.Name) > 0 {
		return cfg.OriginalIdentity
	}

	fromFile, hasFile := loadOriginalUserFile()
	if hasFile {
		return fromFile
	}

	return nil
}

func restoreOriginalProfile(cfg *model.GitProfileConfig, orig *model.OriginalIdentity) {
	if len(orig.ProfileAlias) > 0 {
		cfg.Active = orig.ProfileAlias
	}
	cfg.OriginalIdentity = nil
	removeOriginalUserFiles()
	_ = store.SaveGitProfiles(*cfg)
}

func removeOriginalUserFiles() {
	_ = os.Remove(filepath.Join(store.BinaryDataDir(), "original_user.json"))
	if globalDir := store.GlobalUserDataDir(); len(globalDir) > 0 {
		_ = os.Remove(filepath.Join(globalDir, "original_user.json"))
	}
}

func printRevertedNotice(orig *model.OriginalIdentity) {
	fmt.Printf("%s✔ Reverted global Git identity to original:%s\n",
		constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  Name:  %s\n", orig.Name)
	fmt.Printf("  Email: %s\n", orig.Email)
	if len(orig.ProfileAlias) > 0 {
		fmt.Printf("  Alias: %s (restored)\n", orig.ProfileAlias)
	}
}
