// Package cmdagy — agy_add_read.go implements `gitmap agy add [path|.]` and `gitmap agy add-read [path|.]`.
package cmdagy

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

// AgyAddResult represents the structured result of registering a project in Antigravity.
type AgyAddResult struct {
	ProjectID    string `json:"projectId"`
	ProjectName  string `json:"projectName"`
	ProjectAlias string `json:"projectAlias"`
	ProjectPath  string `json:"projectPath"`
	IsReadQueued bool   `json:"isReadQueued"`
	IsDryRun     bool   `json:"isDryRun"`
	QueueFile    string `json:"queueFile,omitempty"`
}

var (
	agyAddDryRun bool
	agyAddJSON   bool
)

// AgyAddCmd registers a repository path (or current directory '.') in Antigravity.
var AgyAddCmd = &cobra.Command{
	Use:     "add [path|.]",
	Aliases: []string{"add-project"},
	Short:   "Register current directory (.) or path in Antigravity workspace registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyAddPathOrLegacy(args)
	},
}

// AgyAddReadCmd registers a repository path and immediately injects the Read Memory prompt.
var AgyAddReadCmd = &cobra.Command{
	Use:     "add-read [path|.]",
	Aliases: []string{"ar", "add-and-read"},
	Short:   "Register project in Antigravity and immediately dispatch Read Memory prompt",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyAddRead(args)
	},
}

func init() {
	AgyAddCmd.Flags().BoolVarP(&agyAddDryRun, "dry-run", "d", false, "Preview registration without modifying disk")
	AgyAddCmd.Flags().BoolVar(&agyAddJSON, "json", false, "Output result in JSON format")
	AgyAddReadCmd.Flags().BoolVarP(&agyAddDryRun, "dry-run", "d", false, "Preview registration and read prompt without modifying disk")
	AgyAddReadCmd.Flags().BoolVar(&agyAddJSON, "json", false, "Output result in JSON format")
	agyAddCmd.RunE = AgyAddCmd.RunE
}

func runAgyAddPathOrLegacy(args []string) error {
	if len(args) > 0 && strings.EqualFold(args[0], "read") {
		return runAgyAddRead(args[1:])
	}
	if isLegacyAddIdAndName(args) {
		return createProjectFile(args[0], args[1])
	}
	return executeAgyAddFlow(args, false)
}

func runAgyAddRead(args []string) error {
	return executeAgyAddFlow(args, true)
}

