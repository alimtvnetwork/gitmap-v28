# Subtask 02: Help Dashboard (HD) Web UI & Command Catalog Expansion Plan

> **Subtask ID:** `Subtask-02`  
> **Parent Plan:** `.ai-memory/plans/git-prime-machine-os-version-and-ui-help-docs-improvement.md`  
> **Target Path:** `.ai-memory/plans/subtasks/git-prime-machine-os-version-and-ui-help-docs-improvement/02-hd-docs-ui-and-commands-spec.md`  
> **Status:** `PENDING`  
> **Assigned Subsystem:** `src/` (Help Dashboard Web Frontend)  
> **Related Architecture Spec:** `02-spec/21-app/git-prime-machine-os-version-and-ui-help-docs-improvement/01-architecture-spec.md`  
> **Owned Implementation Files:**  
> - `src/data/commands.ts`  
> - `src/pages/HelpIndex.tsx`  
> - `src/components/docs/DocsSidebar.tsx`  
> - `src/App.tsx`  

---

## 1. Objectives & Scope

1. **Command Catalog Registry Expansion (`src/data/commands.ts`):**
   Register 9 newly introduced, high-impact CLI commands into the strongly typed TypeScript command catalog:
   - `merge-ai` (`ma`): AI-assisted multi-source repository/directory merging with collision avoidance (`01_`/`02_`).
   - `agent-ai` (`aai`): Autonomous multi-agent coordination, Split-DB synchronization, lease management.
   - `heatmap` (`hm`): Real-time agent concurrency and file contention visualizer.
   - `machine` (`mac`): Cross-platform machine identity, hardware specs, and deep OS version detection.
   - `full-summary` (`fsum`): Comprehensive high-density terminal dashboard of repo metrics and git status.
   - `new-commands` (`nc`): Interactive catalog and version changelog of newly introduced GitMap commands.
   - `tree` (`tr`): ASCII directory tree visualizer with git status tracking and file exclusion.
   - `which-format` (`wf`): JSON schema detector and validator for GitMap export files, profiles, and ledgers.
   - `cpar` (`cp`): High-speed concurrent parallel execution engine for git operations across nodes.
2. **Help Index Navigation Updates (`src/pages/HelpIndex.tsx`):**
   Add categorized entries for each of the new commands in `src/pages/HelpIndex.tsx` under their canonical sections with `StatusType.Shipped` status pills, aliases, and direct router links.
3. **Sidebar Navigation Integration (`src/components/docs/DocsSidebar.tsx`):**
   Register intuitive sidebar navigation items with relevant Lucide React icons (`Brain`, `GitMerge`, `Flame`, `Monitor`, `BarChart3`, `Sparkles`, `FolderTree`, `FileCode`, `Zap`) ensuring seamless discoverability.
4. **Router Validation (`src/App.tsx`):**
   Verify all routes resolve properly through `src/App.tsx` without broken router links or 404 dead ends.

---

## 2. Detailed Technical Deliverables

### 2.1 Command Metadata Specification (`src/data/commands.ts`)

Each of the 9 commands must be registered as a complete `CommandDef` item conforming to:
```typescript
export interface CommandDef {
  name: string;
  alias?: string;
  description: string;
  usage?: string;
  flags?: CommandFlag[];
  examples?: CommandExample[];
  howToProceed?: { step: number; title: string; action: string }[];
  notes?: string[];
  category: string;
  seeAlso?: CommandSeeAlso[];
}
```

#### Detailed Definitions:

