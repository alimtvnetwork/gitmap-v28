// Data tables for the CommitIn docs page. Extracted so the page
// component itself stays under the project-wide <200-lines rule.

export const commitInFlags = [
  { flag: "-d, --default", def: "off", desc: "Load the default profile bound to <source>" },
  { flag: "--profile <name>", def: "—", desc: "Load .gitmap/commit-in/profiles/<name>.json" },
  { flag: "--save-profile <name>", def: "—", desc: "Persist this run's resolved settings as a profile" },
  { flag: "--save-profile-overwrite", def: "off", desc: "Allow --save-profile to overwrite" },
  { flag: "--set-default", def: "off", desc: "Mark the saved profile as default for <source>" },
  { flag: "--author-name <s>", def: "—", desc: "Override author name (requires --author-email)" },
  { flag: "--author-email <s>", def: "—", desc: "Override author email (requires --author-name)" },
  { flag: "--conflict <mode>", def: "ForceMerge", desc: "ForceMerge or Prompt" },
  { flag: "--exclude <csv>", def: "—", desc: "Per-commit exclude list (trailing / = folder)" },
  { flag: "--message-exclude <csv>", def: "—", desc: "Kind:Value rules: StartsWith: / EndsWith: / Contains:" },
  { flag: "--message-prefix <csv>", def: "—", desc: "Random-pick pool prepended to every body" },
  { flag: "--message-suffix <csv>", def: "—", desc: "Random-pick pool appended to every body" },
  { flag: "--title-prefix <s>", def: "—", desc: "Prepended to the FIRST line only" },
  { flag: "--title-suffix <s>", def: "—", desc: "Appended to the FIRST line only" },
  { flag: "--override-messages <csv>", def: "—", desc: "Replaces the entire message (random pick)" },
  { flag: "--override-only-weak", def: "off", desc: "Override only when the title's first word is weak" },
  { flag: "--weak-words <csv>", def: "change,update,updates", desc: "First-word triggers for override" },
  { flag: "--function-intel on|off", def: "off", desc: "Append per-language new-function block" },
  { flag: "--languages <csv>", def: "Go", desc: "Languages scanned when intel is on" },
  { flag: "--tags <mode>", def: "Annotated", desc: "Mirror source tags: Annotated | All | None" },
  { flag: "--no-release-branch", def: "off", desc: "Suppress auto release/<tag> branch for semver tags" },
  { flag: "--release-branch-prefix <s>", def: "release/", desc: "Override the auto release-branch prefix (must end with /)" },
  { flag: "--no-prompt", def: "off", desc: "Refuse interactive prompts; exit MissingAnswer if unset" },
  { flag: "--dry-run", def: "off", desc: "Plan only; never run git commit" },
  { flag: "--keep-temp", def: "off", desc: "Keep .gitmap/temp/<runId>/ after exit" },
];

export const commitInExitCodes = [
  { code: "0", meaning: "Ok — every walked commit was Created or Skipped" },
  { code: "1", meaning: "PartiallyFailed — at least one commit failed but others succeeded" },
  { code: "2", meaning: "BadArgs — flag / positional validation failed" },
  { code: "3", meaning: "SourceUnusable — <source> could not be resolved or initialized" },
  { code: "4", meaning: "InputUnusable — at least one input could not be cloned / opened" },
  { code: "5", meaning: "DbFailed — SQLite migration or write failed" },
  { code: "6", meaning: "ProfileMissing — --profile / --default lookup empty" },
  { code: "7", meaning: "MissingAnswer — --no-prompt set but a required value was unset" },
  { code: "8", meaning: "ConflictAborted — Prompt mode and the user aborted the merge" },
  { code: "9", meaning: "LockBusy — another commit-in run holds the workspace lock" },
  { code: "10", meaning: "FunctionIntel — a per-language detector panicked" },
];

export const commitInAutoInit = [
  { when: "An https:// or git@ URL", then: "git clone <url> into the derived folder name" },
  { when: "An existing path with .git/", then: "Reuse the repo in place — never re-init" },
  { when: "An existing folder, NO .git/", then: "git init in place (your files are kept untouched)" },
  { when: "A path that does not exist", then: "mkdir -p <path> && git init <path>" },
];

export const commitInProfileJson = `{
  "Name": "Default",
  "SchemaVersion": 1,
  "SourceRepoPath": "/abs/path/to/canonical",
  "IsDefault": true,
  "ConflictMode": "ForceMerge",
  "Author": {
    "Name": "Jane Doe",
    "Email": "jane@example.com"
  },
  "Exclusions": [
    { "Kind": "PathFolder", "Value": "node_modules" },
    { "Kind": "PathFolder", "Value": "dist" },
    { "Kind": "PathFile",   "Value": "secrets.env" }
  ],
  "MessageRules": [
    { "Kind": "StartsWith", "Value": "Signed-off-by:" },
    { "Kind": "Contains",   "Value": "[skip ci]" },
    { "Kind": "EndsWith",   "Value": "(cherry picked from commit)" }
  ],
  "MessagePrefix":   ["chore:", "feat:", "fix:"],
  "MessageSuffix":   [],
  "TitlePrefix":     "",
  "TitleSuffix":     " — via gitmap-v28",
  "OverrideMessages": ["Improve module", "Refine implementation"],
  "OverrideOnlyWeak": true,
  "WeakWords":        ["change", "update", "updates", "misc"],
  "FunctionIntel": {
    "IsEnabled": true,
    "Languages": ["Go", "TypeScript", "Python"]
  }
}`;
