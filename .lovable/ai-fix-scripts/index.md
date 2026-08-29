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

<details>
<summary>02-guideline-autofixer.py</summary>

**Purpose:** Automates repository-wide coding guideline compliance, formatting empty lines before returns, removing explicit boolean comparisons (`== true`), standardizing error handling envelopes, and enforcing PascalCase naming.

**Usage Examples:**
```bash
python .lovable/ai-fix-scripts/02-guideline-autofixer.py <file1> <file2>
```
</details>

<details>
<summary>03-cicd-local-runner.py</summary>

**Purpose:** Auto-generated concurrent local CI/CD test runner with worker pool architecture, real-time job timing, graceful error handling, and aggregated log reporting.

**Usage Examples:**
```bash
python .lovable/ai-fix-scripts/03-cicd-local-runner.py
```
</details>

<details>
<summary>04-relative-path-fixer.py</summary>

**Purpose:** Autonomously converts all absolute filesystem paths (`D:\...`, `C:\...`, `/home/...`) and `file:///` URIs across the codebase to strictly relative Git root paths.

**Usage Examples:**
```bash
python .lovable/ai-fix-scripts/04-relative-path-fixer.py
```
</details>

<details>
<summary>05-naming-autofixer.py</summary>

**Purpose:** Scans and refactors variable and boolean naming violations across Go, TypeScript, and Python codebases, replacing bare `ok` identifiers with domain-specific affirmative booleans, eliminating negative booleans (`hasNo*`, `isNot*`), and enforcing positive framing with inverted guard clauses.

**Usage Examples:**
```bash
python .lovable/ai-fix-scripts/05-naming-autofixer.py
python .lovable/ai-fix-scripts/05-naming-autofixer.py gitmap/cmd
```
</details>

