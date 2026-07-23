package constants

// commit-in (`commit-in` / `cin`) constants. See spec/03-commit-in/.
// Per Core memory rules: no magic strings, PascalCase, all CLI IDs in
// constants_cli.go (the command tokens themselves live there; this
// file owns flags, exit codes, enum names, messages, and paths).

// ---- Flag long names (kebab-case as parsed from argv) -------------
const (
	CommitInFlagDefault              = "default"
	CommitInFlagDefaultShort         = "d"
	CommitInFlagProfile              = "profile"
	CommitInFlagSaveProfile          = "save-profile"
	CommitInFlagSaveProfileOverwrite = "save-profile-overwrite"
	CommitInFlagSetDefault           = "set-default"
	CommitInFlagAuthorName           = "author-name"
	CommitInFlagAuthorEmail          = "author-email"
	CommitInFlagConflict             = "conflict"
	CommitInFlagExclude              = "exclude"
	CommitInFlagMessageExclude       = "message-exclude"
	CommitInFlagMessagePrefix        = "message-prefix"
	CommitInFlagMessageSuffix        = "message-suffix"
	CommitInFlagTitlePrefix          = "title-prefix"
	CommitInFlagTitleSuffix          = "title-suffix"
	CommitInFlagOverrideMessages     = "override-messages"
	CommitInFlagOverrideOnlyWeak     = "override-only-weak"
	CommitInFlagWeakWords            = "weak-words"
	CommitInFlagFunctionIntel        = "function-intel"
	CommitInFlagLanguages            = "languages"
	CommitInFlagNoPrompt             = "no-prompt"
	CommitInFlagDryRun               = "dry-run"
	CommitInFlagKeepTemp             = "keep-temp"
	CommitInFlagNoReleaseBranch      = "no-release-branch"
)

// CommitInDescNoReleaseBranch — flag default is OFF (i.e. release
// branches ARE created). Setting the flag suppresses creation even
// when a mirrored tag matches VersionTagPattern. See spec §08 §2.5
// + spec §09 §9.4 (`MirroredReleaseBranch` NULL contract).
const CommitInDescNoReleaseBranch = "Suppress auto release-branch creation for version tags"

// ---- Flag descriptions (rendered in --help) -----------------------
const (
	CommitInDescDefault              = "Load the default profile bound to <source>"
	CommitInDescProfile              = "Load a named profile from .gitmap/commit-in/profiles/<name>.json"
	CommitInDescSaveProfile          = "Persist this run's resolved answers as a named profile"
	CommitInDescSaveProfileOverwrite = "Allow --save-profile to overwrite an existing same-named profile"
	CommitInDescSetDefault           = "Mark the saved profile as the default for <source>"
	CommitInDescAuthorName           = "Override author identity (Name); requires --author-email"
	CommitInDescAuthorEmail          = "Override author identity (Email); requires --author-name"
	CommitInDescConflict             = "Conflict mode when replaying commits: ForceMerge | Prompt"
	CommitInDescExclude              = "CSV of relative paths to exclude (trailing / => folder)"
	CommitInDescMessageExclude       = "CSV of Kind:Value rules (StartsWith:|EndsWith:|Contains:)"
	CommitInDescMessagePrefix        = "CSV pool prepended to every commit body (random pick per commit)"
	CommitInDescMessageSuffix        = "CSV pool appended to every commit body (random pick per commit)"
	CommitInDescTitlePrefix          = "String prepended to the first line of every commit message"
	CommitInDescTitleSuffix          = "String appended to the first line of every commit message"
	CommitInDescOverrideMessages     = "CSV pool that replaces the entire message (random pick per commit)"
	CommitInDescOverrideOnlyWeak     = "Only apply --override-messages when the title's first word is weak"
	CommitInDescWeakWords            = "CSV of weak first-words triggering override (default: change,update,updates)"
	CommitInDescFunctionIntel        = "Toggle per-language new-function detection block: on | off"
	CommitInDescLanguages            = "CSV of languages to scan when --function-intel on"
	CommitInDescNoPrompt             = "Refuse interactive prompts; exit MissingAnswer when a value is unset"
	CommitInDescDryRun               = "Walk and plan only; never run git commit (DB rows logged as Skipped)"
	CommitInDescKeepTemp             = "Don't delete <.gitmap>/temp/<runId>/ on exit"
)

