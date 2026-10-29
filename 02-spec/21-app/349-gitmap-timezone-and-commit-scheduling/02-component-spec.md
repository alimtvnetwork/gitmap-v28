# 349 - GitMap Timezone and Commit Scheduling: Component Specification

- **Spec ID:** `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/02-component-spec.md`
- **Subsystem:** `cli/timeutil`, `cli/cmdcommit`, `cli/cmdui`, `src/components/settings`, `src/pages`, `cli/helpdoc`
- **Status:** APPROVED
- **Target Version:** `v6.630.0`
- **Related Specs:**
  - `02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/01-architecture-spec.md`
  - `02-spec/13-generic-cli/14-date-formatting.md`
  - `02-spec/02-coding-guidelines/readme.md`
  - `02-spec/03-error-manage/readme.md`
- **Related Plans:**
  - `.ai-memory/plans/subtasks/349-gitmap-timezone-and-commit-scheduling/01-timezone-scheduling-architecture.md`
  - `.ai-memory/plans/subtasks/349-gitmap-timezone-and-commit-scheduling/02-ui-help-and-settings.md`

---

## 1. Executive Summary & Component Scope

This component specification establishes the programmatic contracts, interface definitions, data structures, and concrete implementation boundaries for **GitMap Timezone Resolution and Commit Window Scheduling**.

### 1.1 Scope Boundaries & Core Deliverables
The implementation spans the CLI engine, terminal UX, embedded Go UI, and React web dashboard across seven discrete component modules:
1. **Timezone Resolution Engine (`cli/timeutil/timezone.go`):** Resolves target timezones across IANA database identifiers, UTC aliases, local machine settings, and ISO numeric offsets (`+08:00`, `-0400`), embedding Go's `tzdata` package for offline resiliency.
2. **Commit Window Policy & Scheduling Helper (`cli/timeutil/schedule.go`):** Computes weekday morning cutoff enforcement (`< 08:00 AM`), weekend 24-hour unconstrained passes, and next available morning window targets (`07:30 AM`).
3. **Commit Flag Parser & Domain Metadata (`cli/cmdcommit/commit_flags.go`):** Parses `--timezone` / `-tz`, `--schedule` / `-sc`, and `--auto-window` with zero regex and ambient fallback enrichment.
4. **Interactive Terminal Prompting (`cli/cmdcommit/commit_prompt.go`):** Renders interactive prompts for omitted commit messages, missing timezones, and closed window resolution dialogs when running in interactive TTY sessions.
5. **Embedded Go Settings UI Markup (`cli/cmdui/`):** Manages user preferences for `commit.timezone`, `commit.schedule_active`, and `commit.weekday_cutoff` within the embedded HTML/JS server and `ui_settings.json`.
6. **React Web Dashboard Settings UI (`src/components/settings/CommitScheduleCard.tsx` and `src/pages/Settings.tsx`):** Provides a visual settings card with searchable timezone dropdowns, scheduling toggles, cutoff hour selectors, and real-time timezone status previews.
7. **Terminal Rich Help & Documentation (`cli/cmdcommit/commit_help_menu.go` and `cli/helpdoc/commit.md`):** Documents new scheduling and timezone capabilities in ANSI two-column help menus and markdown documentation.

### 1.2 Non-Negotiable Invariants & Guardrails
- **Strict Relative Git Paths:** All file references, documentation links, and code paths MUST use relative paths (`02-spec/...`, `.ai-memory/...`, `cli/...`, `src/...`). Absolute paths (`D:\...`, `C:\...`) are strictly prohibited.
- **Positive Booleans Only:** All structural flags, conditional switches, and type contracts MUST begin with `is` or `has` (`isAllowed`, `isWeekend`, `isWeekday`, `isMorningAllowed`, `isAutoWindow`, `hasCustomTz`, `isValid`). Banned prefixes: `can`, `should`, `was`, `will`, and negative booleans (`isNotAllowed`, `isBlocked`).
- **15-Line Function Limit:** Functions MUST NOT exceed 15 lines of executable code (<= 8 lines strongly preferred). Complex operations must be decomposed into single-responsibility sub-helpers.
- **Leaf Error Return Sovereignty:** Leaf functions returning `error` or `(Result, error)` MUST return the error directly. Leaf functions MUST NEVER call exit handlers (`cliexit.HandleError`) or `os.Exit`. Outer command dispatchers in `cli/cmd/root.go` handle termination and DB audit capture.
- **Zero Nested Ifs:** Flatten all conditional branching with guard clauses and early returns.
- **Return & Brace Styling:** Exactly one blank line after a closing brace `}` and exactly one blank line before a `return` statement.
- **Universal UTC Persistence:** All underlying database timestamps and Git commit metadata (`GIT_AUTHOR_DATE`, `GIT_COMMITTER_DATE`) are formatted in UTC (`t.UTC()`). Target timezones are utilized exclusively for business policy evaluation and user display formatting.

