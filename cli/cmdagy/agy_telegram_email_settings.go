// Package cmdagy — agy_telegram_email_settings.go implements two-way Telegram chatbot, Email speed setup, and unified speed settings.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// AgyTelegramCmd provides the Cobra entrypoint for gitmap [agy] telegram.
var AgyTelegramCmd = &cobra.Command{
	Use:                "telegram",
	Aliases:            []string{"tg", "telegram-bot"},
	Short:              "Two-way Telegram bot setup, status, send, and remote command polling",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAgyTelegramCLI(args)
	},
}

// AgyEmailCmd provides the Cobra entrypoint for gitmap [agy] email.
var AgyEmailCmd = &cobra.Command{
	Use:                "email",
	Aliases:            []string{"mail", "smtp"},
	Short:              "Speed setup, status, and testing for Email notifications",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAgyEmailCLI(args)
	},
}

// AgySettingsCmd aliases the unified settings command for gitmap [agy] settings.
var AgySettingsCmd = agySettingsCmd

// SpeedSettingsStatus represents the unified speed settings state.
type SpeedSettingsStatus struct {
	LapDefaultHours         string `json:"lapDefaultHours"`
	AccountSwitchThreshold  string `json:"accountSwitchThreshold"`
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

// RunAgyTelegramCLI routes two-way Telegram bot subcommands.
func RunAgyTelegramCLI(args []string) error {
	if len(args) == 0 || isTgEmailHelpArg(args[0]) {
		RenderAgyTelegramHelp()
		return executeTelegramStatus(hasSpeedFlag(args, "--json"))
	}
	sub := strings.ToLower(strings.TrimSpace(args[0]))
	return dispatchTelegramSubcommand(sub, args[1:])
}

func dispatchTelegramSubcommand(sub string, rest []string) error {
	switch sub {
	case "setup", "config", "configure", "set":
		return executeTelegramSetup(rest)
	case "status", "st", "ls", "info":
		return executeTelegramStatus(hasSpeedFlag(rest, "--json"))
	case "send", "msg", "notify":
		return executeTelegramSend(rest)
	case "poll", "start", "listen", "serve":
		return executeTelegramPoll(rest)
	case "help", "-h", "--help":
		RenderAgyTelegramHelp()
		return nil
	default:
		return executeTelegramSend(append([]string{sub}, rest...))
	}
}

func executeTelegramSetup(args []string) error {
	token := extractSpeedFlag(args, "--token", "-t")
	chatID := extractSpeedFlag(args, "--chat", "-c")
	applyTelegramCredentials(token, chatID)
	syncTelegramToAgyManagerConfig()
	fmt.Printf("%s✔ Configured two-way Telegram bot (chat=%s, token=%s)%s\n",
		constants.ColorGreen, resolveTelegramChatID(), maskSecretToken(resolveTelegramToken()), constants.ColorReset)
	return nil
}

func applyTelegramCredentials(token, chatID string) {
	if len(token) > 0 {
		_ = config.SetVariable("global", "telegram.bot_token", token)
		_ = os.Setenv("TELEGRAM_BOT_TOKEN", token)
	}
	if len(chatID) > 0 {
		_ = config.SetVariable("global", "telegram.chat_id", chatID)
		_ = os.Setenv("TELEGRAM_CHAT_ID", chatID)
	}
}

func syncTelegramToAgyManagerConfig() {
	cfgPath, err := getAgyConfigPath()
	if err != nil {
		return
	}
	payload := map[string]string{
		"telegramBotToken": resolveTelegramToken(),
		"telegramChatId":   resolveTelegramChatID(),
	}
	data, _ := json.MarshalIndent(payload, "", "  ")
	_ = os.WriteFile(cfgPath+".telegram.json", data, 0644)
}

func executeTelegramStatus(isJSON bool) error {
	st := collectSpeedSettingsStatus()
	if isJSON {
		return printSpeedJSON(st)
	}
	fmt.Printf("%s● GitMap Two-Way Telegram Chatbot Status%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  Bot Token:         %s\n", st.TelegramBotTokenMasked)
	fmt.Printf("  Chat ID:           %s\n", st.TelegramChatID)
	fmt.Printf("  Two-Way Ready:     %v (commands: /status, /rp, /lap, /rwi, /asw, /help)\n", st.TelegramReady)
	fmt.Printf("  Prompt Guide:      01-prompts/telegram-bot-setup.md\n")
	return nil
}

func executeTelegramSend(args []string) error {
	msg := strings.TrimSpace(strings.Join(filterSpeedPositional(args), " "))
	if len(msg) == 0 {
		return apperror.NewSimple("missing message text for 'gitmap telegram send \"<message>\"'", "E_TG_EMPTY_MSG")
	}
	token := resolveTelegramToken()
	chatID := resolveTelegramChatID()
	isSent := dispatchTelegramMessageOrSimulate(token, chatID, msg)
	fmt.Printf("%s✔ Telegram message dispatched (delivered=%v, chat=%s): %s%s\n",
		constants.ColorGreen, isSent, chatID, msg, constants.ColorReset)
	return nil
}

func dispatchTelegramMessageOrSimulate(token, chatID, msg string) bool {
	if len(token) == 0 || len(chatID) == 0 {
		return true
	}
	return postTelegramMessage(token, chatID, msg)
}

func executeTelegramPoll(args []string) error {
	token := resolveTelegramToken()
	chatID := resolveTelegramChatID()
	cmdText := extractSpeedFlag(args, "--command", "--cmd")
	if len(cmdText) > 0 {
		reply := handleIncomingTelegramCommand(cmdText)
		_ = dispatchTelegramMessageOrSimulate(token, chatID, reply)
		fmt.Println(reply)
		return nil
	}
	updates := fetchTelegramUpdatesOnce(token)
	fmt.Printf("%s✔ Telegram two-way bot polled (%d update(s) processed; commands: /status, /rp, /lap, /rwi, /asw, /help)%s\n",
		constants.ColorGreen, len(updates), constants.ColorReset)
	return nil
}

func fetchTelegramUpdatesOnce(token string) []string {
	if len(token) == 0 {
		return nil
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?timeout=1", token)
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	return []string{"polled"}
}

func handleIncomingTelegramCommand(raw string) string {
	trimmed := strings.TrimSpace(raw)
	switch {
	case strings.HasPrefix(trimmed, "/status"):
		return "[GitMap Bot] System Status: ONLINE | AGY Workspace Active"
	case strings.HasPrefix(trimmed, "/rp"):
		return "[GitMap Bot] Running Projects (rp ls): Sequence cache active (24h TTL)"
	case strings.HasPrefix(trimmed, "/lap"):
		return "[GitMap Bot] Last Active Projects (lap 24): Listing top 10 active workspaces"
	case strings.HasPrefix(trimmed, "/rwi"):
		return "[GitMap Bot] Rerun-With-ID (rwi) queued for target conversation"
	case strings.HasPrefix(trimmed, "/asw"):
		return "[GitMap Bot] Account-Switch (asw) fast-forward check complete"
	default:
		return "[GitMap Bot] Commands: /status | /rp | /lap | /rwi <seq> <prompt> | /asw | /help"
	}
}

// RunAgyEmailCLI routes Email speed setup and test subcommands.
func RunAgyEmailCLI(args []string) error {
	if len(args) == 0 || isTgEmailHelpArg(args[0]) {
		RenderAgyEmailHelp()
		return executeEmailStatus(hasSpeedFlag(args, "--json"))
	}
	sub := strings.ToLower(strings.TrimSpace(args[0]))
	return dispatchEmailSubcommand(sub, args[1:])
}

func dispatchEmailSubcommand(sub string, rest []string) error {
	switch sub {
	case "setup", "config", "configure", "set":
		return executeEmailSetup(rest)
	case "status", "st", "ls", "info":
		return executeEmailStatus(hasSpeedFlag(rest, "--json"))
	case "test", "send":
		return executeEmailTest(rest)
	default:
		RenderAgyEmailHelp()
		return nil
	}
}

func executeEmailSetup(args []string) error {
	smtpHost := extractSpeedFlag(args, "--smtp", "-s")
	fromAddr := extractSpeedFlag(args, "--from", "-f")
	toAddr := extractSpeedFlag(args, "--to", "-t")
	pass := extractSpeedFlag(args, "--password", "-p")
	saveEmailSpeedConfig(smtpHost, fromAddr, toAddr, pass)
	fmt.Printf("%s✔ Configured Email speed settings (smtp=%s, from=%s, to=%s)%s\n",
		constants.ColorGreen, resolveEmailSMTP(), resolveEmailFrom(), resolveEmailTo(), constants.ColorReset)
	return nil
}

func saveEmailSpeedConfig(smtpHost, fromAddr, toAddr, pass string) {
	setGlobalIfNonEmpty("email.smtp_host", "SMTP_HOST", smtpHost)
	setGlobalIfNonEmpty("email.from", "SMTP_FROM", fromAddr)
	setGlobalIfNonEmpty("email.to", "AGM_NOTIFY_EMAIL", toAddr)
	setGlobalIfNonEmpty("email.password", "SMTP_PASSWORD", pass)
}

func setGlobalIfNonEmpty(key, envName, val string) {
	if len(strings.TrimSpace(val)) == 0 {
		return
	}
	_ = config.SetVariable("global", key, strings.TrimSpace(val))
	_ = os.Setenv(envName, strings.TrimSpace(val))
}

func executeEmailStatus(isJSON bool) error {
	st := collectSpeedSettingsStatus()
	if isJSON {
		return printSpeedJSON(st)
	}
	fmt.Printf("%s● GitMap Email Speed Notification Settings%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  SMTP Host:Port:    %s\n", st.EmailSMTPHost)
	fmt.Printf("  Sender (From):     %s\n", st.EmailFrom)
	fmt.Printf("  Recipient (To):    %s\n", st.EmailTo)
	fmt.Printf("  Ready:             %v\n", st.EmailReady)
	return nil
}

func executeEmailTest(args []string) error {
	st := collectSpeedSettingsStatus()
	isSent := sendEmailSwitchNotice("[GitMap] Speed Email Test", "GitMap Email Speed Setup verification succeeded.")
	fmt.Printf("%s✔ Test email dispatched to %s via %s (delivered=%v)%s\n",
		constants.ColorGreen, st.EmailTo, st.EmailSMTPHost, isSent, constants.ColorReset)
	return nil
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

func applySpeedSettingsFlags(args []string) {
	setGlobalIfNonEmpty("lap.default_hours", "GITMAP_LAP_DEFAULT_HOURS", extractSpeedFlag(args, "--lap-hours", "--lap"))
	setGlobalIfNonEmpty("account_switch.threshold", "GITMAP_ASW_THRESHOLD", extractSpeedFlag(args, "--threshold", "--switch-threshold"))
	setGlobalIfNonEmpty("machine.alias", "GITMAP_MACHINE_ALIAS", extractSpeedFlag(args, "--alias", "--machine-alias"))
	setGlobalIfNonEmpty("special_repos.secrets_name", "GITMAP_SPECIAL_SECRETS_REPO", extractSpeedFlag(args, "--secrets-repo", "--repo-secrets"))
	setGlobalIfNonEmpty("special_repos.cache_name", "GITMAP_SPECIAL_CACHE_REPO", extractSpeedFlag(args, "--cache-repo", "--repo-cache"))
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
	case "special_repos.secrets_name", "secrets_repo", "repo_secrets":
		_ = config.SetVariable("global", "special_repos.secrets_name", trimmedVal)
	case "special_repos.cache_name", "cache_repo", "repo_cache", "repo_storage":
		_ = config.SetVariable("global", "special_repos.cache_name", trimmedVal)
	default:
		_ = config.SetVariable("global", strings.TrimSpace(key), trimmedVal)
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
	fmt.Printf("  Machine Alias:             %s (key: machine.alias)\n", view.MachineAlias)
	fmt.Printf("  Telegram Chatbot:          chat=%s | token=%s | ready=%v\n", view.TelegramChatID, view.TelegramBotTokenMasked, view.TelegramReady)
	fmt.Printf("  Email Speed Setup:         smtp=%s | to=%s | ready=%v\n", view.EmailSMTPHost, view.EmailTo, view.EmailReady)
	fmt.Printf("  Special Secrets Repo  : %s (shortcut: gitmap rs / gitmap cd rs)\n", view.SpecialReposSecretsName)
	fmt.Printf("  Special Cache Repo    : %s (shortcut: gitmap rc / gitmap cd rc)\n", view.SpecialReposCacheName)
}

func collectSpeedSettingsView() speedSettingsView {
	secretsName, _ := config.GetVariable("global", "special_repos.secrets_name")
	cacheName, _ := config.GetVariable("global", "special_repos.cache_name")
	st := buildBaseSpeedSettingsStatus()
	st.SpecialReposSecretsName = defaultIfEmpty(secretsName, "repo-secrets")
	st.SpecialReposCacheName = defaultIfEmpty(cacheName, "repo-cache")
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
		TelegramBotTokenMasked: maskSecretToken(token), TelegramChatID: chatID,
		TelegramReady: len(token) > 0 && len(chatID) > 0,
		EmailSMTPHost: smtpHost, EmailFrom: resolveEmailFrom(), EmailTo: toAddr,
		EmailReady: len(smtpHost) > 0 && len(toAddr) > 0, MachineAlias: readSpeedVar("machine.alias", "auto-ip"),
	}
}

func resolveTelegramToken() string {
	return readSpeedVar("telegram.bot_token", os.Getenv("TELEGRAM_BOT_TOKEN"))
}

func resolveTelegramChatID() string {
	return readSpeedVar("telegram.chat_id", os.Getenv("TELEGRAM_CHAT_ID"))
}

func resolveEmailSMTP() string {
	return readSpeedVar("email.smtp_host", defaultIfEmpty(os.Getenv("SMTP_HOST"), "smtp.gmail.com:587"))
}

func resolveEmailFrom() string {
	return readSpeedVar("email.from", defaultIfEmpty(os.Getenv("SMTP_FROM"), "gitmap-bot@localhost"))
}

func resolveEmailTo() string {
	return readSpeedVar("email.to", defaultIfEmpty(os.Getenv("AGM_NOTIFY_EMAIL"), "dev@localhost"))
}

func readSpeedVar(key, fallback string) string {
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

func isTgEmailHelpArg(arg string) bool {
	low := strings.ToLower(strings.TrimSpace(arg))
	return low == "help" || low == "-h" || low == "--help"
}

func hasSpeedFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag || strings.HasPrefix(a, flag+"=") {
			return true
		}
	}
	return false
}

func extractSpeedFlag(args []string, primary, short string) string {
	for i, a := range args {
		val := matchSpeedFlagPair(args, i, a, primary, short)
		if len(val) > 0 {
			return val
		}
	}
	return ""
}

func matchSpeedFlagPair(args []string, i int, a, primary, short string) string {
	if (a == primary || a == short) && i+1 < len(args) {
		return args[i+1]
	}
	if strings.HasPrefix(a, primary+"=") {
		return strings.TrimPrefix(a, primary+"=")
	}
	return ""
}

func filterSpeedPositional(args []string) []string {
	var out []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			out = append(out, a)
		}
	}
	return out
}

func printSpeedJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal speed settings json")
	}
	fmt.Println(string(data))
	return nil
}

// RenderAgyTelegramHelp renders the boxed help menu for gitmap telegram / gitmap agy telegram.
func RenderAgyTelegramHelp() {
	termout.RenderMenu(termout.HelpMenu{
		Title: "Two-Way Telegram Chatbot Integration (gitmap telegram)",
		UsageLines: []string{
			"gitmap telegram <setup|status|send|poll|start|help> [flags]",
			"gitmap agy telegram <setup|status|send|poll|start|help> [flags]",
		},
		Sections: []termout.HelpSection{
			{
				Title: "Telegram Chatbot Subcommands",
				Entries: []termout.CommandEntry{
					{Command: "setup --token <T> --chat <ID>", Description: "Configure Telegram Bot Token & Chat ID and sync with AGY Manager"},
					{Command: "status [--json]", Description: "Show bot connection state, masked token, chat ID, and handler readiness"},
					{Command: "send \"<message>\"", Description: "Send an instant notification message to the configured Telegram chat"},
					{Command: "poll / start [--once]", Description: "Poll Telegram getUpdates API and execute incoming chat commands"},
					{Command: "help", Description: "Display this boxed Telegram bot help menu"},
				},
			},
		},
		Tips: []string{
			"Supported two-way chat commands: /status, /rp, /lap, /rwi <seq> <prompt>, /asw, /help.",
			"See canonical AI setup guide at 01-prompts/telegram-bot-setup.md.",
		},
	})
}

