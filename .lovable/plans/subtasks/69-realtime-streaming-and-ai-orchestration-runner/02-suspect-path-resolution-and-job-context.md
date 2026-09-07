# Subtask 02: Suspect Path Resolution, CWD Propagation & Multi-Tier File Extraction

## Objective
Enhance `JobResult` to propagate `cwd` and `env_overrides`, and implement a robust multi-tier suspect file extractor that correctly parses Windows drive letters, TypeScript formats, Python tracebacks, and falls back to git delta and gate specs.

## Requirements
1. **Job Execution Context Propagation**:
   - Add `cwd: str | None = None` and `env_overrides: dict[str, str] | None = None` to `JobResult`.
   - Propagate `cwd` and `env` from `submit_job_futures` and `run_job` into `JobResult`.
2. **Multi-Pattern Suspect File Extractor**:
   - Pattern 1: Standard colon syntax (`path/to/file.ext:line:col`) including Windows drive letters (`D:\...`).
   - Pattern 2: TypeScript / MSBuild parentheses syntax (`src/App.tsx(45,12)`).
   - Pattern 3: Python traceback syntax (`File "path/to/file.py", line 42`).
   - Pattern 4: Go panic and race stack traces (`\s+path/to/file.go:line`).
3. **Path Normalization & CWD Resolution**:
   - Resolve relative paths against `cwd` if set (e.g., prefixing `gitmap/` if gate ran inside `cwd="gitmap"` and file exists there).
   - Ignore internal vendor and system paths (`node_modules/`, `vendor/`, `.git/`, `.tmp/`, `go/pkg/mod/`, `AppData/`, `site-packages/`).
4. **Semantic Git Delta Fallback**:
   - If error output contains no explicit file paths, extract files by intersecting `repo_delta` with the gate's `relevant_patterns`, falling back to `tool_scripts` and `configs`.
5. **Coding Guidelines**:
   - All functions $\le 15$ lines.
   - Blank line before every return statement.
   - Affirmative booleans (`is_*`, `has_*`).
   - Zero swallowed exceptions.

## Target Files
- `03-ai-scripts/06-cicd-local-runner.py`
- `.lovable/ai-fix-scripts/06-cicd-local-runner.py`

## Acceptance Criteria
- [x] `JobResult` contains `cwd` and `env_overrides`.
- [x] Suspect file extractor parses Windows drive letters (`D:\...`) without truncation.
- [x] TypeScript parens `(line,col)` and Python tracebacks are converted to normalized `path:line` format.
- [x] Relative paths from gates running in `cwd="gitmap"` are resolved to repository-relative paths.
- [x] Non-syntax failures fall back to relevant git delta files or gate tool scripts.
