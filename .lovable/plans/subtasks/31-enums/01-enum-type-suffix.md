# Subtask 01 - Enum `*Type` Suffix Refactoring

## Parent Specification
[31-constants-and-enums-audit.md](.lovable/plans/pending/31-constants-and-enums-audit.md)

## Acceptance Criteria & Requirements
- Update `gitmap/cmd/find_files.go` to define `type MatchKindType int` and `type MatchKind = MatchKindType`.
- Update `gitmap/cliexit/report.go` to define `type OutputModeType int` and `type OutputMode = OutputModeType`.
- Update `gitmap/archive/create.go` to define `type CompressionModeType string` and `type CompressionMode = CompressionModeType`.
- Update `gitmap/cmd/folder/folder.go` to define `type OutputFormatType string` and `type OutputFormat = OutputFormatType`.
- Ensure all Go, TypeScript, PHP, and Python enums end with `Type`.
