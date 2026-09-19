package cmdautomation

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RelPathOptions configures relative path auditing and sanitization.
type RelPathOptions struct {
	Dir        string   `json:"dir"`
	IsFixMode  bool     `json:"isFixMode"`
	Extensions []string `json:"extensions"`
	IsJson     bool     `json:"isJson"`
}

// PathViolation records an absolute path or file URI found in a file.
type PathViolation struct {
	File       string `json:"file"`
	LineNumber int    `json:"lineNumber"`
	RawPath    string `json:"rawPath"`
}

// RelPathResult aggregates findings from relative path audits.
type RelPathResult struct {
	ScannedFiles int             `json:"scannedFiles"`
	Violations   []PathViolation `json:"violations"`
	FixedCount   int             `json:"fixedCount"`
	Duration     time.Duration   `json:"duration"`
	IsClean      bool            `json:"isClean"`
}

// NamingOptions configures boolean and naming convention audits.
type NamingOptions struct {
	Dir        string   `json:"dir"`
	Extensions []string `json:"extensions"`
	IsJson     bool     `json:"isJson"`
}

// NamingViolation records a naming or boolean comparison issue.
type NamingViolation struct {
	File        string `json:"file"`
	LineNumber  int    `json:"lineNumber"`
	Kind        string `json:"kind"`
	LineContent string `json:"lineContent"`
}

// NamingResult aggregates naming audit findings.
type NamingResult struct {
	ScannedFiles int               `json:"scannedFiles"`
	Violations   []NamingViolation `json:"violations"`
	Duration     time.Duration     `json:"duration"`
	IsClean      bool              `json:"isClean"`
}

// RuleAuditOptions configures specialized rule checkers (result-wrapper, params, enums).
type RuleAuditOptions struct {
	Dir      string `json:"dir"`
	RuleName string `json:"ruleName"`
	IsJson   bool   `json:"isJson"`
}

// RuleViolation represents a single guideline infraction.
type RuleViolation struct {
	File        string `json:"file"`
	LineNumber  int    `json:"lineNumber"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

// RuleAuditResult stores aggregated rule audit results.
type RuleAuditResult struct {
	ScannedFiles int             `json:"scannedFiles"`
	Violations   []RuleViolation `json:"violations"`
	Duration     time.Duration   `json:"duration"`
	IsClean      bool            `json:"isClean"`
}

type (
	// RelPathResultMonad wraps RelPathResult with AppError.
	RelPathResultMonad = result.Result[RelPathResult]

	// NamingResultMonad wraps NamingResult with AppError.
	NamingResultMonad = result.Result[NamingResult]

	// RuleAuditResultMonad wraps RuleAuditResult with AppError.
	RuleAuditResultMonad = result.Result[RuleAuditResult]
)
