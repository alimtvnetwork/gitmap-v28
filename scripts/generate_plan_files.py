import json
import os

with open("audit_findings.json", "r", encoding="utf-8") as f:
    findings = json.load(f)

enum_items = [x for x in findings if x["category"] == "enum-suffix"]
swallowed_items = [x for x in findings if x["category"] == "swallowed-error"]
inverted_items = [x for x in findings if x["category"] == "inverted-bool"]
nested_items = [x for x in findings if x["category"] == "nested-if"]
long_items = [x for x in findings if x["category"] == "long-func"]

selected = []
# 1. All enum suffixes (12)
selected.extend(enum_items)
# 2. All swallowed errors (1)
selected.extend(swallowed_items)
# 3. 87 inverted booleans
selected.extend(inverted_items[:87])
# 4. 25 nested ifs
selected.extend(nested_items[:25])
# 5. 25 long funcs
selected.extend(long_items[:25])

assert len(selected) == 150, f"Expected 150 items, got {len(selected)}"

os.makedirs(".lovable/plans/pending", exist_ok=True)
os.makedirs(".lovable/plans/subtasks/01-coding-guideline-fixes", exist_ok=True)

# 1. Generate .lovable/plans/pending/01-coding-guideline-fixes.md
plan_lines = [
    "# Plan: Coding Guideline Audit & Enforcement (v4)",
    "",
    "- slug: 01-coding-guideline-fixes",
    "- status: pending",
    "- steps_count: 150",
    "- created: 2026-08-29",
    "",
    "## Context",
    "Comprehensive, deep multi-stage audit of the entire codebase for coding guideline violations, boolean anti-patterns, missing enum suffixes, cyclomatic nesting, and error-handling flaws across Go backend packages. Structured into exactly 150 granular steps.",
    "",
    "## Fallout & Blast Radius Analysis",
    "- **Enum Suffix Renames**: Modifying type definitions requires updating all call sites across CLI commands, tests, and struct field definitions. Blast radius is contained within Go internal packages.",
    "- **Boolean Inversions**: Changing `!isX` to explicit `isX == false` or `isMissing` preserves runtime semantics without breaking API contracts or downstream scripts.",
    "- **Nested If Flattening & Function Extraction**: Uses guard clauses and single-responsibility helper extractions to strictly comply with the <= 15-line function cap.",
    "- **CI/CD Guard**: No CI/CD workflows, GitHub Actions, or validation rules will be bypassed.",
    "",
    "## Enqueued Granular Tasks (Steps 1 to 150)",
    ""
]

for idx, item in enumerate(selected, 1):
    plan_lines.append(f"{idx}. **{item['category']}**: `{item['file']}:{item['line']}` - {item['desc']} **Fix**: {item['fix']}")

with open(".lovable/plans/pending/01-coding-guideline-fixes.md", "w", encoding="utf-8") as f:
    f.write("\n".join(plan_lines) + "\n")

# 2. Generate subtasks for 3 concurrent subagents
# Subtask 1: Steps 1-50 (Enums, Swallowed, Inverted 1-37)
# Subtask 2: Steps 51-100 (Inverted 38-87)
# Subtask 3: Steps 101-150 (Nested Ifs & Long Functions)

subtasks_def = [
    ("01-enums-and-inverted-bools-part1.md", 0, 50, "Subagent 1: Enum Suffixes, Swallowed Errors & Inverted Booleans (Part 1)"),
    ("02-inverted-bools-part2.md", 50, 100, "Subagent 2: Inverted Booleans & Negation Refactoring (Part 2)"),
    ("03-nested-ifs-and-function-caps.md", 100, 150, "Subagent 3: Nested If Flattening & Function Line Cap Extractions")
]

for fname, start_idx, end_idx, title in subtasks_def:
    chunk = selected[start_idx:end_idx]
    st_lines = [
        f"# Subtask: {title}",
        "",
        f"- parent_plan: 01-coding-guideline-fixes.md",
        f"- steps: {start_idx + 1} to {end_idx}",
        f"- status: pending",
        "",
        "## Execution Instructions",
        "- Strictly enforce the <= 15 lines rule for any modified or extracted functions.",
        "- Wrap errors with `apperror.Wrap` and domain context.",
        "- Never use negative booleans or `!isX` inverted logic.",
        "- Place exactly one blank line before `return` statements and after closing `}` braces.",
        "",
        "## Granular Steps",
        ""
    ]
    for idx, item in enumerate(chunk, start_idx + 1):
        st_lines.append(f"{idx}. **{item['category']}**: `{item['file']}:{item['line']}` - {item['desc']}\n   - **Action**: {item['fix']}")

    with open(f".lovable/plans/subtasks/01-coding-guideline-fixes/{fname}", "w", encoding="utf-8") as f:
        f.write("\n".join(st_lines) + "\n")

print("Generated 01-coding-guideline-fixes.md and 3 subtasks successfully.")
