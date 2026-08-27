# 03-ignore-and-add: Ignore management and Common Files generation

## 1. Context and Problem Statement
The user requested four new commands for gitmap:
1. gitmap ignore <pattern>: Adds a pattern to .gitignore cleanly without duplicates.
2. gitmap ignore-rm <pattern>: Rewrites git history to remove matching files, then adds the pattern to .gitignore.
3. gitmap add common-attr: Creates a common .gitattributes file.
4. gitmap add common-ignore: Creates a common .gitignore file.

## 2. Architecture & Design

### Command ignore & ignore-rm
- **Location**: gitmap/cmd/ignore
- **Logic for ignore**:
  - Open .gitignore (create if doesn't exist).
  - Check if the pattern is already present.
  - Append to the end, ensuring proper newline spacing.
- **Logic for ignore-rm**:
  - Run the same history rewrite logic developed in git-rm (e.g., git filter-branch --index-filter 'git rm --cached --ignore-unmatch -r <pattern>' ...).
  - Then call the ignore logic to append the pattern.

### Command dd
- **Location**: gitmap/cmd/add
- **Logic**:
  - Switch on the argument (common-attr or common-ignore).
  - Write embedded boilerplate/template files for .gitattributes and .gitignore.
  - For .gitattributes: include common text auto-crlf settings, LFS configs, etc.
  - For .gitignore: include common OS files (.DS_Store, Thumbs.db), IDE files (.vscode/, .idea/), etc.

## 3. Subtasks (to be created in .lovable/plans/subtasks/03-ignore-and-add/)

1. **Subtask 1: Scaffold ignore & ignore-rm Commands**
   - Register top-level commands.
   - Implement ignore logic (file appending).
   - Implement ignore-rm logic (history rewrite + file appending).
2. **Subtask 2: Scaffold `add` Command (COMPLETED)**
   - Register top-level command.
   - Implement common-attr writing logic.
   - Implement common-ignore writing logic.
