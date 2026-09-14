package cmdinstall

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func readIconFileContent(filePath string) []byte {
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() || info.Size() == 0 {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}

	return data
}

func findSourceIconBytesInDir(dir string) []byte {
	candidates := []string{
		filepath.Join(dir, "resources", "app", "resources", "icon.png"),
		filepath.Join(dir, "resources", "icon.png"),
		filepath.Join(dir, "icon.png"),
	}

	for _, p := range candidates {
		if data := readIconFileContent(p); len(data) > 0 {
			return data
		}
	}

	return nil
}

func resolveSourceIconBytes(opts IconDeployOptions) ([]byte, error) {
	if data := findSourceIconBytesInDir(opts.SourceDir); len(data) > 0 {
		return data, nil
	}

	if len(opts.RawBytes) > 0 {
		return opts.RawBytes, nil
	}

	return nil, apperror.NewSimple("no source icon data found", "E9000")
}

func writeHicolorThemeFile(indexPath string) error {
	if err := os.WriteFile(indexPath, []byte(HicolorIndexThemeContent), 0644); err != nil {
		return apperror.WrapSimple(err, "write.hicolorIndexTheme")
	}

	return nil
}

func ensureHicolorThemeIndex(hicolorDir string) error {
	indexPath := filepath.Join(hicolorDir, "index.theme")
	if _, err := os.Stat(indexPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(hicolorDir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir.hicolor")
	}

	return writeHicolorThemeFile(indexPath)
}

func writeResolutionIcon(hicolorDir, iconName string, icon ResizedIcon) error {
	resDir := filepath.Join(hicolorDir, fmt.Sprintf("%dx%d", icon.Size, icon.Size), "apps")
	if err := os.MkdirAll(resDir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir.iconRes")
	}

	destPath := filepath.Join(resDir, iconName+".png")
	if err := os.WriteFile(destPath, icon.Data, 0644); err != nil {
		return apperror.WrapSimple(err, "write.iconFile")
	}

	return nil
}

func deployIconsToHicolorDir(hicolorDir, iconName string, icons []ResizedIcon) error {
	if err := ensureHicolorThemeIndex(hicolorDir); err != nil {
		return err
	}

	for _, icon := range icons {
		if err := writeResolutionIcon(hicolorDir, iconName, icon); err != nil {
			return err
		}
	}

	return nil
}

func writePixmapsFile(destPath string, fallbackData []byte) error {
	if err := os.WriteFile(destPath, fallbackData, 0644); err != nil {
		return apperror.WrapSimple(err, "write.pixmapsIcon")
	}

	return nil
}

func deployPixmapsFallback(pixmapsDir, iconName string, fallbackData []byte) error {
	if len(fallbackData) == 0 {
		return nil
	}

	if err := os.MkdirAll(pixmapsDir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir.pixmaps")
	}

	return writePixmapsFile(filepath.Join(pixmapsDir, iconName+".png"), fallbackData)
}

func findDefaultIconData(icons []ResizedIcon) []byte {
	for _, ic := range icons {
		if ic.Size == 256 {
			return ic.Data
		}
	}

	if len(icons) > 0 {
		return icons[len(icons)-1].Data
	}

	return nil
}

func runGtkUpdateIconCache(hicolorDir string) {
	bin, err := exec.LookPath("gtk-update-icon-cache")
	if err != nil {
		return
	}

	_ = exec.Command(bin, "-f", "-t", hicolorDir).Run()
}

func isWritableDirectory(dir string) bool {
	probe := filepath.Join(dir, ".probe_perm")
	if err := os.WriteFile(probe, []byte(""), 0644); err != nil {
		return false
	}

	_ = os.Remove(probe)

	return true
}

func resolveHicolorDirectories(isUserOnly bool) []string {
	home, _ := os.UserHomeDir()
	dirs := []string{filepath.Join(home, ".local", "share", "icons", "hicolor")}

	if isUserOnly {
		return dirs
	}

	if isWritableDirectory("/usr/share/icons/hicolor") {
		dirs = append(dirs, "/usr/share/icons/hicolor")
	}

	return dirs
}

func resolvePixmapsDirectories(isUserOnly bool) []string {
	home, _ := os.UserHomeDir()
	dirs := []string{filepath.Join(home, ".local", "share", "pixmaps")}

	if isUserOnly {
		return dirs
	}

	if isWritableDirectory("/usr/share/pixmaps") {
		dirs = append(dirs, "/usr/share/pixmaps")
	}

	return dirs
}

func deployAllHicolorTrees(hicolorDirs []string, iconName string, icons []ResizedIcon) error {
	for _, hicolorDir := range hicolorDirs {
		if err := deployIconsToHicolorDir(hicolorDir, iconName, icons); err != nil {
			return err
		}

		runGtkUpdateIconCache(hicolorDir)
	}

	return nil
}

func deployAllPixmapsTrees(pixmapsDirs []string, iconName string, iconData []byte) error {
	for _, pixDir := range pixmapsDirs {
		if err := deployPixmapsFallback(pixDir, iconName, iconData); err != nil {
			return err
		}
	}

	return nil
}

func refreshDesktopDatabases() {
	home, _ := os.UserHomeDir()
	userApps := filepath.Join(home, ".local", "share", "applications")
	updateDesktopDatabase(userApps)
	updateDesktopDatabase("/usr/share/applications")
}

func prepareAndResizeIcons(opts IconDeployOptions) ([]ResizedIcon, error) {
	rawBytes, err := resolveSourceIconBytes(opts)
	if err != nil {
		return nil, err
	}

	return GenerateMultiResolutionIcons(rawBytes)
}

func DeployMultiResolutionIcons(opts IconDeployOptions) error {
	icons, err := prepareAndResizeIcons(opts)
	if err != nil {
		return err
	}

	if err := deployAllHicolorTrees(resolveHicolorDirectories(opts.IsUserOnly), opts.IconName, icons); err != nil {
		return err
	}

	_ = deployAllPixmapsTrees(resolvePixmapsDirectories(opts.IsUserOnly), opts.IconName, findDefaultIconData(icons))
	refreshDesktopDatabases()

	return nil
}
