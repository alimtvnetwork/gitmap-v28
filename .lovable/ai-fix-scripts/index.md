# AI Fix Scripts Memory

This directory stores reusable helper scripts that AIs can invoke natively rather than generating temporary single-use code for complex fixes.

<details>
<summary>01-file-manipulator.py</summary>

**Purpose:** Automates mass file lowercasing and automatic filename sequence numbering using standard libraries. It defaults to ignoring `.git` and `node_modules` and prefers `git mv` to preserve git history.

**Usage Examples:**
```bash
python .lovable/ai-fix-scripts/01-file-manipulator.py lowercase ./src --except "vendor/*, *.png"
python .lovable/ai-fix-scripts/01-file-manipulator.py fix-seq-files ./docs --order-by-time
python .lovable/ai-fix-scripts/01-file-manipulator.py fix-seq-files ./docs --order-by-az --keep-old-order --pin "readme=00,draft=01"
```
</details>