// ---- Exit codes (spec §2.7) ---------------------------------------
const (
	CommitInExitOk              = 0
	CommitInExitPartiallyFailed = 1
	CommitInExitBadArgs         = 2
	CommitInExitSourceUnusable  = 3
	CommitInExitInputUnusable   = 4
	CommitInExitDbFailed        = 5
	CommitInExitProfileMissing  = 6
	CommitInExitMissingAnswer   = 7
	CommitInExitConflictAborted = 8
	CommitInExitLockBusy        = 9
	CommitInExitFunctionIntel   = 10
)

// ---- Enum: ConflictMode -------------------------------------------
const (
	CommitInConflictModeForceMerge = "ForceMerge"
	CommitInConflictModePrompt     = "Prompt"
)

// ---- Enum: InputKeyword (spec §2.4) -------------------------------
const (
	CommitInInputKeywordAll      = "all"
	CommitInInputKeywordTailDash = "-" // followed by digits, e.g. "-5"
)

// ---- Enum: InputKind ----------------------------------------------
const (
	CommitInInputKindLocalFolder      = "LocalFolder"
	CommitInInputKindGitUrl           = "GitUrl"
	CommitInInputKindVersionedSibling = "VersionedSibling"
)

// ---- Enum: RunStatus ----------------------------------------------
const (
	CommitInRunStatusPending         = "Pending"
	CommitInRunStatusRunning         = "Running"
	CommitInRunStatusCompleted       = "Completed"
	CommitInRunStatusFailed          = "Failed"
	CommitInRunStatusPartiallyFailed = "PartiallyFailed"
)

// ---- Enum: CommitOutcome ------------------------------------------
const (
	CommitInOutcomeCreated = "Created"
	CommitInOutcomeSkipped = "Skipped"
	CommitInOutcomeFailed  = "Failed"
)

// ---- Enum: SkipReason ---------------------------------------------
const (
	CommitInSkipReasonDuplicateSourceSha     = "DuplicateSourceSha"
	CommitInSkipReasonExcludedAllFiles       = "ExcludedAllFiles"
	CommitInSkipReasonEmptyAfterMessageRules = "EmptyAfterMessageRules"
	CommitInSkipReasonDryRun                 = "DryRun"
)

// ---- Enum: ExclusionKind ------------------------------------------
const (
	CommitInExclusionKindPathFolder = "PathFolder"
	CommitInExclusionKindPathFile   = "PathFile"
)

// ---- Enum: MessageRuleKind ----------------------------------------
const (
	CommitInMessageRuleKindStartsWith = "StartsWith"
	CommitInMessageRuleKindEndsWith   = "EndsWith"
	CommitInMessageRuleKindContains   = "Contains"
)

// ---- Enum: FunctionIntelLanguage ----------------------------------
const (
	CommitInLanguageGo         = "Go"
	CommitInLanguageJavaScript = "JavaScript"
	CommitInLanguageTypeScript = "TypeScript"
	CommitInLanguageRust       = "Rust"
	CommitInLanguagePython     = "Python"
	CommitInLanguagePhp        = "Php"
	CommitInLanguageJava       = "Java"
	CommitInLanguageCSharp     = "CSharp"
)

// ---- Defaults ------------------------------------------------------
const (
	CommitInDefaultConflictMode  = CommitInConflictModeForceMerge
	CommitInDefaultFunctionIntel = "off"
	CommitInDefaultLanguagesCsv  = CommitInLanguageGo
	CommitInDefaultWeakWordsCsv  = "change,update,updates"
	CommitInFunctionIntelOn      = "on"
	CommitInFunctionIntelOff     = "off"
	CommitInMessageRuleKindSep   = ":"
	CommitInPathFolderSep        = "/"
	CommitInCsvSep               = ","
)

// ---- Filesystem layout (under <.gitmap>/) -------------------------
const (
	CommitInDirRoot         = "commit-in"
	CommitInDirProfiles     = "commit-in/profiles"
	CommitInDirTemp         = "temp"
	CommitInProfileFileExt  = ".json"
	CommitInLockFileName    = "commit-in.lock"
	CommitInTempInputFormat = "%d-%s" // <orderIndex>-<basename>
)

