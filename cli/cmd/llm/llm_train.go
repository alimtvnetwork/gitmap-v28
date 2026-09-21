package llm

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunTrain executes the chained curriculum training workflow for LLMs.
func RunTrain(args []string) *apperror.AppError {
	opts, err := parseTrainFlags(args)
	if err != nil {
		return err
	}
	if opts.IsHelp {
		return nil
	}
	printAttribution()
	if err := handleSkillGeneration(opts); err != nil {
		return err
	}
	printCurriculumSummary()

	return nil
}

func printCurriculumSummary() {
	printChainedSequence()
	printOperationalGuidelines()
}

func handleSkillGeneration(opts TrainOptions) *apperror.AppError {
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

	err := fs.Parse(args)
	if err == flag.ErrHelp {
		return TrainOptions{IsHelp: true}, nil
	}
	if err != nil {
		return TrainOptions{}, apperror.WrapSimple(err, "parse train flags")
	}

	return TrainOptions{
		IsTextOnly: *isTextOnly,
		SkillPath:  *skillPath,
	}, nil
}

func printTrainUsage() {
	fmt.Println("Usage: gitmap llm train [flags] (alias: gitmap llm chain)")
	fmt.Println()
	fmt.Println("Executes 4-stage autonomous LLM chained onboarding curriculum and generates")
	fmt.Println("the official Antigravity skill file (.agents/skills/gitmap/SKILL.md).")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -skill-path string   Path for generated Antigravity skill (default \".agents/skills/gitmap/SKILL.md\")")
	fmt.Println("  -text-only           Print curriculum text without writing skill file")
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
	fmt.Println("7. gitmap install --list                        - Discover toolchains, profiles, and runtimes")
	fmt.Println("8. gitmap cluster --help                        - Multi-node SSH and cluster execution")
	fmt.Println("9. gitmap cargo status                          - Inspect Rust and Cargo toolchain status")
	fmt.Println("10. gitmap db status                            - Check repository SQLite database health")
	fmt.Println()
}

func printOperationalGuidelines() {
	fmt.Println(OperationalDirectivesText)
	fmt.Println("======================================================================")
}
