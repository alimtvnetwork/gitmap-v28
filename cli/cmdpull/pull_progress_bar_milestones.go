package cmdpull

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
