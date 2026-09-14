package cmdpull

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/glyphs"
	"github.com/alimtvnetwork/gitmap-v28/cli/theme"
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

// Start marks progress beginning and launches active background ticker.
func (p *PullProgressBar) Start() {
	p.startOnce.Do(func() {
		p.initTickerLifecycle()
		p.tickerWg.Add(1)
		go p.runTickerLoop()
	})
}

func (p *PullProgressBar) initTickerLifecycle() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.startTime = time.Now()
	p.ticker = time.NewTicker(80 * time.Millisecond)
	p.renderLocked()
}

func (p *PullProgressBar) runTickerLoop() {
	defer p.tickerWg.Done()
	for {
		select {
		case <-p.stopChan:
			return
		case <-p.ticker.C:
			p.tick()
		}
	}
}

func (p *PullProgressBar) tick() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isStopped {
		return
	}
	p.spinnerIdx++
	if p.isTTY {
		p.renderLocked()
	}
}

// Stop finalizes progress and adds trailing newline if needed.
func (p *PullProgressBar) Stop() {
	p.stopOnce.Do(func() {
		p.shutdownTicker()
		p.flushFinalOutput()
	})
}

func (p *PullProgressBar) shutdownTicker() {
	p.mu.Lock()
	p.isStopped = true
	if p.ticker != nil {
		p.ticker.Stop()
	}
	if p.stopChan != nil {
		close(p.stopChan)
	}
	p.mu.Unlock()
	p.tickerWg.Wait()
}

func (p *PullProgressBar) flushFinalOutput() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isTTY && !p.isQuiet {
		fmt.Fprintln(p.out)
	}
}

// RegisterWorker records an active worker slot.
func (p *PullProgressBar) RegisterWorker(workerID int, repoName string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.activeWorkers[workerID] = &WorkerSlotState{
		WorkerID:     workerID,
		RepoName:     repoName,
		Step:         PullStepTypeInspecting,
		Description:  "inspecting",
		SubStepIndex: 1,
		SubStepTotal: 4,
		StartTime:    time.Now(),
	}
}

// UpdateWorkerProgress updates the state and sub-step of a worker slot.
func (p *PullProgressBar) UpdateWorkerProgress(workerID int, step PullStepType, desc string, percent int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	slot, isFound := p.activeWorkers[workerID]
	if !isFound {
		return
	}
	slot.Step = step
	slot.Description = desc
	slot.Percent = percent
	slot.SubStepIndex = StepSubStepIndex(step)
}

// UnregisterWorker removes a worker from the active registry.
func (p *PullProgressBar) UnregisterWorker(workerID int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.activeWorkers, workerID)
}

// ActiveWorkers returns a slice of active worker slots.
func (p *PullProgressBar) ActiveWorkers() []*WorkerSlotState {
	p.mu.Lock()
	defer p.mu.Unlock()
	slots := make([]*WorkerSlotState, 0, len(p.activeWorkers))
	for _, slot := range p.activeWorkers {
		slots = append(slots, slot)
	}

	return slots
}

// UpdateRepoStep updates active repository step and refreshes render.
func (p *PullProgressBar) UpdateRepoStep(repoName string, step PullStepType, desc string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.activeRepo = repoName
	p.activeStep = step
	p.activeDesc = desc
	p.subStepIndex = StepSubStepIndex(step)
	p.renderLocked()
}

// SetActivePercent sets the live streaming percent for single repo.
func (p *PullProgressBar) SetActivePercent(percent int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.activePercent = percent
}

// SetSubStep sets the explicit single-repo sub-step milestone indices.
func (p *PullProgressBar) SetSubStep(index, total int, desc string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.subStepIndex = index
	p.subStepTotal = total
	if desc != "" {
		p.activeDesc = desc
	}
}

// SubStepMilestone returns current sub-step index and total.
func (p *PullProgressBar) SubStepMilestone() (int, int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.subStepIndex, p.subStepTotal
}

// CompleteRepo records completed repository state and advances bar.
func (p *PullProgressBar) CompleteRepo(state *PullRepoState) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.completed++
	if state != nil {
		p.states = append(p.states, state)
		p.applyStateMetrics(state)
	}
	p.renderLocked()
}

func (p *PullProgressBar) applyStateMetrics(state *PullRepoState) {
	p.activeRepo = state.RepoName
	p.activeStep = state.Step
	p.activeDesc = state.Changes
	p.subStepIndex = 4
	p.incrementStepCounters(state.Step)
	if state.Step == PullStepTypeError && p.isStopOnFail {
		p.isStopped = true
	}
}

func (p *PullProgressBar) incrementStepCounters(step PullStepType) {
	if p.incrementSuccessCounters(step) {
		return
	}
	p.incrementFailureOrSkipCounters(step)
}

func (p *PullProgressBar) incrementSuccessCounters(step PullStepType) bool {
	if step == PullStepTypeUpToDate {
		p.upToDate++
		p.succeeded++
		return true
	}
	if step == PullStepTypeFastForward || step == PullStepTypeMerging {
		p.succeeded++
		return true
	}

	return false
}

func (p *PullProgressBar) incrementFailureOrSkipCounters(step PullStepType) {
	if step == PullStepTypeError || step == PullStepTypeConflict {
		p.failed++
		return
	}
	if step == PullStepTypeSkipped {
		p.skipped++
	}
}

// Render redraws the bar immediately.
func (p *PullProgressBar) Render() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.renderLocked()
}

func (p *PullProgressBar) renderLocked() {
	if p.isQuiet {
		return
	}
	if p.isTTY {
		p.renderTTY()
		return
	}
	p.renderNonTTY()
}

