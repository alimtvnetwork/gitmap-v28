package cmdai

import (
	"context"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunAiFix executes a target autofix routine or composite pipeline.
func RunAiFix(target string, extraArgs []string) *apperror.AppError {
	return RunAiFixContext(context.Background(), target, extraArgs)
}

// RunAiFixContext executes an autofix routine under the provided context.
func RunAiFixContext(ctx context.Context, target string, extraArgs []string) *apperror.AppError {
	def, resolveErr := ResolveFixTarget(target)
	hasResolveErr := resolveErr != nil
	if hasResolveErr {
		return resolveErr
	}

	isAll := def.Target == FixTargetAll
	if isAll {
		return RunAiFixAll(ctx, extraArgs)
	}

	return RunAiFixTarget(ctx, def, extraArgs)
}

// RunAiFixTarget executes a single resolved autofix target.
func RunAiFixTarget(ctx context.Context, def FixTargetDef, extraArgs []string) *apperror.AppError {
	renderFixTargetHeader(def)
	args := MergeFixArgs(def, extraArgs)

	return RunAiScriptContext(ctx, def.ScriptToken, args)
}

// RunAiFixAll sequentially executes all registered repository autofix targets.
func RunAiFixAll(ctx context.Context, extraArgs []string) *apperror.AppError {
	renderFixAllStart()
	targets := AllFixTargetDefs()
	for _, t := range targets {
		fixErr := RunAiFixTarget(ctx, t, extraArgs)
		hasFixErr := fixErr != nil
		if hasFixErr {
			return fixErr
		}
	}

	renderFixAllSuccess()

	return nil
}

func renderFixTargetHeader(def FixTargetDef) {
	fmt.Printf("\n▶ AI Fix: %s (%s)\n", def.Target, def.Description)
}

func renderFixAllStart() {
	fmt.Println("\n🚀 Launching Composite AI Autofix Suite...")
}

func renderFixAllSuccess() {
	fmt.Println("\n✔ Composite AI Autofix Suite completed successfully.")
}