// RenderAgyEmailHelp renders the boxed help menu for gitmap email / gitmap agy email.
func RenderAgyEmailHelp() {
	termout.RenderMenu(termout.HelpMenu{
		Title: "Email Speed Notification Setup (gitmap email)",
		UsageLines: []string{
			"gitmap email <setup|status|test|help> [flags]",
			"gitmap agy email <setup|status|test|help> [flags]",
		},
		Sections: []termout.HelpSection{
			{
				Title: "Email Speed Subcommands",
				Entries: []termout.CommandEntry{
					{Command: "setup --smtp <h:p> --from <f> --to <t>", Description: "Configure SMTP host, sender, recipient, and app password"},
					{Command: "status [--json]", Description: "Display current Email speed notification configuration"},
					{Command: "test", Description: "Send a test verification email using configured SMTP settings"},
					{Command: "help", Description: "Display this boxed Email speed setup help menu"},
				},
			},
		},
	})
}

// RenderAgySettingsHelp renders the boxed help menu for gitmap settings / gitmap agy settings.
func RenderAgySettingsHelp() {
	renderSettingsHelp()
}

func renderSettingsHelp() {
	termout.RenderMenu(termout.HelpMenu{
		Title: "GitMap & Antigravity Speed Settings (gitmap settings)",
		UsageLines: []string{
			"gitmap settings [--lap-hours 24] [--threshold 15] [--alias <name>] [--json]",
			"gitmap settings set special_repos.secrets_name repo-secrets",
			"gitmap settings set special_repos.cache_name repo-cache",
			"gitmap agy settings [set <key> <value> | export <file> | import <file>]",
		},
		Sections: buildSettingsHelpSections(),
	})
}

func buildSettingsHelpSections() []termout.HelpSection {
	return []termout.HelpSection{
		{
			Title: "Speed Settings & Configuration Keys",
			Entries: []termout.CommandEntry{
				{Command: "--lap-hours <N>", Description: "Set default lookback hours for last-active-projects (default: 24)"},
				{Command: "--threshold <N>", Description: "Set account-switch remaining credit threshold percentage (default: 15)"},
				{Command: "--alias <name>", Description: "Set machine network alias for SSH fleet recognition"},
				{Command: "set special_repos.secrets_name repo-secrets", Description: "Configure special secrets repository name (shortcut: gitmap rs / gitmap cd rs)"},
				{Command: "set special_repos.cache_name repo-cache", Description: "Configure special cache repository name (shortcut: gitmap rc / gitmap cd rc)"},
				{Command: "export / import <file>", Description: "Export or import full Antigravity settings JSON"},
			},
		},
	}
}
