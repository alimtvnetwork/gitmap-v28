// Package cmdprompt — prompt_batch_progress.go tracks multi-repo prompt installation progress.
package cmdprompt

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
)

func NewPromptBatchProgress(count int) *cloner.BatchProgress {
	return cloner.NewBatchProgress(count, "Prompt-Architect", false)
}
