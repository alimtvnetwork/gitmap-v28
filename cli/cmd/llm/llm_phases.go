package llm

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// TrainPhase is one executable stage of the llm train curriculum.
type TrainPhase interface {
	Name() string
	Run(opts TrainOptions) *apperror.AppError
}

// TrainPhaseResult is the outcome line of one phase for the closing summary.
type TrainPhaseResult struct {
	Name   string
	Status string // "ok" | "skipped" | "error"
	Detail string
}

// trainPhases is the ordered curriculum registry:
// Discovery -> Refactoring -> Verification -> Semantic Commit -> Telemetry -> Heal & Fix.
var trainPhases = []TrainPhase{
	discoveryPhase{},
	refactorPhase{},
	verificationPhase{},
	commitPhase{},
	telemetryPhase{},
	healPhase{},
}

// runTrainPhases executes the registry in order, printing a per-stage header
// for each phase and a closing summary table. It stops at the first error.
func runTrainPhases(opts TrainOptions) *apperror.AppError {
	results := make([]TrainPhaseResult, 0, len(trainPhases))
	for i, phase := range trainPhases {
		printTrainStageHeader(i+1, phase.Name())
		if err := phase.Run(opts); err != nil {
			results = append(results, TrainPhaseResult{Name: phase.Name(), Status: "error", Detail: err.Error()})
			printTrainPhaseSummary(results)
			return err
		}
		results = append(results, TrainPhaseResult{Name: phase.Name(), Status: "ok", Detail: "completed"})
	}
	printTrainPhaseSummary(results)
	return nil
}

func printTrainStageHeader(stage int, name string) {
	fmt.Println("======================================================================")
	fmt.Printf("STAGE %d — %s\n", stage, name)
	fmt.Println("======================================================================")
}

func printTrainPhaseSummary(results []TrainPhaseResult) {
	fmt.Println("======================================================================")
	fmt.Println("TRAIN PHASE SUMMARY")
	fmt.Println("======================================================================")
	fmt.Println("PHASE | STATUS | DETAIL")
	for _, r := range results {
		fmt.Printf("%s | %s | %s\n", r.Name, r.Status, r.Detail)
	}
	fmt.Println()
}

// discoveryPhase prints the curriculum attribution and the skill-ingestion
// mandate, then generates the Antigravity skill (phase-agnostic file write).
type discoveryPhase struct{}

func (discoveryPhase) Name() string { return "1. Discovery" }

func (discoveryPhase) Run(opts TrainOptions) *apperror.AppError {
	printAttribution()
	fmt.Print(RenderSkillCreationDirective())
	return handleSkillGeneration(opts)
}

// refactorPhase prints the authoritative public documentation registry.
type refactorPhase struct{}

func (refactorPhase) Name() string { return "2. Refactoring" }

func (refactorPhase) Run(opts TrainOptions) *apperror.AppError {
	fmt.Print(RenderPublicDocLinksText())
	return nil
}

// verificationPhase prints the recursive Git network learning directives.
type verificationPhase struct{}

func (verificationPhase) Name() string { return "3. Verification" }

func (verificationPhase) Run(opts TrainOptions) *apperror.AppError {
	fmt.Print(RenderRecursiveInstructions())
	return nil
}

// commitPhase prints the chained discovery sequence.
type commitPhase struct{}

func (commitPhase) Name() string { return "4. Semantic Commit" }

func (commitPhase) Run(opts TrainOptions) *apperror.AppError {
	printChainedSequence()
	return nil
}

// telemetryPhase prints the operational best-practice guardrail directives.
type telemetryPhase struct{}

func (telemetryPhase) Name() string { return "5. Telemetry" }

func (telemetryPhase) Run(opts TrainOptions) *apperror.AppError {
	printOperationalGuidelines()
	return nil
}
