package cmdagent

import (
	"fmt"
	"strings"
)

// isHelpRequested checks if any token in args is a help flag or keyword.
func isHelpRequested(args []string) bool {
	for _, a := range args {
		trimmed := strings.ToLower(strings.TrimSpace(a))
		if trimmed == "-h" || trimmed == "--help" || trimmed == "help" || trimmed == "-help" {
			return true
		}
	}

	return false
}

func printAgentResetHelp() {
	fmt.Println(`
  ╔══════════════════════════════════════════════════╗
  ║                gitmap agent reset                ║
  ╚══════════════════════════════════════════════════╝

  Reset agent split databases and telemetry while preserving task run directories.

  Usage:
    gitmap agent reset [flags]

  Flags:
    -d, --dir <path>       Target run directory (default: .ai-memory/temp-agents/)
    -t, --task-id <id>     Target parent task ID or slug
    -y, --yes, --force     Confirm reset without prompt
    -j, --json             Output result in JSON format
    -h, --help             Show this help menu

  Examples:
    gitmap agent reset
    gitmap agent reset -d .ai-memory/temp-agents/
    gitmap agent reset -t task-12345
    gitmap agent reset -y`)
}

func printAgentClearHelp() {
	fmt.Println(`
  ╔══════════════════════════════════════════════════╗
  ║                gitmap agent clear                ║
  ╚══════════════════════════════════════════════════╝

  Clear or archive completed agent task runs.

  Usage:
    gitmap agent clear [flags]

  Flags:
    -d, --dir <path>       Target run directory (default: .ai-memory/temp-agents/)
    -t, --task-id <id>     Target parent task ID to clear
    -y, --yes, --force     Confirm clear without prompt
    -j, --json             Output result in JSON format
    -h, --help             Show this help menu

  Examples:
    gitmap agent clear
    gitmap agent clear -t task-12345
    gitmap agent clear -y`)
}

func printAgentTempClearHelp() {
	fmt.Println(`
  ╔══════════════════════════════════════════════════╗
  ║              gitmap agent temp-clear             ║
  ╚══════════════════════════════════════════════════╝

  Delete entire ephemeral .ai-memory/temp-agents/ directory.

  Aliases:
    clean-temp, purge-temp

  Usage:
    gitmap agent temp-clear [flags]

  Flags:
    -d, --dir <path>       Target run directory (default: .ai-memory/temp-agents/)
    -y, --yes, --force     Confirm deletion without prompt
    -j, --json             Output result in JSON format
    -h, --help             Show this help menu

  Examples:
    gitmap agent temp-clear
    gitmap agent temp-clear -y`)
}

func printAgentLogHelp() {
	fmt.Println(`
  ╔══════════════════════════════════════════════════╗
  ║                 gitmap agent log                 ║
  ╚══════════════════════════════════════════════════╝

  Record granular in-flight telemetry into agent split database.

  Aliases:
    log-action, action

  Usage:
    gitmap agent log [flags]

  Flags:
    -t, --task-id <id>     Parent task ID or slug (required)
    -s, --subtask <id>     Subtask ID (required)
    -a, --agent <role>     Agent role completing action (default: Worker)
    --action <type>        Action type (exec, read, write, edit, delete, inspect)
    --file <path>          Target file path
    --details <msg>        Action details or payload
    --status <status>      Status (IN_PROGRESS, DONE, FAILED)
    -d, --dir <path>       Target directory
    -h, --help             Show this help menu

  Examples:
    gitmap agent log -t task-123 -s sub-01 --action exec --file main.go
    gitmap agent log -t task-123 -s sub-01 --action write --file config.json`)
}

func printAgentCrashedHelp() {
	fmt.Println(`
  ╔══════════════════════════════════════════════════╗
  ║               gitmap agent crashed               ║
  ╚══════════════════════════════════════════════════╝

  Detect crashed or abandoned agents and display autopsy report.

  Aliases:
    crash

  Usage:
    gitmap agent crashed [flags]

  Flags:
    -t, --task-id <id>     Parent task ID or slug
    -d, --dir <path>       Target directory
    --db <path>            Direct path to agent SQLite database
    -j, --json             Output report in JSON format
    -h, --help             Show this help menu

  Examples:
    gitmap agent crashed
    gitmap agent crashed -t task-12345
    gitmap agent crashed --json`)
}

func printAgentDiagnoseHelp() {
	fmt.Println(`
  ╔══════════════════════════════════════════════════╗
  ║               gitmap agent diagnose              ║
  ╚══════════════════════════════════════════════════╝

  Perform deep forensic inspection of agent tasks and logs.

  Aliases:
    diag, autopsy

  Usage:
    gitmap agent diagnose [flags]

  Flags:
    -t, --task-id <id>     Parent task ID or slug
    -d, --dir <path>       Target directory
    --db <path>            Direct path to agent SQLite database
    -j, --json             Output report in JSON format
    -h, --help             Show this help menu

  Examples:
    gitmap agent diagnose
    gitmap agent diagnose -t task-12345
    gitmap agent diagnose --json`)
}
