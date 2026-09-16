package cmdpull

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/glyphs"
	"github.com/alimtvnetwork/gitmap-v28/cli/theme"
)

// NewPullProgressBar initializes a new progress bar.
func NewPullProgressBar(total int, isQuiet, isStopOnFail bool) *PullProgressBar {
	isTerminal := resolveTerminalMode()
	isSafe := glyphs.Resolve() == glyphs.ModeSafe

	return &PullProgressBar{
		total:         total,
		isQuiet:       isQuiet,
		isStopOnFail:  isStopOnFail,
		isTTY:         isTerminal,
		isSafe:        isSafe,
		out:           os.Stdout,
		activeStep:    PullStepTypePending,
		subStepIndex:  1,
		subStepTotal:  4,
		activeWorkers: make(map[int]*WorkerSlotState),
		stopChan:      make(chan struct{}),
	}
}

func resolveTerminalMode() bool {
	return theme.IsStdoutTTY()
}

// Start marks progress beginning and launches active background ticker.
func (p *PullProgressBar) Start() {
	p.startOnce.Do(func() {
		p.initTickerLifecycle()
		p.tickerWg.Add(1)
		go p.runTickerLoop()
	})
}

// Stop finalizes progress and adds trailing newline if needed.
func (p *PullProgressBar) Stop() {
	p.stopOnce.Do(func() {
		p.shutdownTicker()
		p.flushFinalOutput()
	})
}