1. **`merge-ai` (Category: `"movemerge"`):**
   - **Name:** `merge-ai`, **Alias:** `ma`
   - **Description:** AI-assisted multi-source directory and repository merge engine with non-destructive collision sequencing (`01_`/`02_`) and automated AI consolidation handoff.
   - **Usage:** `gitmap merge-ai <target_dir> <src1> [src2...] [--url <repo_url>] [--config <config.json>]`
   - **Flags:**
     - `--destination, -d <path>`: Destination target directory for the consolidated merge.
     - `--source, -s <path>`: Source folder or archive to ingest (repeatable or comma-delimited).
     - `--url <git_url>`: Remote Git repository URL to clone and stage into the merge.
     - `--config <path>`: JSON manifest specifying complex multi-source merge inputs.
     - `--no-ai`: Stage files into numbered collision buckets without invoking AI consolidation.
   - **Examples:**
     - `gitmap merge-ai ./unified ./feature-a ./feature-b`: Merge two local branches into unified.
     - `gitmap ma -d ./consolidated -s ./app1,./app2 --url https://github.com/org/core.git`: Multi-source merge.
   - **See Also:** `merge-both`, `diff`, `mv`.

2. **`agent-ai` (Category: `"antigravity"`):**
   - **Name:** `agent-ai`, **Alias:** `aai`
   - **Description:** Autonomous multi-agent coordination engine managing task decomposition, Split-DB synchronization (`ai_agents.db` + `agent-task.db`), lease fencing, and pre-touch action logs.
   - **Usage:** `gitmap agent-ai <subcommand> [flags]`
   - **Flags:**
     - `--task <slug>`: Target task slug identifier.
     - `--status <state>`: Filter or set subtask status (`PENDING`, `IN_PROGRESS`, `DONE`, `FAILED`).
     - `--lease <machine_id>`: Claim or query task lease reservation.
     - `--json`: Output machine-readable JSON representation.
   - **Examples:**
     - `gitmap agent-ai status`: View global agent ledger and running workers.
     - `gitmap aai lease --task refactor-auth`: Inspect or acquire task lease.
   - **See Also:** `heatmap`, `tasks`, `prompt-run`.

3. **`heatmap` (Category: `"history"`):**
   - **Name:** `heatmap`, **Alias:** `hm`
   - **Description:** Real-time agent concurrency and file contention visualizer highlighting active file touches, collision blast radius, and high-contention repository hotspots.
   - **Usage:** `gitmap heatmap [discuss|live] [--window <duration>] [--json]`
   - **Flags:**
     - `--window <duration>`: Historical sliding window for file contention (default: `15m`).
     - `--threshold <n>`: Minimum simultaneous touches to trigger a hotspot warning (default: `2`).
     - `--json`: Output raw contention metrics as JSON.
     - `--live`: Stream real-time file touches directly to the terminal.
   - **Examples:**
     - `gitmap heatmap discuss`: Render terminal summary of agent file contention.
     - `gitmap hm --live`: Open real-time terminal stream.
   - **See Also:** `agent-ai`, `stats`, `watch`.

4. **`machine` (Category: `"tools"`):**
   - **Name:** `machine`, **Alias:** `mac`
   - **Description:** Inspect local host and cluster fleet machine identity, network adapters, hardware specifications, and deep operating system version details.
   - **Usage:** `gitmap machine [subcommand] [--ssh] [--json]`
   - **Flags:**
     - `--ssh, -s`: Query and display registered fleet nodes across the SSH cluster.
     - `--json, -j`: Output structured JSON with `osVersionDetails` object.
     - `--set-alias <name>`: Assign a custom persistent alias to the local machine.
     - `--set-name <name>`: Override the advertised machine hostname.
   - **Examples:**
     - `gitmap machine`: Display local hierarchical machine identity and OS tree.
     - `gitmap mac --ssh`: List all cluster fleet nodes with OS and GitMap versions.
     - `gitmap machine --json`: Export machine identity for automated tooling.
   - **See Also:** `alias`, `ssh`, `doctor`.

5. **`full-summary` (Category: `"history"`):**
   - **Name:** `full-summary`, **Alias:** `fsum`
   - **Description:** Comprehensive high-density terminal dashboard displaying branch topologies, dirty working trees, untracked files, ahead/behind counters, and repository health scores across scanned directories.
   - **Usage:** `gitmap full-summary [dir] [--depth <n>] [--filter <expr>]`
   - **Flags:**
     - `--depth <n>`: Directory recursion limit (default: `3`).
     - `--filter <expr>`: Filter repositories by name, language, or branch.
     - `--export <path>`: Export summary dashboard to Markdown or HTML report.
   - **Examples:**
     - `gitmap full-summary`: Render full summary for current workspace.
     - `gitmap fsum ~/projects --depth 2`: Summary of projects directory.
   - **See Also:** `status`, `stats`, `watch`.

