package cmdos

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunSSHExecFn is the pluggable runner for remote SSH command execution.
var RunSSHExecFn func(args []string) error

// DockPosition represents cardinal screen edge orientation for desktop dock.
type DockPosition string

const (
	DockPositionBottom DockPosition = "bottom"
	DockPositionLeft   DockPosition = "left"
	DockPositionRight  DockPosition = "right"
	DockPositionTop    DockPosition = "top"
	DockPositionCenter DockPosition = "center"
)

// DockConfig holds status metadata about the desktop dock.
type DockConfig struct {
	Position     string `json:"position"`
	RawPosition  string `json:"rawPosition"`
	IsPanelMode  bool   `json:"isPanelMode"`
	HasAppsAtTop bool   `json:"hasAppsAtTop"`
	OS           string `json:"os"`
}

// DockOperator defines platform-specific querying and setting operations.
type DockOperator interface {
	GetDockConfig() (DockConfig, error)
	SetDockPosition(pos DockPosition) error
}

// DockCLIOptions holds parsed CLI parameters.
type DockCLIOptions struct {
	TargetNode  string
	Position    DockPosition
	HasPosition bool
	IsJSON      bool
	IsHelp      bool
}

// NormalizeDockPosition validates and normalizes raw dock position strings.
func NormalizeDockPosition(raw string) (DockPosition, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "bottom", "b", "down":
		return DockPositionBottom, nil
	case "left", "l":
		return DockPositionLeft, nil
	case "right", "r":
		return DockPositionRight, nil
	case "top", "t", "up":
		return DockPositionTop, nil
	case "center", "c":
		return DockPositionCenter, nil
	default:
		return "", apperror.NewSimple("invalid dock position: "+raw, "E_INVALID_POS")
	}
}

// CheckDockFlag parses help, json, and node flags from CLI args.
func CheckDockFlag(args []string, i int, opts *DockCLIOptions, cleaned *[]string) (int, bool) {
	a := args[i]
	if a == "-h" || a == "--help" || a == "help" {
		opts.IsHelp = true
		return 0, true
	}
	if a == "--json" || a == "-j" {
		opts.IsJSON = true
		*cleaned = append(*cleaned, a)
		return 0, true
	}
	return checkDockNodeFlag(args, i, opts)
}

func checkDockNodeFlag(args []string, i int, opts *DockCLIOptions) (int, bool) {
	a := args[i]
	if strings.HasPrefix(a, "--node=") || strings.HasPrefix(a, "-node=") {
		opts.TargetNode = strings.SplitN(a, "=", 2)[1]
		return 0, true
	}
	if (a == "--node" || a == "-node") && i+1 < len(args) {
		opts.TargetNode = args[i+1]
		return 1, true
	}
	return 0, false
}
