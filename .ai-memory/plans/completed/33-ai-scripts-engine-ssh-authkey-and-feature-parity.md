# Consolidated Plan 33: AI Scripts Engine, SSH Authorized Key Deployment & Feature Parity

## Origin & Execution Metadata
- **Originating User Request**: "Can you please confirm that you have SSH update all Gitmap on all machines and SSH join execute, these are running all right? Uh, it can run the commands and then get back to us with nice UI terminal and everything. So these are all done. Please confirm that. Um, and also the cluster servers clients. So these are also working good. The cluster section also works good. I can able to, uh, deploy the authorized key. This command needs to be there. Um, also the ZSH, I could change themes, add automation in ZSH, which you should have learned from the Kubernetes, uh, repo by Adem. We should have integrated similar feature. Confirm all of these are done properly. And also now we can interact with the Antigravity tool to send the errors, fix the errors, add templates for prompts and so on. Uh, current prompts will be injected inside it. Also, at the same time, uh, I want you to understand the Python scripts that we have, um, for the AI script section, and I want you to have a generic way of creating, uh, this type of functionality inside Gitmap so that we can run and fix stuff, um, from the Gitmap. So I want you to learn this code very effectively, um, then make a detailed plan. That's another step. And then you implement the functionalities into Gitmap, and also add those functionalities into the help text terminal UI, and also the UI of the HD doc."
- **Execution Workflow**: Parent Task N-Step Continuous Loop & Multi-Agent Orchestration (v2.2.0, N=450).
- **Execution Summary**: Successfully audited existing features, fixed `--all` fleet update parsing, implemented hardened authorized key deployment, restored 100% mutual cluster/SC parity, integrated Adem's Kubernetes ZSH automation & 41 theme catalog, and built a native Go AI scripts discovery, live execution, and autofix suite (`gitmap ai` / `gitmap scripts`) with AST parity and documentation.
- **Completed Loops/Steps**: 5 subtasks across 2 execution phases.

---

## 1. Subtask 01: SSH Fleet Update & Authorized Key Deployment Engine
- **Fixed `--all` / `-a` Flag Handling**:
  - In `cli/cmdssh/ssh_update_remote.go` and `cli/cmdssh/ssh_install_remote.go`, refactored argument parsing to normalize `--all` and `-a` into target `"all"` rather than falsely interpreting it as a package name.
- **Implemented `gitmap ssh auth-key deploy` / `gitmap ssh copy-id`**:
  - Created `cli/cmdssh/ssh_auth_key_deploy.go`.
  - Discovers public keys via explicit `-i <path>` flag or standard fallback hierarchy (`~/.ssh/id_ed25519.pub`, `~/.ssh/id_rsa.pub`, `~/.ssh/id_ecdsa.pub`).
  - Hardens remote permissions to `0700` for `~/.ssh` and `0600` for `~/.ssh/authorized_keys`.
  - Ensures idempotent deduplication via `grep -qF` before appending.
  - Cross-platform PowerShell parity for Windows remote nodes.
  - Wired into `gitmap ssh`, `gitmap sj`, `gitmap cluster`, and `gitmap sc`.

---

## 2. Subtask 02: Cluster & SC Mutual Parity and ZSH Automation
- **Restored 100% Cluster & SC Mutual Parity**:
  - In `cli/cmd/rootcore.go`, updated `dispatchSCNodeOps`, `dispatchServersClients`, and `dispatchClients` to route `node`, `bootstrap`, and `k8s` subcommands directly to `cmdssh.RouteClusterNodeCLI`, `cmdssh.RunClusterBootstrapCLI`, and `cmdssh.RouteClusterK8sCLI`.
  - Eliminated the gap where `gitmap sc node ...` failed.
- **ZSH Workspace Scaffolding & Automation (Adem Kubernetes Patterns)**:
  - In `cli/cmdssh/cluster_node_recipes.go`, enhanced `createUserScriptTemplate` to scaffold `$HOMEDIR/scripts`, `$HOMEDIR/gitlab`, `$HOMEDIR/github`, and `$HOMEDIR/.ssh`.
  - Automated unattended cloning of `zsh-autosuggestions` into `$HOMEDIR/.oh-my-zsh/custom/plugins/zsh-autosuggestions`.
- **41 Oh-My-Zsh Themes Catalog**:
  - Created `cli/cmdzsh/theme_catalog.go` declaring the 41 supported Oh-My-Zsh themes, `IsSupportedTheme` validator, and `DisplayThemesCatalog`.
  - Wired `gitmap zsh themes` in `cli/cmdzsh/zsh_cmd.go`.

---

