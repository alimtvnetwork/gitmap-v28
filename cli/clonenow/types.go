package clonenow

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

// CloneIdempotentParams captures input arguments for idempotent clone dispatching.
type CloneIdempotentParams struct {
	Row     Row
	URL     string
	AbsDest string
	Cwd     string
	Policy  string
	State   existingRepoState
}

// ConcurrentDispatchParams captures parameters for dispatching concurrent clone workers.
type ConcurrentDispatchParams struct {
	Plan      Plan
	Cwd       string
	BeforeRow BeforeRowHook
	Workers   int
	Out       []Result
}

// ConcurrentWorkerParams captures arguments passed to a concurrent worker loop.
type ConcurrentWorkerParams struct {
	Jobs <-chan concurrentJob
	Plan Plan
	Cwd  string
	Out  []Result
	Wg   *sync.WaitGroup
}

// GitCloneParams encapsulates arguments for running a git clone command.
type GitCloneParams struct {
	Row  Row
	URL  string
	Dest string
	Cwd  string
}

// ProgressWriteParams encapsulates parameters for rendering clone progress lines.
type ProgressWriteParams struct {
	Writer       io.Writer
	CurrentIndex int
	TotalCount   int
	Result       Result
}
