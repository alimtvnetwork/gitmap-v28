package heatmap

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RcLogResponseEnvelope represents the full structured response for heatmap log and update operations.
type RcLogResponseEnvelope struct {
	Status        string             `json:"status"`
	ActionId      int64              `json:"action_id"`
	WhereSaved    RcLogWhereSaved    `json:"where_saved"`
	DbInfo        RcLogDbInfo        `json:"db_info"`
	WhatSaved     RcLogWhatSaved     `json:"what_saved"`
	AiExpectation RcLogAiExpectation `json:"ai_expectation"`
}

// RcLogWhereSaved details the physical and relative persistence location.
type RcLogWhereSaved struct {
	DbPath         string `json:"db_path"`
	RelativeDbPath string `json:"relative_db_path"`
	RepoSecretsUrl string `json:"repo_secrets_url"`
}

// RcLogDbInfo holds database engine, table name, and namespace metadata.
type RcLogDbInfo struct {
	Driver   string `json:"driver"`
	Table    string `json:"table"`
	AiName   string `json:"ai_name"`
	RepoSlug string `json:"repo_slug"`
	Path     string `json:"path"`
}

// RcLogWhatSaved captures the exact action payload persisted to the database.
type RcLogWhatSaved struct {
	ActionId       int64  `json:"action_id"`
	ParentLogId    *int64 `json:"parent_log_id,omitempty"`
	TaskSlug       string `json:"task_slug"`
	RepoName       string `json:"repo_name"`
	ActionType     string `json:"action_type"`
	FilePath       string `json:"file_path"`
	StartLine      int    `json:"start_line"`
	EndLine        int    `json:"end_line"`
	CodeSnippet    string `json:"code_snippet,omitempty"`
	CommandText    string `json:"command_text,omitempty"`
	ReasoningState string `json:"reasoning_state"`
	ReasoningNotes string `json:"reasoning_notes,omitempty"`
	ExecutionDepth int    `json:"execution_depth"`
	ExecutionOrder int    `json:"execution_order"`
}

// RcLogAiExpectation tracks expectations, outcome, and satisfaction status.
type RcLogAiExpectation struct {
	ExpectedOutcome  string `json:"expected_outcome"`
	ActualOutcome    string `json:"actual_outcome"`
	IsExpectationMet bool   `json:"is_expectation_met"`
	Status           string `json:"status"`
}

// BuildRcLogEnvelope constructs the structured envelope from an action record.
func BuildRcLogEnvelope(rec store.HeatmapActionLogRecord, dbPath, aiName, repoSlug string) RcLogResponseEnvelope {
	relPath := computeRelativeDbPath(dbPath)
	secretsUrl := ResolveCompanionRepoSecretsUrl()
	expStatus := computeExpectationStatus(rec.ExpectedOutcome, rec.ActualOutcome, rec.IsExpectationMet)

	return RcLogResponseEnvelope{
		Status:   "success",
		ActionId: rec.HeatmapActionLogId,
		WhereSaved: RcLogWhereSaved{
			DbPath:         dbPath,
			RelativeDbPath: relPath,
			RepoSecretsUrl: secretsUrl,
		},
		DbInfo: RcLogDbInfo{
			Driver:   "sqlite",
			Table:    "HeatmapActionLog",
			AiName:   aiName,
			RepoSlug: repoSlug,
			Path:     dbPath,
		},
		WhatSaved: RcLogWhatSaved{
			ActionId:       rec.HeatmapActionLogId,
			ParentLogId:    rec.ParentLogId,
			TaskSlug:       rec.TaskSlug,
			RepoName:       repoSlug,
			ActionType:     rec.ActionType,
			FilePath:       rec.FilePath,
			StartLine:      rec.StartLine,
			EndLine:        rec.EndLine,
			CodeSnippet:    rec.CodeSnippet,
			CommandText:    rec.CommandText,
			ReasoningState: rec.ReasoningState,
			ReasoningNotes: rec.ReasoningNotes,
			ExecutionDepth: rec.ExecutionDepth,
			ExecutionOrder: rec.ExecutionOrder,
		},
		AiExpectation: RcLogAiExpectation{
			ExpectedOutcome:  rec.ExpectedOutcome,
			ActualOutcome:    rec.ActualOutcome,
			IsExpectationMet: rec.IsExpectationMet,
			Status:           expStatus,
		},
	}
}

// BuildHeatmapLogEnvelope constructs the structured envelope for local heatmap records.
func BuildHeatmapLogEnvelope(rec store.HeatmapActionLogRecord, dbPath string) RcLogResponseEnvelope {
	aiName := store.DetectCurrentAIName()
	repoSlug := rec.RepoName
	if repoSlug == "" {
		repoSlug = store.ResolveAmbientRepoSlug("")
	}

	return BuildRcLogEnvelope(rec, dbPath, aiName, repoSlug)
}

// ResolveCompanionRepoSecretsUrl returns the remote origin URL for repo-secrets.
func ResolveCompanionRepoSecretsUrl() string {
	db, err := store.OpenSpecialReposSplitDB()
	if err == nil {
		defer db.Close()
		if rec, getErr := db.GetSpecialRepo("rs"); getErr == nil && rec.RemoteURL != "" {
			return rec.RemoteURL
		}
	}

	for _, cand := range []string{`D:\work\repo-secrets`, `repo-secrets`, `../repo-secrets`} {
		gitDir := filepath.Join(cand, ".git")
		if info, statErr := os.Stat(gitDir); statErr == nil && info.IsDir() {
			if u := readGitRemoteUrl(cand); u != "" {
				return u
			}
		}
	}

	return "https://github.com/alimtvnetwork/repo-secrets.git"
}

func readGitRemoteUrl(dir string) string {
	cmd := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func computeRelativeDbPath(absPath string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return filepath.ToSlash(absPath)
	}

	rel, err := filepath.Rel(cwd, absPath)
	if err != nil {
		return filepath.ToSlash(absPath)
	}

	return filepath.ToSlash(rel)
}

func computeExpectationStatus(expected, actual string, isMet bool) string {
	if !isMet {
		return "FAILED"
	}

	if actual != "" {
		return "MET"
	}

	if expected != "" {
		return "PENDING"
	}

	return "MET"
}
