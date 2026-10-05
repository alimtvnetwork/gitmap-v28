package cmdos

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var defaultDockOperator DockOperator

func getDockOperator() DockOperator {
	if defaultDockOperator != nil {
		return defaultDockOperator
	}
	return newPlatformDockEngine()
}

// RunOSDockCommand routes os dock, panel, start-menu, and taskbar commands.
func RunOSDockCommand(args []string) error {
	opts, cleaned, err := parseDockCLIOptions(args)
	if err != nil {
		return err
	}
	if opts.IsHelp {
		fmt.Println("Usage: gitmap os dock [bottom|left|right|top] [--node <alias>] [--json]")
		return nil
	}
	if opts.TargetNode != "" {
		return delegateRemoteDock(opts.TargetNode, cleaned)
	}
	if opts.HasPosition {
		return executeSetDock(opts.Position)
	}
	return executeQueryDock(opts.IsJSON)
}

func parseDockCLIOptions(args []string) (DockCLIOptions, []string, error) {
	var opts DockCLIOptions
	var cleaned []string
	for i := 0; i < len(args); i++ {
		adv, err := parseDockArg(args, i, &opts, &cleaned)
		if err != nil {
			return opts, nil, err
		}
		i += adv
	}
	return opts, cleaned, nil
}

func parseDockArg(args []string, i int, opts *DockCLIOptions, cleaned *[]string) (int, error) {
	adv, isFlag := CheckDockFlag(args, i, opts, cleaned)
	if isFlag {
		return adv, nil
	}
	pos, err := NormalizeDockPosition(args[i])
	if err != nil {
		return 0, err
	}
	opts.Position, opts.HasPosition = pos, true
	*cleaned = append(*cleaned, string(pos))
	return 0, nil
}

func delegateRemoteDock(node string, cleaned []string) error {
	remoteCmd := strings.TrimSpace("gitmap os dock " + strings.Join(cleaned, " "))
	fmt.Printf("● Delegating dock configuration to remote node '%s'...\n", node)
	if RunSSHExecFn == nil {
		return apperror.NewSimple("ssh runner not configured", "E_SSH_RUNNER")
	}
	return RunSSHExecFn([]string{node, remoteCmd})
}

func executeSetDock(pos DockPosition) error {
	if err := getDockOperator().SetDockPosition(pos); err != nil {
		return err
	}
	fmt.Printf("✔ Desktop dock position updated to: %s\n", pos)
	return nil
}

func executeQueryDock(isJSON bool) error {
	cfg, err := getDockOperator().GetDockConfig()
	if err != nil {
		return err
	}
	if isJSON {
		data, _ := json.MarshalIndent(cfg, "", "  ")
		fmt.Println(string(data))
		return nil
	}
	fmt.Printf("▶ Desktop Dock: OS=%s Pos=%s Raw=%s Panel=%t\n", cfg.OS, cfg.Position, cfg.RawPosition, cfg.IsPanelMode)
	return nil
}
