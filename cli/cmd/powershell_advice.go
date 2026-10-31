package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// PowerShellTargetKind identifies the semantic classification of a CLI target.
type PowerShellTargetKind string

const (
	TargetKindDirectory     PowerShellTargetKind = "DIRECTORY"
	TargetKindDocument      PowerShellTargetKind = "DOCUMENT"
	TargetKindUnknownCmdlet PowerShellTargetKind = "UNKNOWN_CMDLET"
	TargetKindValidCommand  PowerShellTargetKind = "VALID_COMMAND"
)

// PowerShellAdviceOptions configures the diagnostic card rendering.
type PowerShellAdviceOptions struct {
	RawTarget    string
	TargetKind   PowerShellTargetKind
	ResolvedPath string
	ContextNote  string
}

// ReplacementOption represents a suggested alternative command.
type ReplacementOption struct {
	Command     string
	Description string
}

var documentExtensions = map[string]bool{
	".md": true, ".markdown": true, ".txt": true,
	".rtf": true, ".doc": true, ".docx": true,
	".json": true, ".jsonl": true, ".csv": true, ".tsv": true,
	".yml": true, ".yaml": true, ".toml": true, ".xml": true,
	".sql": true, ".graphql": true, ".html": true, ".css": true, ".log": true,
	".go": true, ".rs": true, ".py": true, ".ts": true,
	".tsx": true, ".js": true, ".jsx": true, ".c": true,
	".cpp": true, ".cs": true, ".java": true,
}

var classDescriptions = map[PowerShellTargetKind]string{
	TargetKindDirectory:     "DIRECTORY (Filesystem folder cannot be executed as script)",
	TargetKindDocument:      "DOCUMENT (Document file cannot be executed as script)",
	TargetKindUnknownCmdlet: "UNKNOWN_CMDLET (Term or path not recognized as valid cmdlet)",
}

func isDocumentFileExtension(ext string) bool {
	cleanExt := strings.ToLower(strings.TrimSpace(ext))
	if cleanExt == "" {
		return false
	}

	if !strings.HasPrefix(cleanExt, ".") {
		cleanExt = "." + cleanExt
	}

	return documentExtensions[cleanExt]
}

func isExecutableScriptExtension(ext string) bool {
	clean := strings.ToLower(strings.TrimSpace(ext))
	if clean == "" {
		return false
	}

	if !strings.HasPrefix(clean, ".") {
		clean = "." + clean
	}

	return clean == ".ps1" || clean == ".psm1" || clean == ".bat" || clean == ".cmd" || clean == ".exe" || clean == ".sh"
}

