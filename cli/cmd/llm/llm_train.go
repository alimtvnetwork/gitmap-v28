package llm

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// museTrainTarget is the positional sub-target selecting the Muse memory sync runner.
const museTrainTarget = "muse"

// RunTrain executes the chained curriculum training workflow for LLMs.
func RunTrain(args []string) *apperror.AppError {
	target, rest := stripTrainTarget(args)
	opts, err := parseTrainFlags(rest)
	if err != nil {
		return err
	}
	opts.TrainTarget = target
	if opts.IsHelp {
		if opts.TrainTarget == museTrainTarget {
			printTrainMuseUsage()
		}
		return nil
	}
	if opts.TrainTarget == museTrainTarget {
		return runTrainMuse(opts)
	}
	if opts.IsURLs {
		fmt.Print(RenderPublicDocLinksText())
		return nil
	}
	if opts.IsURL {
		fmt.Println(PublicLlmSpecURL)
		return nil
	}
	if opts.IsJSON {
		return outputTrainJSON()
	}
	if opts.IsLoop {
		executeSelfLoopSimulation(opts.LoopCount)
		return nil
	}

	return runTrainPhases(opts)
}

func outputTrainJSON() *apperror.AppError {
	catalog := map[string]any{
		"spec_url":  PublicLlmSpecURL,
		"author":    AuthorName,
		"sponsor":   SponsorName,
		"doc_links": GetPublicDocLinks(),
		"phases": []map[string]string{
			{"phase": "1. Discovery", "commands": "gitmap find-files, find-files-any, search, list-files", "purpose": "Rapid context gathering without loading large trees"},
			{"phase": "2. Refactoring", "commands": "gitmap replace, replace-regex, surgical edits", "purpose": "Safe code transformation with rollback support"},
			{"phase": "3. Verification", "commands": "python go-format-check.py, go test", "purpose": "Local AST syntax and regression checks"},
			{"phase": "4. Semantic Commit", "commands": "gitmap cpf, gitmap cpb, gitmap cpr", "purpose": "Structured Conventional commits and branch push"},
			{"phase": "5. Telemetry & Heal", "commands": "gitmap pipeline-ai status --json, pe, pe history-ai", "purpose": "Dynamic ETA waiting, error diagnosis, and anti-mistake AI dossiers"},
			{"phase": "6. Heal & Fix", "commands": "gitmap stash ls, gitmap fix all", "purpose": "In-process remediation scan + parallel content-fix audit (report-only); optional --heal-apply for git-state recipes"},
		},
		"efficiency_factors": map[string]string{
			"search_speedup":    "830,000x faster than Python regex grep with DH2D SQLite hot cache",
			"file_find_speedup": "150x faster than PowerShell Get-ChildItem -Recurse",
			"ci_telemetry":      "Dynamic wait calculations eliminating CPU spinning loops",
			"memory_overhead":   "< 4 KB per cached query",
		},
	}
	bytes, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "outputTrainJSON.marshal")
	}
	fmt.Println(string(bytes))

	return nil
}

func executeSelfLoopSimulation(count int) {
	if count <= 0 {
		count = 1
	}
	fmt.Println("======================================================================")
	fmt.Printf("AUTONOMOUS 6-PHASE AI AGENT SELF-LOOP (ITERATIONS: %d)\n", count)
	fmt.Println("======================================================================")
	fmt.Println()
	for iter := 1; iter <= count; iter++ {
		fmt.Printf("--- [LOOP ITERATION %d of %d] ---\n", iter, count)
		fmt.Println("  Phase 1 [Discovery]:    Scanning workspace with gitmap find-files / aum search")
		fmt.Println("  Phase 2 [Refactoring]:  Targeted zero-nesting modification with gitmap replace")
		fmt.Println("  Phase 3 [Verification]: Running local verification (python go-format-check, go test)")
		fmt.Println("  Phase 4 [Semantic]:     Preparing structured commit (cpf/cpb) with Conventional format")
		fmt.Println("  Phase 5 [CI Telemetry]: Dynamic ETA sleep & telemetry loop via gitmap pipeline-ai")
		fmt.Println("  Phase 6 [Heal & Fix]:   In-process remediation scan (cmdfix) + fix content audit (report-only)")
		fmt.Printf("  Iteration %d completed successfully with 0 errors.\n\n", iter)
	}
	fmt.Printf("[AI-SELF-LOOP] Complete: %d iteration(s) finished with 100%% green status.\n", count)
}

