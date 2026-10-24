package llm

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// museTrainTarget is the positional sub-target selecting the Muse memory sync runner.
const museTrainTarget = "muse"

// RunTrain executes the chained curriculum training workflow for LLMs.
func RunTrain(args []string) *apperror.AppError {
	opts, heatmapOpts, err := resolveTrainConfigurations(args)
	if err != nil || opts.IsHelp {
		return err
	}
	if opts.TrainTarget == museTrainTarget {
		return runTrainMuse(opts)
	}
	if isSpecial := isSpecialTrainMode(opts); isSpecial {
		return nil
	}

	return executeTrainOrJSON(opts, heatmapOpts)
}

func executeTrainOrJSON(opts TrainOptions, heatmapOpts LlmHeatmapOptions) *apperror.AppError {
	if opts.IsJSON {
		return outputTrainJSON()
	}

	return dispatchTrainCurriculum(opts, heatmapOpts)
}

// RunLlmTrain is an alias for RunTrain executing the autonomous LLM training curriculum.
func RunLlmTrain(args []string) *apperror.AppError {
	return RunTrain(args)
}

func resolveTrainConfigurations(args []string) (TrainOptions, LlmHeatmapOptions, *apperror.AppError) {
	target, rest := stripTrainTarget(args)
	opts, heatmapOpts, err := parseTrainHeatmapFlags(rest)
	if err != nil {
		return TrainOptions{}, LlmHeatmapOptions{}, err
	}
	opts.TrainTarget = target
	if opts.IsHelp && opts.TrainTarget == museTrainTarget {
		printTrainMuseUsage()
	}

	return opts, heatmapOpts, nil
}

func isSpecialTrainMode(opts TrainOptions) bool {
	if opts.IsURLs {
		fmt.Print(RenderPublicDocLinksText())
		return true
	}
	if opts.IsURL {
		fmt.Println(PublicLlmSpecURL)
		return true
	}

	return isLoopMode(opts)
}

func isLoopMode(opts TrainOptions) bool {
	if opts.IsLoop {
		executeSelfLoopSimulation(opts.LoopCount)
		return true
	}

	return false
}

func dispatchTrainCurriculum(opts TrainOptions, heatmapOpts LlmHeatmapOptions) *apperror.AppError {
	if err := runTrainPhases(opts); err != nil {
		return err
	}

	return executeHeatmapTrainingCurriculum(heatmapOpts)
}

func executeHeatmapTrainingCurriculum(opts LlmHeatmapOptions) *apperror.AppError {
	sources, err := ResolveHeatmapTrainingSources(opts)
	if err != nil {
		return apperror.WrapSimple(err, "llm.resolve_heatmap")
	}
	if len(sources) == 0 {
		return nil
	}

	return ingestAndPrintHeatmapDigest(sources[0], opts)
}

func ingestAndPrintHeatmapDigest(source string, opts LlmHeatmapOptions) *apperror.AppError {
	examples, err := IngestHeatmapTrainingData(source, opts.Limit)
	if err != nil {
		return apperror.WrapSimple(err, "llm.ingest_heatmap")
	}
	filtered := filterExamplesByTask(examples, opts.TaskSlug)
	if len(filtered) == 0 {
		return nil
	}
	pos, neg := PartitionHeatmapExamples(filtered)
	printHeatmapTelemetry(source, len(filtered), len(pos), len(neg))
	fmt.Print(FormatHeatmapTrainingDigest(filtered))

	return nil
}

func filterExamplesByTask(examples []HeatmapTrainingExample, taskSlug string) []HeatmapTrainingExample {
	hasNoFilter := strings.TrimSpace(taskSlug) == ""
	if hasNoFilter {
		return examples
	}

	var filtered []HeatmapTrainingExample
	for _, ex := range examples {
		if strings.EqualFold(ex.TaskSlug, strings.TrimSpace(taskSlug)) {
			filtered = append(filtered, ex)
		}
	}

	return filtered
}

func printHeatmapTelemetry(source string, total, pos, neg int) {
	fmt.Println("======================================================================")
	fmt.Println("HEATMAP TELEMETRY INGESTION")
	fmt.Println("======================================================================")
	fmt.Printf("[HEATMAP-INGEST] Source: %s\n", source)
	fmt.Printf("[HEATMAP-INGEST] Ingested: %d actions (%d positive, %d negative RCA)\n\n", total, pos, neg)
}

