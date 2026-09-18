package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func assemblePrimaryAndFollowup(opts AgyFixOptions, payload cmdpipeline.PipelineErrorLogsPayload, errorReport string) (string, string, string, string) {
	gitLog := ExtractGitLog("", 5)
	promptContent, promptSource := LoadCanonicalRcaPrompt(opts.CustomPrompt, opts.IsNoRelease)
	primary := AssembleRcaFixPayload(payload.Repo, payload.RunId, payload.Sha, gitLog, errorReport, promptContent)
	followup := BuildVerificationFollowupPrompt(payload.Repo, payload.RunId, payload.Sha)

	return primary, followup, promptContent, promptSource
}

func finalizeFixFeedback(repo, promptSource, errorReport, promptContent, primaryPayload string, hasFailures, noClip bool) {
	renderAgyFixFeedback(repo, promptSource, errorReport, promptContent, primaryPayload, hasFailures, noClip)
	renderQueuedVerificationNotice()
}

func persistAndRecordAgyFix(p AgyFixDispatchParams, primary, followup string) error {
	if writeErr := persistFixPromptPayload(primary, p.Opts.OutputFile, p.Opts.IsNoClipboard); writeErr != nil {
		return apperror.WrapSimple(writeErr, "persist prompt payload")
	}

	_ = StageVerificationFollowupPrompt(primary, followup)
	_ = RecordSentErrorSignature(p.StorePath, p.Sig, p.Payload.Repo, p.Payload.RunId, p.Payload.Sha, p.ErrHash)

	return nil
}

func executeFixPayloadDispatch(p AgyFixDispatchParams, promptContent, promptSource, primary, followup string) error {
	if err := persistAndRecordAgyFix(p, primary, followup); err != nil {
		return err
	}

	finalizeFixFeedback(p.Payload.Repo, promptSource, p.ErrorReport, promptContent, primary, p.HasFailures, p.Opts.IsNoClipboard)
	repoDir := p.TargetDir
	if len(repoDir) == 0 {
		repoDir = resolveProjectRootDir()
	}
	if ok, msg := InjectAgyFixTask(repoDir, toAbsPath(resolveActiveAgyPromptPath()), p.Opts.IsNoInject); ok {
		fmt.Printf("  %s✔ %s%s\n\n", constants.ColorGreen, msg, constants.ColorReset)
	} else if !p.Opts.IsNoInject {
		fmt.Printf("  %sℹ Antigravity Injection: %s%s\n\n", constants.ColorYellow, msg, constants.ColorReset)
	}

	return nil
}

func dispatchAgyFixPrepared(p AgyFixDispatchParams) error {
	primary, followup, promptContent, promptSource := assemblePrimaryAndFollowup(p.Opts, p.Payload, p.ErrorReport)
	if p.Opts.IsDryRun {
		renderAgyFixDryRun(p.Payload.Repo, promptSource, p.ErrorReport, promptContent, primary, p.HasFailures)

		return nil
	}

	return executeFixPayloadDispatch(p, promptContent, promptSource, primary, followup)
}

func dispatchCandidateFix(cand ProjectFixCandidate, opts AgyFixOptions) error {
	payload := cmdpipeline.PipelineErrorLogsPayload{
		Repo: cand.RepoSlug, RunId: cand.RunID, Sha: cand.SHA, ErrorLogs: cand.ErrorReport,
	}
	sig, errHash := ComputeErrorSignature(cand.RepoSlug, cand.RunID, cand.SHA, cand.ErrorReport)
	storePath := sentAgyErrorsStorePath()
	params := AgyFixDispatchParams{
		Opts: opts, StorePath: storePath, Sig: sig, ErrHash: errHash,
		Payload: payload, ErrorReport: cand.ErrorReport, HasFailures: cand.HasFailures,
		TargetDir: cand.Path,
	}

	return dispatchAgyFixPrepared(params)
}

func executeSingleAgyFix(opts AgyFixOptions) error {
	payload, errorReport, hasFailures := cmdpipeline.FetchPipelineErrorReportWithMeta(opts.Repo, opts.IsDetailed)
	isDup, storePath, sig, errHash := checkAgyFixDuplicate(opts, payload, errorReport)
	if isDup {
		return nil
	}

	params := AgyFixDispatchParams{
		Opts: opts, StorePath: storePath, Sig: sig, ErrHash: errHash,
		Payload: payload, ErrorReport: errorReport, HasFailures: hasFailures,
		TargetDir: resolveProjectRootDir(),
	}

	return dispatchAgyFixPrepared(params)
}
