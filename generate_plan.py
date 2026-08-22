import os
import csv
from collections import defaultdict

def ensure_dir(d):
    if not os.path.exists(d):
        os.makedirs(d)

ensure_dir(".lovable/plans/pending")
ensure_dir(".lovable/plans/subtasks/01-coding-guideline-fixes")
ensure_dir(".lovable/plans/subtasks/02-error-management-fixes")

# Read Go issues
go_issues = []
with open("audit_results_go.csv", "r", encoding="utf-8") as f:
    for line in f:
        parts = line.strip().split('|')
        if len(parts) >= 4:
            go_issues.append({
                'file': parts[0],
                'line': int(parts[1]),
                'type': parts[2],
                'content': parts[3]
            })

# Categorize for Error Management (02)
em_categories = {
    'swallowed error': [],
    'panic': [],
    'raw error return': [],
    'fmt.Errorf': []
}

cg_categories = {
    'bare bool return': []
}

for issue in go_issues:
    t = issue['type']
    if t in em_categories:
        em_categories[t].append(issue)
    elif t in cg_categories:
        cg_categories[t].append(issue)

# Read TS issues
ts_issues = []
with open("audit_results_ts.csv", "r", encoding="utf-8") as f:
    for line in f:
        parts = line.strip().split('|')
        if len(parts) >= 4:
            ts_issues.append({
                'file': parts[0],
                'line': int(parts[1]),
                'type': parts[2],
                'content': parts[3]
            })

for issue in ts_issues:
    if issue['type'] == 'swallowed error (catch {})':
        em_categories['swallowed error'].append(issue)
    else:
        if issue['type'] not in cg_categories:
            cg_categories[issue['type']] = []
        cg_categories[issue['type']].append(issue)

# Generate 02 Master Plan
plan_02_content = """# Error Management Migration Plan

This plan documents all required fixes to migrate the codebase to strict error management using `AppError`.

## Rules to Enforce:
1. **Never swallow errors.** Every catch/ignore logs the operation name and key inputs, then rethrows or returns a typed error.
2. **Wrap, do not lose.** Wrap the original error with an operation label and context (`apperror.Wrap(err, 'op', ctx)`).
3. **Typed errors only.** No bare `panic('msg')` or raw `fmt.Errorf` without wrapping.
4. Every variable needs to be captured in an error log/path.
5. All functions returning errors must return `*apperror.AppError`, not raw `error`.

## Subtasks Overview:
"""

subtask_idx = 1
for cat_name, items in em_categories.items():
    if not items: continue
    
    # Split into chunks of 150
    chunk_size = 150
    chunks = [items[i:i + chunk_size] for i in range(0, len(items), chunk_size)]
    
    for i, chunk in enumerate(chunks):
        slug = cat_name.replace(' ', '-').replace('.', '')
        subtask_filename = f"{subtask_idx:02d}-{slug}-part{i+1}.md"
        plan_02_content += f"- [ ] `{subtask_filename}`: Fix {len(chunk)} `{cat_name}` issues.\n"
        
        # Write subtask file
        subtask_path = f".lovable/plans/subtasks/02-error-management-fixes/{subtask_filename}"
        with open(subtask_path, "w", encoding="utf-8") as sf:
            sf.write(f"# Fix {cat_name} (Part {i+1})\n\n")
            sf.write(f"Total items: {len(chunk)}\n\n")
            sf.write("## Files to Modify\n\n")
            for item in chunk:
                sf.write(f"- `{item['file']}:{item['line']}`: `{item['content']}`\n")
        
        subtask_idx += 1

with open(".lovable/plans/pending/02-error-management-fixes.md", "w", encoding="utf-8") as f:
    f.write(plan_02_content)

# Generate 01 Master Plan
plan_01_content = """# Coding Guidelines Audit Plan

This plan documents all required fixes to align the codebase with standard coding guidelines.

## Rules to Enforce:
- Nested ifs
- Boolean variables not starting with `is`, `has`, `can`, `should`, etc.
- Inverted booleans like `isNot*` instead of `is*`
- Magic strings and numbers
- Missing Enum/Type definitions
- Monolithic functions exceeding 15 lines
- `return` statements without preceding blank lines.
- Golang functions returning bare bools instead of a wrapped Result object.

## Subtasks Overview:
"""

subtask_idx = 1
for cat_name, items in cg_categories.items():
    if not items: continue
    
    chunk_size = 150
    chunks = [items[i:i + chunk_size] for i in range(0, len(items), chunk_size)]
    
    for i, chunk in enumerate(chunks):
        slug = cat_name.replace(' ', '-').replace('.', '')
        subtask_filename = f"{subtask_idx:02d}-{slug}-part{i+1}.md"
        plan_01_content += f"- [ ] `{subtask_filename}`: Fix {len(chunk)} `{cat_name}` issues.\n"
        
        subtask_path = f".lovable/plans/subtasks/01-coding-guideline-fixes/{subtask_filename}"
        with open(subtask_path, "w", encoding="utf-8") as sf:
            sf.write(f"# Fix {cat_name} (Part {i+1})\n\n")
            sf.write(f"Total items: {len(chunk)}\n\n")
            sf.write("## Files to Modify\n\n")
            for item in chunk:
                sf.write(f"- `{item['file']}:{item['line']}`: `{item['content']}`\n")
        
        subtask_idx += 1

with open(".lovable/plans/pending/01-coding-guideline-fixes.md", "w", encoding="utf-8") as f:
    f.write(plan_01_content)

print("Generated plans successfully.")
