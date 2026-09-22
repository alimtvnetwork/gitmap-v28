package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// InjectAgyFixTask attempts dispatch of the fix prompt to active IDE or Antigravity CLI.
func InjectAgyFixTask(repoDir, absPayloadPath string, isSkipInject bool) AgyInjectionResult {
	if isSkipInject {
		return makeSkipInjectionResult()
	}

	repoRoot := resolveTargetRepoRoot(repoDir)
	promptPath := stageActivePromptInRepo(repoRoot, absPayloadPath)
	ideProcRes := DetectRunningAntigravityIDE()
	pid := resolveActiveOrZeroPID(ideProcRes)

	return handleIDEInjection(repoRoot, promptPath, pid)
}

func resolveActiveOrZeroPID(ideProcRes result.Result[AgyProcessInfo]) int {
	if ideProcRes.IsSuccess() {
		return ideProcRes.Value.PID
	}

	return 0
}

func handleIDEInjection(repoRoot, promptPath string, pid int) AgyInjectionResult {
	promptContent := readPromptContentOrDefault(promptPath)
	dispatchRes := DispatchPromptToAntigravity(repoRoot, promptPath, "Fix CI/CD Pipeline Errors with 4-Part RCA", promptContent, pid)
	if dispatchRes.IsSuccess {
		return dispatchRes
	}

	convStatus := DetectConversationExecutionStatus(repoRoot)
	isRunning := convStatus == AgyConvStatusRunning
	if isRunning {
		_, _ = EnqueuePrompt("pipeline_fix", "Fix CI/CD Pipeline Errors with 4-Part RCA", promptContent)

		return makeQueuedSuccessResult(pid, repoRoot, promptPath)
	}

	return dispatchRes
}

// InjectAgyPrompt injects any prompt into active Antigravity IDE or CLI using queue protocol.
func InjectAgyPrompt(repoDir, promptText, title, promptType string, isSkipInject bool) AgyInjectionResult {
	if isSkipInject {
		return makeSkipInjectionResult()
	}

	repoRoot := resolveTargetRepoRoot(repoDir)
	targetPrompt := filepath.Join(repoRoot, activeAgyPromptRelativePath)
	writePromptFile(targetPrompt, promptText)
	ideProcRes := DetectRunningAntigravityIDE()
	pid := resolveActiveOrZeroPID(ideProcRes)

	return handleIDEInjectionWithPayload(repoRoot, targetPrompt, promptText, title, promptType, pid)
}

func handleIDEInjectionWithPayload(repoRoot, promptPath, promptText, title, promptType string, pid int) AgyInjectionResult {
	dispatchRes := DispatchPromptToAntigravity(repoRoot, promptPath, title, promptText, pid)
	if dispatchRes.IsSuccess {
		return dispatchRes
	}

	convStatus := DetectConversationExecutionStatus(repoRoot)
	isRunning := convStatus == AgyConvStatusRunning
	if isRunning {
		_, _ = EnqueuePrompt(promptType, title, promptText)

		return makeQueuedSuccessResult(pid, repoRoot, promptPath)
	}

	return dispatchRes
}

func makeQueuedSuccessResult(pid int, repoDir, promptPath string) AgyInjectionResult {
	msg := formatQueuedSuccessMessage(pid)

	return AgyInjectionResult{
		IsSuccess:  true,
		Mode:       AgyInjectionModeQueued,
		PID:        pid,
		Message:    msg,
		PromptPath: promptPath,
		RepoDir:    repoDir,
	}
}

func formatQueuedSuccessMessage(pid int) string {
	hasPID := pid > 0
	if hasPID {
		return fmt.Sprintf("Active Antigravity IDE is currently busy (RUNNING); prompt queued in agy-prompt-queue.json (PID: %d)", pid)
	}

	return "Antigravity IDE offline/ready; prompt queued in agy-prompt-queue.json"
}

func resolveTargetRepoRoot(repoDir string) string {
	startPath := resolveInitialPath(repoDir)
	root, err := gitutil.RepoRoot(startPath)
	hasRoot := err == nil && len(root) > 0
	if hasRoot {
		return root
	}

	return toAbsPath(startPath)
}

func resolveInitialPath(repoDir string) string {
	hasDir := len(repoDir) > 0
	if hasDir {
		return repoDir
	}

	return resolveProjectRootDir()
}

