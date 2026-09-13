//go:build darwin

package cmdinstall

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/tempdir"
)

func checkDarwinDirExists(dir string) bool {
	info, err := os.Stat(dir)

	return err == nil && info.IsDir()
}

func isWritableDarwinDir(dir string) bool {
	if !checkDarwinDirExists(dir) {
		return false
	}

	testFile := filepath.Join(dir, ".gitmap-write-test")
	if err := os.WriteFile(testFile, []byte("1"), 0644); err != nil {
		return false
	}

	_ = os.Remove(testFile)

	return true
}

func resolveDarwinBinDir() string {
	if hasWritePerm := isWritableDarwinDir("/usr/local/bin"); hasWritePerm {
		return "/usr/local/bin"
	}

	home, _ := os.UserHomeDir()

	return filepath.Join(home, ".local", "bin")
}

func resolveDarwinInstallDirs(prefix string) (string, string) {
	if prefix != "" {
		return prefix, resolveDarwinBinDir()
	}

	home, _ := os.UserHomeDir()
	if hasWritePerm := isWritableDarwinDir("/Applications"); hasWritePerm {
		return "/Applications", resolveDarwinBinDir()
	}

	return filepath.Join(home, "Applications"), resolveDarwinBinDir()
}

func attachDarwinDiskImage(dmgPath, mountPoint string) error {
	cmd := exec.Command("hdiutil", "attach", dmgPath, "-nobrowse", "-quiet", "-mountpoint", mountPoint)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "hdiutil.attach")
	}

	return nil
}

func detachDarwinDiskImage(mountPoint string) {
	if mountPoint == "" {
		return
	}

	_ = exec.Command("hdiutil", "detach", mountPoint, "-quiet").Run()
}

func findDarwinAppBundle(mountPoint string) (string, error) {
	entries, err := os.ReadDir(mountPoint)
	if err != nil {
		return "", apperror.WrapSimple(err, "readDir.mountPoint")
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".app") {
			return filepath.Join(mountPoint, entry.Name()), nil
		}
	}

	return "", apperror.NewSimple("no .app bundle found inside DMG", "E9000")
}

func copyDarwinAppBundle(srcApp, destApp string) error {
	_ = os.RemoveAll(destApp)
	cmd := exec.Command("ditto", srcApp, destApp)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "ditto.copyApp")
	}

	return nil
}

func removeDarwinQuarantine(appPath string) {
	_ = exec.Command("xattr", "-dr", "com.apple.quarantine", appPath).Run()
}

func isValidDarwinBinaryCandidate(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		return false
	}

	return !info.IsDir() && info.Size() > 0
}

func getDarwinBinaryCandidates(appPath string) []string {
	return []string{
		filepath.Join(appPath, "Contents", "Resources", "app", "bin", "antigravity"),
		filepath.Join(appPath, "Contents", "Resources", "app", "bin", "agy"),
		filepath.Join(appPath, "Contents", "MacOS", "Antigravity"),
		filepath.Join(appPath, "Contents", "MacOS", "antigravity"),
	}
}

func resolveDarwinCliBinary(appPath string) (string, error) {
	candidates := getDarwinBinaryCandidates(appPath)
	for _, cand := range candidates {
		if isValidDarwinBinaryCandidate(cand) {
			return cand, nil
		}
	}

	return "", apperror.NewSimple("Antigravity CLI executable not found inside app bundle", "E9000")
}

func symlinkDarwinBinary(targetBinary, symlinkPath string) error {
	_ = os.Remove(symlinkPath)
	if err := os.Symlink(targetBinary, symlinkPath); err != nil {
		return apperror.WrapSimple(err, "symlink.darwin")
	}

	return nil
}

func createDualDarwinSymlinks(cliBinary, binDir string) error {
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir.darwinBinDir")
	}

	if err := symlinkDarwinBinary(cliBinary, filepath.Join(binDir, "antigravity")); err != nil {
		return err
	}

	return symlinkDarwinBinary(cliBinary, filepath.Join(binDir, "agy"))
}

func finalizeDarwinCliSymlinks(destApp, binDir string) error {
	cliBinary, errBin := resolveDarwinCliBinary(destApp)
	if errBin != nil {
		return errBin
	}

	_ = os.Chmod(cliBinary, 0755)
	ensureDirInPath(binDir)

	return createDualDarwinSymlinks(cliBinary, binDir)
}

func copyAndUnquarantineBundle(srcApp, destApp string) error {
	if err := copyDarwinAppBundle(srcApp, destApp); err != nil {
		return err
	}

	removeDarwinQuarantine(destApp)

	return nil
}

func completeDarwinDeployment(mountPoint, destDir, binDir string) error {
	srcApp, errSrc := findDarwinAppBundle(mountPoint)
	if errSrc != nil {
		return errSrc
	}

	destApp := filepath.Join(destDir, filepath.Base(srcApp))
	if err := copyAndUnquarantineBundle(srcApp, destApp); err != nil {
		return err
	}

	return finalizeDarwinCliSymlinks(destApp, binDir)
}

func mountAndDeployDarwin(dmgPath, destDir, binDir string) error {
	mountPoint, errMount := os.MkdirTemp("", "antigravity-mount-*")
	if errMount != nil {
		return apperror.WrapSimple(errMount, "mkdirtemp.darwinMount")
	}
	defer os.RemoveAll(mountPoint)
	defer detachDarwinDiskImage(mountPoint)

	if err := attachDarwinDiskImage(dmgPath, mountPoint); err != nil {
		return err
	}

	return completeDarwinDeployment(mountPoint, destDir, binDir)
}

func deployAntigravityDesktopDarwin(dmgPath, prefix string) error {
	destDir, binDir := resolveDarwinInstallDirs(prefix)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir.darwinDestDir")
	}

	return mountAndDeployDarwin(dmgPath, destDir, binDir)
}

func installAntigravityDesktopPlatform(opts installOptions) error {
	url := getAntigravityDesktopDownloadUrl("darwin")
	tempDmg := filepath.Join(tempdir.RepoTempDir("downloads"), "Antigravity.dmg")
	defer os.Remove(tempDmg)

	if err := downloadFileToDest(url, tempDmg); err != nil {
		return apperror.WrapSimple(err, "download.darwinDmg")
	}

	return deployAntigravityDesktopDarwin(tempDmg, opts.Prefix)
}

func getDarwinAppPaths() []string {
	home, _ := os.UserHomeDir()

	return []string{
		"/Applications/Antigravity.app",
		filepath.Join(home, "Applications", "Antigravity.app"),
		"/usr/local/bin/antigravity",
		filepath.Join(home, ".local", "bin", "antigravity"),
	}
}

func findInstalledAntigravityDesktopPath() (string, bool) {
	for _, p := range getDarwinAppPaths() {
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}

	if p, err := exec.LookPath("antigravity"); err == nil {
		return p, true
	}

	return "", false
}
