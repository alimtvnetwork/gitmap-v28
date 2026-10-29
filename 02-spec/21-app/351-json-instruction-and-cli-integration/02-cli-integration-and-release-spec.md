# 351.02 - JSON Instruction and CLI Integration: CLI Integration and Release Specification

- **Spec ID:** `02-spec/21-app/351-json-instruction-and-cli-integration/02-cli-integration-and-release-spec.md`
- **Subsystem:** `cli/cmd`, `cli/cmdspace`, `cli/linkexpander`, `go-email-sender/main`, `go-email-sender/internal/linkexpander`
- **Status:** APPROVED
- **Target Versions:** 
  - `gitmap`: `v6.629.0`
  - `go-email-sender` / `go-email-reader`: `v1.3.0`
- **Related Specs:**
  - `02-spec/21-app/351-json-instruction-and-cli-integration/01-json-instruction-schema-and-format-spec.md`
  - `02-spec/21-app/276-email-integration/01-scope-and-plan.md`
  - `02-spec/02-coding-guidelines/readme.md`
  - `02-spec/03-error-manage/readme.md`
  - `02-spec/01-spec-authoring-guide/06-cli-module-template.md`
- **Related Plans:**
  - `.ai-memory/plans/subtasks/351-json-instruction-and-cli-integration/02-cli-integration-and-release.md`

---

## 1. Executive Summary & Problem Context

In automated operational pipelines, enterprise systems frequently require processing batch instructions that combine **URL expansion and transformation** with **conditional email notifications**. Previously, developers and automated agents had to perform link expansion as a separate manual step via `gitmap space email expand` or `go-email-reader expand`, capture stdout, parse links, and then invoke individual email commands with explicit `--to`, `--subject`, and `--body` parameters.

This specification defines the architectural integration, CLI ergonomics, automated dispatch workflows, and multi-repository release ceremony for **JSON Instruction & CLI Integration**:

1. **Top-Level `gitmap send-email` CLI Exposure:** Expose `send-email` (and alias `sendemail`) directly at the top level of the `gitmap` CLI alongside `gitmap space send-email`, backed by the unified rich help displayer.
2. **Instruction Flag Expansion:** Introduce `--instruction <file.json>` and shorthand `-i <file.json>` across `gitmap send-email`, `gitmap space send-email`, and `gitmap space email expand`, enabling batch execution from structured JSON files.
3. **Validation Rule Relaxation:** Dynamically relax mandatory `--to` and `--subject` command-line requirements when `--instruction` is supplied, resolving recipients and subject templates directly from the input JSON document while allowing explicit flags to serve as global fallbacks or overrides.
4. **Dual-Repo Binary Integration:** Extend `go-email-sender` and `go-email-reader` `main.go` (`cmdSend` and `cmdExpand`) to ingest JSON instructions, resolve link expansions, interpolate body templates, and trigger automated SMTP dispatch via `internal/linkexpander`.
5. **Governed Dispatch Workflow:** Enforce strict recipient throttling (maximum 5 recipients per link), variable interpolation (`{{url}}`, `{{original}}`, `{{expanded}}`), and structured machine-readable JSON output reporting execution and delivery status per recipient.
6. **Minor Version Bump & Release Ceremony:** Standardize the release lifecycle across `gitmap` (`v6.629.0`) and `go-email-sender` (`v1.3.0`), including changelog updates, binary compilation, and semantic commit enforcement.

---

## 2. Architectural Topology & Component Interaction

The end-to-end integration operates across two core codebases: `gitmap` (the orchestrator and developer CLI) and `go-email-sender` / `go-email-reader` (the email delivery engine and link expansion runner).

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                   User / Automation CI                                 │
│      gitmap send-email -i batch.json --json   OR   go-email-reader send -i batch.json    │
└───────────────────────────────────────────┬─────────────────────────────────────────────┘
                                            │
                                            ▼
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                                   gitmap CLI Engine                                     │
│  cli/cmd/rootutility.go: Routes `send-email` & `sendemail` -> cmdspace.RunSendEmail     │
│  cli/cmd/helpcheck.go:   Binds CmdSpaceSendEmail -> cmdspace.SpaceSendEmailHelpDisplay  │
│  cli/cmdspace/send_email.go: Parses --instruction / -i, relaxes required --to / --subj  │
└───────────────────────────────────────────┬─────────────────────────────────────────────┘
                                            │
                    ┌───────────────────────┴───────────────────────┐
                    │ (If go-email-reader binary found)             │ (If built-in fallback)
                    ▼                                               ▼
