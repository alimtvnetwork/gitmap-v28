# 181-gitmap-ignore-and-cache-engine

## User Request (Verbatim)
gitmap pull all (pa)
gitmap pull all ssh (pas)


gitmap fix ignore all  [-y]# will fix on all repos to fix gitignroe issues
gitmap fix-ignore-all (fia)  [-y] # will fix on all repos to fix gitignroe issues
gitmap fix ignores all ssh [-y] # will fix on all repos to fix gitignroe issues
gitmap fix-ignores-all-ssh (fias)  [-y] # will fix on all repos to fix gitignroe issues for all nodes, current node will run as is, follow the gitmap pas formula for this and write in the spec as a term for Gitmap PAS formula so that can be reffered by to any AI properly
gitmap commit-push-all-repos (cpar) [-y]# commit all repos if there is any pending
gitmap commit-push-all-repos (cpar) --review(r) # show the pending commits for review first if approved then commits and push
gitmap commit-push-all-repos (cpar) --review --commit-only(co) # show the pending commits for review first if approved then commits and push
gitmap see commit pending
gitmap see git-ignore/ignore/ig issues
gitmap see errors # same as gitmap errors
gitmap see history # same gitmap history

gitmap see errors ssh (ses) # same as gitmap errors, ssh follow gitmap PAS
gitmap pull-all-ssh(pas) # do pull all right now ly/optimize
gitmap ignore/ig add/scan/scan-ssh(ss)/remove/edit/action/ls/help/ui/app/add-group/remove-group/rm-grp/set-default-group/add-grp-to-default(agtd) <name of the group>/apply ./connect-group-with-repo (cgwp) /export/import
gitmap ignore connect-group-with-repo (cgwp) <group-name> <path1>,<alias of the repo> --add-with-default(awd)
gitmap repo-manage ui

gitmap cache create .
gitmap cache  create <relpath1>,<abs path2>
gitmap cache  create "a.json", "b.json"

gitmap cache ls/add/create/remove/rm/help
gitmap cache search "text search" "*.md" [--lines 10] [--limit 20]
gitmap cache search "text search" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]
gitmap cache search "text search" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]
gitmap cache search-multi "text search", "multi *" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]
gitmap cache search-multi-grep "regex search", "multi *" -file-pattern (fp) "a*.md", "a*.md" [--lines 10] [--limit 20]
gitmap cache recache/reconcile/sync
gitmap history ssh
gitmap nodes histories/history

Important instructions must follow
... (omitting long text for brevity, see original prompt) ...
Screenshot: ![Screenshot](../../../assets/screenshots/181-gitmap-ignore-cache-engine-01.png)

## Overview & Domain Context
This feature introduces major enhancements to Gitmap focusing on cross-repo operations (pulling, committing, ignore management) and a split-DB caching engine to speed up AUM operations. We are also optimizing `pull all` and introducing SSH variations utilizing the "Gitmap PAS Formula" for delicate worker execution on nodes.

## Extracted Actionable Task List
- **Task-01**: Implement `gitmap pull all` optimization (decouple ignore checks) and `gitmap pull all ssh` (PAS) formula.
- **Task-02**: Implement `gitmap fix ignore all` (fia) and `gitmap fix ignores all ssh` (fias).
- **Task-03**: Implement `gitmap commit-push-all-repos` (cpar) with `--review`.
- **Task-04**: Implement Gitmap Ignore Group Engine (`gitmap ignore add/group/connect/export/import`).
- **Task-05**: Implement Gitmap Cache Engine Split-DB (`gitmap cache create/search/reconcile`).
- **Task-06**: Implement Gitmap "See" Commands (`gitmap c commit pending/ignore issues/errors/history`).
- **Task-07**: Investigate and fix the UI/pull bug captured in the screenshot.

## System Architecture / Blueprint
- **Ignore Engine**: Maintain a central SQLite DB for groups and paths. When `fix ignore all` runs, check duplicates and if ignored files are already committed.
- **Cache Engine**: Split-DB architecture (Root DB for repo info + root files, separate DB per subfolder). Used for fast regex and text searches. Reconcile on outdated file mod times.
- **PAS Formula**: For SSH nodes, execute cross-repo tasks with low concurrency (2 workers, 2 async operations) to prevent high load. Use task queues to track history.
- **Task Queue**: All actions enqueue to task servers for status tracking, history, and undo/redo capabilities.
