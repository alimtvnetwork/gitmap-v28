# Plan 88: Pipeline Error Multi-Line Extraction Intelligence & Declarative JSON Format Profiles

Spec Reference: [02-spec/21-app/139-pe-custom-format-and-error-capture.md](../../../02-spec/21-app/139-pe-custom-format-and-error-capture.md)

## Architectural Context & Blast Radius
GitMap's pipeline error extractor (`gitmap pe`) previously isolated single failure markers and truncated single-line warnings. In Rust/Tauri/Node workflows, compiler warnings span multiple lines (`warning:` followed by source code pointers `-->`, line indicators `|`, code lines, and `= note:` annotations), and bundler failures (`failed to bundle project ... does not exist`) were previously missed because only `Process completed with exit code 1.` was logged.

This plan upgraded the default log parser to preserve multi-line context blocks by default, and introduced a declarative JSON format engine with registry storage and CLI management suite.

## Deliverables & Outcomes
- **Visual Assets & Spec 139 Authoring:**
  - Ingested screenshots into `assets/screenshots/pe-pipeline-error-01.png`, `02.png`, `03.png`.
  - Authored canonical Spec 139 in `02-spec/21-app/139-pe-custom-format-and-error-capture.md`.
  - Registered Spec 139 in `02-spec/21-app/01-index.md`.
- **Default Multi-Line Warning & Bundler Extraction:**
  - Expanded `failureMarkers` with bundler, binary copy, missing path, and compiler error patterns.
  - Implemented multi-line warning block capture in `processLogLine` and `isWarningContinuationLine`.
  - Updated `isToolProgressNoise` to strip Vite, Browserslist, and Tauri banner noise.
  - Prioritized bundler failures over generic exit codes in `isStrongerSummary`.
- **Declarative JSON Format Profile Engine & CLI Commands:**
  - Implemented `PEFormatProfile` model with `strip_prefixes`, `strip_contains`, `skip_empty`, `warning_markers`, `error_markers`, and `capture_until_markers`.
  - Implemented `BuiltinTauriProfile()` and `DefaultPEFormatProfile()`.
  - Implemented format storage in `<dataDir>/pipeline/formats/<name>.json`.
  - Supported full CLI suite:
    - `gitmap pe -f <format.json|alias>`
    - `gitmap pe -f <alias> -test <filepath>`
    - `gitmap pe -f <alias> -test-commit <sha> [-repo <path>]`
    - `gitmap pe add-format <file.json> [alias]`
    - `gitmap pe rm-format <name|alias>`
    - `gitmap pe add-all <folder-path>`
    - `gitmap pe list-formats`
    - `gitmap pe preview-format <name|file>`
  - Updated `gitmap pe help` and `gitmap pe --help` with comprehensive documentation and examples.
- **Verification & Quality Gates:**
  - Added unit tests in `cli/cmdpipeline/pipeline_format_test.go` and `cli/cmdpipeline/pipeline_error_extract_test.go`.
  - Verified 0 nested ifs via `check-nested-ifs.py`.
  - Verified boolean and enum compliance via `check-enum-and-boolean.py`.
  - Verified error management via `check-error-management.py`.
  - Formatted Go code via `26-go-code-formatter.py`.
  - Executed end-to-end CLI tests with real compiled `bin/gitmap.exe`.