---

## 2. Component Dependency Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            User Invocations                                 │
│  - CLI: gitmap commit "feat" --tz "America/New_York" --auto-window          │
│  - Web UI: http://localhost:49152/settings (CommitScheduleCard)             │
└──────────────────────────────────────┬──────────────────────────────────────┘
                                       │
                                       ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Target Timezone Resolver                            │
│                     cli/timeutil/timezone.go                                │
│   - ResolveLocation(raw string) (*time.Location, error)                     │
│   - IANA database lookup (_ "time/tzdata" embedded)                         │
│   - Shorthand aliases ("utc", "est", "pst", "local")                        │
│   - Numeric offset parser (+08:00, -0400, +0530)                            │
└──────────────────┬──────────────────────────────────────┬───────────────────┘
                   │                                      │
                   ▼                                      ▼
┌──────────────────────────────────────┐ ┌────────────────────────────────────┐
│      Commit Window Evaluator         │ │      CLI Argument Extractor        │
│      cli/timeutil/schedule.go        │ │      cli/cmdcommit/commit_flags.go │
│  - IsCommitWindowAllowed(...)        │ │  - Parse --timezone, -tz           │
│  - IsWeekend(...) / IsWeekday(...)   │ │  - Parse --schedule, -sc           │
│  - IsMorningAllowed(t, cutoff)       │ │  - Parse --auto-window             │
│  - CalculateNextWindow(...)          │ │  - Parse --non-interactive, -y     │
└──────────────────┬───────────────────┘ └────────────────┬───────────────────┘
                   │                                      │
                   └──────────────────┬───────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Execution Dispatch & Prompting Engine                    │
│                      cli/cmdcommit/commit_prompt.go                         │
│   - PromptCommitMessage: fallback when commit message is omitted            │
│   - PromptTimezone: interactive selection when timezone is needed           │
│   - PromptWindowResolution: interactive dialog when window is closed        │
└──────────────────┬──────────────────────────────────────┬───────────────────┘
                   │                                      │
                   ▼                                      ▼
┌──────────────────────────────────────┐ ┌────────────────────────────────────┐
│      Embedded Go UI Settings         │ │      React Settings Dashboard      │
│      cli/cmdui/                      │ │      src/components/settings/      │
│  - commit.timezone                   │ │  - CommitScheduleCard.tsx          │
│  - commit.schedule_active            │ │  - src/pages/Settings.tsx          │
│  - commit.weekday_cutoff             │ │  - Live target timezone clock      │
└──────────────────┬───────────────────┘ └────────────────┬───────────────────┘
                   │                                      │
                   └──────────────────┬───────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                    Terminal Rich Help & Documentation                       │
│   - cli/cmdcommit/commit_help_menu.go: ANSI two-column help menu            │
│   - cli/helpdoc/commit.md: Markdown reference guide and syntax              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Detailed Component Contracts

### 3.1 Component 1: `cli/timeutil/timezone.go`

- **Target Path:** `cli/timeutil/timezone.go`
- **Package:** `package timeutil`
- **Imports:** `time`, `strings`, `strconv`, `_ "time/tzdata"`, `github.com/alimtvnetwork/gitmap-v28/cli/apperror`, `github.com/alimtvnetwork/gitmap-v28/cli/constants`
- **Architectural Role:** Pure leaf timezone resolution engine providing deterministic parsing of IANA timezone names, common aliases, system local defaults, and fixed numeric offset strings.

#### 3.1.1 Exported API Signatures

```go
// ResolveLocation resolves a raw string into a valid *time.Location.
// Supports IANA identifiers ("America/New_York"), UTC, Local ("" or "local"),
// and numeric offsets ("+08:00", "-0400", "+0530", "-07").
func ResolveLocation(raw string) (*time.Location, error)
```

