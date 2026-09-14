package macro

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type contextKey string

const callStackKey contextKey = "macroCallStack"

const maxMacroRecursionDepth = 10

// ParseRecurseCommand identifies if a step is a macro recursion invocation.
func ParseRecurseCommand(raw string) (bool, RecurseOpts) {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "gitmap ") {
		trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "gitmap"))
	}

	if isCall, opts := matchCallSyntax(trimmed); isCall {
		return true, opts
	}

	return matchMacroRunSyntax(trimmed)
}

func matchCallSyntax(trimmed string) (bool, RecurseOpts) {
	if !strings.HasPrefix(trimmed, "call ") {
		return false, RecurseOpts{}
	}

	rem := strings.TrimSpace(strings.TrimPrefix(trimmed, "call"))

	return true, parseRecurseTokens(rem)
}

func matchMacroRunSyntax(trimmed string) (bool, RecurseOpts) {
	if !strings.HasPrefix(trimmed, "macro run ") {
		return false, RecurseOpts{}
	}

	rem := strings.TrimSpace(strings.TrimPrefix(trimmed, "macro run"))

	return true, parseRecurseTokens(rem)
}

func parseRecurseTokens(rem string) RecurseOpts {
	parts := strings.Fields(rem)
	var opts RecurseOpts
	var nonFlags []string

	for i := 0; i < len(parts); i++ {
		p := parts[i]
		if (p == "-d" || p == "--delay") && i+1 < len(parts) {
			d, _ := time.ParseDuration(parts[i+1])
			opts.Delay = d
			i++
			continue
		}
		nonFlags = append(nonFlags, p)
	}

	if len(nonFlags) > 0 {
		opts.TargetMacro = cleanAsyncCommand(nonFlags[0])
	}

	return opts
}

func executeRecurseStep(
	ctx context.Context,
	step MacroStep,
	recOpts RecurseOpts,
	dt *DirTracker,
	start time.Time,
	opts ExecOptions,
	idx int,
) (StepExecution, error) {
	if execOptsCheck := opts.DryRun; execOptsCheck {
		return executeDryRunStep(step, idx, 1, opts, dt), nil
	}

	callStack := extractCallStack(ctx)
	if err := validateRecursionSafety(callStack, recOpts.TargetMacro); err != nil {
		return StepExecution{CommandLine: step.CommandLine, Status: "failed", Error: err.Error()}, err
	}

	applyRecurseDelay(recOpts.Delay, opts)
	childMacro, err := LoadMacro(recOpts.TargetMacro)
	if err != nil {
		return StepExecution{CommandLine: step.CommandLine, Status: "failed", Error: err.Error()}, apperror.WrapSimple(err, "load recursive macro")
	}

	childCtx := buildChildRecursionContext(ctx, callStack, recOpts.TargetMacro)
	execErr := Execute(childCtx, childMacro, opts)

	return buildRecurseResult(step, recOpts.TargetMacro, dt.CurrentDir, time.Since(start), execErr)
}

func extractCallStack(ctx context.Context) []string {
	if stack, ok := ctx.Value(callStackKey).([]string); ok {
		return stack
	}

	return []string{}
}

func validateRecursionSafety(stack []string, target string) error {
	if len(stack) >= maxMacroRecursionDepth {
		msg := fmt.Sprintf("maximum macro recursion depth (%d) exceeded", maxMacroRecursionDepth)

		return apperror.NewSimple(msg, "E_MACRO_RECURSION_DEPTH_EXCEEDED")
	}

	for _, s := range stack {
		if strings.EqualFold(s, target) {
			msg := fmt.Sprintf("macro cycle detected: %s -> %s", strings.Join(stack, " -> "), target)

			return apperror.NewSimple(msg, "E_MACRO_CYCLE_DETECTED")
		}
	}

	return nil
}

func applyRecurseDelay(d time.Duration, opts ExecOptions) {
	if d <= 0 {
		return
	}

	if !isStructuredOutput(opts) {
		fmt.Printf("  %s⏳ delaying recursive invocation by %v...%s\n", constants.ColorYellow, d, constants.ColorReset)
	}

	time.Sleep(d)
}

func buildChildRecursionContext(
	ctx context.Context,
	stack []string,
	target string,
) context.Context {
	newStack := make([]string, len(stack), len(stack)+1)
	copy(newStack, stack)
	newStack = append(newStack, target)

	return context.WithValue(ctx, callStackKey, newStack)
}

func buildRecurseResult(
	step MacroStep,
	target, dir string,
	elapsed time.Duration,
	err error,
) (StepExecution, error) {
	status := "success"
	errStr := ""
	if err != nil {
		status = "failed"
		errStr = err.Error()
	}

	return StepExecution{
		StepNum:        step.StepNum,
		CommandLine:    fmt.Sprintf("call %s", target),
		WorkingDir:     dir,
		Status:         status,
		ExitCode:       0,
		ElapsedSeconds: elapsed.Seconds(),
		Logs:           []string{fmt.Sprintf("executed child macro %q", target)},
		Error:          errStr,
		ErrorLogs:      []string{},
	}, err
}
