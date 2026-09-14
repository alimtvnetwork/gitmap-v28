package macro

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func isSudoUserHomeValid(sudoUser string) bool {
	if sudoUser == "" || sudoUser == "root" {
		return false
	}
	userHome := filepath.Join("/home", sudoUser)
	info, err := os.Stat(userHome)

	return err == nil && info.IsDir()
}

func resolveEffectiveHomeDir() string {
	sudoUser := os.Getenv("SUDO_USER")
	if isSudoUserHomeValid(sudoUser) {
		return filepath.Join("/home", sudoUser)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return home
}

func resolveXdgConfigMacroDir(home string) string {
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig != "" {
		return filepath.Join(xdgConfig, "gitmap", "macros")
	}

	return filepath.Join(home, ".config", "gitmap", "macros")
}

func resolveXdgDataMacroDir(home string) string {
	xdgData := os.Getenv("XDG_DATA_HOME")
	if xdgData != "" {
		return filepath.Join(xdgData, "gitmap", "macros")
	}

	return filepath.Join(home, ".local", "share", "gitmap", "macros")
}

func resolveCurrentUserName() string {
	userName := os.Getenv("USER")
	if userName != "" {
		return userName
	}
	userWin := os.Getenv("USERNAME")
	if userWin != "" {
		return userWin
	}

	return "user"
}

func resolveTempMacroDir() string {
	userName := resolveCurrentUserName()
	dirName := fmt.Sprintf(".gitmap-%s", userName)

	return filepath.Join(os.TempDir(), dirName, "macros")
}

func candidateMacroDirs() []string {
	var candidates []string
	home := resolveEffectiveHomeDir()
	if home != "" {
		candidates = append(candidates, filepath.Join(home, constants.GitMapDir, "macros"))
		candidates = append(candidates, resolveXdgConfigMacroDir(home))
		candidates = append(candidates, resolveXdgDataMacroDir(home))
	}
	candidates = append(candidates, filepath.Join(".", constants.GitMapDir, "macros"))
	candidates = append(candidates, resolveTempMacroDir())

	return deduplicateDirs(candidates)
}

func deduplicateDirs(dirs []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, dir := range dirs {
		clean := filepath.Clean(dir)
		_, isSeen := seen[clean]
		if isSeen {
			continue
		}
		seen[clean] = true
		out = append(out, clean)
	}

	return out
}

func isDirWritable(dir string) bool {
	probe := filepath.Join(dir, fmt.Sprintf(".perm_probe_%d", time.Now().UnixNano()))
	file, err := os.OpenFile(probe, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return false
	}
	_ = file.Close()
	_ = os.Remove(probe)

	return true
}

func probeAndPrepareDir(dir string) bool {
	_ = os.MkdirAll(dir, 0755)
	if !isDirWritable(dir) {
		attemptPermissionHealing(dir)
		_ = os.MkdirAll(dir, 0755)
	}
	if !isDirWritable(dir) {
		return false
	}
	restoreSudoOwnership(dir)

	return true
}

func resolveWritableMacroDir() (string, error) {
	for _, dir := range candidateMacroDirs() {
		if probeAndPrepareDir(dir) {
			return dir, nil
		}
	}

	return "", apperror.NewSimple("no writable macro directory found", "E_NO_WRITABLE_MACRO_DIR")
}