#### 3.1.2 Internal Helper Decomposition (<= 15 Lines per Function)

```go
package timeutil

import (
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ResolveLocation resolves a raw string into a valid *time.Location.
func ResolveLocation(raw string) (*time.Location, error) {
	clean := strings.TrimSpace(raw)
	if clean == "" || strings.EqualFold(clean, "local") {
		return time.Local, nil
	}

	norm := normalizeTimezoneAlias(clean)
	if loc, hasOffset, err := parseFixedOffsetLocation(norm); hasOffset {
		return loc, err
	}

	loc, err := time.LoadLocation(norm)
	if err != nil {
		return nil, apperror.Wrap(err, "invalid timezone: "+clean, constants.ExitUsageError)
	}

	return loc, nil
}

func normalizeTimezoneAlias(s string) string {
	lower := strings.ToLower(s)
	switch lower {
	case "utc", "zulu", "z":
		return "UTC"
	case "est", "edt":
		return "America/New_York"
	case "cst", "cdt":
		return "America/Chicago"
	case "mst", "mdt":
		return "America/Denver"
	case "pst", "pdt":
		return "America/Los_Angeles"
	case "gmt":
		return "Etc/GMT"
	default:
		return s
	}
}

func parseFixedOffsetLocation(raw string) (*time.Location, bool, error) {
	if !strings.HasPrefix(raw, "+") && !strings.HasPrefix(raw, "-") {
		return nil, false, nil
	}

	totalSec, isParsed := parseOffsetSeconds(raw)
	if !isParsed {
		return nil, true, apperror.New("invalid offset format: "+raw, constants.ExitUsageError)
	}

	return time.FixedZone(raw, totalSec), true, nil
}

func parseOffsetSeconds(s string) (int, bool) {
	isPositive := strings.HasPrefix(s, "+")
	body := strings.TrimPrefix(strings.TrimPrefix(s, "+"), "-")
	body = strings.ReplaceAll(body, ":", "")

	hours, mins, isValid := splitHourMinute(body)
	if !isValid {
		return 0, false
	}

	sec := (hours * 3600) + (mins * 60)
	if !isPositive {
		sec = -sec
	}

	return sec, true
}

func splitHourMinute(s string) (int, int, bool) {
	if len(s) == 2 {
		h, err := strconv.Atoi(s)
		return h, 0, err == nil && h >= 0 && h <= 23
	}

	if len(s) == 4 {
		h, errH := strconv.Atoi(s[:2])
		m, errM := strconv.Atoi(s[2:])
		isValid := errH == nil && errM == nil && h >= 0 && h <= 23 && m >= 0 && m <= 59
		return h, m, isValid
	}

	return 0, 0, false
}
```

---

### 3.2 Component 2: `cli/timeutil/schedule.go`

- **Target Path:** `cli/timeutil/schedule.go`
- **Package:** `package timeutil`
- **Imports:** `time`
- **Architectural Role:** Pure business policy evaluation functions checking commit timing against the weekday morning restriction (< 08:00 AM) and weekend 24h unconstrained pass.

#### 3.2.1 Exported API Signatures

```go
// IsCommitWindowAllowed checks whether the timestamp in the given location satisfies commit policies.
func IsCommitWindowAllowed(t time.Time, loc *time.Location) bool

// IsWeekend returns true if the timestamp falls on Saturday or Sunday in the target location.
func IsWeekend(t time.Time) bool

// IsWeekday returns true if the timestamp falls between Monday and Friday in the target location.
func IsWeekday(t time.Time) bool

// IsMorningAllowed checks whether the hour is strictly less than cutoffHour (e.g. 8).
func IsMorningAllowed(t time.Time, cutoffHour int) bool

// CalculateNextWindow computes the earliest subsequent open commit slot.
func CalculateNextWindow(t time.Time, loc *time.Location, cutoffHour int) time.Time
```

#### 3.2.2 Concrete Implementation Decomposition (<= 15 Lines)

