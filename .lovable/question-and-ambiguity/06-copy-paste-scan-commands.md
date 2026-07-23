# 06 — Copy-paste scan commands per scenario

## Original task

> Insert copy-paste scan command examples for each scenario including `--config`, `--mode`, and `--output csv` so I can reproduce the outputs quickly.

## Ambiguity

1. **Replace or supplement?** The existing Examples A/B/C and the edge-case section already each have a one-line `gitmap scan <root> --output csv` block. Should I (a) edit those in place to include the full triple, or (b) add a new dedicated "copy-paste recipes" section that mirrors each scenario with the full triple?
2. **Which `--config` to reference?** Options: (i) inline a literal `gitmap.config.json` JSON snippet in every block, (ii) point to the `data/config.json` example already present further down the README, or (iii) write `--config <your-config>` placeholder.
3. **Should I include `--output-path`?** It's not explicitly in the user's request but it's the missing piece for the "land on a known artifact path" use case the user implies with "reproduce the outputs quickly".

## Options considered

### Q1 — replace vs supplement

| Option | Pros | Cons |
|---|---|---|
| **A. Add a dedicated "Copy-paste scan commands per scenario" section after Example C** (recommended) | Keeps existing minimal examples focused on the marker / depth logic; the new section is a single scannable place for "what do I actually run?"; doesn't bloat each example with flags that aren't relevant to the rule it's demonstrating | Slight duplication of scan invocations |
| B. Edit each Example A/B/C and the edge-case block in place | No duplication | Each example becomes 5+ lines of flags before the layout shows, hurting the "look at this rule" focus they were designed for |
| C. Add the full triple only to the edge-case section (the most recent) | Minimal change | Doesn't satisfy "for each scenario" — the user said all of them |

### Q2 — `--config` reference style

| Option | Pros | Cons |
|---|---|---|
| **A. Reference an external `./gitmap.config.json` and link to the existing exclude-list section** (recommended) | Avoids inline-JSON repetition; users can copy one config and reuse across all four blocks; cross-reference establishes the canonical example | Requires the user to scroll down once to see the config |
| B. Inline JSON in every block | Self-contained | Quadruples the visual weight of the section; users will copy stale snippets |
| C. Use `<your-config>` placeholder | Smallest footprint | Defeats "copy-paste ready" — users have to substitute |

### Q3 — include `--output-path`?

| Option | Pros | Cons |
|---|---|---|
| **A. Show `--output-path` only on the "second variant" of each scenario** (recommended) | Demonstrates the flag exists without making every block long; first variant uses the default path so the reader sees both behaviors | Mild inconsistency |
| B. Always include `--output-path` | Maximum explicitness | Doubles the line length of every command |
| C. Never include it | Shortest | Users hit the default-path artifact at `./.gitmap/output/gitmap.csv` and may not realize it's tunable |

## Recommendation

- Q1 → **A**: add a dedicated H4 after Example C.
- Q2 → **A**: reference one external config.
- Q3 → **A**: show `--output-path` on the secondary variant of each scenario.

## Decision taken

Inserted a new H4 section **"Copy-paste scan commands per scenario"** between Example C and the "Reading at-cap CSV rows" section. The section provides one copy-paste block per existing scenario:

1. **Reproduce Example A** — HTTPS variant + SSH variant (with `--output-path ./reports/markers`); also clarifies that both `httpsUrl` / `sshUrl` columns are always populated regardless of `--mode` (only `cloneInstruction` switches), verified against `gitmap/mapper/mapper.go::buildOneRecord` lines 89–94.
2. **Reproduce Example B** — primary scan + two `--output-path`-redirected re-scans of `~/work/main-repo/modules` and `~/work/main-repo/vendor` to catalog the rule-2-hidden subtrees.
3. **Reproduce Example C** — primary scan + at-cap-directory re-scan, with a forward link to the "Reading at-cap CSV rows" section for the additive-upsert rationale.
4. **Reproduce the edge-case layout** — single command + a `find … -name .git` / `awk` / `comm -23` audit pipe so users can verify "what was silently dropped" in CI.

All four blocks reference one external `./gitmap.config.json` (no inline JSON repetition). The full triple `--config + --mode + --output csv` appears in every primary command; `--output-path` appears on every secondary command.

## Self-correction note

Initially considered claiming `--mode ssh` produces a CSV where `httpsUrl` is empty. Verified against `gitmap/mapper/mapper.go::buildOneRecord` (lines 89–94) before publishing — both URL columns are always emitted via `toHTTPS()` / `toSSH()`; `--mode` only routes which one feeds `selectCloneURL` and thus the `cloneInstruction` column. The documentation now reflects this accurately.