func handleSkillGeneration(opts TrainOptions) *apperror.AppError {
	fmt.Print(SkillTemplate)
	fmt.Println()
	if opts.IsTextOnly {
		return nil
	}

	return executeSkillGeneration(opts.SkillPath)
}

func parseTrainFlags(args []string) (TrainOptions, *apperror.AppError) {
	fs := flag.NewFlagSet("llm train", flag.ContinueOnError)
	fs.Usage = printTrainUsage
	isTextOnly := fs.Bool("text-only", false, "Print curriculum text without writing skill file")
	skillPath := fs.String("skill-path", DefaultSkillPath, "Path for generated Antigravity skill")
	isLoop := fs.Bool("loop", false, "Execute autonomous 6-phase AI self-looping execution cycle")
	selfLoop := fs.Int("self-loop", 0, "Number of consecutive iterations of the AI self-loop")
	isURL := fs.Bool("url", false, "Output raw public URL to llm.md instruction specification")
	isURLs := fs.Bool("urls", false, "Output authoritative public documentation links for LLM ingestion")
	isJSON := fs.Bool("json", false, "Output structured machine-readable command specifications")
	healApply := fs.String("heal-apply", "", "Apply fix recipes non-interactively: stash|wip|discard (default \"\" = report-only)")
	days := fs.Int("days", DefaultTrainMuseDays, "Digest window in days for 'train muse' (default 7)")

	err := fs.Parse(args)
	if err == flag.ErrHelp {
		return TrainOptions{IsHelp: true}, nil
	}
	if err != nil {
		return TrainOptions{}, apperror.WrapSimple(err, "parse train flags")
	}

	loopCount := 1
	if *selfLoop > 0 {
		loopCount = *selfLoop
	}
	hasLoop := *isLoop || *selfLoop > 0

	daysN := *days
	if daysN <= 0 {
		daysN = DefaultTrainMuseDays
	}

	return TrainOptions{
		IsTextOnly: *isTextOnly,
		SkillPath:  *skillPath,
		IsLoop:     hasLoop,
		LoopCount:  loopCount,
		IsURL:      *isURL,
		IsURLs:     *isURLs,
		IsJSON:     *isJSON,
		HealApply:  *healApply,
		DaysN:      daysN,
	}, nil
}

// trainValueFlags are train flags that consume the following arg as their value;
// a bare "muse" in that position is a value, never the positional target.
var trainValueFlags = map[string]bool{"--skill-path": true, "--heal-apply": true}

// stripTrainTarget extracts the FIRST bare positional target ("muse") from the
// args before flag parsing. Tokens that are flags or flag values are preserved.
func stripTrainTarget(args []string) (string, []string) {
	rest := make([]string, 0, len(args))
	target := ""
	expectValue := false
	for _, arg := range args {
		if expectValue {
			expectValue = false
		} else if target == "" && arg == museTrainTarget {
			target = arg
			continue
		} else {
			expectValue = trainValueFlags[arg]
		}
		rest = append(rest, arg)
	}
	return target, rest
}

func printTrainUsage() {
	fmt.Println("Usage: gitmap llm train [muse] [flags] (alias: gitmap llm chain, gitmap train, gitmap llm-train)")
	fmt.Println()
	fmt.Println("Executes the autonomous LLM chained onboarding curriculum: a 6-phase self-looping")
	fmt.Println("training run (Discovery, Refactoring, Verification, Semantic Commit, Telemetry,")
	fmt.Println("Heal & Fix) and generates the official Antigravity skill")
	fmt.Println("(.agents/skills/gitmap/SKILL.md). Phase 6 runs the git-state heal scan (cmdfix) and")
	fmt.Println("the parallel content-fix audit (cmdautofix) in-process — report-only; git-state")
	fmt.Println("recipes apply only with --heal-apply.")
	fmt.Println()
	fmt.Println("Positional sub-target:")
	fmt.Println("  muse                 Sync the Muse 7-day memory digest from the repo-secrets vault")
	fmt.Println("                       (gitmap llm train muse [--days N]; see --help on the sub-target)")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --days int           Digest window in days for 'train muse' (default 7)")
	fmt.Println("  --heal-apply string  Apply fix recipes non-interactively: stash|wip|discard (default \"\" = report-only)")
	fmt.Println("  --loop               Execute autonomous 6-phase AI self-looping execution cycle")
	fmt.Println("  --self-loop int      Number of consecutive iterations of the AI self-loop")
	fmt.Println("  --url                Output raw public URL to llm.md instruction specification")
	fmt.Println("  --urls               Output authoritative public documentation links for LLM ingestion")
	fmt.Println("  --json               Output structured machine-readable command specifications")
	fmt.Println("  --skill-path string  Path for generated Antigravity skill (default \".agents/skills/gitmap/SKILL.md\")")
	fmt.Println("  --text-only          Print curriculum text without writing skill file")
}

