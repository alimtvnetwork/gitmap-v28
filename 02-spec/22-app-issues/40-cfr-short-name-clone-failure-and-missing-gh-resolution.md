# Issue 40: CFR Short-Name Clone Failure & Missing GitHub CLI Resolution

> **Issue Status:** Resolved  
> **Traceability IDs:** Task-08  
> **Canonical Path:** `02-spec/22-app-issues/40-cfr-short-name-clone-failure-and-missing-gh-resolution.md`  
> **Parent Spec:** `02-spec/22-app-issues/01-index.md`

---

## 1. Reproduction & Symptoms

When running `gitmap cfr <name>` or `gitmap clone <name>` with a bare repository slug:
```text
PS D:\work\project-watch-pro> gitmap cfr pwp-mobile
Cloning pwp-mobile into pwp-mobile...
  [clone] target free, cloning directly into D:\work\project-watch-pro\pwp-mobile

  ▸ git clone  pwp-mobile
    target  D:\work\project-watch-pro\pwp-mobile
    exec    git clone pwp-mobile D:\work\project-watch-pro\pwp-mobile
fatal: repository 'pwp-mobile' does not exist

  ✖ git clone failed
    command  git clone pwp-mobile D:\work\project-watch-pro\pwp-mobile
    exit     128
    error    exit status 128
```

---

## 2. Root Cause Analysis (RCA)

`executeDirectClone` assumes the provided input argument is always a valid git URL (starts with `http://`, `https://`, `git@`, etc.). When a short repository name like `pwp-mobile` is provided, `git clone pwp-mobile` fails with `fatal: repository 'pwp-mobile' does not exist` (exit code 128).

GitMap did not query GitHub CLI (`gh repo list` or `gh repo view`) or the local SQLite cache to resolve the full remote URL before cloning. Furthermore, cloned repositories were not automatically synchronized with VS Code Project Manager and GitHub Desktop.

---

## 3. Grounded Code Fix

1. **Short-Name Resolution:** Before passing the target to `git clone`, detect if it is a bare repository slug without protocol.
2. **GitHub CLI Query:** Query `gh repo list --json name,url,nameWithOwner` to resolve matching repositories for the authenticated user/organization.
3. **SQLite `repodb` Cache:** Persist discovered repository names and URLs into a local SQLite repository cache to accelerate subsequent queries.
4. **Registration Automation:** Ensure successful direct clones automatically register in GitHub Desktop and VS Code Project Manager (`projects.json`).
