package constants

// Help group headers.
const (
	HelpGroupScanning    = "  Scanning & Discovery:"
	HelpGroupCloning     = "  Cloning & Sync:"
	HelpGroupGitOps      = "  Git Operations:"
	HelpGroupNavigation  = "  Navigation & Organization:"
	HelpGroupRelease     = "  Release & Versioning:"
	HelpGroupReleaseInfo = "  Release History & Info:"
	HelpGroupData        = "  Data, Profiles & Bookmarks:"
	HelpGroupHistory     = "  History & Stats:"
	HelpGroupAmendGroup  = "  Author Amendment:"
	HelpGroupProject     = "  Project Detection:"
	HelpGroupSSH         = "  SSH Key Management:"
	HelpGroupZip         = "  Zip Groups (Release Archives):"
	HelpGroupEnvTools    = "  Environment & Tools:"
	HelpGroupTasks       = "  File-Sync Tasks:"
	HelpGroupUtilities   = "  Utilities:"
	HelpGroupVisualize   = "  Visualization:"
	HelpGroupCommitXfer  = "  Commit Transfer (replay between repos):"
	HelpGroupChromeProf  = "  Chrome Profile (copy / export / import / list / delete):"
	HelpGroupTemplates   = "  Templates & Scaffolding (.gitignore / .gitattributes / LFS):"

	HelpAddIgnore      = "  add ignore [langs...]      Merge curated .gitignore block into ./.gitignore (idempotent, marker-block aware)"
	HelpAddAttributes  = "  add attributes [langs...]  Merge curated .gitattributes block into ./.gitattributes (idempotent, marker-block aware)"
	HelpAddLFSInstall  = "  add lfs-install            Run 'git lfs install --local' and merge the lfs/common .gitattributes block"
	HelpTemplatesInit  = "  templates init (tpl ti)    Scaffold .gitignore + .gitattributes for one or more languages [--lfs] [--dry-run] [--force]"
	HelpTemplatesList  = "  templates list (tpl tl)    List every available template (kind, lang, source: user/embed, path)"
	HelpTemplatesShow  = "  templates show (tpl ts)    Print one template (overlay > embed) to stdout, audit-trail header included"
	HelpTemplatesDiff  = "  templates diff (tpl td)    Preview what add ignore/add attributes would change; exit codes mirror diff(1)"
	HelpSync           = "  sync (sy) <target>         Union-merge curated defaults: ignore | attributes | lfs-install | prettier-ignore | prettier-rc | all  [--dry-run] [--force]"
	HelpCommons        = "  commons (co)               Shortcut for 'sync all' — add/dedupe curated .gitignore, .gitattributes, .prettierignore, .prettierrc + git lfs install  [--dry-run]"

	HelpGroupHint    = "  Run any command with --help or -h for detailed usage and examples."
	HelpGroupExample = "  Quick start:"
	HelpExampleScan  = "    $ gitmap scan ~/projects"
	HelpExampleList  = "    $ gitmap ls"
	HelpExamplePull  = "    $ gitmap pull my-api"
	HelpExampleCD    = "    $ gitmap cd my-api"
	HelpCompactHint  = "  Use --compact, --groups, --filter <q> (-f), or --json for scripting (v5.42.0+)."

	HelpAlias    = "  alias (a) <sub>     Assign short names to repos (set, remove, list, show, suggest)"
	HelpSSH      = "  ssh <sub>           Generate, list, and manage SSH keys for Git authentication"
	HelpZipGroup = "  zip-group (z) <sub> Manage named file collections for release ZIP archives"

	// Compact-mode lines: command (alias) only.
	CompactScanning   = "  scan (s), rescan (rsc), rescan-subtree (rss), list (ls)"
	CompactCloning    = "  clone (c), clone-next (cn), desktop-sync (ds), github-desktop (gd)"
	CompactGitOps     = "  pull (p), exec (x), status (st), watch (w), has-any-updates, latest-branch (lb)"
	CompactNavigation = "  cd (go), group (g), multi-group (mg), alias (a), diff-profiles (dp)"
	CompactRelease    = "  release (r), pull-release (pr), release-self (rs), release-branch (rb), temp-release"
	CompactRelInfo    = "  changelog (cl), changelog-generate, list-versions (lv), list-releases (lr), release-pending (rp), revert, clear-release-json (crj), prune"
	CompactData       = "  export (ex), import (im), profile (pf), bookmark (bk), rm (remove/del), db-reset"
	CompactHistory    = "  history (hi), history-reset (hr), stats (ss)"
	CompactAmend      = "  amend (am), amend-list (al)"
	CompactProject    = "  go-repos (gr), node-repos (nr), react-repos (rr), cpp-repos (cr), csharp-repos (csr)"
	CompactSSH        = "  ssh"
	CompactZip        = "  zip-group (z)"
	CompactEnvTools   = "  env, install (in), uninstall (un)"
	CompactTasks      = "  task"
	CompactVisualize  = "  dashboard (db)"
	CompactCommitXfer = "  commit-right (cmr) — LIVE,  commit-left (cml), commit-both (cmb) — scaffolds"
	CompactUtilities  = "  setup, doctor, update, update-cleanup, version (v), completion (cmp), interactive (i), docs (d), help-dashboard (hd), gomod (gm), seo-write (sw), fix-repo (fr), make-public, make-private, clone-fix-repo (cfr), clone-fix-repo-pub (cfrp), help"

	CompactNoMatchFmt = "  No group matching '%s'. Showing all groups:\n"
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
	"history",
	"amend",
	"project",
	"ssh",
	"zip",
	"environment",
	"tasks",
	"visualization",
	"commit-transfer",
	"utilities",
}
