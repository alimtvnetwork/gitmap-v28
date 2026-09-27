package cmdpipeline

import (
	"sync"
)

var (
	explicitAgyFixOnce sync.Once
	rawAgyFixRunnerFn  func(args []string) error
)

// EnsureAgyFixRunnerExplicitOnly wraps PipelineAgyFixRunner to prevent automatic background dispatches.
func EnsureAgyFixRunnerExplicitOnly() {
	explicitAgyFixOnce.Do(func() {
		setupExplicitOnlyRunner()
	})
}

func setupExplicitOnlyRunner() {
	if PipelineAgyFixRunner == nil {
		return
	}
	rawAgyFixRunnerFn = PipelineAgyFixRunner
	PipelineAgyFixRunner = func(args []string) error {
		if !IsExplicitAgyFixArgs(args) {
			return nil
		}
		if rawAgyFixRunnerFn != nil {
			return rawAgyFixRunnerFn(args)
		}
		return nil
	}
}

// IsExplicitAgyFixArgs verifies if pipeline fix invocation was explicitly requested by user.
func IsExplicitAgyFixArgs(args []string) bool {
	if len(args) == 0 {
		return false
	}
	return IsPipelineFixAgyArgs(args)
}

// ExecutePipelineCmd wraps standard pipeline invocation ensuring read-only safety for pe.
func ExecutePipelineCmd(args []string) error {
	EnsureAgyFixRunnerExplicitOnly()
	return runPipeline(args)
}

func init() {
	EnsureAgyFixRunnerExplicitOnly()
}