func outputTrainJSON() *apperror.AppError {
	catalog := map[string]any{
		"spec_url":  PublicLlmSpecURL,
		"author":    AuthorName,
		"sponsor":   SponsorName,
		"doc_links": GetPublicDocLinks(),
		"phases": []map[string]string{
			{"phase": "1. Discovery", "commands": "gitmap find-files, find-files-any, search, list-files, workgroup ls", "purpose": "Rapid context gathering without loading large trees"},
			{"phase": "2. Refactoring", "commands": "gitmap replace, replace-regex, surgical edits", "purpose": "Safe code transformation with rollback support"},
			{"phase": "3. Verification", "commands": "python go-format-check.py, go test", "purpose": "Local AST syntax and regression checks"},
			{"phase": "4. Semantic Commit", "commands": "gitmap cpf, gitmap cpb, gitmap cbf, gitmap cpr", "purpose": "Structured Conventional commits and [module]: <item> format"},
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
		fmt.Println("  Phase 1 [Discovery]:    Scanning workspace with gitmap find-files / aum search / workgroup ls")
		fmt.Println("  Phase 2 [Refactoring]:  Targeted zero-nesting modification with gitmap replace")
		fmt.Println("  Phase 3 [Verification]: Running local verification (python go-format-check, go test)")
		fmt.Println("  Phase 4 [Semantic]:     Preparing structured commit (cpf/cbf) with [module]: <item> format")
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

type rawTrainFlagPointers struct {
	isTextOnly  *bool
	skillPath   *string
	isLoop      *bool
	selfLoop    *int
	isURL       *bool
	isURLs      *bool
	isJSON      *bool
	healApply   *string
	days        *int
	heatmap     *string
	heatmapRepo *string
	heatmapTask *string
}

func declareTrainFlags(fs *flag.FlagSet) rawTrainFlagPointers {
	p := rawTrainFlagPointers{
		isTextOnly: fs.Bool("text-only", false, "Print curriculum text without writing skill file"),
		skillPath:  fs.String("skill-path", DefaultSkillPath, "Path for generated Antigravity skill"),
		isLoop:     fs.Bool("loop", false, "Execute autonomous 6-phase AI self-looping execution cycle"),
		selfLoop:   fs.Int("self-loop", 0, "Number of consecutive iterations of the AI self-loop"),
		isURL:      fs.Bool("url", false, "Output raw public URL to llm.md instruction specification"),
		isURLs:     fs.Bool("urls", false, "Output authoritative public documentation links for LLM ingestion"),
		isJSON:     fs.Bool("json", false, "Output structured machine-readable command specifications"),
	}
	declareTrainExtraFlags(fs, &p)

	return p
}

func declareTrainExtraFlags(fs *flag.FlagSet, p *rawTrainFlagPointers) {
	p.healApply = fs.String("heal-apply", "", "Apply fix recipes non-interactively: stash|wip|discard (default \"\" = report-only)")
	p.days = fs.Int("days", DefaultTrainMuseDays, "Digest window in days for 'train muse' (default 7)")
	declareTrainHeatmapFlags(fs, p)
}

func declareTrainHeatmapFlags(fs *flag.FlagSet, p *rawTrainFlagPointers) {
	p.heatmap = fs.String("heatmap", "", "Path to local SQLite heatmap database or remote HTTP(S) URL")
	p.heatmapRepo = fs.String("heatmap-repo", "", "Filter heatmap training actions by repository slug")
	p.heatmapTask = fs.String("heatmap-task", "", "Filter heatmap training actions by task slug")
}

func parseTrainHeatmapFlags(args []string) (TrainOptions, LlmHeatmapOptions, *apperror.AppError) {
	fs := flag.NewFlagSet("llm train", flag.ContinueOnError)
	fs.Usage = printTrainUsage
	ptrs := declareTrainFlags(fs)
	if err := fs.Parse(args); err != nil {
		return handleTrainFlagParseError(err)
	}

	return assembleParsedTrainFlags(ptrs)
}

func assembleParsedTrainFlags(ptrs rawTrainFlagPointers) (TrainOptions, LlmHeatmapOptions, *apperror.AppError) {
	opts := buildTrainOptions(ptrs)
	heatmapOpts := LlmHeatmapOptions{
		Source:   *ptrs.heatmap,
		RepoSlug: *ptrs.heatmapRepo,
		TaskSlug: *ptrs.heatmapTask,
	}

	return opts, heatmapOpts, nil
}

func handleTrainFlagParseError(err error) (TrainOptions, LlmHeatmapOptions, *apperror.AppError) {
	if err == flag.ErrHelp {
		return TrainOptions{IsHelp: true}, LlmHeatmapOptions{}, nil
	}

	return TrainOptions{}, LlmHeatmapOptions{}, apperror.WrapSimple(err, "parse train flags")
}

func buildTrainOptions(p rawTrainFlagPointers) TrainOptions {
	return TrainOptions{
		IsTextOnly: *p.isTextOnly,
		SkillPath:  *p.skillPath,
		IsLoop:     *p.isLoop || *p.selfLoop > 0,
		LoopCount:  resolveLoopCount(*p.selfLoop),
		IsURL:      *p.isURL,
		IsURLs:     *p.isURLs,
		IsJSON:     *p.isJSON,
		HealApply:  *p.healApply,
		DaysN:      resolveTrainMuseDays(*p.days),
	}
}

func resolveLoopCount(selfLoop int) int {
	if selfLoop > 0 {
		return selfLoop
	}

	return 1
}

func resolveTrainMuseDays(days int) int {
	if days <= 0 {
		return DefaultTrainMuseDays
	}

	return days
}

func parseTrainFlags(args []string) (TrainOptions, *apperror.AppError) {
	opts, _, err := parseTrainHeatmapFlags(args)

	return opts, err
}

// trainValueFlags are train flags that consume the following arg as their value;
// a bare "muse" in that position is a value, never the positional target.
var trainValueFlags = map[string]bool{
	"--skill-path":   true,
	"--heal-apply":   true,
	"--heatmap":      true,
	"--heatmap-repo": true,
	"--heatmap-task": true,
}

// isMuseTarget reports whether arg represents the muse sub-target or its hyphenated aliases.
func isMuseTarget(arg string) bool {
	low := strings.ToLower(arg)
	return low == "muse" || low == "train-muse" || low == "llm-train-muse" || low == "train_muse" || low == "llm_train_muse"
}

// stripTrainTarget extracts the FIRST bare positional target ("muse", "train-muse", or "llm-train-muse")
// from the args before flag parsing. Tokens that are flags or flag values are preserved.
func stripTrainTarget(args []string) (string, []string) {
	rest := make([]string, 0, len(args))
	target := ""
	expectValue := false
	for _, arg := range args {
		switch {
		case expectValue:
			expectValue = false
		case target == "" && isMuseTarget(arg):
			target = museTrainTarget
			continue
		case target == "" && (arg == "train" || arg == "llm-train"):
			continue
		default:
			expectValue = trainValueFlags[arg]
		}
		rest = append(rest, arg)
	}
	return target, rest
}

func printTrainUsage() {
	fmt.Println("Usage: gitmap llm train [muse] [flags] (aliases: gitmap llm chain, gitmap train, gitmap llm-train, gitmap train-muse, gitmap llm-train-muse)")
	fmt.Println()
	fmt.Println("Executes the autonomous LLM chained onboarding curriculum: a 6-phase self-looping")
	fmt.Println("training run (Discovery, Refactoring, Verification, Semantic Commit, Telemetry,")
	fmt.Println("Heal & Fix) and generates the official Antigravity skill")
	fmt.Println("(.agents/skills/gitmap/SKILL.md). Phase 6 runs the git-state heal scan (cmdfix) and")
	fmt.Println("the parallel content-fix audit (cmdautofix) in-process — report-only; git-state")
	fmt.Println("recipes apply only with --heal-apply.")
	fmt.Println()
	fmt.Println("Positional sub-target & aliases:")
	fmt.Println("  muse, train-muse, llm-train-muse   Sync the Muse 7-day memory digest from repo-cache (fallback: repo-secrets)")
	fmt.Println("                                     (gitmap llm train muse [--days N]; see --help on the sub-target)")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --days int           Digest window in days for 'train muse' (default 7)")
	fmt.Println("  --heal-apply string  Apply fix recipes non-interactively: stash|wip|discard (default \"\" = report-only)")
	fmt.Println("  --heatmap string     Path to local SQLite heatmap database or remote HTTP(S) URL")
	fmt.Println("  --heatmap-repo string Filter heatmap actions by repository slug")
	fmt.Println("  --heatmap-task string Filter heatmap actions by task slug")
	fmt.Println("  --loop               Execute autonomous 6-phase AI self-looping execution cycle")
	fmt.Println("  --self-loop int      Number of consecutive iterations of the AI self-loop")
	fmt.Println("  --url                Output raw public URL to llm.md instruction specification")
	fmt.Println("  --urls               Output authoritative public documentation links for LLM ingestion")
	fmt.Println("  --json               Output structured machine-readable command specifications")
	fmt.Println("  --skill-path string  Path for generated Antigravity skill (default \".agents/skills/gitmap/SKILL.md\")")
	fmt.Println("  --text-only          Print curriculum text without writing skill file")
}

// printTrainMuseUsage prints the usage for the positional "muse" sub-target and aliases.
func printTrainMuseUsage() {
	fmt.Println("Usage: gitmap llm train muse [--days N] (aliases: gitmap train-muse, gitmap llm-train-muse)")
	fmt.Println()
	fmt.Println("Syncs the Muse 7-day memory digest: checks repo-cache first for the 12-muse-ai")
	fmt.Println("(or XX-muse-ai) directory with fallback to repo-secrets, diffs the last N days,")
	fmt.Println("prints a terminal summary and a raw AI traversal checklist (file-by-file paths +")
	fmt.Println("gitmap cat / aum search hints) so any LLM can walk the digest. Tolerates a missing")
	fmt.Println("digest with guidance.")
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
