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
	printAttribution()
	if err := handleSkillGeneration(opts); err != nil {
		return err
	}
	printChainedSequence()
	printOperationalGuidelines()
	return nil
}

func handleSkillGeneration(opts TrainOptions) *apperror.AppError {
	if opts.IsTextOnly {
		return nil
	}
	return executeSkillGeneration(opts.SkillPath)
}

func parseTrainFlags(args []string) (TrainOptions, *apperror.AppError) {
	fs := flag.NewFlagSet("llm train", flag.ContinueOnError)
	isTextOnly := fs.Bool("text-only", false, "Print curriculum text without writing skill file")
	skillPath := fs.String("skill-path", DefaultSkillPath, "Path for generated Antigravity skill")

	if err := fs.Parse(args); err != nil {
		return TrainOptions{}, apperror.WrapSimple(err, "parse train flags")
	}
	return TrainOptions{
		IsTextOnly: *isTextOnly,
		SkillPath:  *skillPath,
	}, nil
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
	fmt.Println("1. gitmap aum search --help          - Multi-core search and regex filters")
	fmt.Println("2. gitmap aum guard                  - 500 KB limit, large JSONs, binary null-byte probe")
	fmt.Println("3. gitmap aum sequence --help        - Numbering gap detector & H1 title validator")
	fmt.Println("4. gitmap aum exclude list           - Query persistent search exclusions from SQLite")
	fmt.Println("5. gitmap pipeline-ai status --json  - Check remote CI workflow state and dynamic ETA")
	fmt.Println("6. gitmap install --list             - Discover toolchains, profiles, and runtimes")
	fmt.Println("7. gitmap cluster --help             - Multi-node SSH and cluster execution")
	fmt.Println("8. gitmap cargo status               - Inspect Rust and Cargo toolchain status")
	fmt.Println("9. gitmap db status                  - Check repository SQLite database health")
	fmt.Println()
}

func printOperationalGuidelines() {
	fmt.Println(OperationalDirectivesText)
	fmt.Println("======================================================================")
}
