# Issue: Macro CMD Syntax Error & Test Failures
**Date:** 2026-08-25
**Module:** gitmap/cmd and gitmap/helptext

## Description
A top-level Go declaration syntax error occurred (cmd/macro_cmd.go:124:2: syntax error: unexpected keyword type after top level declaration) because automated code generation strictly appended strings without verifying structural newlines. A parallel issue occurred during parsing for the se command because eorderFlagsBeforeArgs was aggressively attempting to hoist flags (-p) out of nested positional command arguments (mkdir -p foo).
Also, two help text markdown files (chrome-profile-export.md and chrome-profile-import.md) lacked the required ## Examples section containing a code block, causing gitmap/helptext integration tests to fail.

## Root Cause Analysis
1. **Scaffold Templating flaw**: The orchestration script wrote structs via string concatenation (}type ...). If the target file lacked a trailing newline, it resulted in structurally invalid syntax like }type TypeFuncf0956a91 struct{}.
2. **Greedy Flag Reordering**: gitmap/cmd/releaseargs.go partitioned all flags ahead of positional commands. When passed arbitrary bash scripts for SSH Execution (e.g. mkdir -p), it incorrectly captured the -p flag into the se binary's flag parser, causing a lag provided but not defined panic.
3. **Missing Examples Constraint**: Golden tests dynamically mandate that every .md file in gitmap/helptext contains a ## Examples section followed by a fenced code block.

## Resolution
1. Parsed gitmap/cmd/*.go files using Get-Content to locate and surgically correct }type  into }\n\ntype  structurally.
2. Removed the eorderFlagsBeforeArgs() wrapper from the parseSEFlags in gitmap/cmd/sshexec.go. SSH Execution arguments are explicitly allowed to have unbounded nested flags, so strict linear parsing via s.Parse(args) is the correct semantic.
3. Appended ## Examples and fenced code blocks for chrome-profile-export and chrome-profile-import inside gitmap/helptext.
4. Tests fully re-run.

Status: Resolved
