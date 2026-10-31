package cmdlogin

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	loginconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdlogin/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	appconstants "github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/usercontext"
)

const (
	aiAutoBotName  = "gitmap automation bot - runs autonomously"
	aiAutoBotEmail = "info@rgitmap.com"
	aiAutoBotAlias = "autobot"
)

func runUseAIAuto(hasRestrictionTime bool) error {
	return runUseAIAutoWithAgent("", true, hasRestrictionTime)
}

func runUseAIAutoWithAgent(agentName string, isGlobal bool, hasRestrictionTime bool) error {
	effectiveAgent := resolveEffectiveAgent(agentName)
	if !isGlobal && len(effectiveAgent) == 0 {
		return formatMissingAgentNameError()
	}

	botName := resolveAIBotName(effectiveAgent)
	botEmail := config.GetAIAutoDefaultEmail()
	cfg := prepareAIAutoProfiles()
	targetGlobal := resolveTargetGlobal(isGlobal)
	if err := usercontext.SetGitIdentity(botName, botEmail, targetGlobal); err != nil {
		return err
	}

	return applyAIAutoPostSwitch(&cfg, botName, targetGlobal, hasRestrictionTime)
}

func resolveEffectiveAgent(agentName string) string {
	clean := strings.Trim(strings.TrimSpace(agentName), "()[]{}")
	if len(clean) > 0 {
		return clean
	}

	return config.GetAIAutoDefaultAgent()
}

func formatMissingAgentNameError() error {
	msg := strings.Join([]string{
		"AI agent name is mandatory for 'gitmap login set ai-auto <agent-name>'.",
		"",
		"💡 Examples to teach on the fly:",
		"  • gitmap login set ai-auto muse",
		"  • gitmap login set ai-auto gemini",
		"  • gitmap login set ai-auto cursor-model",
		"  • gitmap login set ai-auto antigravity",
		"  • gitmap login set ai-auto claude",
		"",
		"To configure a generic machine-wide bot without an agent name, use:",
		"  • gitmap login set global ai-auto",
		"",
		"Or configure a default AI agent for this machine in Settings:",
		"  • gitmap settings set ai_auto.default_agent <agent-name> (e.g. antigravity, gemini, cursor-model, muse)",
	}, "\n")

	return apperror.NewValidationError(msg)
}

func prepareAIAutoProfiles() model.GitProfileConfig {
	cfg, _ := store.LoadGitProfiles()
	attachExistingOriginalBackup(&cfg)
	backupOriginalUser(&cfg)

	return cfg
}

func resolveTargetGlobal(isGlobal bool) bool {
	if isGlobal {
		return true
	}

	return !isInsideGitRepo()
}

func isInsideGitRepo() bool {
	cmd := exec.Command(loginconstants.GitBin, loginconstants.GitRevParse, loginconstants.GitInsideWorkTree)

	return cmd.Run() == nil
}

func applyAIAutoPostSwitch(cfg *model.GitProfileConfig, botName string, isGlobal bool, hasRestriction bool) error {
	updateAndActivateAutobotProfile(cfg, botName)
	printAIAutoSwitchedNotice(botName, isGlobal)
	if hasRestriction {
		return applyRestrictionTimeDefaults()
	}

	return nil
}

func resolveAIBotName(agentName string) string {
	clean := strings.Trim(strings.TrimSpace(agentName), "()[]{}")
	if len(clean) == 0 {
		return config.GetAIAutoDefaultName()
	}

	return fmt.Sprintf("gitmap automation bot - %s - runs autonomously", clean)
}

func isAIBotIdentity(name string) bool {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 {
		return true
	}

	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "gitmap automation") {
		return true
	}

	return strings.Contains(lower, "autonomously") || strings.Contains(lower, "autonmously")
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
	if isAIBotIdentity(curIdentity.Name) {
		return
	}

	if cfg.OriginalIdentity != nil && len(cfg.OriginalIdentity.Name) > 0 {
		return
	}

	persistOriginalBackup(cfg, curIdentity)
}

func persistOriginalBackup(cfg *model.GitProfileConfig, identity usercontext.GitIdentity) {
	orig := model.OriginalIdentity{
		Name:         identity.Name,
		Email:        identity.Email,
		ProfileAlias: cfg.Active,
		SavedAt:      time.Now(),
	}
	cfg.OriginalIdentity = &orig
	saveOriginalUserFile(orig)
	printBackupSavedNotice(orig)
}

func printBackupSavedNotice(orig model.OriginalIdentity) {
	fmt.Printf("%s✔ Saved original Git identity: %s <%s>%s\n",
		appconstants.ColorGreen, orig.Name, orig.Email, appconstants.ColorReset)
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

func updateAndActivateAutobotProfile(cfg *model.GitProfileConfig, botName string) {
	prof := model.GitProfile{
		ID:         aiAutoBotAlias,
		Name:       botName,
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

func printAIAutoSwitchedNotice(botName string, isGlobal bool) {
	fmt.Printf("%s✔ Switched %s Git identity to AI Automation:%s\n",
		appconstants.ColorGreen, scopeLabel(isGlobal), appconstants.ColorReset)
	fmt.Printf("  Name:  %s\n", botName)
	fmt.Printf("  Email: %s\n", aiAutoBotEmail)
	fmt.Printf("  Alias: %s (active)\n", aiAutoBotAlias)
}

func scopeLabel(isGlobal bool) string {
	if isGlobal {
		return "global"
	}

	return "local"
}

func applyRestrictionTimeDefaults() error {
	fmt.Printf("%s✔ Applied time restriction defaults:%s\n",
		appconstants.ColorGreen, appconstants.ColorReset)
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
			appconstants.ColorYellow, appconstants.ColorReset)
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
		appconstants.ColorGreen, appconstants.ColorReset)
	fmt.Printf("  Name:  %s\n", orig.Name)
	fmt.Printf("  Email: %s\n", orig.Email)
	if len(orig.ProfileAlias) > 0 {
		fmt.Printf("  Alias: %s (restored)\n", orig.ProfileAlias)
	}
}
