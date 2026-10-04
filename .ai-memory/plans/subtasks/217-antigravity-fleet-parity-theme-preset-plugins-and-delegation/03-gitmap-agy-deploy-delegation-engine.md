# Subtask 217.3: GitMap Automated Antigravity IDE Delegation Engine

- **Parent Plan:** `pending/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation.md`
- **Spec Reference:** [01-architecture-spec.md](../../../../02-spec/21-app/217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `cli/cmdagy/agy_deploy_cmd.go`, `cli/cmdagy/agy_deploy_types.go`, `cli/cmdssh/ssh_deploy_router.go`, `cli/cmd/roottooling.go`, `cli/cmdagy/agy_cmd.go`

---

## 1. Objective

Implement the `gitmap agy deploy <node>` command and underlying delegation engine to autonomously package, sanitize, transport, and deploy Antigravity IDE configurations (UI themes, permission presets, 4 official plugins, 43 skills, and sanitized instance registries) to remote SSH fleet nodes. Wire the command into `gitmap deploy ide <node>` and `gitmap migrate host <node>` to guarantee perpetual fleet parity in a single CLI invocation.

---

## 2. Architectural Design & Component Breakdown

### 2.1 File Placement & Responsibility Matrix

| File Path | Package | Responsibility |
| :--- | :--- | :--- |
| `cli/cmdagy/agy_deploy_types.go` | `cmdagy` | Structs for deployment options, items, metrics, and JSON telemetry envelopes. Strictly positive booleans (`isAll`, `hasTheme`, `hasPreset`, `hasPlugins`, `hasSkills`, `isForce`, `isDryRun`, `isJSON`). |
| `cli/cmdagy/agy_deploy_cmd.go` | `cmdagy` | CLI argument parsing, SSH target resolution, remote bundling/deploy logic, rich terminal table rendering, and JSON output formatting. |
| `cli/cmdssh/ssh_deploy_router.go` | `cmdssh` | Subcommand routing for `gitmap ssh deploy agy` and `gitmap deploy ide`. |
| `cli/cmd/roottooling.go` | `cmd` | Root CLI dispatch registration for `deploy-ide` and `agy-deploy`. |
| `cli/cmdagy/agy_cmd.go` | `cmdagy` | Cobra subcommand registration (`AgyCmd.AddCommand(agyDeployCmd)`). |

### 2.2 Data Structures & Contracts (`agy_deploy_types.go`)

```go
package cmdagy

import "time"

// AgyDeployOptions holds parsed CLI parameters for Antigravity deployment.
type AgyDeployOptions struct {
	TargetNode string
	IsAll      bool
	HasTheme   bool
	HasPreset  bool
	HasPlugins bool
	HasSkills  bool
	IsForce    bool
	IsDryRun   bool
	IsJSON     bool
}

// AgyDeployItemResult represents the status of an individual deployment component.
type AgyDeployItemResult struct {
	Component string `json:"component"`
	Status    string `json:"status"` // SUCCESS, SKIPPED, FAILED
	Details   string `json:"details"`
}

// AgyDeployMetrics holds numerical counts and timings.
type AgyDeployMetrics struct {
	PluginsCount        int   `json:"pluginsCount"`
	SkillsCount         int   `json:"skillsCount"`
	ProjectsUpdated     int   `json:"projectsUpdated"`
	SanitizedPathsCount int   `json:"sanitizedPathsCount"`
	BytesTransferred    int64 `json:"bytesTransferred"`
	DurationMs          int64 `json:"durationMs"`
}

// AgyDeployResultJSON defines the structured machine-readable output contract.
type AgyDeployResultJSON struct {
	IsSuccess          bool                  `json:"success"`
	Node               string                `json:"node"`
	IP                 string                `json:"ip"`
	Timestamp          time.Time             `json:"timestamp"`
	DeployedComponents map[string]bool       `json:"deployedComponents"`
	Items              []AgyDeployItemResult `json:"items"`
	Metrics            AgyDeployMetrics      `json:"metrics"`
	Errors             []string              `json:"errors"`
}
```

---

## 3. Step-by-Step Implementation Details

### Step 3.1: Define Data Types in `cli/cmdagy/agy_deploy_types.go`
1. Create `cli/cmdagy/agy_deploy_types.go`.
2. Implement `AgyDeployOptions`, `AgyDeployItemResult`, `AgyDeployMetrics`, and `AgyDeployResultJSON`.
3. Adhere strictly to the positive boolean convention (`is*`, `has*`).

### Step 3.2: Implement `cli/cmdagy/agy_deploy_cmd.go`
1. Implement `RunAgyDeployCLI(args []string) error`:
   - Parse flags: `--all`, `--theme`, `--preset`, `--plugins`, `--skills`, `--target`, `--force`, `--dry-run`, `--json`, `-y`.
   - If no target node is provided as a flag or positional argument, render interactive or targeted help.
   - If `--help` or `-h` is passed, render rich Lipgloss/ANSI help output.
2. Implement Target Node Dialing:
   - Resolve target SSH connection using `SSHConnectionsFetcher` and `SSHNodeDialer`.
3. Implement Component Packaging & Transfer:
   - **Theme Deployment:** Sync `customThemeSeedsDark` (`#19191C`, `#BD93F9`, `#F8F8F2`) and VS Code `settings.json`.
   - **Preset Deployment:** Inject `CASCADE_COMMANDS_AUTO_EXECUTION_EAGER` and `AGENT_SETTING_POLICY_ALLOW` into `config.json` and iterate `/home/a/.gemini/config/projects/*.json`.
   - **Plugins Deployment:** Copy `~/.gemini/config/plugins/` (all 4 official plugins).
   - **Skills Deployment:** Validate and link all 43 skills.
   - **Sanitization:** Remap Windows backslashes in `instances.json` to `/home/a/.config/Antigravity` and execute remote cleanup (`rm -rf /home/a/C:*`).
4. Implement Output Formatter:
   - If `opts.IsJSON`: Marshal `AgyDeployResultJSON` and print to `os.Stdout`.
   - Else: Render formatted terminal summary with status icons (`✓`, `✖`, `ℹ`).

### Step 3.3: Wire Routing in `cli/cmdssh/ssh_deploy_router.go`
1. Open `cli/cmdssh/ssh_deploy_router.go`.
2. In `RunSSHDeployRouterCLI(args []string)`:
   - Add cases for `"ide"`, `"agy"`, `"antigravity"`:
     ```go
     case "ide", "agy", "antigravity":
         return cmdagy.RunAgyDeployCLI(append([]string{"--all"}, args[1:]...))
     ```
3. Update `printDeployHelp()`:
   - Add documentation lines for `gitmap deploy ide <node>` and `gitmap agy deploy <node>`.

### Step 3.4: Register in `cli/cmd/roottooling.go` and `cli/cmdagy/agy_cmd.go`
1. In `cli/cmd/roottooling.go`:
   - Under `toolingOpsEntries()`, register:
     ```go
     {[]string{"deploy-ide", "ide-deploy"}, func() error { return cmdagy.RunAgyDeployCLI(append([]string{"--all"}, argsTail()...)) }},
     {[]string{"agy-deploy"}, func() error { return cmdagy.RunAgyDeployCLI(argsTail()) }},
     ```
2. In `cli/cmdagy/agy_cmd.go`:
   - Define `agyDeployCmd` with Cobra syntax and flags.
   - Add `AgyCmd.AddCommand(agyDeployCmd)`.

---

## 4. Verification & Acceptance Criteria

- **Flag Parsing Compliance:** `gitmap agy deploy u1 --all --json` accurately populates `opts.IsAll == true` and `opts.IsJSON == true`.
- **Help Documentation:** `gitmap agy deploy --help` and `gitmap deploy ide --help` output comprehensive descriptions and examples without errors.
- **Positive Booleans Enforced:** All newly authored structs and variables strictly employ positive prefixes (`is*`, `has*`).
- **Strict Coding Guidelines:** Error handling uses structured `apperror.AppError` and `result.Result` types; zero runtime panics.
- **Strict Zero-Build/Test Constraint:** File authoring completed in pure specification mode without executing `go build` or `go test`.