┌──────────────────────────────────────┐        ┌─────────────────────────────────────────┐
│      go-email-reader / sender        │        │      gitmap Built-in SMTP Sender        │
│  main.go (cmdSend / cmdExpand)       │        │  cli/cmdspace/send_email_smtp.go        │
│  - Parses --instruction / -i         │        │  - Built-in SMTP fallback               │
│  - Reads JSON instruction payload    │        │  - Evaluates instruction entries        │
└──────────────────┬───────────────────┘        └───────────────────┬─────────────────────┘
                   │                                                │
                   └────────────────────────┬───────────────────────┘
                                            │
                                            ▼
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                              Link Expander & Dispatch Engine                            │
│                  (cli/linkexpander/ AND internal/linkexpander/)                         │
│  1. Ingest JSON Instruction -> Validate Schema & Max 5 Recipients Guard                │
│  2. Expand URLs via HTTP HEAD/GET & Domain Regex (e.g. serviceactivation.in)            │
│  3. Interpolate Variables: {{url}}, {{original}}, {{expanded}} into Body/Subject        │
│  4. Evaluate isAutoSend -> If true, dispatch SMTP email per recipient                   │
│  5. Emit Consolidated JSON Output (IsSuccess, URLs, Links, Per-Recipient Status)        │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. GitMap CLI Integration

### 3.1 Top-Level Exposure: `gitmap send-email`

Currently, `send-email` is accessible primarily under the `space` namespace (`gitmap space send-email`). To provide first-class ergonomics for developers and CI automation scripts, `send-email` and its alias `sendemail` MUST be registered at the top level of `gitmap`.

#### Implementation Contract (`cli/cmd/rootutility.go`)

In `cli/cmd/rootutility.go`, add a dedicated entry group or register directly in `utilityToolEntries()`:

```go
func utilityToolEntries() []dispatchEntry {
	return []dispatchEntry{
		// ... existing tool entries ...
		{[]string{"send-email", "sendemail"}, runSendEmailUtility},
		// ... remaining entries ...
	}
}

func runSendEmailUtility() error {
	checkHelp("send-email", argsTail())

	return cmdspace.RunSendEmail(argsTail())
}
```

- **Arguments Forwarding:** `argsTail()` strips the command token (`send-email` or `sendemail`) and forwards all trailing flags directly into `cmdspace.RunSendEmail`.
- **Error Sovereignty:** `runSendEmailUtility` returns the `error` from `cmdspace.RunSendEmail` directly, allowing the root dispatcher in `cli/cmd/root.go` to handle command audits and exit codes.

### 3.2 Rich Help Integration (`cli/cmd/helpcheck.go`)

Running `gitmap send-email --help` or `gitmap sendemail --help` MUST render the rich ANSI help displayer already implemented in `cmdspace.SpaceSendEmailHelpDisplay()`.

#### Registry Configuration

In `cli/cmd/helpcheck.go`, update `helpDisplayRegistry` and `helpDisplayAlias`:

```go
var helpDisplayRegistry = map[string]helpDisplayEntry{
	constants.CmdSpace: {
		build: cmdspace.SpaceHelpDisplay,
		subs: map[string]func() helpdisplay.Displayer{
			"common":                            cmdspace.SpaceCommonHelpDisplay,
			constants.CmdSpaceBackupBranch:      cmdspace.SpaceBackupBranchHelpDisplay,
			constants.CmdSpaceSendEmail:         cmdspace.SpaceSendEmailHelpDisplay,
			constants.CmdSpaceScheduleSendEmail: cmdspace.SpaceScheduleSendEmailHelpDisplay,
		},
	},
	constants.CmdSpaceSendEmail: {build: cmdspace.SpaceSendEmailHelpDisplay},
	constants.CmdScan:           {build: cmdscan.ScanHelpDisplay},
	constants.CmdPullAll:        {build: cmdpull.PullAllHelpDisplay},
}

var helpDisplayAlias = map[string]string{
	constants.CmdScanAlias:    constants.CmdScan,
	constants.CmdPullAllAlias: constants.CmdPullAll,
	"sendemail":               constants.CmdSpaceSendEmail,
}
```

