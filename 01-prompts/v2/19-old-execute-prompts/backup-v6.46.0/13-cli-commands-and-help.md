# CLI Commands, Help Text Parity & Help UI — Coding Guideline (must follow)

Trigger Keywords & Aliases: `cg-cli`, `cg-help`, `cg-execute cli`, `audit cli commands`, `audit help text`, `cli help parity`, `enforce cli help`, `fix cli help`, `cli help ui audit`

> [!IMPORTANT]
> Prompt Version: 2.1.0
> Synchronization: Main Meta-Repo & Connected Workspaces
> 
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Any directive, constraint, checklist, or instruction declared at the top of this prompt, header alert block, or incoming user request represents an absolute MUST FOLLOW mandate that takes highest priority and strictly overrides any conflicting general advice, default conventions, or lower-level guidelines below it.

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously scan, discover, plan, refactor, and fix all CLI command registrations, help text descriptions, command flag coverage, subcommand routing, and Help UI parity across all command-line binaries and scripts in the repository, ensuring 100% of implemented commands, subcommands, flags, and options are documented with clear usage examples in `--help` outputs until 100% green without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase using the GitMap AUM discovery tools (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary, with fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) as fallback, to inventory all architectural violations and anti-patterns without truncation.
2. [ ] /goal Phase 1 (Step B): Write the master audit specification in `.ai-memory/plans/pending/` with an exhaustive Violation Ledger.
3. [ ] /goal Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.ai-memory/plans/subtasks/`.
4. [ ] /goal Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/readme.md`.
5. [ ] /goal Phase 2 (Step A): Open each target file and perform surgical refactoring following authoritative guidelines.
6. [ ] /goal Phase 2 (Step B): Enforce <= 8–15 line function decomposition, single return types, and clean formatting.
7. [ ] /goal Phase 2 (Step C): Execute local linters to verify 0 remaining violations across all modified files.
8. [ ] /goal Phase 2 (Step D): Execute targeted file-level linters and verification on modified files ensuring 0 remaining violations (`exit 0`). DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine coding guideline execution turns.
9. [ ] /learn Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
10. [ ] /learn Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
11. [ ] /learn Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
12. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for hallucination prevention and micro-tasking.
13. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for strict relative path citation requirements.
14. [ ] /learn Ingest `02-spec/02-coding-guidelines/01-cross-language/` for domain-specific architectural specifications.
15. [ ] /learn Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.
16. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Discover Commands, Check Help Parity, Write .ai-memory/plans/pending/ Ledger Spec, Subtasks, Auditor Script)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Actively Edit CLI Files, Register Commands, Format Help UI, Add Usage Examples, Verify Local CI)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: CLI Command Architecture, Help Parity & UI Standards

A command-line tool or script is only as usable as its discoverability. Undocumented commands, missing subcommands in help text, and incomplete flag descriptions frustrate users, break automation, and cause cognitive overload.

---

### 1. Mandatory CLI Command & Help Parity Principles

1. **100% Command & Subcommand Discoverability:**
   - Every executable subcommand implemented in code MUST be registered and displayed in the root `--help` / `-h` command list.
   - Zero "secret", orphaned, or unlisted commands unless explicitly marked and documented as internal debug flags.
2. **Comprehensive Flag & Option Documentation:**
   - Every flag (e.g. `--config`, `-v`, `--timeout`) MUST have a human-readable description, expected data type, default value (if any), and shorthand alias.
3. **Usage Examples in Help Text:**
   - Every command and subcommand MUST include at least one concrete, real-world usage example in its `--help` output (e.g. `Example: mycli user create --email user@example.com --role admin`).
4. **Standard Help UI Layout:**
   - Help text MUST follow a clean, consistent hierarchical layout:
     - `NAME / USAGE:` Binary name and syntax synopsis.
     - `DESCRIPTION:` 1–2 sentence explanation of command purpose.
     - `COMMANDS / SUBCOMMANDS:` Alphabetical or logical list of available subcommands with 1-line summaries.
     - `OPTIONS / FLAGS:` Formatted table of supported flags and options.
     - `EXAMPLES:` Practical terminal invocations.
5. **Help Invocation Parity:**
   - All standard help flags MUST work identically: `--help`, `-h`, `help <command>`, and invoking the binary without required arguments should display help or a concise error pointing to `--help`.
6. **Unknown Command Error Handling:**
   - If an invalid command is passed, the CLI MUST output a clear error message, suggest closest matching commands if available, and direct the user to `--help`.

