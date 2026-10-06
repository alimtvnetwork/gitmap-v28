# Subtask 07: Image Audit, Email Blurring, Text Screenshot Cleanup and Git History Scrub

> **Parent Plan:** [230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md](../../pending/230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md)  
> **Spec Reference:** [02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/02-component-and-cli-spec.md](../../../../02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/02-component-and-cli-spec.md)  
> **Status:** `QUEUED`  
> **Target Subsystems:**  
> - `assets/screenshots/`  
> - `.ai-memory/assets/`  
> - `02-spec/`  
> - `.git/` history  

---

## 1. Technical Objective

Execute a comprehensive security and hygiene audit across all 141+ image files in the repository. Specifically:
1. Identify and remove redundant screenshots that merely capture terminal text, diffs, or console output, transcribing essential diagnostic details into standard Markdown code blocks.
2. Detect screenshots that expose developer email addresses or personal identifiers (PII), such as GitHub commit feeds or account settings.
3. Apply a strict blurring/redaction protocol to sensitive areas and save the redacted assets under **entirely new filenames** to guarantee that old, sensitive file paths never persist.
4. Formulate and verify a `git-filter-repo` rewrite protocol to permanently eradicate unblurred assets and leaked credentials from historical Git commits and packfiles.

---

## 2. Artifact Modification Inventory

| Category / Path | Action | Description & Target Assets |
|:---|:---|:---|
| `assets/screenshots/` | **Purge** | Remove 50+ text-only terminal screenshots (e.g., `pe-pipeline-error-*.png`, `commit-pull-dry-run-telemetry.png`, `pipeline-failed-shows-pass-01.png`). |
| `.ai-memory/assets/error-handling/` | **Purge** | Remove text diff screenshots (`01-releasepull-diff.png`, `02-reinstall-diff.png`, `03-rootadd-diff.png`) after ensuring diffs exist in text format. |
| `.ai-memory/assets/screenshots/` | **Redact & Rename** | Redact developer email from `14-github-238-commits.png`; save as `sanitized-github-commit-history.png`. |
| `assets/screenshots/` | **Redact & Rename** | Redact developer email from `commit-pull-prompt-notes.png` and `media_1791131607200.png`; save under distinct sanitized names. |
| `02-spec/` & `.ai-memory/` | **Update Links** | Update markdown references to point to new sanitized image filenames and remove broken links to deleted text screenshots. |
| Git History (`.git/`) | **Scrub Plan** | Formulate `git-filter-repo` paths list to rewrite commit history and remove all unredacted blobs permanently. |

---

## 3. Step-by-Step Implementation Plan

### Step 3.1: Catalog & Audit Codebase Images
1. Enumerate all image files across repository directories:
   - Run inventory script to group by format, size, and location.
   - Classify into: Diagram, UI/Branding, Text/Terminal, PII/Sensitive.

### Step 3.2: Text Screenshot Removal
1. Identify screenshots whose sole content is CLI terminal text.
2. For each identified image:
   - Check if relevant terminal output is already documented in spec markdown.
   - If not documented, transcribe the terminal text or error message into markdown code blocks.
   - Safely remove the image file using `git rm`.

### Step 3.3: Email Redaction & File Renaming Protocol
1. For images containing personal email addresses:
   - Isolate pixel bounding box containing email strings (`@gmail.com`, `@users.noreply.github.com`, etc.).
   - Apply a 24px radius Gaussian blur or an opaque solid fill (`#18181b`) matching surrounding card tones.
   - Ensure surrounding diagnostic context (commit messages, branch names, error tags) remains legible.
2. **Mandatory Renaming Rule:**
   - Save the redacted image under a new lowercase kebab-case name (e.g., `assets/screenshots/redacted-commit-notes.png`).
   - The unblurred filename is permanently slated for Git history eradication.

### Step 3.4: Reference Link Synchronization
1. Scan all markdown files in `02-spec/` and `.ai-memory/` for references to purged or renamed images.
2. Replace old image links with new sanitized paths or remove redundant image tags where text transcripts replace screenshots.

### Step 3.5: Git History Scrub Protocol (`git-filter-repo`)
1. Compile manifest of sensitive paths to purge:
   ```text
   assets/screenshots/commit-pull-prompt-notes.png
   assets/screenshots/media_1791131607200.png
   .ai-memory/assets/screenshots/14-github-238-commits.png
   assets/screenshots/_N8xDMM-6ylG.png
   assets/screenshots/0J1Th8lwMNIL.png
   assets/screenshots/8BnieUEKidCr.png
   assets/screenshots/MNRD-mOPioTv.png
   assets/screenshots/w66J-xV1NN3E.png
   ```
2. Prepare isolated backup clone: `git clone --mirror . ../gitmap-pre-scrub-backup.git`.
3. Execute `git-filter-repo --paths-from-file /tmp/git-purge-paths.txt --invert-paths --force`.
4. Run aggressive garbage collection:
   ```bash
   git reflog expire --expire=now --all
   git gc --prune=now --aggressive
   ```

---

## 4. Verification Commands & Expected Output

```bash
# 1. Count remaining images and verify categorization
find . -type f \( -name "*.png" -o -name "*.jpg" \) ! -path "*/.git/*" | wc -l

# 2. Confirm zero occurrences of purged image references in markdown
grep -rn "commit-pull-prompt-notes.png" 02-spec/ .ai-memory/ || echo "No stale references"

# 3. Confirm blob absence in git log history
git log --all --full-history -- "**/commit-pull-prompt-notes.png" || echo "Git history scrub verified"
```

---

## 5. Acceptance Criteria

- [ ] All 141+ codebase images audited and categorized.
- [ ] Text-only screenshots purged and replaced with markdown code blocks where needed.
- [ ] Images exposing developer email addresses redacted and saved under new non-colliding filenames.
- [ ] Markdown links updated across `02-spec/` and `.ai-memory/`.
- [ ] Git history rewrite procedure with `git-filter-repo` verified and ready for execution.