When `tryInterceptCommandHelp("send-email", args)` or `tryInterceptCommandHelp("sendemail", args)` executes, `resolveHelpCommand` normalizes `"sendemail"` to `"send-email"`, locates `constants.CmdSpaceSendEmail` in `helpDisplayRegistry`, and invokes `SpaceSendEmailHelpDisplay()`.

### 3.3 Instruction Flag Expansion (`--instruction` / `-i`)

The `--instruction` flag and `-i` shorthand flag MUST be added to:
1. `gitmap send-email` (top-level)
2. `gitmap space send-email`
3. `gitmap space email expand`

#### Flag Constants & Options Structure (`cli/cmdspace/send_email.go`)

```go
const (
	flagSendEmailTo               = "--to"
	flagSendEmailSubject          = "--subject"
	flagSendEmailBody             = "--body"
	flagSendEmailFile             = "--file"
	flagSendEmailHTML             = "--html"
	flagSendEmailExpand           = "--expand-links"
	flagSendEmailExpandShort      = "--expand"
	flagSendEmailInstruction      = "--instruction"
	flagSendEmailInstructionShort = "-i"
	flagSendEmailJSON             = "--json"
	flagSendEmailHelp             = "--help"
	flagSendEmailHelpShort        = "-h"
)

type sendEmailOptions struct {
	To              string
	Subject         string
	Body            string
	File            string
	InstructionFile string
	IsHTML          bool
	IsExpandLinks   bool
	IsJSON          bool
	HasInstruction  bool
}
```

#### Argument Parsing Logic

Update `isKnownSendEmailFlag` and `assignSendEmailFlag`:

```go
func isKnownSendEmailFlag(name string) bool {
	switch name {
	case flagSendEmailTo, flagSendEmailSubject, flagSendEmailBody, flagSendEmailFile, flagSendEmailHTML,
		flagSendEmailExpand, flagSendEmailExpandShort, flagSendEmailInstruction, flagSendEmailInstructionShort,
		flagSendEmailJSON, flagSendEmailHelp, flagSendEmailHelpShort:
		return true
	}

	return false
}

func assignSendEmailFlag(name, val string, opts *sendEmailOptions) error {
	switch name {
	case flagSendEmailTo:
		opts.To = val
	case flagSendEmailSubject:
		opts.Subject = val
	case flagSendEmailBody:
		opts.Body = val
	case flagSendEmailFile:
		opts.File = val
	case flagSendEmailInstruction, flagSendEmailInstructionShort:
		opts.InstructionFile = val
		opts.HasInstruction = len(strings.TrimSpace(val)) > 0
	default:
		return apperror.NewValidationError(fmt.Sprintf("unknown flag: %s", name))
	}

	return nil
}
```

#### Instruction Support in `space email expand` (`cli/cmdspace/space_email.go`)

In `runSpaceEmailExpand(args []string)`:
1. Parse flags for `--instruction` or `-i`.
2. If `--instruction <file.json>` is provided without positional URLs, load and expand all URLs defined in the instruction document.
3. If `--json` is supplied, output the full expansion manifest as JSON; otherwise, format line-by-line expansions to standard output.

### 3.4 Relaxed Validation Rules

In standard `send-email` mode, `--to` and `--subject` are mandatory. When `--instruction` is specified:
1. **Mandatory Check Bypass:** If `opts.HasInstruction` is `true`, omit the requirement for CLI `--to` and `--subject`.
2. **Fallback / Override Behavior:** If CLI `--to` or `--subject` flags ARE provided alongside `--instruction`, they act as ambient defaults for any link instruction entry that does not specify its own recipients or subject.

```go
func validateEmailMessage(msg EmailMessage, hasInstruction bool) error {
	if hasInstruction {
		return nil
	}

	if len(msg.To) == 0 {
		return apperror.NewSimple("missing --to: usage: gitmap space send-email --to <addr[,addr]> --subject <s> --body <b>", "E_SEND_EMAIL_NO_TO")
	}

	if len(strings.TrimSpace(msg.Subject)) == 0 {
		return apperror.NewSimple("missing subject: pass --subject or provide it in --file", "E_SEND_EMAIL_NO_SUBJECT")
	}

	return nil
}
```

---

## 4. `go-email-sender` / `go-email-reader` CLI Integration

