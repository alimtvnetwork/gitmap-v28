# 26-go-125-gofmt-and-goimports-drift

## Error Summary

The CI/CD pipeline failed at `bash ../.github/scripts/go-format-check.sh` with:
`Error: The following .go files are not gofmt-clean:` across 36 files.

## Root Cause Analysis

- After bumping to Go 1.25, the Go 1.25 toolchain's `gofmt` introduced stricter doc-comment and directive spacing rules (e.g. enforcing blank `//` comment lines before directive comments like `//nolint:unused` and normalizing comment indentations on `// return nil`).
- In addition, earlier edits in Windows environments used standard string replacements without running full recursive `gofmt` and `goimports -local "github.com/alimtvnetwork/gitmap-v28/gitmap"` across subpackages.

## Solution Applied

1. Executed a recursive `gofmt -w` pass across all `.go` files in the repository to format all doc comments and code blocks according to Go 1.25 standards.
2. Executed a recursive `goimports -w -local "github.com/alimtvnetwork/gitmap-v28/gitmap"` pass across all `.go` files to reorder import blocks into stdlib, third-party, and first-party.
3. Verified locally that `gofmt -l .` and `goimports -l -local ...` return zero files.
4. Ran `.lovable/ai-fix-scripts/02-cicd-local-runner.py` verifying full passing status for both Go Vet and Full Suite Lint.

## What NOT to Repeat

- ALWAYS run `gofmt -w` and `goimports -w -local` recursively across all files after upgrading the Go toolchain.
