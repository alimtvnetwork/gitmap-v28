# Interactive Macro Builder Verification & Context Wrapper Staging

Sequence: 002
CapturedUtc: 2026-09-09T05:15:00Z
Span: 2 user prompts
Topic: Interactive macro builder commands implementation, test verification, and context wrapper staging

---

## User Instructions (verbatim)

### 1.

> the ls should show the dir files and everythng and also the PWD is missing everything where we are and that needs to be display above the line and cna be enable and disbaled, and laso fix the other running stuff like add new commands in it, like replace, find, search, etc, clear???

### 2.

> # Conversation Log & Context Wrapper — Engineering Workflow (must follow)
> 
> > **Prompt Version:** 2.1.0
> > **Synchronization:** Main Meta-Repo & Connected Workspaces
> 
> > Purpose: Before acting on the *next* prompt the user pastes, the AI must (a) persist the full chat so far to disk as a numbered Markdown log, (b) rewrite and improve the user's follow-up prompt in place, (c) confirm that the project's coding guidelines are captured in memory, and (d) surface any ambiguity. The AI does NOT execute the follow-up prompt yet.
> 
> > When to run: The user pastes this prompt and says some variant of "here is the prompt, don't act on it yet" or "rewrite it first". After the AI completes all steps below and the user explicitly says "go" / "now act" / "execute", the AI may then run the rewritten prompt.
> 
> ---
> 
> ## STRICT AVOIDANCE: Never Disable CI/CD
> 
> > [!CAUTION]
> > **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> > Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

---

## Assistant Actions Summary (one bullet per turn, no chain-of-thought)

- Implemented dynamic PWD banner above the step prompt in `gitmap/cmd/macro_add_interactive.go`, rendered via `printMacroPromptPwd()` when enabled.
- Added session override and persistent storage for PWD visibility in `gitmap/uipref/uipref.go` (`IsMacroPwdVisible()`, `SetMacroPwdOverride()`, `SetMacroShowPwd()`) and wired `--pwd` and `--no-pwd` flags in `gitmap/cmd/macro_add.go`.
- Implemented `gitmap/cmd/macro_add_helpers.go` containing in-builder directory listing (`ls`, `dir`), PWD controls (`pwd on`, `pwd off`, `:pwd`), in-builder search (`find <pattern>`, `search <query>`, `grep <query>`), text replacement (`replace <old> <new> [glob]`), directory navigation (`cd <path>`), and in-builder help table.
- Added interactive command recording hooks (`+add` to record last inspected command, `add <cmd>` to record explicitly without executing).
- Refactored `gitmap/cmd/macro_add_interactive.go` into concise functions ($\le 15$ lines, zero nested `if` blocks, positive booleans, blank line before returns).
- Authored test suites in `gitmap/cmd/macro_add_helpers_test.go` and `gitmap/uipref/uipref_test.go`, verifying PWD toggling, directory listing formatting, find/search/replace utilities, and step line dispatching (all passing).
- Verified repository linters (`check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-relative-paths.py`, `check-error-management.py`) with 0 violations across all scanned files.
- Completed Plan 83, moved it to `.lovable/plans/completed/83-interactive-macro-builder-pwd-ls-search.md`, and marked all 4 subtasks complete.
- Updated `walkthrough.md` with complete documentation of delivered interactive commands and test logs.
- Persisted conversation turn to `conversation/002-macro-builder-completion-and-context-wrapper.md` and staged prompt rewrite in `prompts/002-gitmap-interactive-builder-and-clean-sync.md`.

---

## Outcomes / Decisions

- Plan 83 transitioned from pending to completed in `.lovable/plans/01-index.md`.
- Interactive macro builder UX now provides real-time directory awareness and file manipulation helpers without corrupting the macro definition session.
- Staged rewritten prompt and updated prompt library index at `01-prompts/01-prompt-library-setup/01-prompt-library-setup.md`.

---

## Open Threads (carry-over)

- Staging phase active: awaiting user explicit approval ("go" / "execute") before running any follow-up actions or git push.