The companion binary `go-email-reader` (located in `d:\work\go-email-sender`) performs link expansion and SMTP transmission.

### 4.1 Flag Additions in `main.go`

In `cmdSend(args []string)`:
```go
func cmdSend(args []string) error {
	fs := flag.NewFlagSet("send", flag.ExitOnError)
	from := fs.String("from", "", "stored account email to send from (defaults to first stored account)")
	to := fs.String("to", "", "recipient(s), comma-separated")
	subject := fs.String("subject", "", "subject")
	body := fs.String("body", "", "plain-text body")
	file := fs.String("file", "", "read subject/body from file")
	isHTML := fs.Bool("html", false, "send body as text/html")
	templateName := fs.String("template", "", "use a stored template (see `template add`)")
	vars := fs.String("var", "", "template variables as k=v,k=v")
	expandLinks := fs.Bool("expand-links", false, "expand shortened or redirecting URLs in email body/subject before sending")
	expandAlias := fs.Bool("expand", false, "alias for --expand-links")
	instruction := fs.String("instruction", "", "path to JSON instruction file containing URLs and autoSend configurations")
	instructionAlias := fs.String("i", "", "alias for --instruction")
	asJSON := fs.Bool("json", false, "output send result as JSON")
	
	if err := fs.Parse(args); err != nil {
		return err
	}

	instructionPath := *instruction
	if instructionPath == "" {
		instructionPath = *instructionAlias
	}

	if instructionPath == "" && *to == "" {
		return fmt.Errorf("missing recipient: pass --to <address> or --instruction <file.json>")
	}
	
	// If instruction provided, route to instruction runner:
	if instructionPath != "" {
		return runSendInstruction(instructionPath, *from, *asJSON)
	}

	// ... continue standard send flow ...
}
```

In `cmdExpand(args []string)`:
```go
func cmdExpand(args []string) error {
	fs := flag.NewFlagSet("expand", flag.ExitOnError)
	instruction := fs.String("instruction", "", "path to JSON instruction file")
	instructionAlias := fs.String("i", "", "alias for --instruction")
	asJSON := fs.Bool("json", false, "output as JSON")
	
	if err := fs.Parse(args); err != nil {
		return err
	}

	instructionPath := *instruction
	if instructionPath == "" {
		instructionPath = *instructionAlias
	}

	if instructionPath != "" {
		return runExpandInstruction(instructionPath, *asJSON)
	}

	// ... continue standard positional expand flow ...
}
```

### 4.2 Link Expander Engine Wiring (`internal/linkexpander`)

Create `internal/linkexpander/instruction.go` to provide instruction file parsing, validation, and execution orchestration.

```go
package linkexpander

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// LinkInstruction represents a single link processing directive.
type LinkInstruction struct {
	URL         string   `json:"url"`
	IsAutoSend  bool     `json:"autoSend"`
	Recipients  []string `json:"recipients"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	IsHTML      bool     `json:"isHTML"`
}

// InstructionDocument represents the root structure of an instruction file.
type InstructionDocument struct {
	Version      string            `json:"version,omitempty"`
	Defaults     InstructionDefaults `json:"defaults,omitempty"`
	Instructions []LinkInstruction `json:"instructions"`
}

// InstructionDefaults provides fallback values for individual instructions.
type InstructionDefaults struct {
	Recipients []string `json:"recipients,omitempty"`
	Subject    string   `json:"subject,omitempty"`
	Body       string   `json:"body,omitempty"`
	IsHTML     bool     `json:"isHTML,omitempty"`
}

// ParseInstructionFile reads and validates a JSON instruction file from disk.
func ParseInstructionFile(filePath string) (InstructionDocument, error) {
	var doc InstructionDocument

	content, err := os.ReadFile(filePath)
	if err != nil {
		return doc, fmt.Errorf("failed to read instruction file %q: %w", filePath, err)
	}

	if err := json.Unmarshal(content, &doc); err != nil {
		return doc, fmt.Errorf("invalid instruction JSON format: %w", err)
	}

	if err := ValidateInstructionDoc(doc); err != nil {
		return doc, err
	}

	return doc, nil
}
```

### 4.3 SMTP Credential Retrieval & AutoSend Dispatch

When `autoSend` is `true`:
1. Retrieve decrypted SMTP credentials from the local SQLite key-value/keystore via `loadDecrypted(st, sender)`.
2. Format message using `maildef.Message`.
3. Dispatch email using `mail.SendMail(...)`.
4. Capture delivery response or error without halting subsequent recipients or instructions.

---

## 5. Automatic Email Dispatch Workflow

### 5.1 Link-to-Recipient Cardinality (Up to 5 Recipients Invariant)

To guard against spam, accidental mass-mailing, and API rate-limiting, **each link instruction is strictly constrained to a maximum of 5 recipients**.

#### Invariant Validation Rule:
```go
const MaxRecipientsPerLink = 5

