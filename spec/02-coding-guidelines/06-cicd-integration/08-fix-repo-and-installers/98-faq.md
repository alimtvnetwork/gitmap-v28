# FAQ — Fix-Repo & Installers

> **Version:** 1.0.0 · **Updated:** 2026-04-28

---

### Q1. When do I use `install.sh` vs `release-install.sh`?

- **`install.sh`** — implicit mode (latest from a branch) is the default.
  Use this from the repo's `raw.githubusercontent.com` URL when you want
  the always-current version, or with `--version vX.Y.Z` when you want a
  pin.
- **`release-install.sh`** — pinned mode only. Use this when you download
  the asset from a GitHub Release page; the tag is baked in at release
  time so the same asset always installs the same version forever.

### Q2. How do I install a specific older version?

```bash
./install.sh --version v3.21.0

# OR

curl -fsSL https://github.com/<o>/<r>/releases/download/v3.21.0/release-install.sh | bash
```

### Q3. Can I run the installer offline?

Yes — pass `--offline` (alias `--use-local-archive`). It will skip every
network operation and require a pre-staged archive at the expected path.
If any code path tries to hit the network in offline mode, exit code `2`
is raised.

### Q4. How do I rename my fork from v17 → v18 without breaking links?

1. Rename the GitHub repo (or fork into a new `-v18` repo).
2. Run `./fix-repo.sh` (Bash) or `.\fix-repo.ps1` (PowerShell).
3. Commit the result. URLs, install one-liners, badges — everything that
   contains `<base>-v17` is rewritten to `<base>-v18`.

`fix-repo` reads its identity from `git remote get-url origin`, so make
sure the rename is reflected in your local clone first
(`git remote set-url origin https://github.com/<owner>/<base>-v18`).

### Q5. Why doesn't `fix-repo` touch `coding-guidelines-v170`?

The numeric-overflow guard. If the token were a plain substring with no
trailing-digit check, `v17` would corrupt `v170`, `v171`, etc. The guard
is mandated by §5.3 of the contract.

### Q6. How do I keep only the last 5 fix-repo logs?

```bash
./install.sh --run-fix-repo --max-fix-repo-logs 5

# or via env:

INSTALL_MAX_FIX_REPO_LOGS=5 ./install.sh --run-fix-repo
```

### Q7. The visibility-change script asks me to confirm — can I skip?

Only when going `private → public`. Pass `--yes` (or `-Yes` on
PowerShell) to skip. The prompt cannot be bypassed in any other
direction because there is no destructive transition to confirm.

### Q8. Can I use `visibility-change` on Bitbucket / Gitea?

Not yet. Provider detection only knows GitHub and GitLab. PRs welcome —
add a branch to `scripts/visibility-change/provider.sh` and wire up the
matching CLI in `apply.sh`.

### Q9. Where do I report a CI/CD installer bug?

Run `./linters-cicd/run-all.sh --path . --output report.sarif` and open
an issue with the SARIF file attached. For installer-specific bugs,
also include the output of `./install.sh --dry-run --version vX.Y.Z`.
