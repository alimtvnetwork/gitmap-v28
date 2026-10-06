# Subtask 01: Antigravity IDE Projects Ingestion & Pinned Projects Suite

> **Parent Plan:** `.ai-memory/plans/231-antigravity-ide-projects-and-repo-secrets-restore.md`  
> **Spec Reference:** `02-spec/21-app/231-antigravity-ide-projects-and-repo-secrets-restore/01-architecture-spec.md`  
> **Component Spec Reference:** `02-spec/21-app/231-antigravity-ide-projects-and-repo-secrets-restore/02-component-spec.md`  
> **Status:** `PENDING`  
> **Worker Assignment:** Worker 01  
> **Target Subsystems:** Antigravity IDE Configuration Store, Dual-Profile Engine, Pinned Projects Store (`pinned_projects.json`), GitMap AGY CLI Suite  
> **Relative Affected Paths:**  
> - `.antigravity_tools/instances/<instance-id>/home/.gemini/config/projects/*.json`  
> - `.antigravity_tools/instances/<instance-id>/home/.gemini/config/pinned_projects.json`  
> - `.gemini/config/projects/*.json`  
> - `.gemini/config/pinned_projects.json`  

---

## 1. Technical Objective

Ingest all 78 Git repositories across the developer workspace into Google Antigravity IDE project descriptors, organize them into 10 logical project sections for IDE sidebar navigation, configure the dual-tier pinned projects store (`pinned_projects.json`) with 23 repositories (7 Tier 1 Core and 16 Tier 2 Ecosystem), and ensure seamless dual-write synchronization between the active instance sandbox profile and the host global administrator profile.

---

## 2. Implementation Execution Plan

### Step 1: Pre-Execution Discovery & Workspace Reconciliation
1. Run `gitmap agy scan` to discover all active Git repository roots across the workspace.
2. Confirm the exact inventory count of 78 repositories (42 direct root repositories and 36 nested repositories across 9 namespaces: `02-prompts/`, `03-aukgo/`, `aukgit/`, `chris/`, `commit-fix/`, `presentations-repos/`, `project-watch-pro/`, `seo-writing/`, `web-system/`).
3. Audit existing registered project descriptors in the active instance profile (`config/projects/*.json`) to record existing UUIDs and permission grants. Specifically, ensure the following core UUIDs are strictly preserved:
   - `gitmap`: `18161e7d-2758-4bb9-82ab-317cabda949a`
   - `scripts-fixer`: `650cd964-3809-4b5f-9114-4709111c7c7c`
   - `wp-html-automate`: `d4c6435c-c96d-4afc-adb2-c12b504734ba`

### Step 2: Project Descriptors Ingestion & Dual-Profile Synchronization
1. Generate complete Antigravity project descriptor JSON files for all 78 repositories.
2. Ensure each descriptor conforms to the canonical Antigravity schema:
   - `id`: Preserved existing UUID or newly generated RFC 4122 UUID v4.
   - `name`: Clean slug identifier matching the repository folder name.
   - `projectResources.resources[0].gitFolder.folderUri`: Resolved RFC 3986 URI for the local environment.
   - `projectResources.resources[0].gitFolder.defaultBranch`: Probed default branch (`main` or `master`).
   - `permissionGrants`: `{ "v2Migrated": true }`.
   - `settings`: `{ "autoExecutionPolicy": "CASCADE_COMMANDS_AUTO_EXECUTION_EAGER" }`.
   - `isWorkspaceOnly`: `false` (affirmative boolean indicating full IDE availability).
3. Execute dual-write synchronization:
   - Write all 78 descriptors into the active instance profile: `.antigravity_tools/instances/<instance-id>/home/.gemini/config/projects/<id>.json`.
   - Mirror all 78 descriptors into the global administrator profile: `.gemini/config/projects/<id>.json`.
4. Enforce file permissions (`0644` / `rw-r--r--`) across all descriptor files.

### Step 3: Project Sections Configuration for Sidebar Navigation
1. Categorize all 78 registered repositories into the 10 logical project sections specified in `02-spec/21-app/231-antigravity-ide-projects-and-repo-secrets-restore/01-architecture-spec.md`:
   - **Section 1: Core Infrastructure & Toolchains** (7 repos): `gitmap`, `scripts-fixer`, `wp-html-automate`, `coding-guidelines`, `repo-secrets`, `antigravity-manager`, `repo-cache`.
   - **Section 2: AI Prompts & Agent Architecture** (6 repos): `02-prompts/ai-empathy-prompt-tuner`, `02-prompts/prompt-architect`, `02-prompts/prompts-connect`, `commit-fix/prompt-architect-v2`, `commit-fix/prompt-architect-v3`, `ui-prompts-cat`.
   - **Section 3: Go Core & Foundation Libraries** (7 repos): `03-aukgo/core`, `03-aukgo/enum`, `03-aukgo/errorwrapper`, `03-aukgo/extendcore`, `03-aukgo/pathhelper`, `03-aukgo/rediswrapper`, `03-aukgo/strhelper`.
   - **Section 4: System Utilities & Workstation Automation** (7 repos): `chris/winutil`, `chris/linutil`, `chris/ChrisTitusTech`, `chris/image-tools`, `macro-ahk`, `password-reset-helper-powershell`, `cat-my`.
   - **Section 5: DevOps, Clustering & Network Infrastructure** (6 repos): `aukgit/common-linux-installer`, `aukgit/kubernetes-training`, `network-fixer-hub`, `gitlogger-new`, `test-repo-bootcamp-v1`, `slack-ai-agent`.
   - **Section 6: Web Platforms, Backend & Workflow Systems** (8 repos): `letsmarknow`, `letsmarknow-ui`, `awansoft-v10`, `gstack`, `workflowy`, `workflowy-ui`, `laravel-automation`, `lara-publishing`.
   - **Section 7: Presentation Decks & Visual Design Systems** (15 repos): 13 in `presentations-repos/`, `slide-sensei-37`, `img-pdf`.
   - **Section 8: Content Automation, SEO & Media Production** (6 repos): `seo-packages`, `seo-writing/alim-seo-writing`, `yt-descriptions`, `kita-social-media-content-calender`, `punam-case-studies-v1`, `movie-cli`.
   - **Section 9: WordPress Ecosystem & Publishing Plugins** (5 repos): `wp-exam`, `wp-git-log`, `wp-link-manager`, `wp-onboarding`, `lara-licensing`.
   - **Section 10: Fleet Monitoring, Portfolios & Digital Identity** (11 repos): `project-watch-pro/project-watch-pro`, `project-watch-pro/pwp-mobile`, `alim-cv`, `alim-karim-profile`, `alim-status-sample`, `digital-name-card`, `aukgit/alim.karim.profile`, `web-system/sweet-digs-finder`, `spec-builder`, `riseup-asia-website-project`, `icon-coding-guidelines`.