func ValidateInstructionDoc(doc InstructionDocument) error {
	if len(doc.Instructions) == 0 {
		return fmt.Errorf("instruction file contains no instructions")
	}

	for idx, inst := range doc.Instructions {
		if strings.TrimSpace(inst.URL) == "" {
			return fmt.Errorf("instruction [%d]: url is required", idx)
		}

		effectiveRecipients := inst.Recipients
		if len(effectiveRecipients) == 0 && len(doc.Defaults.Recipients) > 0 {
			effectiveRecipients = doc.Defaults.Recipients
		}

		if len(effectiveRecipients) > MaxRecipientsPerLink {
			return fmt.Errorf("instruction [%d] for %q exceeds maximum of %d recipients (got %d)",
				idx, inst.URL, MaxRecipientsPerLink, len(effectiveRecipients))
		}
	}

	return nil
}
```

Any instruction file containing more than 5 recipients for a given URL will fail validation immediately with a clear explanatory error before any network requests or link expansions are performed.

### 5.2 Body Template Variable Interpolation

When preparing the subject and body for automated sending, the engine interpolates dynamic placeholders:

| Variable | Description | Example Replacement |
|:---|:---|:---|
| `{{url}}` | Active resolved destination URL (expanded if redirected, else original) | `https://serviceactivation.google.com/subscription/new/usr_123` |
| `{{original}}` | Raw URL before expansion | `https://serviceactivation.in/subscription/new/usr_123` |
| `{{expanded}}` | Fully expanded target URL | `https://serviceactivation.google.com/subscription/new/usr_123` |

#### Interpolation Helper Contract
```go
// InterpolateTemplate replaces {{url}}, {{original}}, and {{expanded}} placeholders.
func InterpolateTemplate(template string, link ExpandedLink) string {
	res := strings.ReplaceAll(template, "{{url}}", link.Expanded)
	res = strings.ReplaceAll(res, "{{original}}", link.Original)
	res = strings.ReplaceAll(res, "{{expanded}}", link.Expanded)

	return res
}
```

### 5.3 Clean JSON Reporting Schema

When invoked with `--json`, the command MUST produce deterministic, machine-parseable JSON reporting the status of every URL expansion and email recipient.

```json
{
  "isSuccess": true,
  "totalLinks": 2,
  "expandedCount": 2,
  "sentCount": 3,
  "failedCount": 0,
  "results": [
    {
      "originalUrl": "https://serviceactivation.in/subscription/new/sub_alpha",
      "expandedUrl": "https://serviceactivation.google.com/subscription/new/sub_alpha",
      "isExpanded": true,
      "isAutoSent": true,
      "subject": "Activation Ready: https://serviceactivation.google.com/subscription/new/sub_alpha",
      "recipients": [
        {
          "email": "dev1@example.com",
          "isSuccess": true,
          "status": "sent",
          "sentAt": "2026-10-29T10:20:00Z"
        },
        {
          "email": "ops@example.com",
          "isSuccess": true,
          "status": "sent",
          "sentAt": "2026-10-29T10:20:01Z"
        }
      ]
    },
    {
      "originalUrl": "https://example.com/shortlink",
      "expandedUrl": "https://example.com/dest?item=99",
      "isExpanded": true,
      "isAutoSent": true,
      "subject": "Direct Link: https://example.com/dest?item=99",
      "recipients": [
        {
          "email": "lead@example.com",
          "isSuccess": true,
          "status": "sent",
          "sentAt": "2026-10-29T10:20:02Z"
        }
      ]
    }
  ]
}
```

#### Go Struct Schema Contract (`cli/linkexpander/types.go` & `internal/linkexpander/types.go`)

