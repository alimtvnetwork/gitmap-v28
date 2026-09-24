package cmdagy

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	newConvModelFlag   string
	newConvTitleFlag   string
	newConvProfileFlag string
	newConvFileFlag    string

	sendMsgTitleFlag string
	sendMsgFirstFlag bool
	sendMsgFileFlag  string
)

var agyNewConvCmd = &cobra.Command{
	Use:     "new-conversation [flags] [prompt]",
	Aliases: []string{"new-conv", "nc"},
	Short:   "Create a new conversation session in Antigravity IDE via agentapi",
	Long: `Create a new conversation session in Antigravity IDE using agentapi.

Examples:
  gitmap agy new-conversation "Review recent git commits"
  gitmap agy new-conv --model=pro --title="Refactor" "Please refactor the db queries"
  gitmap agy nc -f prompt.txt
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		appErr := runAgyNewConversation(args)
		if appErr != nil {
			return appErr
		}
		return nil
	},
}

var agySendMessageCmd = &cobra.Command{
	Use:     "send-message [flags] [recipient_id] [content]",
	Aliases: []string{"send", "msg"},
	Short:   "Send a message to an active Antigravity session via agentapi",
	Long: `Send a prompt or message into an existing Antigravity session via agentapi.

Targeting:
  If recipient_id is omitted or specified as "first" / "latest", or --first is passed,
  the command automatically targets the most recently active conversation.

Examples:
  gitmap agy send-message 2bc23757-bdea-41a1-9b13-65d61f7365e1 "Continue with task 2"
  gitmap agy send-message first "Inspect the failing test logs"
  gitmap agy send --first "Check pipeline status"
  gitmap agy msg "Quick update on build"
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		appErr := runAgySendMessage(args)
		if appErr != nil {
			return appErr
		}
		return nil
	},
}

func init() {
	agyNewConvCmd.Flags().StringVarP(&newConvModelFlag, "model", "m", "", "Model type (flash_lite, flash, pro)")
	agyNewConvCmd.Flags().StringVarP(&newConvTitleFlag, "title", "t", "", "Conversation title")
	agyNewConvCmd.Flags().StringVar(&newConvProfileFlag, "profile", "", "Profile to use for session")
	agyNewConvCmd.Flags().StringVarP(&newConvFileFlag, "file", "f", "", "Read prompt from file")

	agySendMessageCmd.Flags().StringVarP(&sendMsgTitleFlag, "title", "t", "", "Message title")
	agySendMessageCmd.Flags().BoolVarP(&sendMsgFirstFlag, "first", "1", false, "Target the first/most recent active conversation")
	agySendMessageCmd.Flags().StringVarP(&sendMsgFileFlag, "file", "f", "", "Read message content from file")
}

func runAgyNewConversation(args []string) *apperror.AppError {
	prompt, err := resolvePromptContent(args, newConvFileFlag)
	if err != nil {
		return apperror.WrapSimple(err, "resolve prompt content")
	}

	res := AgentAPINewConversationWithOptions(newConvTitleFlag, newConvModelFlag, newConvProfileFlag, prompt)
	if res.IsFailure() {
		return res.Err
	}

	renderNewConvSuccess(res.Value, newConvTitleFlag, newConvModelFlag)
	return nil
}

func runAgySendMessage(args []string) *apperror.AppError {
	recipientID, content, err := parseSendMessageArgs(args, sendMsgFileFlag, sendMsgFirstFlag)
	if err != nil {
		return apperror.WrapSimple(err, "parse send-message arguments")
	}

	res := AgentAPISendMessage(recipientID, sendMsgTitleFlag, content)
	if res.IsFailure() {
		return res.Err
	}

	renderSendMessageSuccess(recipientID, sendMsgTitleFlag, content)
	return nil
}

func resolvePromptContent(args []string, filePath string) (string, error) {
	if filePath != "" {
		return readPromptFromFile(filePath)
	}

	if len(args) == 0 {
		return "", apperror.NewSimple("prompt content is required (pass as argument or via --file)", "E9035")
	}

	return strings.TrimSpace(strings.Join(args, " ")), nil
}

func readPromptFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", apperror.WrapSimple(err, "read prompt file")
	}
	return strings.TrimSpace(string(data)), nil
}

func parseSendMessageArgs(args []string, filePath string, isFirst bool) (string, string, error) {
	fileContent, hasFileContent := loadContentFromFile(filePath)
	if isFirst || len(args) == 1 {
		return resolveSingleOrFirstTarget(args, fileContent, hasFileContent)
	}

	if len(args) >= 2 {
		recipient := args[0]
		content := selectContent(fileContent, hasFileContent, args[1:])
		targetID, err := resolveRecipientID(recipient)
		return targetID, content, err
	}

	if hasFileContent {
		targetID, err := getLatestActiveConversationID()
		return targetID, fileContent, err
	}

	return "", "", apperror.NewSimple("recipient ID and message content are required", "E9036")
}

func selectContent(fileContent string, hasFileContent bool, tokens []string) string {
	if hasFileContent {
		return fileContent
	}
	return strings.TrimSpace(strings.Join(tokens, " "))
}

func loadContentFromFile(filePath string) (string, bool) {
	if filePath == "" {
		return "", false
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(data)), true
}

func resolveSingleOrFirstTarget(args []string, fileContent string, hasFileContent bool) (string, string, error) {
	targetID, err := getLatestActiveConversationID()
	if err != nil {
		return "", "", err
	}

	if hasFileContent {
		return targetID, fileContent, nil
	}

	if len(args) == 0 {
		return "", "", apperror.NewSimple("content is required when targeting first conversation", "E9037")
	}

	return targetID, resolveArgContent(args), nil
}

func resolveArgContent(args []string) string {
	firstTok := strings.ToLower(strings.TrimSpace(args[0]))
	if firstTok == "first" || firstTok == "latest" {
		return "Verify task completion."
	}
	return strings.TrimSpace(strings.Join(args, " "))
}

func resolveRecipientID(token string) (string, error) {
	trimmed := strings.TrimSpace(token)
	lower := strings.ToLower(trimmed)
	if lower == "first" || lower == "latest" || lower == "" {
		return getLatestActiveConversationID()
	}
	return trimmed, nil
}

func getLatestActiveConversationID() (string, error) {
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return "", apperror.WrapSimple(err, "conversation summaries db path")
	}

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return "", apperror.WrapSimple(openErr, "open conversation summaries db")
	}
	defer conn.Close()

	query := "SELECT conversation_id FROM conversation_summaries WHERE (killed IS NULL OR killed = 0) ORDER BY last_modified_time DESC LIMIT 1"
	row := conn.QueryRow(query)
	var convID string
	if scanErr := row.Scan(&convID); scanErr != nil {
		return "", apperror.NewSimple("no active conversation found in conversation_summaries.db", "E9038")
	}

	return convID, nil
}

func renderNewConvSuccess(convID, title, model string) {
	successMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Render("✔")
	fmt.Printf("\n  %s Created new Antigravity session via agentapi!\n", successMark)
	fmt.Printf("    • Conversation ID: %s%s%s\n", constants.ColorCyan, convID, constants.ColorReset)
	if title != "" {
		fmt.Printf("    • Title:           %s\n", title)
	}
	if model != "" {
		fmt.Printf("    • Model:           %s\n", model)
	}
	fmt.Println()
}

func renderSendMessageSuccess(recipientID, title, content string) {
	successMark := lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Render("✔")
	fmt.Printf("\n  %s Sent message to Antigravity session via agentapi!\n", successMark)
	fmt.Printf("    • Recipient ID:    %s%s%s\n", constants.ColorCyan, recipientID, constants.ColorReset)
	if title != "" {
		fmt.Printf("    • Title:           %s\n", title)
	}
	preview := content
	if len(preview) > 80 {
		preview = preview[:77] + "..."
	}
	fmt.Printf("    • Content Preview: %s\n\n", strings.ReplaceAll(preview, "\n", " "))
}
