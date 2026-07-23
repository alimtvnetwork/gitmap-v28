# Clone: Replace Existing Target Folder Safely

**Status:** Draft
**Audience:** Any AI or human implementing the `gitmap-v28 clone <url>` flow (or any command that materialises a Git working copy into a target directory).
**Scope:** Generic. Applies wherever an installer/cloner needs to replace an existing on-disk folder without aborting on "already exists".

---

## 1. Purpose

When `gitmap-v28 clone <url>` (or `clone-next`, or any flatten-style clone) targets a folder that **already exists** on disk, the current behaviour aborts with:

```
Error: target folder already exists: D:\wp-work\riseup-asia\scripts-fixer
```

This is wrong. The user almost always means *"replace it with a fresh clone"*. The flow MUST:

1. Detect the existing folder.
2. Try to remove it cleanly.
3. If that fails (Windows file locks, in-use binaries, permission issues), **clone into a temp folder first**, then atomically swap contents — never leave the user without either the old or the new copy.
4. Log every step with a clear, grep-friendly prefix so the user sees what's happening and why.
5. Only abort when **both** strategies fail.

---

## 2. Terminology

| Term            | Meaning                                                          |
|-----------------|------------------------------------------------------------------|
| **Target**      | The final on-disk folder where the clone must end up.            |
| **Tempclone**   | Sibling folder used as a staging area: `<target>.gitmap-tmp-<rand>`. |
| **Lock failure**| `os.RemoveAll` returns a Windows sharing-violation / EBUSY / EACCES. |
| **Swap**        | Move tempclone contents into target, replacing prior contents.   |

---

## 3. Algorithm

```
INPUT:  url, target
OUTPUT: target populated with a fresh clone of url

if target does NOT exist:
    log "[clone] target free, cloning directly into <target>"
    git clone <url> <target>
    return

# Target exists. Try the fast path first.
log "[clone] target exists: <target>"
log "[clone] strategy 1/2 — direct remove + clone"

err := removeAll(target)
if err == nil:
    git clone <url> <target>
    return

# Direct removal failed (Windows lock, in-use, permission).
log "[clone] strategy 1/2 failed: <err>"
log "[clone] strategy 2/2 — temp-clone then swap-in-place"

tempclone := "<target>.gitmap-tmp-<rand>"
ensure tempclone does not exist (remove if stale)
git clone <url> <tempclone>            # never touches target

# Empty target's contents (don't remove the folder itself — it may be the cwd
# of the user's shell or another locked handle to the directory itself).
emptyDirectoryContents(target)         # see §4

# Move every entry from tempclone → target, then drop tempclone.
moveDirectoryContents(tempclone, target)
removeAll(tempclone)

log "[clone] swap complete; target now points at fresh clone"
return

# Both strategies failed:
log "[clone] both strategies failed"
exit non-zero with a clear actionable error
```

---

## 4. Why "empty contents" instead of "remove folder"

