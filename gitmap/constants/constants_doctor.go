package constants

// Doctor command messages.
const (
	DoctorBannerFmt       = "\n  gitmap doctor (v%s)\n"
	DoctorBannerRule      = "  ──────────────────────────────────────────"
	DoctorIssuesFmt       = "  Found %d issue(s). See recommendations above.\n"
	DoctorFixPathTip      = "  Tip: run 'gitmap doctor --fix-path' to auto-sync the PATH binary.\n\n"
	DoctorAllPassed       = "  All checks passed."
	DoctorFixBannerFmt    = "\n  gitmap doctor --fix-path (v%s)\n"
	DoctorActivePathFmt   = "  Active PATH:  %s (%s)\n"
	DoctorDeployedFmt     = "  Deployed:     %s (%s)\n"
	DoctorSyncingFmt      = "  Syncing %s -> %s...\n"
	DoctorRetryFmt        = "  [%d/%d] File in use, retrying...\n"
	DoctorRenamedMsg      = "  Renamed locked binary to .old, copying fresh..."
	DoctorKillingMsg      = "  Attempting to stop stale gitmap processes..."
	DoctorKilledFmt       = "  Stopped process(es): %s\n"
	DoctorSyncFailTitle   = "Could not sync PATH binary after all fallback attempts"
	DoctorSyncFailDetail  = "The file is still locked by another process."
	DoctorSyncFailFix1    = "Close all terminals and apps using gitmap, then run:"
	DoctorSyncFailFix2Fmt = "  Copy-Item \"%s\" \"%s\" -Force"
	DoctorFixFlagDesc     = "Sync the active PATH binary from the deployed binary"
	DoctorOKPathFmt       = "PATH binary synced successfully: %s"
	DoctorWarnSyncFmt     = "Synced but version mismatch: got %s, expected %s"
	DoctorNotOnPath       = "gitmap not found on PATH"
	DoctorNoSync          = "Cannot sync — no active binary to replace."
	DoctorAddPathFix      = "Add your deploy directory to PATH first."
	DoctorCannotResolve   = "Cannot resolve deployed binary"
	DoctorAlreadySynced   = "PATH already points to the deployed binary. Nothing to sync."
	DoctorVersionsMatch   = "Versions already match (%s). No sync needed."
	DoctorRepoPathMissing = "RepoPath not embedded"
	DoctorRepoPathDetail  = "Binary was not built with run.ps1. Self-update will not work."
	DoctorRepoPathFix     = "Rebuild with: .\\run.ps1"
	DoctorRepoPathOKFmt   = "RepoPath: %s"
	DoctorPathBinaryFmt   = "PATH binary: %s (%s)"
	DoctorPathMissTitle   = "gitmap not found on PATH"
	DoctorPathMissDetail  = "The gitmap binary is not accessible from your terminal."
	DoctorPathMissFix     = "Add your deploy directory to PATH (check deployPath in powershell.json or run 'gitmap installed-dir')"
	DoctorDeployReadFail  = "Cannot read powershell.json"
	DoctorDeployReadDet   = "Deploy path detection unavailable."
	DoctorNoDeployPath    = "No deployPath in powershell.json"
	DoctorNoDeployDet     = "Deploy target not configured."
	DoctorDeployNotFound  = "Deployed binary not found"
	DoctorDeployRunFix    = "Run: .\\run.ps1"
	DoctorDeployOKFmt     = "Deployed binary: %s (%s)"
	DoctorGitMissTitle    = "git not found on PATH"
	DoctorGitMissDetail   = "Git is required for most gitmap commands."
	DoctorGitOKFmt        = "Git: %s (%s)"
	DoctorGitOKPathFmt    = "Git: %s (version unknown)"
	DoctorGoWarn          = "Go not found on PATH (needed only for building from source)"
	DoctorGoOKFmt         = "Go: %s"
	DoctorGoOKPathFmt     = "Go: %s (version unknown)"
	DoctorChangelogWarn   = "CHANGELOG.md not found (changelog command will not work)"
	DoctorChangelogOK     = "CHANGELOG.md present"
	DoctorVersionMismatch = "PATH binary version mismatch"
	DoctorVMismatchFmt    = "PATH: %s, Source: %s"
	DoctorVMismatchFix    = "Run: gitmap update  or  gitmap doctor --fix-path"
	DoctorDeployMismatch  = "Deployed binary version mismatch"
	DoctorDMismatchFmt    = "Deployed: %s, Source: %s"
	DoctorDMismatchFix    = "Run: .\\run.ps1 -NoPull"
	DoctorBinariesDiffer  = "PATH and deployed binaries differ"
	DoctorBDifferFmt      = "PATH: %s (%s), Deployed: %s (%s)"
	DoctorBDifferFix      = "Run: gitmap doctor --fix-path"
	DoctorSourceOKFmt     = "Source version: %s (all binaries match)"
	DoctorResolveNoRepo   = "RepoPath not embedded — rebuild with run.ps1"
	DoctorResolveNoRead   = "cannot read powershell.json: %v"
	DoctorResolveNoDeploy = "no deployPath in powershell.json"
	DoctorResolveNotFound = "deployed binary not found: %s"
	DoctorDefaultBinary   = "gitmap.exe"
)