func stageActivePromptInRepo(repoRoot, absPayloadPath string) string {
	targetPrompt := filepath.Join(repoRoot, activeAgyPromptRelativePath)
	isSame := isSamePath(targetPrompt, absPayloadPath)
	if isSame {
		return targetPrompt
	}

	content := readPromptContentOrDefault(absPayloadPath)
	writePromptFile(targetPrompt, content)

	return targetPrompt
}

func isSamePath(pathA, pathB string) bool {
	cleanA := filepath.Clean(pathA)
	cleanB := filepath.Clean(pathB)

	return strings.EqualFold(cleanA, cleanB)
}

func readPromptContentOrDefault(path string) string {
	hasEmpty := len(path) == 0
	if hasEmpty {
		return ""
	}

	data, err := os.ReadFile(path)
	hasErr := err != nil
	if hasErr {
		return ""
	}

	return string(data)
}

func writePromptFile(path, content string) {
	hasEmpty := len(content) == 0
	if hasEmpty {
		return
	}

	_ = os.MkdirAll(filepath.Dir(path), 0755)
	_ = os.WriteFile(path, []byte(content), 0644)
}

func makeSkipInjectionResult() AgyInjectionResult {
	return AgyInjectionResult{
		IsSuccess: false,
		Mode:      AgyInjectionModeNone,
		Message:   "direct injection skipped by flag (--no-inject)",
	}
}

// DetectConversationExecutionStatus detects if the conversation for repoRoot is IDLE or RUNNING.
func DetectConversationExecutionStatus(repoRoot string) AgyConvStatusType {
	state := DetectConversationExecutionState(repoRoot)

	return state.Status
}

// DetectConversationExecutionState inspects transcript.jsonl for the active conversation.
func DetectConversationExecutionState(repoRoot string) AgyConversationExecutionState {
	convID, transcriptPath := findActiveConvTranscript(repoRoot)
	hasTranscript := len(transcriptPath) > 0
	if hasTranscript == false {
		return AgyConversationExecutionState{Status: AgyConvStatusUnknown}
	}

	status := detectTranscriptStatus(transcriptPath)

	return AgyConversationExecutionState{
		ConvID:         convID,
		Status:         status,
		TranscriptPath: transcriptPath,
	}
}

func resolveExistingConvTranscript(convID string) (string, bool) {
	hasConvID := len(convID) > 0
	if hasConvID == false {
		return "", false
	}
	path := resolveTranscriptPathForConv(convID)

	return path, checkFileExists(path)
}

func findActiveConvTranscript(repoRoot string) (string, string) {
	convID := findMatchingActiveConvID(repoRoot)
	path, isFound := resolveExistingConvTranscript(convID)
	if isFound {
		return convID, path
	}

	return findLatestModifiedTranscript()
}

func resolveTranscriptPathForConv(convID string) string {
	brainDir, err := GetBrainLogsDirPath()
	hasErr := err != nil
	if hasErr {
		return ""
	}

	fullPath := filepath.Join(brainDir, convID, ".system_generated", "logs", "transcript_full.jsonl")
	if checkFileExists(fullPath) {
		return fullPath
	}

	return filepath.Join(brainDir, convID, ".system_generated", "logs", "transcript.jsonl")
}

func checkFileExists(path string) bool {
	hasEmpty := len(path) == 0
	if hasEmpty {
		return false
	}
	info, err := os.Stat(path)
	isFoundFile := err == nil && info.IsDir() == false

	return isFoundFile
}

func findMatchingActiveConvID(repoRoot string) string {
	conv, err := SelectMatchingConversation(repoRoot)
	hasErr := err != nil
	if hasErr {
		return ""
	}

	return conv.ID
}

func findLatestModifiedTranscript() (string, string) {
	brainDir, entries, isFound := readBrainDirEntries()
	if isFound {
		return scanEntriesForLatestTranscript(brainDir, entries)
	}

	return "", ""
}

func readBrainDirEntries() (string, []os.DirEntry, bool) {
	brainDir, err := GetBrainLogsDirPath()
	hasDirErr := err != nil
	if hasDirErr {
		return "", nil, false
	}
	entries, err := os.ReadDir(brainDir)
	hasReadErr := err != nil
	if hasReadErr {
		return "", nil, false
	}

	return brainDir, entries, true
}

func scanEntriesForLatestTranscript(brainDir string, entries []os.DirEntry) (string, string) {
	var latestPath string
	var latestConvID string
	var latestTime time.Time
	for _, e := range entries {
		if e.IsDir() {
			latestConvID, latestPath, latestTime = inspectConvDirForLatest(brainDir, e.Name(), latestConvID, latestPath, latestTime)
		}
	}

	return latestConvID, latestPath
}

