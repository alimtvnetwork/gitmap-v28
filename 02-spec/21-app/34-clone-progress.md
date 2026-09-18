# Clone Progress — Visual Feedback for Multi-Repo Clone

## Overview

The clone command gains a progress indicator that shows the current
repo being cloned, a counter, and elapsed time during multi-repo
clone operations.

---

## How It Works

1. Before cloning starts, count total repos from the input file.
2. For each repo, print a progress line to stderr.
3. After each clone completes, update the counter.
4. On completion, print a summary with total time.

---

## Progress Output

```
[  1/24]  Cloning api-gateway ...
[  2/24]  Cloning frontend-app ... done (3.2s)
[  3/24]  Cloning shared-lib ... done (1.8s)
[  4/24]  Cloning docs-site ...
```

### Completion Summary

```
Clone complete: 24/24 repos in 2m 14s
  Cloned: 20 | Pulled: 3 | Failed: 1
```

### Error Handling

Failed repos show inline errors and are collected for the summary:

```
[  5/24]  Cloning broken-repo ... FAILED (exit 128)
```

---

## Flags

| Flag | Description |
|------|-------------|
| `--quiet` | Suppress progress output (existing flag, extended) |

When `--quiet` is set, no progress lines are printed. Only the
final summary appears.

---

## Implementation Notes

- Progress goes to stderr so stdout piping is unaffected.
- Counter is right-aligned with padding for clean alignment.
- Elapsed time per repo shown after completion.
- Total elapsed time in the summary uses `time.Since()`.
- No external dependencies; plain `fmt.Fprintf(os.Stderr, ...)`.

---

## File Layout

| File | Purpose |
|------|---------|
| `cloner/progress.go` | Progress tracker struct and display methods |
| `constants/constants_clone.go` | New progress format strings (add to existing) |

Changes to existing files:
- `cloner/cloner.go` — integrate progress tracker into clone loop

---

## Constraints

- Progress must not interfere with `--verbose` debug logging.
- All files under 200 lines, all functions 8–15 lines.

## See Also

- [Cloner](05-cloner.md) — File-based and direct URL clone behavior
- [Clone Next](59-clone-next.md) — Version iteration cloning workflow
- [Clone-Next Flatten](87-clone-next-flatten.md) — `--flatten` flag with DB version tracking
