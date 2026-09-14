//go:build linux

package cmdinstall

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func assertStaleNamePresent(t *testing.T, target string) {
	for _, name := range AntigravityStaleLauncherNames {
		if name == target {
			return
		}
	}
	t.Fatalf("expected stale launcher name %s not found", target)
}

func TestAntigravityStaleLauncherNames(t *testing.T) {
	if len(AntigravityStaleLauncherNames) != 5 {
		t.Fatalf("expected 5 stale launcher names, got %d", len(AntigravityStaleLauncherNames))
	}
	assertStaleNamePresent(t, "antigravity-ide.desktop")
	assertStaleNamePresent(t, "Google Antigravity.desktop")
	assertStaleNamePresent(t, "Google-Antigravity.desktop")
	assertStaleNamePresent(t, "Antigravity.desktop")
	assertStaleNamePresent(t, "antigravity.desktop")
}

func assertContentContains(t *testing.T, content, substring string) {
	if !strings.Contains(content, substring) {
		t.Fatalf("expected content to contain %q, but got: %s", substring, content)
	}
}

func TestBuildAntigravityDesktopContent(t *testing.T) {
	content := buildAntigravityDesktopContent("/usr/bin/antigravity", "my-icon")
	assertContentContains(t, content, "[Desktop Entry]")
	assertContentContains(t, content, "StartupWMClass=Antigravity")
	assertContentContains(t, content, "Exec=/usr/bin/antigravity %U")
	assertContentContains(t, content, "Icon=my-icon")
	assertContentContains(t, content, "Type=Application")
}

func TestBuildAntigravityDesktopContentDefaultIcon(t *testing.T) {
	content := buildAntigravityDesktopContent("/bin/antigravity", "")
	assertContentContains(t, content, "Icon=antigravity")
}

func createStaleTestFiles(t *testing.T, dir string) {
	for _, name := range AntigravityStaleLauncherNames {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("stale"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}
	_ = os.WriteFile(filepath.Join(dir, "keep.desktop"), []byte("keep"), 0644)
}

func assertStaleFilesRemoved(t *testing.T, dir string) {
	for _, name := range AntigravityStaleLauncherNames {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			t.Fatalf("expected stale file %s to be removed", name)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "keep.desktop")); err != nil {
		t.Fatalf("expected non-stale file to be preserved: %v", err)
	}
}

func TestPurgeStaleLaunchersFromDir(t *testing.T) {
	tmpDir := t.TempDir()
	createStaleTestFiles(t, tmpDir)
	purgeStaleLaunchersFromDir(tmpDir)
	assertStaleFilesRemoved(t, tmpDir)
}

func createDeployTestConfig(userAppDir, desktopDir string) AntigravityLauncherDeployConfig {
	return AntigravityLauncherDeployConfig{
		BinPath:      "/opt/antigravity/antigravity",
		UserAppDir:   userAppDir,
		DesktopDir:   desktopDir,
		IconName:     "antigravity",
		IsPrivileged: false,
	}
}

func assertDesktopFileValid(t *testing.T, filePath string) {
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("expected desktop file %s to exist: %v", filePath, err)
	}
	if info.Mode().Perm()&0111 == 0 {
		t.Fatalf("expected desktop file %s to have executable permissions, got %v", filePath, info.Mode())
	}
}

func TestDeployCanonicalAntigravityLauncher(t *testing.T) {
	userAppDir := t.TempDir()
	desktopDir := t.TempDir()
	cfg := createDeployTestConfig(userAppDir, desktopDir)

	if err := DeployCanonicalAntigravityLauncher(cfg); err != nil {
		t.Fatalf("DeployCanonicalAntigravityLauncher failed: %v", err)
	}
	assertDesktopFileValid(t, filepath.Join(userAppDir, "antigravity.desktop"))
	assertDesktopFileValid(t, filepath.Join(desktopDir, "antigravity.desktop"))
}

func TestMarkDesktopFileTrustedNoPanic(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "antigravity.desktop")
	_ = os.WriteFile(testFile, []byte("[Desktop Entry]\n"), 0755)
	markDesktopFileTrusted(testFile)
}

func TestBuildLauncherDeployConfig(t *testing.T) {
	cfg := buildLauncherDeployConfig("/bin/antigravity", "/custom/apps")
	if cfg.UserAppDir != "/custom/apps" {
		t.Fatalf("expected /custom/apps, got %s", cfg.UserAppDir)
	}
	if cfg.BinPath != "/bin/antigravity" {
		t.Fatalf("expected /bin/antigravity, got %s", cfg.BinPath)
	}
}

func TestCheckDirExists(t *testing.T) {
	tmpDir := t.TempDir()
	if !checkDirExists(tmpDir) {
		t.Fatalf("expected %s to exist", tmpDir)
	}
	if checkDirExists(filepath.Join(tmpDir, "non-existent")) {
		t.Fatalf("expected non-existent dir to not exist")
	}
}

func TestIsDirWritable(t *testing.T) {
	tmpDir := t.TempDir()
	if !isDirWritable(tmpDir) {
		t.Fatalf("expected temp dir %s to be writable", tmpDir)
	}
}
