# Specification: GitMap Node & Prompt System Issues Remediation

## 1. Overview & Strategic Intent
This specification remediates critical operational friction across GitMap's Node Fleet, Prompt Management, Telegram Remote Bot, and Supabase integration:
1. **Node & SSH Fleet Telemetry:** Enforces GitHub node aliasing, exposes bound accounts, instance counts, and active instances across CLI and Web UI, and fixes broken SSH connection/ping handlers.
2. **Interactive Prompt Tree Suite:** Overhauls the Prompt Tree section with a visual project/node tree, wiring fully functional action handlers for Send, Backup, Restore, Reset, and Refresh.
3. **Telegram Remote Dispatch & Priority Ordering:** Fixes Telegram machine target identification and UI send dispatch, pins running and queued prompts to the top of prompt queues, and provides distinct "RUNNING" and "QUEUED"/"ENQUEUED" badges.
4. **Supabase Error Backoff & Transaction Throttling:** Enforces a mandatory 5-minute suppression window after any Supabase error, relaxes rapid UI polling to 5-minute cycles, and enforces a hard $\ge 30$s rate limit on account information and credit syncing.
5. **Prompt Preview Header Polish & CLM Caching:** De-clutters repeated instance metadata (e.g. `default-copy-814`) and prompt text into an expandable "Details" panel, caches CLM conversation lookups to eliminate screen flickering, and repairs template import/export search indexing.
6. **Isolated Testing & Minor Release:** Protects the active running instance (`default-copy-814`) by restricting testing to isolated test instances, followed by full minor version release ceremony and green CI/CD verification.

## 2. User Request (Verbatim)
```markdown
# GitMap Node and Prompt System Issues: high priority instruction, non-negotiable task

A couple of serious issues. When we have the `gitmap` nodes, or we can connect, we should have the name or the aliasing based on GitHub nodes. It has been several times asked to fix your node information, and it was not fixed. Make sure that we fix that issue. The node information, connecting the nodes, sending node, and knowing which node actually uses which account and how many instances are there. That is not yet visible. If I do the prompt machine, then what instance, which account? I don't have the way in the SSH section. This is stupid, and also your connection does not work yet. I think I requested several times, and it was not really respected. Make sure that we respect that. I asked you to work on the prompt tree section. Most of the things are still not improved. The send button, getting back up, and restoring the prompt still does not work. The prompt tree reset also does not do anything. Refresh also does not do anything. Make sure you focus on this, understand the issues, and resolve those issues. The prompt tree view. The next issue is in the Telegram section. The Telegram prompt sending is very buggy, it does not work. It cannot identify the proper machine. The send button is not working properly. I don't see that the running prompt should be on top. Running in your UI shows two, three different running. One is running on top, another is running on the bottom. It could be all the running on top or down, usually the running ones or the recent ones should be on the top, and enqueued one also should be on the top. Enqueued items should be marked with enqueued or queued. These are not respected yet. For the Supabase, if you have a Supabase error, make sure that the Supabase does not check in the next five minutes. It is trying to connect the Supabase over and over again, so it should not get. Why so many transactions happen in this UI? Can it please reduce this transaction? I requested that the transaction should happen in every five minutes or so. Updating the account information, the syncing of the credits, it shouldn't be done less than 30 seconds. This is also not respected. The prompt preview. The preview minimization is nice, but also in the header section, you can see a lot of information which is repeated again and again, like the default copy and things like that. They should go into the details button. User prompt. These things should be in the details button, which is not respected yet. I want you to check the backup and restore prompt. There are glitches on the screen when you are updating and checking. We need to check this. So many retrievals inside this CLM thing. We need to fix it. The import/export button, the indexing, the index is still not fixed. It's been several, more than 10 times, I believe I have given you this. It is not respected. Make sure you list out all these issues, fix all these issues. At the end, verify all these tasks as well, and then check the build and other stuff, and also do a minor bump and release, and also check the CI/CD until green. Make sure you test this stuff with a new instance or the test instance, which is already there, so that the running thing should not be closed.

## task: GitMap Node and Prompt System Issues
## slug: gitmap-node-and-prompt-system-issues

# Actionable Items Must Follow Non-Negotiable

1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/gitmap-node-and-prompt-system-issues/ and enqueue plan task in .ai-memory/plans/gitmap-node-and-prompt-system-issues.md (subtasks in .ai-memory/plans/subtasks/gitmap-node-and-prompt-system-issues/) first
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Fix node information visibility and connection issues, ensuring nodes are correctly identified with their respective accounts and instances.
6. Resolve issues in the prompt tree section, including send button functionality, backup and restore, prompt tree reset, and refresh actions.
7. Address Telegram prompt sending bugs, ensuring proper machine identification and correct UI display of running and enqueued prompts.
8. Implement Supabase error handling to prevent frequent connection attempts and reduce transaction frequency.
9. Refactor prompt preview header to minimize repeated information and move details to a dedicated button.
10. Verify and fix backup and restore prompt glitches, and address retrieval issues in the CLM section.
11. Correct import/export button functionality and indexing issues.
12. Conduct thorough testing on new or test instances, verify all tasks, perform a minor version bump, and ensure CI/CD pipeline is green.

## Must follow and spawn agent using

@[\1](file;.agents/skills/\1)

## Additional Instructions

- /plan first before doing the work to reduce the credits.
- /learn from @[\1](file;.agents/skills/\1) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
```

