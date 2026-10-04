# AI Agent Task Orchestrator & Multi-Tier Split-DB Architecture Specification

## User Request (Verbatim)

```text
# High Priority Instruction

Okay. If you read the latest coding guideline, in the coding guideline, we stated that from now on, the execute task in inner steps, the agents that is created. So these agents actually write to SQLite DB using the Python script. So look into the Python script and also look into the execute task in inner steps prompt that how the JSONs are created, how these AI agents is going to communicate with each other, what is the file structure, folder structure it's going to keep. And this folder structure can be changed from the settings. Okay, by default, we're going to follow the convention that we have that anyone can go through or any LLM can be trained through, and it can actually change the folder structure using the commands. Okay? And let's start. So first we will try to understand how these agents are going to create the temp folder based on the task, and inside this, how it's going to create its own database, root database, and its split DBs to communicate their task. And also root DB would know for the task. And there should be a root of the root DB, like temp folder should have agents.db that actually contains how many task has been created, when the task is completed, how many, let's say, pending agents or cached agents, all type of root level information it should know. So every agent is going to communicate on its own database and also report back to its root DB. So if we look back inside the temp folder, in the AI memory, we have the ai_agents DB, which will be the split DB. So I want you to write this spec inside the system so that it is later on understood how it's going to work. Now, anytime the agents, let's say, spawn based on the task, the task will create its own folder with a sequence. Once the task folder is created, inside this it would have its root DB, and each one of the agent will have its own agent task slack somewhere like this and have a DB to communicate. And the main root DB for that task would know about as many agent as it spawns and slugs and tasks are created. It would have this information. Now, if any sub-agent is failed or did not complete its database or task, then inside the databases, in the sub-tasks databases for that agent, it would actually write all types of logs so that we can investigate later on. Remember that. So the first thing is that we wanted to have our own system where agents or we will have commands, command line commands, that can be easily used to create a database and agent. Agents can communicate using the terminal commands by itself, and they will have their own format, how they wanted to enqueue a task, complete a task, and then finish the task. If the task is not finished, it will log as incomplete. Who made it incomplete and how far it went, we should have a mechanism inside the `gitmap` to know. We could do an AI. So we will have a command like `gitmap` ai-agents. Okay? Inside the AI agents, we will have task. So we will have a root level task. Okay? Inside the root level task, we will have sub-level task. Okay? So we can add sub-task. So if we give the sub-task, then we have to give the parent `Task Id`. Remember that. And we can always do a ls that will tell us how many tasks there was, how many tasks actually recently completed. Usually, it's going to give like 10 or 20. A very short amount. But we can also do a limiting trace or all to see all type of tasks that it has completed. Okay. Yeah, so a task can be created, parent task. When a task is created, it creates its own folder and own root DB. And besides agent, it's going to create the sub-tasks databases for that agent. So each agent will have its own `Id` and database table, and it will be enqueued to that root DB that is created for that task, and the main DB, which is outside ai_agents DB for that repo inside the temp directory. It would know where things are created. And we can trace back to using ls list item help to see all kind of commands. These are there. And also creating task, creating sub-task, completing task, and also logging inside the task. So agents should have all these type of commands to track what files they're reading, what they're searching, so things like that. So there should be enum, so that would have bunch of things like what commands they are running, are they doing a search, are they doing reading, writing to a file, which parts, which lines they are writing to that file. So all kinds of things that we should have from the command line that agents should tap in and do this, and we should able to view this as a visualizer. So you should also create a visualizer. If you do the AI agents UI, that's basically going to view the whole database from the SQLite table and other information nicely in a UI view that we can see in the browser. And also we can delete these tasks, which is done or archived. So long before the task which are done, this should not be indivisible usually, but we can see this in the UI mode or also if we want to see all the tasks that is completed. And again, this can be removed because these are ephemeral or temporary. This does not need to be committed. Remember that. Yeah. So this is how the agents could do the work, and we can also trace AI agents who crashed. That will tell us whichever agent actually created or started a task and did not finish it, and where it failed, how it failed. We should be able to know because agent should be able to let us know it started the reading and then crashed. Okay, that was the last step. So all the timing we need to have so that we can understand how it is going on. We could see all kind of audit trail in the UI mode. We can do drill down, tree view, things like that. You need to consider all types of things. So this is a very big task. I hope you understand. You look into the overall view and try to start implementing it.

Okay. Once you apply all this, I want you to update in the coding guideline the execute task in NS steps for the Python script to `gitmap` with specific commands. And also I want you to create a detailed skill of the `gitmap` so that the AI agents can use easily and also create their stuff, like this task engagement, what they're doing very efficiently. Make sure that also there needs to be a clear command on all or a specific task-based clear command as well for the AI agent section. Okay? Think about this as well, and reset option as well. Also, there should be a temp clear option that will basically remove that temp folder automatically from this directory. So there is that, and also at the same time, I want you to create a skill file for the `gitmap` and also a prompt. I think a prompt is already there in the coding guideline. You look into that prompt for the coding guideline learning skill, and also you update if there is anything that you find wrong or is missing, you can add more that agent might be needed in the terminal use to do things faster. Try to include that as well. Do you understand the task?

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Look into the Python script and execute task in inner steps prompt to understand JSON creation and AI agent communication.
4. Understand and document the file and folder structure for agent tasks, including temp folder and database setup.
5. Implement a system for agents to create databases and communicate via terminal commands.
6. Develop a visualizer for AI agents UI to view database information in a browser.
7. Update the coding guideline for execute task in NS steps with specific `gitmap` commands.
8. Create a detailed skill file for `gitmap` for AI agents.
9. Implement clear and reset options for AI agent tasks and temp folders.
10. Review and update the coding guideline prompt for learning skills.

## Must follow and spawn agent using

@[.agents/skills/execute-parent-task-with-n-steps-v6]
```

