# Subtask 03: GitMap CLI `ide setup` Command, Subcommand Dispatch, and Flag Parser

> **Subtask ID:** 352-03  
> **Sequence ID:** 352  
> **Status:** PENDING  
> **Module Path:** `cli/cmdide/`  
> **Parent Plan:** `.ai-memory/plans/pending/352-cursor-like-ide-setup-from-gitmap.md`  
> **Specification Reference:** `02-spec/21-app/352-cursor-like-ide-setup-from-gitmap/01-cursor-like-ide-setup-spec.md`  

---

## 1. Objective & Scope

Integrate the Cursor-like IDE setup into the core GitMap CLI workflow:
1. Wire `gitmap ide setup` (along with aliases `cursor-like`, `configure`) into `cli/cmdide/ide_cmd.go`.
2. Implement a robust flag parser in `cli/cmdide/setup_flags.go` supporting `--provider`, `--model`, `--base-url`, `--api-key`, `--dry-run`, `--json`, `--test`, `--force`, and `--backup`.
3. Provide rich terminal feedback in `cli/cmdide/setup_view.go` using Lipgloss styling, including provider details, configured target tables, measured latency, and editor launch recommendations.
4. Support machine-readable `--json` output formatted as a standard `result.Result[IDESetupSummary]` envelope.
5. Adhere to strict Error Management: Leaf functions return `*apperror.AppError` directly without calling `os.Exit`.

---

## 2. Command Routing & Subcommand Dispatch (`cli/cmdide/ide_cmd.go`)

Update `dispatchSubcommand` to recognize setup subcommands:

```go
func dispatchSubcommand(subcmd string, args []string) error {
	switch subcmd {
	case "setup", "cursor-like", "configure":
		return RunIDESetup(args)
	case "add", "register":
		return runIDEAdd(args)
	case "sync":
		return runIDESync(args)
	case "remove", "rm":
		return runIDERemove(args)
	case "list", "ls":
		return runIDEList(args)
	case "status":
		return runIDEStatus(args)
	case "help", "-h", "--help":
		return RunHelp()
	default:
		return apperror.NewValidation("cmdide.RunIDE", "E1000", "unknown ide subcommand: "+subcmd)
	}
}
```

---

## 3. Flag Parsing Architecture (`cli/cmdide/setup_flags.go`)

Flattened parsing loop adhering to positive boolean conventions (`is*`, `has*`):

```go
package cmdide

import (
	"strings"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/setup"
)

func defaultSetupOptions() setup.IDESetupOptions {
	return setup.IDESetupOptions{
		Target:          enums.TargetAllAI,
		Provider:        enums.ProviderKimi,
		IsBackupEnabled: true,
		IsForceCreate:   true,
	}
}

func parseSetupOptions(args []string) (setup.IDESetupOptions, error) {
	opts := defaultSetupOptions()
	for i := 0; i < len(args); i++ {
		step, isFlag, err := parseSetupFlag(&opts, args, i)
		if err != nil {
			return opts, err
		}
		if isFlag {
			i += step
			continue
		}
		opts.Target = resolveTargetArg(args[i])
	}

	applyPresetDefaults(&opts)
	opts.APIKey = setup.ResolveAPIKey(opts.Provider, opts.APIKey)

	return opts, nil
}

func parseSetupFlag(opts *setup.IDESetupOptions, args []string, idx int) (int, bool, error) {
	arg := args[idx]
	if step, hasVal, err := parseValueFlags(opts, args, idx); hasVal {
		return step, true, err
	}

	return parseBoolFlags(opts, arg)
}

func parseBoolFlags(opts *setup.IDESetupOptions, arg string) (int, bool, error) {
	switch arg {
	case "--dry-run":
		opts.IsDryRun = true
		return 0, true, nil
	case "--json":
		opts.IsJSON = true
		return 0, true, nil
	case "--quiet", "-q":
		opts.IsQuiet = true
		return 0, true, nil
	case "--test":
		opts.IsTestEnabled = true
		return 0, true, nil
	case "--force", "-f":
		opts.IsForceCreate = true
		return 0, true, nil
	case "--backup":
		opts.IsBackupEnabled = true
		return 0, true, nil
	default:
		return 0, false, nil
	}
}

func parseValueFlags(opts *setup.IDESetupOptions, args []string, idx int) (int, bool, error) {
	if idx+1 >= len(args) {
		return 0, false, nil
	}

	switch args[idx] {
	case "--provider", "-p":
		p, err := enums.ParseIDEProviderType(args[idx+1])
		if err != nil {
			return 0, true, err
		}
		opts.Provider = p
		return 1, true, nil
	case "--model", "-m":
		opts.Model = args[idx+1]
		return 1, true, nil
	case "--base-url", "-u":
		opts.BaseURL = args[idx+1]
		return 1, true, nil
	case "--api-key", "-k":
		opts.APIKey = args[idx+1]
		return 1, true, nil
	default:
		return 0, false, nil
	}
}

func resolveTargetArg(raw string) enums.IDETargetType {
	clean := strings.ToLower(strings.TrimSpace(raw))
	switch clean {
	case "cline", "composer":
		return enums.TargetCline
	case "continue", "tab":
		return enums.TargetContinue
	case "zed":
		return enums.TargetZed
	default:
		return enums.TargetAllAI
	}
}

func applyPresetDefaults(opts *setup.IDESetupOptions) {
	preset, exists := setup.GetProviderPresets()[opts.Provider]
	if !exists {
		return
	}

	if opts.BaseURL == "" {
		opts.BaseURL = preset.DefaultBaseURL
	}
	if opts.Model == "" {
		opts.Model = preset.DefaultModel
	}
}
```