func isDirectoryTarget(token string) bool {
	clean := strings.Trim(token, "\"'")
	if clean == "" || strings.HasPrefix(clean, "-") {
		return false
	}

	fi, err := os.Stat(clean)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func isDocumentFileTarget(token string) bool {
	clean := strings.Trim(token, "\"'")
	if clean == "" || strings.HasPrefix(clean, "-") {
		return false
	}

	fi, err := os.Stat(clean)
	if err != nil || fi.IsDir() {
		return false
	}

	return isDocumentFileExtension(filepath.Ext(clean))
}

func resolveNonExecutableTargetKind(clean string) (PowerShellTargetKind, string) {
	if isDirectoryTarget(clean) {
		return TargetKindDirectory, clean
	}

	if isDocumentFileTarget(clean) {
		return TargetKindDocument, clean
	}

	return TargetKindValidCommand, clean
}

func detectPowerShellTarget(token string) (PowerShellTargetKind, string) {
	clean := strings.Trim(token, "\"'")
	if clean == "" || strings.HasPrefix(clean, "-") {
		return TargetKindValidCommand, ""
	}

	if isExecutableScriptExtension(filepath.Ext(clean)) {
		return TargetKindValidCommand, clean
	}

	return resolveNonExecutableTargetKind(clean)
}

func resolveClassDescription(kind PowerShellTargetKind) string {
	desc, hasDesc := classDescriptions[kind]
	if hasDesc {
		return desc
	}

	return "UNKNOWN (Invalid execution target)"
}

func renderCardTitle(title string) string {
	return termout.SafeTokyoStyle(termout.TokyoPrimary).Bold(true).Render(title)
}

func renderTargetText(target string) string {
	return termout.SafeTokyoStyle(termout.TokyoError).Bold(true).Render(target)
}

func renderClassText(desc string) string {
	return termout.SafeTokyoStyle(termout.TokyoWarn).Render(desc)
}

func renderContextNote(note string) string {
	if note == "" {
		return ""
	}

	return fmt.Sprintf("Note:    %s", termout.SafeTokyoStyle(termout.TokyoMuted).Render(note))
}

func resolveTargetString(opts PowerShellAdviceOptions) string {
	if opts.RawTarget != "" {
		return opts.RawTarget
	}

	return opts.ResolvedPath
}

func buildTargetHeader(opts PowerShellAdviceOptions) string {
	classDesc := resolveClassDescription(opts.TargetKind)
	target := resolveTargetString(opts)
	title := renderCardTitle("⚡ GITMAP AI GUARD: Invalid PowerShell Target Detected")
	lines := fmt.Sprintf("%s\n\nTarget:  %s\nClass:   %s", title, renderTargetText(target), renderClassText(classDesc))
	if opts.ContextNote != "" {
		return fmt.Sprintf("%s\n%s", lines, renderContextNote(opts.ContextNote))
	}

	return lines
}

func buildAdviceQuestions() string {
	header := termout.SafeTokyoStyle(termout.TokyoWarn).Bold(true).Render("🤔 AI Reflection Questions:")
	q1 := "1. Why did you run this? What was your expectation?"
	q2 := "2. Why did you run this command? Did you mistake 'gitmap ps' for directory listing or task inspection?"
	q3 := "3. What was your intended outcome (list files, read doc, run script)?"

	return fmt.Sprintf("%s\n%s\n%s\n%s", header, q1, q2, q3)
}

func buildActualBehaviorNote() string {
	header := termout.SafeTokyoStyle(termout.TokyoPrimary).Bold(true).Render("ℹ️ What 'gitmap ps' Actually Does:")
	b1 := "• Executes commands via PowerShell Core ('pwsh -NoProfile -Command ...')."
	b2 := "• Expects valid PowerShell expressions, cmdlets, or executable scripts."
	b3 := "• Does NOT navigate folders, display directory contents, or view files."
	b4 := "• In command mode vs expression mode, syntax parser triggers ParserError on paths like '02-spec' due to '-' operator."

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s", header, b1, b2, b3, b4)
}

func buildPowerShellDiagnosticCard(opts PowerShellAdviceOptions) string {
	header := buildTargetHeader(opts)
	questions := buildAdviceQuestions()
	actual := buildActualBehaviorNote()
	replacements := buildNativeReplacements(opts)
	learning := build3StepLearningActions()
	div := renderCardDivider()

	content := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
		header, div, questions, div, actual, div, replacements, div, learning,
	)

	return termout.TokyoBoxStyle(content)
}

func renderPowerShellDiagnosticCard(opts PowerShellAdviceOptions) {
	card := buildPowerShellDiagnosticCard(opts)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, card)
	fmt.Fprintln(os.Stderr)
}

func buildPreflightDetails(target string, kind PowerShellTargetKind) map[string]any {
	return map[string]any{
		"reported": true,
		"target":   target,
		"kind":     string(kind),
		"exitCode": int(cliexit.ExitCodeUsageError),
	}
}

func buildPreflightAbortError(target string, kind PowerShellTargetKind) error {
	msg := fmt.Sprintf("invalid powershell target %s: %s cannot be executed directly", target, strings.ToLower(string(kind)))
	ctx := buildPreflightDetails(target, kind)

	return apperror.NewWithDetails("powershell", "E_PS_PATH_INTERCEPTED", msg, "cmdpowershell", apperror.ErrorTypeAbort, apperror.SeverityWarn, ctx)
}

func reportInvalidPreflightTarget(candidate string, kind PowerShellTargetKind, resolved string) error {
	opts := PowerShellAdviceOptions{RawTarget: candidate, TargetKind: kind, ResolvedPath: resolved}
	renderPowerShellDiagnosticCard(opts)

	return buildPreflightAbortError(candidate, kind)
}

func inspectPowerShellPreflight(args []string) error {
	candidate := extractPreflightCandidate(args)
	if candidate == "" {
		return nil
	}

	kind, resolved := detectPowerShellTarget(candidate)
	if kind == TargetKindValidCommand {
		return nil
	}

	return reportInvalidPreflightTarget(candidate, kind, resolved)
}

func extractPreflightCandidate(args []string) string {
	if len(args) == 0 {
		return ""
	}

	if args[0] == "-Command" || args[0] == "-c" {
		return extractTokenAfterCommandFlag(args)
	}

	if strings.HasPrefix(args[0], "-") {
		return ""
	}

	return args[0]
}

func extractTokenAfterCommandFlag(args []string) string {
	if len(args) < 2 {
		return ""
	}

	return args[1]
}
