// Package cmd — prompt_batch_progress.go tracks multi-repo prompt installation progress.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloner"
)

func NewPromptBatchProgress(count int) *cloner.BatchProgress {
	return cloner.NewBatchProgress(count, "Prompt-Architect", false)
}
