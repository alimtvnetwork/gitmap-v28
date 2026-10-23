package cmdagy

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func findPromptRecordByID(records []PromptRecord, id string) *PromptRecord {
	for _, rec := range records {
		if rec.ID == id {
			r := rec
			return &r
		}
	}

	return nil
}

func dispatchUrgentResend(payload PromptPayload, targetNode string) (*PromptResult, error) {
	if _, err := SendPromptImmediate(payload); err != nil {
		return nil, err
	}

	return &PromptResult{
		PromptID:      fmt.Sprintf("resend-%d", time.Now().UnixNano()),
		Status:        "delivered",
		TargetNode:    targetNode,
		ProjectName:   payload.EffectiveTarget(),
		DispatchedAt:  time.Now().Format(time.RFC3339),
		OutputPreview: "Prompt re-dispatched immediately to " + targetNode,
	}, nil
}

// ResendPrompt dispatches an existing prompt from history or template.
func ResendPrompt(req PromptResendRequest) (*PromptResult, error) {
	history, _ := GetPromptHistory()
	targetRecord := findPromptRecordByID(history, req.PromptID)
	if targetRecord == nil {
		saved, _ := LoadSavedPrompts()
		targetRecord = findPromptRecordByID(saved, req.PromptID)
	}
	if targetRecord == nil {
		return nil, apperror.NewWithDetails("agy.resend", "E404", "prompt not found", "cmdagy", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
	}

	targetNode := targetRecord.TargetNode
	if req.OverrideNode != "" {
		targetNode = req.OverrideNode
	}
	if targetNode == "" {
		targetNode, _ = os.Hostname()
	}

	payload := PromptPayload{
		ProjectName:   targetRecord.ProjectName,
		ProjectTarget: targetRecord.ProjectTarget,
		TargetNode:    targetNode,
		Title:         "Resend: " + targetRecord.Title,
		PromptText:    targetRecord.PromptText,
		IsUrgent:      req.IsUrgent,
	}

	if req.IsUrgent {
		return dispatchUrgentResend(payload, targetNode)
	}

	queuedRec, err := EnqueuePromptPayload(payload)
	if err != nil {
		return nil, err
	}
	return &PromptResult{
		PromptID:      queuedRec.ID,
		Status:        "queued",
		TargetNode:    targetNode,
		ProjectName:   payload.EffectiveTarget(),
		DispatchedAt:  time.Now().Format(time.RFC3339),
		OutputPreview: fmt.Sprintf("Prompt enqueued at position %d", queuedRec.QueuePosition),
	}, nil
}

var (
	queueWorkerOnce sync.Once
	queueWorkerStop chan struct{}
)

// StartPromptQueueWorker begins continuous background processing of queued prompts.
func StartPromptQueueWorker() {
	queueWorkerOnce.Do(func() {
		queueWorkerStop = make(chan struct{})
		go runPromptQueueWorkerLoop(queueWorkerStop)
	})
}

// StopPromptQueueWorker signals the background queue worker to terminate.
func StopPromptQueueWorker() {
	if queueWorkerStop != nil {
		close(queueWorkerStop)
		queueWorkerStop = nil
	}
}

func runPromptQueueWorkerLoop(stopChan <-chan struct{}) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	concurrencyLimit := make(chan struct{}, 2)

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			processNextQueuedPrompts(concurrencyLimit)
		}
	}
}

func processNextQueuedPrompts(sem chan struct{}) {
	queuePath := getTempAgyFilePath("agy-prompt-queue.json")
	records, err := readRecordsFromFile(queuePath)
	if err != nil || len(records) == 0 {
		return
	}

	var remaining []PromptRecord
	for i, rec := range records {
		switch {
		case rec.Status == "queued" && len(sem) < cap(sem):
			sem <- struct{}{}
			go func(r PromptRecord) {
				defer func() { <-sem }()
				executeQueuedPrompt(r)
			}(rec)
		case rec.Status == "queued":
			rec.QueuePosition = len(remaining) + 1
			remaining = append(remaining, rec)
		case i > 0:
			remaining = append(remaining, rec)
		}
	}
	_ = writeRecordsToFile(queuePath, remaining)
}

func executeQueuedPrompt(rec PromptRecord) {
	payload := PromptPayload{
		ProjectName:   rec.ProjectName,
		ProjectTarget: rec.ProjectTarget,
		TargetNode:    rec.TargetNode,
		Title:         rec.Title,
		PromptText:    rec.PromptText,
	}
	_, _ = SendPromptImmediate(payload)
	rec.Status = "completed"
	rec.CompletedAt = time.Now().Format(time.RFC3339)
	_ = AppendPromptHistory(rec)
}
