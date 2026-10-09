# Subtask 04 — Document `gitmap fix` in the Skills (program 246)

## Objective

Add a "Workspace Heal & Fix" section to the gitmap skill so LLMs learn the
`gitmap fix` engine. The skill has THREE homes that must stay byte-identical:

1. `cli/cmd/llm/llm_skill.go` — the `SkillTemplate` const (generator source)
2. `.agents/skills/gitmap/SKILL.md` — the checked-in Antigravity skill
3. `.cursor/skills/gitmap/skill.md` — the Cursor mirror (verified present 2026-10-09)

## Read first

- `02-spec/21-app/246-parallel-fix-command/03-llm-train-integration.md` (§5: the §10 reference)
- `cli/cmdfix/fix_cmd.go` — engine entry, actions, routing keywords (ground truth)

## Owned files (ONLY these)

- `cli/cmd/llm/llm_skill.go` (edit `SkillTemplate` only — insert §10)
- `.agents/skills/gitmap/SKILL.md` (insert identical §10)
- `.cursor/skills/gitmap/skill.md` (insert identical §10)

## Exact section to insert (after §9, before "Command Replacement Matrix")

```markdown
### 10. Workspace Heal & Fix (gitmap fix)
- `gitmap fix` — Show pending remediation summary for dirty/diverged tracked repos (report-only).
- `gitmap fix ls` — Tabular list of repos needing remediation with per-repo recipes.
- `gitmap fix <repo> <action>` — Apply a recipe to one repo. Actions: `stash` (s/1), `wip` (w/2), `discard` (d/3).
- `gitmap fix all <action>` — Apply one action across all pending repos.
- `gitmap fix --prompt` — Interactive per-repo remediation picker.
- `gitmap stash|wip|discard` — Root aliases routing into the same engine with the action pre-selected.
- `gitmap fix agy` — Route Antigravity pipeline errors through the agy fixer.
- `gitmap fix ignore|fia` — Batch .gitignore remediation (local or `--ssh` fleet).
- `gitmap llm train` phase 6 ("Heal & Fix") runs this engine in-process, report-only
  by default; `gitmap llm train --heal-apply stash|wip|discard` applies non-interactively.
```

Adapt the backtick-escaping to the Go const style used in `SkillTemplate`
(`` ` + "`...`" + ` `` concatenation) — the RENDERED markdown in all three files
must be byte-identical.

## Ground-truth notes for the text above (do not invent beyond these)

- Engine entry: `cmdfix.RunFix(args, aliasOverride)`; root aliases `fix`, `stash`,
  `wip`, `discard` (`cli/cmd/rootcore.go:79-82`).
- Actions resolve via `parseRecipeIndex`: `1|stash|s`, `2|wip|w`, `3|discard|clean|d`.
- Routing keywords inside `RunFix`: `agy|aef|agy-errors-fix` → agy fixer;
  `ignore|ignores|gitignore|fia|fias|fix-ignore-all` → ignore dispatch;
  `ls` → `runFixLs`; `--prompt|-p|interactive` → interactive remediation.
- Empty remediation state triggers live discovery of tracked repos (`handleLiveDiscoveredIssues`).

## Done criteria

- [ ] §10 inserted in `SkillTemplate`, `.agents/skills/gitmap/SKILL.md`, and
      `.cursor/skills/gitmap/skill.md` — rendered markdown byte-identical across all three.
- [ ] No other skill sections modified; no file deletions; no alias changes.
- [ ] Report the byte-identity check method used (e.g. `cmp` of rendered outputs).

## Hard rules

- Search ONLY via `gitmap aum search` / `gitmap find` / `gitmap cat`. TOTAL BAN on
  grep, rg, git grep, Select-String.
- NEVER run git commands. NEVER run `go build` / `go test` (lead verifies centrally).
- Relative paths only. Small files only (this section, not new files).
- The lead regenerates both skill files via `gitmap llm train` after merge and
  confirms byte-identity with the generator output — flag any drift in your report.