```go
package timeutil

import (
	"time"
)

// DefaultWeekdayCutoffHour specifies the 08:00 AM boundary.
const DefaultWeekdayCutoffHour = 8

// IsCommitWindowAllowed checks whether the timestamp in loc satisfies commit policies.
func IsCommitWindowAllowed(t time.Time, loc *time.Location) bool {
	target := t.In(loc)
	if IsWeekend(target) {
		return true
	}

	return IsMorningAllowed(target, DefaultWeekdayCutoffHour)
}

// IsWeekend returns true if the day is Saturday or Sunday.
func IsWeekend(t time.Time) bool {
	w := t.Weekday()

	return w == time.Saturday || w == time.Sunday
}

// IsWeekday returns true if the day is Monday through Friday.
func IsWeekday(t time.Time) bool {
	return !IsWeekend(t)
}

// IsMorningAllowed checks whether the hour is strictly less than cutoffHour.
func IsMorningAllowed(t time.Time, cutoffHour int) bool {
	return t.Hour() < cutoffHour
}

// CalculateNextWindow computes the earliest valid target window time.
func CalculateNextWindow(t time.Time, loc *time.Location, cutoffHour int) time.Time {
	target := t.In(loc)
	if IsWeekend(target) {
		return target
	}
	if IsMorningAllowed(target, cutoffHour) {
		return target
	}

	return computeNextMorningTarget(target)
}

func computeNextMorningTarget(target time.Time) time.Time {
	w := target.Weekday()
	if w == time.Friday {
		return time.Date(target.Year(), target.Month(), target.Day()+1, 0, 0, 0, 0, target.Location())
	}

	return time.Date(target.Year(), target.Month(), target.Day()+1, 7, 30, 0, 0, target.Location())
}
```

---

### 3.3 Component 3: `cli/cmdcommit/commit_flags.go`

- **Target Path:** `cli/cmdcommit/commit_flags.go`
- **Package:** `package cmdcommit`
- **Imports:** `strings`, `github.com/alimtvnetwork/gitmap-v28/cli/timeutil`
- **Architectural Role:** Command flag extraction and ambient provenance resolution for `--timezone` / `-tz`, `--schedule` / `-sc`, and `--auto-window`.

#### 3.3.1 Domain Model Extensions in `cli/cmdcommit/commit_types.go`

```go
// CommitMetadata encapsulates operational flags and ambient AI provenance.
type CommitMetadata struct {
	TaskSlug       string
	AgentID        string
	HeatmapID      string
	AIName         string
	SpecPath       string
	Author         string
	TargetTimezone string
	ScheduleSpec   string
	IsAIMode       bool
	HasPush        bool
	IsDryRun       bool
	IsAutoWindow   bool
	HasTimezone    bool
	HasSchedule    bool
	IsNonInteractive bool
	HasForceWindow bool
}
```

#### 3.3.2 Flag Constants Extension

```go
const (
	FlagTimezoneLong  = "--timezone"
	FlagTimezoneShort = "-tz"
	FlagScheduleLong  = "--schedule"
	FlagScheduleShort = "-sc"
	FlagAutoWindow    = "--auto-window"
	FlagNonInteractive = "--non-interactive"
	FlagForceWindow   = "--force-window"

	EnvGitmapTimezone = "GITMAP_TIMEZONE"
	EnvGitmapSchedule = "GITMAP_SCHEDULE"
)
```

#### 3.3.3 Parsing Helper Routines (<= 15 Lines)

```go
func isTimezoneFlag(arg, val string, meta *CommitMetadata) bool {
	if arg != FlagTimezoneLong && arg != FlagTimezoneShort {
		return false
	}
	meta.TargetTimezone = strings.TrimSpace(val)
	meta.HasTimezone = true

	return true
}

func isScheduleFlag(arg, val string, meta *CommitMetadata) bool {
	if arg != FlagScheduleLong && arg != FlagScheduleShort {
		return false
	}
	meta.ScheduleSpec = strings.TrimSpace(val)
	meta.HasSchedule = true

	return true
}

func isSchedulingStandaloneFlag(arg string, meta *CommitMetadata) bool {
	switch arg {
	case FlagAutoWindow:
		meta.IsAutoWindow = true
		return true
	case FlagNonInteractive, "-y":
		meta.IsNonInteractive = true
		return true
	case FlagForceWindow:
		meta.HasForceWindow = true
		return true
	default:
		return false
	}
}
```

---

### 3.4 Component 4: `cli/cmdcommit/commit_prompt.go`