6. **`new-commands` (Category: `"tools"`):**
   - **Name:** `new-commands`, **Alias:** `nc`
   - **Description:** Interactive command discovery tool cataloging recently introduced GitMap CLI commands, version milestones, category groupings, and usage cheat sheets.
   - **Usage:** `gitmap new-commands [--since <version>] [--category <key>]`
   - **Flags:**
     - `--since <version>`: Filter commands introduced since specified release (e.g. `v6.500.0`).
     - `--category <key>`: Filter by command category (`tools`, `movemerge`, etc.).
     - `--search <term>`: Search command descriptions and usage strings.
   - **Examples:**
     - `gitmap new-commands`: List commands introduced in recent releases.
     - `gitmap nc --since v6.520.0`: Show latest command additions.
   - **See Also:** `generic-cli`, `help-dashboard`, `commands`.

7. **`tree` (Category: `"scanning"`):**
   - **Name:** `tree`, **Alias:** `tr`
   - **Description:** Unicode tree visualizer rendering project directory structures with integrated Git status badges, ignored file masking, and custom depth controls.
   - **Usage:** `gitmap tree [dir] [--depth <n>] [--all] [--no-ignore]`
   - **Flags:**
     - `--depth <n>`: Maximum directory traversal depth (default: `2`).
     - `--all, -a`: Show hidden dotfiles and system directories.
     - `--no-ignore`: Display files ignored by `.gitignore`.
   - **Examples:**
     - `gitmap tree`: Show 2-level directory tree of current repository.
     - `gitmap tr ./src --depth 3`: Inspect src directory hierarchy.
   - **See Also:** `scan`, `rescan`, `fspath`.

8. **`which-format` (Category: `"data"`):**
   - **Name:** `which-format`, **Alias:** `wf`
   - **Description:** JSON schema inspector and validator that analyzes arbitrary JSON files, identifies their GitMap schema family (scan export, bookmark, profile, task ledger), and validates payload integrity.
   - **Usage:** `gitmap which-format [file.json...] [--dir <directory>]`
   - **Flags:**
     - `--dir <directory>`: Batch inspect all JSON files in specified folder.
     - `--strict`: Enable strict schema validation against known interfaces.
     - `--json`: Output inspection findings as JSON.
   - **Examples:**
     - `gitmap which-format ./export.json`: Identify format of exported data file.
     - `gitmap wf --dir ./backups`: Batch inspect backup payloads.
   - **See Also:** `import`, `export`, `profile`.

9. **`cpar` (Category: `"monitoring"`):**
   - **Name:** `cpar`, **Alias:** `cp`
   - **Description:** High-speed concurrent parallel execution engine executing arbitrary git commands or scripts across dozens of repositories simultaneously with throttled worker pools.
   - **Usage:** `gitmap cpar "<command>" [--concurrency <n>] [--timeout <sec>]`
   - **Flags:**
     - `--concurrency <n>`: Maximum concurrent worker routines (default: CPU cores * 2).
     - `--timeout <sec>`: Command execution timeout per repository.
     - `--continue-on-error`: Prevent single repository failure from halting batch run.
   - **Examples:**
     - `gitmap cpar "git status -s"`: Concurrently inspect status across all repos.
     - `gitmap cp "git fetch --all" --concurrency 8`: Parallel fetch across 8 workers.
   - **See Also:** `batch-actions`, `cluster`, `pipeline`.

---

