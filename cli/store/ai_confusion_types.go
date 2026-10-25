package store

import (
	"time"
)

// AiConfusionFileRecord represents a referenced file associated with an AI confusion episode.
type AiConfusionFileRecord struct {
	AiConfusionFileId int64  `json:"aiConfusionFileId,omitempty"`
	AiConfusionId     int64  `json:"aiConfusionId,omitempty"`
	ConfusionUuid     string `json:"confusionUuid"`
	FilePath          string `json:"filePath"`
	FileType          string `json:"fileType"`
	LineNumbers       string `json:"lineNumbers"`
	CodeSnippet       string `json:"codeSnippet"`
	WhyModified       string `json:"whyModified"`
	CreatedAt         string `json:"createdAt,omitempty"`
}

// AiConfusionRecord represents a persistent AI confusion episode in SQLite.
type AiConfusionRecord struct {
	AiConfusionId     int64                   `json:"aiConfusionId,omitempty"`
	ConfusionUuid     string                  `json:"confusionUuid"`
	AgentSlug         string                  `json:"agentSlug"`
	FileNames         string                  `json:"fileNames"`
	MdFiles           string                  `json:"mdFiles"`
	Spec              string                  `json:"spec"`
	Task              string                  `json:"task"`
	Subtask           string                  `json:"subtask"`
	HeatmapInfo       string                  `json:"heatmapInfo"`
	LineNumbers       string                  `json:"lineNumbers"`
	CodeSnippet       string                  `json:"codeSnippet"`
	Reason            string                  `json:"reason"`
	InitialIntention  string                  `json:"initialIntention"`
	ReasoningThinking string                  `json:"reasoningThinking"`
	ConfusionQuestion string                  `json:"confusionQuestion"`
	IsCompleted       bool                    `json:"isCompleted"`
	ReferencedFiles   []AiConfusionFileRecord `json:"referencedFiles,omitempty"`
	CreatedAt         time.Time               `json:"createdAt"`
	UpdatedAt         time.Time               `json:"updatedAt"`
}

// AiConfusionFilter provides structured query constraints for listing confusion episodes.
type AiConfusionFilter struct {
	Id          string `json:"id,omitempty"`
	AgentSlug   string `json:"agentSlug,omitempty"`
	Task        string `json:"task,omitempty"`
	IsCompleted *bool  `json:"isCompleted,omitempty"`
	Limit       int    `json:"limit,omitempty"`
}

// AiConfusionFileInput captures file details in the JSON ingestion payload.
type AiConfusionFileInput struct {
	FilePath    string `json:"filePath"`
	FileType    string `json:"fileType,omitempty"`
	LineNumbers string `json:"lineNumbers,omitempty"`
	CodeSnippet string `json:"codeSnippet,omitempty"`
	WhyModified string `json:"whyModified,omitempty"`
}

// ToRecord converts an ingested file input to an AiConfusionFileRecord.
func (in AiConfusionFileInput) ToRecord(uuid string) AiConfusionFileRecord {
	fileType := resolveFirstNonEmpty(in.FileType, "SOURCE")

	return AiConfusionFileRecord{
		ConfusionUuid: uuid,
		FilePath:      in.FilePath,
		FileType:      fileType,
		LineNumbers:   in.LineNumbers,
		CodeSnippet:   in.CodeSnippet,
		WhyModified:   in.WhyModified,
	}
}

// AiConfusionJSONInput defines the ingestion payload when adding episodes via JSON files.
type AiConfusionJSONInput struct {
	ConfusionUuid     string                 `json:"confusionUuid,omitempty"`
	AgentSlug         string                 `json:"agentSlug,omitempty"`
	Agent             string                 `json:"agent,omitempty"`
	FileNames         string                 `json:"fileNames,omitempty"`
	FilePath          string                 `json:"filePath,omitempty"`
	MdFiles           string                 `json:"mdFiles,omitempty"`
	Spec              string                 `json:"spec,omitempty"`
	SpecPath          string                 `json:"specPath,omitempty"`
	Task              string                 `json:"task,omitempty"`
	TaskSlug          string                 `json:"taskSlug,omitempty"`
	Subtask           string                 `json:"subtask,omitempty"`
	SubtaskSlug       string                 `json:"subtaskSlug,omitempty"`
	HeatmapInfo       string                 `json:"heatmapInfo,omitempty"`
	LineNumbers       string                 `json:"lineNumbers,omitempty"`
	CodeSnippet       string                 `json:"codeSnippet,omitempty"`
	Reason            string                 `json:"reason,omitempty"`
	InitialIntention  string                 `json:"initialIntention,omitempty"`
	AiIntention       string                 `json:"aiIntention,omitempty"`
	ReasoningThinking string                 `json:"reasoningThinking,omitempty"`
	AiThinking        string                 `json:"aiThinking,omitempty"`
	ConfusionQuestion string                 `json:"confusionQuestion,omitempty"`
	Question          string                 `json:"question,omitempty"`
	IsCompleted       bool                   `json:"isCompleted,omitempty"`
	ReferencedFiles   []AiConfusionFileInput `json:"referencedFiles,omitempty"`
}

// ToRecord converts an ingested JSON input to a normalized AiConfusionRecord.
func (in AiConfusionJSONInput) ToRecord() AiConfusionRecord {
	rec := in.populateBaseRecord()
	in.populateNarrativeFields(&rec)
	rec.ReferencedFiles = in.populateReferencedFiles(rec.ConfusionUuid)

	return rec
}

func (in AiConfusionJSONInput) populateReferencedFiles(uuid string) []AiConfusionFileRecord {
	hasFiles := len(in.ReferencedFiles) > 0
	if !hasFiles {
		return nil
	}

	files := make([]AiConfusionFileRecord, 0, len(in.ReferencedFiles))
	for _, f := range in.ReferencedFiles {
		files = append(files, f.ToRecord(uuid))
	}

	return files
}

func (in AiConfusionJSONInput) populateBaseRecord() AiConfusionRecord {
	rec := in.populateIdentityRecord()
	in.populateContextRecord(&rec)

	return rec
}

func (in AiConfusionJSONInput) populateIdentityRecord() AiConfusionRecord {
	return AiConfusionRecord{
		ConfusionUuid: in.ConfusionUuid,
		AgentSlug:     resolveFirstNonEmpty(in.AgentSlug, in.Agent),
		FileNames:     resolveFirstNonEmpty(in.FileNames, in.FilePath),
		MdFiles:       in.MdFiles,
		Spec:          resolveFirstNonEmpty(in.Spec, in.SpecPath),
		Task:          resolveFirstNonEmpty(in.Task, in.TaskSlug),
		Subtask:       resolveFirstNonEmpty(in.Subtask, in.SubtaskSlug),
	}
}

func (in AiConfusionJSONInput) populateContextRecord(rec *AiConfusionRecord) {
	rec.HeatmapInfo = in.HeatmapInfo
	rec.LineNumbers = in.LineNumbers
	rec.CodeSnippet = in.CodeSnippet
	rec.Reason = in.Reason
	rec.IsCompleted = in.IsCompleted
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