## 3. Architectural Blueprint

### 3.1 Node & SSH Fleet Telemetry, Account Mapping & Ping Handler
- **GitHub Node Aliasing & Account Resolution:**
  - In `cli/cmdnodes/`, extend `UnifiedFleetNode` and `NodeSummaryResult` with:
    - `Account`: Active GitHub / AGM account (e.g. email or username) mapped to this node.
    - `InstanceCount`: Number of detected IDE / AGY instances running or configured on that node.
    - `ActiveInstances`: String summary of instances (e.g. `default-copy-814 (running)`).
  - In `cli/cmdssh/`, enrich connection resolution so that `gitmap nodes ls` and `gitmap ssh` table outputs display the GitHub node identity, account email, and instance count.
  - In `cli/cmdnodes/nodes_agy_prompt.go` (`printPromptDispatchConfirmation` and remote execution), output the exact target node, account, and instance being invoked.
- **Web UI SSH Connection & Ping Handler:**
  - In `cli/cmdui/ui_types.go`, extend `NodeSummary`:
    ```go
    type NodeSummary struct {
        NodeId        string   `json:"nodeId"`
        Alias         string   `json:"alias"`
        Host          string   `json:"host"`
        Port          int      `json:"port"`
        IsOnline      bool     `json:"isOnline"`
        OS            string   `json:"os"`
        Account       string   `json:"account"`
        Instances     int      `json:"instances"`
        InstanceList  []string `json:"instanceList"`
    }
    ```
  - In `cli/cmdui/ui_server_ssh.go`, populate `Account` and `Instances` from AGM account registry and instance discovery cache.
  - Add `/api/ssh/test` or `/api/ssh/ping` endpoint handler in `ui_server_ssh.go` calling `cmdssh.PingNode(alias)`.
  - In `cli/cmdui/ui_assets_js_nodes.go`, implement missing `testNode(alias)` function to execute ping, update online/offline badge dynamically, and show feedback toast.
  - Update `cli/cmdui/ui_assets_markup_ops.go` to add `Account` and `Instances` columns to the Nodes table.

### 3.2 Interactive Prompt Tree Suite
- **API Endpoints (`cli/cmdagy/agy_ui.go`):**
  - `/api/prompts/tree`: Returns structured hierarchical tree with node aliases, projects, active prompts, queued prompts, and backup batches.
  - `/api/prompts/backup`: Triggers snapshot into `PromptBackupRepository` (`data/backup-prompts/sql.db`) and returns batch ID.
  - `/api/prompts/restore`: Restores queued prompts from specified batch ID.
  - `/api/prompts/reset`: Clears transient queue/errors and resets the prompt tree.
  - `/api/prompts/refresh`: Performs synchronous rescan of active workspaces.
- **Frontend Prompt Tree UI (`cli/cmdagy/agy_ui_html.go` & `cli/cmdui/ui_assets_js_prompts_help.go`):**
  - Replace naive `alert(msg)` with interactive modal/panel displaying:
    - Visual collapsible project nodes (`📁 <projectName> [Node: <alias>]`).
    - Active prompts marked with status badges (`RUNNING`).
    - Queued prompts marked with status badges (`QUEUED`).
    - Action toolbar containing:
      * **Send:** Open dispatch dialog targeting the selected node/project.
      * **Backup:** Create immediate snapshot batch with confirmation toast.
      * **Restore:** Choose recent batch to restore with confirmation toast.
      * **Reset:** Clear completed/stale prompts.
      * **Refresh:** Re-fetch tree instantly.

### 3.3 Telegram Remote Bot & Prompt Priority Sorting
- **Telegram Machine Identification & Send Dispatch:**
  - In `cli/cmdagy/agy_telegram_subcommands.go`:
    - In `handleIncomingTelegramCommand(raw)`: Support explicit machine targeting `/status <machine>`, `/rp <machine>`, `/rwi <machine> <seq> <prompt>`.
    - If machine parameter is omitted or matches current machine name/alias, execute locally.
    - If machine parameter specifies a remote fleet node, forward via SSH dispatch.
  - In Web UI Settings (`cli/cmdui/ui_assets_markup_settings_a.go` & `cli/cmdui/ui_assets_js_settings.go`):
    - Add a "Test Telegram Dispatch" button that calls `/api/telegram/send` with the current token and chat ID, verifying connectivity.