```go
type RecipientSendStatus struct {
	Email     string `json:"email"`
	IsSuccess bool   `json:"isSuccess"`
	Status    string `json:"status"` // "sent", "failed", "skipped"
	SentAt    string `json:"sentAt,omitempty"`
	Error     string `json:"error,omitempty"`
}

type LinkProcessResult struct {
	OriginalURL   string                `json:"originalUrl"`
	ExpandedURL   string                `json:"expandedUrl"`
	IsExpanded    bool                  `json:"isExpanded"`
	IsAutoSent    bool                  `json:"isAutoSent"`
	Subject       string                `json:"subject,omitempty"`
	Recipients    []RecipientSendStatus `json:"recipients,omitempty"`
	Error         string                `json:"error,omitempty"`
}

type InstructionBatchResult struct {
	IsSuccess     bool                `json:"isSuccess"`
	TotalLinks    int                 `json:"totalLinks"`
	ExpandedCount int                 `json:"expandedCount"`
	SentCount     int                 `json:"sentCount"`
	FailedCount   int                 `json:"failedCount"`
	Results       []LinkProcessResult `json:"results"`
	Error         string              `json:"error,omitempty"`
}
```

### 5.4 Partial Failure Isolation

If one recipient or SMTP call encounters an error (e.g. invalid recipient address or transient connection timeout):
1. The failed recipient is recorded with `isSuccess: false`, `status: "failed"`, and the detailed `error` string.
2. The loop proceeds immediately to the remaining recipients and remaining links.
3. The overall `InstructionBatchResult.IsSuccess` is set to `true` if all succeeded, or `false` if any link/recipient failed, while returning code 0 in JSON mode so CI callers can parse the detailed error breakdown.

---

## 6. Coding Guidelines & Strict Compliance Enforcements

Every line of code added for this feature MUST strictly adhere to repository guidelines:

### 6.1 Error Management & Leaf Sovereignty
- **Leaf Functions:** Any function in `cli/cmdspace`, `cli/linkexpander`, or `internal/linkexpander` that returns an `error` MUST return it directly via `return err` or `return apperror.Wrap(err, ...)`. Under no circumstances may a leaf function call `cliexit.HandleError`, `cliexit.Exit`, or `os.Exit`.
- **Top-Level Exit Handling:** Only `dispatch` in `cli/cmd/root.go` and `main()` handle process termination and DB audit logging (`finishCommandAudit`).
- **Typed Error Enums:** Use `constants.ExitUsageError` and `constants.ExitGeneralError` rather than raw integer codes.

### 6.2 Boolean Naming Discipline
- All boolean variables, struct fields, and functions MUST start with `is` or `has` (`isAutoSend`, `isHTML`, `isExpanded`, `isSuccess`, `hasInstruction`, `hasChanges`).
- TOTAL BAN on prefixes: `can`, `should`, `was`, `will`, and negative booleans (`isNotAutoSend`).

### 6.3 15-Line Function Limit & Guard Clause Flattening
- Every function MUST remain <= 15 lines of executable code (<= 8 lines preferred).
- Eliminate all nested conditionals using guard clauses and early returns.

```go
// CORRECT: Guard clause with early return, <= 8 lines
func validateRecipientCount(recipients []string) error {
	if len(recipients) > MaxRecipientsPerLink {
		return apperror.NewValidationError(fmt.Sprintf("exceeds max %d recipients", MaxRecipientsPerLink))
	}

	return nil
}
```

### 6.4 Whitespace and Return Formatting
- Exactly one blank line after every closing brace `}`.
- Exactly one blank line before every `return` statement.

---

## 7. Testing & Verification Matrix