---

### 2. Multi-Language CLI Help Implementations

#### 2a. Go (Cobra CLI Framework)

```go
// ❌ WRONG: Missing Short/Long descriptions, missing examples, unregistered subcommands
var userCmd = &cobra.Command{
    Use: "user",
    Run: func(cmd *cobra.Command, args []string) {
        // ...
    },
}

// ✅ CORRECT: Complete command definition with Short, Long, Example, and Flags
package cmd

import (
    "github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
    Use:   "user [command]",
    Short: "Manage system user accounts and credentials",
    Long: `Provides administrative commands to create, inspect, update,
and revoke user accounts and role-based access controls.`,
    Example: `  # Create a new administrator account
  mycli user create --username alice --role admin

  # List active users with JSON output
  mycli user list --status active --format json`,
    Args: cobra.NoArgs,
}

func init() {
    rootCmd.AddCommand(userCmd)
    userCmd.AddCommand(userCreateCmd)
    userCmd.AddCommand(userListCmd)
    userCmd.AddCommand(userDeleteCmd)
}
```

```go
// ✅ REQUIRED: Nested Subcommand Tree Example (e.g., gitmap ssh join, ssh keygen, ssh test)
package cmd

import (
    "github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
    Use:   "ssh [command]",
    Short: "Manage SSH keys, agent forwarding, and remote node connections",
    Long: `Provides a comprehensive suite of SSH subcommands to generate keys,
join clusters, verify tunnel connectivity, and configure authorized keys.`,
    Example: `  # Join a cluster via SSH tunnel
  gitmap ssh join --host node-01.internal --port 22

  # Test SSH key authentication
  gitmap ssh test --key ~/.ssh/id_ed25519 --user git`,
    Args: cobra.NoArgs,
}

var sshJoinCmd = &cobra.Command{
    Use:   "join",
    Short: "Connect and join a remote cluster node via SSH tunnel",
    Example: "  gitmap ssh join --host node-01.internal --port 22",
    RunE:  runSshJoin,
}

var sshTestCmd = &cobra.Command{
    Use:   "test",
    Short: "Verify SSH key connectivity and credentials against a remote host",
    Example: "  gitmap ssh test --key ~/.ssh/id_ed25519 --user git",
    RunE:  runSshTest,
}

func init() {
    rootCmd.AddCommand(sshCmd)
    // Mandatory: Register all nested subcommands to parent sshCmd
    sshCmd.AddCommand(sshJoinCmd)
    sshCmd.AddCommand(sshTestCmd)
}
```

---

#### 2b. TypeScript / Node.js (Commander.js)

```typescript
// ❌ WRONG: Squeezed commands without descriptions or examples
program
    .command('audit')
    .action(runAudit);

// ✅ CORRECT: Fully documented command with description, options, and help examples
import { Command } from 'commander';

export function registerAuditCommand(program: Command): void {
    program
        .command('audit')
        .description('Scan codebase for coding guideline and architectural violations')
        .option('-c, --config <path>', 'Path to custom audit configuration file', 'architect.config.json')
        .option('-f, --format <type>', 'Output format: text, json, or markdown', 'text')
        .option('--strict', 'Treat guideline warnings as blocking build errors', false)
        .addHelpText('after', `
Examples:
  $ mycli audit
  $ mycli audit --format markdown --strict
  $ mycli audit --config ./config/strict-rules.json
`)
        .action(async (options) => {
            await executeAudit(options);
        });
}
```

---

#### 2c. Python (Click / Argparse)

```python
# ❌ WRONG: Undocumented arguments and missing help
import click

@click.group()
def cli():
    pass

@cli.command()
@click.argument("target")
def build(target):
    pass

# ✅ CORRECT: Rich help metadata, options documentation, and epilog examples
import click

@click.group(
    help="Prompt Architect CLI — Multi-agent engineering tooling and automation."
)
@click.version_option(version="1.35.0", prog_name="prompt-architect")
def cli() -> None:
    """Root entry point for Prompt Architect commands."""
    pass

@cli.command(
    name="build",
    short_help="Compile and package target artifacts.",
    help="Builds specified target modules into standalone release packages."
)
@click.argument("target", type=click.STRING)
@click.option(
    "-o", "--output",
    type=click.Path(),
    default="dist/",
    show_default=True,
    help="Directory where compiled release artifacts will be written."
)
@click.option(
    "--optimize/--no-optimize",
    default=True,
    show_default=True,
    help="Enable compiler optimizations and minification."
)
def build(target: str, output: str, optimize: bool) -> None:
    """Execute the build pipeline for the given target."""
    execute_build(target, output, optimize)
```

