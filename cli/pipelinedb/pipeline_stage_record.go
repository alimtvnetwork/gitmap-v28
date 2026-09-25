package pipelinedb

// PipelineJobRecord represents an individual workflow stage/job execution record.
type PipelineJobRecord struct {
	JobId           uint64 `json:"jobId"`
	RunId           uint64 `json:"runId"`
	RepoSlug        string `json:"repoSlug"`
	JobName         string `json:"jobName"`
	Status          string `json:"status"`
	Conclusion      string `json:"conclusion"`
	StartedAt       string `json:"startedAt"`
	CompletedAt     string `json:"completedAt"`
	DurationSeconds int    `json:"durationSeconds"`
	JobUrl          string `json:"jobUrl"`
	CreatedAt       string `json:"createdAt"`
}

// PipelineSegmentRecord represents an individual step/segment within a workflow job.
type PipelineSegmentRecord struct {
	SegmentId       int64  `json:"segmentId"`
	RunId           uint64 `json:"runId"`
	JobName         string `json:"jobName"`
	StepName        string `json:"stepName"`
	StepNumber      int    `json:"stepNumber"`
	Status          string `json:"status"`
	Conclusion      string `json:"conclusion"`
	DurationSeconds int    `json:"durationSeconds"`
	CreatedAt       string `json:"createdAt"`
}

// PipelineStageSummary aggregates single-stage durations and calculates concurrency metrics.
type PipelineStageSummary struct {
	RunId             uint64              `json:"runId"`
	RepoSlug          string              `json:"repoSlug"`
	WorkflowName      string              `json:"workflowName"`
	Status            string              `json:"status"`
	Conclusion        string              `json:"conclusion"`
	WallClockSeconds  int                 `json:"wallClockSeconds"`
	StageSumSeconds   int                 `json:"stageSumSeconds"`
	SpeedupMultiplier float64             `json:"speedupMultiplier"`
	Jobs              []PipelineJobRecord `json:"jobs"`
}
