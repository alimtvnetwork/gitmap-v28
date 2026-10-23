package cmdcommit

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RenderCommitInfoReport outputs the commit report as JSON or ASCII tree.
func RenderCommitInfoReport(w io.Writer, rep *CommitInfoReport, isJson bool) error {
	if isJson {
		return renderCommitInfoJSON(w, rep)
	}

	renderCommitInfoTree(w, rep)

	return nil
}

func renderCommitInfoJSON(w io.Writer, rep *CommitInfoReport) error {
	data, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal commit info json")
	}

	_, writeErr := fmt.Fprintln(w, string(data))
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "write commit info json")
	}

	return nil
}

func renderCommitInfoTree(w io.Writer, rep *CommitInfoReport) {
	renderCommitHeader(w, rep)
	renderCommitMeta(w, rep)
	renderAiAttribution(w, rep)
	renderOriginSection(w, rep)
	renderTaskSpecSection(w, rep)
	renderFileDeltasSection(w, rep)
	renderSummaryFooter(w, rep)
}

func renderCommitHeader(w io.Writer, rep *CommitInfoReport) {
	fmt.Fprintf(w, "%s%scommit %s%s (short: %s%s%s)\n",
		constants.ColorBold, constants.ColorNeonGreen,
		rep.FullSha, constants.ColorReset,
		constants.ColorNeonGold, rep.ShortSha,
		constants.ColorReset,
	)
}

func renderCommitMeta(w io.Writer, rep *CommitInfoReport) {
	b := constants.TreeBranch
	c := constants.ColorCyan
	r := constants.ColorReset
	fmt.Fprintf(w, "%s %sAuthor:%s      %s <%s>\n", b, c, r, rep.AuthorName, rep.AuthorEmail)
	fmt.Fprintf(w, "%s %sDate:%s        %s (%s)\n", b, c, r, rep.FormattedDate, rep.RelativeDate)
	fmt.Fprintf(w, "%s %sBranch:%s      %s\n", b, c, r, rep.BranchName)
	fmt.Fprintf(w, "%s %sMessage:%s     %s\n", b, c, r, rep.CommitMessage)
}

func renderOriginSection(w io.Writer, rep *CommitInfoReport) {
	b := constants.TreeBranch
	c := constants.ColorCyan
	r := constants.ColorReset
	color := getOriginColor(rep.Provenance.OriginType)
	fmt.Fprintf(w, "%s %sOrigin:%s      %s%s%s\n", b, c, r, color, rep.Provenance.AttributionText, r)
}

func getOriginColor(origin OriginType) string {
	switch origin {
	case OriginLocalAI:
		return constants.ColorMagenta
	case OriginForeignAI:
		return constants.ColorYellow
	default:
		return constants.ColorDim
	}
}

func renderAiAttribution(w io.Writer, rep *CommitInfoReport) {
	if !rep.Provenance.IsAiCommitted {
		return
	}

	p := rep.Provenance
	pipe := constants.TreePipe
	fmt.Fprintf(w, "%s %sAI Attribution:%s\n", constants.TreeBranch, constants.ColorCyan, constants.ColorReset)
	renderAiAttrLine(w, pipe, constants.TreeBranch, "Model", p.AiModel)
	renderAiAttrLine(w, pipe, constants.TreeBranch, "Tool Name", p.ToolName)
	renderAiAttrLine(w, pipe, constants.TreeBranch, "Transmission ID", p.TransmissionId)
	renderAiAttrLine(w, pipe, constants.TreeBranch, "Agent ID", p.AgentId)
	renderAiAttrLine(w, pipe, constants.TreeCorner, "Heatmap ID", p.HeatmapId)
}

func renderAiAttrLine(w io.Writer, pipe, branch, label, val string) {
	if len(val) == 0 {
		val = "-"
	}

	fmt.Fprintf(w, "%s%s %-16s %s\n", pipe, branch, label+":", val)
}

