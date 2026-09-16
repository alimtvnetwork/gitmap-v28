package cmdpull

import (
	"time"
)

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
