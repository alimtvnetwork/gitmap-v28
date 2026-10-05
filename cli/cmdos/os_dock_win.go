//go:build windows

package cmdos

import (
	"fmt"

	"golang.org/x/sys/windows/registry"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const explorerAdvancedKey = `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`

type windowsDockEngine struct{}

func newPlatformDockEngine() DockOperator {
	return &windowsDockEngine{}
}

func (w *windowsDockEngine) GetDockConfig() (DockConfig, error) {
	cfg := DockConfig{OS: "windows", IsPanelMode: true, Position: "bottom", RawPosition: "bottom"}
	val, err := readWindowsTaskbarAl()
	if err == nil {
		if val == 0 {
			cfg.Position = "left"
			cfg.RawPosition = "0"
		} else {
			cfg.Position = "center"
			cfg.RawPosition = "1"
		}
	}
	return cfg, nil
}

func (w *windowsDockEngine) SetDockPosition(pos DockPosition) error {
	val := uint32(1)
	if pos == DockPositionLeft {
		val = 0
	}
	if pos == DockPositionTop || pos == DockPositionRight {
		fmt.Printf("ℹ Note: Windows shell fixes taskbar to bottom edge; alignment updated.\n")
	}
	return writeWindowsTaskbarAl(val)
}

func readWindowsTaskbarAl() (uint64, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, explorerAdvancedKey, registry.QUERY_VALUE)
	if err != nil {
		return 1, err
	}
	defer k.Close()

	val, _, err := k.GetIntegerValue("TaskbarAl")
	if err != nil {
		return 1, err
	}

	return val, nil
}

func writeWindowsTaskbarAl(val uint32) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, explorerAdvancedKey, registry.SET_VALUE)
	if err != nil {
		return apperror.WrapSimple(err, "failed to open Explorer Advanced registry key")
	}
	defer k.Close()

	return k.SetDWordValue("TaskbarAl", val)
}