// Doctor binary and tool lookup names.
const (
	GitMapBin            = "gitmap"
	GoBin                = "go"
	GoVersionArg         = "version"
	PowershellConfigFile = "powershell.json"
	JSONKeyDeployPath    = "deployPath"
	JSONKeyBuildOutput   = "buildOutput"
	JSONKeyBinaryName    = "binaryName"
	BackupSuffix         = ".old"
)

// GitMapSubdir and GitMapCliSubdir are now defined in deploy_manifest.go and
// populated from the embedded deploy-manifest.json (single source of truth).

// Bare-invocation banner suppression (v3.6.0+).
const (
	FlagNoBanner       = "--no-banner"
	EnvGitMapQuiet     = "GITMAP_QUIET"
	EnvGitMapQuietTrue = "1"
)

// Bare-invocation binary readout labels (v3.6.0+).
const (
	BinaryReadoutActive   = "  Active binary:    %s\n"
	BinaryReadoutDeployed = "  Deployed binary:  %s\n"
	BinaryReadoutConfig   = "  Config binary:    %s\n"
	BinaryReadoutMissing  = "(not found)"
)

// Doctor format markers.
const (
	DoctorOKFmt    = "  %s[OK]%s %s\n"
	DoctorIssueFmt = "  %s[!!]%s %s\n"
	DoctorFixFmt   = "       %sFix:%s %s\n"
	DoctorWarnFmt  = "  %s[--]%s %s\n"
	DoctorDetail   = "       %s\n"
)

// Doctor config validation messages.
const (
	DoctorConfigMissing = "config.json not found (using defaults)"
	DoctorConfigInvalid = "config.json is not valid JSON"
	DoctorConfigOKFmt   = "Config: %s"
)

// Doctor database validation messages.
const (
	DoctorDBOpenFail    = "Database cannot be opened"
	DoctorDBMigrateFail = "Database migration failed"
	DoctorDBOK          = "Database: %s"
)

// Doctor lock file messages.
const (
	DoctorLockNone   = "No stale lock file"
	DoctorLockExists = "Lock file exists — another gitmap may be running (or stale)"
)

// Doctor network messages.
const (
	DoctorNetworkOK      = "Network: github.com reachable"
	DoctorNetworkOffline = "Network: github.com unreachable (offline mode)"
)

// Doctor digital signature messages.
const (
	DoctorSignTitle      = "Digital signature"
	DoctorSignOKFmt      = "Signed: %s (%s)"
	DoctorSignUnsigned   = "Binary is not digitally signed"
	DoctorSignUnsignDet  = "Users will see SmartScreen warnings on Windows."
	DoctorSignUnsignFix  = "See spec/03-general/05-code-signing.md for SignPath.io setup"
	DoctorSignSkipUnix   = "Signature check: skipped (not Windows)"
	DoctorSignNoPath     = "Signature check: skipped (binary not on PATH)"
	DoctorSignCheckFail  = "Signature check: could not verify (PowerShell unavailable)"
	DoctorSignInvalidFmt = "Signature invalid: %s"
	DoctorSignInvalidDet = "The binary's signature is present but not valid."
)

// Doctor legacy directory messages.
const (
	DoctorLegacyDirsOK = "No legacy directories (.release/, gitmap-output/, .deployed/)"
)

// Doctor duplicate binary messages.
const (
	DoctorDupBinOK    = "Single gitmap binary on PATH"
	DoctorDupBinTitle = "Multiple gitmap binaries on PATH"
)

// Doctor Release↔Repo integrity messages.
const (
	DoctorNoOrphans        = "Release↔Repo integrity: no orphaned rows"
	DoctorOrphanedReleases = "%d orphaned Release row(s) with invalid RepoId"
	DoctorOrphanedDetail   = "These releases reference a Repo that no longer exists in the database."
	DoctorOrphanedFix      = "Run: gitmap list-releases (re-imports from .gitmap/release/v*.json with valid RepoId)"
	DoctorReposNoReleases  = "%d repo(s) have no releases (run 'gitmap list-releases' in each repo to populate)"
	DoctorIntegrityFail    = "Release↔Repo integrity check failed: %v"
)

// Doctor setup config messages.
const (
	DoctorSetupConfigMissing = "git-setup.json not found (setup will fail without --config)"
	DoctorSetupConfigOKFmt   = "Setup config: %s"
)

// Doctor command wrapper messages.
const (
	DoctorWrapperOK        = "Command wrapper active (gitmap cd will change directory)"
	DoctorWrapperNotLoaded = "Command wrapper not loaded — gitmap cd prints path but cannot change directory"
	DoctorWrapperFix       = "Run: gitmap setup, then restart terminal or reload profile (. $PROFILE / source ~/.bashrc / source ~/.zshrc)"
)

// Doctor VS Code Project Manager check messages (v3.41.0+).
const (
	DoctorVSCodePMOKFmt        = "VS Code Project Manager: %s"
	DoctorVSCodePMNoVSCode     = "VS Code user-data dir not found — projects.json sync will be skipped (install VS Code or set APPDATA / HOME / XDG_CONFIG_HOME)"
	DoctorVSCodePMNoExtension  = "alefragnani.project-manager extension not installed — projects.json sync will be skipped (install the extension and re-run)"
	DoctorVSCodePMUnknownTitle = "VS Code Project Manager check failed"
)
