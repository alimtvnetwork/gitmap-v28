# Subtask 94.01: Generator Upgrades & Tag-Driven Code Generation

## Goal
Upgrade `03-ai-scripts/30-db-struct-enum-generator.py` to support explicit `db:` struct tags, eliminate hardcoded struct skips, and generate typed mutation helpers (`Insert`, `Update`, `DeleteById`) on `*DbRepo`.

## Files Impacted
- `03-ai-scripts/30-db-struct-enum-generator.py`
- `gitmap/pipelinedb/enums/`
- `gitmap/pipelinedb/`

## Acceptance Criteria
1. Struct parser parses `db:"ColumnName"` tag from struct fields; defaults to field name if tag is absent.
2. Removes hardcoded `if s_name == "PipelineSplitDb"` hack by checking for explicit entity markers or skipping connection wrapper structs.
3. Extends `<Model>DbRepo` generation to include typed mutation methods (`Insert`, `Update`, `DeleteById`) utilizing `dbengine.DbWrapper`.
4. Regenerates `gitmap/pipelinedb/` cleanly with no regression.
5. All generated functions $\le 15$ lines, zero nested ifs, affirmative booleans only.
