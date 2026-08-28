# AI Implementation Spec: File Manipulation Commands (Lowercase & Fix Sequence)

## Overview
This specification provides strict instructions for an AI to implement robust file manipulation commands in a CLI tool (e.g., GitMap). The target capabilities include mass-renaming files to lowercase (while respecting Git tracking and ignore patterns) and automatically re-sequencing files in folders based on various ordering strategies.

## Non-Negotiable Rules for the Executing AI
1. **Generic Implementation**: Do not hardcode paths to specific local folders or assume framework internals outside of standard CLI parsing and filesystem utilities.
2. **Boolean Naming**: All boolean variables MUST begin with `is`, `has`, `can`, or `should`.
3. **No Garbage Names**: Do not use generic variables like `data`, `obj`, `temp`.
4. **Git Awareness**: Operations that move or rename files MUST be aware of the underlying Git index (e.g., using `git mv` instead of a plain `os.Rename`) so that the Git history retains the continuity of the file.
5. **Windows Path Normalization**: The implementation MUST normalize paths and handle Windows long-path limitations (e.g., prefixing `\\?\` internally if the host language requires it).

---

## Command 1: Lowercase Renamer
**Command Pattern**:
`cli-tool lowercase <source_pattern> <target_pattern> -except "<paths>"`

**Capabilities required**:
1. Take a source filename/pattern (e.g., `"OLD.md"`) and rename instances to a target (e.g., `"old.md"`).
2. Support an `-except` flag taking a comma-separated list of paths or globs to ignore (e.g., `"/path/file", "node_modules/*", ".git/*"`).
3. **Default Ignore Profile**: Provide a flag (e.g., `-ignore default`) that automatically excludes standard volatile/system directories (e.g., `node_modules`, `.git`).
4. **Git History Fix**: The rename operation MUST use the native Git API or shell out to `git mv` if the file is tracked. This prevents the filesystem from showing an "untracked file" and a "deleted file" while keeping the rename atomic in the git history.
5. **Help/UI Examples**: The CLI's help output MUST include explicit examples of these capabilities.

**Example Help Output Checklist**:
- [ ] Show `cli-tool lowercase "OLD.md" "old.md" -except "node_modules/*"`
- [ ] Show `cli-tool lowercase -ignore default`

---

## Command 2: Fix File Sequencing (`fix-seq-files` / `fsf`)
**Command Pattern**:
`cli-tool fix-seq-files <folder1> <folder2> [flags]`

**Capabilities required**:
1. Scan specified folders and identify numbered/sequenced files (e.g., `01-draft.md`, `02-notes.md`).
2. **Order Strategies**:
   - `-orderbytime`: Re-sequence files sequentially based on their filesystem modification or creation time.
   - `-orderbyaz`: Re-sequence files alphabetically based on the non-sequence portion of their name.
3. **Keep Old Order**: Support a `-keep-old-order` flag. If two files end up with colliding sequence numbers, fallback to alphabetization or time to resolve the tie without destroying the existing relative sequence order.
4. **Fixated / Pinned Sequences**: Allow users to explicitly pin or "fixate" a specific sequence number to a specific filename base (e.g., `-pin "draft=01,notes=02"`). If a new file is added, it correctly increments around the pinned files.
5. **Path Normalization**: The system MUST normalize absolute/relative paths and gracefully handle Windows MAX_PATH limits.

**Example Help Output Checklist**:
- [ ] Show `cli-tool fix-seq-files /folder1, folder2 -orderbytime`
- [ ] Show `cli-tool fix-seq-files /folder1 -orderbyaz -keep-old-order`
- [ ] Show an example pinning a sequence: `cli-tool fsf /folder1 -pin "readme=00"`

---

## Execution Checklist for the AI
Before submitting the code, the executing AI MUST verify:
- [ ] I have implemented `lowercase` with Git-native moving (`git mv`).
- [ ] I have implemented `-ignore default` to skip `.git` and `node_modules`.
- [ ] I have implemented `fix-seq-files` with both `-orderbytime` and `-orderbyaz`.
- [ ] I have implemented the sequence fixation (pinning) logic and added it to the help text.
- [ ] I have handled Windows long paths properly by utilizing path normalization before traversing directories.
- [ ] I have added complete examples to the CLI help interface.
- [ ] I did NOT leave any `TODO` placeholders in the code.
