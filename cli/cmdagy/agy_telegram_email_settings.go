// Package cmdagy — agy_telegram_email_settings.go implements unified speed settings and core CLI bindings.
package cmdagy

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AgySettingsCmd aliases the unified settings command for gitmap [agy] settings.
var AgySettingsCmd = agySettingsCmd

// SpeedSettingsStatus represents the unified speed settings state.
type SpeedSettingsStatus struct {
	LapDefaultHours         string `json:"lapDefaultHours"`
	AccountSwitchThreshold  string `json:"accountSwitchThreshold"`
	NodeResourceThreshold   string `json:"nodeResourceThreshold"`
	PaeRecencyWindowHours   string `json:"paeRecencyWindowHours"`
	TelegramBotTokenMasked  string `json:"telegramBotTokenMasked"`
	TelegramChatID          string `json:"telegramChatId"`
	TelegramReady           bool   `json:"telegramReady"`
	EmailSMTPHost           string `json:"emailSmtpHost"`
	EmailFrom               string `json:"emailFrom"`
	EmailTo                 string `json:"emailTo"`
	EmailReady              bool   `json:"emailReady"`
	MachineAlias            string `json:"machineAlias"`
	SpecialReposSecretsName string `json:"specialReposSecretsName"`
	SpecialReposCacheName   string `json:"specialReposCacheName"`
	AIAutoDefaultName       string `json:"aiAutoDefaultName"`
	AIAutoDefaultEmail      string `json:"aiAutoDefaultEmail"`
	AIAutoDefaultAgent      string `json:"aiAutoDefaultAgent"`
}

type speedSettingsView = SpeedSettingsStatus

func init() {
	AgyCmd.AddCommand(AgyTelegramCmd)
	AgyCmd.AddCommand(AgyEmailCmd)
	agySettingsCmd.DisableFlagParsing = true
	agySettingsCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return RunAgySettingsCLI(args)
	}
}

// RunAgySettingsCLI routes unified speed settings (lap default hours, threshold, telegram, email, alias) and export/import.
func RunAgySettingsCLI(args []string) error {
	if len(args) > 0 && strings.EqualFold(args[0], "export") {
		return runSettingsExportBridge(args[1:])
	}
	if len(args) > 0 && strings.EqualFold(args[0], "import") {
		return runSettingsImportBridge(args[1:])
	}
	if len(args) > 0 && isTgEmailHelpArg(args[0]) {
		RenderAgySettingsHelp()
		return nil
	}
	applySpeedSettingsFlags(args)
	return renderSpeedSettingsOverview(hasSpeedFlag(args, "--json"))
}

func runSettingsExportBridge(args []string) error {
	outPath := "antigravity_settings.json"
	if len(args) > 0 {
		outPath = args[0]
	}
	return executeAgySettingsExport(outPath)
}

func runSettingsImportBridge(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("usage: gitmap agy settings import <file.json>", "E_SETTINGS_IMPORT")
	}
	return executeAgySettingsImport(args[0])
}

func saveDBSetting(key, val string) {
	dbConn, errOpen := store.OpenDefault()
	if errOpen == nil {
		defer dbConn.Close()
		_ = dbConn.SetSetting(key, val)
	}
}

func readDBSetting(key string) string {
	dbConn, errOpen := store.OpenDefault()
	if errOpen == nil {
		defer dbConn.Close()
		val := dbConn.GetSetting(key)
		return val
	}
	return ""
}

