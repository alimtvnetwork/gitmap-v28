package cmdos

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos/enums"
)

// RunSSHExecFn is the pluggable runner for remote SSH command execution.
var RunSSHExecFn func(args []string) error

// DockPositionType represents cardinal screen edge orientation for desktop dock.
type DockPositionType = enums.DockPositionType

// DockPosition is an alias for DockPositionType.
type DockPosition = DockPositionType

const (
	DockPositionBottom = enums.DockPositionBottom
	DockPositionLeft   = enums.DockPositionLeft
	DockPositionRight  = enums.DockPositionRight
	DockPositionTop    = enums.DockPositionTop
	DockPositionCenter = enums.DockPositionCenter
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
	SetDockPosition(pos DockPositionType) error
}

// DockCLIOptions holds parsed CLI parameters.
type DockCLIOptions struct {
	TargetNode  string
	Position    DockPositionType
	HasPosition bool
	IsJSON      bool
	IsHelp      bool
}

// NormalizeDockPosition validates and normalizes raw dock position strings.
// It delegates to the canonical enums parser, which accepts canonical
// positions and their shorthand aliases, and returns its enriched app error.
func NormalizeDockPosition(raw string) (DockPositionType, error) {
	pos, err := enums.ParseDockPositionType(raw)
	if err != nil {
		return "", err
	}

	return pos, nil
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
