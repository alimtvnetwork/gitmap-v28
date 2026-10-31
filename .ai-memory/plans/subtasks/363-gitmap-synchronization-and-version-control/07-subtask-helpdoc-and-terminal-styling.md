# Subtask 07: Help Documentation & Terminal UI Tokyo Cards Integration

- **Parent Plan:** [363-gitmap-synchronization-and-version-control](../../pending/363-gitmap-synchronization-and-version-control.md)
- **Subtask ID:** `TASK-07`
- **Slug:** `gitmap-synchronization-and-version-control - Task 07`
- **Target Package:** `cli/helpdoc/`, `cli/richhelp/`
- **Status:** `pending`
- **Assigned Worker:** UI & Documentation Subagent

---

## 1. Goal

Author rich markdown documentation for all synchronization and Muse commands, and implement the Tokyo Night themed Lipgloss card rendering engine for terminal summary tables, conflict alerts, and interactive help dialogs.

---

## 2. In Scope

1. **Markdown Help Documentation (`cli/helpdoc/`):**
   - Author authoritative help documents:
     - `cli/helpdoc/cg_sync.md`: Coding guidelines synchronization, options, conflict behavior.
     - `cli/helpdoc/prompt_sync.md`: Prompt architect synchronization, V6/V7 templates.
     - `cli/helpdoc/skill_sync.md`: Antigravity skill synchronization, discovery, and manifests.
     - `cli/helpdoc/muse.md`: Muse onboarding, training vault links, and memory workflows.
   - Register entries in `cli/helpdoc/` topic mappings and `rich_help_dispatcher.go`.

2. **Tokyo Night Lipgloss Card Engine (`cli/richhelp/tokyo_cards.go`):**
   - Theme Tokens:
     - Header Title: Tokyo Purple (`#bd93f9`, Bold)
     - Accents & Subtitles: Cyan (`#8be9fd`)
     - Clean Success Badges: Spring Green (`#50fa7b`)
     - Human Conflict Alerts: Vibrant Amber / Gold (`#ffb86c`, Bold border)
     - Muted Comments / Paths: Tokyo Comment Grey (`#6272a4`)
     - Card Borders: Subtle Lavender (`#44475a`)
   - Card Components:
     - `RenderSyncSummaryCard(report SyncReport)`: Displays total scanned assets, updated count, skipped count, and conflict count in a dual-column bordered card.
     - `RenderConflictAlertCard(conflict ConflictEntry)`: Renders high-visibility warning box displaying relative path, disk modification time, and actionable remediation switches (`--force`, `--backup`).
     - `RenderPromptCard(promptText string, stats PromptStats)`: Formats on-the-fly execution prompts with token count and clipboard status badges.

3. **Terminal Help Dispatcher Parity:**
   - Ensure `gitmap thelp cg sync`, `gitmap thelp prompt sync`, `gitmap thelp skill sync`, and `gitmap thelp muse` render complete interactive cards.

---

## 3. Out of Scope

- Modifying Web UI frontend CSS files (this subtask focuses on CLI Terminal UX and helpdoc markdown).
- Full builds or test suite execution.

---

## 4. Architectural Invariants

- **Positive Booleans:** `hasBorder`, `isCompact`, `hasConflicts`, `isTokyoTheme`.
- **Function Line Limit:** Lipgloss styling helpers <= 15 lines per function.
- **Relative Paths:** Markdown documentation must only reference 100% relative paths.
- **Error Sovereignty:** Documentation rendering must return errors upward without exiting.

---

## 5. Verification Checklist

- [ ] All four markdown help docs exist under `cli/helpdoc/` with complete flags and examples.
- [ ] Help documents registered in `rich_help_dispatcher.go`.
- [ ] Tokyo Night Lipgloss styles render high-contrast visual cards in terminal.
- [ ] Conflict alert cards highlight human-modified files in amber with clear instructions.
- [ ] `gitmap thelp muse` and `gitmap thelp cg sync` display properly formatted cards.
