package cmdagy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

// RunAGYRestore deserializes an exported snapshot JSON and re-enqueues projects and conversations.
func RunAGYRestore(snapPath string) (int, int, error) {
	resolvedPath, pathErr := resolveRestorePath(snapPath)
	hasPathErr := pathErr != nil
	if hasPathErr {
		return 0, 0, pathErr
	}

	snap, readErr := readSnapshotPayload(resolvedPath)
	hasReadErr := readErr != nil
	if hasReadErr {
		return 0, 0, readErr
	}

	restoredProj := restoreProjectsFromSnapshot(snap.Projects)
	restoredConv := restoreConversationsFromSnapshot(snap.Conversations)

	fmt.Printf("%s[GitMap]%s Successfully restored %d project(s) and %d conversation(s) from:\n  %s\n",
		constants.ColorGreen, constants.ColorReset, restoredProj, restoredConv, resolvedPath)

	return restoredProj, restoredConv, nil
}

func resolveRestorePath(snapPath string) (string, error) {
	hasCustom := snapPath != ""
	if hasCustom {
		return filepath.Abs(snapPath)
	}
	return findLatestSnapshotPath()
}

func findLatestSnapshotPath() (string, error) {
	home, err := os.UserHomeDir()
	hasErr := err != nil
	if hasErr {
		return "", err
	}
	snapDir := filepath.Join(home, ".gitmap")
	entries, readErr := os.ReadDir(snapDir)
	hasReadErr := readErr != nil
	if hasReadErr {
		return "", fmt.Errorf("no snapshot directory found at %s: %w", snapDir, readErr)
	}
	var snaps []string
	for _, e := range entries {
		isMatch := !e.IsDir() && strings.HasPrefix(e.Name(), "agy-snapshot-") && strings.HasSuffix(e.Name(), ".json")
		if isMatch {
			snaps = append(snaps, filepath.Join(snapDir, e.Name()))
		}
	}
	hasNone := len(snaps) == 0
	if hasNone {
		return "", fmt.Errorf("no AGY restore snapshot JSON found in %s", snapDir)
	}
	sort.Strings(snaps)
	return snaps[len(snaps)-1], nil
}

func readSnapshotPayload(path string) (AGYRestoreSnapshot, error) {
	data, err := os.ReadFile(path)
	hasErr := err != nil
	if hasErr {
		return AGYRestoreSnapshot{}, err
	}
	var snap AGYRestoreSnapshot
	unmarshalErr := json.Unmarshal(data, &snap)
	hasUnmarshalErr := unmarshalErr != nil
	if hasUnmarshalErr {
		return AGYRestoreSnapshot{}, unmarshalErr
	}
	return snap, nil
}

func restoreProjectsFromSnapshot(projects []AGYSnapshotProject) int {
	dirPath, err := getProjectsDirPath()
	hasErr := err != nil
	if hasErr {
		return 0
	}
	_ = os.MkdirAll(dirPath, 0755)

	count := 0
	for _, p := range projects {
		isOk := restoreSingleProject(p, dirPath)
		if isOk {
			count++
		}
	}
	return count
}

func restoreSingleProject(p AGYSnapshotProject, dir string) bool {
	hasEmpty := p.ID == ""
	if hasEmpty {
		return false
	}
	projObj := buildProjectObject(p)
	data, marshalErr := json.MarshalIndent(projObj, "", "  ")
	hasMarshalErr := marshalErr != nil
	if hasMarshalErr {
		return false
	}
	outPath := filepath.Join(dir, p.ID+".json")
	writeErr := os.WriteFile(outPath, data, 0644)
	return writeErr == nil
}

func buildProjectObject(p AGYSnapshotProject) AgyProject {
	branch := p.Branch
	hasBranch := branch != ""
	if !hasBranch {
		branch = "main"
	}
	return AgyProject{
		ID:        p.ID,
		Name:      p.Name,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		ProjectResources: &AgyProjectResources{
			Resources: []AgyResource{
				{
					GitFolder: &AgyGitFolder{
						FolderURI:     pathToFolderURI(p.Workspace),
						DefaultBranch: branch,
					},
				},
			},
		},
	}
}

func pathToFolderURI(path string) string {
	clean := filepath.ToSlash(filepath.Clean(path))
	return "file:///" + clean
}

func restoreConversationsFromSnapshot(convs []AGYSnapshotConv) int {
	dbPath, err := getConversationSummariesDBPath()
	hasErr := err != nil
	if hasErr {
		return 0
	}
	_ = os.MkdirAll(filepath.Dir(dbPath), 0755)
	conn, openErr := store.OpenSQLiteDB(dbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return 0
	}
	defer conn.Close()

	initConvSummariesTable(conn)
	return insertConvRows(conn, convs)
}

func initConvSummariesTable(conn *sql.DB) {
	stmt := `CREATE TABLE IF NOT EXISTS conversation_summaries (
		conversation_id TEXT PRIMARY KEY,
		title TEXT,
		workspace_uris TEXT,
		project_id TEXT,
		step_count INTEGER,
		last_modified_time TIMESTAMP
	);`
	_, _ = conn.Exec(stmt)
}

func insertConvRows(conn *sql.DB, convs []AGYSnapshotConv) int {
	count := 0
	query := `INSERT OR REPLACE INTO conversation_summaries
		(conversation_id, title, workspace_uris, project_id, step_count, last_modified_time)
		VALUES (?, ?, ?, ?, ?, ?);`
	for _, c := range convs {
		hasEmpty := c.ID == ""
		if hasEmpty {
			continue
		}
		uri := pathToFolderURI(c.WorkspacePath)
		_, err := conn.Exec(query, c.ID, c.Title, uri, c.ProjectID, c.StepCount, time.Now())
		if err == nil {
			count++
		}
	}
	return count
}

var agyRestoreCmd = &cobra.Command{
	Use:     "restore [snapshot-path]",
	Aliases: []string{"reenq", "re-enqueue", "import-snapshot"},
	Short:   "Restore Antigravity projects and conversations from snapshot JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		snapPath := ""
		hasArg := len(args) > 0
		if hasArg {
			snapPath = args[0]
		}
		_, _, err := RunAGYRestore(snapPath)
		return err
	},
}

func init() {
	AgyCmd.AddCommand(agyRestoreCmd)
}
