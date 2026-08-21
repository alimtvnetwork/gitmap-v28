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

type llmDocsOptions struct {
	toStdout bool
	format   string
	sections string
}

func parseLLMDocsFlags(args []string) (llmDocsOptions, error) {
	fs := flag.NewFlagSet(constants.CmdLLMDocs, flag.ExitOnError)
	toStdout := fs.Bool(constants.FlagLLMDocsStdout, false, constants.FlagDescLLMDocsStdout)
	format := fs.String(constants.FlagLLMDocsFormat, constants.FormatMarkdown, constants.FlagDescLLMDocsFormat)
	sections := fs.String(constants.FlagLLMDocsSections, "", constants.FlagDescLLMDocsSections)
	if err := fs.Parse(reorderFlagsBeforeArgs(args)); err != nil {
		return llmDocsOptions{}, err
	}
	if *format != constants.FormatMarkdown && *format != constants.FormatJSON {
		return llmDocsOptions{}, fmt.Errorf(constants.ErrLLMDocsFormat, *format)
	}
	return llmDocsOptions{toStdout: *toStdout, format: *format, sections: *sections}, nil
}

func llmDocsExt(format string) string {
	if format == constants.FormatJSON {
		return constants.ExtJSON
	}

	return constants.ExtMD
}

func writeLLMDocsFile(content, format string) {
	fmt.Print(constants.MsgLLMDocsGenning)
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrLLMDocsWrite, err)
		os.Exit(1)
	}
	outPath := filepath.Join(wd, "LLM"+llmDocsExt(format))
	if writeErr := os.WriteFile(outPath, []byte(content), constants.FilePermission); writeErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrLLMDocsWrite, writeErr)
		os.Exit(1)
	}
	fmt.Printf(constants.MsgLLMDocsWritten, outPath)
}

// runLLMDocs generates LLM.md or prints to stdout with --stdout.
func runLLMDocs(args []string) {
	checkHelp(constants.CmdLLMDocs, args)
	opts, err := parseLLMDocsFlags(args)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	sectionSet := parseSections(opts.sections)
	content := buildLLMOutput(opts.format, sectionSet)
	if opts.toStdout {
		fmt.Print(content)

		return
	}
	writeLLMDocsFile(content, opts.format)
}

func collectValidSections() map[string]bool {
	valid := make(map[string]bool)
	for _, s := range strings.Split(constants.LLMDocsValidSections, ",") {
		valid[s] = true
	}

	return valid
}

// parseSections converts the comma-separated --sections value into a set.
// An empty string means all sections are included.
func parseSections(raw string) map[string]bool {
	if raw == "" {
		return nil
	}
	valid := collectValidSections()
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
	if format == constants.FormatJSON {
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

func appendLLMSectionsFirstHalf(sb *strings.Builder, sections map[string]bool) {
	if wantSection(sections, llmDocsKeyArchitecture) {
		writeLLMArchitecture(sb)
	}
	if wantSection(sections, llmDocsKeyCommands) {
		writeLLMCommands(sb)
	}
	if wantSection(sections, llmDocsKeyFlags) {
		writeLLMGlobalFlags(sb)
	}
	if wantSection(sections, llmDocsKeyConventions) {
		writeLLMCodingConventions(sb)
	}
}

func appendLLMSectionsSecondHalf(sb *strings.Builder, sections map[string]bool) {
	if wantSection(sections, llmDocsKeyStructure) {
		writeLLMProjectStructure(sb)
	}
	if wantSection(sections, llmDocsKeyDatabase) {
		writeLLMDatabase(sb)
	}
	if wantSection(sections, llmDocsKeyInstallation) {
		writeLLMInstallation(sb)
	}
	if wantSection(sections, llmDocsKeyPatterns) {
		writeLLMPatterns(sb)
	}
}

// buildLLMDocument assembles the complete LLM.md content dynamically.
func buildLLMDocument(sections map[string]bool) string {
	var sb strings.Builder
	writeLLMHeader(&sb)
	appendLLMSectionsFirstHalf(&sb, sections)
	appendLLMSectionsSecondHalf(&sb, sections)

	return sb.String()
}