---

## 1. System Architecture: Three-Tier Multi-Agent SQLite Hierarchy

```text
.ai-memory/temp-agents/                       <-- Root Temp Directory (Configurable via gitmap config set agent.temp-dir)
├── ai_agents.db                              <-- TIER 1: Repository Root Master Agent DB (Global registry of all tasks, agents, lifecycles)
└── <sequence>-<task-slug>/                   <-- TIER 2: Run-Scoped Parent Task Directory (e.g. 79-ai-agent-task-orchestrator-and-split)
    ├── task.db (or agent-task.db)            <-- Task Root DB (ParentTask, Subtask manifest, claiming, rollups)
    ├── ledger.md                             <-- Markdown ledger for human/LLM transparency
    └── agents/                               <-- TIER 3: Agent-Scoped Split DBs
        ├── worker-01.db                      <-- Dedicated Agent Telemetry, ActionLog & File-touch DB
        └── worker-02.db                      <-- Dedicated Agent Telemetry, ActionLog & File-touch DB
```

### 1.1 Tier 1: Repository Root DB (`ai_agents.db`)
- **Location:** `.ai-memory/temp-agents/ai_agents.db`
- **Purpose:** Acts as the central ledger for all parent tasks, agent runs, overall execution statistics, and crash forensics across the entire repository.
- **Tables:**
  - `ParentTaskRegistry`: Stores `ParentTaskId`, `TaskSlug`, `TaskName`, `RunDirectory`, `RootDbPath`, `Status` (`ACTIVE`, `COMPLETED`, `FAILED`), `TotalStepsBudget`, `CompletedSteps`, `SpawnedAgentCount`, `CreatedAt`, `UpdatedAt`.
  - `AgentRegistry`: Stores `AgentId`, `ParentTaskId`, `AgentRole`, `AgentSlug`, `SplitDbPath`, `Status` (`IDLE`, `WORKING`, `COMPLETED`, `CRASHED`), `LastHeartbeatAt`.
  - `GlobalLifecycleMetrics`: Aggregated run counts, average completion times, failure rates.

### 1.2 Tier 2: Parent Task Root DB (`agent-task.db`)
- **Location:** `.ai-memory/temp-agents/<nn>-<slug>/agent-task.db`
- **Purpose:** Coordinates subtask allocation, dependency order, worker claiming, and task rollup metrics for a single parent task execution.
- **Tables:**
  - `ParentTask`: Task parameters, total budget, current wave, state.
  - `Subtask`: `SubtaskId`, `ParentTaskId`, `TaskCode` (`Task-01`), `Title`, `AssignedAgentRole`, `OwnedFilesJson`, `Status` (`PENDING`, `IN_PROGRESS`, `DONE`, `FAILED`), `Evidence`, timestamps.
  - `SubtaskAuditRollup`: Consolidated view of completed actions, file modifications, and diagnostic summaries.

