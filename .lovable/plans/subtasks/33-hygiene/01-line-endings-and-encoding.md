# Subtask 01 - Line Endings & Encoding Standards

## Parent Specification
[33-code-hygiene-and-file-standards-audit.md](.lovable/plans/pending/33-code-hygiene-and-file-standards-audit.md)

## Acceptance Criteria & Requirements
- Enforce strict Unix LF (`\n`, `0x0A`) across all source and markdown files. Zero CRLF.
- Enforce UTF-8 without BOM encoding across all files. Zero `\xef\xbb\xbf` headers.
- Ensure every file terminates with exactly one newline at EOF.
