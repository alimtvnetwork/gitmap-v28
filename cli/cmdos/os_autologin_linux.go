//go:build linux

package cmdos

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type linuxAutoLoginEngine struct{}

func newPlatformAutoLoginEngine() AutoLoginOperator {
	return &linuxAutoLoginEngine{}
}

func (l *linuxAutoLoginEngine) Configure(cfg AutoLoginConfig) error {
	if isFilePresent(gdm3ConfigPath) {
		return configureGDM3(cfg.Username)
	}
	if isDirPresent(lightdmConfigDir) {
		return configureLightDM(cfg.Username)
	}
	return apperror.NewSimple("unsupported Linux display manager (neither GDM3 nor LightDM detected)", "E_LINUX_DM_UNSUPPORTED")
}

func (l *linuxAutoLoginEngine) Disable() error {
	if isFilePresent(lightdmFile) {
		_ = os.Remove(lightdmFile)
	}
	if isFilePresent(gdm3ConfigPath) {
		return disableGDM3()
	}
	return nil
}

func (l *linuxAutoLoginEngine) Status() (AutoLoginStatus, error) {
	if isFilePresent(gdm3ConfigPath) {
		return readGDM3Status()
	}
	if isFilePresent(lightdmFile) {
		return readLightDMStatus()
	}
	return AutoLoginStatus{IsEnabled: false, DisplayManager: "None"}, nil
}

func isFilePresent(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isDirPresent(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
