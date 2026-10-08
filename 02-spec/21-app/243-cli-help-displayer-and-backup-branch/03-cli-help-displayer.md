# Spec 243.3 — CLI Help Displayer System

## Goal
One DRY, themeable, object-oriented help-rendering system replacing hand-maintained
help builders. The terminal output users see today stays the same; the internals
become a proper object model with theming power.

## Core types (new package `cli/helpdisplay/`)

```go
// Displayer is the interface every help display implements.
type Displayer interface {
    Render(ctx RenderContext) string
    Print(ctx RenderContext)
}

// HelpDisplay is the concrete struct, bound to Displayer via constructor binding
// (coding-guideline Go binding principle — read 2 existing examples, cite them).
type HelpDisplay struct {
    header      string
    groups      []CommandHelpGroup
    suggestions []Suggestion
    theme       *Theme
}

type CommandHelper struct {
    command     string
    description string
    example     string
    url         string
    subHelpers  []CommandHelper
}

type CommandHelpGroup struct {
    header   string
    commands []CommandHelper
    hints    []string
}

type Suggestion struct {
    text      string
    variables map[string]string // interpolated from RenderContext at render time
}

type RenderContext struct {
    values map[string]string // cached: e.g. detected IPs, repo names, versions
}

type Theme struct {
    headerColor      Color
    commandColor     Color // left side: the command text
    descriptionColor Color // right side: description text
    hintColor        Color
    parent           *Theme // inheritance chain
}
```

## Rules
- **Binding**: `NewHelpDisplay(...) Displayer` constructor binds struct → interface.
  Follow the binding pattern used by sibling packages per the coding guidelines.
- **Colors**: reuse `cli/theme/` and `cli/glyphs/` — no new color engine. Theme
  properties cover header, command (left), description (right), hints.
- **Inheritance + override**: item-level color overrides group-level, which overrides
  display-level, which overrides theme-level, which falls back to the global default.
  An explicit override always wins; unset inherits. This is what makes per-help-text
  animation/tweaks possible on top of a shared theme.
- **Suggestions with variables**: `Suggestion.variables` interpolate from
  `RenderContext.values` at render time. Context values are cached (e.g. detected
  machine IPs) so help shows real values (`ssh user@192.168.1.22`) instead of
  placeholders (`ssh user@<ip>`).
- **Auto-display**: `--help` on any command builds its `HelpDisplay` and prints it.
  A `CommandHelper` with `subHelpers` renders its sub-display automatically on
  `--help <sub>`.
- **DRY**: the structs are the source of truth. `helptext/*.md` topics are GENERATED
  from them (generator location: implementer's choice — `03-ai-scripts/` or a gitmap
  subcommand — documented in the plan before building).

## Migration (incremental — no big bang)
1. Build the package + theme wiring; prove `go build` clean.
2. Port 2–3 pilot commands (`space`, `scan`, `pull-all`): output must match today's
   rendering (modulo intentional improvements) — diff the terminal output as evidence.
3. New/rewritten commands use the displayer from day one (see spec 04 §5).
4. Migrate remaining builders opportunistically, duplicated/outdated first.

## Implementer must verify in code before building
- The exact binding pattern from 2 sibling packages (cite file:line in the plan).
- How `cli/theme/` exposes colors today — reuse, don't wrap unnecessarily.
- Where `--help` is dispatched per command today, to hook auto-display.
