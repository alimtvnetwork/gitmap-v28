// Package cmd — llmdocs.go generates a consolidated LLM.md reference file.
package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runLLMDocs generates LLM.md or prints to stdout with --stdout.
func runLLMDocs(args []string) {
	checkHelp("llm-docs", args)

	fs := flag.NewFlagSet("llm-docs", flag.ExitOnError)
	toStdout := fs.Bool(constants.FlagLLMDocsStdout, false, constants.FlagDescLLMDocsStdout)
	format := fs.String(constants.FlagLLMDocsFormat, "markdown", constants.FlagDescLLMDocsFormat)
	sections := fs.String(constants.FlagLLMDocsSections, "", constants.FlagDescLLMDocsSections)

	reordered := reorderFlagsBeforeArgs(args)

	if err := fs.Parse(reordered); err != nil {
		fmt.Fprintf(os.Stderr, "llm-docs: %v\n", err)
		os.Exit(1)
	}

	if *format != "markdown" && *format != "json" {
		fmt.Fprintf(os.Stderr, constants.ErrLLMDocsFormat, *format)
		os.Exit(1)
	}

	sectionSet := parseSections(*sections)

	content := buildLLMOutput(*format, sectionSet)

	if *toStdout {
		fmt.Print(content)

		return
	}

	fmt.Print(constants.MsgLLMDocsGenning)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrLLMDocsWrite, err)
		os.Exit(1)
	}

	ext := ".md"
	if *format == "json" {
		ext = ".json"
	}

	outPath := filepath.Join(wd, "LLM"+ext)

	if writeErr := os.WriteFile(outPath, []byte(content), 0o644); writeErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrLLMDocsWrite, writeErr)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgLLMDocsWritten, outPath)
}

// parseSections converts the comma-separated --sections value into a set.
// An empty string means all sections are included.
func parseSections(raw string) map[string]bool {
	if raw == "" {
		return nil
	}

	valid := make(map[string]bool)
	for _, s := range strings.Split(constants.LLMDocsValidSections, ",") {
		valid[s] = true
	}

	set := make(map[string]bool)

	for _, s := range strings.Split(raw, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		if !valid[s] {
			fmt.Fprintf(os.Stderr, constants.ErrLLMDocsSections, s)
			os.Exit(1)
		}

		set[s] = true
	}

	return set
}

// wantSection returns true if the section should be included.
func wantSection(set map[string]bool, name string) bool {
	if set == nil {
		return true
	}

	return set[name]
}

// buildLLMOutput returns the document in the requested format.
func buildLLMOutput(format string, sections map[string]bool) string {
	if format == "json" {
		return buildLLMJSON(sections)
	}

	return buildLLMDocument(sections)
}

// buildLLMJSON assembles a JSON representation of the LLM reference.
// Routed through stablejson for compile-time key-order guarantees.
func buildLLMJSON(sections map[string]bool) string {
	var buf bytes.Buffer
	if err := encodeLLMDocsJSON(&buf, sections); err != nil {
		return "{}\n"
	}

	return buf.String()
}

// buildLLMDocument assembles the complete LLM.md content dynamically.
func buildLLMDocument(sections map[string]bool) string {
	var sb strings.Builder

	writeLLMHeader(&sb)

	if wantSection(sections, "architecture") {
		writeLLMArchitecture(&sb)
	}

	if wantSection(sections, "commands") {
		writeLLMCommands(&sb)
	}

	if wantSection(sections, "flags") {
		writeLLMGlobalFlags(&sb)
	}

	if wantSection(sections, "conventions") {
		writeLLMCodingConventions(&sb)
	}

	if wantSection(sections, "structure") {
		writeLLMProjectStructure(&sb)
	}

	if wantSection(sections, "database") {
		writeLLMDatabase(&sb)
	}

	if wantSection(sections, "installation") {
		writeLLMInstallation(&sb)
	}

	if wantSection(sections, "patterns") {
		writeLLMPatterns(&sb)
	}

	return sb.String()
}
