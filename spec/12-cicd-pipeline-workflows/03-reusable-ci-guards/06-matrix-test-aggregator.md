# Pattern 06 — Matrix Test Aggregator

**Version:** 1.0.0
**Updated:** 2026-04-21

---

## Problem Statement

A test matrix runs N parallel jobs (sharded by package, by Go version,
by OS, etc.). Each job uploads its raw test output as an artifact. A
reviewer who wants to know "which tests failed?" must:

1. Open each of N artifact archives
2. Search each for `--- FAIL` (or framework-specific marker)
3. Manually correlate failure reasons across files

The aggregator pattern collapses this into **one summary step** that
downloads all matrix artifacts, parses each, and emits a single
copy-paste-ready failure report.

---

## Algorithm

```
INPUT:  results_dir/test-results-<suite>/test-output.txt  (one per matrix shard)
OUTPUT: exit 0 if all suites pass, exit 1 with copy-paste failure report

1. for each subdir under results_dir matching `test-results-*`:
     suite = basename without prefix
     pass_count = grep -c '<PASS_MARKER>' file
     fail_count = grep -c '<FAIL_MARKER>' file

     if fail_count > 0:
        for each failing test:
            extract test name
            extract failure reason (assertion message + file:line)
        record suite section
        overall = 1

     emit one-line per-suite summary (✅ / ❌)

2. if overall == 1:
     print fenced FAILURE REPORT block
     exit 1
   else:
     exit 0
```

The **failure-reason extraction** is framework-specific. The reference
Go implementation captures lines between `=== RUN <test>` and
`--- FAIL: <test>` and filters for assertion patterns
(`expected|got|Error|FAIL|panic|undefined|mismatch`) plus any
`*.go:NN:` source pointers, capping at 10 lines per failure.

---

## Contract

| Aspect | Value |
|--------|-------|
| Invocation | `bash test-summary.sh <results_dir>` |
| Exit `0` | All suites passed |
| Exit `1` | At least one suite has failing tests (full report on stdout) |
| Exit `2` | Results dir missing or contains no shard subdirs |
| stdout | Per-suite summary + copy-paste failure block + `::error::` annotation |

---

## Configuration Surface

| Variable | Purpose | Example |
|----------|---------|---------|
| `RESULTS_DIR` | Where matrix artifacts were downloaded | `./test-artifacts` |
| `SHARD_PATTERN` | Glob for shard subdirs | `test-results-*` |
| `OUTPUT_FILENAME` | Filename inside each shard dir | `test-output.txt` |
| `PASS_MARKER` | Per-test PASS line pattern | `^--- PASS:` (Go) |
| `FAIL_MARKER` | Per-test FAIL line pattern | `^--- FAIL:` (Go) |
| `EXTRACT_REASON_FN` | Function returning failure context per test | See per-framework section |

---

## Adaptations

### Go (`go test -v`)

- `PASS_MARKER`: `^--- PASS:`
- `FAIL_MARKER`: `^--- FAIL:`
- Test-name extraction: `sed 's/^--- FAIL: //' | sed 's/ (.*$//'`
- Reason extraction: lines between `=== RUN <test>` and
  `--- FAIL: <test>` matching `\.go:[0-9]+:` or assertion keywords

### Node / Vitest (with `--reporter=verbose`)

- `PASS_MARKER`: `^ ✓ ` (vitest tick)
- `FAIL_MARKER`: `^ ✗ ` or `^ × `
- Test-name extraction: strip the marker prefix
- Reason extraction: lines following `FAIL  <suite>` until the next
  test or blank line; capture `at <file>:<line>` pointers

### Python (pytest with `-v`)

- `PASS_MARKER`: ` PASSED$`
- `FAIL_MARKER`: ` FAILED$`
- Reason extraction: between `___ <test> ___` headers and the next
  `___` divider; capture `assert` lines and tracebacks

### Rust (`cargo test`)

- `PASS_MARKER`: ` ... ok$`
- `FAIL_MARKER`: ` ... FAILED$`
- Reason extraction: from the `failures:` block at end of output;
  capture `panicked at` and `assertion failed` lines

---

## Output Format

```
=========================================
  ALL TEST RESULTS
=========================================

✅ unit-shard-01: 142 passed
✅ unit-shard-02: 138 passed
❌ integration-shard-01: 2 failed, 41 passed
✅ integration-shard-02: 38 passed

=========================================
  FAILURE REPORT (copy-paste ready)
=========================================

-----------------------------------------
  Suite: integration-shard-01 (2 failed)
-----------------------------------------
    --- FAIL: TestReleaseFlow_Bumps
        release_test.go:88: expected version "v1.2.0", got "v1.1.0"
        release_test.go:91: Error: bump did not apply

    --- FAIL: TestReleaseFlow_Tags
        release_test.go:131: expected tag "v1.2.0" to exist, got nothing

=========================================
::error::Some test suites failed — see failure report above.
```

---

## Workflow Wiring

```yaml
jobs:
  test:
    strategy:
      matrix:
        shard: [01, 02, 03, 04]
    steps:
      - run: <test-command> | tee test-output.txt
      - uses: actions/upload-artifact@v4
        with:
          name: test-results-${{ matrix.shard }}
          path: test-output.txt

  summarize:
    needs: test
    if: always()
    steps:
      - uses: actions/download-artifact@v4
        with:
          path: ./test-artifacts
          pattern: test-results-*
      - run: bash .github/scripts/test-summary.sh ./test-artifacts
```

`if: always()` ensures the summary runs even when shards fail.

---

## Failure Modes

| Pitfall | Fix |
|---------|-----|
| `set -e` aborts on `grep -c` returning non-zero (no matches) | Use `set -uo pipefail` (omit `-e`) and append `|| true` to count operations |
| Reason extractor captures hundreds of stack-trace lines | Cap with `head -10` per test |
| Non-ASCII test names break `sed` patterns | Use `awk` with explicit field separators OR ensure tests use ASCII names |
| `download-artifact` with `pattern:` requires `actions/download-artifact@v4` | Pin to `@v4`; v3 had different API |

---

## Reference Implementation

See `.github/scripts/test-summary.sh` in `gitmap-v6`.

---

## Cross-References

- [00-overview.md](./00-overview.md)
- [Shared Conventions](../01-shared-conventions.md) — Pinned `actions/download-artifact@v4`

---

*Matrix test aggregator — v1.0.0 — 2026-04-21*