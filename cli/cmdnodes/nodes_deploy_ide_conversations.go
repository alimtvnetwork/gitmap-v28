package cmdnodes

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"golang.org/x/crypto/ssh"
)

// DeployConversationsForRepo finds matching conversations locally and syncs them to remote.
func DeployConversationsForRepo(client *ssh.Client, conn db.SSHConnection, localRepoPath string) (int, error) {
	convDir, brainDir, isFound := resolveLocalGeminiDirs()
	if !isFound {
		return 0, nil
	}

	matchingIDs := findMatchingConversationIDs(convDir, brainDir, localRepoPath)
	if len(matchingIDs) == 0 {
		return 0, nil
	}

	return deployMatchingConversations(client, conn, convDir, brainDir, matchingIDs)
}

func deployMatchingConversations(client *ssh.Client, conn db.SSHConnection, convDir, brainDir string, matchingIDs []string) (int, error) {
	tarData, errPkg := packageConversationsArchive(convDir, brainDir, matchingIDs)
	if errPkg != nil {
		return 0, apperror.WrapSimple(errPkg, "package conversations failed")
	}

	return streamAndExtractConversations(client, conn, tarData, len(matchingIDs))
}

func findMatchingConversationIDs(convDir, brainDir, localRepoPath string) []string {
	entries, errRead := os.ReadDir(convDir)
	if errRead != nil {
		return nil
	}

	cleanTarget := strings.ToLower(filepath.ToSlash(filepath.Clean(localRepoPath)))
	slug := strings.ToLower(filepath.Base(cleanTarget))

	return collectMatchingIDs(entries, convDir, brainDir, cleanTarget, slug)
}

func collectMatchingIDs(entries []os.DirEntry, convDir, brainDir, cleanTarget, slug string) []string {
	var matching []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".db") {
			continue
		}

		convID := strings.TrimSuffix(e.Name(), ".db")
		if isConvMatchingRepo(convDir, brainDir, convID, cleanTarget, slug) {
			matching = append(matching, convID)
		}
	}

	return matching
}

func isConvMatchingRepo(convDir, brainDir, convID, cleanTarget, slug string) bool {
	dbPath := filepath.Join(convDir, convID+".db")
	if hasConvDBMatch(dbPath, cleanTarget, slug) {
		return true
	}

	return hasConvTranscriptMatch(brainDir, convID, cleanTarget, slug)
}

func hasConvDBMatch(dbPath, cleanTarget, slug string) bool {
	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return false
	}
	conn := openRes.Value

	defer conn.Close()
	blob := readTrajectoryBlob(conn)
	if len(blob) == 0 {
		return false
	}

	return hasTargetOrSlug(string(blob), cleanTarget, slug)
}

func readTrajectoryBlob(conn *sql.DB) []byte {
	var blob []byte
	row := conn.QueryRow("SELECT data FROM trajectory_metadata_blob WHERE id='main'")
	if row.Scan(&blob) != nil {
		return nil
	}

	return blob
}

func hasTargetOrSlug(text, cleanTarget, slug string) bool {
	lowerText := strings.ToLower(text)

	return strings.Contains(lowerText, cleanTarget) || strings.Contains(lowerText, slug)
}

func hasConvTranscriptMatch(brainDir, convID, cleanTarget, slug string) bool {
	logDir := filepath.Join(brainDir, convID, ".system_generated", "logs")
	for _, fname := range []string{"transcript_full.jsonl", "transcript.jsonl"} {
		fpath := filepath.Join(logDir, fname)
		if hasFileKeywords(fpath, cleanTarget, slug) {
			return true
		}
	}

	return false
}

func isFileReadable(fpath string) bool {
	stat, err := os.Stat(fpath)

	return err == nil && !stat.IsDir()
}

func readOpenFile(f *os.File, maxBytes int) []byte {
	buf := make([]byte, maxBytes)
	n, _ := f.Read(buf)

	return buf[:n]
}