func isLegacyAddIdAndName(args []string) bool {
	if len(args) != 2 {
		return false
	}
	if args[0] == "." || args[0] == ".." || strings.ContainsAny(args[0], `/\`) {
		return false
	}
	return !checkDirExists(args[0])
}

func executeAgyAddFlow(args []string, isReadMode bool) error {
	targetPath, err := resolveAddTargetRepoPath(args)
	if err != nil {
		return err
	}
	res, regErr := registerAndOptionallyRead(targetPath, isReadMode)
	if regErr != nil {
		return regErr
	}
	return renderAgyAddResult(res)
}

func resolveAddTargetRepoPath(args []string) (string, error) {
	raw := "."
	if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
		raw = strings.TrimSpace(args[0])
	}
	if raw == "." {
		return resolveCurrentWorkingGitRoot()
	}
	abs, err := filepath.Abs(raw)
	if err != nil {
		return "", apperror.WrapSimple(err, "resolve target path")
	}
	return findEnclosingGitRoot(abs), nil
}

func resolveCurrentWorkingGitRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", apperror.WrapSimple(err, "resolve current directory")
	}
	return findEnclosingGitRoot(cwd), nil
}

func findEnclosingGitRoot(startDir string) string {
	curr := filepath.Clean(startDir)
	for i := 0; i < 12; i++ {
		if checkDirExists(filepath.Join(curr, ".git")) {
			return curr
		}
		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}
	return filepath.Clean(startDir)
}

func registerAndOptionallyRead(projectPath string, isReadMode bool) (*AgyAddResult, error) {
	name := filepath.Base(filepath.Clean(projectPath))
	alias := strings.ToLower(name)
	projID := deriveOrLookupProjectID(projectPath)
	res := &AgyAddResult{
		ProjectID:    projID,
		ProjectName:  name,
		ProjectAlias: alias,
		ProjectPath:  projectPath,
		IsReadQueued: isReadMode,
		IsDryRun:     agyAddDryRun,
	}
	if agyAddDryRun {
		return res, nil
	}
	return performLiveAddAndRead(res, isReadMode)
}

func performLiveAddAndRead(res *AgyAddResult, isReadMode bool) (*AgyAddResult, error) {
	if err := persistProjectInAgyRegistry(res); err != nil {
		return nil, err
	}
	if !isReadMode {
		return res, nil
	}
	queueFile, err := enqueueAndDispatchReadMemory(res.ProjectPath, res.ProjectName)
	if err != nil {
		return nil, err
	}
	res.QueueFile = queueFile
	return res, nil
}

func deriveOrLookupProjectID(projectPath string) string {
	if id := discoverFreshProjectId(projectPath); id != "" {
		return id
	}
	sum := sha1.Sum([]byte(strings.ToLower(filepath.Clean(projectPath))))
	hexStr := hex.EncodeToString(sum[:])
	return fmt.Sprintf("%s-%s-%s-%s-%s", hexStr[0:8], hexStr[8:12], hexStr[12:16], hexStr[16:20], hexStr[20:32])
}

func persistProjectInAgyRegistry(res *AgyAddResult) error {
	configDir, err := getProjectsDirPath()
	if err != nil {
		return apperror.WrapSimple(err, "resolve projects dir")
	}
	_ = os.MkdirAll(configDir, 0755)
	_ = workspacesync.SyncAntigravity(res.ProjectPath, res.ProjectName)
	if syncedID := discoverFreshProjectId(res.ProjectPath); syncedID != "" {
		res.ProjectID = syncedID
		return nil
	}
	return writeFullAgyProjectEntry(configDir, res)
}

func writeFullAgyProjectEntry(configDir string, res *AgyAddResult) error {
	proj := AgyProject{
		ID:        res.ProjectID,
		Name:      res.ProjectName,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339Nano),
		ProjectResources: &AgyProjectResources{
			Resources: []AgyResource{
				{GitFolder: &AgyGitFolder{FolderURI: buildFolderURI(res.ProjectPath), DefaultBranch: "main"}},
			},
		},
	}
	data, err := json.MarshalIndent(proj, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal project json")
	}
	filePath := filepath.Join(configDir, res.ProjectID+".json")
	return os.WriteFile(filePath, data, 0644)
}

func enqueueAndDispatchReadMemory(projectPath, projectName string) (string, error) {
	promptText := loadReadMemoryPromptText(projectPath)
	title := fmt.Sprintf("Read Memory: %s", projectName)
	qFile, err := writeReadPromptToProjectQueue(projectPath, title, promptText)
	if err != nil {
		return "", err
	}
	_ = InjectAgyPrompt(projectPath, promptText, title, "add_read", false)
	return qFile, nil
}

func loadReadMemoryPromptText(projectPath string) string {
	candidates := []string{
		filepath.Join(projectPath, "01-prompts", "read.md"),
		filepath.Join(projectPath, ".ai-memory", "prompts", "read.md"),
		filepath.Join("01-prompts", "read.md"),
	}
	for _, p := range candidates {
		data, err := os.ReadFile(p)
		if err == nil && len(strings.TrimSpace(string(data))) > 0 {
			return strings.TrimSpace(string(data))
		}
	}
	return defaultReadMemoryPrompt
}

func writeReadPromptToProjectQueue(projectPath, title, promptText string) (string, error) {
	qDir := filepath.Join(projectPath, ".ai-memory", "temp")
	if err := os.MkdirAll(qDir, 0755); err != nil {
		return "", apperror.WrapSimple(err, "mkdir .ai-memory/temp")
	}
	qFile := filepath.Join(qDir, "agy-prompt-queue.json")
	qData := buildInitialQueueFilePayload(title, promptText)
	raw, err := json.MarshalIndent(qData, "", "  ")
	if err != nil {
		return "", apperror.WrapSimple(err, "marshal prompt queue")
	}
	if writeErr := os.WriteFile(qFile, raw, 0644); writeErr != nil {
		return "", apperror.WrapSimple(writeErr, "write prompt queue")
	}
	return qFile, nil
}

func buildInitialQueueFilePayload(title, promptText string) AgyPromptQueueFile {
	now := time.Now().UTC().Format(time.RFC3339)
	entry := AgyPromptQueueEntry{
		ID:        1,
		Title:     title,
		Prompt:    promptText,
		Type:      "read_memory",
		Status:    "queued",
		CreatedAt: now,
	}
	return AgyPromptQueueFile{
		Active:    &entry,
		Queued:    []AgyPromptQueueEntry{entry},
		UpdatedAt: now,
	}
}

func renderAgyAddResult(res *AgyAddResult) error {
	if agyAddJSON {
		return printAgyAddJSON(res)
	}
	printAgyAddBox(res)
	return nil
}

func printAgyAddJSON(res *AgyAddResult) error {
	data, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal add result")
	}
	fmt.Println(string(data))
	return nil
}

func printAgyAddBox(res *AgyAddResult) {
	modeLabel := "ADD PROJECT"
	if res.IsReadQueued {
		modeLabel = "ADD & READ MEMORY"
	}
	fmt.Printf("\n  %s╔══ ANTIGRAVITY %s ══════════════════════════════════════╗%s\n", constants.ColorCyan, modeLabel, constants.ColorReset)
	fmt.Printf("  │ Project Name:  %-52s │\n", res.ProjectName)
	fmt.Printf("  │ Project ID:    %-52s │\n", res.ProjectID)
	fmt.Printf("  │ Project Alias: %-52s │\n", res.ProjectAlias)
	fmt.Printf("  │ Project Path:  %-52s │\n", res.ProjectPath)
	if res.IsReadQueued {
		fmt.Printf("  │ Read Memory:   %-52s │\n", "Dispatched & Enqueued (01-prompts/read.md)")
	}
	fmt.Printf("  %s╚═════════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorCyan, constants.ColorReset)
}
