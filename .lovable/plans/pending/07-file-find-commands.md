# 07 File Find Commands Implementation Spec

## 1. Problem Context
The user requires a robust suite of `find` commands that search for file names inside a repository using the previously built Split DB architecture. It must support exact matches, wildcards, regex, limit flags, and content-reading extensions.

## 2. Core Requirements
### 2.1 File Find Variants (Base)
- `gitmap find "exact"`
- `gitmap find "*starts"`
- `gitmap find "ends*"`
- `gitmap find "mid*dle"`
- `gitmap find "*contains*"`
- Uses analytical cache. If not, uses repo sqlite DB.

### 2.2 Regex & Read Variants
- `gitmap find-regex "<regex>"` (Regex file name match)
- `gitmap find-read "<query>"` (Matches file names, then outputs file contents)
- `gitmap find-read-json "<query>"` (JSON payload)
- `gitmap find-regex-read "<regex>"`
- `gitmap find-regex-read-json "<regex>"`

### 2.3 Limits and Help
- `--limit <n>` flag to bound the result sets.
- `gitmap find help`, `gitmap search help`, `gitmap regex help` that organize and resuggest all related commands.

## 3. Subtasks Execution
1. **01-cli-registry.md**: Add `CmdFind`, `CmdFindRegex`, `CmdFindRead`, `CmdFindReadJson`, `CmdFindRegexRead`, `CmdFindRegexReadJson`, `CmdFindHelp`, `CmdSearchHelp`, `CmdRegexHelp` and the `--limit` flag.
2. **02-find-engine.md**: Implement `searcher/finder.go` with SQL LIKE combinations and `lazyregex` for filenames. Implement cache tables or `SearchCache` reuse.
3. **03-read-engine.md**: Implement file content extraction. Since big files (>300KB) aren't in SQLite, read them from disk if matched.
4. **04-help-routers.md**: Write the help categorization printouts.
5. **05-wire-and-test.md**: Hook up `find_entry.go` to the engine logic and pass all AST tests.
