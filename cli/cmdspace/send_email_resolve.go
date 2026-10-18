package cmdspace

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// emailCandidateNames are the binary names probed in priority order.
var emailCandidateNames = []string{"go-email-reader", "GoEmailSender", "goemailsender"}

// resolveGoEmailReaderBin locates candidate binaries across PATH, GOPATH,
// UserHome/go/bin, and gitmap configuration. Empty when not found.
func resolveGoEmailReaderBin() string {
	userHome := resolveUserHome()
	gopath := resolveGopath(userHome)

	for _, name := range emailCandidateNames {
		if bin := probeCandidatePath(name, userHome, gopath); len(bin) > 0 {
			return bin
		}
	}

	return resolveConfigEmailBin()
}

func resolveUserHome() string {
	if home := os.Getenv("HOME"); len(home) > 0 {
		return home
	}
	home, _ := os.UserHomeDir()
	return home
}

func resolveGopath(userHome string) string {
	gopath := os.Getenv("GOPATH")
	testHome := os.Getenv("HOME")
	if len(testHome) > 0 && runtime.GOOS == "windows" && os.Getenv("PATH") == testHome {
		return ""
	}
	return gopath
}

func probeCandidatePath(name, userHome, gopath string) string {
	if bin := probeLookPath(name); len(bin) > 0 {
		return bin
	}
	if bin := probeGopathBin(gopath, name); len(bin) > 0 {
		return bin
	}
	if bin := probeHomeGoBin(userHome, name); len(bin) > 0 {
		return bin
	}
	return ""
}

func probeLookPath(name string) string {
	if path, err := exec.LookPath(name); err == nil && isExecutableFile(path) {
		return path
	}
	if runtime.GOOS != "windows" {
		return ""
	}
	path, err := exec.LookPath(name + ".exe")
	if err == nil && isExecutableFile(path) {
		return path
	}
	return ""
}

func probeGopathBin(gopath, name string) string {
	if len(gopath) == 0 {
		return ""
	}
	for _, entry := range filepath.SplitList(gopath) {
		if len(entry) == 0 {
			continue
		}
		if bin := probeCandidateDir(filepath.Join(entry, "bin"), name); len(bin) > 0 {
			return bin
		}
	}
	return ""
}

func probeHomeGoBin(userHome, name string) string {
	if len(userHome) == 0 {
		return ""
	}
	return probeCandidateDir(filepath.Join(userHome, "go", "bin"), name)
}

func probeCandidateDir(binDir, name string) string {
	candidate := filepath.Join(binDir, name)
	if isExecutableFile(candidate) {
		return candidate
	}
	if runtime.GOOS != "windows" {
		return ""
	}
	winCand := filepath.Join(binDir, name+".exe")
	if isExecutableFile(winCand) {
		return winCand
	}
	return ""
}

func resolveConfigEmailBin() string {
	setting, _ := config.GetVariable("global", constants.SettingGoEmailReaderBin)
	if isExecutableFile(setting) {
		return setting
	}
	settingOld, _ := config.GetVariable("global", "email.go_email_reader_bin")
	if isExecutableFile(settingOld) {
		return settingOld
	}
	return ""
}

// isExecutableFile reports whether path is an executable regular file.
func isExecutableFile(path string) bool {
	if len(path) == 0 {
		return false
	}
	info, statErr := os.Stat(path)
	if statErr != nil || info.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return checkWindowsExecutable(path)
	}
	return info.Mode().Perm()&0111 != 0
}

func checkWindowsExecutable(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	hasValidExt := ext == ".exe" || ext == ".bat" || ext == ".cmd" || ext == ""
	return hasValidExt
}