### 1.3 Tier 3: Agent-Scoped Split DBs (`agents/<agent-slug>.db`)
- **Location:** `.ai-memory/temp-agents/<nn>-<slug>/agents/<agent-slug>.db`
- **Purpose:** Eliminates concurrency write lock contention. Each autonomous worker records its high-frequency actions directly into its own private SQLite database.
- **Tables:**
  - `AgentActionLog`:
    - `ActionLogId`: Primary Key
    - `SubtaskId`: Referenced Subtask ID
    - `ActionType`: Enum string (`SEARCH`, `READ`, `WRITE`, `EXEC`, `LINT`, `CHECK`, `CLAIM`, `COMPLETE`, `FAIL`, `CRASH`)
    - `TargetFile`: Relative file path touched
    - `StartLine`: Target start line (for surgical edits)
    - `EndLine`: Target end line (for surgical edits)
    - `QueryOrCommand`: Search pattern or shell command executed
    - `ActionDetails`: Detailed context and intention
    - `DurationMs`: Execution duration in milliseconds
    - `Status`: `IN_PROGRESS`, `SUCCESS`, `FAILED`
    - `ErrorMessage`: Error or stack trace if failed/crashed
    - `CreatedAt`: ISO 8601 UTC timestamp

---

## 2. CLI Command Suite Specification: `gitmap agent` (and `gitmap ai-agents`)

The CLI commands provide native Go execution, bypassing the need for Python runner scripts:

```text
gitmap agent [subcommand] [flags]
gitmap ai-agents [subcommand] [flags]

Subcommands:
  task init --name <name> [--budget <n>] [--dir <dir>]     Initialize parent task & registers in ai_agents.db
  task ls [--limit <n>] [--all] [--status <s>]              List active and completed tasks
  task status [--task-id <id>]                             Show comprehensive task status & metrics
  
  subtask add --parent <id> --json <json>                  Enqueue subtasks to task database
  subtask claim --agent <role> [--task-id <id>]            Claim next available pending subtask atomically
  subtask start <subtask-id> --agent <role>                Start claimed subtask
  subtask complete <subtask-id> --agent <role> --evidence  Mark subtask completed with evidence
  subtask fail <subtask-id> --agent <role> --reason <r>    Mark subtask failed with RCA details
  subtask ls <task-id>                                     List subtasks under a parent task
  
  log --agent <role> --subtask <id> --action <type> ...    Log granular action telemetry (file read/write/search)
  crashed [--task-id <id>]                                 Detect crashed/abandoned agents & show autopsy
  diagnose [--task-id <id>]                                Deep forensic inspection of in-flight actions
  
  ui [--port <p>] [--browse]                               Launch interactive web visualizer & dashboard
  
  clear [--task-id <id>] [-y]                              Clear specific task or completed tasks
  reset [-y]                                               Reset all agent databases in repository
  temp-clear [-y]                                          Delete entire .ai-memory/temp-agents/ folder
```

---

## 3. Interactive Web Visualizer Dashboard (`gitmap agent ui`)

A browser-based dashboard hosted via `gitmap agent ui` featuring:
1. **Fleet & Task Tree View:** Interactive hierarchical visualization from Root `ai_agents.db` down to Task runs, Subtasks, and Agent Split DBs.
2. **Real-time Live Status Cards:** Active tasks, in-flight agents, completion percentages, steps budget consumption.
3. **Audit Trail & Action Log Explorer:** Filterable timeline of all agent reads, writes, searches, and commands with line ranges and durations.
4. **Crash Autopsy Forensics:** Highlighted red indicators showing exact last action, target file, and reason for any incomplete/crashed subtasks.
5. **Lifecycle Management:** Clean, archive, or delete ephemeral agent runs directly from the browser UI.

---

## 4. Configurable Temp Folder Architecture

- Default: `.ai-memory/temp-agents/`
- Configuration via:
  - CLI: `gitmap config set agent.temp-dir <path>`
  - Environment variable: `GITMAP_AGENT_TEMP_DIR`
  - Command flag: `--dir <path>`
- Dynamic resolution ensures full compatibility across diverse workspace configurations and LLM workflows.
