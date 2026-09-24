// Package cmd — llmdocs.go generates a consolidated LLM.md reference file.
package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
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
		cliexit.HandleError(apperror.WrapSimple(err, constants.ErrLLMDocsWrite), 1)
	}

	outPath := filepath.Join(wd, "LLM"+llmDocsExt(format))
	if writeErr := os.WriteFile(outPath, []byte(content), constants.FilePermission); writeErr != nil {
		cliexit.HandleError(apperror.WrapSimple(writeErr, constants.ErrLLMDocsWrite), 1)
	}

	fmt.Printf(constants.MsgLLMDocsWritten, outPath)
}

// runLLMDocs generates LLM.md or prints to stdout with --stdout.
func runLLMDocs(args []string) error {
	checkHelp(constants.CmdLLMDocs, args)
	opts, err := parseLLMDocsFlags(args)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())

		return apperror.WrapSimple(err, "parse llmdocs flags")
	}

	sectionSet := parseSections(opts.sections)
	content := buildLLMOutput(opts.format, sectionSet)
	if opts.toStdout {
		fmt.Print(content)

		return nil
	}

	writeLLMDocsFile(content, opts.format)

	return nil
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
			var empty map[string]bool

			return empty
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
	writeLLMTrainMultiStepCurriculum(&sb)

	return sb.String()
}

func writeLLMTrainMultiStepCurriculum(sb *strings.Builder) {
	sb.WriteString("\n## Continuous Multi-Step AI Loop & Fleet Orchestration Curriculum (`gitmap llm-train`)\n\n")
	sb.WriteString("Follow this 6-Stage Continuous Sequencing Loop during autonomous multi-node and repository execution:\n\n")
	sb.WriteString("### Loop Stage 0: Pre-Flight Workspace Sync & In-Memory AI Telemetry Ping\n")
	sb.WriteString("1. Always execute `git pull --rebase origin main` before inspecting or modifying files.\n")
	sb.WriteString("2. Probe the **In-Memory AI Server** (`gitmap ai-server`) bound on the first free port in `[47831, 47832, 47833, 47834]`:\n")
	sb.WriteString("   - `GET /api/v1/ai/ping` — Instant zero-disk JSON liveness response.\n")
	sb.WriteString("   - `GET /api/v1/ai/status` — Active version, bound port, and `DH2D` hot-cache count.\n")
	sb.WriteString("   - `GET /api/v1/ai/search-cache?q=<token>` — Sub-40µs RAM lookup of frequent repository queries.\n\n")
	sb.WriteString("### Loop Stage 1: Native AUM Search (`DH2D` SQLite History & Hot-Cache Optimization)\n")
	sb.WriteString("- Never use slow unindexed pipeline searches (`Get-ChildItem -Recurse | Select-String` takes ~14.85s; Python grep takes ~33.20s).\n")
	sb.WriteString("- Use `gitmap search \"<query>\"` (`0.82 ms` cold, `0.04 ms` hot-cached). Every search generates a deterministic **`DH2D-<HEX>`** SQL ID in `SearchSplitDB` (`SearchHotCache`) and automatically promotes queries with `HitCount >= 2` into RAM.\n")
	sb.WriteString("- Inspect top searches and `DH2D` IDs via `gitmap search history` or `gitmap aum search-history`.\n\n")
	sb.WriteString("### Loop Stage 2: SSH Nodes Export/Import JSON, One-Liner & Node-Config Fleet Deployment\n")
	sb.WriteString("- `gitmap ssh nodes export-json [path]` — Exports all SSH nodes to `gitmap-ssh-nodes.json` (and `gitmap-ssh.json`) or a custom file/directory.\n")
	sb.WriteString("- `gitmap ssh nodes import-json [path]` — Imports SSH nodes from `gitmap-ssh-nodes.json` (or `gitmap-ssh.json`) or `--base64 <payload>`.\n")
	sb.WriteString("- `gitmap ssh export-oneliner` — Generates a single-line copy-paste command (`gitmap ssh nodes import-json --base64 \"...\"`) to onboard any machine instantly.\n")
	sb.WriteString("- `gitmap ssh deploy node-config` (or `gitmap ssh deploy nc`) `--except id,ip,alias` — Deploys all nodes & worker IDs across remote SSH machines in parallel.\n\n")
	sb.WriteString("### Loop Stage 3: Bidirectional SSH & AGY Fleet Execution\n")
	sb.WriteString("- Both `gitmap ssh agy <subcmd>` and `gitmap agy ssh <subcmd>` execute any AGY or GitMap command across all SSH nodes concurrently.\n")
	sb.WriteString("- Exclude specific nodes by numeric ID (`1`), worker ID (`worker-1`), IP address, or alias using `--except`, `--excep`, `--accept`, or `-e`.\n\n")
	sb.WriteString("### Loop Stage 4: AGY Rerun (`[1|2|3|4|all|queue]`) with Image Re-Injection & Completion Prefix\n")
	sb.WriteString("- `gitmap agy rerun 1` — Restarts Antigravity IDE for Project #1 and replays its active prompt with all attached image file paths (`media` / `.png` / `.jpg`).\n")
	sb.WriteString("- `gitmap agy rerun all` — Restarts and replays active prompts across all running projects (`1..4`).\n")
	sb.WriteString("- `gitmap agy rerun queue` — Re-injects the active prompt plus up to 5 queued prompts prefixed with `\"is it completed properly and released??\"`.\n\n")
	sb.WriteString("### Loop Stage 5: AGM (Antigravity-Manager) Fleet SSH Updates\n")
	sb.WriteString("- `gitmap agm update ssh --except id,ip,alias`\n")
	sb.WriteString("- `gitmap agm update-all-nodes --except id,ip,alias`\n")
	sb.WriteString("- `gitmap ssh update agm --except id,ip,alias`\n")
}
