package orchestrator

import (
	"fmt"
	"io"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/dedupe"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/message"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/replay"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/runlog"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/walk"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/workspace"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// processOneCommit runs one source commit through dedupe → build →
// replay → record. Errors at any stage become Failed counters but
// never abort the loop (spec §3.4: best-effort partial success).
func processOneCommit(ctx *runContext, staged workspace.StagedInput, c walk.SourceCommit, pick func(int) int, stdout io.Writer) {
	inputRepoID, srcID, ok := persistSource(ctx, staged, c, stdout)
	if !ok {
		return
	}
	if handleDedupe(ctx, srcID, c, stdout) {
		return
	}
	keptFiles := applyExclusions(c.Files, ctx.Resolved.Exclusions)
	if len(c.Files) > 0 && len(keptFiles) == 0 {
		recordSkip(ctx, srcID, constants.CommitInSkipReasonExcludedAllFiles, stdout, c.Sha)
		return
	}
	c.Files = keptFiles
	intelBlock := renderFunctionIntel(staged.WorkPath, c, keptFiles, ctx.Resolved.FunctionIntel)
	finalMsg := buildMessage(ctx, c, intelBlock, pick)
	if finalMsg.IsEmpty {
		recordSkip(ctx, srcID, constants.CommitInSkipReasonEmptyAfterMessageRules, stdout, c.Sha)
		return
	}
	doReplayAndRecord(ctx, staged, c, finalMsg.Message, inputRepoID, srcID, stdout)
}

// persistSource inserts InputRepo (once per staged input via cache)
// and the SourceCommit row. Returns (inputRepoID, sourceCommitID, ok).
func persistSource(ctx *runContext, staged workspace.StagedInput, c walk.SourceCommit, stdout io.Writer) (int64, int64, bool) {
	inputRepoID, err := ctx.inputRepoID(staged)
	if err != nil {
		fmt.Fprintf(stdout, constants.CommitInErrDbWrite, err)
		ctx.Counters.Failed++
		return 0, 0, false
	}
	row := toSourceCommitRow(c)
	srcID, err := runlog.InsertSourceCommit(ctx.DB.Conn(), inputRepoID, row)
	if err != nil {
		fmt.Fprintf(stdout, constants.CommitInErrDbWrite, err)
		ctx.Counters.Failed++
		return 0, 0, false
	}
	return inputRepoID, srcID, true
}

func toSourceCommitRow(c walk.SourceCommit) runlog.SourceCommitRow {
	return runlog.SourceCommitRow{
		OrderIndex:           c.OrderIndex,
		Sha:                  c.Sha,
		AuthorName:           c.AuthorName,
		AuthorEmail:          c.AuthorEmail,
		AuthorDateRFC3339:    c.AuthorDate.Format(time.RFC3339),
		CommitterDateRFC3339: c.CommitterDate.Format(time.RFC3339),
		OriginalMessage:      c.OriginalMessage,
		Files:                c.Files,
	}
}

// handleDedupe checks ShaMap and records a Skip when it's a hit.
// Returns true when the caller should stop processing this commit.
func handleDedupe(ctx *runContext, srcID int64, c walk.SourceCommit, stdout io.Writer) bool {
	v, err := dedupe.Lookup(ctx.DB.Conn(), c.Sha)
	if err != nil {
		fmt.Fprintf(stdout, constants.CommitInErrDbWrite, err)
		ctx.Counters.Failed++
		return true
	}
	if !v.IsHit {
		return false
	}
	prev := v.PreviousRewrittenId
	_ = runlog.RecordSkip(ctx.DB.Conn(), ctx.RunID, srcID, constants.CommitInSkipReasonDuplicateSourceSha, &prev)
	fmt.Fprintf(stdout, constants.CommitInMsgCommitSkip, c.Sha, constants.CommitInSkipReasonDuplicateSourceSha)
	ctx.Counters.Skipped++
	return true
}

func buildMessage(ctx *runContext, c walk.SourceCommit, intelBlock string, pick func(int) int) message.Result {
	return message.Build(message.Inputs{
		OriginalMessage: c.OriginalMessage,
		FunctionIntel:   intelBlock,
		Resolved:        ctx.Resolved,
		PickIndex:       pick,
	})
}

func recordSkip(ctx *runContext, srcID int64, reason string, stdout io.Writer, sha string) {
	_ = runlog.RecordSkip(ctx.DB.Conn(), ctx.RunID, srcID, reason, nil)
	fmt.Fprintf(stdout, constants.CommitInMsgCommitSkip, sha, reason)
	ctx.Counters.Skipped++
}

func doReplayAndRecord(ctx *runContext, staged workspace.StagedInput, c walk.SourceCommit, msg string, inputRepoID, srcID int64, stdout io.Writer) {
	_ = inputRepoID
	if ctx.Raw.IsDryRun {
		recordSkip(ctx, srcID, constants.CommitInSkipReasonDryRun, stdout, c.Sha)
		return
	}
	plan := buildReplayPlan(ctx, staged, c, msg)
	if abort, skip := conflictCheck(ctx, plan, c, stdout); abort || skip {
		if abort {
			recordFail(ctx, srcID, c, msg, errConflictAborted, stdout)
		}
		return
	}
	res, err := replay.ApplyCommit(plan, false)
	if err != nil {
		recordFail(ctx, srcID, c, msg, err, stdout)
		return
	}
	recordCreated(ctx, srcID, c, msg, res.NewSha, stdout)
}

func buildReplayPlan(ctx *runContext, staged workspace.StagedInput, c walk.SourceCommit, msg string) replay.Plan {
	return replay.Plan{
		SourceRepoDir: staged.WorkPath,
		TargetRepoDir: ctx.Source.Path,
		SourceSha:     c.Sha,
		Files:         c.Files,
		Message:       msg,
		AuthorName:    pickAuthorName(ctx, c),
		AuthorEmail:   pickAuthorEmail(ctx, c),
		AuthorDate:    c.AuthorDate,
		CommitterDate: c.CommitterDate,
	}
}

func pickAuthorName(ctx *runContext, c walk.SourceCommit) string {
	if ctx.Resolved.Author != nil && ctx.Resolved.Author.Name != "" {
		return ctx.Resolved.Author.Name
	}
	return c.AuthorName
}

func pickAuthorEmail(ctx *runContext, c walk.SourceCommit) string {
	if ctx.Resolved.Author != nil && ctx.Resolved.Author.Email != "" {
		return ctx.Resolved.Author.Email
	}
	return c.AuthorEmail
}

func recordCreated(ctx *runContext, srcID int64, c walk.SourceCommit, msg, newSha string, stdout io.Writer) {
	row := runlog.RewrittenRow{
		NewSha:               newSha,
		SourceSha:            c.Sha,
		FinalMessage:         msg,
		AuthorName:           pickAuthorNameRow(ctx, c),
		AuthorEmail:          pickAuthorEmailRow(ctx, c),
		AuthorDateRFC3339:    c.AuthorDate.Format(time.RFC3339),
		CommitterDateRFC3339: c.CommitterDate.Format(time.RFC3339),
		Outcome:              constants.CommitInOutcomeCreated,
	}
	if _, err := runlog.RecordRewritten(ctx.DB.Conn(), ctx.RunID, srcID, row); err != nil {
		fmt.Fprintf(stdout, constants.CommitInErrDbWrite, err)
	}
	fmt.Fprintf(stdout, constants.CommitInMsgCommitOk, c.Sha, newSha, firstLine(msg))
	ctx.Counters.Created++
}

func recordFail(ctx *runContext, srcID int64, c walk.SourceCommit, msg string, cause error, stdout io.Writer) {
	row := runlog.RewrittenRow{
		SourceSha:            c.Sha,
		FinalMessage:         msg,
		AuthorName:           pickAuthorNameRow(ctx, c),
		AuthorEmail:          pickAuthorEmailRow(ctx, c),
		AuthorDateRFC3339:    c.AuthorDate.Format(time.RFC3339),
		CommitterDateRFC3339: c.CommitterDate.Format(time.RFC3339),
		Outcome:              constants.CommitInOutcomeFailed,
	}
	_, _ = runlog.RecordRewritten(ctx.DB.Conn(), ctx.RunID, srcID, row)
	fmt.Fprintf(stdout, constants.CommitInMsgCommitFail, c.Sha, cause)
	ctx.Counters.Failed++
}

func pickAuthorNameRow(ctx *runContext, c walk.SourceCommit) string {
	return pickAuthorName(ctx, c)
}

func pickAuthorEmailRow(ctx *runContext, c walk.SourceCommit) string {
	return pickAuthorEmail(ctx, c)
}
