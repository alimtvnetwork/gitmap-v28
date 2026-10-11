//go:build windows

package cmdos

import (
	"golang.org/x/sys/windows/registry"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const winlogonRegPath = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`

type windowsAutoLoginEngine struct{}

func newPlatformAutoLoginEngine() AutoLoginOperator {
	return &windowsAutoLoginEngine{}
}

func writeWinlogonValues(k registry.Key, cfg AutoLoginConfig) error {
	if err := k.SetStringValue("AutoAdminLogon", "1"); err != nil {
		return apperror.WrapSimple(err, "failed to set AutoAdminLogon")
	}
	if err := k.SetStringValue("DefaultUserName", cfg.Username); err != nil {
		return apperror.WrapSimple(err, "failed to set DefaultUserName")
	}
	return k.SetStringValue("DefaultDomainName", cfg.Domain)
}

func (w *windowsAutoLoginEngine) Configure(cfg AutoLoginConfig) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, winlogonRegPath, registry.SET_VALUE)
	if err != nil {
		return apperror.WrapSimple(err, "failed to open Winlogon registry key (run as Administrator)")
	}
	defer k.Close()
	return writeWinlogonValues(k, cfg)
}

func (w *windowsAutoLoginEngine) Disable() error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, winlogonRegPath, registry.SET_VALUE)
	if err != nil {
		return apperror.WrapSimple(err, "failed to open Winlogon registry key (run as Administrator)")
	}
	defer k.Close()

	if err := k.SetStringValue("AutoAdminLogon", "0"); err != nil {
		return apperror.WrapSimple(err, "failed to clear AutoAdminLogon")
	}
	_ = k.DeleteValue("DefaultPassword")
	return nil
}

func readWinlogonIdentity(k registry.Key) (string, string, bool) {
	autoVal, _, _ := k.GetStringValue("AutoAdminLogon")
	userVal, _, _ := k.GetStringValue("DefaultUserName")
	domainVal, _, _ := k.GetStringValue("DefaultDomainName")
	hasAutoAdminLogon := autoVal == "1"
	return userVal, domainVal, hasAutoAdminLogon
}

func (w *windowsAutoLoginEngine) Status() (AutoLoginStatus, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, winlogonRegPath, registry.QUERY_VALUE)
	if err != nil {
		return AutoLoginStatus{}, apperror.WrapSimple(err, "failed to read Winlogon registry key")
	}
	defer k.Close()

	userVal, domainVal, hasAutoAdminLogon := readWinlogonIdentity(k)
	return AutoLoginStatus{
		IsEnabled:      hasAutoAdminLogon,
		Username:       userVal,
		Domain:         domainVal,
		HasPassword:    false,
		DisplayManager: "Winlogon",
	}, nil
}
