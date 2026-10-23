package llm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// PartitionHeatmapExamples partitions records into positive reinforcement and negative anti-pattern examples.
func PartitionHeatmapExamples(examples []HeatmapTrainingExample) ([]HeatmapTrainingExample, []HeatmapTrainingExample) {
	var positive []HeatmapTrainingExample
	var negative []HeatmapTrainingExample

	for _, ex := range examples {
		if ex.IsExpectationMet {
			positive = append(positive, ex)
		} else {
			negative = append(negative, ex)
		}
	}

	return positive, negative
}

// FormatHeatmapTrainingDigest formats a structured Markdown digest partitioned into positive and negative traces.
func FormatHeatmapTrainingDigest(examples []HeatmapTrainingExample) string {
	hasExamples := len(examples) > 0
	if !hasExamples {
		return "No heatmap training actions recorded.\n"
	}

	positive, negative := PartitionHeatmapExamples(examples)

	var sb strings.Builder
	writeDigestHeader(&sb, len(examples), len(positive), len(negative))
	writePositiveDigest(&sb, positive)
	writeNegativeDigest(&sb, negative)

	return sb.String()
}

func writeDigestHeader(sb *strings.Builder, total, pos, neg int) {
	sb.WriteString("======================================================================\n")
	sb.WriteString("HEATMAP COGNITIVE DECISION TRACES & TRAINING DIGEST\n")
	sb.WriteString("======================================================================\n")
	sb.WriteString(fmt.Sprintf("Total Telemetry Actions: %d | Positive Reinforcements: %d | Negative RCA: %d\n\n", total, pos, neg))
}

func writePositiveDigest(sb *strings.Builder, positive []HeatmapTrainingExample) {
	hasPositive := len(positive) > 0
	if !hasPositive {
		return
	}

	sb.WriteString("### [Positive Reinforcement Traces (Verified Successes)]\n")
	renderPositiveItems(sb, positive)
	sb.WriteString("\n")
}

func renderPositiveItems(sb *strings.Builder, positive []HeatmapTrainingExample) {
	limit := resolveDigestLimit(len(positive))
	for i := 0; i < limit; i++ {
		p := positive[i]
		sb.WriteString(fmt.Sprintf("- [%s] Task: %s | Action: %s | Tool: %s | File: %s\n", p.ReasoningState, p.TaskSlug, p.ActionType, p.ToolName, p.FilePath))
		if p.ExpectedOutcome != "" {
			sb.WriteString(fmt.Sprintf("  Expected: %s -> Actual: %s\n", p.ExpectedOutcome, p.ActualOutcome))
		}
	}
}

func writeNegativeDigest(sb *strings.Builder, negative []HeatmapTrainingExample) {
	hasNegative := len(negative) > 0
	if !hasNegative {
		return
	}

	sb.WriteString("### [Negative Anti-Pattern Traces (RCA Prevention)]\n")
	renderNegativeItems(sb, negative)
	sb.WriteString("\n")
}

func renderNegativeItems(sb *strings.Builder, negative []HeatmapTrainingExample) {
	limit := resolveDigestLimit(len(negative))
	for i := 0; i < limit; i++ {
		n := negative[i]
		sb.WriteString(fmt.Sprintf("- [%s] Task: %s | Action: %s | Tool: %s | File: %s\n", n.ReasoningState, n.TaskSlug, n.ActionType, n.ToolName, n.FilePath))
		if n.ExpectedOutcome != "" {
			sb.WriteString(fmt.Sprintf("  Divergence: Expected %q, got %q\n", n.ExpectedOutcome, n.ActualOutcome))
		}
	}
}

func resolveDigestLimit(count int) int {
	if count < 10 {
		return count
	}

	return 10
}

// FormatHeatmapTrainingJsonl formats examples as newline-delimited JSON.
func FormatHeatmapTrainingJsonl(examples []HeatmapTrainingExample) (string, error) {
	var sb strings.Builder
	for _, ex := range examples {
		bytes, err := json.Marshal(ex)
		if err != nil {
			return "", apperror.WrapSimple(err, "llm.heatmap.marshal_jsonl")
		}
		sb.Write(bytes)
		sb.WriteByte('\n')
	}

	return sb.String(), nil
}
