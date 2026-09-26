# Zero-Storage Actions Purge & Storage Governance — Workflow (must follow)

Trigger Keywords & Aliases: `purge actions storage`, `zero storage mandate`, `purge artifacts`, `purge actions cache`, `clean actions storage`, `github actions storage purge`

> [!IMPORTANT]
> Prompt Version: 1.0.0
> Synchronization: Main Meta-Repo & Connected Workspaces
>
> **Zero-Storage Actions Mandate:**
> GitHub Actions shared storage (artifacts and caches) must be maintained at 0.0 GB (or strictly within the 0.5 GB free quota). Any workflow generating artifacts must enforce `retention-days: 1` and clean up post-release. Any repository hoarding unpurged artifacts or multi-gigabyte build caches must be purged immediately using the autonomous purge automation.

- [A] Purge Automation Script: `03-ai-scripts/34-purge-github-actions-artifacts.py`
- [B] Purge Workflow: `.github/workflows/purge-actions-artifacts.yml`
- [C] Governing Spec: `02-spec/12-cicd-pipeline-workflows/`
- [D] Prior Issues: `.ai-memory/cicd-issues/`

---

## Autonomous Storage Purge & Hygiene Workflow

### Mandatory Storage Governance Rules

1. **Top-Instruction Precedence:** Any custom instruction provided in chat or incoming prompt strictly overrides default guidelines.
2. **Defensive Git Pull First:** ALWAYS run `git pull` on the target repository before modifying or pushing workflows/scripts to prevent merge conflicts.
3. **Dual Purge Scope:** Purge BOTH remote GitHub Actions artifacts and stale GitHub Actions caches (`actions/caches`). Never leave multi-gigabyte compiler target caches accumulating indefinitely.
4. **Enforce `retention-days: 1` on Intermediate Artifacts:**
   Every `actions/upload-artifact@v4` step used solely for cross-job handoffs (e.g. multi-platform build handoffs to a release job) MUST declare:
   ```yaml
   retention-days: 1
   ```
5. **Post-Release Automated Cleanup:**
   Release workflows publishing packages to GitHub Releases must execute an immediate post-release artifact deletion step:
   ```yaml
   - name: Purge Run Artifacts Post-Release
     if: always() && startsWith(github.ref, 'refs/tags/')
     env:
       GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
     run: |
       python3 03-ai-scripts/34-purge-github-actions-artifacts.py --repo "${{ github.repository }}" --artifacts-only || true
   ```
6. **Ban on `cache-all-crates: "true"` and Unbounded Object Dumps:**
   Never cache complete local build target directories containing compiled objects (`target/`, `~/.cache/go-build`, etc.) across multiple matrix runners. Rely on package/dependency caches (`go.sum`, `Cargo.lock`, `package-lock.json`) instead.
7. **Scheduled Purge Workflow:**
   Ensure each repository includes `.github/workflows/purge-actions-artifacts.yml` configured to trigger on schedule (nightly at 02:00 UTC) and manual `workflow_dispatch`.
8. **Final Step Atomic Commit & Immediate Push:**
   Stage all changed scripts and workflow definitions together and commit atomically (`git commit -m "ci(actions): enforce zero-storage mandate and add purge automation"`). Immediately push to origin (`git push origin <branch>`).

---

## Execution Prompt (Copy & Paste Trigger)

```markdown
- [A] Purge Automation Script: `03-ai-scripts/34-purge-github-actions-artifacts.py`
- [B] Purge Workflow: `.github/workflows/purge-actions-artifacts.yml`

Execute the Zero-Storage Actions Mandate on repository `<target-repo-or-path>`.
1. Run `git pull` to ensure local synchronization with 0 merge conflicts.
2. Inspect `.github/workflows/` for unbounded `upload-artifact` (missing `retention-days: 1`), bloated build caches (`cache-all-crates`, massive object dumps), or missing purge triggers.
3. Eliminate excessive storage definitions and enforce `retention-days: 1`.
4. Install/update `[A]` and `[B]`.
5. Execute `python [A] --repo <owner/repo>` to immediately purge all stored artifacts and stale caches via GitHub API.
6. Verify YAML and test scripts.
7. Commit atomically and push immediately to `origin <branch>`.
```
