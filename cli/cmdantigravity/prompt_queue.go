package cmdantigravity

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

// DefaultReadMemoryPrompt is the canonical "read all first" prompt for newly registered projects.
const DefaultReadMemoryPrompt = "Execute enhanced Read Memory protocol. Defensively load memory, specs, constraints, and pending plans before taking action."

// PromptQueueFile models the JSON storage for Antigravity prompt queues.
type PromptQueueFile struct {
	Active    *PromptQueueEntry  `json:"active,omitempty"`
	Queued    []PromptQueueEntry `json:"queued,omitempty"`
	UpdatedAt string             `json:"updatedAt"`
}

// PromptQueueEntry records a single prompt inside the queue.
type PromptQueueEntry struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Prompt    string `json:"prompt"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

type rawProjectConfig struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	ProjectResources *rawProjectResources `json:"projectResources,omitempty"`
}

type rawProjectResources struct {
	Resources []rawResource `json:"resources,omitempty"`
}

type rawResource struct {
	GitFolder *rawGitFolder `json:"gitFolder,omitempty"`
}

type rawGitFolder struct {
	FolderURI     string `json:"folderUri,omitempty"`
	DefaultBranch string `json:"defaultBranch,omitempty"`
}

// AssemblePrompt combines template content and additional text with prefix or suffix ordering.
func AssemblePrompt(templateContent, text string, isSuffixMode bool) string {
	tpl := strings.TrimSpace(templateContent)
	txt := strings.TrimSpace(text)
	hasTpl := len(tpl) > 0
	hasTxt := len(txt) > 0
	if hasTpl && hasTxt {
		return formatCombinedPrompt(tpl, txt, isSuffixMode)
	}
	if hasTpl {
		return tpl
	}

	return txt
}

func formatCombinedPrompt(tpl, txt string, isSuffixMode bool) string {
	if isSuffixMode {
		return txt + "\n\n" + tpl
	}

	return tpl + "\n\n" + txt
}

// IsRepoRegisteredInAgy checks whether repoRoot is configured in ~/.gemini/config/projects/.
func IsRepoRegisteredInAgy(repoRoot string) bool {
	configDir := resolveAgyProjectsConfigDir()
	entries, err := os.ReadDir(configDir)
	hasErr := err != nil
	if hasErr {
		return false
	}
	cleanRepo := filepath.ToSlash(filepath.Clean(repoRoot))

	return matchRepoInProjectEntries(configDir, entries, cleanRepo)
}

func resolveAgyProjectsConfigDir() string {
	home, err := os.UserHomeDir()
	hasHome := err == nil
	if hasHome {
		return filepath.Join(home, ".gemini", "config", "projects")
	}

	return ""
}

func matchRepoInProjectEntries(dir string, entries []os.DirEntry, cleanRepo string) bool {
	for _, entry := range entries {
		isJson := filepath.Ext(entry.Name()) == ".json"
		if isJson && matchProjectFileURI(filepath.Join(dir, entry.Name()), cleanRepo) {
			return true
		}
	}

	return false
}

func matchProjectFileURI(filePath, cleanRepo string) bool {
	data, err := os.ReadFile(filePath)
	hasData := err == nil
	if hasData == false {
		return false
	}
	var cfg rawProjectConfig
	hasCfg := json.Unmarshal(data, &cfg) == nil
	if hasCfg {
		return checkResourceFolderURI(cfg.ProjectResources, cleanRepo)
	}

	return false
}

func checkResourceFolderURI(res *rawProjectResources, cleanRepo string) bool {
	hasRes := res != nil && len(res.Resources) > 0
	if hasRes == false {
		return false
	}
	for _, r := range res.Resources {
		isMatch := r.GitFolder != nil && isFolderURIMatch(r.GitFolder.FolderURI, cleanRepo)
		if isMatch {
			return true
		}
	}

	return false
}

func isFolderURIMatch(rawURI, cleanRepo string) bool {
	parsed := parseFolderURI(rawURI)

	return strings.EqualFold(parsed, cleanRepo)
}

func parseFolderURI(rawURI string) string {
	decoded, err := url.PathUnescape(rawURI)
	hasDecoded := err == nil
	if hasDecoded == false {
		decoded = rawURI
	}
	clean := strings.TrimPrefix(decoded, "file:///")
	clean = strings.TrimPrefix(clean, "file://")

	return filepath.ToSlash(filepath.Clean(clean))
}

// AutoRegisterRepoInAgy registers repoRoot in Antigravity via workspacesync.
func AutoRegisterRepoInAgy(repoRoot string) error {
	baseName := filepath.Base(repoRoot)
	isSynced := workspacesync.SyncAntigravity(repoRoot, baseName)
	if isSynced {
		return nil
	}

	return apperror.NewSimple("failed to auto-register repository in Antigravity: "+repoRoot, "E9014")
}

// EnqueueWithDualQueuePolicy handles auto-registration and dual-queue or single-queue prompt execution.
func EnqueueWithDualQueuePolicy(repoRoot, templateName, assembledPrompt string) error {
	isRegistered := IsRepoRegisteredInAgy(repoRoot)
	if isRegistered {
		return processRegisteredRepoPrompt(repoRoot, templateName, assembledPrompt)
	}

	return processUnregisteredRepoPrompt(repoRoot, templateName, assembledPrompt)
}

func processUnregisteredRepoPrompt(repoRoot, templateName, assembledPrompt string) error {
	regErr := AutoRegisterRepoInAgy(repoRoot)
	if regErr != nil {
		return regErr
	}
	renderAutoRegisterNotice(repoRoot)
	dualErr := enqueueDualPrompts(repoRoot, templateName, assembledPrompt)
	if dualErr != nil {
		return dualErr
	}
	stageAndCopyActivePrompt(repoRoot, assembledPrompt)

	return nil
}

func processRegisteredRepoPrompt(repoRoot, templateName, assembledPrompt string) error {
	title := resolveUserPromptTitle(templateName)
	qErr := enqueueSinglePrompt(repoRoot, "user_prompt", title, assembledPrompt)
	if qErr != nil {
		return qErr
	}
	renderSingleQueueNotice(repoRoot)
	stageAndCopyActivePrompt(repoRoot, assembledPrompt)

	return nil
}

func renderAutoRegisterNotice(repoRoot string) {
	fmt.Printf("  %s✔ Auto-registered repository '%s' in Antigravity%s\n",
		constants.ColorGreen, filepath.Base(repoRoot), constants.ColorReset)
	fmt.Printf("  %s✔ Dual-queued: [1] Read Memory protocol, [2] User Prompt%s\n",
		constants.ColorGreen, constants.ColorReset)
}

func renderSingleQueueNotice(repoRoot string) {
	fmt.Printf("  %s✔ Enqueued prompt for repository '%s'%s\n",
		constants.ColorGreen, filepath.Base(repoRoot), constants.ColorReset)
}

func resolveUserPromptTitle(templateName string) string {
	hasName := len(strings.TrimSpace(templateName)) > 0
	if hasName {
		return "Prompt: " + templateName
	}

	return "User Prompt"
}

func resolvePromptQueuePath(repoRoot string) string {
	return filepath.Join(repoRoot, ".ai-memory", "temp", "agy-prompt-queue.json")
}

func resolveActivePromptPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".ai-memory", "temp", "active-agy-pipeline-fix-prompt.txt")
}

func loadPromptQueueFile(queuePath string) PromptQueueFile {
	data, err := os.ReadFile(queuePath)
	hasData := err == nil
	if hasData == false {
		return PromptQueueFile{Queued: make([]PromptQueueEntry, 0)}
	}
	var q PromptQueueFile
	hasParsed := json.Unmarshal(data, &q) == nil
	if hasParsed {
		return q
	}

	return PromptQueueFile{Queued: make([]PromptQueueEntry, 0)}
}

func savePromptQueueFile(queuePath string, q PromptQueueFile) error {
	_ = os.MkdirAll(filepath.Dir(queuePath), 0755)
	q.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	data, err := json.MarshalIndent(q, "", "  ")
	hasErr := err != nil
	if hasErr {
		return apperror.WrapSimple(err, "marshal prompt queue")
	}
	writeErr := os.WriteFile(queuePath, data, 0644)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "write prompt queue")
	}

	return nil
}

func enqueueDualPrompts(repoRoot, templateName, assembledPrompt string) error {
	queuePath := resolvePromptQueuePath(repoRoot)
	q := loadPromptQueueFile(queuePath)
	now := time.Now().UTC().Format(time.RFC3339)
	entry1 := makeQueueEntry(computeNextID(q), "read_memory", "Read Memory Protocol", DefaultReadMemoryPrompt, now)
	entry2 := makeQueueEntry(entry1.ID+1, "user_prompt", resolveUserPromptTitle(templateName), assembledPrompt, now)
	q = appendDualEntries(q, entry1, entry2)

	return savePromptQueueFile(queuePath, q)
}

func enqueueSinglePrompt(repoRoot, pType, title, promptText string) error {
	queuePath := resolvePromptQueuePath(repoRoot)
	q := loadPromptQueueFile(queuePath)
	now := time.Now().UTC().Format(time.RFC3339)
	entry := makeQueueEntry(computeNextID(q), pType, title, promptText, now)
	q = appendSingleEntry(q, entry)

	return savePromptQueueFile(queuePath, q)
}

func makeQueueEntry(id int, pType, title, text, now string) PromptQueueEntry {
	return PromptQueueEntry{
		ID:        id,
		Type:      pType,
		Title:     title,
		Prompt:    text,
		Status:    "queued",
		CreatedAt: now,
	}
}

func appendDualEntries(q PromptQueueFile, e1, e2 PromptQueueEntry) PromptQueueFile {
	hasActive := q.Active != nil
	if hasActive {
		q.Queued = append(q.Queued, e1, e2)

		return q
	}
	q.Active = &e1
	q.Queued = append(q.Queued, e2)

	return q
}

func appendSingleEntry(q PromptQueueFile, e PromptQueueEntry) PromptQueueFile {
	hasActive := q.Active != nil
	if hasActive {
		q.Queued = append(q.Queued, e)

		return q
	}
	q.Active = &e

	return q
}

func computeNextID(q PromptQueueFile) int {
	nextID := 1
	hasActive := q.Active != nil && q.Active.ID >= nextID
	if hasActive {
		nextID = q.Active.ID + 1
	}
	for _, item := range q.Queued {
		hasHigher := item.ID >= nextID
		if hasHigher {
			nextID = item.ID + 1
		}
	}

	return nextID
}

func stageAndCopyActivePrompt(repoRoot, content string) {
	stagePromptOnDisk(repoRoot, content)
	copyPromptToClipboard(content)
}

func stagePromptOnDisk(repoRoot, content string) {
	promptPath := resolveActivePromptPath(repoRoot)
	_ = os.MkdirAll(filepath.Dir(promptPath), 0755)
	_ = os.WriteFile(promptPath, []byte(content), 0644)
}

func copyPromptToClipboard(content string) {
	hasContent := len(content) > 0
	if hasContent {
		_ = clipboard.WriteAll(content)
	}
}
