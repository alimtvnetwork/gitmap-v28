# RCA: clone-sync.md Missing Fenced Code Blocks (TestEveryHelpFileHasExamples)

## Error Summary

The `github.com/alimtvnetwork/gitmap-v28/gitmap/helptext` test suite failed due to `TestEveryHelpFileHasExamples`. The error reported:
`help files missing ## Examples section (1): - clone-sync.md`.

## Root Cause Analysis

1. We recently introduced the `clone-sync.md` help file.
2. The help file contained an `## Examples` section, but used 4-space indented code blocks instead of fenced code blocks (````bash`).
3. The `TestEveryHelpFileHasExamples` strictly checks for `\n## Examples` AND the presence of the ```` token (`strings.Contains(md[idx:], "```")`).
4. Because fenced code blocks were absent, the test evaluated the condition as false and failed the CI pipeline.

## Solution Applied

Modified `gitmap/helptext/clone-sync.md` to use markdown fenced code blocks (````bash ... ````) instead of 4-space indentations for the commands in the `## Examples` section. 

## What NOT to Repeat

- Never use 4-space indentations for code blocks in `helptext/` markdown files. Always use fenced code blocks (```` ``` ````) to satisfy the `TestEveryHelpFileHasExamples` runnable-examples CI gate.
