package constants

// Fix-repo flags help section. Rendered by cmd.printUsageFixRepoFlags
// as part of `gitmap help` so the -2/-3/-5/--all family is discoverable
// from the top-level helpline (not just from `gitmap help fix-repo`).
const (
	HelpFixRepoFlags     = "Fix-repo flags:"
	HelpFRMode2          = "  -2 (default)        Rewrite last 2 prior versions (v(K-2)..v(K-1) -> vK)"
	HelpFRMode3          = "  -3, -5              Widen rewrite window to last 3 or 5 prior versions"
	HelpFRMode5          = "  " + ColorCyan + ColorCyan + "--all" + ColorReset + ColorReset + "               Rewrite every prior version v1..v(K-1) -> vK"
	HelpFRAll            = "  " + ColorCyan + ColorCyan + "--dry-run" + ColorReset + ColorReset + "           Preview changes without writing files (-DryRun)"
	HelpFRDryRun         = "  " + ColorCyan + ColorCyan + "--strict" + ColorReset + ColorReset + "            Run `go test` on touched packages; exit 9 on fail"
	HelpFRVerbose        = "  " + ColorCyan + ColorCyan + "--verbose" + ColorReset + ColorReset + "           Log each modified file and replacement count"
	HelpFRConfig         = "  " + ColorCyan + ColorCyan + "--config" + ColorReset + ColorReset + " <path>     Custom fix-repo.config.json path"
	HelpFRStrict         = "  " + ColorCyan + ColorCyan + "--restrict" + ColorReset + ColorReset + " <mode>   Narrow scope: no-version (nv) avoids bare-base"
	HelpFRRestrict       = "      example:        gitmap fix-repo -2 " + ColorCyan + ColorCyan + "--restrict" + ColorReset + ColorReset + " nv"
	HelpFRExample1       = "      example:        gitmap fr -3 " + ColorCyan + ColorCyan + "--dry-run" + ColorReset + ColorReset
	HelpFRExample2       = "  " + ColorCyan + ColorCyan + "--gofmt-max-cmd-len" + ColorReset + ColorReset + "  Cap gofmt batch argv length (default: 30000)"
	HelpFRGofmtMaxCmdLen = "  exit codes:         0 ok | 2 not-repo | 4 no-vN | 6 bad-flag | 9 test-fail"
	HelpFixRepoExitCodes = ""
)
