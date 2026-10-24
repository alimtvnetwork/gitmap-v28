package cmdagy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

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

	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return nil
	}
	conn := openRes.Value
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
