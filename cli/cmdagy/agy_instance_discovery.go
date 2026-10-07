package cmdagy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type diskInstancesFile struct {
	ActiveInstanceID string              `json:"active_instance_id"`
	Instances        []diskInstanceEntry `json:"instances"`
}

type diskInstanceEntry struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	DataDir        string  `json:"data_dir"`
	ExecutablePath *string `json:"executable_path"`
	ExtensionsDir  *string `json:"extensions_dir"`
	BoundAccountID string  `json:"bound_account_id"`
	BoundEmail     string  `json:"bound_email"`
	CreatedAt      int64   `json:"created_at"`
	LastUsed       int64   `json:"last_used"`
	IsDefault      bool    `json:"is_default"`
	PID            *int    `json:"pid"`
	SeqNum         int     `json:"seq_num"`
}

// DiscoverAllAgyInstances enumerates primary, registered, and active Antigravity instances.
func DiscoverAllAgyInstances() ([]AgyInstanceInfo, error) {
	var instances []AgyInstanceInfo
	seenIDs := make(map[string]bool)

	primaryInst, err := discoverPrimaryInstance()
	if err == nil {
		instances = append(instances, primaryInst)
		seenIDs[primaryInst.InstanceID] = true
		seenIDs["default"] = true
	}

	diskInstances := loadInstancesFromDiskRegistries()
	procInstances := scanRunningProcessesForInstances()

	for _, disk := range diskInstances {
		if seenIDs[disk.InstanceID] {
			continue
		}
		seenIDs[disk.InstanceID] = true

		for _, p := range procInstances {
			isMatch := p.InstanceID == disk.InstanceID || (disk.ConfigDir != "" && strings.Contains(p.ConfigDir, disk.ConfigDir))
			if isMatch {
				mergeProcMatch(&disk, p)
				break
			}
		}

		if disk.SummariesDBPath != "" {
			disk.ActiveWorkspaces = queryDistinctWorkspacesFromDB(disk.SummariesDBPath)
		}
		if len(disk.ActiveWorkspaces) == 0 && primaryInst.SummariesDBPath != "" {
			disk.ActiveWorkspaces = primaryInst.ActiveWorkspaces
		}

		instances = append(instances, disk)
	}

	for _, proc := range procInstances {
		if seenIDs[proc.InstanceID] {
			continue
		}
		seenIDs[proc.InstanceID] = true
		if proc.SummariesDBPath != "" {
			proc.ActiveWorkspaces = queryDistinctWorkspacesFromDB(proc.SummariesDBPath)
		}
		instances = append(instances, proc)
	}

	return instances, nil
}

func discoverPrimaryInstance() (AgyInstanceInfo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return AgyInstanceInfo{}, apperror.WrapSimple(err, "discover_primary_instance")
	}

	configDir := filepath.Join(home, ".gemini", "antigravity")
	brainDir := filepath.Join(configDir, "brain")
	summariesDB := filepath.Join(configDir, "conversation_summaries.db")

	inst := AgyInstanceInfo{
		InstanceID:       "primary",
		InstanceName:     "Default Antigravity Workspace",
		ConfigDir:        configDir,
		BrainDir:         brainDir,
		SummariesDBPath:  summariesDB,
		ActiveWorkspaces: queryDistinctWorkspacesFromDB(summariesDB),
		IsPrimary:        true,
		IsRunning:        false,
		LastActiveAt:     time.Now().UTC().Format(time.RFC3339),
	}

	procs := scanRunningProcessesForInstances()
	for _, p := range procs {
		if p.IsPrimary || strings.Contains(p.ConfigDir, ".config/Antigravity") || p.ConfigDir == configDir {
			inst.IsRunning = true
			inst.ProcessID = p.ProcessID
			inst.LanguageServer = p.LanguageServer
			break
		}
	}

	if !inst.IsRunning && len(procs) > 0 {
		attachFirstRunningProcess(&inst, procs)
	}

	return inst, nil
}