## 3. Subtask 03: Native AI Scripts Catalog & Types Architecture
- **Domain Types & Data Models (`cli/cmdai/types.go`)**:
  - Declared `ScriptCategoryType` (`engine`, `format`, `ci`, `audit`, `discovery`, `plans`, `git`, `release`, `database`, `migration`, `inventory`).
  - Declared `ScriptMetadata`, `PythonRuntime`, `ScriptListOptions`, `FixTargetType`, and typed Result aliases.
- **Cross-Platform Python 3 Discovery (`cli/cmdai/ai_python_detector.go`)**:
  - Memoized runtime detection inspecting `GITMAP_PYTHON`, `PYTHON_EXE`, Windows `python.exe`/`py.exe`, and Unix `python3`.
  - Validates candidates with a 3-second context probe (`--version`) to reject Windows App Store 0-byte stubs.
- **Master Catalog Registry (`cli/cmdai/ai_catalog.go`)**:
  - Registered all 44 Python scripts from `03-ai-scripts/` with numbers, filenames, slugs, aliases, categories, and autofix flags.
- **Listing & Table Formatting (`cli/cmdai/ai_list.go`, `cli/cmdai/ai_list_render.go`)**:
  - Implemented `RunAiList` with category filtering, fix-only filtering, search keywords, JSON streaming, and auto-aligned terminal table rendering via `termtable`.

---

## 4. Subtask 04: Native AI Scripts Execution & Autofix Engine
- **Unbuffered Live Streaming Runner (`cli/cmdai/ai_exec.go`)**:
  - Anchors working directory to repository root containing `03-ai-scripts` or `.git`.
  - Pipes process stdout, stderr, and stdin directly to `os.Stdout`, `os.Stderr`, and `os.Stdin` with `PYTHONUNBUFFERED=1`, preserving ANSI colors and real-time progress bars.
  - Extracts process exit codes into structured `*apperror.AppError`.
- **Script Target Resolution & Forwarding (`cli/cmdai/ai_run.go`)**:
  - Supports numeric tokens (`06`, `6`), slugs (`cicd-local-runner`), and aliases (`ci`, `runner`).
  - Forwards extra CLI arguments directly to the target script.
- **Autofix Dispatch Pipeline (`cli/cmdai/ai_fix_targets.go`, `cli/cmdai/ai_fix.go`)**:
  - Mapped standard targets: `guidelines`, `newlines`, `paths`, `naming`, `encoding`, `spelling`, `plans`, `all`.
  - `gitmap ai fix all` runs an automated sequential repair pipeline.
- **Cobra Command Hierarchy (`cli/cmdai/ai_cmd.go`)**:
  - `gitmap ai list` (`gitmap ai ls`)
  - `gitmap ai run <target> [args...]`
  - `gitmap ai fix <target> [args...]`
  - Direct dispatch: `gitmap ai 06 --no-tests` automatically routes to `run 06 --no-tests`.

---

## 5. Subtask 05: CLI Wiring, AST Constants, Helptext & HD Docs
- **AST Constants & Parity (`cli/constants/constants_cli.go`, `cli/constants/cmd_constants_test.go`)**:
  - Registered `CmdAi = "ai"`, `CmdAiAlias = "scripts"`, and `HelpAi`.
  - Added to `topLevelCmds()` satisfying AST parity tests (`TestTopLevelCmdRegistryMatchesAST`, `TestTopLevelCmdConstantsAreUnique`).
- **Root Dispatcher (`cli/cmd/root.go`)**:
  - Wired `dispatchGeneralCommands` routing `"ai"` and `"scripts"` to `cmdai.DispatchAi`.
- **Embedded Help Text (`cli/helptext/ai.md`)**:
  - Authored concise help (<120 lines) with top-level `## Examples` and fenced code blocks.
  - Registered topic summaries in `cli/helptext/catalog.go` and alias in `cli/helptext/print.go`.
- **LLM / HD Documentation (`cli/llm.md`, `cli/helptext/llm.md`, `llm.md`)**:
  - Updated Phase 3 verification diagrams, added Section 5 documenting the Native AI Scripts Catalog & Autofix Engine, and updated command cheat sheets.

---

## Verification Results
- `check-boolean-guidelines.py`: PASS (0 violations across 2,391 source files).
- `check-nested-ifs.py`: PASS (0 depth-2 nested ifs).
- `check-error-management.py`: PASS (0 bare panics, exits, or swallowed errors across 3,237 source files).
- `check-newline-styling.py`: PASS (100% strict Unix LF line endings).
- `check-constants-naming.py`: PASS (all 5,740 constants pass naming rules).
- `check-constants-collisions.py`: PASS (0 collisions across 5,995 identifiers).
- `09-cli-help-auditor.py`: PASS (all 3,364 CLI source files passed help audit).
- Function length verification: PASS (all functions $\le 15$ lines).
- File length verification: PASS (all new files $\le 200$ lines).
