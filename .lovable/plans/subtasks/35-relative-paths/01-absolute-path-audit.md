# Subtask 01 - Absolute Path Audit & Inventory

## Parent Specification
[35-relative-paths-audit.md](.lovable/plans/pending/35-relative-paths-audit.md)

## Acceptance Criteria & Requirements
- Scan repository files for absolute filesystem paths (`/absolute/path/to/...`, `C:\Users\...`, `/home/...`).
- Verify that markdown links, citations, and subtask paths use strictly relative Git root paths.
- Enforce 0 `file:///` URIs across all documentation.
