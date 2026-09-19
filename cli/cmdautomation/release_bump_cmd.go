package cmdautomation

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	releaseBumpCmd = &cobra.Command{
		Use:     "release-bump [major|minor|patch]",
		Aliases: []string{"bump", "semver-bump"},
		Short:   "Compute next semantic version and update repository manifests",
		RunE:    runReleaseBumpCmd,
	}

	releaseBumpOpts ReleaseBumpOptions
)

func runReleaseBumpCmd(cmd *cobra.Command, args []string) error {
	hasArgs := len(args) > 0
	if hasArgs {
		releaseBumpOpts.BumpType = args[0]
	}
	monad := RunReleaseBump(releaseBumpOpts)
	isFail := monad.IsFailure()
	if isFail {
		return monad.Err
	}
	res := monad.Value
	renderReleaseBumpResult(res, releaseBumpOpts.IsJson, releaseBumpOpts.IsDryRun)
	return nil
}

func renderReleaseBumpResult(res ReleaseBumpResult, isJson, isDryRun bool) {
	if isJson {
		printReleaseBumpJson(res)
		return
	}
	printReleaseBumpTerminal(res, isDryRun)
}

func printReleaseBumpJson(res ReleaseBumpResult) {
	bytes, err := json.MarshalIndent(res, "", "  ")
	isSuccess := err == nil
	if isSuccess {
		fmt.Println(string(bytes))
	}
}

func printReleaseBumpTerminal(res ReleaseBumpResult, isDryRun bool) {
	printReleaseBumpHeader(res, isDryRun)
	printReleaseBumpDetails(res)
	printReleaseBumpFiles(res.UpdatedFiles)
	printReleaseBumpStatus(res, isDryRun)
}

func printReleaseBumpHeader(res ReleaseBumpResult, isDryRun bool) {
	modeLabel := "RELEASE BUMP"
	if isDryRun {
		modeLabel = "DRY RUN - RELEASE BUMP PREVIEW"
	}
	fmt.Printf("\n%s[%s]%s\n", constants.ColorBold, modeLabel, constants.ColorReset)
	fmt.Printf("  Previous Version: v%s\n", res.PreviousVersion)
	fmt.Printf("  New Version:      v%s\n", res.NewVersion)
	fmt.Printf("  Bump Tier:        %s\n", res.BumpType)
}

func printReleaseBumpDetails(res ReleaseBumpResult) {
	fmt.Printf("  Release Branch:   %s\n", res.ReleaseBranch)
	fmt.Printf("  Tag Name:         %s (Tagged: %t)\n", res.TagName, res.IsTagged)
	fmt.Printf("  Remote Push:      %t\n", res.IsPushed)
	fmt.Printf("  Duration:         %s\n", res.Duration)
}

func printReleaseBumpFiles(files []string) {
	fmt.Printf("\n%sUpdated Manifests:%s\n", constants.ColorCyan, constants.ColorReset)
	for _, f := range files {
		fmt.Printf("  - %s\n", f)
	}
}

func printReleaseBumpStatus(res ReleaseBumpResult, isDryRun bool) {
	if isDryRun {
		fmt.Printf("\n%sℹ️ DRY RUN COMPLETE: Manifests evaluated for v%s without disk changes.%s\n\n",
			constants.ColorCyan, res.NewVersion, constants.ColorReset)
		return
	}
	fmt.Printf("\n%s✅ SUCCESS: Version bumped to v%s across all SSoT manifests.%s\n\n",
		constants.ColorGreen, res.NewVersion, constants.ColorReset)
}

func init() {
	AutomationCmd.AddCommand(releaseBumpCmd)
	initReleaseBumpFlags()
}

func initReleaseBumpFlags() {
	releaseBumpCmd.Flags().StringVarP(&releaseBumpOpts.TargetVersion, "version", "v", "", "Explicit version string to bump to (e.g. 6.263.0)")
	releaseBumpCmd.Flags().BoolVarP(&releaseBumpOpts.IsDryRun, "dry-run", "d", false, "Preview version bump without modifying files")
	releaseBumpCmd.Flags().BoolVar(&releaseBumpOpts.IsTag, "tag", false, "Create git release tag (vX.Y.Z)")
	releaseBumpCmd.Flags().BoolVar(&releaseBumpOpts.IsPush, "push", false, "Push release branch and tag to remote")
	releaseBumpCmd.Flags().BoolVar(&releaseBumpOpts.IsSkipTests, "skip-tests", false, "Skip pre-release quality gate checks")
	releaseBumpCmd.Flags().StringVarP(&releaseBumpOpts.Scope, "scope", "s", "Automated release orchestration", "Scope description for release commit")
	releaseBumpCmd.Flags().StringSliceVarP(&releaseBumpOpts.Bullets, "bullet", "b", nil, "Changelog bullet points (repeatable)")
	releaseBumpCmd.Flags().BoolVar(&releaseBumpOpts.IsJson, "json", false, "Output results as machine-readable JSON")
}