---

#### 2d. PHP (Symfony Console)

```php
// ❌ WRONG: Missing help text and argument descriptions
class MigrateCommand extends Command {
    protected static $defaultName = 'db:migrate';
}

// ✅ CORRECT: Expressive configure() with full help, arguments, and options
namespace App\Commands;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputOption;

class MigrateCommand extends Command {
    protected static $defaultName = 'db:migrate';

    protected function configure(): void {
        $this
            ->setDescription('Executes pending database schema migrations')
            ->setHelp(<<<'EOF'
The <info>%command.name%</info> command runs all outstanding database migrations:

  <info>php %command.full_name%</info>

To roll back the last migration batch:

  <info>php %command.full_name% --rollback</info>
EOF
            )
            ->addOption(
                'rollback',
                'r',
                InputOption::VALUE_NONE,
                'Roll back the most recent batch of executed migrations'
            )
            ->addOption(
                'dry-run',
                null,
                InputOption::VALUE_NONE,
                'Simulate schema execution without persisting database changes'
            );
    }
}
```

---

## 3. The Phase 1 CLI Command Audit Ledger Format

In Phase 1, you MUST generate `.ai-memory/plans/pending/XX-cli-commands-help-audit.md` containing the following master inventory table:

```markdown
| Command / Script | Implemented Subcommands | Registered in Help UI? | Flag Coverage % | Missing Help Text / Examples | Planned Fix | Status |
|---|---|:---:|:---:|---|---|:---:|
| `cmd/user.go` | `create`, `list`, `delete` | ⚠️ Missing `delete` | 60% | Missing example for `user create` | Register `userDeleteCmd` and add examples | PENDING |
| `src/cli/audit.ts` | `audit` | ✅ YES | 80% | Missing description for `--strict` | Document `--strict` option in command | PENDING |
| `scripts/deploy.py` | `deploy` | ❌ NO | 0% | Missing `--help` parser in script | Migrate to `argparse` with complete help | PENDING |
```

---

---

## Continuous 2-Phase Self-Loop & 2-Agent Concurrency Architecture

To guarantee full execution without stopping after planning mode, the master orchestrator MUST enforce this continuous 2-phase loop:

### 1. 2-Agent Concurrency & Strict `.ai-memory/` Bounding

- **2-Agent Limit (Max 2 Threads Each):** When dispatching work, spawn **at most 2 sub-agents concurrently**, with **no more than 2 threads per agent**.
- **Strict Folder Bounding (`.ai-memory/`):** Subagents can ONLY write planning files, subtasks, status reports, and logs inside `.ai-memory/` (`.ai-memory/plans/`, `.ai-memory/readme.md`, `.ai-memory/memory/issues/`).
- **Context Diet:** Provide subagents with minimal instructions (e.g. "Read subtask file `.ai-memory/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md` and execute it"). Do not paste huge files into agent prompts.

### 2. Phase 1: Planning Mode & Subtask Generation (Steps 1 .. N/2)

### Fast File Discovery & Reading via Python Toolchain (Mandatory Acceleration)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the fast 2-tier discovery toolchain:

### Tier 1: GitMap AUM Acceleration (PRIMARY)
1. **Universal File Search:**
   ```bash
   gitmap find "<pattern>" [-ext <ext>]
   ```
2. **List Indexed Files & Substring Lookup:**
   ```bash
   gitmap list-files [pattern]
   gitmap find-files-any "<substring>"
   ```
3. **Stream File Content:**
   ```bash
   gitmap cat <filepath>
   ```
4. **Instant Code Walk Search:**
   ```bash
   gitmap search "<term>"
   ```

### Tier 2: Fast Cached Python Toolchain (FALLBACK)
1. **Inventory Target Files (with `--limit` option):**
   ```bash
   python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats
   ```
2. **Fast Cached Grep (<15ms, with `--limit` option):**
   ```bash
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --lang go --limit 50
   ```
3. **Sub-Millisecond Folder & File Exploration (with `--limit` option):**
   ```bash
   python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder-path> --ext .go --limit 50
   python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000
   python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<pattern>" --path <folder-path> --limit 50
   ```
