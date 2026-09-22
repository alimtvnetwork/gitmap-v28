//go:build windows

package cmdos

import (
	"golang.org/x/sys/windows/registry"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	clsidBaseKey    = `Software\Classes\CLSID\{86ca1aa0-34aa-4e8b-a509-50c905bae2a2}`
	clsidInprocKey  = clsidBaseKey + `\InprocServer32`
	startMenuRegKey = `SYSTEM\ControlSet001\Control\FeatureManagement\Overrides\8\3036241548`
)

func applyContextMenuRegistry(isClassic bool) error {
	if isClassic {
		return enableClassicContextMenu()
	}
	return revertModernContextMenu()
}

func enableClassicContextMenu() error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, clsidInprocKey, registry.SET_VALUE)
	if err != nil {
		return apperror.WrapSimple(err, "failed to create CLSID InprocServer32 key")
	}
	defer k.Close()
	if err := k.SetStringValue("", ""); err != nil {
		return apperror.WrapSimple(err, "failed to set default empty value for classic context menu")
	}
	return nil
}

func revertModernContextMenu() error {
	_ = registry.DeleteKey(registry.CURRENT_USER, clsidInprocKey)
	_ = registry.DeleteKey(registry.CURRENT_USER, clsidBaseKey)
	return nil
}

func isClassicContextMenuActive() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, clsidInprocKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	val, _, err := k.GetStringValue("")
	return err == nil && len(val) == 0
}

func applyStartMenuRegistry(isClassic bool) error {
	if isClassic {
		return setStartMenuEnabledState(1)
	}
	return setStartMenuEnabledState(0)
}

func setStartMenuEnabledState(state uint32) error {
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, startMenuRegKey, registry.SET_VALUE)
	if err != nil {
		return apperror.WrapSimple(err, "failed to set Start Menu registry key (run as Administrator)")
	}
	defer k.Close()
	if err := k.SetDWordValue("EnabledState", state); err != nil {
		return apperror.WrapSimple(err, "failed to write EnabledState DWORD")
	}
	return nil
}

func isClassicStartMenuActive() bool {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, startMenuRegKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()
	val, _, err := k.GetIntegerValue("EnabledState")
	return err == nil && val == 1
}
