package cmdssh

import (
	"archive/tar"
	"bytes"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// StreamFileToRemote streams binary data to a remote host over an SSH channel.
// It prioritizes streaming tar archives directly to stdin for extraction,
// falling back to direct stdin file piping if remote tar is unavailable.
func StreamFileToRemote(client *ssh.Client, remotePath string, data []byte, osType string) error {
	isWin := isWindowsOS(osType)
	dir, filename := splitRemoteDestPath(remotePath, isWin)

	errTar := streamTarToRemote(client, dir, filename, data, isWin)
	if errTar == nil {
		return nil
	}

	return streamDirectToRemote(client, remotePath, data, isWin)
}

func splitWindowsRemoteDestPath(remotePath string) (string, string) {
	norm := strings.ReplaceAll(remotePath, "/", "\\")
	idx := strings.LastIndex(norm, "\\")
	if idx >= 0 {
		return norm[:idx], norm[idx+1:]
	}
	return "C:\\Windows\\Temp", norm
}

func splitRemoteDestPath(remotePath string, isWin bool) (string, string) {
	if isWin {
		return splitWindowsRemoteDestPath(remotePath)
	}

	idx := strings.LastIndex(remotePath, "/")
	if idx >= 0 {
		return remotePath[:idx], remotePath[idx+1:]
	}
	return "/tmp", remotePath
}

func ensureRemoteDir(client *ssh.Client, dir string, isWin bool) error {
	if dir == "" || dir == "." {
		return nil
	}
	if isWin {
		cmd := fmt.Sprintf("powershell -NoProfile -Command \"if (-not (Test-Path '%s')) { New-Item -ItemType Directory -Force -Path '%s' | Out-Null }\"", dir, dir)
		_, err := crypto.RunCommand(client, cmd, "")
		return err
	}
	cmd := fmt.Sprintf("mkdir -p '%s'", dir)
	_, err := crypto.RunCommand(client, cmd, "")
	return err
}

func buildTarExtractCmd(dir string, isWin bool) string {
	if isWin {
		return fmt.Sprintf("cmd.exe /c tar.exe -xf - -C \"%s\"", dir)
	}
	return fmt.Sprintf("sh -c \"tar -xf - -C '%s'\"", dir)
}

func buildDirectStreamCmd(remotePath string, isWin bool) string {
	if isWin {
		return fmt.Sprintf("powershell -NoProfile -Command \"$p = '%s'; $dir = [System.IO.Path]::GetDirectoryName($p); if ($dir -and -not (Test-Path $dir)) { [System.IO.Directory]::CreateDirectory($dir) | Out-Null }; $in = [System.Console]::OpenStandardInput(); $out = [System.IO.File]::Create($p); $in.CopyTo($out); $out.Close(); $in.Close()\"", remotePath)
	}
	return fmt.Sprintf("sh -c \"dir=$(dirname '%s'); mkdir -p \\\"$dir\\\" && cat > '%s' && chmod 755 '%s'\"", remotePath, remotePath, remotePath)
}

func streamTarToRemote(client *ssh.Client, dir, filename string, data []byte, isWin bool) error {
	_ = ensureRemoteDir(client, dir, isWin)

	session, err := client.NewSession()
	if err != nil {
		return apperror.WrapSimple(err, "streamTarToRemote.NewSession")
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return apperror.WrapSimple(err, "streamTarToRemote.StdinPipe")
	}

	var stderr bytes.Buffer
	session.Stderr = &stderr
	cmd := buildTarExtractCmd(dir, isWin)
	if err := session.Start(cmd); err != nil {
		return apperror.WrapSimple(err, "streamTarToRemote.Start")
	}

	tw := tar.NewWriter(stdin)
	hdr := &tar.Header{
		Name:    filename,
		Mode:    0755,
		Size:    int64(len(data)),
		ModTime: time.Now(),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		_ = stdin.Close()
		return apperror.WrapSimple(err, "streamTarToRemote.WriteHeader")
	}
	if _, err := tw.Write(data); err != nil {
		_ = stdin.Close()
		return apperror.WrapSimple(err, "streamTarToRemote.Write")
	}
	_ = tw.Close()
	_ = stdin.Close()

	if err := session.Wait(); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("streamTarToRemote.Wait: %s", stderr.String()))
	}
	return nil
}

func streamDirectToRemote(client *ssh.Client, remotePath string, data []byte, isWin bool) error {
	session, err := client.NewSession()
	if err != nil {
		return apperror.WrapSimple(err, "streamDirectToRemote.NewSession")
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return apperror.WrapSimple(err, "streamDirectToRemote.StdinPipe")
	}

	var stderr bytes.Buffer
	session.Stderr = &stderr
	cmd := buildDirectStreamCmd(remotePath, isWin)
	if err := session.Start(cmd); err != nil {
		return apperror.WrapSimple(err, "streamDirectToRemote.Start")
	}

	if _, err := stdin.Write(data); err != nil {
		_ = stdin.Close()
		return apperror.WrapSimple(err, "streamDirectToRemote.Write")
	}
	_ = stdin.Close()

	if err := session.Wait(); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("streamDirectToRemote.Wait: %s", stderr.String()))
	}
	return nil
}

// ProbeRemoteOSTypeForTest exports probeRemoteOSType for testing.
func ProbeRemoteOSTypeForTest(client *ssh.Client) string {
	return probeRemoteOSType(client)
}

// FetchAllSSHConnectionsForTest exports fetchAllSSHConnections for testing.
func FetchAllSSHConnectionsForTest() ([]db.SSHConnection, error) {
	return fetchAllSSHConnections()
}

// ConnectSSHClientForTest exports connectSSHClient for testing.
func ConnectSSHClientForTest(c db.SSHConnection) (*ssh.Client, bool) {
	return connectSSHClient(c)
}