func renderTaskSpecSection(w io.Writer, rep *CommitInfoReport) {
	p := rep.Provenance
	hasTask := len(p.TaskSlug) > 0 || len(rep.SpecPath) > 0
	if !hasTask {
		return
	}

	pipe := constants.TreePipe
	fmt.Fprintf(w, "%s %sTask & Specification:%s\n", constants.TreeBranch, constants.ColorCyan, constants.ColorReset)
	renderTaskSpecLine(w, pipe, constants.TreeBranch, "Task Slug", p.TaskSlug)
	renderTaskSpecLine(w, pipe, constants.TreeBranch, "Spec Link", rep.SpecPath)
	subDesc := formatSubtaskDesc(rep)
	renderTaskSpecLine(w, pipe, constants.TreeBranch, "Subtasks", subDesc)
	planDesc := fmt.Sprintf("%d plan(s)", rep.PlanCount)
	renderTaskSpecLine(w, pipe, constants.TreeCorner, "Plans", planDesc)
}

func formatSubtaskDesc(rep *CommitInfoReport) string {
	code := rep.CrossReferences.SubtaskCode
	status := rep.CrossReferences.SubtaskStatus
	if len(code) > 0 {
		return fmt.Sprintf("%s (%s)", code, status)
	}

	return fmt.Sprintf("%d subtask(s)", rep.SubtaskCount)
}

func renderTaskSpecLine(w io.Writer, pipe, branch, label, val string) {
	if len(val) == 0 {
		val = "-"
	}

	fmt.Fprintf(w, "%s%s %-14s %s\n", pipe, branch, label+":", val)
}

func renderFileDeltasSection(w io.Writer, rep *CommitInfoReport) {
	hasFiles := len(rep.FileDeltas) > 0
	if !hasFiles {
		return
	}

	fmt.Fprintf(w, "%s %sFiles Modified & Cognitive Intention:%s\n", constants.TreeBranch, constants.ColorCyan, constants.ColorReset)
	pipe := constants.TreePipe
	renderCognitiveIntention(w, pipe, rep)
	for idx, delta := range rep.FileDeltas {
		isLast := idx == len(rep.FileDeltas)-1
		branch := constants.TreeBranch
		if isLast {
			branch = constants.TreeCorner
		}
		renderSingleFileDelta(w, pipe, branch, delta)
	}
}

func renderCognitiveIntention(w io.Writer, pipe string, rep *CommitInfoReport) {
	if len(rep.CognitiveIntention) > 0 {
		fmt.Fprintf(w, "%s%s %-20s %s\n", pipe, constants.TreeBranch, "Cognitive Intention:", rep.CognitiveIntention)
	}
	if rep.IsExpectationMet {
		fmt.Fprintf(w, "%s%s %-20s %sYes%s\n", pipe, constants.TreeBranch, "Expectation Met:", constants.ColorGreen, constants.ColorReset)
	}
}

func renderSingleFileDelta(w io.Writer, pipe, branch string, delta CommitFileDelta) {
	cGreen := constants.ColorGreen
	cRed := constants.ColorRed
	cReset := constants.ColorReset
	if delta.IsBinary {
		fmt.Fprintf(w, "%s%s %-36s (binary)\n", pipe, branch, delta.FilePath)

		return
	}

	diffStr := fmt.Sprintf("(%s+%d%s, %s-%d%s)", cGreen, delta.Insertions, cReset, cRed, delta.Deletions, cReset)
	fmt.Fprintf(w, "%s%s %-36s %s\n", pipe, branch, delta.FilePath, diffStr)
}

func renderSummaryFooter(w io.Writer, rep *CommitInfoReport) {
	fmt.Fprintf(w, "%s %sSummary:%s    %d files changed, %d insertions(+), %d deletions(-)\n",
		constants.TreeCorner,
		constants.ColorCyan,
		constants.ColorReset,
		rep.TotalFiles,
		rep.TotalInsertions,
		rep.TotalDeletions,
	)
}