func applySpeedSettingsFlags(args []string) {
	setGlobalIfNonEmpty("lap.default_hours", "GITMAP_LAP_DEFAULT_HOURS", extractSpeedFlag(args, "--lap-hours", "--lap"))
	setGlobalIfNonEmpty("account_switch.threshold", "GITMAP_ASW_THRESHOLD", extractSpeedFlag(args, "--switch-threshold", ""))
	setGlobalIfNonEmpty("machine.alias", "GITMAP_MACHINE_ALIAS", extractSpeedFlag(args, "--alias", "--machine-alias"))
	setGlobalIfNonEmpty("special_repos.secrets_name", "GITMAP_SPECIAL_SECRETS_REPO", extractSpeedFlag(args, "--secrets-repo", "--repo-secrets"))
	setGlobalIfNonEmpty("special_repos.cache_name", "GITMAP_SPECIAL_CACHE_REPO", extractSpeedFlag(args, "--cache-repo", "--repo-cache"))

	nodeThresh := extractSpeedFlag(args, "--node-threshold", "-t")
	if len(nodeThresh) == 0 {
		nodeThresh = extractSpeedFlag(args, "--threshold", "")
	}
	if len(nodeThresh) > 0 {
		setGlobalIfNonEmpty(constants.SettingNodeResourceThreshold, "GITMAP_NODE_RESOURCE_THRESHOLD", nodeThresh)
		saveDBSetting(constants.SettingNodeResourceThreshold, nodeThresh)
	}

	paeWindow := extractSpeedFlag(args, "--pae-window", "-r")
	if len(paeWindow) == 0 {
		paeWindow = extractSpeedFlag(args, "--recency-hours", "--pae-hours")
	}
	if len(paeWindow) > 0 {
		setGlobalIfNonEmpty(constants.SettingPullActivityWindowHours, "GITMAP_PAE_RECENCY_HOURS", paeWindow)
		saveDBSetting(constants.SettingPullActivityWindowHours, paeWindow)
	}

	applyPositionalSettingsSet(args)
}

func applyPositionalSettingsSet(args []string) {
	pos := filterSpeedPositional(args)
	if len(pos) >= 3 && strings.EqualFold(pos[0], "set") {
		applySettingKeyValue(pos[1], pos[2])
	}
}

func applySettingKeyValue(key, val string) {
	norm := strings.ToLower(strings.TrimSpace(key))
	trimmedVal := strings.TrimSpace(val)
	switch norm {
	case "node_resource_threshold", "node-threshold", "node_threshold":
		_ = config.SetVariable("global", constants.SettingNodeResourceThreshold, trimmedVal)
		saveDBSetting(constants.SettingNodeResourceThreshold, trimmedVal)
	case "pae_recency_window_hours", "pae-window", "pae_window", "recency_hours", "recency-hours", "pull_activity_window_hours":
		_ = config.SetVariable("global", constants.SettingPullActivityWindowHours, trimmedVal)
		saveDBSetting(constants.SettingPullActivityWindowHours, trimmedVal)
	case "special_repos.secrets_name", "secrets_repo", "repo_secrets":
		_ = config.SetVariable("global", "special_repos.secrets_name", trimmedVal)
	case "special_repos.cache_name", "cache_repo", "repo_cache", "repo_storage":
		_ = config.SetVariable("global", "special_repos.cache_name", trimmedVal)
	case "ai_auto.default_name", "ai_bot_name", "bot_name":
		_ = config.SetAIAutoSetting(config.SettingAIAutoDefaultName, trimmedVal)
	case "ai_auto.default_email", "ai_bot_email", "bot_email":
		_ = config.SetAIAutoSetting(config.SettingAIAutoDefaultEmail, trimmedVal)
	case "ai_auto.default_agent", "ai_agent", "default_agent":
		_ = config.SetAIAutoSetting(config.SettingAIAutoDefaultAgent, trimmedVal)
	default:
		_ = config.SetVariable("global", strings.TrimSpace(key), trimmedVal)
		saveDBSetting(strings.TrimSpace(key), trimmedVal)
	}
}

func renderSpeedSettingsOverview(isJSON bool) error {
	view := collectSpeedSettingsView()
	if isJSON {
		return printSpeedJSON(view)
	}
	runSettingsShow(view)
	return nil
}