### 2.2 Help Index Updates (`src/pages/HelpIndex.tsx`)
In `src/pages/HelpIndex.tsx`, update the `sections` array to include each command under its corresponding section:
- Under section `"merge"`: Add entry for `merge-ai` (`to: "/merge-ai"`, `cmd: "merge-ai"`, `alias: "ma"`, `status: StatusType.Shipped`).
- Under section `"diff"` or new `"agent"` / `"tools"`: Add entries for `agent-ai`, `heatmap`, `machine`, `full-summary`, `new-commands`, `tree`, `which-format`, and `cpar`.
- Ensure all entries declare `status: StatusType.Shipped` with positive visual styling.

### 2.3 Sidebar Navigation Items (`src/components/docs/DocsSidebar.tsx`)
In `src/components/docs/DocsSidebar.tsx`, add navigation entries to `navItems`:
```typescript
{ title: "merge-ai (cmd)", url: "/merge-ai", icon: GitMerge },
{ title: "agent-ai (cmd)", url: "/agent-ai", icon: Boxes },
{ title: "heatmap (cmd)", url: "/heatmap", icon: Activity },
{ title: "machine (cmd)", url: "/machine", icon: Monitor },
{ title: "full-summary (cmd)", url: "/full-summary", icon: BarChart3 },
{ title: "new-commands (cmd)", url: "/new-commands", icon: Terminal },
{ title: "tree (cmd)", url: "/tree", icon: FolderOpen },
{ title: "which-format (cmd)", url: "/which-format", icon: FileText },
{ title: "cpar (cmd)", url: "/cpar", icon: Zap },
```

### 2.4 Application Routes & Build Integrity (`src/App.tsx`)
Verify in `src/App.tsx` that route endpoints are configured. For commands that leverage existing command detail templates (e.g. `GenericCLIPage` or dedicated command view components), ensure route definitions exist so clicking sidebar and help index pills does not trigger `NotFound`.

---

## 3. Acceptance Criteria & Invariants

- **`AC-06` (Catalog Completeness):** All 9 commands (`merge-ai`, `agent-ai`, `heatmap`, `machine`, `full-summary`, `new-commands`, `tree`, `which-format`, `cpar`) are present in `src/data/commands.ts` with valid descriptions, usage strings, flags, and examples.
- **`AC-07` (Type-Safe CommandDef):** The updated `src/data/commands.ts` passes TypeScript compiler validation (`tsc --noEmit` / `pnpm build`) without type errors or missing required properties.
- **`AC-08` (Help Index Shipped Badges):** `src/pages/HelpIndex.tsx` renders entries for the new commands displaying the green `Shipped` status pill (`StatusType.Shipped`).
- **`AC-09` (Sidebar Nav Parity):** `src/components/docs/DocsSidebar.tsx` includes accessible navigation links with distinct icons for the new commands.
- **`AC-10` (Zero 404 Route Breakages):** All registered command links in the Help Dashboard point to valid, reachable React routes.
- **`INV-01` (Affirmative Booleans):** Status and flags avoid negative polarity.
- **`INV-03` (Strict Relative Paths):** All imported paths and component references strictly use project-relative paths.
- **`INV-05` (Offline-First Bundle):** No external assets or remote CDN links introduced.

---

## 4. Verification Protocol & Build Execution

### 4.1 Frontend Build & Typecheck
Verify compilation and type integrity:
```bash
# 1. Type check TypeScript definitions
pnpm tsc --noEmit

# 2. Build Vite production bundle
pnpm build
```

### 4.2 Interactive Verification via Help Dashboard
Verify local serving through the compiled Go binary:
```bash
# Launch help dashboard server and verify command catalog in browser
gitmap help-dashboard --port 8080
# or alias
gitmap hd
```

### 4.3 Evidence Checks
- [ ] `src/data/commands.ts` exports all 9 new commands.
- [ ] Vite production build completes with exit code 0 (`dist/` directory generated).
- [ ] Navigating to `/commands` in Help Dashboard displays `merge-ai`, `agent-ai`, `machine`, etc. in search results and filter lists.
- [ ] Navigating to `/` (Help Index) shows the new commands with `Shipped` badges.
- [ ] Sidebar links navigate to valid pages without triggering 404 views.