func readFilePrefix(fpath string, maxBytes int) []byte {
	if !isFileReadable(fpath) {
		return nil
	}

	f, errOpen := os.Open(fpath)
	if errOpen != nil {
		return nil
	}

	defer f.Close()

	return readOpenFile(f, maxBytes)
}

func hasFileKeywords(fpath, cleanTarget, slug string) bool {
	data := readFilePrefix(fpath, 65536)
	if len(data) == 0 {
		return false
	}

	return hasTargetOrSlug(string(data), cleanTarget, slug)
}

func packageConversationsArchive(convDir, brainDir string, convIDs []string) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	for _, id := range convIDs {
		archiveConvDBFiles(tw, convDir, id)
		archiveBrainFolder(tw, brainDir, id)
	}

	_ = tw.Close()
	_ = gw.Close()

	return buf.Bytes(), nil
}

func archiveConvDBFiles(tw *tar.Writer, convDir, id string) {
	for _, ext := range []string{".db", ".db-wal", ".db-shm"} {
		fPath := filepath.Join(convDir, id+ext)
		stat, err := os.Stat(fPath)
		if err == nil && !stat.IsDir() {
			writeTarFileEntry(tw, "conversations/"+id+ext, fPath)
		}
	}
}

func archiveBrainFolder(tw *tar.Writer, brainDir, id string) {
	src := filepath.Join(brainDir, id)
	stat, err := os.Stat(src)
	if err != nil || !stat.IsDir() {
		return
	}

	_ = filepath.Walk(src, buildBrainWalkFn(tw, brainDir))
}

func archiveBrainPath(tw *tar.Writer, brainDir, path string, info os.FileInfo) {
	if info.IsDir() {
		return
	}

	rel, relErr := filepath.Rel(brainDir, path)
	if relErr != nil {
		return
	}

	writeTarFileEntry(tw, "brain/"+filepath.ToSlash(rel), path)
}

func buildBrainWalkFn(tw *tar.Writer, brainDir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		archiveBrainPath(tw, brainDir, path, info)

		return nil
	}
}

func streamAndExtractConversations(client *ssh.Client, conn db.SSHConnection, tarData []byte, count int) (int, error) {
	stagingPath := resolveConvStagingPath(conn)
	if errStream := cmdssh.StreamFileToRemote(client, stagingPath, tarData, conn.OS); errStream != nil {
		return 0, apperror.WrapSimple(errStream, "stream conversations failed")
	}

	return runRemoteConvExtract(client, conn, stagingPath, count)
}

func runRemoteConvExtract(client *ssh.Client, conn db.SSHConnection, stagingPath string, count int) (int, error) {
	cmd, shell := resolveConvExtractCommand(conn, stagingPath)
	if _, errExec := secrets.RunCommand(client, cmd, shell); errExec != nil {
		return 0, apperror.WrapSimple(errExec, "extract conversations failed")
	}

	return count, nil
}

func resolveConvStagingPath(conn db.SSHConnection) string {
	if isWindowsNode(conn) {
		return `C:\Windows\Temp\convs_deploy.tar.gz`
	}

	return "/tmp/convs_deploy.tar.gz"
}

func resolveConvExtractCommand(conn db.SSHConnection, stagingPath string) (string, string) {
	if isWindowsNode(conn) {
		cmd := fmt.Sprintf(`powershell.exe -NoProfile -Command "$dest = Join-Path $env:USERPROFILE '.gemini\antigravity'; if (-not (Test-Path $dest)) { New-Item -ItemType Directory -Force -Path $dest | Out-Null }; tar.exe -xzf '%s' -C $dest; Remove-Item '%s' -Force -ErrorAction SilentlyContinue"`, stagingPath, stagingPath)

		return cmd, "ps"
	}

	cmd := fmt.Sprintf(`sh -c "mkdir -p ~/.gemini/antigravity && tar -xzf %s -C ~/.gemini/antigravity && rm -f %s"`, stagingPath, stagingPath)

	return cmd, "sh"
}
