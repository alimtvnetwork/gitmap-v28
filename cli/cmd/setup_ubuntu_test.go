package cmd

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type mockDirInfo struct{ name string }

func (m mockDirInfo) Name() string       { return m.name }
func (m mockDirInfo) Size() int64        { return 0 }
func (m mockDirInfo) Mode() fs.FileMode  { return fs.ModeDir | 0755 }
func (m mockDirInfo) ModTime() time.Time { return time.Now() }
func (m mockDirInfo) IsDir() bool        { return true }
func (m mockDirInfo) Sys() any           { return nil }

func TestFindZshBinary_Mock(t *testing.T) {
	origLook := lookPathFunc
	defer func() { lookPathFunc = origLook }()

	lookPathFunc = func(file string) (string, error) {
		if file == "zsh" {
			return "/usr/bin/zsh", nil
		}

		return "", exec.ErrNotFound
	}

	path, isFound := findZshBinary()
	if !isFound || path != "/usr/bin/zsh" {
		t.Fatalf("expected /usr/bin/zsh found, got path=%s, isFound=%v", path, isFound)
	}
}

func TestFindZshBinary_NotFound(t *testing.T) {
	origLook := lookPathFunc
	origStat := statPathFunc
	defer func() {
		lookPathFunc = origLook
		statPathFunc = origStat
	}()

	lookPathFunc = func(file string) (string, error) {
		return "", exec.ErrNotFound
	}

	statPathFunc = func(path string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	_, isFound := findZshBinary()
	if isFound {
		t.Fatal("expected zsh not found")
	}
}

func TestIsOhMyZshInstalled_Mock(t *testing.T) {
	origStat := statPathFunc
	origHome := userHomeDirFunc
	defer func() {
		statPathFunc = origStat
		userHomeDirFunc = origHome
	}()

	userHomeDirFunc = func() (string, error) {
		return "/home/mockuser", nil
	}

	statPathFunc = func(path string) (os.FileInfo, error) {
		if strings.HasSuffix(path, ".oh-my-zsh") {
			return mockDirInfo{name: ".oh-my-zsh"}, nil
		}

		return nil, os.ErrNotExist
	}

	if !isOhMyZshInstalled() {
		t.Fatal("expected isOhMyZshInstalled to return true")
	}
}

func TestEnsureZshUbuntuStep_AlreadyInstalled(t *testing.T) {
	origLook := lookPathFunc
	origStat := statPathFunc
	origHome := userHomeDirFunc
	defer func() {
		lookPathFunc = origLook
		statPathFunc = origStat
		userHomeDirFunc = origHome
	}()

	lookPathFunc = func(file string) (string, error) {
		return "/bin/zsh", nil
	}

	userHomeDirFunc = func() (string, error) {
		return "/home/mockuser", nil
	}

	statPathFunc = func(path string) (os.FileInfo, error) {
		return mockDirInfo{name: ".oh-my-zsh"}, nil
	}

	// Should not block or panic, prints installed status and returns.
	ensureZshUbuntuStep(false, false)
}

func TestEnsureZshUbuntuStep_SkipFlag(t *testing.T) {
	// When skip is set to true, should return immediately without checking.
	ensureZshUbuntuStep(false, true)
}

func TestEnsureZshUbuntuStep_SkipEnv(t *testing.T) {
	origEnv := getenvFunc
	defer func() { getenvFunc = origEnv }()

	getenvFunc = func(key string) string {
		if key == "GITMAP_SKIP_ZSH" {
			return "1"
		}

		return ""
	}

	ensureZshUbuntuStep(false, false)
}

func TestConfigureZshTheme_AlreadyConfigured(t *testing.T) {
	origHome := userHomeDirFunc
	origRead := osReadFileHook
	origWrite := osWriteFileHook
	defer func() {
		userHomeDirFunc = origHome
		osReadFileHook = origRead
		osWriteFileHook = origWrite
	}()

	userHomeDirFunc = func() (string, error) {
		return "/home/mockuser", nil
	}

	osReadFileHook = func(name string) ([]byte, error) {
		return []byte(`ZSH_THEME="agnoster"`), nil
	}

	isWritten := false
	osWriteFileHook = func(name string, data []byte, perm os.FileMode) error {
		isWritten = true

		return nil
	}

	configureZshTheme()
	if isWritten {
		t.Fatal("expected no write when theme is already agnoster")
	}
}

func TestConfigureZshTheme_ReplaceRobbyrussell(t *testing.T) {
	origHome := userHomeDirFunc
	origRead := osReadFileHook
	origWrite := osWriteFileHook
	defer func() {
		userHomeDirFunc = origHome
		osReadFileHook = origRead
		osWriteFileHook = origWrite
	}()

	userHomeDirFunc = func() (string, error) {
		return "/home/mockuser", nil
	}

	osReadFileHook = func(name string) ([]byte, error) {
		return []byte(`ZSH_THEME="robbyrussell"`), nil
	}

	var writtenData string
	osWriteFileHook = func(name string, data []byte, perm os.FileMode) error {
		writtenData = string(data)

		return nil
	}

	configureZshTheme()
	if !strings.Contains(writtenData, `ZSH_THEME="agnoster"`) {
		t.Fatalf("expected agnoster written, got %s", writtenData)
	}
}

func TestParseSetupFlags_SkipZsh(t *testing.T) {
	args := []string{"--skip-zsh"}
	_, dryRun, _, isSkip := parseSetupFlags(args)
	if dryRun || !isSkip {
		t.Fatalf("expected dryRun=false isSkip=true, got dryRun=%v isSkip=%v", dryRun, isSkip)
	}
}

func TestIsStdinTerminal_Error(t *testing.T) {
	origStat := stdinStatFunc
	defer func() { stdinStatFunc = origStat }()

	stdinStatFunc = func() (os.FileInfo, error) {
		return nil, errors.New("not a terminal")
	}

	if isStdinTerminal() {
		t.Fatal("expected false on stdin stat error")
	}
}