// ---- URL detection (spec §2.3 step 1) -----------------------------
const (
	CommitInUrlPrefixHttps = "https://"
	CommitInUrlPrefixHttp  = "http://"
	CommitInUrlPrefixSshAt = "git@"
	CommitInUrlPrefixSsh   = "ssh://"
	CommitInUrlPrefixGit   = "git://"
	CommitInUrlSuffixGit   = ".git"
)

// ---- Phase banners (STDERR) ---------------------------------------
const (
	CommitInMsgPhaseResolveSource = "▸ commit-in: resolving <source> %q\n"
	CommitInMsgPhaseInitSource    = "▸ commit-in: git init in %s\n"
	CommitInMsgPhaseCloneSource   = "▸ commit-in: cloning %s into %s\n"
	CommitInMsgPhaseDiscoverAll   = "▸ commit-in: discovering versioned siblings of %s\n"
	CommitInMsgPhaseStageInputs   = "▸ commit-in: staging %d input(s) under %s\n"
	CommitInMsgPhaseWalk          = "▸ commit-in: walking %d commit(s) chronologically\n"
	CommitInMsgPhaseReplay        = "▸ commit-in: replaying into %s\n"
	CommitInMsgPhaseSaveProfile   = "▸ commit-in: saving profile %q\n"
	CommitInMsgCommitOk           = "✓ commit-in: %s -> %s  %s\n"
	CommitInMsgCommitSkip         = "▣ commit-in: skip %s  reason=%s\n"
	CommitInMsgCommitFail         = "✗ commit-in: %s  %v\n"
	CommitInMsgSummaryLine        = "commit-in: run=%d created=%d skipped=%d failed=%d\n"
)

// ---- Errors (one line per error path; zero-swallow rule) ----------
const (
	CommitInErrBadArgs           = "commit-in: bad args: %s\n"
	CommitInErrSourceClone       = "commit-in: source: clone failed: %v\n"
	CommitInErrSourceInit        = "commit-in: source: git init failed: %v\n"
	CommitInErrSourceMkdir       = "commit-in: source: mkdir failed: %v\n"
	CommitInErrInputClone        = "commit-in: input %q: clone failed: %v\n"
	CommitInErrInputOpen         = "commit-in: input %q: open failed: %v\n"
	CommitInErrInputMixedKeyword = "commit-in: keyword %q must appear alone"
	CommitInErrDbMigrate         = "commit-in: db: migration failed: %v\n"
	CommitInErrDbWrite           = "commit-in: db: write failed: %v\n"
	CommitInErrProfileMissing    = "commit-in: profile %q not found\n"
	CommitInErrProfileBadJson    = "commit-in: profile %q: invalid json: %v\n"
	CommitInErrMissingAnswer     = "commit-in: --no-prompt set but %s is unset\n"
	CommitInErrConflictAborted   = "commit-in: conflict aborted by user at commit %s\n"
	CommitInErrLockBusy          = "commit-in: another commit-in run is in progress (lock at %s)\n"
	CommitInErrFunctionIntel     = "commit-in: function-intel: %s parser failed on %s: %v\n"
	CommitInErrAuthorPair        = "commit-in: --author-name and --author-email must be set together"
	CommitInErrConflictMode      = "commit-in: --conflict must be ForceMerge or Prompt, got %q"
	CommitInErrFunctionIntelArg  = "commit-in: --function-intel must be on or off, got %q"
	CommitInErrUnknownLanguage   = "commit-in: unknown language %q (supported: Go,JavaScript,TypeScript,Rust,Python,Php,Java,CSharp)"
	CommitInErrMessageRuleShape  = "commit-in: message rule %q: expected Kind:Value with Kind in StartsWith|EndsWith|Contains"
	CommitInErrSaveProfileExists = "commit-in: profile %q already exists; pass --save-profile-overwrite to replace"
	CommitInErrInputKeyword      = "commit-in: keyword %q invalid (use 'all' or '-N' where N>=1)"
)
