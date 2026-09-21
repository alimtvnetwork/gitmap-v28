package cmdssh

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func getAskPassScriptExt() string {
	isWindows := runtime.GOOS == "windows"
	if isWindows {
		return ".bat"
	}
	return ".sh"
}

func buildAskPassScriptContent() string {
	isWindows := runtime.GOOS == "windows"
	if isWindows {
		return "@echo off\r\n<nul set /p=\"%GITMAP_SSH_PASS%\"\r\n"
	}
	return "#!/bin/sh\nprintf \"%s\" \"$GITMAP_SSH_PASS\"\n"
}

func resolveAskPassTempPath() string {
	ts := time.Now().UnixNano()
	filename := fmt.Sprintf("gitmap-askpass-%d%s", ts, getAskPassScriptExt())
	baseDir := filepath.Join(os.TempDir(), "gitmap", "ssh")
	_ = os.MkdirAll(baseDir, 0700)

	return filepath.Join(baseDir, filename)
}

func writeAskPassScript(targetPath string, content string) error {
	err := os.WriteFile(targetPath, []byte(content), 0700)
	if err != nil {
		return apperror.WrapSimple(err, "writeAskPassScript")
	}
	return nil
}

// CreateAskPassScript creates a platform-appropriate helper script for SSH_ASKPASS.
func CreateAskPassScript() (string, func(), error) {
	scriptPath := resolveAskPassTempPath()
	content := buildAskPassScriptContent()
	if err := writeAskPassScript(scriptPath, content); err != nil {
		return "", func() {}, err
	}
	cleanup := func() {
		_ = os.Remove(scriptPath)
	}
	return scriptPath, cleanup, nil
}

// BuildAskPassEnv constructs process environment variables ensuring OpenSSH uses AskPass.
func BuildAskPassEnv(baseEnv []string, scriptPath, password string) []string {
	env := append([]string{}, baseEnv...)
	env = append(env, fmt.Sprintf("SSH_ASKPASS=%s", scriptPath))
	env = append(env, "SSH_ASKPASS_REQUIRE=force")
	env = append(env, "DISPLAY=1")
	env = append(env, fmt.Sprintf("GITMAP_SSH_PASS=%s", password))
	return env
}
