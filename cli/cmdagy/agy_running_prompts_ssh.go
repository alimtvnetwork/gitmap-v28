package cmdagy

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AggregateSSHRunningPromptsLs queries cluster SSH nodes and aggregates running prompts.
func AggregateSSHRunningPromptsLs(local []store.RunningPromptRecord, limit, wordCount int, isFull, isJSON bool) error {
	tagLocalPrompts(local)
	remote := fetchClusterSSHRunningPromptsLs(wordCount, isFull)
	all := mergePromptRecords(local, remote)
	all = applyPromptsLimit(all, limit)
	if isJSON {
		return printJSON(all)
	}
	RenderSSHRunningPromptsTable(all, isFull)
	return nil
}

func mergePromptRecords(local, remote []store.RunningPromptRecord) []store.RunningPromptRecord {
	merged := make([]store.RunningPromptRecord, 0, len(local)+len(remote))
	merged = append(merged, local...)
	merged = append(merged, remote...)
	return merged
}

func tagLocalPrompts(items []store.RunningPromptRecord) {
	for i := range items {
		items[i].Node = "local"
		items[i].Host = "127.0.0.1"
	}
}

func fetchClusterSSHRunningPromptsLs(wordCount int, isFull bool) []store.RunningPromptRecord {
	if SSHConnectionsFetcher == nil {
		return nil
	}
	conns, err := SSHConnectionsFetcher()
	if err != nil || len(conns) == 0 {
		return nil
	}
	var aggregated []store.RunningPromptRecord
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, c := range conns {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			records := querySingleNodePrompts(conn, wordCount, isFull)
			mu.Lock()
			aggregated = append(aggregated, records...)
			mu.Unlock()
		}(c)
	}
	wg.Wait()
	return aggregated
}

func querySingleNodePrompts(c db.SSHConnection, wordCount int, isFull bool) []store.RunningPromptRecord {
	client, err := dialSSHNodeClient(c)
	if err != nil {
		return []store.RunningPromptRecord{buildOfflinePromptRecord(c)}
	}
	defer client.Close()
	cmdStr := fmt.Sprintf("gitmap agy running-prompts ls --json --wc %d", wordCount)
	if isFull {
		cmdStr += " --full"
	}
	out, runErr := crypto.RunCommand(client, cmdStr, resolveNodeShell(c.OS))
	if runErr != nil {
		return []store.RunningPromptRecord{buildNodeErrorPromptRecord(c, runErr)}
	}
	return parseRemotePrompts(c, out)
}

func parseRemotePrompts(c db.SSHConnection, out string) []store.RunningPromptRecord {
	cleanJSON := extractJSONArrayFromOutput(out)
	if len(cleanJSON) == 0 {
		return nil
	}
	var records []store.RunningPromptRecord
	if err := json.Unmarshal([]byte(cleanJSON), &records); err != nil {
		return []store.RunningPromptRecord{buildNodeErrorPromptRecord(c, err)}
	}
	for i := range records {
		records[i].Node = c.Alias
		records[i].Host = c.IPAddress
	}
	return records
}

func extractJSONArrayFromOutput(out string) string {
	start := strings.Index(out, "[")
	end := strings.LastIndex(out, "]")
	if start >= 0 && end > start {
		return out[start : end+1]
	}
	return ""
}

func extractJSONObjectFromOutput(out string) string {
	start := strings.Index(out, "{")
	end := strings.LastIndex(out, "}")
	if start >= 0 && end > start {
		return out[start : end+1]
	}
	return ""
}

func buildOfflinePromptRecord(c db.SSHConnection) store.RunningPromptRecord {
	return store.RunningPromptRecord{
		ProjectName: "(unreachable)",
		Node:        c.Alias,
		Host:        c.IPAddress,
		Status:      store.PromptStatusEnqueued,
		Snippet:     "SSH connection failed or node offline",
	}
}

func buildNodeErrorPromptRecord(c db.SSHConnection, err error) store.RunningPromptRecord {
	return store.RunningPromptRecord{
		ProjectName: "(node error)",
		Node:        c.Alias,
		Host:        c.IPAddress,
		Status:      store.PromptStatusEnqueued,
		Snippet:     TruncateErr(err.Error(), 40),
	}
}

// TruncateErr trims error messages to fit in table previews.
func TruncateErr(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}

// RenderSSHRunningPromptsTable prints active and queued prompts with a NODE column.
func RenderSSHRunningPromptsTable(items []store.RunningPromptRecord, isFull bool) {
	fmt.Println()
	fmt.Printf("  %s● Antigravity Cross-Node Running & Queued Prompts (%d items)%s\n", constants.ColorCyan, len(items), constants.ColorReset)
	fmt.Printf("  %s%-12s %-18s %-10s %-8s %s%s\n", constants.ColorWhite, "NODE", "PROJECT", "STATUS", "WORDS", "PROMPT", constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("─", 90), constants.ColorReset)
	if len(items) == 0 {
		fmt.Printf("  %sNo active or queued prompts found across nodes.%s\n\n", constants.ColorDim, constants.ColorReset)
		return
	}
	for _, it := range items {
		printSSHRunningPromptRow(it, isFull)
	}
	fmt.Println()
}

