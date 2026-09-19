package cmdai

import (
	"context"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunAiScript resolves a script token and executes it with forwarded arguments.
func RunAiScript(token string, forwardArgs []string) *apperror.AppError {
	return RunAiScriptContext(context.Background(), token, forwardArgs)
}

// RunAiScriptContext executes a resolved AI script under the provided context.
func RunAiScriptContext(ctx context.Context, token string, forwardArgs []string) *apperror.AppError {
	meta, findErr := FindScriptByToken(token)
	hasFindErr := findErr != nil
	if hasFindErr {
		return findErr
	}

	scriptPath, pathErr := ResolveScriptPath(meta.Filename)
	hasPathErr := pathErr != nil
	if hasPathErr {
		return pathErr
	}

	args := resolveForwardArgs(meta, forwardArgs)

	return ExecuteScriptStreaming(ctx, scriptPath, args)
}

// RunAiScriptDirect executes a script directly by filename without catalog lookup.
func RunAiScriptDirect(filename string, forwardArgs []string) *apperror.AppError {
	scriptPath, pathErr := ResolveScriptPath(filename)
	hasPathErr := pathErr != nil
	if hasPathErr {
		return pathErr
	}

	return ExecuteScriptStreaming(context.Background(), scriptPath, forwardArgs)
}

func resolveForwardArgs(meta ScriptMetadata, userArgs []string) []string {
	hasUserArgs := len(userArgs) > 0
	if hasUserArgs {
		return userArgs
	}

	return nil
}
