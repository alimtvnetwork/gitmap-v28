# FAQ — Fix-Repo & Installers (AI Execution Prompt)

> **/goal** Provide authoritative, actionable answers and operational guidance for fix-repo scripts, version token replacement, and multi-channel installer operations.
> **/learn** Understand differences between implicit and release installer modes, offline installation requirements, numeric-overflow guards (`v17` vs `v170`), visibility change confirmation flows, and log pruning policies.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Clarify usage boundaries between `install.sh` (implicit/latest) and `release-install.sh` (pinned tag).
- [ ] `/learn` Master `--offline` installation requirements and prevent unexpected network calls (exit code 2).
- [ ] `/goal` Verify numeric-overflow guards preventing false positive version rewrites (e.g., `v17` must not match `v170`).
- [ ] `/learn` Enforce interactive confirmation guards during visibility transitions unless explicitly overridden with `--yes`.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

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

---

## Cross-References

- [Fix-Repo & Installers Index](./readme.md)
- [Fix-Repo Contract](./02-fix-repo-contract.md)
- [Installer Contract](./03-installer-contract.md)
- [Visibility Change Contract](./04-visibility-change-contract.md)
- [Acceptance Criteria](./97-acceptance-criteria.md)

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-FIX-FAQ-001: Fix-Repo and Installer Scripts FAQ Conformance

- [ ] Installer modes distinguish implicit branch tracking (`install.sh`) from immutable release downloads (`release-install.sh`).
- [ ] Offline flag (`--offline` / `--use-local-archive`) semantics enforce zero network activity and deterministic error handling.
- [ ] Numeric-overflow guard rationale is documented to prevent substring corruption of overlapping version tokens.
- [ ] Script logging limits (`--max-fix-repo-logs`) and environment variable overrides are clearly specified.
- [ ] Cross-references point to authoritative installer contracts and acceptance criteria registries.

**Given** Repository migration scripts, installer harnesses, and user operational queries.
**When** Audited against this fix-repo and installer FAQ specification.
**Then** All installer behaviors, flags, and guard mechanisms comply with specification requirements with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/06-cicd-integration/08-fix-repo-and-installers/98-faq.md --check-only
```
**Expected:** exit 0. Zero violations.
