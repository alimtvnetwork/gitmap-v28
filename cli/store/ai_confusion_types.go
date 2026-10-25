package store

import (
	"time"
)

// AiConfusionRecord represents a persistent AI confusion episode in SQLite.
type AiConfusionRecord struct {
	AiConfusionId     int64     `json:"aiConfusionId,omitempty"`
	ConfusionUuid     string    `json:"confusionUuid"`
	AgentSlug         string    `json:"agentSlug"`
	FileNames         string    `json:"fileNames"`
	MdFiles           string    `json:"mdFiles"`
	Spec              string    `json:"spec"`
	Task              string    `json:"task"`
	Subtask           string    `json:"subtask"`
	HeatmapInfo       string    `json:"heatmapInfo"`
	LineNumbers       string    `json:"lineNumbers"`
	CodeSnippet       string    `json:"codeSnippet"`
	Reason            string    `json:"reason"`
	InitialIntention  string    `json:"initialIntention"`
	ReasoningThinking string    `json:"reasoningThinking"`
	ConfusionQuestion string    `json:"confusionQuestion"`
	IsCompleted       bool      `json:"isCompleted"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// AiConfusionFilter provides structured query constraints for listing confusion episodes.
type AiConfusionFilter struct {
	Id          string `json:"id,omitempty"`
	AgentSlug   string `json:"agentSlug,omitempty"`
	Task        string `json:"task,omitempty"`
	IsCompleted *bool  `json:"isCompleted,omitempty"`
	Limit       int    `json:"limit,omitempty"`
}

// AiConfusionJSONInput defines the ingestion payload when adding episodes via JSON files.
type AiConfusionJSONInput struct {
	ConfusionUuid     string `json:"confusionUuid,omitempty"`
	AgentSlug         string `json:"agentSlug,omitempty"`
	Agent             string `json:"agent,omitempty"`
	FileNames         string `json:"fileNames,omitempty"`
	FilePath          string `json:"filePath,omitempty"`
	MdFiles           string `json:"mdFiles,omitempty"`
	Spec              string `json:"spec,omitempty"`
	SpecPath          string `json:"specPath,omitempty"`
	Task              string `json:"task,omitempty"`
	TaskSlug          string `json:"taskSlug,omitempty"`
	Subtask           string `json:"subtask,omitempty"`
	SubtaskSlug       string `json:"subtaskSlug,omitempty"`
	HeatmapInfo       string `json:"heatmapInfo,omitempty"`
	LineNumbers       string `json:"lineNumbers,omitempty"`
	CodeSnippet       string `json:"codeSnippet,omitempty"`
	Reason            string `json:"reason,omitempty"`
	InitialIntention  string `json:"initialIntention,omitempty"`
	AiIntention       string `json:"aiIntention,omitempty"`
	ReasoningThinking string `json:"reasoningThinking,omitempty"`
	AiThinking        string `json:"aiThinking,omitempty"`
	ConfusionQuestion string `json:"confusionQuestion,omitempty"`
	Question          string `json:"question,omitempty"`
	IsCompleted       bool   `json:"isCompleted,omitempty"`
}

// ToRecord converts an ingested JSON input to a normalized AiConfusionRecord.
func (in AiConfusionJSONInput) ToRecord() AiConfusionRecord {
	rec := in.populateBaseRecord()
	in.populateNarrativeFields(&rec)

	return rec
}

func (in AiConfusionJSONInput) populateBaseRecord() AiConfusionRecord {
	return AiConfusionRecord{
		ConfusionUuid: in.ConfusionUuid,
		AgentSlug:     resolveFirstNonEmpty(in.AgentSlug, in.Agent),
		FileNames:     resolveFirstNonEmpty(in.FileNames, in.FilePath),
		MdFiles:       in.MdFiles,
		Spec:          resolveFirstNonEmpty(in.Spec, in.SpecPath),
		Task:          resolveFirstNonEmpty(in.Task, in.TaskSlug),
		Subtask:       resolveFirstNonEmpty(in.Subtask, in.SubtaskSlug),
		HeatmapInfo:   in.HeatmapInfo,
		LineNumbers:   in.LineNumbers,
		CodeSnippet:   in.CodeSnippet,
		Reason:        in.Reason,
		IsCompleted:   in.IsCompleted,
	}
}

func (in AiConfusionJSONInput) populateNarrativeFields(rec *AiConfusionRecord) {
	rec.InitialIntention = resolveFirstNonEmpty(in.InitialIntention, in.AiIntention)
	rec.ReasoningThinking = resolveFirstNonEmpty(in.ReasoningThinking, in.AiThinking)
	rec.ConfusionQuestion = resolveFirstNonEmpty(in.ConfusionQuestion, in.Question)
}

func resolveFirstNonEmpty(a, b string) string {
	hasA := len(a) > 0
	if hasA {
		return a
	}

	return b
}