4. **Subsystem & Topology Overview:**
   ```bash
   python 03-ai-scripts/18-codebase-topology-discoverer.py --summary
   ```
Do not rely on standard search tools with 50-item truncation when discovering repository-wide violations.

- Spawn 2 planning subagents to scan the codebase for target guideline violations.
- Write the master architectural specification in `.ai-memory/plans/pending/xx-audit.md` with an exhaustive Violation Ledger table.
- Decompose the master plan into granular subtasks in `.ai-memory/plans/subtasks/xx-<parent-slug>/01-<subtask-title>.md`, `02-<subtask-title>.md`, etc.
- **MANDATORY AUTO-LOOP (DO NOT STOP):** Once Phase 1 planning completes, the master orchestrator **MUST NOT STOP or ask the user for confirmation**. It MUST immediately self-loop and transition directly into Phase 2 execution mode.

### 3. Phase 2: Execution Mode & Parallel Refactoring (Steps N/2+1 .. N)

- Spawn 2 execution subagents (max 2 threads each) to execute subtasks in parallel on disjoint files.
- Subagents refactor code following all coding guidelines (<= 8–15 line functions, single return types, universal `*AppError` wrapping, Unix LF line endings).
- Move completed subtasks from `.ai-memory/plans/subtasks/` to `.ai-memory/plans/completed/` and update `.ai-memory/plans/readme.md`.
- **Failure Memory & Feedback Loop:** If a subagent fails:
  - Rollback dirty working tree and log error details to `.ai-memory/plan.md` and `.ai-memory/memory/issues/xx-failure.md`.
  - The next subagent spawned MUST read the previous failure log first, record it as a pending memory task, and implement the necessary fix.
- Execute targeted local linters on modified files ensuring `exit 0` before concluding. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine loops.

## Strict In-Repository Execution & `.ai-memory/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.ai-memory/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.ai-memory/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.ai-memory/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`, `07-relative-path-fixer.py`, `05-naming-autofixer.py`, `09-cli-help-auditor.py`).
>    - RCA & Issue Logs: `.ai-memory/memory/issues/` and `.ai-memory/cicd-issues/`.
>    - Execution Plans & Subtasks: `.ai-memory/plans/pending/`, `.ai-memory/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.ai-memory/coding-guidelines.md`.
> 3. **Worker Pool & Log Aggregation Architecture:** All local runners and test orchestrators must use a concurrent worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without canceling sibling workers, and print a consolidated final summary with full stdout/stderr error logs for failed jobs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
> 5. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.

---

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/readme.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict .ai-memory/ Folder Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **CLI Help Auditor Script:** Use `python 03-ai-scripts/09-cli-help-auditor.py` to scan for CLI entry points, parse `--help` outputs, and verify command registrations.
- [ ] **Go Generate Sync:** If you modify Go constants, enums, or stringers, you MUST run `go generate ./...` in the relevant directory (e.g., `cd gitmap && go generate ./...`) and commit the resulting generated files to prevent CI drift.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/readme.md` using sequential script naming (e.g. `09-cli-help-auditor.py`). For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

---

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new code changes.
- [ ] Sub-agents are actively assigned disjoint files verified against `.ai-memory/readme.md`.
- [ ] Completed tasks were `mv`'d to `plans/completed/` and `.ai-memory/plans/readme.md` was updated.
- [ ] 3-strike rule respected: failed tasks cleanly rolled back and logged to `last-failure.md`.
- [ ] **Strict Relative Git Paths:** All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths (`/absolute/path/to/...`, `/absolute/path/to/...`) or `file:///` URIs.
- [ ] **100% CLI Command Coverage in Help:** Every implemented command and subcommand is registered and visible in `--help`.
- [ ] **Flag Documentation:** All flags/options have descriptions, types, and defaults documented.
- [ ] **Usage Examples:** Every command help includes at least one concrete terminal example.
- [ ] **Standard Help UI:** Follows clean Name, Usage, Commands, Options, and Examples layout.
- [ ] **LF Line Endings (`\n`):** All files use Unix LF line endings. Zero CRLF (`\r\n`).
- [ ] **UTF-8 Encoding (No BOM):** All files encoded in UTF-8 without BOM.
- [ ] **Single Trailing Newline:** Every file ends with exactly one terminating newline (`\n`).
- [ ] **Blank Line Before `if`:** Exactly one blank line precedes every `if` statement (unless at the very top of a block).
- [ ] **Blank Line After `}`:** Exactly one blank line follows every closing brace `}` (unless closing the enclosing block).
- [ ] **Blank Line Before `return`:** Exactly one blank line precedes `return` / `throw` in multi-line blocks.
- [ ] **Zero Double Blank Lines:** No `\n\n\n` in code or markdown.
- [ ] **Markdown Heading Spacing:** Exactly one blank line before and after headings (no leading blank line on line 1).
- [ ] **Zero Nested `if`:** All conditionals flattened to depth 0 using guard clauses and early returns.
- [ ] **Function Sizing:** All functions <= 8 lines preferred (hard cap 15 lines).
- [ ] Coding Guidelines & Master Consolidated File: I have fully read, checked, and strictly enforced every file in `02-spec/02-coding-guidelines/`, as well as the master consolidated coding guideline file at `.ai-memory/coding-guidelines.md`.

