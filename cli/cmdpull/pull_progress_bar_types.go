package cmdpull

import (
	"io"
	"sync"
	"time"
)

var brailleFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
var asciiFrames = []string{"-", "\\", "|", "/"}

// WorkerSlotState tracks active execution progress for a single worker thread.
type WorkerSlotState struct {
	WorkerID     int
	RepoName     string
	Step         PullStepType
	Description  string
	SubStepIndex int
	SubStepTotal int
	Percent      int
	StartTime    time.Time
}

// PullProgressBar manages thread-safe terminal progress tracking for batch pulls.
type PullProgressBar struct {
	mu            sync.Mutex
	startOnce     sync.Once
	stopOnce      sync.Once
	total         int
	completed     int
	succeeded     int
	failed        int
	skipped       int
	upToDate      int
	isQuiet       bool
	isStopOnFail  bool
	isStopped     bool
	isTTY         bool
	isSafe        bool
	activeRepo    string
	activeStep    PullStepType
	activeDesc    string
	activePercent int
	subStepIndex  int
	subStepTotal  int
	spinnerIdx    int
	ticker        *time.Ticker
	stopChan      chan struct{}
	tickerWg      sync.WaitGroup
	startTime     time.Time
	states        []*PullRepoState
	activeWorkers map[int]*WorkerSlotState
	out           io.Writer
}

// SetOutput overrides the destination writer (used in tests).
func (p *PullProgressBar) SetOutput(w io.Writer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.out = w
}

// SetTTY toggles interactive TTY mode.
func (p *PullProgressBar) SetTTY(isTTY bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isTTY = isTTY
}