func mergeProcMatch(disk *AgyInstanceInfo, p AgyInstanceInfo) {
	disk.IsRunning = true
	if p.ProcessID > 0 {
		disk.ProcessID = p.ProcessID
	}
	if p.LanguageServer != "" {
		disk.LanguageServer = p.LanguageServer
	}
}

func attachFirstRunningProcess(inst *AgyInstanceInfo, procs []AgyInstanceInfo) {
	for _, p := range procs {
		if p.ProcessID > 0 {
			inst.IsRunning = true
			inst.ProcessID = p.ProcessID
			inst.LanguageServer = p.LanguageServer
			break
		}
	}
}

func loadInstancesFromDiskRegistries() []AgyInstanceInfo {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	candidateFiles := []string{
		filepath.Join(home, ".antigravity_tools", "instances", "instances.json"),
		filepath.Join(home, ".gemini", "antigravity", "instances.json"),
	}

	var results []AgyInstanceInfo
	for _, candidate := range candidateFiles {
		data, readErr := os.ReadFile(candidate)
		if readErr != nil {
			continue
		}

		var filePayload diskInstancesFile
		if jsonErr := json.Unmarshal(data, &filePayload); jsonErr != nil {
			continue
		}

		for _, item := range filePayload.Instances {
			if item.ID == "default" || item.ID == "primary" {
				continue
			}

			brainDir := filepath.Join(item.DataDir, "brain")
			if _, statErr := os.Stat(brainDir); statErr != nil {
				brainDir = filepath.Join(home, ".gemini", "antigravity", "brain")
			}

			summariesDB := filepath.Join(item.DataDir, "conversation_summaries.db")
			if _, statErr := os.Stat(summariesDB); statErr != nil {
				summariesDB = filepath.Join(home, ".gemini", "antigravity", "conversation_summaries.db")
			}

			pid := 0
			if item.PID != nil {
				pid = *item.PID
			}

			lastActive := time.Now().UTC().Format(time.RFC3339)
			if item.LastUsed > 0 {
				lastActive = time.Unix(item.LastUsed, 0).UTC().Format(time.RFC3339)
			}

			info := AgyInstanceInfo{
				InstanceID:      item.ID,
				InstanceName:    item.Name,
				ProcessID:       pid,
				ConfigDir:       item.DataDir,
				BrainDir:        brainDir,
				SummariesDBPath: summariesDB,
				IsPrimary:       item.IsDefault,
				IsRunning:       isProcessAlive(pid),
				LastActiveAt:    lastActive,
			}
			results = append(results, info)
		}
	}

	return results
}

func scanRunningProcessesForInstances() []AgyInstanceInfo {
	if runtime.GOOS == "windows" {
		return scanRunningProcessesWindows()
	}
	return scanRunningProcessesUnix()
}

func scanRunningProcessesUnix() []AgyInstanceInfo {
	cmd := exec.Command("ps", "-eo", "pid,args")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var found []AgyInstanceInfo
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.Contains(trimmed, "antigravity") && !strings.Contains(trimmed, "language_server") {
			continue
		}

		parts := strings.Fields(trimmed)
		if len(parts) < 2 {
			continue
		}

		pid, parseErr := strconv.Atoi(parts[0])
		if parseErr != nil {
			continue
		}

		info := parseProcessCommandLine(pid, trimmed)
		if info != nil {
			found = append(found, *info)
		}
	}

	return found
}

func scanRunningProcessesWindows() []AgyInstanceInfo {
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		"Get-CimInstance Win32_Process | Where-Object { $_.Name -like '*antigravity*' -or $_.CommandLine -like '*language_server*' } | Select-Object ProcessId, CommandLine | ConvertTo-Json -Compress")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	items := unmarshalWinProcItems(output)

	var found []AgyInstanceInfo
	for _, it := range items {
		info := parseProcessCommandLine(it.ProcessID, it.CommandLine)
		if info != nil {
			found = append(found, *info)
		}
	}

	return found
}