func inspectConvDirForLatest(brainDir, name, curID, curPath string, curTime time.Time) (string, string, time.Time) {
	tPath := resolveTranscriptPathForConv(name)
	info, err := os.Stat(tPath)
	isNewer := err == nil && info.ModTime().After(curTime)
	if isNewer {
		return name, tPath, info.ModTime()
	}

	return curID, curPath, curTime
}

type rawTranscriptStepForStatus struct {
	StepIndex int               `json:"step_index"`
	Type      string            `json:"type"`
	Source    string            `json:"source"`
	ToolCalls []json.RawMessage `json:"tool_calls"`
	Status    string            `json:"status"`
}

func detectTranscriptStatus(transcriptPath string) AgyConvStatusType {
	step, hasStep := readLastTranscriptStep(transcriptPath)
	if hasStep {
		return parseConvStatusFromStep(step)
	}

	return AgyConvStatusUnknown
}

func parseConvStatusFromStep(step rawTranscriptStepForStatus) AgyConvStatusType {
	if strings.EqualFold(step.Status, "DONE") || strings.EqualFold(step.Status, "COMPLETED") || strings.EqualFold(step.Status, "IDLE") {
		return AgyConvStatusIdle
	}
	stepType := strings.ToUpper(step.Type)
	source := strings.ToUpper(step.Source)
	if stepType == "USER_INPUT" {
		return AgyConvStatusRunning
	}
	if isModelStatusStep(stepType, source) {
		return resolveModelStepStatus(len(step.ToolCalls) > 0 && !strings.EqualFold(step.Status, "DONE"))
	}

	return AgyConvStatusIdle
}

func resolveModelStepStatus(hasToolCalls bool) AgyConvStatusType {
	if hasToolCalls {
		return AgyConvStatusRunning
	}

	return AgyConvStatusIdle
}

func isModelStatusStep(stepType, source string) bool {
	return stepType == "MODEL" || stepType == "PLANNER_RESPONSE" || source == "MODEL"
}

func openAndStatTranscript(path string) (*os.File, int64, bool) {
	f, err := os.Open(path)
	hasErr := err != nil
	if hasErr {
		return nil, 0, false
	}
	stat, statErr := f.Stat()
	hasStatErr := statErr != nil || stat.Size() == 0
	if hasStatErr {
		_ = f.Close()
		return nil, 0, false
	}

	return f, stat.Size(), true
}

func readLastTranscriptStep(transcriptPath string) (rawTranscriptStepForStatus, bool) {
	f, size, isReady := openAndStatTranscript(transcriptPath)
	if isReady {
		defer f.Close()
		return readTranscriptBuffer(f, size)
	}

	return rawTranscriptStepForStatus{}, false
}

func readTranscriptBuffer(f *os.File, size int64) (rawTranscriptStepForStatus, bool) {
	buf, isReadSuccess := readEndChunk(f, size)
	if isReadSuccess {
		return parseLastStepFromBuffer(buf)
	}

	return rawTranscriptStepForStatus{}, false
}

func calculateChunkOffset(totalSize, maxSize int64) (int64, int64) {
	readSize := maxSize
	if totalSize < readSize {
		readSize = totalSize
	}

	return totalSize - readSize, readSize
}

func readEndChunk(f *os.File, totalSize int64) ([]byte, bool) {
	offset, readSize := calculateChunkOffset(totalSize, 256*1024)
	if _, err := f.Seek(offset, 0); err != nil {
		return nil, false
	}
	buf := make([]byte, readSize)
	n, readErr := f.Read(buf)
	hasReadErr := readErr != nil || n == 0
	if hasReadErr {
		return nil, false
	}

	return buf[:n], true
}

func parseLastStepFromBuffer(buf []byte) (rawTranscriptStepForStatus, bool) {
	lines := strings.Split(string(buf), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		step, isFound := parseSingleTranscriptLine(lines[i])
		if isFound {
			return step, true
		}
	}

	return rawTranscriptStepForStatus{}, false
}

func parseSingleTranscriptLine(line string) (rawTranscriptStepForStatus, bool) {
	trimmed := strings.TrimSpace(line)
	hasEmpty := len(trimmed) == 0
	if hasEmpty {
		return rawTranscriptStepForStatus{}, false
	}
	var step rawTranscriptStepForStatus
	err := json.Unmarshal([]byte(trimmed), &step)
	hasData := err == nil && (len(step.Type) > 0 || len(step.Source) > 0)
	if hasData {
		return step, true
	}

	return rawTranscriptStepForStatus{}, false
}
