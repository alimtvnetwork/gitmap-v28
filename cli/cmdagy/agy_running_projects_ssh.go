package cmdagy

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"golang.org/x/crypto/ssh"
)

// SSHConnectionsFetcher is an injectable function to query cluster SSH connections without package cycles.
var SSHConnectionsFetcher func() ([]db.SSHConnection, error)

// AggregateSSHRunningProjects queries remote SSH cluster nodes and aggregates with local.
func AggregateSSHRunningProjects(local []RunningProjectRecord, isJSON bool, filePath string) error {
	remoteRecords := fetchClusterSSHRunningProjects()
	allRecords := append(local, remoteRecords...)
	if filePath != "" {
		return writeRunningProjectsToFile(allRecords, filePath)
	}
	if isJSON {
		return renderRunningProjectsJSON(allRecords)
	}
	renderSSHRunningProjectsTable(allRecords)
	return nil
}

func fetchClusterSSHRunningProjects() []RunningProjectRecord {
	if SSHConnectionsFetcher == nil {
		return nil
	}
	conns, err := SSHConnectionsFetcher()
	if err != nil || len(conns) == 0 {
		return nil
	}
	var aggregated []RunningProjectRecord
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, c := range conns {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			records := querySingleSSHNodeRunningProjects(conn)
			mu.Lock()
			aggregated = append(aggregated, records...)
			mu.Unlock()
		}(c)
	}
	wg.Wait()
	return aggregated
}

func querySingleSSHNodeRunningProjects(c db.SSHConnection) []RunningProjectRecord {
	client, err := dialSSHNodeClient(c)
	if err != nil {
		return []RunningProjectRecord{buildOfflineNodeRecord(c)}
	}
	defer client.Close()
	cmdStr := "gitmap agy running-projects --json"
	out, runErr := crypto.RunCommand(client, cmdStr, resolveNodeShell(c.OS))
	if runErr != nil {
		return []RunningProjectRecord{buildNodeErrorRecord(c, runErr)}
	}
	return parseRemoteRunningProjects(c, out)
}

// SSHNodeDialer is an injectable function to dial cluster SSH connections with full credential fallbacks.
var SSHNodeDialer func(c db.SSHConnection) (*ssh.Client, error)

func dialSSHNodeClient(c db.SSHConnection) (*ssh.Client, error) {
	if SSHNodeDialer != nil {
		return SSHNodeDialer(c)
	}
	plainPass, _ := crypto.DecryptStoredPassword(c.EncryptedPassword)
	return crypto.ConnectWithFallback(c.IPAddress, c.Username, c.KeyPath, plainPass)
}

func resolveNodeShell(osName string) string {
	low := strings.ToLower(osName)
	if strings.Contains(low, "win") {
		return "cmd"
	}
	return "sh"
}

func buildOfflineNodeRecord(c db.SSHConnection) RunningProjectRecord {
	return RunningProjectRecord{
		ProjectName:   "(node unreachable)",
		Node:          c.Alias,
		Host:          c.IPAddress,
		Status:        "OFFLINE",
		PromptPreview: "SSH connection failed or node unreachable",
	}
}

func buildNodeErrorRecord(c db.SSHConnection, err error) RunningProjectRecord {
	return RunningProjectRecord{
		ProjectName:   "(remote error)",
		Node:          c.Alias,
		Host:          c.IPAddress,
		Status:        "ERROR",
		PromptPreview: CompactWords(err.Error(), 8),
	}
}

func parseRemoteRunningProjects(c db.SSHConnection, out string) []RunningProjectRecord {
	jsonStart := strings.Index(out, "[")
	jsonEnd := strings.LastIndex(out, "]")
	if jsonStart < 0 || jsonEnd <= jsonStart {
		return nil
	}
	var list []RunningProjectRecord
	if err := json.Unmarshal([]byte(out[jsonStart:jsonEnd+1]), &list); err != nil {
		return nil
	}
	for i := range list {
		list[i].Node = c.Alias
		list[i].Host = c.IPAddress
	}
	return list
}

func renderSSHRunningProjectsTable(records []RunningProjectRecord) {
	fmt.Println()
	fmt.Printf("  %s%s CLUSTER & LOCAL RUNNING PROJECTS (%d discovered) %s%s\n",
		constants.ColorCyan, "╔════", len(records), "════╗", constants.ColorReset)
	fmt.Printf("  %s%-18s  %-12s  %-15s  %-14s  %-8s  %s%s\n",
		constants.ColorWhite, "PROJECT", "NODE", "HOST", "STATUS", "QUEUED", "PREVIEW", constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("─", 95), constants.ColorReset)
	if len(records) == 0 {
		fmt.Printf("  %sNo active or queued projects found across cluster.%s\n\n", constants.ColorDim, constants.ColorReset)
		return
	}
	for _, r := range records {
		printSSHRunningProjectRow(r)
	}
	fmt.Println()
}

func printSSHRunningProjectRow(r RunningProjectRecord) {
	node := resolveFallbackValue(r.Node, "local")
	host := resolveFallbackValue(r.Host, "127.0.0.1")
	color := resolveRunningColorSSH(r.HasActivePrompt)
	preview := truncateString(r.PromptPreview, 30)
	fmt.Printf("  %-18s  %-12s  %-15s  %s%-14s%s  %-8d  %s\n",
		truncateString(r.ProjectName, 18), node, host, color, r.Status, constants.ColorReset, r.QueuedCount, preview)
}

func resolveFallbackValue(val, fallback string) string {
	if val != "" {
		return val
	}
	return fallback
}

func resolveRunningColorSSH(hasActive bool) string {
	if hasActive {
		return constants.ColorGreen + "\033[1m"
	}
	return constants.ColorYellow
}