1. [ ] /learn and apply as a /goal `.ai-memory/coding-guidelines.md` and also make sure the agent rules are created in the repo to read in the future quickly.

- [ ] `python linter-scripts/check-newline-styling.py` and `python linter-scripts/check-markdown-header-spacing.py` exited with code 0.
- [ ] **Atomic File Recording:** Modified files recorded to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`). Builds, tests, and CI runner are deferred to CI/CD fix.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Strict Relative Git Paths: All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths or `file:///` URIs.
- [ ] Master Guidelines: I have fully read and strictly enforced `02-spec/02-coding-guidelines/` and `.ai-memory/coding-guidelines.md`.
- [ ] CLI Help Parity: All commands, subcommands, and flags are documented in `--help`.
- [ ] LF Line Endings & UTF-8 (No BOM): Verified Unix LF and UTF-8 across all files.
- [ ] Blank Line Before `if`: Verified blank line before every `if` statement across all modified files.
- [ ] Blank Line After `}`: Verified blank line after every closing brace `}` followed by code.
- [ ] Blank Line Before `return`: Verified blank line before every `return`/`throw` in multi-line blocks.
- [ ] Zero Nested `if`: Zero nested `if` statements (depth > 1).

1. [ ] /learn the section as a /goal [AI Fix Scripts Memory](#ai-fix-scripts-memory)

- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## Mandatory Linter & CI/CD Integration

1. **Linter Scripts:** `linter-scripts/check-newline-styling.py`, `linter-scripts/check-function-lengths.py`, `linter-scripts/check-markdown-header-spacing.py`
2. **Local Run Command:** `python 03-ai-scripts/09-cli-help-auditor.py`
3. **Autofixer Command:** `python 03-ai-scripts/05-guideline-autofixer.py <file>`
4. **CI/CD Integration (`.github/workflows/ci.yml`):**
   ```yaml
   - name: Validate CLI Commands & Help Parity
     run: |
       python 03-ai-scripts/09-cli-help-auditor.py
       python linter-scripts/check-newline-styling.py
       python linter-scripts/check-markdown-header-spacing.py
   ```
5. **Runner Registration (`03-ai-scripts/06-cicd-local-runner.py`):**
   ```python
   JOBS = {
       "CLI Help Parity Check": [sys.executable, "03-ai-scripts/09-cli-help-auditor.py"],
       "Newline Styling Check": [sys.executable, "linter-scripts/check-newline-styling.py"],
       "Markdown Header Check": [sys.executable, "linter-scripts/check-markdown-header-spacing.py"],
   }
   ```

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## Strictly Avoid: No Automatic Releases, No Test Running & No Full CI/CD Runner in Routine Turns (Strict Policy)

- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- **NO RELEASES (Strict Policy):** You MUST NOT bump versions, update changelogs, or cut a release at the end of this task. Commits must remain standard development commits. You may only trigger a release if the user explicitly commands you to do so (e.g., "cut a release" or "bump the version").
- **NO TEST RUNNING (Strict Policy):** Test execution is strictly disabled. You MUST NOT execute unit tests, integration tests, or test suites unless explicitly commanded by the repository owner.
- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns. Running the heavy 28-38 gate pipeline across the entire repository wastes massive amounts of time and scans unrelated files. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **Test Inventory & Recent Changes Tracking:** Whenever any file is modified, append its repository-relative path to `.ai-memory/temp/recent-file-changes.json` under atomic file lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <path>`), cross-referencing `.ai-memory/test-inventory.json` so associated tests are known for future release verification.
