package cmdui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

var syntaxLangMap = map[string]string{
	".go":   "go",
	".json": "json",
	".md":   "markdown",
	".js":   "javascript",
	".ts":   "javascript",
	".jsx":  "javascript",
	".tsx":  "javascript",
	".html": "html",
	".htm":  "html",
	".css":  "css",
	".py":   "python",
	".sh":   "shell",
	".bash": "shell",
	".ps1":  "shell",
}

// ReadRemoteOrLocalFile reads file content locally or from an SSH node.
func ReadRemoteOrLocalFile(nodeAlias, filePath string) (RemoteFileContent, error) {
	isLocal := nodeAlias == "" || strings.EqualFold(nodeAlias, "local")

	if isLocal {
		return readLocalFile(filePath)
	}

	return readRemoteSSHFile(nodeAlias, filePath)
}

func readLocalFile(filePath string) (RemoteFileContent, error) {
	data, err := os.ReadFile(filePath)

	if err != nil {
		return RemoteFileContent{FilePath: filePath, IsSuccess: false, Error: err.Error()}, apperror.WrapSimple(err, "read_local_file")
	}

	return RemoteFileContent{
		NodeAlias: "local",
		FilePath:  filePath,
		Content:   string(data),
		Language:  detectSyntaxLanguage(filePath),
		IsSuccess: true,
	}, nil
}

func readRemoteSSHFile(nodeAlias, filePath string) (RemoteFileContent, error) {
	conn, hasConn := findSSHNodeByAlias(nodeAlias)

	if hasConn == false {
		return RemoteFileContent{NodeAlias: nodeAlias, FilePath: filePath, IsSuccess: false, Error: "node not found"}, apperror.NewSimple("read_remote", "node_not_found")
	}

	return executeRemoteRead(conn, nodeAlias, filePath)
}

func executeRemoteRead(conn db.SSHConnection, nodeAlias, filePath string) (RemoteFileContent, error) {
	client, isConnected := cmdssh.ConnectSSHClient(conn, fmt.Sprintf("[%s]", nodeAlias))

	if isConnected == false {
		return RemoteFileContent{NodeAlias: nodeAlias, FilePath: filePath, IsSuccess: false, Error: "ssh connect failed"}, apperror.NewSimple("read_remote", "connect_failed")
	}
	defer client.Close()

	cmd := fmt.Sprintf("cat %q 2>/dev/null || type %q", filePath, filePath)
	out, runErr := crypto.RunCommand(client, cmd, "sh")

	if runErr != nil {
		return RemoteFileContent{NodeAlias: nodeAlias, FilePath: filePath, IsSuccess: false, Error: runErr.Error()}, apperror.WrapSimple(runErr, "remote_read")
	}

	return RemoteFileContent{
		NodeAlias: nodeAlias,
		FilePath:  filePath,
		Content:   out,
		Language:  detectSyntaxLanguage(filePath),
		IsSuccess: true,
	}, nil
}

// SaveRemoteOrLocalFile saves file content locally or to an SSH node.
func SaveRemoteOrLocalFile(nodeAlias, filePath, content string) error {
	isLocal := nodeAlias == "" || strings.EqualFold(nodeAlias, "local")

	if isLocal {
		return os.WriteFile(filePath, []byte(content), 0644)
	}

	return saveRemoteSSHFile(nodeAlias, filePath, content)
}

func saveRemoteSSHFile(nodeAlias, filePath, content string) error {
	conn, hasConn := findSSHNodeByAlias(nodeAlias)

	if hasConn == false {
		return apperror.NewSimple("save_remote", "node_not_found")
	}

	return executeRemoteSave(conn, nodeAlias, filePath, content)
}

func executeRemoteSave(conn db.SSHConnection, nodeAlias, filePath, content string) error {
	client, isConnected := cmdssh.ConnectSSHClient(conn, fmt.Sprintf("[%s]", nodeAlias))

	if isConnected == false {
		return apperror.NewSimple("save_remote", "connect_failed")
	}
	defer client.Close()

	escapedContent := strings.ReplaceAll(content, "'", "'\\''")
	cmd := fmt.Sprintf("printf '%%s' '%s' > %q", escapedContent, filePath)
	_, runErr := crypto.RunCommand(client, cmd, "sh")

	if runErr != nil {
		return apperror.WrapSimple(runErr, "remote_save")
	}

	return nil
}

func findSSHNodeByAlias(alias string) (db.SSHConnection, bool) {
	conns, err := db.LoadAllSSHConnections()

	if err != nil {
		return db.SSHConnection{}, false
	}

	for _, c := range conns {
		if strings.EqualFold(c.Alias, alias) || c.IPAddress == alias {
			return c, true
		}
	}

	return db.SSHConnection{}, false
}

func detectSyntaxLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	lang, hasLang := syntaxLangMap[ext]

	if hasLang {
		return lang
	}

	return "plaintext"
}