---

## 4. Setup Execution Pipeline (`cli/cmdide/setup_cmd.go`)

Orchestrate pre-flight verification, generators, and output rendering:

```go
package cmdide

import (
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdide/setup"
)

// IDESetupSummary details the results of configuring editors.
type IDESetupSummary struct {
	Provider              enums.IDEProviderType `json:"provider"`
	Model                 string                `json:"model"`
	BaseURL               string                `json:"baseUrl"`
	IsTestPassed          bool                  `json:"isTestPassed"`
	TestLatencyMs         int64                 `json:"testLatencyMs,omitempty"`
	ConfiguredTargets     []enums.IDETargetType `json:"configuredTargets"`
	BackupPaths           []string              `json:"backupPaths,omitempty"`
	IsDryRun              bool                  `json:"isDryRun"`
}

// RunIDESetup executes the setup workflow.
func RunIDESetup(args []string) error {
	opts, err := parseSetupOptions(args)
	if err != nil {
		return err
	}

	summary := IDESetupSummary{
		Provider: opts.Provider,
		Model:    opts.Model,
		BaseURL:  opts.BaseURL,
		IsDryRun: opts.IsDryRun,
	}

	if opts.IsTestEnabled {
		testRes := setup.TestConnection(setup.TestConnectionRequest{
			Provider: opts.Provider,
			BaseURL:  opts.BaseURL,
			APIKey:   opts.APIKey,
			Model:    opts.Model,
		})
		summary.IsTestPassed = testRes.IsSuccess
		summary.TestLatencyMs = testRes.LatencyMs
	}

	resolver := setup.NewPathResolver()
	if execErr := executeTargetGenerators(resolver, opts, &summary); execErr != nil {
		return execErr
	}

	if opts.IsJSON {
		return outputJSONSummary(summary)
	}

	renderTerminalSummary(summary)
	return nil
}

func executeTargetGenerators(r *setup.PathResolver, opts setup.IDESetupOptions, s *IDESetupSummary) error {
	if opts.Target == enums.TargetCline || opts.Target == enums.TargetAllAI {
		if err := setup.GenerateClineConfig(r, opts); err != nil {
			return err
		}
		s.ConfiguredTargets = append(s.ConfiguredTargets, enums.TargetCline)
	}

	if opts.Target == enums.TargetContinue || opts.Target == enums.TargetAllAI {
		if err := setup.GenerateContinueConfig(r, opts); err != nil {
			return err
		}
		s.ConfiguredTargets = append(s.ConfiguredTargets, enums.TargetContinue)
	}

	if opts.Target == enums.TargetZed || opts.Target == enums.TargetAllAI {
		if err := setup.GenerateZedConfig(r, opts); err != nil {
			return err
		}
		s.ConfiguredTargets = append(s.ConfiguredTargets, enums.TargetZed)
	}

	return nil
}

func outputJSONSummary(s IDESetupSummary) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}
```

---

## 5. Terminal Feedback & Visualizer (`cli/cmdide/setup_view.go`)

Render structured Lipgloss tables and instructions:

```go
package cmdide

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
)

func renderTerminalSummary(s IDESetupSummary) {
	fmt.Println()
	fmt.Println(titleStyle.Render("⚡ GitMap Cursor-Like IDE Setup Complete"))
	fmt.Printf("  • Provider:  %s (%s)\n", s.Provider, s.Model)
	fmt.Printf("  • Base URL:  %s\n", s.BaseURL)

	if s.IsTestPassed {
		fmt.Printf("  • Test Ping: %s (%dms)\n", successStyle.Render("✓ Connected"), s.TestLatencyMs)
	}

	fmt.Println()
	fmt.Println(titleStyle.Render("Configured IDE Targets:"))
	for _, target := range s.ConfiguredTargets {
		renderTargetRow(target, s.IsDryRun)
	}

	fmt.Println()
	fmt.Println(dimStyle.Render("Next steps: Launch your editor via 'code .' or 'zed .' to start coding with AI!"))
	fmt.Println()
}

func renderTargetRow(target fmt.Stringer, isDryRun bool) {
	modeTag := successStyle.Render("✓ Updated")
	if isDryRun {
		modeTag = dimStyle.Render("● Dry Run (Preview Only)")
	}

	fmt.Printf("  • %-12s %s\n", target.String(), modeTag)
}
```

---

## 6. Verification Criteria & Acceptance Tests

1. **`setup_flags_test.go`**:
   - Asserts `--provider kimi` correctly populates `opts.Provider = enums.ProviderKimi` and `opts.BaseURL = "https://api.moonshot.cn/v1"`.
   - Asserts `--dry-run`, `--json`, `--test`, and `--force` set corresponding positive booleans.
   - Asserts unknown provider fails validation with an `E1001` error code.
2. **`setup_cmd_test.go`**:
   - Executes `RunIDESetup` with `--dry-run` and validates no filesystem side-effects occur.
   - Executes `RunIDESetup` with `--json` and parses standard `IDESetupSummary` JSON structure from stdout.
3. **Coding Guidelines Gate**:
   - Every function <= 15 lines of executable code.
   - 100% positive booleans (`is*`, `has*`).
   - Zero nested ifs.
   - Errors bubble up to `RunIDESetup` without calling `os.Exit`.