- **Target Path:** `cli/cmdcommit/commit_prompt.go`
- **Package:** `package cmdcommit`
- **Imports:** `bufio`, `fmt`, `os`, `strings`, `github.com/alimtvnetwork/gitmap-v28/cli/timeutil`, `github.com/charmbracelet/lipgloss`
- **Architectural Role:** Interactive CLI prompting routines for developer terminal sessions. Prompts for commit messages when omitted, resolves target timezones interactively, and displays structured Lipgloss resolution cards when a closed window is encountered.

#### 3.4.1 Type Declarations

```go
// WindowResolutionChoice represents user selection during closed window prompts.
type WindowResolutionChoice int

const (
	ChoiceAutoSchedule WindowResolutionChoice = 1
	ChoiceChangeTimezone WindowResolutionChoice = 2
	ChoiceForceCommit WindowResolutionChoice = 3
	ChoiceCancel WindowResolutionChoice = 4
)
```

#### 3.4.2 Exported API Signatures

```go
// PromptCommitMessage interactively asks the user for a commit message if none was provided.
func PromptCommitMessage(initialMsg string) (string, error)

// PromptTargetTimezone prompts the user for a timezone string or confirms current default.
func PromptTargetTimezone(defaultTz string) (string, error)

// PromptWindowResolution renders an interactive menu when the commit window is closed.
func PromptWindowResolution(eval timeutil.CommitWindowEvaluation) (WindowResolutionChoice, error)
```

#### 3.4.3 Concrete Implementation Decomposition (<= 15 Lines)

```go
package cmdcommit

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/timeutil"
	"github.com/charmbracelet/lipgloss"
)

// PromptCommitMessage prompts for commit message if empty.
func PromptCommitMessage(initialMsg string) (string, error) {
	clean := strings.TrimSpace(initialMsg)
	if clean != "" {
		return clean, nil
	}

	fmt.Print("Enter commit message: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", apperror.Wrap(err, "failed reading commit message", constants.ExitGeneralError)
	}

	msg := strings.TrimSpace(input)
	if msg == "" {
		return "", apperror.New("commit message cannot be empty", constants.ExitUsageError)
	}

	return msg, nil
}

// PromptTargetTimezone prompts for target timezone.
func PromptTargetTimezone(defaultTz string) (string, error) {
	fmt.Printf("Enter target timezone [%s]: ", defaultTz)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", apperror.Wrap(err, "failed reading timezone", constants.ExitGeneralError)
	}

	res := strings.TrimSpace(input)
	if res == "" {
		return defaultTz, nil
	}

	return res, nil
}

// PromptWindowResolution renders Lipgloss dialog for closed window.
func PromptWindowResolution(eval timeutil.CommitWindowEvaluation) (WindowResolutionChoice, error) {
	renderWindowClosedBanner(eval)
	renderWindowChoices()

	fmt.Print("Select option [1-4] (default 1): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return ChoiceCancel, apperror.Wrap(err, "failed reading option", constants.ExitGeneralError)
	}

	return parseResolutionChoice(strings.TrimSpace(input)), nil
}

func renderWindowClosedBanner(eval timeutil.CommitWindowEvaluation) {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFA500")).
		BorderStyle(lipgloss.RoundedBorder()).
		Padding(0, 1)

	msg := fmt.Sprintf("⚠️ COMMIT WINDOW CLOSED\nTarget: %s\nNext Window: %s",
		eval.TargetLocation.String(),
		eval.NextAvailableWindow.Format("02-Jan-2006 03:04 PM"))

	fmt.Println(style.Render(msg))
}

func renderWindowChoices() {
	fmt.Println("  [1] Auto-schedule for next morning window (07:30 AM)")
	fmt.Println("  [2] Change target timezone")
	fmt.Println("  [3] Force commit immediately (audit logged)")
	fmt.Println("  [4] Cancel commit operation")
}

func parseResolutionChoice(in string) WindowResolutionChoice {
	switch in {
	case "2":
		return ChoiceChangeTimezone
	case "3":
		return ChoiceForceCommit
	case "4":
		return ChoiceCancel
	default:
		return ChoiceAutoSchedule
	}
}
```

---

### 3.5 Component 5: Embedded Go UI Settings Markup (`cli/cmdui/`)

