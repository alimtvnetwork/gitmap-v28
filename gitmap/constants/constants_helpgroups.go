package constants

// Help group headers.
const (
	HelpGroupScanning     = "  Scanning & Discovery:"
	HelpGroupCloning      = "  Cloning & Sync:"
	HelpGroupGitOps       = "  Git Operations:"
	HelpGroupNavigation   = "  Navigation & Organization:"
	HelpGroupRelease      = "  Release & Versioning:"
	HelpGroupReleaseInfo  = "  Release History & Info:"
	HelpGroupData         = "  Data, Profiles & Bookmarks:"
	HelpGroupImportExport = "  Import / Export:"
	HelpGroupHistory      = "  History & Stats:"
	HelpGroupAmendGroup   = "  Author Amendment:"
	HelpGroupProject      = "  Project Detection:"
	HelpGroupSSH          = "  SSH Key Management:"
	HelpGroupZip          = "  Zip Groups (Release Archives):"
	HelpGroupEnvTools     = "  Environment & Tools:"
	HelpGroupTasks        = "  File-Sync Tasks:"
	HelpGroupUtilities    = "  Utilities:"
	HelpGroupVisualize    = "  Visualization:"
	HelpGroupCommitXfer   = "  Commit Transfer (replay between repos):"
	HelpGroupChromeProf   = "  Chrome Profile (copy / export / import / list / delete):"
	HelpGroupTemplates    = "  Templates & Scaffolding (.gitignore / .gitattributes / LFS):"
	HelpGroupCluster      = "  Cluster & Delegation (multi-machine networks):"
	HelpGroupInstallers   = "  Installers & Macros:"
	HelpGroupIntegrations = "  Integrations (VS Code, Antigravity, Scheduler):"

	HelpAddIgnore     = "  add ignore [langs...]      Merge curated .gitignore block into ./.gitignore (idempotent, marker-block aware)"
	HelpAddAttributes = "  add attributes [langs...]  Merge curated .gitattributes block into ./.gitattributes (idempotent, marker-block aware)"
	HelpAddLFSInstall = "  add lfs-install            Run 'git lfs install --local' and merge the lfs/common .gitattributes block"
	HelpTemplatesInit = "  templates init (tpl ti)    Scaffold .gitignore + .gitattributes for one or more languages [--lfs] [--dry-run] [--force]"
	HelpTemplatesList = "  templates list (tpl tl)    List every available template (kind, lang, source: user/embed, path)"
	HelpTemplatesShow = "  templates show (tpl ts)    Print one template (overlay > embed) to stdout, audit-trail header included"
	HelpTemplatesDiff = "  templates diff (tpl td)    Preview what add ignore/add attributes would change; exit codes mirror diff(1)"
	HelpSync          = "  sync (sy) <target>         Union-merge curated defaults: ignore | attributes | lfs-install | prettier-ignore | prettier-rc | all  [--dry-run] [--force]"
	HelpCommons       = "  commons (co)               Shortcut for 'sync all' — add/dedupe curated .gitignore, .gitattributes, .prettierignore, .prettierrc + git lfs install  [--dry-run]"

	HelpServersClients = "  servers-clients (sc) <sub>       Broadcast commands across all server + client nodes"
	HelpClients        = "  clients <sub>              Broadcast commands across client nodes only"
	HelpCluster        = "  cluster <sub>              Manage cluster nodes, history, exports, passwords"
	HelpExecute        = "  execute (exec) <name>      Replay recorded macro steps with dry-run/verbose options"

	HelpGroupHint    = "  Run any command with --help or -h for detailed usage and examples."
	HelpGroupExample = "  Quick start:"
	HelpExampleScan  = "    $ gitmap scan ~/projects"
	HelpExampleList  = "    $ gitmap ls"
	HelpExamplePull  = "    $ gitmap pull my-api"
	HelpExampleCD    = "    $ gitmap cd my-api"
	HelpCompactHint  = "  Use --compact, --groups, --filter <q> (-f), or --json for scripting (v5.42.0+)."

	HelpAlias           = "  alias (a) <sub>     Assign short names to repos (set, remove, list, show, suggest)"
	HelpSSH             = "  ssh <sub>           Generate, list, and manage SSH keys for Git authentication"
	HelpSSHJoin         = "  sj (ssh-join) [ls|rm|history|auth]   Join an SSH network and manage connections (use --help to expand)"
	HelpCodingGuideline = "  cg <sub>            Install Coding Guidelines (v24) into a repo"
	HelpZipGroup        = "  zip-group (z) <sub>       Manage named file collections for release ZIP archives"
	HelpMV              = "  mv (move) <s..> <d..>       Relocate repo directory with VSCode & GitHub Desktop sync"
	HelpImportExport    = "  import-export (ie)  Export or import gitmap tracked repos, aliases, and groups"
	HelpExportSummary   = "  export              Export tracked repos and settings to a JSON snapshot"
	HelpImportSummary   = "  import              Import a JSON snapshot to restore tracking state"

	// Compact-mode lines: command (alias) only.
	CompactScanning     = "  scan (s), rescan (rsc), rescan-subtree (rss), list (ls)"
	CompactCloning      = "  clone (c), clone-next (cn), desktop-sync (ds), github-desktop (gd)"
	CompactGitOps       = "  pull (p), exec (x), status (st), watch (w), has-any-updates, latest-branch (lb)"
	CompactNavigation   = "  cd (go), group (g), multi-group (mg), alias (a), diff-profiles (dp)"
	CompactRelease      = "  release (r), pull-release (pr), release-self (rs), release-branch (rb), temp-release"
	CompactRelInfo      = "  changelog (cl), changelog-generate, list-versions (lv), list-releases (lr), release-pending (rp), revert, clear-release-json (crj), prune"
	CompactData         = "  export (ex), import (im), profile (pf), bookmark (bk), mv (move), rm (remove/del), db-reset"
	CompactImportExport = "  import-export (ie), export, import"
	CompactHistory      = "  history (hi), history-reset (hr), stats (ss)"
	CompactAmend        = "  amend (am), amend-list (al)"
	CompactProject      = "  go-repos (gr), node-repos (nr), react-repos (rr), cpp-repos (cr), csharp-repos (csr)"
	CompactSSH          = "  ssh, sj (ssh-join)"
	CompactZip          = "  zip-group (z)"
	CompactEnvTools     = "  env, install (in), uninstall (un), installer, cg"
	CompactTasks        = "  task, macro"
	CompactVisualize    = "  dashboard (db)"
	CompactCommitXfer   = "  commit-right (cmr) — LIVE,  commit-left (cml), commit-both (cmb) — scaffolds"
	CompactCluster      = "  servers-clients (sc), clients, cluster"
	CompactUtilities    = "  setup, doctor, update, update-cleanup, version (v), completion (cmp), interactive (i), docs (d), help-dashboard (hd), gomod (gm), seo-write (sw), fix-repo (fr), make-public, make-private, clone-fix-repo (cfr), clone-fix-repo-pub (cfrp), help"

	CompactNoMatchFmt = "  No group matching '%s'. Showing all groups:\n"
	HelpInstaller     = "  installer (in) [run|ls|rm|history|edit]   Manage installer scripts and history (use --help to expand)"
	HelpMacro         = "  macro (m) [run|ls|rm|history|edit]   Manage and execute macros (use --help to expand)"
	HelpSchedule      = "  schedule (sc) [add|ls|rm|pause|resume]   Schedule tasks and run jobs (use --help to expand)"
	HelpVSCode        = "  vscode (vsc) [ls|add|rm|pap|plugins]   Manage VS Code PM integrations (use --help to expand)"
	HelpAgy           = "  antigravity (ag) [ls|add|rm|clear|open|prompt|rw|sync|pap|ep|ip|stats|plugin]   Manage workspaces (use --help to expand)"
	HelpPipeline      = "  pipeline (pl) [status|waittime|error-logs|logs|help]   Monitor CI/CD workflows and ETA (use --help to expand)"
	HelpUI            = "  ui [settings|pipeline|terminal]   Open interactive settings and dashboard in browser"
)

// HelpGroupKeys returns short keywords for tab-completion of group filtering.
var HelpGroupKeys = []string{
	"scanning",
	"cloning",
	"gitops",
	"navigation",
	"release",
	"release-info",
	"data",
	"import-export",
	"history",
	"amend",
	"project",
	"ssh",
	"zip",
	"environment",
	"tasks",
	"visualization",
	"commit-transfer",
	"cluster",
	"utilities",
}
