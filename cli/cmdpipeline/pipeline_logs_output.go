package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

func writeOrRenderErrorLogs(params ErrorLogOutputParams) error {
	contentToWrite, err := formatErrorLogContent(params)
	if err != nil {
		return err
	}
	_ = persistAutoErrorReport(params)

	if len(params.TempFile) > 0 || len(params.FilePath) > 0 {
		return writeErrorLogsToDisk(params, contentToWrite)
	}

	return dispatchErrorLogPresentation(params, contentToWrite)
}

func dispatchErrorLogPresentation(params ErrorLogOutputParams, content string) error {
	if params.IsJSON {
		return outputJSONErrorLogs(content)
	}
	if params.HasSuppressOutputLog {
		printSuppressedStagingNotice(params.Payload.SavedReportFile)

		return nil
	}
	renderErrorLogsTerminal(params.Payload)

	hasFailure := params.Payload.Conclusion == "failure" || len(params.Payload.FailedRuns) > 0
	if hasFailure && params.WantFix && PipelineAgyFixRunner != nil {
		fmt.Printf("\n  🚀 Automatically dispatching CI/CD fix to Antigravity IDE...\n")
		return PipelineAgyFixRunner([]string{params.Payload.Repo, "--force"})
	}

	return nil
}

func outputJSONErrorLogs(content string) error {
	fmt.Println(content)
	if isClipboardWriteAllowed(false) {
		_ = writeClipboard(content)
	}

	return nil
}

func printSuppressedStagingNotice(reportFile string) {
	if len(reportFile) > 0 {
		fmt.Printf("  ✓ Error logs staged to %s (suppressed terminal output via --no-output-log)\n", reportFile)

		return
	}

	fmt.Println("  ✓ Error logs staged to filesystem (suppressed terminal output via --no-output-log)")
}

func writeErrorLogsToDisk(params ErrorLogOutputParams, content string) error {
	if isClipboardWriteAllowed(false) {
		_ = writeClipboard(content)
	}
	if len(params.TempFile) > 0 {
		targetPath := filepath.Join(resolveTempDir(), params.TempFile)

		return writeContentToFile(targetPath, content)
	}

	return writeContentToFile(params.FilePath, content)
}

func persistAutoErrorReport(params ErrorLogOutputParams) error {
	if !isFailingConclusion(params.Payload.Conclusion) && len(params.Payload.FailedRuns) == 0 {
		clearLocalErrorLogs()

		return nil
	}

	return saveActiveErrorReport(params.Payload)
}

func saveActiveErrorReport(p PipelineErrorLogsPayload) error {
	reportContent := p.ErrorLogs
	if len(reportContent) == 0 {
		reportContent = p.CombinedErrors
	}

	if len(reportContent) == 0 {
		return nil
	}

	_, err := writeCombinedErrorReportForRepo(p.Repo, reportContent)

	return err
}

func formatErrorLogContent(params ErrorLogOutputParams) (string, error) {
	if !params.IsJSON {
		return params.Payload.ErrorLogs, nil
	}

	b, err := json.MarshalIndent(params.Payload, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