- **Target Path:** `cli/cmdui/ui_assets_markup_settings_b.go`, `cli/cmdui/ui_assets_js_settings.go`, `cli/cmdui/ui_server_settings.go`
- **Package:** `package cmdui`
- **Architectural Role:** Embedded Go UI components rendering settings forms for timezone and commit scheduling, storing values in `SettingsData.Attributes`:
  - `commit.timezone`: e.g. `"America/New_York"`, `"UTC"`, `"Local"`.
  - `commit.schedule_active`: `"true"` / `"false"`.
  - `commit.weekday_cutoff`: `"8"`.

#### 3.5.1 HTML Markup Template Snippet (`ui_assets_markup_settings_b.go`)

```html
<div class="settings-card" id="card-commit-schedule">
  <div class="card-header">
    <div class="card-title">
      <span class="card-icon">🕒</span>
      <h3>Commit Timezone & Schedule</h3>
    </div>
    <span class="badge" id="schedule-status-badge">Active</span>
  </div>
  <div class="card-body">
    <div class="form-row">
      <label for="commit-tz-input">Target Timezone</label>
      <input type="text" id="commit-tz-input" placeholder="e.g. America/New_York, UTC, Local" class="form-input" />
      <span class="helper-text">Used to evaluate weekday commit window boundaries.</span>
    </div>
    <div class="form-row checkbox-row">
      <input type="checkbox" id="commit-schedule-active" class="form-checkbox" />
      <label for="commit-schedule-active">Enforce Weekday Morning Restriction (&lt; 08:00 AM)</label>
    </div>
    <div class="form-row">
      <label for="commit-cutoff-hour">Weekday Cutoff Hour (24h format)</label>
      <input type="number" id="commit-cutoff-hour" min="1" max="12" value="8" class="form-input" />
      <span class="helper-text">Commits on Monday through Friday must occur before this hour.</span>
    </div>
  </div>
</div>
```

#### 3.5.2 JavaScript State Hydration & Save Snippet (`ui_assets_js_settings.go`)

```javascript
function loadCommitScheduleSettings(settings) {
  const attrs = settings.attributes || {};
  const tzInput = document.getElementById('commit-tz-input');
  const activeCheck = document.getElementById('commit-schedule-active');
  const cutoffInput = document.getElementById('commit-cutoff-hour');

  if (tzInput) tzInput.value = attrs['commit.timezone'] || 'Local';
  if (activeCheck) activeCheck.checked = attrs['commit.schedule_active'] === 'true';
  if (cutoffInput) cutoffInput.value = attrs['commit.weekday_cutoff'] || '8';
}

function collectCommitScheduleSettings(settings) {
  if (!settings.attributes) settings.attributes = {};
  const tzInput = document.getElementById('commit-tz-input');
  const activeCheck = document.getElementById('commit-schedule-active');
  const cutoffInput = document.getElementById('commit-cutoff-hour');

  if (tzInput) settings.attributes['commit.timezone'] = tzInput.value.trim() || 'Local';
  if (activeCheck) settings.attributes['commit.schedule_active'] = activeCheck.checked ? 'true' : 'false';
  if (cutoffInput) settings.attributes['commit.weekday_cutoff'] = cutoffInput.value.trim() || '8';
}
```

---

### 3.6 Component 6: React Settings UI (`src/components/settings/CommitScheduleCard.tsx` & `src/pages/Settings.tsx`)

- **Component Path:** `src/components/settings/CommitScheduleCard.tsx`
- **Page Path:** `src/pages/Settings.tsx`
- **Architectural Role:** Full-featured React component card rendering target timezone configurations, searchable IANA lists, live time previews, and scheduling enforcement toggles.

#### 3.6.1 React Component Contract (`CommitScheduleCard.tsx`)

