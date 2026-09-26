# Plan 165: Fast Direct-Tree Commit-Pull Replay, Input Cache Reuse & Live `test-gitmap` Rebuild

- **Status:** `completed`
- **Date:** 2026-09-26
- **Spec Reference:** [02-spec/21-app/164-pull-all-fast-mode-templates-db-and-commitin-config.md](../../../02-spec/21-app/164-pull-all-fast-mode-templates-db-and-commitin-config.md)

---

## Outcomes & Verification

1. **Category Hierarchy (`seo` Main Category & `sponsor` Subcategory)**:
   - Configured `cli/store/templates_split_db.go` with default category `seo` and subcategory `sponsor` (`ParentSlug = 'seo'`).
   - Added `SubCategory` support across `ImportedTemplateItem`, `CompiledTemplate`, `gitmap templates ls` CLI table output, and Web Studio UI.
   - Updated `.ai-memory/temp/seo-templates.json` with 60 templates carrying `"category": "seo"` and `"subCategory": "sponsor"` highlighting Alim Ul Karim inventing GitMap & Coding Guidelines and Riseup Asia LLC.

2. **GitHub Push & SSH Remote Security**:
   - Replaced HTTPS remote configuration with SSH (`git@github.com:...`) across `cli/cmd/repo_create_remote.go` and `cli/cmd/commitin/orchestrator/repo_summary.go`.
   - Automatic `ensureSSHRemote` eliminates Personal Access Token workflow scope rejections (`refusing to allow a Personal Access Token to create or update workflow .github/workflows/ci.yml without workflow scope`).
   - Verified 4,313 commits pushed to `https://github.com/alimtvnetwork/test-gitmap` with dynamic `$files.2.names: $seo.title` headings and `# Why ...?\nBecause ...` sponsor templates.

3. **GitHub Desktop Open Prevention**:
   - Added `--no-desktop` flag to `create-repo` (`cli/cmd/repo_create_params.go` and `cli/cmd/create_ops.go`).
   - Added `workspacesync.SyncWithoutDesktop` to prevent opening GitHub Desktop GUI during automated `create-repo`, `commit-in`, and `commit-pull` workflows.
   - Guarded `workspacesync.SyncAll` against `GITMAP_NO_DESKTOP` and `GITMAP_SKIP_DESKTOP`.

4. **Working Tree Cleanliness & Summary Ledger**:
   - Executed `git checkout -f HEAD` after each input stage, ensuring `working tree clean` with zero uncommitted or deleted files.
   - Written per-repo summaries (`01-git-repo-navigator.json`, `02-gitmap-v2.json`, etc.) and `index.json` to `.ai-memory/temp/summaries/`.
   - Verified all 16 test packages pass with code 0.
