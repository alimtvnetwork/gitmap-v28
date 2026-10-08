# Component & Suggestion Engine Specification: Domain Models, Multi-Tier Resolution & CLI Interception

> **Spec Version:** 1.0.0  
> **Status:** Approved / Grounded  
> **Task Slug:** `241-cicd-enhancements-help-text-grouping-and-suggestion-engine-design`  
> **Target Path:** `02-spec/21-app/241-cicd-enhancements-help-text-grouping-and-suggestion-engine-design/02-component-and-suggestion-spec.md`  

---

## 1. Suggestion Domain Model (`cli/suggestion/`)

### 1.1 Core Domain Types
The centralized suggestion engine provides strongly typed objects to represent remedial actions, alternative commands, and typo corrections:

```go
package suggestion

// SuggestionCategory classifies the underlying intent of a suggested action.
type SuggestionCategory string

const (
    CategoryTypo        SuggestionCategory = "typo"        // Mistyped command or flag
    CategorySubcommand  SuggestionCategory = "subcommand"  // Valid child command of current verb
    CategoryFlag        SuggestionCategory = "flag"        // Supported flag matching user token
    CategoryRemediation SuggestionCategory = "remediation" // Prescriptive step to fix failure
    CategoryIntent      SuggestionCategory = "intent"      // Semantic alias (e.g. docker -> cluster)
    CategoryResource    SuggestionCategory = "resource"    // Repository, branch, or node target
)

// ActionType defines the execution mode of the suggested command.
type ActionType string

const (
    ActionAutoRun ActionType = "auto-run" // Safe to execute immediately if user passes --yes
    ActionManual  ActionType = "manual"   // Requires human parameter adjustment
    ActionCopy    ActionType = "copy"     // Copy-paste command candidate
)

// Suggestion represents an actionable CLI recommendation.
type Suggestion struct {
    Command     string             `json:"command"`
    Description string             `json:"description"`
    Category    SuggestionCategory `json:"category"`
    Confidence  float64            `json:"confidence"` // Normalized score: 0.0 to 1.0
    ActionType  ActionType         `json:"actionType"`
    Context     map[string]any     `json:"context,omitempty"`
}

// SuggestionGroup aggregates related suggestions under an explanatory banner.
type SuggestionGroup struct {
    Title       string       `json:"title"`
    Reason      string       `json:"reason"`
    InputToken  string       `json:"inputToken,omitempty"`
    Suggestions []Suggestion `json:"suggestions"`
}
```

---

## 2. Multi-Tier Resolution Engine

When an unknown token or failed command occurs, the `Engine` resolves candidate suggestions through a 4-tier pipeline:

```text
User Input: "gitmap scann"
   │
   ▼
┌────────────────────────────────────────────────────────┐
│ Tier 1: Canonical Alias Registry                       │
│ Matches known shorthand: s -> status, pe -> pull-error │
│ Match: None                                            │
└──────────────────────────┬─────────────────────────────┘
                           ▼
┌────────────────────────────────────────────────────────┐
│ Tier 2: Levenshtein Distance & QGram Fuzzy Matching     │
│ Computes edit distance against command registry:       │
│ "scann" -> "scan" (Distance = 1, Confidence = 0.95)    │
│ Match: FOUND                                           │
└──────────────────────────┬─────────────────────────────┘
                           ▼
┌────────────────────────────────────────────────────────┐
│ Tier 3: Subsequence / Prefix Matcher                   │
│ "cl" -> "clone", "rele" -> "release"                   │
└──────────────────────────┬─────────────────────────────┘
                           ▼
┌────────────────────────────────────────────────────────┐
│ Tier 4: Semantic Context & Synonym Dictionary          │
│ "docker" -> "cluster", "grep" -> "aum search"          │
│ "remove" -> "rm", "history" -> "log"                   │
└────────────────────────────────────────────────────────┘
```

### 2.1 Engine Interface
```go
type Engine interface {
    // ResolveCommand suggests valid commands for an unrecognized token.
    ResolveCommand(token string) SuggestionGroup

    // ResolveFlag suggests valid flags for a specific command context.
    ResolveFlag(cmdName, flagToken string) SuggestionGroup

    // ResolveRemediation creates context-aware fix suggestions for errors.
    ResolveRemediation(err error, ctx map[string]any) SuggestionGroup
}
```

---

## 3. Terminal Presentation & Renderers

### 3.1 Catppuccin Terminal Box Format (`RenderFormatBox`)
Replaces ugly unformatted error lines with elegant, bordered diagnostic cards:

```text
╭─ Did you mean? ────────────────────────────────────────────────────────╮
│ Command 'scann' is not recognized. Closest matches:                   │
│                                                                        │
│  • gitmap scan       (Confidence: 95%)                                 │
│    Fast repository scanner with parallel discovery                     │
│                                                                        │
│  • gitmap status     (Confidence: 45%)                                 │
│    Check working tree status across all repositories                   │
│                                                                        │
│ Run: gitmap scan --help for options                                   │
╰────────────────────────────────────────────────────────────────────────╯
```

### 3.2 Compact Inline Format (`RenderFormatCompact`)
Used for quick hints and minimal terminal sessions:
```text
Did you mean: gitmap scan (95%) | gitmap status (45%)
```

### 3.3 Structured JSON Format (`RenderFormatJSON`)
Enables AI agent orchestration and IDE tool calls to parse suggestions programmatically:
```json
{
  "title": "Unrecognized Command",
  "reason": "Token 'scann' not found in command registry",
  "inputToken": "scann",
  "suggestions": [
    {
      "command": "gitmap scan",
      "description": "Fast repository scanner with parallel discovery",
      "category": "typo",
      "confidence": 0.95,
      "actionType": "auto-run"
    }
  ]
}
```

---

## 4. Global CLI Dispatch Interceptor & Error Integration

### 4.1 Interception in `cli/cmd/root.go:runDispatch`
Currently, unrecognized commands trigger generic errors or fragmented `fmt.Printf` blocks. The centralized interceptor hooks dispatch routing:

```go
func runDispatch(args []string) error {
    if len(args) == 0 {
        return runDefaultCommand()
    }

    cmdName := args[0]
    handler, exists := registry.Lookup(cmdName)
    if !exists {
        // AUTOMATED INTERCEPTOR: No manual print statements
        group := suggestionEngine.ResolveCommand(cmdName)
        if len(group.Suggestions) > 0 {
            suggestionRenderer.RenderBox(os.Stderr, group)
        }
        return appfault.NewNotFoundError("unrecognized command: " + cmdName)
    }

    return handler(args[1:])
}
```

### 4.2 Integration with `appfault.AppError`
Extend GitMap's universal error envelope with structured suggestions:

```go
package appfault

import "github.com/alimtvnetwork/gitmap-v28/cli/suggestion"

type AppError struct {
    Code        string                  `json:"code"`
    Message     string                  `json:"message"`
    StackTrace  string                  `json:"stackTrace,omitempty"`
    Suggestions []suggestion.Suggestion `json:"suggestions,omitempty"`
}

// Fluent helper to attach suggestions to any error
func (e *AppError) WithSuggestions(items ...suggestion.Suggestion) *AppError {
    e.Suggestions = append(e.Suggestions, items...)
    return e
}
```

When `handleGlobalError` processes an `*AppError` with non-empty suggestions, it automatically delegates to `suggestionRenderer`, completely eliminating ad-hoc `fmt.Printf("   Did you mean: ...")` calls across all 36+ command packages.
