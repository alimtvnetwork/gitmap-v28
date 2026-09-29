# 22-verify-markdown-restructure-and-universal-help.md

Use this prompt in any CI pipeline, local test runner, or autonomous AI agent to rigorously verify that all commands support `gitmap <cmd> help` and that markdown help is restructured in the modern framed double-line box layout.

## Pre-flight Quality Gate Checklist

1. **Verify Universal Help Interception**:
   - Run `gitmap backup help` -> Verify framed double-line box banner `╔═╗` at top, clean `Subcommands:`, `Flags:`, `Examples:`, and zero raw markdown table pipes `|`.
   - Run `gitmap clone help` -> Verify modern framed double-line box header and zero side-effects.
   - Run `gitmap clean help` (dynamic catalog fallback) -> Verify framed box header, `Usage:`, `Description:`, `Examples:`, and exit code 0.
   - Run `gitmap add help` -> Verify box banner even for files without leading `#`.

2. **Verify Unit Tests & Formatting**:
   - Execute `go test ./render/...` -> `TestRenderANSIBoxBannerAndTable` must pass.
   - Execute `go test ./cmd/... -run TestRenderDynamicCommandHelp` -> must pass.
   - Execute `python .github/scripts/go-format-check.py --check-only` -> `All .go files are gofmt-clean`.

3. **Verify LLM Train & Cheatsheets**:
   - Run `gitmap train --help` -> Verify `--loop`, `--url`, and `--json` flags.
   - Run `gitmap train --loop 2` -> Verify simulation cycles complete cleanly with convergence telemetry.

4. **Verify Workstation Hygiene & Cache Purge**:
   - Run `python 03-ai-scripts/42-clean-test-and-build-caches.py` -> All temp files and Go cache cleared.

5. **CI/CD Quality Gate**:
   - Verify all workflows on `main` are 100% green via `gh run list --limit 6`.