```tsx
import React, { useState, useEffect } from 'react';
import { Clock, ShieldCheck, AlertCircle, Calendar } from 'lucide-react';

export interface CommitScheduleCardProps {
  timezone: string;
  isScheduleActive: boolean;
  cutoffHour: number;
  onUpdate: (data: { timezone: string; isScheduleActive: boolean; cutoffHour: number }) => void;
}

const COMMON_TIMEZONES = [
  'Local',
  'UTC',
  'America/New_York',
  'America/Chicago',
  'America/Denver',
  'America/Los_Angeles',
  'Europe/London',
  'Europe/Paris',
  'Europe/Berlin',
  'Asia/Tokyo',
  'Asia/Singapore',
  'Asia/Shanghai',
  'Australia/Sydney',
];

export const CommitScheduleCard: React.FC<CommitScheduleCardProps> = ({
  timezone,
  isScheduleActive,
  cutoffHour,
  onUpdate,
}) => {
  const [selectedTz, setSelectedTz] = useState(timezone || 'Local');
  const [isActive, setIsActive] = useState(isScheduleActive);
  const [cutoff, setCutoff] = useState(cutoffHour || 8);
  const [currentLocalTime, setCurrentLocalTime] = useState('');
  const [isWindowOpen, setIsWindowOpen] = useState(true);

  useEffect(() => {
    const updateClock = () => {
      try {
        const now = new Date();
        const options: Intl.DateTimeFormatOptions = {
          timeZone: selectedTz === 'Local' ? undefined : selectedTz,
          dateStyle: 'medium',
          timeStyle: 'medium',
        };
        const formatter = new Intl.DateTimeFormat('en-US', options);
        setCurrentLocalTime(formatter.format(now));

        // Evaluate window open state
        const day = now.getDay(); // 0 is Sunday, 6 is Saturday
        const isWeekend = day === 0 || day === 6;
        const hour = now.getHours();
        setIsWindowOpen(isWeekend || hour < cutoff);
      } catch {
        setCurrentLocalTime('Invalid Timezone');
        setIsWindowOpen(true);
      }
    };

    updateClock();
    const interval = setInterval(updateClock, 1000);
    return () => clearInterval(interval);
  }, [selectedTz, cutoff]);

  const handleTzChange = (newTz: string) => {
    setSelectedTz(newTz);
    onUpdate({ timezone: newTz, isScheduleActive: isActive, cutoffHour: cutoff });
  };

  const handleActiveToggle = (active: boolean) => {
    setIsActive(active);
    onUpdate({ timezone: selectedTz, isScheduleActive: active, cutoffHour: cutoff });
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 shadow-lg">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-3">
          <Clock className="w-6 h-6 text-indigo-400" />
          <h3 className="text-lg font-semibold text-slate-100">Commit Timezone & Schedule</h3>
        </div>
        <span
          className={`px-3 py-1 text-xs font-medium rounded-full ${
            isActive ? 'bg-indigo-500/10 text-indigo-400 border border-indigo-500/20' : 'bg-slate-800 text-slate-400'
          }`}
        >
          {isActive ? 'Enforced' : 'Disabled'}
        </span>
      </div>

      <div className="space-y-4">
        {/* Live Status Banner */}
        <div className="bg-slate-950/60 border border-slate-800/80 rounded-lg p-4 flex items-center justify-between">
          <div>
            <div className="text-xs text-slate-400">Target Time ({selectedTz})</div>
            <div className="text-base font-mono font-medium text-slate-200">{currentLocalTime}</div>
          </div>
          <div className="flex items-center gap-2">
            {isWindowOpen ? (
              <span className="flex items-center gap-1.5 text-xs text-emerald-400 bg-emerald-500/10 border border-emerald-500/20 px-2.5 py-1 rounded-md">
                <ShieldCheck className="w-3.5 h-3.5" /> Window Open
              </span>
            ) : (
              <span className="flex items-center gap-1.5 text-xs text-amber-400 bg-amber-500/10 border border-amber-500/20 px-2.5 py-1 rounded-md">
                <AlertCircle className="w-3.5 h-3.5" /> Closed (&ge; {cutoff}:00 AM)
              </span>
            )}
          </div>
        </div>

        {/* Timezone Selector */}
        <div>
          <label className="block text-xs font-medium text-slate-400 mb-1">Target Timezone</label>
          <select
            value={selectedTz}
            onChange={(e) => handleTzChange(e.target.value)}
            className="w-full bg-slate-950 border border-slate-800 text-slate-200 text-sm rounded-lg px-3 py-2 focus:ring-1 focus:ring-indigo-500 focus:outline-none"
          >
            {COMMON_TIMEZONES.map((tz) => (
              <option key={tz} value={tz}>
                {tz}
              </option>
            ))}
          </select>
        </div>

        {/* Schedule Active Toggle */}
        <div className="flex items-center justify-between pt-2">
          <div>
            <div className="text-sm font-medium text-slate-200">Enforce Weekday Cutoff</div>
            <div className="text-xs text-slate-400">Block commits on Monday-Friday after cutoff hour</div>
          </div>
          <input
            type="checkbox"
            checked={isActive}
            onChange={(e) => handleActiveToggle(e.target.checked)}
            className="w-5 h-5 accent-indigo-500 rounded border-slate-800 cursor-pointer"
          />
        </div>
      </div>
    </div>
  );
};
```

