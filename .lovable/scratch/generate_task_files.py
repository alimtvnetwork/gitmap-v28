import os

def read_lines(file, limit):
    lines = []
    if os.path.exists(file):
        with open(file, 'r', encoding='utf-8') as f:
            for i, line in enumerate(f):
                if i >= limit:
                    break
                lines.append(line.strip())
    return lines

inverted = read_lines(r"d:\wp-work\riseup-asia\gitmap\.lovable\scratch\inverted_bools.txt", 50)
nested = read_lines(r"d:\wp-work\riseup-asia\gitmap\.lovable\scratch\nested_ifs.txt", 50)
monolithic = read_lines(r"d:\wp-work\riseup-asia\gitmap\.lovable\scratch\monolithic.txt", 50)
bare_bools = read_lines(r"d:\wp-work\riseup-asia\gitmap\.lovable\scratch\bare_bools.txt", 50)

os.makedirs(r"d:\wp-work\riseup-asia\gitmap\.lovable\plans\subtasks\01-coding-guideline-fixes", exist_ok=True)
os.makedirs(r"d:\wp-work\riseup-asia\gitmap\.lovable\plans\pending", exist_ok=True)

# Write 01-inverted-booleans.md
with open(r"d:\wp-work\riseup-asia\gitmap\.lovable\plans\subtasks\01-coding-guideline-fixes\01-inverted-booleans.md", "w") as f:
    f.write("# Task: Fix Inverted Booleans\n\n")
    f.write("According to the Coding Guidelines: 'Always use explicit boolean state variables (e.g., `isFail`) instead of inverting positive ones (`!isSuccess`).'\n\n")
    f.write("## Files to fix\n")
    for line in inverted:
        f.write(f"- {line}\n")

# Write 02-nested-ifs.md
with open(r"d:\wp-work\riseup-asia\gitmap\.lovable\plans\subtasks\01-coding-guideline-fixes\02-nested-ifs.md", "w") as f:
    f.write("# Task: Flatten Nested Ifs\n\n")
    f.write("According to the Coding Guidelines: 'No nested `if` branches; flatten logic using early returns.'\n\n")
    f.write("## Files to fix\n")
    for line in nested:
        f.write(f"- {line}\n")

# Write 03-monolithic-functions.md
with open(r"d:\wp-work\riseup-asia\gitmap\.lovable\plans\subtasks\01-coding-guideline-fixes\03-monolithic-functions.md", "w") as f:
    f.write("# Task: Refactor Monolithic Functions\n\n")
    f.write("According to the Coding Guidelines: Monolithic functions exceeding 15 lines should be refactored.\n\n")
    f.write("## Files to fix\n")
    for line in monolithic:
        f.write(f"- {line}\n")

# Write 04-bare-bools.md
with open(r"d:\wp-work\riseup-asia\gitmap\.lovable\plans\subtasks\01-coding-guideline-fixes\04-bare-bools.md", "w") as f:
    f.write("# Task: Refactor Bare Bools\n\n")
    f.write("According to the Coding Guidelines: Golang functions returning bare bools instead of a wrapped Result object should be refactored.\n\n")
    f.write("## Files to fix\n")
    for line in bare_bools:
        f.write(f"- {line}\n")

# Write master plan
with open(r"d:\wp-work\riseup-asia\gitmap\.lovable\plans\pending\01-coding-guideline-fixes.md", "w") as f:
    f.write("# Master Plan: Coding Guideline Fixes\n\n")
    f.write("This master plan details the codebase-wide audit for coding guideline violations across `gitmap` and `src` directories.\n\n")
    f.write("## Summary of Audit Findings\n")
    f.write("- **Inverted Booleans:** Found and queued 50 instances of inverted booleans (`isNot* := !is*`).\n")
    f.write("- **Nested Ifs:** Found and queued 50 instances of nested `if` statements requiring flattening.\n")
    f.write("- **Monolithic Functions:** Identified 2,200+ monolithic functions; queued 50 for refactoring.\n")
    f.write("- **Bare Bools:** Found 800+ instances of returning bare `true`/`false`; queued 50 for wrapping in `Result` objects.\n")
    f.write("\n## Subtasks\n")
    f.write("- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/01-inverted-booleans.md`\n")
    f.write("- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/02-nested-ifs.md`\n")
    f.write("- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/03-monolithic-functions.md`\n")
    f.write("- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/04-bare-bools.md`\n")
    f.write("\nTotal queued exact file/line changes: ~200.\n")
