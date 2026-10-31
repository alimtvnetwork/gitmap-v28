# Specification 359: Component & CLI Contracts

- **Spec Number:** 359
- **Component:** `cli/cmdlogin`, `cli/cmdcommit`, `.agents/skills`, `01-prompts`
- **Parent Spec:** `01-architecture-spec.md`

---

## 1. CLI Commands & Signatures

### 1.1 `gitmap login set ai-auto` Command Syntax
```bash
# Set repository or global autonomous identity for specific agent
gitmap login set ai-auto <agent-name>
gitmap login set ai-auto muse
gitmap login set ai-auto gemini
gitmap login set ai-auto cursor-model

# Set system-wide global autonomous identity without specific agent name
gitmap login set global ai-auto
gitmap login set ai-auto global
gitmap login set ai-auto --global

# Set global autonomous identity with specific agent name
gitmap login set global ai-auto muse
gitmap login set ai-auto muse --global

# Combine with time-restriction defaults
gitmap login set ai-auto muse --restriction-time
```

### 1.2 `gitmap cpf/cpb/cpc/cfr/cpr` with `--ai-agent`
```bash
# Feature commit authored by specific agent with heatmap attribution
gitmap cpf "Module - summary" --ai-agent muse

# Bug fix commit authored by specific agent
gitmap cbf "Module - summary" --ai-agent gemini

# Release commit with agent attribution
gitmap cpr "v6.633.0 - release summary" --ai-agent muse
```

---

## 2. Component Design & Responsibilities

### 2.1 `cli/cmdlogin` (Authentication & Identity)
- `dispatchSubcommand(args []string)`: Checks if `args[0] == "set"` and delegates to `runLoginSet(args[1:])`.
- `runLoginSet(args []string)`: Parses tokens (`ai-auto`, `global`, `--global`, `agentName`, `--restriction-time`), performs identity configuration, updates profile configuration, and prints user feedback.
- `resolveAIBotName(agentName string) string`: Formats the bot identity name:
  - If `agentName` is empty: `"gitmap automation bot - runs autonomously"`
  - If `agentName` is specified: `"gitmap automation bot - <agentName> - runs autonomously"`
- `backupOriginalUser(cfg *model.GitProfileConfig)`: Guarded by `isAIBotIdentity(name)` to prevent overwriting human identity backups with bot identities.

### 2.2 `cli/cmdcommit` (Commit Flag & Heatmap Attribution)
- `constants/constants.go`: Defines `FlagAIAgent = "--ai-agent"`.
- `commit_types.go`: Defines `FlagAIAgent = commitconstants.FlagAIAgent`.
- `commit_flags.go`: Ingests `--ai-agent <name>` into `meta.AgentID`, `meta.AIName`, and `meta.IsAIMode = true`. If `meta.Author == ""`, defaults to the formatted AI bot author string.
- `commit_history_heatmap.go`: Ensures `ModelName` and `Notes` in `HeatmapActionLog` and `HeatmapCommitLog` contain agent ID/name attribution.

---

## 3. Cross-Repository Synchronization
The following artifacts must be kept in 100% parity across the `gitmap-v28` and `coding-guidelines-v24` repos:
1. `.agents/skills/gitmap/skill.md`
2. `.agents/skills/gitmap-automation-bot/skill.md`
3. `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`
4. `.agents/skills/execute-parent-task-with-n-steps-v7/skill.md`
5. `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`
6. `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md`
