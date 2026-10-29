# Design System Spec Improvements — Repeatable Prompt

> **Prompt version:** 1.3.0
> **Updated:** 28-Oct-26
> **Last id on file:** `DS-IMP-400`. The next run continues at `DS-IMP-401`.
> **Category:** `01-prompts/30-improvements/`

[/goal](slashCommand;goal) Add the next 100 design-system improvements so a weak AI can still produce a professional slide or a compact tree in one shot.

[/learn](slashCommand;learn) Read [gitmap](file;.agents/skills/gitmap) and the spec folder `02-spec/21-app/35-coding-guideline-and-design-system-spec-improvements/` before writing.

## Must do

1. Run `gitmap pull` in this repository before edits. If that fails, run `git pull` and resolve conflicts.
2. Read, in order:
   - `02-spec/21-app/35-coding-guideline-and-design-system-spec-improvements/readme.md`
   - `02-spec/21-app/35-coding-guideline-and-design-system-spec-improvements/01-architecture-and-token-contract.md`
   - `02-spec/21-app/35-coding-guideline-and-design-system-spec-improvements/02-one-hundred-design-improvements.md`
   - `02-spec/21-app/35-coding-guideline-and-design-system-spec-improvements/03-tree-view-and-detail-panel.md`
   - `02-spec/21-app/35-coding-guideline-and-design-system-spec-improvements/04-blind-ai-execution-checklist.md`
   - `02-spec/07-design-system/40-theme-switch.md`
   - `02-spec/07-design-system/20-ai-training-and-checklist-guide.md`
3. Search only with `gitmap aum search`, `gitmap find`, `gitmap cat`, `gitmap lf`. Do not use ripgrep, grep, or Select-String.
4. Append 100 new rows to `02-one-hundred-design-improvements.md`. Continue the last `DS-IMP` id. Do not renumber and do not delete `DS-IMP-001` through the current last id.
5. Each new row needs a Do and a Do not. Cite a relative path that was read. If a sibling repo is the evidence, name it as `<repo>/<relative-path>` with no drive letter and no `file:///` URI.
6. Keep the decision table. Unnamed slides stay `bright-gold-tech`. Do not mix file 30 amber hex, `slides-app` blue, or website JSON into that theme.
7. Tree work stays generic: about 12px rows, 340px pane clamped 280 to 760, 40px collapsed rail, four selection roles, detail bound to the leaf.
8. If `letsmarknow` still has no UI source, write that it is missing. Do not invent components.
9. Update the index row in `02-spec/21-app/readme.md` only if the folder title changes. Update `.ai-memory/plans/pending/35-coding-guideline-and-design-system-spec-improvements.md` with the new id range.
10. Do not restore `.agents/skills-archived/create-slide-deck/` or `.agents/skills-archived/slides-deck-management/`.
11. Do not bump `version.json` and do not edit a changelog unless the user asks for a release.
12. Finish with one GitMap feature commit in hyphen form: module `Design System`, no colon in the message.

## Must not

- Do not paint a purple-blue full-screen gradient, a second accent, or a `scale(1.05)` hover and call it the house style.
- Do not bind the detail pane to the tree root.
- Do not rotate tree chevrons.
- Do not write absolute paths in any file, plan, or release note.
- Do not recode `slides-app/` on a spec-only run. That stays review-first until the user asks for implementation.

## Done when

- One hundred new ids exist after the previous last id.
- A weak model can still answer the 12 probes in `05-benchmark-previous-versus-new.md` from the spec alone.
- The working tree contains only this improvement and its index updates.
