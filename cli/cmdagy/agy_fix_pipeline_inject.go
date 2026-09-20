package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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

	if ideProcRes.IsSuccess() {
		return handleIDEInjection(repoRoot, promptPath, ideProcRes.Value.PID)
	}

	return handleCLIInjectionFallback(repoRoot, promptPath, ideProcRes)
}

func handleIDEInjection(repoRoot, promptPath string, pid int) AgyInjectionResult {
	convStatus := DetectConversationExecutionStatus(repoRoot)
	promptContent := readPromptContentOrDefault(promptPath)

	if convStatus == AgyConvStatusRunning {
		_, _ = EnqueuePrompt("pipeline_fix", "Fix CI/CD Pipeline Errors with 4-Part RCA", promptContent)

		return makeQueuedSuccessResult(pid, repoRoot, promptPath)
	}

	copyClipboardIfNotSkipped(promptContent, false)

	return makeIDESuccessResult(pid, repoRoot, promptPath)
}

func handleCLIInjectionFallback(repoRoot, promptPath string, ideRes result.Result[AgyProcessInfo]) AgyInjectionResult {
	cliRes := ResolveAntigravityCLI()
	if cliRes.IsSuccess() {
		return launchAgyBackgroundRunner(cliRes.Value, repoRoot, promptPath, ideRes)
	}

	return makeNoTargetResult(repoRoot, promptPath)
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

	if ideProcRes.IsSuccess() {
		return handleIDEInjectionWithPayload(repoRoot, targetPrompt, promptText, title, promptType, ideProcRes.Value.PID)
	}

	return handleCLIInjectionFallback(repoRoot, targetPrompt, ideProcRes)
}

func handleIDEInjectionWithPayload(repoRoot, promptPath, promptText, title, promptType string, pid int) AgyInjectionResult {
	convStatus := DetectConversationExecutionStatus(repoRoot)
	if convStatus == AgyConvStatusRunning {
		_, _ = EnqueuePrompt(promptType, title, promptText)

		return makeQueuedSuccessResult(pid, repoRoot, promptPath)
	}

	copyClipboardIfNotSkipped(promptText, false)

	return makeIDESuccessResult(pid, repoRoot, promptPath)
}

func makeQueuedSuccessResult(pid int, repoDir, promptPath string) AgyInjectionResult {
	msg := fmt.Sprintf("Active Antigravity IDE is currently busy (RUNNING); prompt queued in agy-prompt-queue.json (PID: %d)", pid)

	return AgyInjectionResult{
		IsSuccess:  true,
		Mode:       AgyInjectionModeQueued,
		PID:        pid,
		Message:    msg,
		PromptPath: promptPath,
		RepoDir:    repoDir,
	}
}

func resolveTargetRepoRoot(repoDir string) string {
	startPath := resolveInitialPath(repoDir)
	root, err := gitutil.RepoRoot(startPath)
	if err == nil && len(root) > 0 {
		return root
	}

	return toAbsPath(startPath)
}

func resolveInitialPath(repoDir string) string {
	if len(repoDir) > 0 {
		return repoDir
	}

	return resolveProjectRootDir()
}

