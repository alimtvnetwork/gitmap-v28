package cmdos

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos/enums"
)

// OSTUIViewModeType specifies the current mode of the OS TUI dashboard.
type OSTUIViewModeType = enums.OSTUIViewModeType

const (
	OSTUIViewModeNavType  = enums.OSTUIViewModeNavType
	OSTUIViewModeExecType = enums.OSTUIViewModeExecType
	OSTUIViewModeDoneType = enums.OSTUIViewModeDoneType
)

// OSTUIItem represents a toggleable configuration or maintenance action.
type OSTUIItem struct {
	ID          string
	Title       string
	Description string
	Badge       string
	IsChecked   bool
	IsSupported bool
}

// OSTUITab defines a categorized collection of OS dashboard items.
type OSTUITab struct {
	ID          string
	Title       string
	ShortcutKey string
	Items       []OSTUIItem
}

// OSTUIActionResult records execution output for an applied TUI action.
type OSTUIActionResult struct {
	Title     string
	IsSuccess bool
	Message   string
}

// IsFail reports whether the applied action failed.
func (r OSTUIActionResult) IsFail() bool {
	return !r.IsSuccess
}

// IsFailed reports whether the applied action failed.
func (r OSTUIActionResult) IsFailed() bool {
	return !r.IsSuccess
}