type winProcItem struct {
	ProcessID   int    `json:"ProcessId"`
	CommandLine string `json:"CommandLine"`
}

func unmarshalWinProcItems(output []byte) []winProcItem {
	var items []winProcItem
	if err := json.Unmarshal(output, &items); err == nil {
		return items
	}

	var single winProcItem
	if err := json.Unmarshal(output, &single); err == nil {
		return []winProcItem{single}
	}

	return nil
}

func parseProcessCommandLine(pid int, cmdLine string) *AgyInstanceInfo {
	isLanguageServer := strings.Contains(cmdLine, "language_server")
	isAntigravity := strings.Contains(cmdLine, "antigravity")

	if !isLanguageServer && !isAntigravity {
		return nil
	}

	langServer := strings.TrimPrefix(extractCmdLineFlagValue(cmdLine, "--host_bridge_url="), "http://")
	dataDir := extractCmdLineFlagValue(cmdLine, "--user-data-dir=")
	instanceID, instanceName, isPrimary := classifyInstanceDataDir(dataDir, isLanguageServer)

	return &AgyInstanceInfo{
		InstanceID:     instanceID,
		InstanceName:   instanceName,
		ProcessID:      pid,
		LanguageServer: langServer,
		ConfigDir:      dataDir,
		IsPrimary:      isPrimary,
		IsRunning:      true,
		LastActiveAt:   time.Now().UTC().Format(time.RFC3339),
	}
}

func extractCmdLineFlagValue(cmdLine, flag string) string {
	if !strings.Contains(cmdLine, flag) {
		return ""
	}

	parts := strings.Split(cmdLine, flag)
	if len(parts) <= 1 {
		return ""
	}

	fields := strings.Fields(parts[1])
	if len(fields) == 0 {
		return ""
	}

	return fields[0]
}

func classifyInstanceDataDir(dataDir string, isLanguageServer bool) (string, string, bool) {
	if dataDir == "" {
		return "primary", "Antigravity Active Process", isLanguageServer
	}

	seg := resolveInstanceSegment(dataDir)
	if seg != "" {
		return seg, seg, false
	}

	if strings.Contains(dataDir, ".config/Antigravity") {
		return "primary", "Antigravity Active Process", true
	}

	return "primary", "Antigravity Active Process", false
}

func resolveInstanceSegment(dataDir string) string {
	if !strings.Contains(dataDir, "instances/") {
		return ""
	}

	seg := filepath.Base(filepath.Dir(dataDir))
	if seg != "" && seg != "." {
		return seg
	}

	return ""
}

func queryDistinctWorkspacesFromDB(dbPath string) []string {
	if dbPath == "" {
		return nil
	}
	if _, err := os.Stat(dbPath); err != nil {
		return nil
	}

	conn, err := store.OpenSQLiteDB(dbPath)
	if err != nil {
		return nil
	}
	defer conn.Close()

	query := `SELECT DISTINCT workspace_uris FROM conversation_summaries WHERE (killed IS NULL OR killed = 0) AND workspace_uris IS NOT NULL ORDER BY last_modified_time DESC LIMIT 20;`
	rows, qErr := conn.Query(query)
	if qErr != nil {
		return nil
	}
	defer rows.Close()

	seen := make(map[string]bool)
	var workspaces []string
	for rows.Next() {
		var raw sql.NullString
		if scanErr := rows.Scan(&raw); scanErr != nil || !raw.Valid {
			continue
		}
		clean := extractCleanWorkspaceFromURIs(raw.String)
		if clean != "" && !seen[clean] {
			seen[clean] = true
			workspaces = append(workspaces, clean)
		}
	}

	return workspaces
}