func stageActivePromptInRepo(repoRoot, absPayloadPath string) string {
	targetPrompt := filepath.Join(repoRoot, activeAgyPromptRelativePath)
	if isSamePath(targetPrompt, absPayloadPath) {
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
	if len(path) == 0 {
		return ""
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return string(data)
}

func writePromptFile(path, content string) {
	if len(content) == 0 {
		return
	}

	_ = os.MkdirAll(filepath.Dir(path), 0755)
	_ = os.WriteFile(path, []byte(content), 0644)
}

func buildAgyPromptArg(absPayloadPath string) string {
	return fmt.Sprintf("Autonomous CI/CD pipeline fix: follow all directives in %s", absPayloadPath)
}

func attachInjectionLog(cmd *exec.Cmd, repoDir string) {
	if len(repoDir) == 0 {
		return
	}

	cmd.Dir = repoDir
	logPath := filepath.Join(repoDir, ".ai-memory", "temp", "agy-injection.log")
	_ = os.MkdirAll(filepath.Dir(logPath), 0755)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}

	cmd.Stdout = logFile
	cmd.Stderr = logFile
}

func launchAgyBackgroundRunner(
	binPath string,
	repoDir string,
	promptPath string,
	ideRes result.Result[AgyProcessInfo],
) AgyInjectionResult {
	promptArg := buildAgyPromptArg(promptPath)
	cmd := exec.Command(binPath, "--dangerously-skip-permissions", "-p", promptArg)
	attachInjectionLog(cmd, repoDir)
	configureBackgroundProcess(cmd)

	if err := cmd.Start(); err != nil {
		return makeLaunchFailureResult(err, repoDir, promptPath)
	}

	pid := cmd.Process.Pid
	_ = cmd.Process.Release()

	return makeLaunchSuccessResult(pid, repoDir, promptPath, ideRes)
}

func makeLaunchSuccessResult(
	pid int,
	repoDir string,
	promptPath string,
	ideRes result.Result[AgyProcessInfo],
) AgyInjectionResult {
	msg := formatLaunchSuccessMessage(pid, repoDir, ideRes)

	return AgyInjectionResult{
		IsSuccess:  true,
		Mode:       AgyInjectionModeCLI,
		PID:        pid,
		Message:    msg,
		PromptPath: promptPath,
		RepoDir:    repoDir,
	}
}

func formatLaunchSuccessMessage(
	pid int,
	repoDir string,
	ideRes result.Result[AgyProcessInfo],
) string {
	if ideRes.IsSuccess() {
		return fmt.Sprintf("Started Antigravity CLI background runner (PID: %d) with active IDE detected (PID: %d)", pid, ideRes.Value.PID)
	}

	return fmt.Sprintf("Started Antigravity CLI background runner (PID: %d) in %s", pid, repoDir)
}

func makeIDESuccessResult(pid int, repoDir, promptPath string) AgyInjectionResult {
	msg := fmt.Sprintf("Active Antigravity IDE detected (PID: %d); staged fix prompt in %s", pid, promptPath)

	return AgyInjectionResult{
		IsSuccess:  true,
		Mode:       AgyInjectionModeIDE,
		PID:        pid,
		Message:    msg,
		PromptPath: promptPath,
		RepoDir:    repoDir,
	}
}

func makeLaunchFailureResult(err error, repoDir, promptPath string) AgyInjectionResult {
	return AgyInjectionResult{
		IsSuccess:  false,
		Mode:       AgyInjectionModeCLI,
		Message:    fmt.Sprintf("failed to launch agy CLI process: %v", err),
		PromptPath: promptPath,
		RepoDir:    repoDir,
	}
}

func makeSkipInjectionResult() AgyInjectionResult {
	return AgyInjectionResult{
		IsSuccess: false,
		Mode:      AgyInjectionModeNone,
		Message:   "direct injection skipped by flag (--no-inject)",
	}
}

func makeNoTargetResult(repoDir, promptPath string) AgyInjectionResult {
	return AgyInjectionResult{
		IsSuccess:  false,
		Mode:       AgyInjectionModeNone,
		Message:    "antigravity binary (agy) not detected and no active IDE process found",
		PromptPath: promptPath,
		RepoDir:    repoDir,
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
	if len(transcriptPath) == 0 {
		return AgyConversationExecutionState{Status: AgyConvStatusUnknown}
	}

	status := detectTranscriptStatus(transcriptPath)

	return AgyConversationExecutionState{
		ConvID:         convID,
		Status:         status,
		TranscriptPath: transcriptPath,
	}
}

func findActiveConvTranscript(repoRoot string) (string, string) {
	convID := findMatchingActiveConvID(repoRoot)
	if len(convID) > 0 {
		path := resolveTranscriptPathForConv(convID)
		if checkFileExists(path) {
			return convID, path
		}
	}

	return findLatestModifiedTranscript()
}

func resolveTranscriptPathForConv(convID string) string {
	brainDir, err := GetBrainLogsDirPath()
	if err != nil {
		return ""
	}

	fullPath := filepath.Join(brainDir, convID, ".system_generated", "logs", "transcript_full.jsonl")
	if checkFileExists(fullPath) {
		return fullPath
	}

	return filepath.Join(brainDir, convID, ".system_generated", "logs", "transcript.jsonl")
}

func checkFileExists(path string) bool {
	if len(path) == 0 {
		return false
	}
	info, err := os.Stat(path)

	return err == nil && info.IsDir() == false
}

func findMatchingActiveConvID(repoRoot string) string {
	convs, err := scanAllConversations()
	if err != nil || len(convs) == 0 {
		return ""
	}

	pClean := cleanProjectWorkspace(repoRoot)
	var bestConv AgyConvInfo
	hasMatch := false

	for _, c := range convs {
		if isConvPathMatch(pClean, c.CleanPath) {
			bestConv = pickMoreActiveConv(bestConv, c)
			hasMatch = true
		}
	}
	if hasMatch {
		return bestConv.ID
	}

	return ""
}

func pickMoreActiveConv(curr, next AgyConvInfo) AgyConvInfo {
	if next.StepCount > curr.StepCount || next.UserSteps > curr.UserSteps {
		return next
	}

	return curr
}

func findLatestModifiedTranscript() (string, string) {
	brainDir, err := GetBrainLogsDirPath()
	if err != nil {
		return "", ""
	}
	entries, err := os.ReadDir(brainDir)
	if err != nil {
		return "", ""
	}

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
	if err == nil && info.ModTime().After(curTime) {
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
	if hasStep == false {
		return AgyConvStatusUnknown
	}

	return parseConvStatusFromStep(step)
}

func parseConvStatusFromStep(step rawTranscriptStepForStatus) AgyConvStatusType {
	stepType := strings.ToUpper(step.Type)
	source := strings.ToUpper(step.Source)
	hasToolCalls := len(step.ToolCalls) > 0

	if stepType == "USER_INPUT" {
		return AgyConvStatusRunning
	}
	if isModelStatusStep(stepType, source) {
		if hasToolCalls {
			return AgyConvStatusRunning
		}
		return AgyConvStatusIdle
	}
	if stepType == "GENERIC" || source == "GENERIC" {
		return AgyConvStatusRunning
	}

	return AgyConvStatusRunning
}

func isModelStatusStep(stepType, source string) bool {
	return stepType == "MODEL" || stepType == "PLANNER_RESPONSE" || source == "MODEL"
}

func readLastTranscriptStep(transcriptPath string) (rawTranscriptStepForStatus, bool) {
	f, err := os.Open(transcriptPath)
	if err != nil {
		return rawTranscriptStepForStatus{}, false
	}
	defer f.Close()

	stat, statErr := f.Stat()
	if statErr != nil || stat.Size() == 0 {
		return rawTranscriptStepForStatus{}, false
	}

	buf, isReadSuccess := readEndChunk(f, stat.Size())
	if isReadSuccess == false {
		return rawTranscriptStepForStatus{}, false
	}

	return parseLastStepFromBuffer(buf)
}

func readEndChunk(f *os.File, totalSize int64) ([]byte, bool) {
	readSize := int64(256 * 1024)
	if totalSize < readSize {
		readSize = totalSize
	}
	offset := totalSize - readSize
	if _, err := f.Seek(offset, 0); err != nil {
		return nil, false
	}
	buf := make([]byte, readSize)
	n, readErr := f.Read(buf)
	if readErr != nil || n == 0 {
		return nil, false
	}

	return buf[:n], true
}

func parseLastStepFromBuffer(buf []byte) (rawTranscriptStepForStatus, bool) {
	lines := strings.Split(string(buf), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if len(line) == 0 {
			continue
		}
		var step rawTranscriptStepForStatus
		if err := json.Unmarshal([]byte(line), &step); err == nil && (len(step.Type) > 0 || len(step.Source) > 0) {
			return step, true
		}
	}

	return rawTranscriptStepForStatus{}, false
}
