// Package cmdfixreleasetags provides the main entrypoint and orchestration
// for auditing and deleting orphan or broken release tags.
package cmdfixreleasetags

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunFixReleaseTags is the primary orchestrator for the fix-release-tags command.
func RunFixReleaseTags(args []string) error {
	for _, arg := range args {
		if isHelpFlag(arg) {
			RenderFixReleaseTagsHelp()
			return nil
		}
	}

	flags, err := ParseFixReleaseTagsFlags(args)
	if err != nil {
		return err
	}

	if flags.IsHelpRequested {
		RenderFixReleaseTagsHelp()
		return nil
	}

	auditOpts := AuditFilterOptions{}
	report, err := AuditReleaseTags(flags.TargetDirectory, auditOpts)
	if err != nil {
		return fmt.Errorf("release tags audit failed: %w", err)
	}

	if flags.IsJSON {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(report)
	}

	RenderHeaderBanner(flags.TargetDirectory)
	RenderDiagnosticTable(report.Records)

	var candidates []ReleaseTagAuditRecord
	for _, rec := range report.Records {
		if rec.IsEligibleForDeletion {
			candidates = append(candidates, rec)
		}
	}

	if len(candidates) == 0 {
		fmt.Printf("\n%s✓%s All release tags are healthy and accompanied by valid release assets.\n",
			constants.ColorGreen, constants.ColorReset)
		return nil
	}

	if flags.IsDryRun {
		fmt.Printf("\n%s[dry-run]%s Dry-run completed. %d candidate(s) would be deleted.\n",
			constants.ColorYellow, constants.ColorReset, len(candidates))
		return nil
	}

	confirmed, err := ConfirmDeletion(len(candidates), flags.IsConfirmed)
	if err != nil {
		return err
	}

	if !confirmed {
		fmt.Println("\nAborted by user.")
		return nil
	}

	results, err := ExecuteBatchDeletions(flags.TargetDirectory, candidates, false, nil)
	if err != nil {
		return fmt.Errorf("release tags deletion execution failed: %w", err)
	}

	report.DeletionResults = results

	deletedCount := 0
	for _, res := range results {
		if res.IsSuccess {
			deletedCount++
		}
	}

	RenderSummaryBanner(report.Summary.TotalTagsChecked, len(candidates), deletedCount, false)

	return nil
}

// RunCLI forwards execution to RunFixReleaseTags for compatibility with all callers.
func RunCLI(args []string) error {
	return RunFixReleaseTags(args)
}