// ResolveInstance resolves an Antigravity instance by ID or alias with default fallback to primary.
func ResolveInstance(instanceID string) (*AgyInstanceInfo, error) {
	instances, err := DiscoverAllAgyInstances()
	if err != nil {
		return nil, err
	}

	target := strings.TrimSpace(strings.ToLower(instanceID))
	if target == "" || target == "primary" || target == "default" {
		return resolvePrimaryInstance(instances)
	}

	for _, inst := range instances {
		if strings.EqualFold(inst.InstanceID, target) || strings.EqualFold(inst.InstanceName, target) {
			return &inst, nil
		}
	}

	return nil, apperror.NewSimple("resolve_instance", fmt.Sprintf("instance '%s' not found", instanceID))
}

func resolvePrimaryInstance(instances []AgyInstanceInfo) (*AgyInstanceInfo, error) {
	for _, inst := range instances {
		if inst.IsPrimary || inst.InstanceID == "primary" || inst.InstanceID == "default" {
			return &inst, nil
		}
	}
	if len(instances) > 0 {
		return &instances[0], nil
	}
	return nil, apperror.NewSimple("resolve_instance", "primary instance unavailable")
}

// QueryInstancePrompts queries running, queued, and conversation prompts across instances.
func QueryInstancePrompts(opts AgyInstancePromptQueryOptions) (*AgyMultiInstancePromptResponse, error) {
	allInstances, err := DiscoverAllAgyInstances()
	if err != nil {
		return nil, err
	}

	targetID := strings.TrimSpace(opts.InstanceID)
	targetInstances, resErr := resolveTargetInstances(targetID, opts.IsAll, allInstances)
	if resErr != nil {
		return &AgyMultiInstancePromptResponse{
			IsSuccess:      false,
			TotalInstances: 0,
			TotalPrompts:   0,
			Instances:      nil,
			Error:          resErr.Error(),
		}, nil
	}

	limit := opts.Limit
	if limit <= 0 {
		limit = 10
	}
	maxWords := opts.MaxWords
	if maxWords <= 0 {
		maxWords = 100
	}

	var payloads []AgyInstancePromptPayload
	totalPrompts := 0

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, inst := range targetInstances {
		wg.Add(1)
		go func(currentInst AgyInstanceInfo) {
			defer wg.Done()
			payload := collectPayloadForInstance(currentInst, opts.Status, limit, maxWords, opts.IncludeConvs)
			mu.Lock()
			payloads = append(payloads, payload)
			totalPrompts += payload.RunningCount + payload.QueuedCount
			mu.Unlock()
		}(inst)
	}
	wg.Wait()

	return &AgyMultiInstancePromptResponse{
		IsSuccess:      true,
		TotalInstances: len(payloads),
		TotalPrompts:   totalPrompts,
		Instances:      payloads,
	}, nil
}

func resolveTargetInstances(targetID string, isAll bool, allInstances []AgyInstanceInfo) ([]AgyInstanceInfo, error) {
	if isAll || strings.EqualFold(targetID, "all") || targetID == "" {
		return allInstances, nil
	}
	resolved, rErr := ResolveInstance(targetID)
	if rErr != nil {
		return nil, rErr
	}
	return []AgyInstanceInfo{*resolved}, nil
}

func collectPayloadForInstance(inst AgyInstanceInfo, statusFilter string, limit, maxWords int, includeConvs bool) AgyInstancePromptPayload {
	payload := AgyInstancePromptPayload{
		InstanceID:   inst.InstanceID,
		InstanceName: inst.InstanceName,
		Running:      make([]AgyPromptSnapshotItem, 0),
		Queued:       make([]AgyPromptQueueEntry, 0),
		RecentConvs:  make([]AgyConversationPreview, 0),
	}

	queryRunning := statusFilter == "" || statusFilter == "all" || statusFilter == "running"
	queryQueued := statusFilter == "" || statusFilter == "all" || statusFilter == "queued"

	if queryRunning {
		runningItems := truncateSliceByLimit(fetchInstanceRunningPrompts(inst, maxWords), limit)
		payload.Running = runningItems
		payload.RunningCount = len(runningItems)
	}

	if queryQueued {
		queuedItems := truncateSliceByLimit(fetchInstanceQueuedPrompts(inst), limit)
		payload.Queued = queuedItems
		payload.QueuedCount = len(queuedItems)
	}

	if includeConvs {
		payload.RecentConvs = fetchInstanceRecentConversations(inst, limit, maxWords)
	}

	return payload
}

