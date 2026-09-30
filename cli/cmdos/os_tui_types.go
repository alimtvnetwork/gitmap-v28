package cmdos

// OSTUIViewModeType specifies the current mode of the OS TUI dashboard.
type OSTUIViewModeType string

const (
	OSTUIViewModeNavType  OSTUIViewModeType = "nav"
	OSTUIViewModeExecType OSTUIViewModeType = "exec"
	OSTUIViewModeDoneType OSTUIViewModeType = "done"
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
