# Commit Author Information Specification
**Sequence ID:** 335  
**Component:** `cli/cmdcommit`  

---

## 1. Architectural Changes
1. **Metadata Extension:**
   - Update `CommitMetadata` in `cli/cmdcommit/commit_types.go` to include `Author string`.
2. **Flag Parsing & Author Resolution:**
   - Extend `cli/cmdcommit/commit_flags.go` to support `--author` and `-a`.
   - In `resolveAuthorFallback(metadata *CommitMetadata)`, query `git config user.name` and `git config user.email` if `metadata.Author` is not explicitly provided, assembling `"Author Name <author.email@example.com>"`.
3. **Commit Invocation:**
   - In `cli/cmdcommit/commit_push.go`, when `metadata.Author` is non-empty, include `--author=<author>` in `git commit` arguments.
4. **Visual UI Attribution:**
   - In `cli/cmdcommit/commit_ui.go`, add an `Author` field row to `renderCommitPushSummaryCard`, styling the author name and email consistently with terminal design standards.

---

## 2. Behavioral Contract
- `gitmap cp -m "feat: new stuff" --author "Jane Doe <jane@example.com>"` explicitly sets the author.
- Default invocation (`gitmap cpf "module - feat"`) automatically detects and displays author from Git config in the commit summary card.
