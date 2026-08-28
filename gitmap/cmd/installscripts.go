package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// scriptsConfig holds the powershell.json structure for deploy path resolution.
type scriptsConfig struct {
	DeployPath string `json:"deployPath"`
}

type scriptSource struct {
	src  string
	name string
}

func defaultScriptSources(tmpDir string) []scriptSource {
	return []scriptSource{
		{filepath.Join(tmpDir, "gitmap", "scripts", "install.ps1"), "install.ps1"},
		{filepath.Join(tmpDir, "gitmap", "scripts", "install.sh"), "install.sh"},
		{filepath.Join(tmpDir, "gitmap", "scripts", "uninstall.ps1"), "uninstall.ps1"},
		{filepath.Join(tmpDir, "gitmap", "scripts", "Get-LastRelease.ps1"), "Get-LastRelease.ps1"},
		{filepath.Join(tmpDir, "run.ps1"), "run.ps1"},
		{filepath.Join(tmpDir, "run.sh"), "run.sh"},
	}
}

func cloneRepoToTemp() (string, error) {
	repoURL := constants.PrefixHTTPS + constants.GitmapRepoPrefix + constants.ExtGit
	fmt.Printf(constants.MsgScriptsCloning, repoURL)

	tmpDir, err := os.MkdirTemp("", "gitmap-scripts-clone-*")
	if err != nil {
		return "", err
	}

	cloneCmd := exec.Command("git", "clone", "--depth", "1", repoURL, tmpDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr

	if err := cloneCmd.Run(); err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", err
	}
	return tmpDir, nil
}

func copySingleScript(s scriptSource, targetDir string) bool {
	data, err := os.ReadFile(s.src)
	if err != nil {
		fmt.Printf(constants.MsgScriptsSkip, s.name)
		return false
	}
	dest := filepath.Join(targetDir, s.name)
	if err := os.WriteFile(dest, data, constants.DirPermission); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrScriptsCopy, s.name, err)
		return false
	}
	fmt.Printf(constants.MsgScriptsCopied, s.name)
	return true
}

func copyScriptFiles(tmpDir, targetDir string) int {
	copied := 0
	for _, s := range defaultScriptSources(tmpDir) {
		if copySingleScript(s, targetDir) {
			copied++
		}
	}
	return copied
}

// runInstallScripts clones/copies gitmap scripts to a platform-specific folder.
func runInstallScripts() error {
	targetDir := resolveScriptsDir()
	fmt.Printf(constants.MsgScriptsTarget, targetDir)
	if err := os.MkdirAll(targetDir, constants.DirPermission); err != nil {
		return apperror.NewSimple(constants.ErrScriptsMkdir, "E9000")
	}
	tmpDir, err := cloneRepoToTemp()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrScriptsClone, err)
		exitWith(1)
	}
	defer os.RemoveAll(tmpDir)
	copied := copyScriptFiles(tmpDir, targetDir)
	fmt.Println()
	fmt.Printf(constants.MsgScriptsDone, copied, targetDir)
	return nil
}

// resolveScriptsDir returns the target directory for scripts.
// Windows: reads deployPath from powershell.json, defaults to D:\gitmap-scripts.
// Linux/macOS: ~/Desktop/gitmap-scripts.
func resolveScriptsDir() string {
	if runtime.GOOS == constants.PlatformWindows {
		return resolveScriptsDirWindows()
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return filepath.Join(home, "Desktop", "gitmap-scripts")
}

func scriptsCandidatePaths() []string {
	candidates := []string{}
	exe, err := os.Executable()
	resolved, evalErr := filepath.EvalSymlinks(exe)
	if err == nil && evalErr == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(resolved), constants.PowershellConfigFile))
	}
	if err == nil && evalErr != nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), constants.PowershellConfigFile))
	}
	candidates = append(candidates,
		filepath.Join("gitmap", constants.PowershellConfigFile),
		constants.PowershellConfigFile,
	)
	return candidates
}

func resolveDriveFromConfig(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var cfg scriptsConfig
	if err := json.Unmarshal(data, &cfg); err != nil || cfg.DeployPath == "" {
		return ""
	}
	return filepath.VolumeName(cfg.DeployPath)
}

// resolveScriptsDirWindows reads powershell.json for the deploy drive.
func resolveScriptsDirWindows() string {
	for _, path := range scriptsCandidatePaths() {
		if drive := resolveDriveFromConfig(path); drive != "" {
			return filepath.Join(drive+"\\", "gitmap-scripts")
		}
	}

	return `D:\gitmap-scripts`
}
