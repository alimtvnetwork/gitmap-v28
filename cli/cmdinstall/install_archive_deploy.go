package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func resolveLinuxAppDirs(appName string) (string, string, string) {
	if os.Geteuid() == 0 {
		return filepath.Join("/opt", appName), "/usr/local/bin", "/usr/share/applications"
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", appName),
		filepath.Join(home, ".local", "bin"),
		filepath.Join(home, ".local", "share", "applications")
}

func deployArchiveApp(params ArchiveDeployParams) error {
	switch params.Inspection.Strategy {
	case StrategyBinaryApp, StrategySingleGz:
		return deployBinaryApplication(params)
	case StrategyScript:
		return deployViaInstallScript(params)
	case StrategySource:
		return deployViaSourceBuild(params)
	default:
		return apperror.NewSimple("unrecognized archive package structure; no binary or install script found", "E9000")
	}
}

func deployBinaryApplication(params ArchiveDeployParams) error {
	destAppDir, binDir, desktopDir := resolveLinuxAppDirs(params.Opts.AppName)
	_ = os.RemoveAll(destAppDir)
	if err := copyDirectoryTree(params.Inspection.SourceRoot, destAppDir); err != nil {
		return apperror.WrapSimple(err, "archive.deployAppDir")
	}

	relBin, _ := filepath.Rel(params.Inspection.SourceRoot, params.Inspection.BinaryPath)
	targetBin := filepath.Join(destAppDir, relBin)
	_ = os.Chmod(targetBin, 0755)

	if err := linkAppBinary(targetBin, binDir, params.Opts.AppName); err != nil {
		return err
	}
	deployAppDesktopEntry(params, destAppDir, targetBin, desktopDir)
	recordToolInDatabases(params.Opts.AppName, "archive", "archive")
	return nil
}

func linkAppBinary(targetBin, binDir, appName string) error {
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return apperror.WrapSimple(err, "archive.mkdirBin")
	}
	symlinkPath := filepath.Join(binDir, appName)
	_ = os.Remove(symlinkPath)
	if err := os.Symlink(targetBin, symlinkPath); err != nil {
		return apperror.WrapSimple(err, "archive.symlink")
	}
	ensureDirInPath(binDir)
	return nil
}

func deployAppDesktopEntry(params ArchiveDeployParams, destAppDir, targetBin, desktopDir string) {
	_ = os.MkdirAll(desktopDir, 0755)
	destDesktop := filepath.Join(desktopDir, params.Opts.AppName+".desktop")
	if params.Inspection.DesktopPath != "" {
		_ = copyFile(params.Inspection.DesktopPath, destDesktop)
		updateDesktopDatabase(desktopDir)
		return
	}
	icon := params.Inspection.IconPath
	if icon == "" {
		icon = resolveAppIconPath(destAppDir)
	}
	_ = createArchiveDesktopEntry(params.Opts.AppName, targetBin, icon, destDesktop)
	updateDesktopDatabase(desktopDir)
}

func createArchiveDesktopEntry(appName, binPath, iconPath, desktopFile string) error {
	content := fmt.Sprintf("[Desktop Entry]\nName=%s\nExec=%s %%U\nIcon=%s\nType=Application\nStartupNotify=true\nCategories=Utility;Development;\n", appName, binPath, iconPath)
	return os.WriteFile(desktopFile, []byte(content), 0644)
}

func deployViaInstallScript(params ArchiveDeployParams) error {
	script := params.Inspection.ScriptPath
	_ = os.Chmod(script, 0755)
	cmd := exec.Command("bash", script)
	cmd.Dir = params.Inspection.SourceRoot
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "archive.runInstallScript")
	}
	recordToolInDatabases(params.Opts.AppName, "archive", "script")
	return nil
}

func deployViaSourceBuild(params ArchiveDeployParams) error {
	if params.Inspection.HasConfigure {
		conf := exec.Command("./configure")
		conf.Dir = params.Inspection.SourceRoot
		_ = conf.Run()
	}
	makeCmd := exec.Command("make", "install")
	makeCmd.Dir = params.Inspection.SourceRoot
	makeCmd.Stdout, makeCmd.Stderr = os.Stdout, os.Stderr
	if err := makeCmd.Run(); err != nil {
		return apperror.WrapSimple(err, "archive.makeInstall")
	}
	recordToolInDatabases(params.Opts.AppName, "archive", "source")
	return nil
}

func copyDirectoryTree(src, dst string) error {
	cmd := exec.Command("cp", "-r", src, dst)
	if err := cmd.Run(); err == nil {
		return nil
	}
	return os.Rename(src, dst)
}
