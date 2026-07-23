# Git History Rewrite: Removing a File & Pinning a File to Its Latest State

> Research note. Audience: GitMap maintainers / contributors who need to clean
> up history before a public release or before importing a third-party repo.
>
> **Both operations rewrite history.** Every commit hash downstream of the
> rewrite changes. Anyone who has cloned the repo will need to re-clone or do
> a hard reset. Coordinate with collaborators before pushing.

---

## TL;DR

| Goal | Recommended Tool | One-liner |
|------|------------------|-----------|
| Remove a file from **all** history (e.g. leaked secret, huge binary) | `git filter-repo` | `git filter-repo --invert-paths --path secrets.env` |
| Same, simpler for "purge by name/size" | **BFG Repo-Cleaner** | `bfg --delete-files secrets.env` |
| Replace a file's content in every past commit with the **current** version (so history shows only the latest state) | `git filter-repo --blob-callback` (or `--replace-text`) | see Section 2 |
| Legacy / no install allowed | `git filter-branch` | slow, error-prone, **discouraged** by Git itself |

`git filter-repo` is the modern, officially-recommended replacement for
`git filter-branch`. The Git docs themselves point users away from
`filter-branch` due to safety and performance issues.

---

## 0. Preflight (do this every time)

```bash
# 1. Make a fresh mirror clone — NEVER rewrite history in your working repo first.
git clone --mirror git@github.com:org/repo.git repo-rewrite.git
cd repo-rewrite.git

# 2. Confirm you have a backup somewhere off-machine.
# 3. Tell collaborators: "force-push incoming, do not push for the next hour".
```

After the rewrite you will `git push --force` (or `--force-with-lease`) the
mirror back to the remote.

Install once:

```bash
# Recommended:
pip install git-filter-repo
# or: brew install git-filter-repo
# or: scoop install git-filter-repo   (Windows)

# Optional alternative:
# Download BFG jar from https://rtyley.github.io/bfg-repo-cleaner/
```

---

## 1. Question A — Remove a file from **all** history

### 1a. With `git filter-repo` (recommended)

Remove a single file across every commit and every branch:

```bash
git filter-repo --invert-paths --path path/to/secret.env
```

Remove multiple paths in one pass:

```bash
git filter-repo \
  --invert-paths \
  --path path/to/secret.env \
  --path path/to/old-binary.zip \
  --path-glob '*.pem'
```

- `--invert-paths` flips the meaning: "keep everything **except** these paths".
- Without `--invert-paths`, the same flags would *keep only* those paths.
- Tags and all branches are rewritten automatically.

### 1b. With BFG (fastest, simplest for "delete by name")

```bash
java -jar bfg.jar --delete-files secret.env
java -jar bfg.jar --delete-folders node_modules
java -jar bfg.jar --strip-blobs-bigger-than 50M

# BFG only rewrites; it does not GC. Finish with:
git reflog expire --expire=now --all
git gc --prune=now --aggressive
```

BFG is 10–720× faster than `filter-branch` for this exact task, but it
protects your **latest** commit (it refuses to touch HEAD). So first commit a
clean version that no longer references the file, then run BFG to scrub
history.

### 1c. Legacy: `git filter-branch` (avoid if possible)

```bash
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch path/to/secret.env" \
  --prune-empty --tag-name-filter cat -- --all
```

Git's own documentation now opens with a warning recommending `filter-repo`
instead. Only use this when you truly cannot install anything.

### 1d. Push & cleanup

```bash
git push --force --all
git push --force --tags

# Local repo cleanup (after filter-repo / filter-branch):
git reflog expire --expire=now --all
git gc --prune=now --aggressive
```

If the leaked content was a **secret**, also:

1. Rotate / revoke the secret immediately — the rewrite does not delete data
   already cloned, cached by GitHub, or scraped by bots.
2. Ask GitHub Support to purge cached views and pull-request diffs.
3. Audit access logs.

---

## 2. Question B — Pin a file to its **current** state across all history

> "I have a file `X` with 10 historical states. I want every past commit that
> touched `X` to instead contain the version-10 (current) content."

