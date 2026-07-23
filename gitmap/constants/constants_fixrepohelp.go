package constants

// Fix-repo flags help section. Rendered by cmd.printUsageFixRepoFlags
// as part of `gitmap help` so the -2/-3/-5/--all family is discoverable
// from the top-level helpline (not just from `gitmap help fix-repo`).
const (
	HelpFixRepoFlags     = "Fix-repo flags:"
	HelpFRMode2          = "  -2 (default)        Rewrite the last 2 prior versions (v(K-2)..v(K-1) -> vK)"
	HelpFRMode3          = "  -3                  Widen window to last 3 prior versions"
	HelpFRMode5          = "  -5                  Widen window to last 5 prior versions"
	HelpFRAll            = "  --all               Rewrite every prior version v1..v(K-1) -> vK"
	HelpFRDryRun         = "  --dry-run           Preview only; no file is written (PowerShell alias: -DryRun)"
	HelpFRVerbose        = "  --verbose           Print every modified file with replacement count (alias: -Verbose)"
	HelpFRConfig         = "  --config <path>     Override fix-repo.config.json location (alias: -Config <path>)"
	HelpFRStrict         = "  --strict            After rewrite + gofmt, run `go test` on every touched Go package; exit 9 on test failure (alias: -Strict)"
	HelpFRRestrict       = "  --restrict <mode>   Narrow rewrite scope. Values: no-version | nv. Suppresses the v1→v2 bare-base sweep so ONLY {base}-vN tokens are rewritten (alias: -r <mode>)"
	HelpFRExample1       = "      example:        gitmap fix-repo -2 --restrict no-version   # only gitmap-v27 -> gitmap-v27 (skip bare `gitmap`)"
	HelpFRExample2       = "      example:        gitmap fr -2 -r nv                          # short form, same behavior"
	HelpFixRepoExitCodes = "  exit codes:         0 ok | 2 not-a-repo | 3 no-remote | 4 no-vN-suffix | 5 bad-version | 6 bad-flag | 7 write-failed | 8 bad-config | 9 tests-failed"
)
