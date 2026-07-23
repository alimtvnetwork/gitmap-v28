# History Rewrite: `history-purge` & `history-pin`

> Status: Spec, v4.15.0. Implementation lives in
> `gitmap-v28/cmd/historyrewrite*.go`. Research backing this spec:
> `spec/15-research/git-history-rewrite-remove-and-pin-file.md`.

Two CLI commands wrap `git filter-repo` in a **mirror-clone sandbox**
so the user's working repository is never rewritten in place.

| Command           | Alias  | Purpose                                              |
|-------------------|--------|------------------------------------------------------|
| `history-purge`   | `hp`   | Remove file(s) and folder(s) from **all** history    |
| `history-pin`     | `hpin` | Pin file(s) to current content across **all** history |

## 1. Synopsis

```
gitmap-v28 history-purge <path> [<path> ...] [flags]
gitmap-v28 hp            <path> [<path> ...] [flags]

gitmap-v28 history-pin   <path> [<path> ...] [flags]
gitmap-v28 hpin          <path> [<path> ...] [flags]
```

Multiple paths may be passed as separate args, joined by `,`, or
joined by `, ` (comma-space). Quoting is irrelevant — the parser
normalizes all three forms:

```
gitmap-v28 hp secret.env build/cache.bin
gitmap-v28 hp "secret.env, build/cache.bin"
gitmap-v28 hp secret.env,build/cache.bin
```

Folders and files are both accepted. Folder paths are passed to
`filter-repo` as directory prefixes (everything under them is purged
or pinned).

## 2. Mirror-clone sandbox flow

Every invocation runs the same five phases. None of them ever touch
the user's working repository.

1. **Identify origin** — `git remote get-url origin` in the cwd. Hard
   fail (exit `2`) when not in a repo or no `origin` remote.
2. **Mirror-clone** — `git clone --mirror <origin>` into a fresh
   `os.MkdirTemp` sandbox. Cleaned up via `defer` on every exit path.
3. **Filter-repo run** — invoke `git filter-repo` inside the sandbox
   with the per-command flag set (see §3 / §4). Hard fail (exit `5`)
   on non-zero exit from `filter-repo`.
4. **Verification** — re-run the canonical verification one-liner
   from `spec/15-research/.../§3` against the rewritten sandbox. Hard
   fail (exit `6`) when verification disagrees with the requested
   operation.
5. **Push prompt** — print a verification-passed banner (mode, path
   count, sandbox path, remote URL, warning) then prompt
   `Type 'yes' to force-push to <origin> (anything else aborts): `.
   The user must type the literal token `yes` (case-sensitive). Any
   other input — including empty, `y`, `Y`, `YES` — aborts. `--yes`
   skips the prompt; `--no-push` short-circuits before the prompt and
   prints the manual `git push` command instead.

## 3. `history-purge` — remove paths from all history

Wraps `git filter-repo --invert-paths --path P1 --path P2 ...`. After
the run, every commit that previously touched `P` no longer references
it; the blob is unreachable and will be GC'd by `git gc`.

### Verification

For each requested path `P`:

```
git -C <sandbox> log --all --oneline -- <P>   # MUST be empty
```

Any non-empty output exits `6` and aborts before the push prompt.

## 4. `history-pin` — pin paths to current content

Wraps `git filter-repo --blob-callback '<python>'` where the Python
callback reads the **current bytes** of each requested path from the
user's working tree and overwrites every historical blob hash that
previously belonged to that path.

### Algorithm (per path `P`)

1. Read current bytes from working tree: `data := os.ReadFile(P)`.
2. Enumerate every historical blob SHA of `P`:
   ```
   git log --all --pretty=format: --raw -- P | awk '{print $4}'
   ```
3. Stage `(P, data, sha_set)` into a JSON manifest written to a
   tempfile.
4. Generate a single `--blob-callback` Python script that loads the
   manifest, then for each incoming blob: if `blob.original_id` is in
   any path's sha_set, replace `blob.data` with that path's bytes.

### Verification

For each requested path `P`:

```
for sha in $(git -C <sandbox> log --all --pretty=format:%H -- P); do
  git -C <sandbox> show "$sha:P" | sha256sum
done | sort -u | wc -l   # MUST be exactly 1
```

## 5. Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-y`, `--yes` | `false` | Skip the push confirmation prompt; force-push immediately on success. |
| `--no-push` | `false` | Stop after verification; print the manual `git push --force-with-lease` command. Mutually exclusive with `--yes`. |
| `--dry-run` | `false` | Run mirror-clone + filter-repo + verification, then **always** stop without pushing and print the sandbox path for inspection. |
| `--message <s>` | `""` | Rewrite the commit message to `<s>` **only** for commits whose `file_changes` reference one of the requested paths (exact match, or descendants when a path is a folder). Untouched commits keep their original message verbatim. Empty = no message rewrite. Implemented via `filter-repo --commit-callback` that inspects `commit.file_changes` against the requested path set. |
| `--keep-sandbox` | `false` | Don't delete the temp mirror-clone on exit. Path is printed. Useful for debugging a verification failure. |
| `-q`, `--quiet` | `false` | Suppress per-phase progress lines; only print errors and the final summary. |

## 6. Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success (rewrite + verification + push, OR dry-run completed cleanly) |
| `2` | Not in a git repo, or `origin` remote missing |
| `3` | `git filter-repo` not installed (with install hint) |
| `4` | Bad args (zero paths, conflicting flags, unreadable file for `pin`) |
| `5` | `git filter-repo` returned non-zero |
| `6` | Verification disagreed with the requested operation |
| `7` | Push failed |

## 7. Dependency: `git filter-repo`

`filter-repo` is not bundled with Git. On first invocation we shell
out to `git filter-repo --version`. If the binary is missing we exit
`3` and print the OS-appropriate install command:

- Linux: `pip install --user git-filter-repo`
- macOS: `brew install git-filter-repo`
- Windows: `scoop install git-filter-repo` (or `pip install`)

No auto-install — explicit user consent is required for a Python
dependency.

## 8. Safety guarantees (non-negotiable)

- Working repo is **never** rewritten — all mutation happens in a
  `os.MkdirTemp` sandbox cleaned up via `defer`.
- Verification runs **before** the push prompt; a failed verification
  cannot push.
- `--dry-run` short-circuits before push regardless of `--yes`.
- `--no-push` and `--yes` are mutually exclusive (exit `4`).
- The sandbox path is printed on every error path so the user can
  inspect it (`--keep-sandbox` keeps it permanently).

## 9. CI smoke tests

`.github/workflows/history-rewrite-smoke.yml` runs two end-to-end
scenarios on `ubuntu-latest`:

- **purge-secret**: build a temp git repo with 5 commits incl.
  `secret.env`, run `gitmap-v28 hp secret.env --yes --no-push`, assert
  `git log --all -- secret.env` is empty and `git count-objects -v`
  shrunk.
- **pin-current**: build a temp git repo where `X` has 3 historical
  states, run `gitmap-v28 hpin X --yes --no-push`, assert the verification
  loop produces exactly one sha256.

The workflow installs `git-filter-repo` via `pip install
--user git-filter-repo` and runs both scenarios in under 60 s.

## 10. References

- Research note: `spec/15-research/git-history-rewrite-remove-and-pin-file.md`
- `git filter-repo`: https://github.com/newren/git-filter-repo
- Memory: `mem://features/history-rewrite`