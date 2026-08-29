# 02-git-rm-and-folder: Folder Export and Git History Cleaning

## 1. Context and Problem Statement
The user requested two new top-level commands for gitmap:
1. gitmap folder <dir> <output_file> [-exclude <pattern> 0|1]
   - Extracts a relative folder structure into an output file.
   - Output format depends on the file extension (.txt, .md, .json, .yaml).
   - Default output file if none provided (e.g. $path/nothing) is iles.txt.
2. gitmap git-rm <input>
   - Reads the specified input (a single file path, folder, CSV list, or a text/json file containing paths like my-files.json).
   - Removes these files entirely from the Git history.
   - **Crucial Requirement**: It must back up the removed files to the globally installed .gitmap location (e.g., ~/.gitmap/backups/git-rm/), NOT the local repository's .gitmap folder.

## 2. Architecture & Design

### Command older
- **Location**: gitmap/cmd/folder
- **Logic**: Walk the specified directory recursively using standard Go ilepath.WalkDir or s.WalkDir.
- **Filtering**: Apply -exclude globs. The format seems to be -exclude <pattern> <flag>.
- **Output Formats**:
  - 	xt / md: Tree-like text structure.
  - json / yaml: Structured list or hierarchy.

### Command git-rm
- **Location**: gitmap/cmd/gitrm
- **Input Parsing**: Detect if input is a JSON file (parse paths), a TXT file (read lines), a CSV string, or a direct path/folder.
- **Backup**:
  - Resolve global .gitmap location (e.g., ~/.gitmap/backups/git-rm/<repo-name>-<timestamp>).
  - Extract the files as they currently exist at HEAD (or across history if we want full backup, but typically HEAD or latest found is sufficient).
- **History Rewrite**:
  - Use git filter-branch --force --index-filter 'git rm --cached --ignore-unmatch <paths>' --prune-empty --tag-name-filter cat -- --all or similar. Alternatively, generate a fast-export stream and filter it. Given Go and cross-platform needs, git filter-branch might be deprecated, but it's universally available.

## 3. Subtasks (to be created in .lovable/plans/subtasks/02-git-rm-and-folder/)

1. **Subtask 1: Scaffold older Command**
   - Register top-level command older.
   - Implement directory walking and -exclude flag parsing.
   - Implement .txt and .md tree export.
2. **Subtask 2: Implement .json and .yaml exports for older**
   - Implement JSON and YAML marshaling for the tree/list.
3. **Subtask 3: Scaffold git-rm Command & Input Parsing**
   - Register top-level command git-rm.
   - Implement input parser (JSON, TXT, CSV, direct paths).
4. **Subtask 4: Implement Backup & Git History Rewrite**
   - Implement file backup to global ~/.gitmap/backups/git-rm/.
   - Implement history rewrite mechanism (e.g., using git filter-branch or git rev-list + manual reconstruct if small, but ilter-branch is safest for generic git repos).

## 4. Coding Guidelines Checklist
- All is/has/can/should booleans used.
- PascalCase acronyms (Json, Yaml, Id).
- Max 15 lines per function, max 200 lines per file.
- Strict error management using pperror.Wrap.
- All magic strings extracted to constants/.
- Register commands in cmd_constants_test.go or skip.
