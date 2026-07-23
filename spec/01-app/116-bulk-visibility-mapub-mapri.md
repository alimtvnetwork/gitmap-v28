# Bulk wildcard visibility flips — `make-all-public` / `make-all-private` / `MAPUB` / `MAPRI`

**Status:** pending
**Plan:** `.lovable/plans/pending/01-bulk-visibility-mapub-mapri.md`
**Related:** `spec/01-app/113-clone-parent-escape-and-bulk-visibility.md` (single-repo + `vN..vN-count+1` window flips — these new commands are owner-wide and pattern-driven, not version-windowed).

---

## 1. Symptoms / Why

Today the only way to flip many repos under an owner/org is to loop
`gitmap make-public <url>` by hand. The existing single-target commands
(`gitmap/cmd/visibility.go`, `gitmap/cmd/visibilitybulk.go`) cannot:

- enumerate every repo under a given owner;
- filter that list with wildcard patterns (`exact`, `prefix*`,
  `*contains*`, `prefix*suffix`, comma-separated);
- show an interactive numbered table with per-index exclusion;
- persist a per-run audit trail (run row + per-repo result rows).

This spec adds four new top-level commands and the supporting pattern
engine, interactive prompt, and SQLite audit schema.

---

## 2. Command grammar

```
gitmap make-all-public  <target> <patterns> [-Y|--yes]
gitmap make-all-private <target> <patterns> [-Y|--yes]
gitmap MAPUB            <target> <patterns> [-Y|--yes]    # alias of make-all-public
gitmap MAPRI            <target> <patterns> [-Y|--yes]    # alias of make-all-private
```

`<target>` — any of:
- full provider URL: `https://github.com/acme`, `git@gitlab.com:acme.git`
- bare `host/owner`: `github.com/acme`
- folder path: `./acme` or `.` (origin owner extracted from `.git/config`)

`<patterns>` — comma-separated list of wildcard patterns; see §3.

`-Y` / `--yes` — skip BOTH the y/n confirm and the exclusion prompt;
operate on every matched repo.

---

## 3. Pattern rules

Exactly one wildcard char: `*`. No regex, no `?`, no character classes.

| Pattern shape | Meaning                                    |
|---------------|--------------------------------------------|
| `macro`       | exact match                                |
| `macro*`      | prefix                                     |
| `*macro`      | suffix                                     |
| `*macro*`     | contains                                   |
| `lotus*v1`    | prefix + suffix (interior `*` = wildcard)  |
| `a*b*c`       | ordered: `a` … `b` … `c` (greedy)          |
| `*`           | **error** — refuse bare `*` (footgun)      |

Comma list: split on `,`, trim, dedupe, reject empty tokens. Match per
pattern; union across patterns; dedupe by repo name keeping first matcher
(used in the “matched by” column).

---

## 4. Interactive flow

```
$ gitmap make-all-public github.com/acme "macro*,*-vault"
→ Listing repos under acme (gh repo list --limit 1000)…
→ Matched 7 of 142 repos:

   1  macro-frontend         (macro*)
   2  macro-api              (macro*)
   3  macro-docs             (macro*)
   4  secret-vault           (*-vault)
   5  team-vault             (*-vault)
   6  macro-vault            (macro*)        ← first matcher wins
   7  macro-internal         (macro*)

Make these 7 repos PUBLIC? [y/N/exclude] _
```

User input:
- `y` / `Y` → proceed with all 7.
- `n` / `N` / empty → abort, mark every row `Skipped`, exit 0.
- numeric list (`1,3-5`, ranges allowed, `none`, `all`) → set
  `IsExcluded=1` on listed rows, re-prompt with the remaining set.

`-Y` short-circuits both prompts.

---

## 5. SQLite DDL (migration `0NN_gitmap_run_repo_result.sql`)

```sql
CREATE TABLE IF NOT EXISTS GitMapRun (
    GitMapRunId      INTEGER PRIMARY KEY AUTOINCREMENT,
    CommandKind      INTEGER NOT NULL,         -- enum: 1=MakeAllPublic 2=MakeAllPrivate
    TargetRaw        TEXT    NOT NULL,         -- exact arg as typed
    ResolvedOwner    TEXT    NOT NULL,
    ResolvedProvider TEXT    NOT NULL,         -- 'github'|'gitlab'
    PatternsRaw      TEXT    NOT NULL,         -- comma-joined input
    IsAutoConfirm    INTEGER NOT NULL DEFAULT 0,
    StartedAtUnix    INTEGER NOT NULL,
    FinishedAtUnix   INTEGER
);

CREATE TABLE IF NOT EXISTS GitMapRepoResult (
    GitMapRepoResultId INTEGER PRIMARY KEY AUTOINCREMENT,
    GitMapRunId        INTEGER NOT NULL REFERENCES GitMapRun(GitMapRunId) ON DELETE CASCADE,
    RepoName           TEXT    NOT NULL,
    MatchedPattern     TEXT    NOT NULL,
    IsExcluded         INTEGER NOT NULL DEFAULT 0,
    ResultStatus       INTEGER NOT NULL,       -- enum: 1=Pending 2=Succeeded 3=Failed 4=Skipped
    ErrorMessage       TEXT,
    UpdatedAtUnix      INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS IX_GitMapRepoResult_RunId ON GitMapRepoResult(GitMapRunId);
```

Schema obeys Core memory: PascalCase, `INTEGER PRIMARY KEY AUTOINCREMENT`,
`IF NOT EXISTS` for idempotence, single-conn pool (`SetMaxOpenConns(1)`).

### ERD

```text
GitMapRun (1) ────< (n) GitMapRepoResult
   GitMapRunId  ◄──── GitMapRunId  (FK, ON DELETE CASCADE)
```

---

## 6. Exit codes

| Code | Meaning                                          |
|------|--------------------------------------------------|
| 0    | All matched repos succeeded (or user aborted)    |
| 1    | Usage error (bad args, bad pattern, bad target)  |
| 2    | At least one per-repo failure                    |

Matches the convention in `gitmap/cmd/visibility.go`.

---

## 7. Error management (zero-swallow)

- Per-repo failure → `fmt.Fprintf(os.Stderr, "✗ %s: %v\n", repo, err)`
  and `UPDATE GitMapRepoResult … ResultStatus=Failed, ErrorMessage=…`.
- Ctrl-C → `errors.Is(err, context.Canceled)` → mark remaining
  `Pending` rows as `Skipped`, flush, exit 130.
- Tx-open failure on the audit DB does NOT block the visibility flips;
  it logs to stderr with Code Red format
  `"Error: failed to open audit db at %s: %v (operation: begin tx, reason: %s)"`
  and continues without persistence.
- All user-facing strings live in
  `gitmap/constants/constants_visibilitybulk.go` (no magic strings).

---

## 8. Acceptance checklist

- [ ] `rg "make-all-public|make-all-private|MAPUB|MAPRI" gitmap/` → 4 CLI
      ID constants in `constants_cli.go`, 4 help MDs, dispatcher entries,
      handler, tests.
- [ ] `go test ./gitmap/cmd/... ./gitmap/visibility/... ./gitmap/db/... -race -count=1` green.
- [ ] `golangci-lint run ./...` zero findings (v1.64.8).
- [ ] `gitmap regoldens` clean (help fixtures regenerated).
- [ ] Migration applies on a fresh `:memory:` DB and FK cascade works.
- [ ] `-Y` smoke run skips both prompts (stdin stub panics if read).
- [ ] CHANGELOG entry under next minor version + docs-site mirror.
- [ ] `gitmap fix-repo --strict` clean across touched packages.
