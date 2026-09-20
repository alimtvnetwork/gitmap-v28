package cmdagy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

// RunAgyOpen opens Google Antigravity Desktop IDE on the given path.
func RunAgyOpen(path string) error {
	absPath := resolveTargetPath(path)
	_ = workspacesync.SyncAntigravity(absPath, filepath.Base(absPath))
	backupAgyLocation(absPath)
	ideRes := ResolveAntigravityIDE()
	if ideRes.IsFailure() {
		printAntigravityInstallSuggestion()

		return nil
	}

	return startAntigravity(ideRes.Value, absPath)
}

func startAntigravity(exePath, absPath string) error {
	if err := launchAntigravityProcess(exePath, absPath); err != nil {
		fmt.Printf("Error launching Antigravity: %v\n", err)

		return err
	}
	fmt.Printf("✓ Opened Google Antigravity in %s\n", absPath)

	return nil
}

func resolveCwdOrPath(path string) string {
	cwd, err := os.Getwd()
	if err == nil {
		return cwd
	}

	return path
}

func resolveTargetPath(path string) string {
	if path == "" || path == "." {
		return resolveCwdOrPath(path)
	}
	abs, err := filepath.Abs(path)
	if err == nil {
		return abs
	}

	return path
}

func backupAgyLocation(absPath string) {
	db, err := store.OpenDefault()
	if err != nil {
		return
	}
	defer db.Close()

	schema := "CREATE TABLE IF NOT EXISTS agy_locations (path TEXT PRIMARY KEY, opened_at TEXT NOT NULL);"
	_, _ = store.ExecWrapper(db.Conn(), schema).Destruct()
	now := time.Now().UTC().Format(time.RFC3339)
	upsert := "INSERT INTO agy_locations (path, opened_at) VALUES (?, ?) ON CONFLICT(path) DO UPDATE SET opened_at = excluded.opened_at;"
	_, _ = store.ExecWrapper(db.Conn(), upsert, absPath, now).Destruct()
}

func printAntigravityInstallSuggestion() {
	fmt.Println("✖ Google Antigravity Desktop IDE is not installed.")
	fmt.Println("  To install it, run:")
	fmt.Println("    gitmap install antigravity")
	fmt.Println("  Or download from:")
	fmt.Println("    https://antigravity.google/download")
}

func launchAntigravityProcess(exePath, absPath string) error {
	cmd := buildAntigravityCmd(exePath, absPath)

	return cmd.Start()
}

func buildAntigravityCmd(exePath, absPath string) *exec.Cmd {
	cleanPath := filepath.Clean(exePath)
	isMacApp := runtime.GOOS == "darwin" && strings.HasSuffix(cleanPath, ".app")
	if isMacApp {
		return buildMacAppCmd(cleanPath, absPath)
	}

	return exec.Command(cleanPath, absPath)
}

func buildMacAppCmd(appPath, absPath string) *exec.Cmd {
	innerBin := filepath.Join(appPath, "Contents", "MacOS", "Antigravity")
	hasInnerBin := isValidBinaryPath(innerBin)
	if hasInnerBin {
		return exec.Command(innerBin, absPath)
	}

	return exec.Command("open", "-a", appPath, absPath)
}
