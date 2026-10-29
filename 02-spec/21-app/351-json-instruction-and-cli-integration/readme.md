# 351 - JSON Instruction and CLI Integration Specification

This specification suite governs the end-to-end design and implementation of structured JSON instruction ingestion, link expansion, CLI output formats across GitMap and `go-email-sender`, prompt formatter condensation (< 5,000 characters) across Letterly and Cursor suites, and multi-repository fleet synchronization.

## Documents
- `01-architecture-spec.md`: Architectural specification for the canonical JSON instruction schema, validation rules, link expansion pipeline, and multi-binary CLI contract.
- `02-component-spec.md`: Component specifications and implementation contracts for GitMap CLI (`gitmap space send-email --json`, `gitmap email expand`), `go-email-sender` parity, and unit test suites.
- `03-prompt-condensation-and-sync-spec.md`: Character budget enforcement (< 5,000 chars) across 12 Letterly and Cursor prompts, double-declaration elimination, skill link URI standards, and multi-repository synchronization protocol using `03-ai-scripts/38-sync-prompts-skills-scripts.py --force`.