On Windows, a directory **handle** held by any process (including the user's own PowerShell sitting in that directory) blocks `RemoveAll(target)` with a sharing violation, **even when every file inside is removable**.

The trick is to delete each child entry while leaving the directory inode itself in place. The user's shell keeps its handle, but the directory becomes empty and ready to receive new contents.

```
emptyDirectoryContents(dir):
    for each entry in dir:
        try removeAll(dir + "/" + entry)
        on failure:
            log "[clone] could not remove <entry>: <err>"
            collect into failures[]
    if failures non-empty:
        return error("could not empty <dir>: <N> entries failed")
    return ok
```

Best-effort: if a single child file is locked (e.g. `.git/index.lock` held by VS Code), **fail strategy 2 cleanly** rather than leaving a half-empty directory.

---

## 5. Atomic swap semantics

`moveDirectoryContents(tempclone, target)` MUST:

1. Iterate `tempclone`'s entries.
2. For each entry, call `os.Rename(tempclone/entry, target/entry)`.
3. On Windows, `os.Rename` across the same volume is atomic and survives most lock conditions because we're moving *into* the existing dir handle (which is what's locked), not replacing the dir itself.
4. If any rename fails, attempt a copy+remove fallback for that single entry; on persistent failure, log and continue (don't leave the target half-populated silently — surface the list at the end).

After all entries are moved, remove the now-empty `tempclone`.

---

## 6. Logging Format

All steps prefix `[clone]` for grep-ability:

```
  [clone] target exists: D:\wp-work\riseup-asia\scripts-fixer
  [clone] strategy 1/2 — direct remove + clone
  [clone] strategy 1/2 failed: remove ...: The process cannot access the file...
  [clone] strategy 2/2 — temp-clone then swap-in-place
  [clone] cloning into D:\wp-work\riseup-asia\scripts-fixer.gitmap-tmp-7f2a
  [clone] emptying target contents (12 entries)
  [clone] moving 12 entries from temp into target
  [clone] swap complete; target now points at fresh clone (gitmap-v28)
  [clone] cleaned up temp folder
```

On failure of strategy 2, the final line MUST tell the user what to do:

```
  [clone] strategy 2/2 failed: 1 entry could not be replaced (.git/index)
  [clone] hint: close any process holding the target folder and re-run, or
                'cd ..' out of the target in your shell first.
```

---

## 7. Edge Cases

| Case                                                | Behaviour                                                |
|-----------------------------------------------------|----------------------------------------------------------|
| Target does not exist                               | Direct clone. No strategy negotiation.                   |
| Target is the cwd of *this* gitmap-v28 process         | Same as user-shell cwd — strategy 2 handles it via empty-contents. |
| Target is a file, not a directory                   | Remove the file (via `RemoveAll`) and clone fresh.       |
| `tempclone` path already exists (stale)             | Remove it before cloning. Use a random suffix to minimise collisions. |
| `git clone <url> <tempclone>` fails                 | Surface git's stderr, abort. Target left untouched.       |
| Mid-swap failure (some entries moved, some not)     | Log every failure. Exit non-zero. Target is in a known-broken state — this is acceptable because the user explicitly asked for a replace and we can't roll back without doubling disk usage. |
| Cross-volume target & tempclone                     | Use `os.Rename` first; on `EXDEV` fall back to copy+remove. (Same-volume sibling avoids this.) |
| User passes `--no-replace` (future flag)            | Restore the old "abort if exists" behaviour.             |

---

## 8. Configuration Knobs

| Flag             | Default | Purpose                                            |
|------------------|---------|----------------------------------------------------|
| `--replace`      | **on**  | New default. Existing target → replace via §3.     |
| `--no-replace`   | off     | Restore abort-on-exists behaviour.                 |
| `--temp-suffix`  | `.gitmap-tmp-<rand>` | Override the staging folder suffix.   |

`--replace` being default is the right call: every observed user invocation of `gitmap-v28 clone <url>` against an existing folder so far has meant "replace it".

---

## 9. Pending-Task Integration

The replace flow MUST integrate with the existing pending-task system (`createPendingTask` / `completePendingTask` / `failPendingTask`):

- Open a single `clone` task at the start.
- On strategy-1 success: complete the task.
- On strategy-2 success: complete the task with a note `"replaced via temp-swap"`.
- On both-failed: fail the task with the last error.

This ensures `gitmap-v28 update-cleanup` can later sweep stale `*.gitmap-tmp-*` siblings.

---

## 10. Acceptance Checklist

An implementation conforms when:

- [ ] `gitmap-v28 clone <url>` against an existing folder no longer exits with "target folder already exists".
- [ ] Strategy 1 (direct remove + clone) is tried first and succeeds when no locks exist.
- [ ] Strategy 2 (temp-clone + swap) kicks in automatically on strategy-1 failure.
- [ ] Every step logs with the `[clone]` prefix.
- [ ] The temp folder is always removed on success and left in place on partial failure.
- [ ] `--no-replace` opts back into the strict abort behaviour.
- [ ] Pending tasks track both success modes and the failure mode.
- [ ] No data loss when both strategies fail (target is either fully old or fully new — never silently merged).

---

## 11. Reference Implementation Sketch (Go)

```go
func cloneReplacing(url, target string) error {
    if _, err := os.Stat(target); errors.Is(err, fs.ErrNotExist) {
        log("[clone] target free, cloning directly into %s", target)
        return gitClone(url, target)
    }

    log("[clone] target exists: %s", target)
    log("[clone] strategy 1/2 — direct remove + clone")
    if err := os.RemoveAll(target); err == nil {
        return gitClone(url, target)
    } else {
        log("[clone] strategy 1/2 failed: %v", err)
    }

    log("[clone] strategy 2/2 — temp-clone then swap-in-place")
    tmp := target + ".gitmap-tmp-" + randSuffix()
    _ = os.RemoveAll(tmp) // stale safety
    if err := gitClone(url, tmp); err != nil {
        return fmt.Errorf("git clone into temp failed: %w", err)
    }
    defer os.RemoveAll(tmp)

    if err := emptyDirContents(target); err != nil {
        return fmt.Errorf("could not empty target: %w", err)
    }
    if err := moveDirContents(tmp, target); err != nil {
        return fmt.Errorf("swap failed: %w", err)
    }
    log("[clone] swap complete; target now points at fresh clone")
    return nil
}
```

---

## 12. Out of Scope

- Backing up the old folder before replace. (User can `git stash` or copy manually if they care.)
- Merge/rebase semantics. This is a *clone*, not a sync.
- Recovering from corrupt mid-swap state automatically. We surface the failure; the user decides.

---

## See Also

- [05-cloner.md](05-cloner.md) — Existing cloner flow.
- [95-installer-script-find-latest-repo.md](95-installer-script-find-latest-repo.md) — Versioned repo discovery (related: stale URL → newer version).
- [`gitmap-v28/cmd/clonenext.go`](../../gitmap-v28/cmd/clonenext.go) — Already implements a similar fallback for the flattened-folder lock case (v2.87.0+).