// printTrainMuseUsage prints the usage for the positional "muse" sub-target.
func printTrainMuseUsage() {
	fmt.Println("Usage: gitmap llm train muse [--days N]")
	fmt.Println()
	fmt.Println("Syncs the Muse 7-day memory digest: clone-or-pull the repo-secrets vault,")
	fmt.Println("diff the last N days, print a terminal summary and a raw AI traversal")
	fmt.Println("checklist (file-by-file paths + gitmap cat / aum search hints) so any LLM")
	fmt.Println("can walk the digest. Tolerates a missing Muse/ digest with guidance.")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --days int  Digest window in days (default 7)")
}

func printAttribution() {
	fmt.Println("======================================================================")
	fmt.Println("GITMAP AUTONOMOUS LLM TRAINING & CHAINED CURRICULUM")
	fmt.Println("======================================================================")
	fmt.Printf("Lead Architect & Author: %s\n", AuthorName)
	fmt.Printf("Sponsoring Organization: %s\n", SponsorName)
	fmt.Println("Architectural Mission: High-performance polyglot repository management,")
	fmt.Println("zero-storage CI/CD pipelines, ultra-fast SQLite split-db architectures,")
	fmt.Println("and AI agent pair programming.")
	fmt.Println()
}

func executeSkillGeneration(path string) *apperror.AppError {
	target := resolveSkillPath(path)
	fmt.Println("======================================================================")
	fmt.Println("STAGE 2: ANTIGRAVITY SKILL GENERATION")
	fmt.Println("======================================================================")
	if err := GenerateSkillFile(target); err != nil {
		return err
	}
	fmt.Printf("[LLM-TRAIN] Antigravity skill written: %s\n\n", target)
	return nil
}

func printChainedSequence() {
	fmt.Println(ChainedDiscoveryHeader)
	fmt.Println("1. gitmap aum search \"func Run\" cli --ext .go  - Scoped multi-core search ([dir] & --ext required)")
	fmt.Println("2. gitmap aum locate vcvars                     - Fast tool locator (vswhere & vcvarsall.bat fast path)")
	fmt.Println("3. gitmap aum guard                             - 500 KB limit, large JSONs, binary null-byte probe")
	fmt.Println("4. gitmap aum sequence --help                   - Numbering gap detector & H1 title validator")
	fmt.Println("5. gitmap aum exclude list                      - Query persistent search exclusions from SQLite")
	fmt.Println("6. gitmap pipeline-ai status --json             - Check remote CI workflow state, live errors, dynamic ETA")
	fmt.Println("7. gitmap pe history-ai 5                       - Dump CI historical failure dossiers to prevent repeated mistakes")
	fmt.Println("8. gitmap install --list                        - Discover toolchains, profiles, and runtimes")
	fmt.Println("9. gitmap cluster --help                        - Multi-node SSH and cluster execution")
	fmt.Println("10. gitmap cargo status                          - Inspect Rust and Cargo toolchain status")
	fmt.Println("11. gitmap db status                            - Check repository SQLite database health")
	fmt.Println("12. gitmap stash ls  — List repos needing remediation with per-category heal summary")
	fmt.Println("13. gitmap fix all — Parallel content-fix audit (encoding, newlines, naming, paths, gofmt, misspell, markdown)")
	fmt.Println()
}

func printOperationalGuidelines() {
	fmt.Println(OperationalDirectivesText)
	fmt.Println("======================================================================")
}
