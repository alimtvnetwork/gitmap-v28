package macro

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// RecordInteractive starts an interactive shell session and records commands to a macro.
func RecordInteractive(name string) error {
	fmt.Printf("  %s● [REC] Recording session: %q%s\n", constants.ColorRed, name, constants.ColorReset)
	fmt.Println("  Type shell commands. They will execute live and be saved to the macro.")
	fmt.Println("  Type 'stop' or 'exit' when finished.")
	fmt.Println()

	cwd, _ := os.Getwd()
	m := &Macro{
		Name:      name,
		CreatedAt: time.Now(),
		Steps:     make([]MacroStep, 0),
	}

	reader := bufio.NewReader(os.Stdin)
	stepNum := 1

	for {
		fmt.Printf("  %s[rec:%s %d]> %s", constants.ColorCyan, name, stepNum, constants.ColorReset)
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		cmdText := strings.TrimSpace(line)
		if len(cmdText) == 0 {
			continue
		}
		if cmdText == "stop" || cmdText == "exit" || cmdText == "quit" {
			break
		}

		execLive(cmdText, cwd)

		m.Steps = append(m.Steps, MacroStep{
			StepNum:        stepNum,
			CommandLine:    cmdText,
			WorkingDir:     cwd,
			TimeoutSeconds: 300,
		})
		stepNum++
	}

	if err := SaveMacro(m); err != nil {
		return fmt.Errorf("could not save macro: %w", err)
	}

	fmt.Printf("\n  %s✔ Saved macro %q with %d steps.%s\n\n",
		constants.ColorGreen, name, len(m.Steps), constants.ColorReset)
	return nil
}

func execLive(cmdText, dir string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(context.Background(), "powershell", "-NoProfile", "-Command", cmdText)
	} else {
		cmd = exec.CommandContext(context.Background(), "sh", "-c", cmdText)
	}
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