This is **content rewriting**, not path removal. You are not deleting the
file's history entries — you are overwriting the *content* stored at each
historical revision so that, at every commit where `X` exists, its bytes
equal the current bytes.

The cleanest tool is `git filter-repo --blob-callback`.

### 2a. Strategy

1. Compute the SHA-1 (blob id) of every historical version of `X`.
2. Read the current bytes of `X` from your working tree once.
3. In a `--blob-callback`, when the incoming blob's `original_id` matches any
   of those historical ids, replace `blob.data` with the current bytes.

### 2b. Step-by-step

```bash
# (in your normal working clone, just to gather info)

# 1. List every blob hash that file X ever had:
git log --all --pretty=format: --raw -- path/to/X \
  | awk '{print $4}' | grep -v '^$' | sort -u > /tmp/x-blobs.txt

# 2. Save the CURRENT content of X somewhere outside the repo:
cp path/to/X /tmp/x-current.bin

# 3. Switch to a fresh mirror clone (see Preflight).
cd /path/to/repo-rewrite.git
```

### 2c. Run the rewrite

```bash
git filter-repo --blob-callback '
import os

# Load once, cache on the function object.
if not hasattr(blob_callback, "_payload"):
    with open("/tmp/x-current.bin", "rb") as f:
        blob_callback._payload = f.read()
    with open("/tmp/x-blobs.txt", "rb") as f:
        blob_callback._targets = {
            line.strip() for line in f if line.strip()
        }

if blob.original_id in blob_callback._targets:
    blob.data = blob_callback._payload
'
```

Result: every commit that previously introduced or modified `X` still exists,
still touches `X`, and still has the same author/date/message — but the
*content* of `X` at that revision is now identical to the latest state.

### 2d. Variant: simple text substitution

If you only want to scrub specific strings (e.g. an API key) inside the file
rather than overwrite the whole file, prefer `--replace-text`:

```bash
# rules.txt
OLD_API_KEY==>REDACTED
regex:AKIA[0-9A-Z]{16}==>REDACTED

git filter-repo --replace-text rules.txt
```

### 2e. Variant: collapse N historical states into a single commit

If the goal is "I do not even want the *commits* that touched the old
versions of X", that is a different operation: an **interactive rebase** or
`git filter-repo --commit-callback` that drops or squashes those commits.
The blob-callback approach above keeps the commit graph shape intact, which
is usually what you want for auditability.

---

## 3. Verification checklist

After either operation, before force-pushing:

```bash
# Confirm the file is gone (Question A):
git log --all --oneline -- path/to/secret.env   # should be empty

# Confirm every revision of X now matches current (Question B):
for sha in $(git log --all --pretty=format:%H -- path/to/X); do
  git show "$sha:path/to/X" | sha256sum
done | sort -u
# Expect a SINGLE unique sha256 across all commits.

# Repo size sanity:
git count-objects -vH
```

---

## 4. Pitfalls

- **Force-push hazard.** Every downstream clone diverges. Coordinate.
- **GitHub caches.** PR diffs and old refs may still expose removed content.
  Open a support ticket if the file was a leaked secret.
- **Submodules / LFS pointers.** `filter-repo` handles these but verify
  manually before pushing.
- **Signed commits / signed tags.** Signatures break on rewrite — re-sign if
  required.
- **CI / deploy keys pinned to old SHAs.** Update them after the push.
- **`filter-branch` is officially discouraged.** Reach for `filter-repo` or
  BFG first.

---

## 5. References

- `git filter-repo` — <https://github.com/newren/git-filter-repo>
- Real-world recipes — <https://mintlify.com/newren/git-filter-repo/examples/user-scenarios>
- Content-based filtering guide — <https://www.mintlify.com/newren/git-filter-repo/guides/content-based-filtering>
- BFG Repo-Cleaner — <https://rtyley.github.io/bfg-repo-cleaner/>
- Git's own `filter-branch` docs (with the "use filter-repo instead" warning)
  — <https://git-scm.com/docs/git-filter-branch>
