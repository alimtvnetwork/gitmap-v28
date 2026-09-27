package cmdagy

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// DiscoverUnifiedAgyProjects aggregates projects across config JSON files, brain DB summaries, and workspaces.
func DiscoverUnifiedAgyProjects(dirPath string) ([]AgyProject, error) {
	fileProjects, _ := loadAllAgyProjects(dirPath)
	seen := make(map[string]bool)
	var projects []AgyProject
	for _, p := range fileProjects {
		trackAgyProject(&projects, seen, p)
	}
	dbProjects := discoverProjectsFromSummariesDB()
	for _, p := range dbProjects {
		trackAgyProject(&projects, seen, p)
	}
	candProjects := discoverProjectsFromCandidates()
	for _, p := range candProjects {
		trackAgyProject(&projects, seen, p)
	}
	return projects, nil
}

func trackAgyProject(projects *[]AgyProject, seen map[string]bool, p AgyProject) {
	ws := cleanProjectWorkspace(p.GetPath())
	if ws == "" {
		return
	}
	if seen[ws] {
		return
	}
	seen[ws] = true
	*projects = append(*projects, p)
}

func discoverProjectsFromSummariesDB() []AgyProject {
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return nil
	}
	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return nil
	}
	defer conn.Close()
	return queryProjectsFromSummariesConn(conn)
}

func queryProjectsFromSummariesConn(conn *sql.DB) []AgyProject {
	query := `SELECT project_id, title, workspace_uris, MAX(last_modified_time)
		FROM conversation_summaries
		WHERE (killed IS NULL OR killed = 0) AND workspace_uris IS NOT NULL
		GROUP BY workspace_uris ORDER BY MAX(last_modified_time) DESC`
	rows, err := conn.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanSummariesProjects(rows)
}

func scanSummariesProjects(rows *sql.Rows) []AgyProject {
	var projects []AgyProject
	for rows.Next() {
		var pid, title, wsRaw, maxTime sql.NullString
		if err := rows.Scan(&pid, &title, &wsRaw, &maxTime); err != nil {
			continue
		}
		p, isSynthesized := synthesizeProjectFromRow(pid, title, wsRaw, maxTime)
		if isSynthesized {
			projects = append(projects, p)
		}
	}
	return projects
}

func synthesizeProjectFromRow(pid, title, wsRaw, maxTime sql.NullString) (AgyProject, bool) {
	ws := extractCleanWorkspaceFromURIs(wsRaw.String)
	if ws == "" || !checkDirExists(ws) {
		return AgyProject{}, false
	}
	name := filepath.Base(ws)
	if title.Valid && title.String != "" && !strings.Contains(title.String, "Read & Understand") {
		name = title.String
	}
	id := resolveProjectID(pid, ws)
	updated := resolveProjectUpdatedTime(maxTime)
	return createSynthesizedAgyProject(ws, name, id, updated), true
}

func resolveProjectID(pid sql.NullString, ws string) string {
	if pid.Valid && pid.String != "" {
		return pid.String
	}
	h := sha256.Sum256([]byte(strings.ToLower(ws)))
	return hex.EncodeToString(h[:8])
}

func resolveProjectUpdatedTime(maxTime sql.NullString) string {
	if maxTime.Valid && maxTime.String != "" {
		return maxTime.String
	}
	return time.Now().UTC().Format(time.RFC3339)
}

func createSynthesizedAgyProject(ws, name, id, updated string) AgyProject {
	folderURI := "file:///" + filepath.ToSlash(ws)
	return AgyProject{
		ID:        id,
		Name:      name,
		UpdatedAt: updated,
		ProjectResources: &AgyProjectResources{
			Resources: []AgyResource{
				{GitFolder: &AgyGitFolder{FolderURI: folderURI, DefaultBranch: "main"}},
			},
		},
	}
}

func discoverProjectsFromCandidates() []AgyProject {
	cwd, err := os.Getwd()
	if err != nil || cwd == "" || !checkDirExists(cwd) {
		return nil
	}
	h := sha256.Sum256([]byte(strings.ToLower(cwd)))
	id := hex.EncodeToString(h[:8])
	p := createSynthesizedAgyProject(cwd, filepath.Base(cwd), id, time.Now().UTC().Format(time.RFC3339))
	return []AgyProject{p}
}

// FetchClusterSSHProjects queries remote cluster nodes and parses their Antigravity projects.
func FetchClusterSSHProjects() []AgyProject {
	if SSHConnectionsFetcher == nil {
		return nil
	}
	conns, err := SSHConnectionsFetcher()
	if err != nil || len(conns) == 0 {
		return nil
	}
	var aggregated []AgyProject
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, c := range conns {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			projs := querySingleNodeProjects(conn)
			mu.Lock()
			aggregated = append(aggregated, projs...)
			mu.Unlock()
		}(c)
	}
	wg.Wait()
	return aggregated
}

func probeNodeOnline(target string, timeout time.Duration) bool {
	host := strings.TrimSpace(target)
	if host == "" {
		return false
	}
	if !strings.Contains(host, ":") {
		host = net.JoinHostPort(host, "22")
	}
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func querySingleNodeProjects(c db.SSHConnection) []AgyProject {
	if !probeNodeOnline(c.IPAddress, 400*time.Millisecond) {
		return nil
	}
	client, err := dialSSHNodeClient(c)
	if err != nil {
		return nil
	}
	defer client.Close()
	cmdStr := "gitmap agy ls --json"
	out, runErr := crypto.RunCommand(client, cmdStr, resolveNodeShell(c.OS))
	if runErr != nil {
		return nil
	}
	return parseRemoteAgyProjectsJSON(out)
}

func parseRemoteAgyProjectsJSON(raw string) []AgyProject {
	trimmed := strings.TrimSpace(raw)
	startIdx := strings.Index(trimmed, "[")
	endIdx := strings.LastIndex(trimmed, "]")
	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		return nil
	}
	jsonStr := trimmed[startIdx : endIdx+1]
	var projects []AgyProject
	if err := json.Unmarshal([]byte(jsonStr), &projects); err != nil {
		return nil
	}
	return projects
}
