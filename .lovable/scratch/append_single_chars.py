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

single_chars = read_lines(r"d:\wp-work\riseup-asia\gitmap\.lovable\scratch\single_chars.txt", 50)

with open(r"d:\wp-work\riseup-asia\gitmap\.lovable\plans\subtasks\01-coding-guideline-fixes\05-single-character-vars.md", "w") as f:
    f.write("# Task: Rename Single-Character Variables\n\n")
    f.write("According to the Coding Guidelines: 'No single-character variables (`s`, `x`, `d`).'\n\n")
    f.write("## Files to fix\n")
    for line in single_chars:
        f.write(f"- {line}\n")

with open(r"d:\wp-work\riseup-asia\gitmap\.lovable\plans\pending\01-coding-guideline-fixes.md", "a") as f:
    f.write("- **Single-Character Vars:** Found 2,300+ instances of single character variables; queued 50 for renaming.\n")
    f.write("- [ ] `.lovable/plans/subtasks/01-coding-guideline-fixes/05-single-character-vars.md`\n")
