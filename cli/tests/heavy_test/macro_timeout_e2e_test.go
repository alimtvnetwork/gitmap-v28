package heavy_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func TestExecute_StepTimeout(t *testing.T) {
	m := &macro.Macro{
		Name: "test-timeout-macro",
		Steps: []macro.MacroStep{
			{StepNum: 1, CommandLine: getMacroSleepCmd(3), TimeoutSeconds: 1},
		},
	}

	start := time.Now()
	err := macro.Execute(context.Background(), m, macro.ExecOptions{DryRun: false})
	elapsed := time.Since(start)
	if err == nil || elapsed > 2800*time.Millisecond {
		t.Fatalf("expected timeout error under 2.8s, got err: %v, took: %v", err, elapsed)
	}
}

func getMacroSleepCmd(seconds int) string {
	if runtime.GOOS == constants.OSWindows {
		return fmt.Sprintf("Start-Sleep -Seconds %d", seconds)
	}

	return fmt.Sprintf("exec sleep %d", seconds)
}
