package cmdssh

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type sshErrorLogsFlags struct {
	asJSON   bool
	filePath string
	tempFile string
	clear    bool
}

func parseSSHErrorFlags(args []string) sshErrorLogsFlags {
	var f sshErrorLogsFlags
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--json":
			f.asJSON = true
		case arg == "--clear" || arg == "clear":
			f.clear = true
		case arg == "--file" && i+1 < len(args):
			f.filePath = args[i+1]
			i++
		case arg == "--tempfile" && i+1 < len(args):
			f.tempFile = args[i+1]
			i++
		}
	}
	return f
}

func outputSSHLogsData(content []byte, f sshErrorLogsFlags) error {
	if f.filePath != "" {
		return os.WriteFile(f.filePath, content, 0644)
	}
	if f.tempFile != "" {
		p := filepath.Join(".ai-memory", "temp", f.tempFile)
		_ = os.MkdirAll(filepath.Dir(p), 0755)
		return os.WriteFile(p, content, 0644)
	}

	fmt.Println(string(content))
	return nil
}

func loadOrActiveTrace() *SSHExecutionTrace {
	if trace := GetLastSSHTrace(); trace != nil {
		return trace
	}
	logPath := resolveSSHLogPath()
	data, err := os.ReadFile(logPath)
	if err != nil {
		return nil
	}
	var loaded SSHExecutionTrace
	if err := json.Unmarshal(data, &loaded); err != nil {
		return nil
	}
	return &loaded
}

func runSSHErrorLogsCLI(args []string) error {
	flags := parseSSHErrorFlags(args)
	if flags.clear {
		_ = os.Remove(resolveSSHLogPath())
		fmt.Println("✓ SSH error logs cleared.")
		return nil
	}

	trace := loadOrActiveTrace()
	if trace == nil {
		fmt.Println("ℹ No SSH error logs found. Last operations succeeded or no logs recorded.")
		return nil
	}

	if flags.asJSON {
		data, _ := json.MarshalIndent(trace, "", "  ")
		return outputSSHLogsData(data, flags)
	}

	report := strings.TrimRight(trace.FormatTerminalReport(), "\n")
	return outputSSHLogsData([]byte(report), flags)
}
