import os
import re
import json

repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
gitmap_dir = os.path.join(repo_root, "gitmap")
violations = []


def check_type_suffix(line_str, rel_path, idx, violations):
    m = re.match(r"^type\s+([A-Z][a-zA-Z0-9_]*)\s+(string|int|int[0-9]+|uint|uint[0-9]+|byte)\b", line_str)
    if not m:
        return
    type_name = m.group(1)
    if type_name.endswith(("Type", "ID", "Code", "Flags")):
        return
    violations.append({
        "category": "enum-suffix",
        "file": rel_path,
        "line": idx + 1,
        "desc": f"Type alias `{type_name}` acting as enum lacks `Type` suffix.",
        "fix": f"Rename `{type_name}` to `{type_name}Type`."
    })


def audit_file(filepath):
    rel_path = os.path.relpath(filepath, repo_root).replace("\\", "/")
    with open(filepath, "r", encoding="utf-8", errors="ignore") as f:
        lines = f.readlines()
    
    # 1. Enum Definitions lacking Type suffix
    for idx, line in enumerate(lines):
        check_type_suffix(line.strip(), rel_path, idx, violations)
    
    # 2. Function lengths (> 15 lines) in non-test files
    func_start = None
    func_name = ""
    brace_depth = 0
    
    for idx, line in enumerate(lines):
        m = re.match(r"^func\s+(?:\([^)]+\)\s+)?([A-Za-z0-9_]+)\s*\(", line)
        if m and func_start is None:
            func_name = m.group(1)
            func_start = idx + 1
            brace_depth = line.count("{") - line.count("}")
            is_single_line = bool("{" in line and brace_depth == 0)
            func_start = None if is_single_line else func_start
            continue
        
        if func_start is None:
            continue

        brace_depth += line.count("{") - line.count("}")
        if brace_depth > 0:
            continue

        func_len = (idx + 1) - func_start + 1
        if func_len > 15 and not rel_path.endswith("_test.go"):
            violations.append({
                "category": "long-func",
                "file": rel_path,
                "line": func_start,
                "desc": f"Function `{func_name}` exceeds 15 lines ({func_len} lines).",
                "fix": "Extract helper functions, table-driven dispatch, or guard clauses."
            })
        func_start = None

    # 3. Inverted booleans & Negated logic
    for idx, line in enumerate(lines):
        code_part = line.split("//")[0]
        for b_match in re.finditer(r"!(is[A-Z]\w*|has[A-Z]\w*|can[A-Z]\w*|should[A-Z]\w*|was[A-Z]\w*|did[A-Z]\w*|will[A-Z]\w*|must[A-Z]\w*|ok\b|found\b|valid\b|exists\b)", code_part):
            violations.append({
                "category": "inverted-bool",
                "file": rel_path,
                "line": idx + 1,
                "desc": f"Inverted boolean logic: `!{b_match.group(1)}`.",
                "fix": "Extract into positive boolean check or use explicit `== false`."
            })

    # 4. Nested if statements (depth >= 2)
    nesting_stack = []
    for idx, line in enumerate(lines):
        code_part = line.split("//")[0].strip()
        indent_level = len(line) - len(line.lstrip('\t'))
        if not (code_part.startswith("if ") and "{" in code_part):
            while nesting_stack and indent_level < nesting_stack[-1][1]:
                nesting_stack.pop()
            continue

        nesting_stack.append((idx + 1, indent_level))
        if len(nesting_stack) >= 3 and not rel_path.endswith("_test.go"):
            violations.append({
                "category": "nested-if",
                "file": rel_path,
                "line": idx + 1,
                "desc": f"Nested `if` statement detected at indentation depth {len(nesting_stack)}.",
                "fix": "Flatten with early returns or guard clauses."
            })
        while nesting_stack and indent_level < nesting_stack[-1][1]:
            nesting_stack.pop()


for root, dirs, files in os.walk(gitmap_dir):
    for f in files:
        if f.endswith(".go") and not f.endswith(".pb.go"):
            audit_file(os.path.join(root, f))

with open(".lovable/scratch/cg_violations.json", "w", encoding="utf-8") as f:
    json.dump(violations, f, indent=2)

print(f"Audit complete! Found {len(violations)} total violations.")
