# 264 — `scan export` / `scan merge`: acceptance criteria and evidence

## AC1 — `scan export` writes `<machine-slug>/repos.json` from cache

**Given** a scanned repo set in the cache,
**when** `gitmap scan export --out /tmp/x` runs,
**then** `/tmp/x/<hostname-slug>/repos.json` exists and contains the full
cached repo list as a JSON array of scan records (≥1 of `httpsUrl`/`sshUrl`
per row).

**Evidence:** command output + `jq length` of the JSON.

## AC2 — `--machine` overrides the slug

**Given** `--machine "My Laptop!"`,
**when** export runs,
**then** the folder is the sanitized slug (`my-laptop`) and the JSON sits
inside it.

**Evidence:** `ls` of the output dir.

## AC3 — machine slug is stable and documented

**Given** no `--machine`,
**when** export runs twice on the same host,
**then** both runs write the same folder (hostname-derived), and `--help`
states the default.

**Evidence:** two runs, same path; help text captured.

## AC4 — `scan merge` dedupes across folders

**Given** two export folders sharing ≥1 repo (same `httpsUrl`),
**when** `gitmap scan merge /tmp/a /tmp/b --out /tmp/m.json` runs,
**then** `/tmp/m.json` contains each unique repo exactly once (count ==
union count), order preserved, first occurrence wins.

**Evidence:** `jq` counts before/after.

## AC5 — clone side consumes the merged file unchanged

**Given** the merged JSON from AC4,
**when** `gitmap clone-from /tmp/m.json` (dry-run) runs,
**then** it plans clones for every repo in the file with no new code and no
schema conversion.

**Evidence:** dry-run plan output lists the repos.

## AC6 — normal `scan <dir>` is unaffected

**Given** the interception of `export`/`merge`,
**when** `gitmap scan /tmp/somedir` runs,
**then** scanning behaves exactly as before (no flag-parse regression).

**Evidence:** scan output on a temp dir.

## AC7 — build, vet, and help

`go build ./...` exit 0, `go vet ./...` clean, `gitmap scan export --help`
and `gitmap scan merge --help` print usage.
