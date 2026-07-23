package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// printGroupUtilities prints the utility commands.
func printGroupUtilities() {
	fmt.Println()
	fmt.Println(constants.HelpGroupUtilities)
	fmt.Println(constants.HelpSetup)
	fmt.Println(constants.HelpDoctor)
	fmt.Println(constants.HelpUpdate)
	fmt.Println(constants.HelpUpdateCleanup)
	fmt.Println(constants.HelpVersion)
	fmt.Println(constants.HelpCompletion)
	fmt.Println(constants.HelpInteractive)
	fmt.Println(constants.HelpDocs)
	fmt.Println(constants.HelpHelpDash)
	fmt.Println(constants.HelpGoMod)
	fmt.Println(constants.HelpSEOWrite)
	fmt.Println(constants.HelpLLMDocs)
	fmt.Println(constants.HelpFixRepo)
	fmt.Println(constants.HelpMakePublic)
	fmt.Println(constants.HelpMakePrivate)
	fmt.Println(constants.HelpCloneFixRepo)
	fmt.Println(constants.HelpCloneFixRepoPub)
	fmt.Println(constants.HelpCmdOpen)
	fmt.Println(constants.HelpHelp)
}

// printUsageFlagSections prints all flag detail sections.
func printUsageFlagSections() {
	printUsageScanFlags()
	printUsageCloneFlags()
	printUsageReleaseFlags()
	printUsageSEOFlags()
	printUsageAmendFlags()
	printUsageGoModFlags()
	printUsageInteractiveFlags()
	printUsageCloneNextFlags()
	printUsageFixRepoFlags()
}

// printUsageFixRepoFlags prints the fix-repo flags section so the
// -2 / -3 / -5 / --all / --dry-run family is discoverable from the
// top-level `gitmap help` output (not only `gitmap help fix-repo`).
func printUsageFixRepoFlags() {
	fmt.Println()
	fmt.Println(constants.HelpFixRepoFlags)
	fmt.Println(constants.HelpFRMode2)
	fmt.Println(constants.HelpFRMode3)
	fmt.Println(constants.HelpFRMode5)
	fmt.Println(constants.HelpFRAll)
	fmt.Println(constants.HelpFRDryRun)
	fmt.Println(constants.HelpFRVerbose)
	fmt.Println(constants.HelpFRConfig)
	fmt.Println(constants.HelpFRStrict)
	fmt.Println(constants.HelpFRRestrict)
	fmt.Println(constants.HelpFRExample1)
	fmt.Println(constants.HelpFRExample2)
	fmt.Println(constants.HelpFRGofmtMaxCmdLen)
	fmt.Println(constants.HelpFixRepoExitCodes)
}

// printUsageCloneNextFlags prints the clone-next flags section.
func printUsageCloneNextFlags() {
	fmt.Println()
	fmt.Println(constants.HelpCloneNextFlags)
	fmt.Println(constants.HelpCNDelete)
	fmt.Println(constants.HelpCNKeep)
	fmt.Println(constants.HelpCNNoDesktop)
	fmt.Println(constants.HelpCNSSHKey)
	fmt.Println(constants.HelpCNVerbose)
	fmt.Println(constants.HelpCNCreateRemote)
}

// printUsageInteractiveFlags prints the interactive flags section.
func printUsageInteractiveFlags() {
	fmt.Println()
	fmt.Println(constants.HelpInteractiveFlags)
	fmt.Println(constants.HelpRefresh)
}

// printUsageScanFlags prints the scan flags section.
func printUsageScanFlags() {
	fmt.Println()
	fmt.Println(constants.HelpScanFlags)
	fmt.Println(constants.HelpConfig)
	fmt.Println(constants.HelpMode)
	fmt.Println(constants.HelpOutput)
	fmt.Println(constants.HelpOutputPath)
	fmt.Println(constants.HelpOutFile)
	fmt.Println(constants.HelpScanFlagGitHubDesktop)
	fmt.Println(constants.HelpOpen)
	fmt.Println(constants.HelpQuiet)
}

// printUsageCloneFlags prints the clone flags section.
func printUsageCloneFlags() {
	fmt.Println()
	fmt.Println(constants.HelpCloneFlags)
	fmt.Println(constants.HelpTargetDir)
	fmt.Println(constants.HelpSafePull)
	fmt.Println(constants.HelpVerbose)
}

// printUsageReleaseFlags prints the release flags section.
func printUsageReleaseFlags() {
	fmt.Println()
	fmt.Println(constants.HelpReleaseFlags)
	fmt.Println(constants.HelpAssets)
	fmt.Println(constants.HelpCommit)
	fmt.Println(constants.HelpRelBranch)
	fmt.Println(constants.HelpBump)
	fmt.Println(constants.HelpDraft)
	fmt.Println(constants.HelpDryRun)
	fmt.Println(constants.HelpCompressFlag)
	fmt.Println(constants.HelpChecksumsFlag)
	fmt.Println(constants.HelpBin)
	fmt.Println(constants.HelpTargets)
	fmt.Println(constants.HelpListTargets)
}

// printUsageSEOFlags prints the seo-write flags section.
func printUsageSEOFlags() {
	fmt.Println()
	fmt.Println(constants.HelpSEOWriteFlags)
	fmt.Println(constants.HelpSEOCSV)
	fmt.Println(constants.HelpSEOURL)
	fmt.Println(constants.HelpSEOService)
	fmt.Println(constants.HelpSEOArea)
	fmt.Println(constants.HelpSEOCompany)
	fmt.Println(constants.HelpSEOPhone)
	fmt.Println(constants.HelpSEOEmail)
	fmt.Println(constants.HelpSEOAddress)
	fmt.Println(constants.HelpSEOMaxCommits)
	fmt.Println(constants.HelpSEOInterval)
	fmt.Println(constants.HelpSEOFilesFlag)
	fmt.Println(constants.HelpSEORotate)
	fmt.Println(constants.HelpSEODryRunFlag)
	fmt.Println(constants.HelpSEOTemplateF)
	fmt.Println(constants.HelpSEOCreateTpl)
	fmt.Println(constants.HelpSEOAuthorName)
	fmt.Println(constants.HelpSEOAuthorEmail)
}

// printUsageAmendFlags prints the amend flags section.
func printUsageAmendFlags() {
	fmt.Println()
	fmt.Println(constants.HelpAmendFlags)
	fmt.Println(constants.HelpAmendName)
	fmt.Println(constants.HelpAmendEmail)
	fmt.Println(constants.HelpAmendBr)
	fmt.Println(constants.HelpAmendDry)
	fmt.Println(constants.HelpAmendForce)
}

// printUsageGoModFlags prints the gomod flags section.
func printUsageGoModFlags() {
	fmt.Println()
	fmt.Println(constants.HelpGoModFlags)
	fmt.Println(constants.HelpGoModDry)
	fmt.Println(constants.HelpGoModNoMrg)
	fmt.Println(constants.HelpGoModNoTdy)
	fmt.Println(constants.HelpGoModVerb)
	fmt.Println(constants.HelpGoModExt)
}