func runSettingsShow(view speedSettingsView) {
	fmt.Printf("%s● GitMap & Antigravity Unified Speed Settings%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  LAP Default Window:        %sh (key: lap.default_hours)\n", view.LapDefaultHours)
	fmt.Printf("  Account Switch Threshold:  %s%% (key: account_switch.threshold)\n", view.AccountSwitchThreshold)
	fmt.Printf("  Node Resource Threshold:   %s%% (key: %s)\n", view.NodeResourceThreshold, constants.SettingNodeResourceThreshold)
	fmt.Printf("  PAE Recency Window:        %sh (key: %s)\n", view.PaeRecencyWindowHours, constants.SettingPullActivityWindowHours)
	fmt.Printf("  Machine Alias:             %s (key: machine.alias)\n", view.MachineAlias)
	fmt.Printf("  Telegram Chatbot:          chat=%s | token=%s | ready=%v\n", view.TelegramChatID, view.TelegramBotTokenMasked, view.TelegramReady)
	fmt.Printf("  Email Speed Setup:         smtp=%s | to=%s | ready=%v\n", view.EmailSMTPHost, view.EmailTo, view.EmailReady)
	fmt.Printf("  Special Secrets Repo  : %s (shortcut: gitmap rs / gitmap cd rs)\n", view.SpecialReposSecretsName)
	fmt.Printf("  Special Cache Repo    : %s (shortcut: gitmap rc / gitmap cd rc)\n", view.SpecialReposCacheName)
	fmt.Printf("  AI Auto Bot Name:          %s (key: ai_auto.default_name)\n", view.AIAutoDefaultName)
	fmt.Printf("  AI Auto Bot Email:         %s (key: ai_auto.default_email)\n", view.AIAutoDefaultEmail)
	fmt.Printf("  AI Auto Default Agent:     %s (key: ai_auto.default_agent)\n", view.AIAutoDefaultAgent)
}

func collectSpeedSettingsView() speedSettingsView {
	secretsName, _ := config.GetVariable("global", "special_repos.secrets_name")
	cacheName, _ := config.GetVariable("global", "special_repos.cache_name")
	st := buildBaseSpeedSettingsStatus()
	st.SpecialReposSecretsName = defaultIfEmpty(secretsName, "repo-secrets")
	st.SpecialReposCacheName = defaultIfEmpty(cacheName, "repo-cache")
	st.AIAutoDefaultName = config.GetAIAutoDefaultName()
	st.AIAutoDefaultEmail = config.GetAIAutoDefaultEmail()
	st.AIAutoDefaultAgent = config.GetAIAutoDefaultAgent()
	return st
}

func collectSpeedSettingsStatus() SpeedSettingsStatus {
	return collectSpeedSettingsView()
}

func buildBaseSpeedSettingsStatus() SpeedSettingsStatus {
	token, chatID := resolveTelegramToken(), resolveTelegramChatID()
	smtpHost, toAddr := resolveEmailSMTP(), resolveEmailTo()
	return SpeedSettingsStatus{
		LapDefaultHours:        readSpeedVar("lap.default_hours", "24"),
		AccountSwitchThreshold: readSpeedVar("account_switch.threshold", "15"),
		NodeResourceThreshold:  readSpeedVar(constants.SettingNodeResourceThreshold, fmt.Sprintf("%d", constants.DefaultNodeResourceThreshold)),
		PaeRecencyWindowHours:  readSpeedVar(constants.SettingPullActivityWindowHours, fmt.Sprintf("%d", constants.DefaultPullActivityWindowHours)),
		TelegramBotTokenMasked: maskSecretToken(token), TelegramChatID: chatID,
		TelegramReady: len(token) > 0 && len(chatID) > 0,
		EmailSMTPHost: smtpHost, EmailFrom: resolveEmailFrom(), EmailTo: toAddr,
		EmailReady: len(smtpHost) > 0 && len(toAddr) > 0, MachineAlias: readSpeedVar("machine.alias", "auto-ip"),
	}
}

func setGlobalIfNonEmpty(key, envName, val string) {
	if len(strings.TrimSpace(val)) == 0 {
		return
	}
	_ = config.SetVariable("global", key, strings.TrimSpace(val))
	_ = os.Setenv(envName, strings.TrimSpace(val))
}

func readSpeedVar(key, fallback string) string {
	if dbVal := readDBSetting(key); len(strings.TrimSpace(dbVal)) > 0 {
		return strings.TrimSpace(dbVal)
	}
	val, isFound := config.GetVariable("global", key)
	if isFound && len(strings.TrimSpace(val)) > 0 {
		return strings.TrimSpace(val)
	}
	return fallback
}

func defaultIfEmpty(val, fallback string) string {
	if len(strings.TrimSpace(val)) > 0 {
		return strings.TrimSpace(val)
	}
	return fallback
}

func maskSecretToken(token string) string {
	if len(token) == 0 {
		return "(not configured)"
	}
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "..." + token[len(token)-4:]
}