| Test Case ID | Test Target | Verification Description | Expected Outcome |
|:---|:---|:---|:---|
| **TC-CLI-01** | `gitmap send-email --help` | Invoke help on top-level command | Displays rich ANSI HelpDisplay with description, flags, and examples |
| **TC-CLI-02** | `gitmap sendemail -h` | Invoke help on alias | Displays identical rich HelpDisplay as `send-email` |
| **TC-CLI-03** | Flag Parsing `-i` / `--instruction` | Pass instruction file to `gitmap send-email` | Flag accepted, `HasInstruction` parsed as `true` |
| **TC-VAL-01** | Relaxed Flag Validation | Pass `--instruction` without `--to` or `--subject` | Validation succeeds; error `E_SEND_EMAIL_NO_TO` is suppressed |
| **TC-VAL-02** | Recipient Throttling (> 5) | Instruction JSON with 6 recipients for one URL | Fails pre-flight validation with recipient limit error |
| **TC-EXP-01** | URL Expansion in Batch | Instruction JSON with multiple URLs | All URLs expanded and mapped to destination endpoints |
| **TC-TPL-01** | Variable Interpolation | Subject/body using `{{url}}`, `{{original}}`, `{{expanded}}` | Tokens accurately replaced with resolved URLs |
| **TC-OUT-01** | JSON Output Contract | Run with `--json` | Outputs well-formed JSON conforming to `InstructionBatchResult` schema |
| **TC-DEL-01** | Delegate Execution | Run with `go-email-reader` binary installed | Command correctly delegates execution to binary |

---

## 8. Minor Version Bump & Release Ceremony

Upon completing implementation, tests, and guideline compliance verification, execute the standardized minor version bump and release ceremony across both repositories.

### 8.1 `gitmap` Release Lifecycle (`v6.629.0`)

1. **Version Constant Update:**
   Update `cli/constants/constants.go`:
   ```go
   var Version = "6.629.0"
   ```
2. **Changelog Maintenance:**
   Update `CHANGELOG.md` with:
   - Feature: Registered top-level `gitmap send-email` and alias `sendemail` with rich help integration.
   - Feature: Added `--instruction` (`-i`) flag to `send-email` and `space email expand`.
   - Feature: Relaxed mandatory `--to`/`--subject` flag checks when `--instruction` is provided.
   - Feature: Implemented batch URL expansion with automated email dispatch up to 5 recipients per link.
3. **Compilation & Local Verification:**
   Build and smoke-test:
   ```pwsh
   go build -ldflags="-s -w" -o bin/gitmap.exe .
   bin\gitmap.exe version
   bin\gitmap.exe send-email --help
   ```
4. **Semantic Commit & Tagging:**
   ```pwsh
   gitmap cpr "v6.629.0 JSON instruction engine and top-level send-email integration"
   git tag v6.629.0
   ```

### 8.2 `go-email-sender` Release Lifecycle (`v1.3.0`)

1. **Config & Version Manifest Update:**
   Update `d:\work\go-email-sender\run.config.json`:
   ```json
   {
     "name": "go-email-sender-and-reader",
     "projectName": "go-email-sender",
     "version": "1.3.0",
     ...
   }
   ```
2. **Binary Build & Packaging:**
   Compile fresh binaries for both `go-email-reader` and `go-email-sender`:
   ```pwsh
   cd d:\work\go-email-sender
   go build -ldflags="-s -w" -o bin/go-email-reader.exe .
   go build -ldflags="-s -w" -o bin/go-email-sender.exe .
   ```
3. **Commit & Tagging:**
   ```pwsh
   git add run.config.json main.go internal/
   git commit -m "chore(release): v1.3.0 JSON instruction support and autoSend link expander"
   git tag v1.3.0
   ```

---

## 9. Acceptance Criteria & Verification Checklist

- [ ] `gitmap send-email` and alias `sendemail` route cleanly to `cmdspace.RunSendEmail(argsTail())`.
- [ ] `gitmap send-email --help` and `gitmap sendemail --help` display the rich HelpDisplay menu.
- [ ] `--instruction` and `-i` flags are parsed across `gitmap send-email`, `gitmap space send-email`, and `gitmap space email expand`.
- [ ] Mandatory `--to` and `--subject` validation checks are relaxed when `--instruction` is supplied.
- [ ] `go-email-sender` / `go-email-reader` `cmdSend` and `cmdExpand` accept `--instruction` and `-i`.
- [ ] `internal/linkexpander` parses JSON instructions and validates the <= 5 recipients per link ceiling.
- [ ] Template variable interpolation accurately resolves `{{url}}`, `{{original}}`, and `{{expanded}}`.
- [ ] Automated email dispatch functions reliably when `autoSend: true` is present.
- [ ] `--json` produces structured JSON conforming to `InstructionBatchResult`.
- [ ] Code strictly follows 15-line function limits, zero nested ifs, positive booleans, and leaf error return sovereignty.
- [ ] `gitmap` is released at `v6.629.0` and `go-email-sender` at `v1.3.0`.