### Step 4: Pinned Projects Hierarchy Configuration
1. Initialize `pinned_projects.json` with positive boolean properties (`isSyncEnabled: true`, `isPinned: true`).
2. Populate the 23 priority repositories across two defined tiers:
   - **Tier 1: Mandatory Core Orchestrators (7 Repositories):**
     1. `gitmap` (Priority 1)
     2. `scripts-fixer` (Priority 2)
     3. `wp-html-automate` (Priority 3)
     4. `coding-guidelines` (Priority 4)
     5. `repo-secrets` (Priority 5)
     6. `antigravity-manager` (Priority 6)
     7. `repo-cache` (Priority 7)
   - **Tier 2: Key Ecosystem Services (16 Repositories):**
     8. `chris/winutil` (Priority 8)
     9. `chris/linutil` (Priority 9)
     10. `03-aukgo/core` (Priority 10)
     11. `03-aukgo/errorwrapper` (Priority 11)
     12. `02-prompts/prompt-architect` (Priority 12)
     13. `project-watch-pro/project-watch-pro` (Priority 13)
     14. `presentations-repos/hiltrax` (Priority 14)
     15. `presentations-repos/slides-spec` (Priority 15)
     16. `cat-my` (Priority 16)
     17. `movie-cli` (Priority 17)
     18. `network-fixer-hub` (Priority 18)
     19. `laravel-automation` (Priority 19)
     20. `letsmarknow` (Priority 20)
     21. `letsmarknow-ui` (Priority 21)
     22. `workflowy` (Priority 22)
     23. `seo-writing/alim-seo-writing` (Priority 23)
3. Persist `pinned_projects.json` simultaneously to:
   - Active instance profile: `.antigravity_tools/instances/<instance-id>/home/.gemini/config/pinned_projects.json`
   - Global administrator profile: `.gemini/config/pinned_projects.json`

---

## 3. Verification & Validation Commands

Worker 01 must execute the following non-git validation routines to confirm operational integrity:

```bash
# 1. Verify project registration status across all 78 repositories
gitmap agy scan

# 2. Confirm total registered project count in active profile
gitmap agy ls

# 3. Verify pinned projects count and tiering
gitmap agy pins ls --json

# 4. Verify dual-profile descriptor count parity
# Instance profile descriptor count:
ls -1 "$HOME/.gemini/config/projects/"*.json | wc -l
# Global admin profile descriptor count (when on host):
ls -1 "${HOME_ADMIN:-$HOME}/.gemini/config/projects/"*.json 2>/dev/null | wc -l

# 5. Verify pinned_projects.json schema and positive boolean properties
python3 -c "
import json
with open('$HOME/.gemini/config/pinned_projects.json') as f:
    data = json.load(f)
assert data.get('isSyncEnabled') is True, 'isSyncEnabled must be true'
assert len(data.get('projects', [])) == 23, f'Expected 23 pins, found {len(data.get(\"projects\", []))}'
print('Pinned projects store verified successfully: 23 projects')
"
```

---

## 4. Acceptance Criteria Checklist

- [ ] **AC-01**: Exactly 78 repository descriptors generated and active in the instance projects directory (`78 / 78 added`).
- [ ] **AC-02**: Existing project UUIDs (`gitmap`, `scripts-fixer`, `wp-html-automate`) strictly preserved without churn.
- [ ] **AC-03**: Dual-profile parity established between active instance and global profile stores.
- [ ] **AC-04**: 10 logical project sections mapped across all 78 repositories.
- [ ] **AC-05**: Exactly 23 repositories pinned in `pinned_projects.json` (7 Tier 1 Core + 16 Tier 2 Ecosystem).
- [ ] **AC-06**: Positive booleans exclusively utilized (`isPinned`, `isActive`, `hasGitFolder`, `isSyncEnabled`, `isWorkspaceOnly`).
- [ ] **AC-07**: Strictly relative Git paths referenced in all task documentation with zero absolute paths or drive letters.
- [ ] **AC-08**: `gitmap agy ls` and `gitmap agy pins ls` execute cleanly without warnings or errors.
