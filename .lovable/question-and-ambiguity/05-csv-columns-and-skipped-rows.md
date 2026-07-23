# 05 — CSV columns + skipped/edge-case rows in README

## Original task

> Add a README section that shows the exact CSV columns and example rows for skipped non-repos and marker-like edge cases (including the depth column).

## Ambiguity

Two minor points required inference:

1. **What does the CSV look like for a "skipped" row?** The user phrased this as "example rows for skipped non-repos", which could mean (a) literal rows in the CSV that *describe* skips (with a `notes` column explaining why), or (b) showing what the scan emits when the input *contains* skip-able layouts — i.e. the CSV omits them entirely and only the survivors appear.
2. **Which marker-like edge cases to enumerate.** The §1 bullet list (added in task 04) lists 7 negative cases; should the new section repeat all of them or pick a representative subset?

## Options considered

### Q1 — meaning of "rows for skipped non-repos"

| Option | Pros | Cons |
|---|---|---|
| **A. Show the CSV as the scanner currently emits it (skipped layouts simply absent), and use a separate "Why each was rejected" table to explain the silence** (recommended) | Truthful to today's behavior; matches `csv_header_contract_test.go`; doesn't promise a feature that doesn't exist | Less visually punchy than seeing the skip reasons inline |
| B. Mock up rows with `notes="skipped: …"` as if the scanner emitted diagnostic rows | Visually richer | Misleads users — `gitmap scan` does not produce these rows; would create support tickets when users grep for `notes^=skipped` and find nothing |
| C. Add a real `--include-skipped` flag and emit diagnostic rows | Best UX long-term | Out of scope for a docs task; would need scanner changes, contract-test updates, and a version bump for a new column-meaning |

### Q2 — coverage of edge cases

| Option | Pros | Cons |
|---|---|---|
| **A. Cover all 7 negative bullets from §1 plus the depth-cap case, in one combined layout** (recommended) | Single source of truth; one `tree` block + one CSV block + one rejection table; easy to scan | Slightly long |
| B. Pick 3 representative cases (stray text, bare repo, deep) | Shorter | Loses the "predictability" guarantee the user asked for in task 04 |
| C. One mini-example per case | Most thorough | 7+ near-identical code blocks; reader fatigue |

## Recommendation

- **Q1 → Option A** with an explicit note that diagnostic rows are a tracked-but-unimplemented future change. Keeps documentation honest with code.
- **Q2 → Option A**: one combined edge-case layout that covers all §1 negative bullets plus the depth cap, so the section is the predictive reference the user asked for.

## Decision taken

Inserted two new H4 sections into `README.md` between Example C (depth-cap) and the `rescan` section:

1. **"CSV column reference — the 10 columns, in order"** — links to the contract test (`csv_header_contract_test.go`) and the model (`gitmap/model/record.go`), shows the exact header line, and provides a per-column table covering type, source, and meaning. Includes the line-ending and quoting contract (RFC 4180 / `\r\n`).
2. **"Skipped non-repos and marker-like edge cases — what does NOT appear in CSV"** — a single `~/edge/` layout with 11 child directories covering all 7 negative cases from §1 plus the depth cap and 2 positive controls (worktree-link FOUND, nested-under-real/inner FOUND). The example CSV shows only the 3 surviving rows; a "Why each skipped row was rejected" table maps each rejected layout back to the §1 bullet that explains the rejection. A trailing paragraph notes the `notes` column is empty for skipped rows today and that diagnostic-row emission is not yet implemented.

Also corrected the `cloneInstruction` column description and the example CSV rows after grepping `gitmap/mapper/mapper.go` and `gitmap/formatter/formatter_test.go` confirmed the actual format is `git clone -b <branch> <url> <relativePath>` (not `git clone --branch <branch> . <relativePath>` as I'd initially written from memory).

## Self-correction note

First-draft example rows showed `git clone --branch main . real-repo` (positional `.` source). Verified against the codebase that the real format uses the URL from the `httpsUrl`/`sshUrl` columns and `-b` short form, not `--branch` long form. Fixed before finalizing.
