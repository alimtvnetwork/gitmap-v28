package cmdpull

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