- **Prompt List Priority Sorting & Explicit Badging:**
  - In `cli/cmdui/ui_assets_js_prompts_help.go` and `cli/cmdagy/agy_ui_html.go`:
    - Sort all prompt items:
      1. `status === 'running'` (FIRST / Top of list).
      2. `status === 'queued'` or `status === 'enqueued'` (SECOND / Immediately below running).
      3. `status === 'active'` (THIRD).
      4. `status === 'completed'` or `status === 'idle'` (FOURTH).
    - Render badges:
      - Running items: `<span class="badge badge-online">RUNNING</span>`
      - Queued items: `<span class="badge badge-queued">QUEUED</span>` (or `ENQUEUED`)

### 3.4 Supabase 5-Minute Error Backoff & Transaction Throttling
- **Supabase Error Backoff Gate:**
  - In `cli/cmdsupabase/` (e.g. `supabase_db.go` or `supabase_cmd.go`):
    - Maintain package-level `lastSupabaseErrorAt time.Time` and `lastSupabaseErrorMsg string`.
    - When a connection, authentication, or query error occurs, record `lastSupabaseErrorAt = time.Now()`.
    - At entry of any Supabase probe or connection check:
      `if time.Since(lastSupabaseErrorAt) < 5*time.Minute { return cachedResult, apperror.NewSimple("supabase error suppressed (5m backoff active)", "E_SB_BACKOFF") }`
- **UI Transaction Reduction:**
  - In `cli/cmdagy/agy_ui_html.go` and `cli/cmdui/ui_assets_js_*.go`:
    - Set auto-refresh timer to 5 minutes (`300000ms`) instead of 5000ms.
    - Add manual refresh button for on-demand sync.
- **Account & Credit Sync Rate Limit:**
  - In `cli/cmdagy/agy_account_switch_engine.go` or settings API:
    - Enforce rate limit: `if time.Since(lastCreditSyncAt) < 30*time.Second { return cachedCredits, nil }`. Never trigger credit sync more frequently than once every 30 seconds.

### 3.5 Prompt Preview Header & CLM Optimization
- **Header Clean-up with Expandable "Details" Button:**
  - In Prompt Preview / Execution UI:
    - Move repeated metadata (`instance_id` e.g. `default-copy-814`, process PIDs, timestamps, system info, raw prompt text) into an expandable `<details class="prompt-meta-details"><summary class="btn-xs">Details</summary>...</details>` element.
    - Header displays only the essential project slug, concise title, and status badge.
- **CLM Anti-Glitch & Debounce:**
  - In `cli/cmdagy/agy_conv_*.go`:
    - Add a TTL cache (10s) for conversation transcripts and metadata lookups so rapid UI polls do not trigger disk file rereads or console screen flickering.
- **Prompt Template Import/Export Indexing:**
  - In `cli/cmdprompttemplate/template_import_export.go` and `cli/cmdprompt/prompt_export_import.go`:
    - Ensure that after importing templates or prompt records, the search index and catalog store are immediately updated and rebuilt (`RebuildTemplateIndex()`).

### 3.6 Isolated Test Instance Guard & Minor Bump
- Testing must target an isolated temporary test directory or mock instance; never terminate or modify the active production instance `default-copy-814`.
- Run FastGate (`python 03-ai-scripts/50-fastgate.py`).
- Execute minor bump ceremony (`python 03-ai-scripts/37-bump-version.py --minor`).
- Verify CI/CD pipeline until green.

## 4. Acceptance Criteria & Quality Gates
- **AC-1:** Nodes table in CLI (`gitmap nodes`) and Web UI SSH section displays Node Alias, Host, Account, Instances count, and Status.
- **AC-2:** Web UI SSH Ping button (`testNode`) executes cleanly with live response and error handling.
- **AC-3:** Prompt Tree UI renders interactive tree with working Send, Backup, Restore, Reset, and Refresh buttons.
- **AC-4:** Telegram command parser resolves target machines accurately; UI provides functional Telegram test send.
- **AC-5:** Prompts in Web UI and Studio are strictly sorted with Running at top, Queued immediately after, with prominent badges.
- **AC-6:** Supabase connection failures trigger a 5-minute suppression window preventing retry loops.
- **AC-7:** UI background auto-refresh is relaxed to 5 minutes; credit/account sync enforces $\ge 30$s rate limiting.
- **AC-8:** Prompt preview header collapses repeated instance metadata into an expandable Details element.
- **AC-9:** Template import/export operations immediately rebuild search and catalog indices.
- **AC-10:** FastGate passes, minor version bump executed, and CI/CD pipeline monitored to green.