func printSSHRunningPromptRow(it store.RunningPromptRecord, isFull bool) {
	statusColor := constants.ColorGreen
	if it.Status == store.PromptStatusEnqueued {
		statusColor = constants.ColorYellow
	}
	node := it.Node
	if len(node) == 0 {
		node = "local"
	}
	proj := it.ProjectName
	if len(proj) > 18 {
		proj = proj[:15] + "..."
	}
	text := it.Snippet
	if isFull {
		text = it.Prompt
	}
	fmt.Printf("  %-12s %-18s %s%-10s%s %-8d %s\n", node, proj, statusColor, strings.ToUpper(string(it.Status)), constants.ColorReset, it.WordCount, text)
}

// AggregateSSHRunningPromptsBackup delegates backup across all cluster SSH nodes.
func AggregateSSHRunningPromptsBackup(customFile string, isJSON bool) error {
	localErr := RunRunningPromptsBackup(customFile, isJSON, false)
	if SSHConnectionsFetcher == nil {
		return localErr
	}
	conns, err := SSHConnectionsFetcher()
	if err != nil || len(conns) == 0 {
		return localErr
	}
	var wg sync.WaitGroup
	for _, c := range conns {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			delegateNodeBackup(conn, customFile, isJSON)
		}(c)
	}
	wg.Wait()
	return nil
}

func delegateNodeBackup(c db.SSHConnection, customFile string, isJSON bool) {
	client, err := dialSSHNodeClient(c)
	if err != nil {
		printNodeSSHNotice(c.Alias, "connection failed")
		return
	}
	defer client.Close()
	cmdStr := "gitmap agy backup-running-prompts --json"
	if len(customFile) > 0 {
		cmdStr += fmt.Sprintf(" -f %q", customFile)
	}
	out, runErr := crypto.RunCommand(client, cmdStr, resolveNodeShell(c.OS))
	if runErr != nil {
		printNodeSSHNotice(c.Alias, "backup command failed")
		return
	}
	renderRemoteBackupSummary(c.Alias, out, isJSON)
}

func printNodeSSHNotice(node, msg string) {
	fmt.Printf("  %s[%s] Warning: %s%s\n", constants.ColorYellow, node, msg, constants.ColorReset)
}

func renderRemoteBackupSummary(node, out string, isJSON bool) {
	if isJSON {
		fmt.Printf("[%s] %s\n", node, strings.TrimSpace(out))
		return
	}
	clean := extractJSONObjectFromOutput(out)
	var summary store.PromptBackupSummary
	if err := json.Unmarshal([]byte(clean), &summary); err == nil {
		fmt.Printf("%s✔ [%s] Backed up %d prompts (Batch: %s, Running: %d, Enqueued: %d)%s\n",
			constants.ColorGreen, node, summary.TotalPrompts, summary.BatchId,
			summary.RunningCount, summary.EnqueuedCount, constants.ColorReset)
		return
	}
	fmt.Printf("  [%s] %s\n", node, strings.TrimSpace(out))
}

// AggregateSSHRunningPromptsBackupLs lists backup batches across cluster SSH nodes.
func AggregateSSHRunningPromptsBackupLs(customFile string, isJSON bool) error {
	localBatches := fetchLocalBatches(customFile)
	remoteBatches := fetchClusterSSHBackupBatches(customFile)
	all := mergeBatchRecords(localBatches, remoteBatches)
	if isJSON {
		return printJSON(all)
	}
	RenderSSHBackupBatchesTable(all)
	return nil
}

func mergeBatchRecords(local, remote []store.PromptBackupBatchRecord) []store.PromptBackupBatchRecord {
	merged := make([]store.PromptBackupBatchRecord, 0, len(local)+len(remote))
	merged = append(merged, local...)
	merged = append(merged, remote...)
	return merged
}

func fetchLocalBatches(customFile string) []store.PromptBackupBatchRecord {
	dbConn, err := store.OpenBackupPromptsSplitDB(customFile)
	if err != nil {
		return nil
	}
	defer dbConn.Close()
	batches, _ := dbConn.ListBackupBatches()
	for i := range batches {
		batches[i].Node = "local"
		batches[i].Host = "127.0.0.1"
	}
	return batches
}