func truncateSliceByLimit[T any](items []T, limit int) []T {
	if limit > 0 && len(items) > limit {
		return items[:limit]
	}

	return items
}

func fetchInstanceRunningPrompts(inst AgyInstanceInfo, maxWords int) []AgyPromptSnapshotItem {
	dbPath := inst.SummariesDBPath
	if dbPath == "" {
		return nil
	}

	conn, err := store.OpenSQLiteDB(dbPath)
	if err != nil {
		return nil
	}
	defer conn.Close()

	activeConvs, qErr := queryRunningConversations(conn)
	if qErr != nil {
		return nil
	}

	var items []AgyPromptSnapshotItem
	for _, c := range activeConvs {
		preview := queryConversationPreview(conn, c.ConversationID)
		if preview == "" {
			preview = c.Title
		}
		if maxWords > 0 {
			preview = CompactWords(preview, maxWords)
		}

		projName := c.ProjectID
		if projName == "" {
			projName = c.WorkspacePath
		}

		items = append(items, AgyPromptSnapshotItem{
			Type:           "running",
			ProjectName:    projName,
			WorkspacePath:  c.WorkspacePath,
			ConversationID: c.ConversationID,
			Title:          c.Title,
			PromptPreview:  preview,
			Status:         "RUNNING",
			ElapsedSeconds: c.ElapsedSeconds,
		})
	}

	return items
}

func fetchInstanceQueuedPrompts(inst AgyInstanceInfo) []AgyPromptQueueEntry {
	var queued []AgyPromptQueueEntry

	for _, ws := range inst.ActiveWorkspaces {
		queueFile := filepath.Join(ws, ".ai-memory", "temp", "agy-prompt-queue.json")
		data, err := os.ReadFile(queueFile)
		if err != nil {
			continue
		}

		var q AgyPromptQueueFile
		if jsonErr := json.Unmarshal(data, &q); jsonErr == nil {
			queued = append(queued, q.Queued...)
		}
	}

	return queued
}

func fetchInstanceRecentConversations(inst AgyInstanceInfo, limit, maxWords int) []AgyConversationPreview {
	dbPath := inst.SummariesDBPath
	if dbPath == "" {
		return nil
	}

	conn, err := store.OpenSQLiteDB(dbPath)
	if err != nil {
		return nil
	}
	defer conn.Close()

	if limit <= 0 {
		limit = 10
	}

	query := `SELECT conversation_id, title, workspace_uris, preview, step_count, last_modified_time, not_fully_idle
		FROM conversation_summaries
		WHERE (killed IS NULL OR killed = 0)
		ORDER BY last_modified_time DESC
		LIMIT ?;`

	rows, qErr := conn.Query(query, limit)
	if qErr != nil {
		return nil
	}
	defer rows.Close()

	var convs []AgyConversationPreview
	for rows.Next() {
		var cid, title, wsRaw, prev, lastMod sql.NullString
		var stepCount, idle int
		if scanErr := rows.Scan(&cid, &title, &wsRaw, &prev, &stepCount, &lastMod, &idle); scanErr != nil {
			continue
		}

		previewText := prev.String
		if maxWords > 0 && previewText != "" {
			previewText = CompactWords(previewText, maxWords)
		}

		cleanWs := extractCleanWorkspaceFromURIs(wsRaw.String)

		convs = append(convs, AgyConversationPreview{
			ConversationID: cid.String,
			Title:          title.String,
			WorkspacePath:  cleanWs,
			LastPromptText: previewText,
			StepCount:      stepCount,
			UpdatedAt:      lastMod.String,
			IsActive:       idle != 0,
		})
	}

	return convs
}