---

### 3.7 Component 7: Terminal Rich Help & Documentation

- **Help Menu Path:** `cli/cmdcommit/commit_help_menu.go`
- **Markdown Doc Path:** `cli/helpdoc/commit.md`
- **Architectural Role:** Terminal documentation rendering updated flag descriptions and syntax rules.

#### 3.7.1 `cli/cmdcommit/commit_help_menu.go` Additions

```go
func buildCommitFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-tz, --timezone <tz>", Description: "Target timezone for commit window policy (e.g. UTC, America/New_York)"},
		{Command: "-sc, --schedule <time>", Description: "Schedule commit timestamp explicitly"},
		{Command: "--auto-window", Description: "Auto-schedule for next morning window if current window is closed"},
		{Command: "--interleave", Description: "Merge commits by author date order for commit-both"},
		{Command: "--dry-run", Description: "Preview transfer or commit actions without modifying disk"},
		{Command: "--no-push", Description: "Create local commits without pushing to remote"},
		{Command: "-h, --help", Description: "Show this commit help menu"},
	}
}
```

#### 3.7.2 `cli/helpdoc/commit.md` Content Expansion
Document the following sections:
- **Timezone & Window Policy:** Weekday restrictions (< 08:00 AM) and weekend 24h passes.
- **Flags:**
  - `-tz, --timezone <name>`: Target timezone override.
  - `-sc, --schedule <spec>`: Scheduled target timestamp.
  - `--auto-window`: Auto-deferral to next valid slot.
  - `-y, --non-interactive`: Suppress TUI dialogs.
- **Examples:**
  ```bash
  # Commit in specific timezone:
  gitmap commit "feat: add pipeline telemetry" --tz="America/New_York"

  # Auto-defer if outside morning window:
  gitmap commit "chore: update dependencies" --tz="UTC" --auto-window
  ```

---

## 4. Coding Guidelines & Hygiene Enforcement

| Guideline Rule | Compliance Implementation |
|:---|:---|
| **Leaf Error Sovereignty** | Leaf helpers in `cli/timeutil` return `error` or typed structs directly; zero calls to `cliexit.HandleError` or `os.Exit`. |
| **Boolean Naming** | Only `is*` or `has*` prefixes permitted: `isAllowed`, `isWeekend`, `isWeekday`, `isMorningAllowed`, `isAutoWindow`, `hasCustomTz`. |
| **Function Length** | Strictly <= 15 lines per function (<= 8 lines preferred). Split into discrete sub-helpers. |
| **Conditional Depth** | Zero nested `if` statements. Guard clauses and early returns only. |
| **Whitespace Styling** | Exactly one blank line after closing brace `}` and before `return`. |
| **Strict Relative Paths** | All imports and path references use relative repository notations. |

---

## 5. Verification & Testing Matrix

| Component | Test File | Test Cases |
|:---|:---|:---|
| `cli/timeutil/timezone.go` | `cli/timeutil/timezone_test.go` | Test IANA resolution (`America/New_York`), UTC/Local fallbacks, offsets (`+08:00`, `-0400`, `+0530`), and invalid strings. |
| `cli/timeutil/schedule.go` | `cli/timeutil/schedule_test.go` | Test Saturday 23:59 (allowed), Sunday 12:00 (allowed), Monday 07:59 (allowed), Monday 08:00 (closed), and next window 07:30 AM calculation. |
| `cli/cmdcommit/commit_flags.go` | `cli/cmdcommit/commit_flags_test.go` | Test `--tz`, `-tz`, `--schedule`, `-sc`, `--auto-window`, and combined shorthand flags. |
| `cli/cmdcommit/commit_prompt.go` | `cli/cmdcommit/commit_prompt_test.go` | Test mock reader input for empty commit messages and closed window choices [1-4]. |
| Terminal Rich Help | `cli/cmdcommit/commit_cmd_test.go` | Verify help menu string output contains `--timezone` and `--auto-window`. |