func fetchClusterSSHBackupBatches(customFile string) []store.PromptBackupBatchRecord {
	if SSHConnectionsFetcher == nil {
		return nil
	}
	conns, err := SSHConnectionsFetcher()
	if err != nil || len(conns) == 0 {
		return nil
	}
	var aggregated []store.PromptBackupBatchRecord
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, c := range conns {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			batches := querySingleNodeBackupBatches(conn, customFile)
			mu.Lock()
			aggregated = append(aggregated, batches...)
			mu.Unlock()
		}(c)
	}
	wg.Wait()
	return aggregated
}

func querySingleNodeBackupBatches(c db.SSHConnection, customFile string) []store.PromptBackupBatchRecord {
	client, err := dialSSHNodeClient(c)
	if err != nil {
		return nil
	}
	defer client.Close()
	cmdStr := "gitmap agy backup-running-prompts ls --json"
	if len(customFile) > 0 {
		cmdStr += fmt.Sprintf(" -f %q", customFile)
	}
	out, runErr := crypto.RunCommand(client, cmdStr, resolveNodeShell(c.OS))
	if runErr != nil {
		return nil
	}
	return parseRemoteBackupBatches(c, out)
}

func parseRemoteBackupBatches(c db.SSHConnection, out string) []store.PromptBackupBatchRecord {
	clean := extractJSONArrayFromOutput(out)
	var batches []store.PromptBackupBatchRecord
	if err := json.Unmarshal([]byte(clean), &batches); err != nil {
		return nil
	}
	for i := range batches {
		batches[i].Node = c.Alias
		batches[i].Host = c.IPAddress
	}
	return batches
}

// RenderSSHBackupBatchesTable prints backup batches with a NODE column.
func RenderSSHBackupBatchesTable(batches []store.PromptBackupBatchRecord) {
	fmt.Println()
	fmt.Printf("  %s┌── Antigravity Cross-Node Backup Registry (%d batches) ──────────────┐%s\n", constants.ColorCyan, len(batches), constants.ColorReset)
	fmt.Printf("  %s%-12s %-14s %-20s %-6s %-8s %-9s %s%s\n", constants.ColorWhite, "NODE", "Batch ID", "Created At", "Total", "Running", "Enqueued", "Status", constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("─", 85), constants.ColorReset)
	if len(batches) == 0 {
		fmt.Printf("  %sNo backup batches recorded across cluster nodes.%s\n\n", constants.ColorDim, constants.ColorReset)
		return
	}
	for _, b := range batches {
		printSSHBackupBatchRow(b)
	}
	fmt.Println()
}

func printSSHBackupBatchRow(b store.PromptBackupBatchRecord) {
	statusColor := constants.ColorGreen
	if b.Status == "restored" {
		statusColor = constants.ColorYellow
	}
	node := b.Node
	if len(node) == 0 {
		node = "local"
	}
	created := b.CreatedAt
	if len(created) > 19 {
		created = created[:19]
	}
	fmt.Printf("  %-12s %-14s %-20s %-6d %-8d %-9d %s%s%s\n",
		node, b.BatchID, created, b.TotalPrompts, b.RunningCount, b.EnqueuedCount, statusColor, b.Status, constants.ColorReset)
}

// AggregateSSHRunningPromptsRestore delegates restore across all cluster SSH nodes.
func AggregateSSHRunningPromptsRestore(opts store.RestoreOptions) error {
	localOpts := opts
	localOpts.IsSSH = false
	localErr := RunRunningPromptsRestore(localOpts)
	if SSHConnectionsFetcher == nil {
		return localErr
	}
	conns, err := SSHConnectionsFetcher()
	if err != nil || len(conns) == 0 {
		return localErr
	}
	var wg sync.WaitGroup
	for _, c := range conns {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			delegateNodeRestore(conn, opts)
		}(c)
	}
	wg.Wait()
	return nil
}

func delegateNodeRestore(c db.SSHConnection, opts store.RestoreOptions) {
	client, err := dialSSHNodeClient(c)
	if err != nil {
		printNodeSSHNotice(c.Alias, "connection failed")
		return
	}
	defer client.Close()
	cmdStr := "gitmap agy restore-running-prompts --json"
	if opts.IsKeep {
		cmdStr += " --keep"
	}
	if len(opts.TargetFile) > 0 {
		cmdStr += fmt.Sprintf(" -f %q", opts.TargetFile)
	}
	out, runErr := crypto.RunCommand(client, cmdStr, resolveNodeShell(c.OS))
	if runErr != nil {
		printNodeSSHNotice(c.Alias, "restore command failed")
		return
	}
	renderRemoteRestoreResult(c.Alias, out, opts.IsJSON)
}

func renderRemoteRestoreResult(node, out string, isJSON bool) {
	if isJSON {
		fmt.Printf("[%s] %s\n", node, strings.TrimSpace(out))
		return
	}
	fmt.Printf("  [%s] %s\n", node, strings.TrimSpace(out))
}