func (p *PullProgressBar) currentSpinner() string {
	if p.isSafe {
		return asciiFrames[p.spinnerIdx%len(asciiFrames)]
	}

	return brailleFrames[p.spinnerIdx%len(brailleFrames)]
}

func (p *PullProgressBar) renderTTY() {
	if p.total <= 1 {
		p.renderSingleRepoTTY()
		return
	}
	p.renderMultiRepoTTY()
}

func (p *PullProgressBar) renderSingleRepoTTY() {
	percent := p.resolveSingleRepoPercent()
	bar := FormatVisualBar(percent, 100, 20, p.isSafe)
	spinner := p.currentSpinner()
	badge := FormatStepBadge(p.activeStep, p.isSafe)
	stepLabel := fmt.Sprintf("[Step %d/%d] %s", p.subStepIndex, p.subStepTotal, StepMilestoneTitle(p.subStepIndex, p.activeStep))
	elapsed := fmt.Sprintf("%.1fs", time.Since(p.startTime).Seconds())
	desc := p.formatActiveDescription()

	fmt.Fprintf(p.out, "\r\033[K%s %s %3d%% %s | %s %s%s (%s)", spinner, bar, percent, stepLabel, badge, p.activeRepo, desc, elapsed)
}

func (p *PullProgressBar) resolveSingleRepoPercent() int {
	if p.completed >= 1 {
		return 100
	}
	if p.activePercent > 0 {
		return p.activePercent
	}

	return StepMilestonePercent(p.activeStep)
}

func (p *PullProgressBar) formatActiveDescription() string {
	if p.activeDesc == "" {
		return ""
	}

	return ": " + p.activeDesc
}

func (p *PullProgressBar) renderMultiRepoTTY() {
	percent := calcProgressPercent(p.completed, p.total)
	bar := FormatVisualBar(p.completed, p.total, 20, p.isSafe)
	spinner := p.currentSpinner()
	counter := fmt.Sprintf("%3d%% (%d/%d repos)", percent, p.completed, p.total)
	elapsed := fmt.Sprintf("%.1fs", time.Since(p.startTime).Seconds())
	summary := p.formatWorkersSummary()

	fmt.Fprintf(p.out, "\r\033[K%s %s %s | %s (%s)", spinner, bar, counter, summary, elapsed)
}

func (p *PullProgressBar) formatWorkersSummary() string {
	if len(p.activeWorkers) == 0 {
		return p.formatFallbackSummary()
	}

	return p.formatActiveWorkersList()
}

func (p *PullProgressBar) formatFallbackSummary() string {
	badge := FormatStepBadge(p.activeStep, p.isSafe)
	desc := p.formatActiveDescription()

	return fmt.Sprintf("%s %s%s", badge, p.activeRepo, desc)
}

func (p *PullProgressBar) formatActiveWorkersList() string {
	workerIDs := p.sortedWorkerIDs()
	var parts []string
	for _, id := range workerIDs {
		w := p.activeWorkers[id]
		badge := FormatStepBadge(w.Step, p.isSafe)
		desc := formatWorkerSlotDesc(w)
		parts = append(parts, fmt.Sprintf("%s %s%s", badge, w.RepoName, desc))
	}

	return strings.Join(parts, ", ")
}

func formatWorkerSlotDesc(w *WorkerSlotState) string {
	if w.Description == "" {
		return ""
	}

	return ": " + w.Description
}

func (p *PullProgressBar) sortedWorkerIDs() []int {
	ids := make([]int, 0, len(p.activeWorkers))
	for id := range p.activeWorkers {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	return ids
}

func (p *PullProgressBar) renderNonTTY() {
	percent := calcProgressPercent(p.completed, p.total)
	bar := FormatVisualBar(p.completed, p.total, 20, p.isSafe)
	counter := fmt.Sprintf("%3d%% (%d/%d repos)", percent, p.completed, p.total)
	badge := FormatStepBadge(p.activeStep, p.isSafe)
	desc := p.formatActiveDescription()

	fmt.Fprintf(p.out, "%s %s | %s %s%s\n", bar, counter, badge, p.activeRepo, desc)
}

func calcProgressPercent(completed, total int) int {
	if total <= 0 || completed <= 0 {
		return 0
	}
	if completed >= total {
		return 100
	}

	return (completed * 100) / total
}

// FormatVisualBar returns ASCII or block visual bar.
func FormatVisualBar(completed, total, barWidth int, isSafe bool) string {
	if barWidth <= 0 {
		barWidth = 20
	}
	filledLen := calcFilledBarLength(completed, total, barWidth)
	emptyLen := barWidth - filledLen
	fillChar, emptyChar := resolveBarChars(isSafe)

	return fmt.Sprintf("[%s%s]", strings.Repeat(fillChar, filledLen), strings.Repeat(emptyChar, emptyLen))
}

func calcFilledBarLength(completed, total, barWidth int) int {
	if total <= 0 || completed <= 0 {
		return 0
	}
	if completed >= total {
		return barWidth
	}

	return (completed * barWidth) / total
}

func resolveBarChars(isSafe bool) (string, string) {
	if isSafe {
		return constants.ProgressBarFilledSafe, constants.ProgressBarEmptySafe
	}

	return constants.ProgressBarFilledRich, constants.ProgressBarEmptyRich
}

// IsStopped returns true when early termination was triggered.
func (p *PullProgressBar) IsStopped() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.isStopped
}

// States returns all recorded completed states.
func (p *PullProgressBar) States() []*PullRepoState {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.states
}

// Failed returns failure count.
func (p *PullProgressBar) Failed() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.failed
}

// Succeeded returns success count.
func (p *PullProgressBar) Succeeded() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.succeeded
}
