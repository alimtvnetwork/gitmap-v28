import os
import re
import json

repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))
gitmap_dir = os.path.join(repo_root, "gitmap")
violations = []

def audit_file(filepath):
    rel_path = os.path.relpath(filepath, repo_root).replace("\\", "/")
    with open(filepath, "r", encoding="utf-8", errors="ignore") as f:
        lines = f.readlines()
    
    # 1. Enum Definitions lacking Type suffix
    for idx, line in enumerate(lines):
        line_str = line.strip()
        m = re.match(r"^type\s+([A-Z][a-zA-Z0-9_]*)\s+(string|int|int[0-9]+|uint|uint[0-9]+|byte)\b", line_str)
        if m:
            type_name = m.group(1)
            if not (type_name.endswith("Type") or type_name.endswith("ID") or type_name.endswith("Code") or type_name.endswith("Flags")):
                violations.append({
                    "category": "enum-suffix",
                    "file": rel_path,
                    "line": idx + 1,
                    "desc": f"Type alias `{type_name}` acting as enum lacks `Type` suffix.",
                    "fix": f"Rename `{type_name}` to `{type_name}Type`."
                })
    
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
            if "{" in line and brace_depth == 0:
                func_start = None
            continue
        
        if func_start is not None:
            brace_depth += line.count("{") - line.count("}")
            if brace_depth <= 0:
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
        # Inverted booleans like !isX, !hasX, !canX, !shouldX, !ok, !found, !valid, !exists
        for b_match in re.finditer(r"!(is[A-Z]\w*|has[A-Z]\w*|can[A-Z]\w*|should[A-Z]\w*|was[A-Z]\w*|did[A-Z]\w*|will[A-Z]\w*|must[A-Z]\w*|ok\b|found\b|valid\b|exists\b)", code_part):
            violations.append({
                "category": "inverted-bool",
                "file": rel_path,
                "line": idx + 1,
                "desc": f"Inverted boolean logic: `!{b_match.group(1)}`.",
                "fix": f"Extract into positive boolean check or use explicit `== false`."
            })

    # 4. Nested if statements (depth >= 2)
    current_nesting = 0
    nesting_stack = []
    for idx, line in enumerate(lines):
        code_part = line.split("//")[0].strip()
        indent_level = (len(line) - len(line.lstrip('\t')))
        if code_part.startswith("if ") and "{" in code_part:
            nesting_stack.append((idx + 1, indent_level))
            if len(nesting_stack) >= 3 and not rel_path.endswith("_test.go"):
                violations.append({
                    "category": "nested-if",
                    "file": rel_path,
                    "line": idx + 1,
                    "desc": f"Nested `if` statement detected at indentation depth {len(nesting_stack)}.",
                    "fix": "Flatten with early returns or guard clauses."
                })
        # approximate block exit
        while nesting_stack and indent_level < nesting_stack[-1][1]:
            nesting_stack.pop()

    # 5. Swallowed errors (_ = err, or unchecked error in assignment without handling)
    for idx, line in enumerate(lines):
        code_part = line.split("//")[0].strip()
        if re.search(r"\b_\s*=\s*(?:[a-zA-Z0-9_]+\.)?([A-Za-z0-9_]+Error|err|error)\b", code_part):
            violations.append({
                "category": "swallowed-error",
                "file": rel_path,
                "line": idx + 1,
                "desc": f"Swallowed error variable: `{code_part}`.",
                "fix": "Handle or wrap error using `apperror.Wrap`."
            })

for root, _, files in os.walk(gitmap_dir):
    for file in files:
        if file.endswith(".go"):
            audit_file(os.path.join(root, file))

print(f"Total violations found: {len(violations)}")
with open("audit_findings.json", "w", encoding="utf-8") as f:
    json.dump(violations, f, indent=2)
