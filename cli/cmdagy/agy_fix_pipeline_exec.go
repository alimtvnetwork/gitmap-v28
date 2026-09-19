package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func assemblePrimaryAndFollowup(
	opts AgyFixOptions,
	payload cmdpipeline.PipelineErrorLogsPayload,
	errorReport string,
) AgyAssembledPromptPayload {
	gitLog := ExtractGitLog("", 5)
	promptContent, promptSource := LoadCanonicalRcaPrompt(opts.CustomPrompt, opts.IsNoRelease)
	primary := AssembleRcaFixPayload(payload.Repo, payload.RunId, payload.Sha, gitLog, errorReport, promptContent)
	followup := BuildVerificationFollowupPrompt(payload.Repo, payload.RunId, payload.Sha)

	return AgyAssembledPromptPayload{
		Primary:       primary,
		Followup:      followup,
		PromptContent: promptContent,
		PromptSource:  promptSource,
	}
}

func finalizeFixFeedback(
	repo string,
	promptSource string,
	errorReport string,
	promptContent string,
	primaryPayload string,
	hasFailures bool,
	isNoClipboard bool,
) {
	renderAgyFixFeedback(repo, promptSource, errorReport, promptContent, primaryPayload, hasFailures, isNoClipboard)
	renderQueuedVerificationNotice()
}

func persistAndRecordAgyFix(p AgyFixDispatchParams, primary, followup string) *apperror.AppError {
	if writeErr := persistFixPromptPayload(primary, p.Opts.OutputFile, p.Opts.IsNoClipboard); writeErr != nil {
		return apperror.WrapSimple(writeErr, "persist prompt payload")
	}

	_ = StageVerificationFollowupPrompt(primary, followup)
	_ = RecordSentErrorSignature(p.StorePath, p.Sig, p.Payload.Repo, p.Payload.RunId, p.Payload.Sha, p.ErrHash)

	return nil
}

func renderInjectionFeedback(res AgyInjectionResult, isNoInject bool) {
	if res.IsSuccess {
		fmt.Printf("  %s✔ %s%s\n\n", constants.ColorGreen, res.Message, constants.ColorReset)

		return
	}

	if !isNoInject {
		fmt.Printf("  %sℹ Antigravity Injection: %s%s\n\n", constants.ColorYellow, res.Message, constants.ColorReset)
	}
}

func executeFixPayloadDispatch(p AgyFixDispatchParams, assembled AgyAssembledPromptPayload) *apperror.AppError {
	if err := persistAndRecordAgyFix(p, assembled.Primary, assembled.Followup); err != nil {
		return err
	}

	finalizeFixFeedback(p.Payload.Repo, assembled.PromptSource, p.ErrorReport, assembled.PromptContent, assembled.Primary, p.HasFailures, p.Opts.IsNoClipboard)
	repoDir := resolveDispatchRepoDir(p.TargetDir)
	injectRes := InjectAgyFixTask(repoDir, toAbsPath(resolveActiveAgyPromptPath()), p.Opts.IsNoInject)
	renderInjectionFeedback(injectRes, p.Opts.IsNoInject)

	return nil
}

func resolveDispatchRepoDir(targetDir string) string {
	if len(targetDir) > 0 {
		return targetDir
	}

	return resolveProjectRootDir()
}

func dispatchAgyFixPrepared(p AgyFixDispatchParams) *apperror.AppError {
	assembled := assemblePrimaryAndFollowup(p.Opts, p.Payload, p.ErrorReport)
	if p.Opts.IsDryRun {
		renderAgyFixDryRun(p.Payload.Repo, assembled.PromptSource, p.ErrorReport, assembled.PromptContent, assembled.Primary, p.HasFailures)

		return nil
	}

	return executeFixPayloadDispatch(p, assembled)
}

func dispatchCandidateFix(cand ProjectFixCandidate, opts AgyFixOptions) *apperror.AppError {
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

func executeSingleAgyFix(opts AgyFixOptions) *apperror.AppError {
	payload, errorReport, hasFailures := cmdpipeline.FetchPipelineErrorReportWithMeta(opts.Repo, opts.IsDetailed)
	isDup, storePath, sig, errHash := checkAgyFixDuplicate(opts, payload, errorReport)
	if isDup {
		return nil
	}

	params := AgyFixDispatchParams{
		Opts: opts, StorePath: storePath, Sig: sig, ErrHash: errHash,
		Payload: payload, ErrorReport: errorReport, HasFailures: hasFailures,
		TargetDir: resolveDispatchRepoDir(""),
	}

	return dispatchAgyFixPrepared(params)
}
