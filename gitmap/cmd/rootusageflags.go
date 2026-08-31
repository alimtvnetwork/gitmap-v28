package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// printGroupUtilities prints the utility commands.
func printGroupUtilities() {
	renderHeader(constants.HelpGroupUtilities)
	renderLine(constants.HelpDoctor)
	renderLine(constants.HelpUpdate)
	renderLine(constants.HelpUpdateCleanup)
	renderLine(constants.HelpVersion)
	renderLine(constants.HelpCompletion)
	renderLine(constants.HelpInteractive)
	renderLine(constants.HelpDocs)
	renderLine(constants.HelpHelpDash)
	renderLine(constants.HelpGoMod)
	renderLine(constants.HelpServe)
	renderLine(constants.HelpSEOWrite)
	renderLine(constants.HelpLLMDocs)
	renderLine(constants.HelpFixRepo)
	renderLine(constants.HelpMakePublic)
	renderLine(constants.HelpMakePrivate)
	renderLine(constants.HelpCloneFixRepo)
	renderLine(constants.HelpCloneFixRepoPub)
	renderLine(constants.HelpCmdOpen)
	renderLine(constants.HelpHelp)
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
	printUsageServeFlags()
}

// printUsageServeFlags prints the serve flags section.
func printUsageServeFlags() {
	renderFlagHeader("Serve flags:")
	renderLine("  -port <port>        " + constants.FlagDescServePort)
}

// printUsageFixRepoFlags prints the fix-repo flags section so the
// -2 / -3 / -5 / --all / --dry-run family is discoverable from the
// top-level `gitmap help` output (not only `gitmap help fix-repo`).
func printUsageFixRepoFlags() {
	renderFlagHeader(constants.HelpFixRepoFlags)
	renderLine(constants.HelpFRMode2)
	renderLine(constants.HelpFRMode3)
	renderLine(constants.HelpFRMode5)
	renderLine(constants.HelpFRAll)
	renderLine(constants.HelpFRDryRun)
	renderLine(constants.HelpFRVerbose)
	renderLine(constants.HelpFRConfig)
	renderLine(constants.HelpFRStrict)
	renderLine(constants.HelpFRRestrict)
	renderLine(constants.HelpFRExample1)
	renderLine(constants.HelpFRExample2)
	renderLine(constants.HelpFRGofmtMaxCmdLen)
}

// printUsageCloneNextFlags prints the clone-next flags section.
func printUsageCloneNextFlags() {
	renderFlagHeader(constants.HelpCloneNextFlags)
	renderLine(constants.HelpCNDelete)
	renderLine(constants.HelpCNKeep)
	renderLine(constants.HelpCNNoDesktop)
	renderLine(constants.HelpCNSSHKey)
	renderLine(constants.HelpCNVerbose)
	renderLine(constants.HelpCNCreateRemote)
}

// printUsageInteractiveFlags prints the interactive flags section.
func printUsageInteractiveFlags() {
	renderFlagHeader(constants.HelpInteractiveFlags)
	renderLine(constants.HelpRefresh)
}

// printUsageScanFlags prints the scan flags section.
func printUsageScanFlags() {
	renderFlagHeader(constants.HelpScanFlags)
	renderLine(constants.HelpConfig)
	renderLine(constants.HelpMode)
	renderLine(constants.HelpOutput)
	renderLine(constants.HelpOutputPath)
	renderLine(constants.HelpOutFile)
	renderLine(constants.HelpScanFlagGitHubDesktop)
	renderLine(constants.HelpOpen)
	renderLine(constants.HelpQuiet)
}

// printUsageCloneFlags prints the clone flags section.
func printUsageCloneFlags() {
	renderFlagHeader(constants.HelpCloneFlags)
	renderLine(constants.HelpTargetDir)
	renderLine(constants.HelpSafePull)
	renderLine(constants.HelpVerbose)
}

// printUsageReleaseFlags prints the release flags section.
func printUsageReleaseFlags() {
	renderFlagHeader(constants.HelpReleaseFlags)
	renderLine(constants.HelpAssets)
	renderLine(constants.HelpCommit)
	renderLine(constants.HelpRelBranch)
	renderLine(constants.HelpBump)
	renderLine(constants.HelpDraft)
	renderLine(constants.HelpDryRun)
	renderLine(constants.HelpCompressFlag)
	renderLine(constants.HelpChecksumsFlag)
	renderLine(constants.HelpBin)
	renderLine(constants.HelpTargets)
	renderLine(constants.HelpListTargets)
}

// printUsageSEOFlags prints the seo-write flags section.
func printUsageSEOFlags() {
	renderFlagHeader(constants.HelpSEOWriteFlags)
	renderLine(constants.HelpSEOCSV)
	renderLine(constants.HelpSEOURL)
	renderLine(constants.HelpSEOService)
	renderLine(constants.HelpSEOArea)
	renderLine(constants.HelpSEOCompany)
	renderLine(constants.HelpSEOPhone)
	renderLine(constants.HelpSEOEmail)
	renderLine(constants.HelpSEOAddress)
	renderLine(constants.HelpSEOMaxCommits)
	renderLine(constants.HelpSEOInterval)
	renderLine(constants.HelpSEOFilesFlag)
	renderLine(constants.HelpSEORotate)
	renderLine(constants.HelpSEODryRunFlag)
	renderLine(constants.HelpSEOTemplateF)
	renderLine(constants.HelpSEOCreateTpl)
	renderLine(constants.HelpSEOAuthorName)
	renderLine(constants.HelpSEOAuthorEmail)
}

// printUsageAmendFlags prints the amend flags section.
func printUsageAmendFlags() {
	renderFlagHeader(constants.HelpAmendFlags)
	renderLine(constants.HelpAmendName)
	renderLine(constants.HelpAmendEmail)
	renderLine(constants.HelpAmendBr)
	renderLine(constants.HelpAmendDry)
	renderLine(constants.HelpAmendForce)
}

// printUsageGoModFlags prints the gomod flags section.
func printUsageGoModFlags() {
	renderFlagHeader(constants.HelpGoModFlags)
	renderLine(constants.HelpGoModDry)
	renderLine(constants.HelpGoModNoMrg)
	renderLine(constants.HelpGoModNoTdy)
	renderLine(constants.HelpGoModVerb)
	renderLine(constants.HelpGoModExt)
}
