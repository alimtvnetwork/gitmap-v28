package clonefrom

import (
	"io"
	"sync"
)

// ConcurrentExecutionParams encapsulates parameters for concurrent clone execution.
type ConcurrentExecutionParams struct {
	Plan      Plan
	Cwd       string
	Progress  io.Writer
	BeforeRow BeforeRowHook
	Workers   int
}

// CloneFromDispatchParams captures parameters for dispatching concurrent clone workers.
type CloneFromDispatchParams struct {
	Plan      Plan
	Cwd       string
	BeforeRow BeforeRowHook
	Workers   int
	Out       []Result
}

// CloneFromWorkerParams captures arguments passed to a concurrent worker loop.
type CloneFromWorkerParams struct {
	Jobs <-chan concurrentJob
	Cwd  string
	Out  []Result
	Wg   *sync.WaitGroup
}

// BeforeRowInvokeParams encapsulates arguments for invoking the before-row hook.
type BeforeRowInvokeParams struct {
	Hook         BeforeRowHook
	CurrentIndex int
	TotalCount   int
	Row          Row
}

// ProgressWriteParams encapsulates parameters for rendering clone progress lines.
type ProgressWriteParams struct {
	Writer       io.Writer
	CurrentIndex int
	TotalCount   int
	Result       Result
}
