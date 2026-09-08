# Subtask 03: Interactive Helper Commands: Find, Search, and Replace

**Parent Plan:** [83-interactive-macro-builder-pwd-ls-search.md](../../completed/83-interactive-macro-builder-pwd-ls-search.md)  
**Status:** complete  
**Target:** `gitmap/cmd/macro_add_helpers.go`

---

## Objectives

1. Implement `find <pattern>` helper command:
   - Walk current directory tree up to max depth (default 5).
   - Match filenames with glob pattern.
   - Render matching paths relative to PWD.
2. Implement `search <query>` helper command:
   - Search file contents matching query string.
   - Limit to text/source files (ignore binary/vendor/git).
   - Print matching file paths and line numbers with snippets.
3. Implement `replace <old> <new> [target-file-or-glob]` helper command:
   - Safe in-place replacement with backup or preview.
   - Output count of replaced instances.
4. Integrate command parsing into `promptInteractiveMacroSteps` without colliding with regular shell commands.
