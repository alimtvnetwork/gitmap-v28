//go:build windows

package cmdos

import (
	"golang.org/x/sys/windows/registry"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const winlogonRegPath = `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon`

type windowsAutoLoginEngine struct{}

func newPlatformAutoLoginEngine() AutoLoginEngine {
	return &windowsAutoLoginEngine{}
}

func (w *windowsAutoLoginEngine) Configure(cfg AutoLoginConfig) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, winlogonRegPath, registry.SET_VALUE)
	if err != nil {
		return apperror.WrapSimple(err, "failed to open Winlogon registry key (run as Administrator)")
	}
	defer k.Close()

	if err := k.SetStringValue("AutoAdminLogon", "1"); err != nil {
		return apperror.WrapSimple(err, "failed to set AutoAdminLogon")
	}
	if err := k.SetStringValue("DefaultUserName", cfg.Username); err != nil {
		return apperror.WrapSimple(err, "failed to set DefaultUserName")
	}
	if err := k.SetStringValue("DefaultDomainName", cfg.Domain); err != nil {
		return apperror.WrapSimple(err, "failed to set DefaultDomainName")
	}
	if err := k.SetStringValue("DefaultPassword", cfg.Password); err != nil {
		return apperror.WrapSimple(err, "failed to set DefaultPassword")
	}
	return nil
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

func (w *windowsAutoLoginEngine) Status() (AutoLoginStatus, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, winlogonRegPath, registry.QUERY_VALUE)
	if err != nil {
		return AutoLoginStatus{}, apperror.WrapSimple(err, "failed to read Winlogon registry key")
	}
	defer k.Close()

	autoVal, _, _ := k.GetStringValue("AutoAdminLogon")
	userVal, _, _ := k.GetStringValue("DefaultUserName")
	domainVal, _, _ := k.GetStringValue("DefaultDomainName")
	passVal, _, _ := k.GetStringValue("DefaultPassword")

	return AutoLoginStatus{
		IsEnabled:      autoVal == "1",
		Username:       userVal,
		Domain:         domainVal,
		HasPassword:    len(passVal) > 0,
		DisplayManager: "Winlogon",
	}, nil
}
