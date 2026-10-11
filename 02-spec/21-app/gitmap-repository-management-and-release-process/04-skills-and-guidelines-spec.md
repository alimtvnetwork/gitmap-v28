# Skills & Coding Guidelines Specification: Python Subprocess Ban & Native Golang CLI Acceleration

## 1. Executive Summary & Visual Evidence

This specification formalizes the updates required across the repository's agent skills (`.agents/skills/gitmap/SKILL.md`) and canonical coding guidelines (`02-spec/02-coding-guidelines/`). It establishes an absolute ban on AI agents executing Python subprocess wrappers to query GitMap information and mandates direct, native Golang CLI commands (`gitmap rc --path`, `gitmap path rc`, `gitmap rs --path`) that execute in under $5\text{ms}$.

### 1.1 Visual Ground Truth

![Screenshot 02: Python Subprocess Inefficiency](assets/screenshots/gitmap-repository-management-and-release-process-02.png)

### 1.2 Analysis of the Inefficiency
As captured in Screenshot 02, AI coding agents frequently attempt to discover the repository cache path by synthesizing multi-line Python scripts:
```python
# ❌ ERRONEOUS AI ANTI-PATTERN (Captured in Screenshot 02)
gitmap py -c "import subprocess; print(subprocess.check_output(['gitmap', 'cd', 'rc'], text=True))"
```

#### Why This Pattern Fails:
1. **Severe Latency Overhead:** Spawns a Python interpreter process, imports `subprocess`, executes a child `gitmap.exe` process, waits for process termination, and decodes stdout. This burns over $1000\text{ms}$ of agent execution time.
2. **Interactive Shell Warning Collision:** `gitmap cd rc` is engineered for interactive terminal directory switching. When invoked in an unbuffered subprocess without a shell wrapper hook, it emits interactive shell warnings (`Warning: shell wrapper not active...`) and generates shell handoff files (`cmdas.WriteShellHandoff`).
3. **Fragile Output Parsing:** The calling script must strip newline characters, filter out warning banners, and handle cross-platform slash delimiters.
4. **Agent Resource Burn:** Burns context tokens and unnecessary tool steps.

---

## 2. The Native Golang Solution

GitMap provides first-class, ultra-fast Go CLI commands that resolve repository and cache locations directly in $<5\text{ms}$:

```bash
# ✅ MANDATORY NATIVE GOLANG COMMANDS
gitmap rc --path          # Outputs clean directory path of repo-cache directly to stdout
gitmap rc path            # Positional syntax alternative
gitmap rs --path          # Outputs clean directory path of repo-secrets directly to stdout
gitmap path rc            # Universal path resolver for repo-cache
gitmap path rs            # Universal path resolver for repo-secrets
gitmap path <repo-name>   # Universal path resolver for any tracked workspace repo
gitmap path rc --json     # Emits structured JSON: {"name":"repo-cache","path":"...","isValid":true}
```

### 2.1 Execution Contract for Native Path Queries
- **Sub-5ms Execution:** Resolves directory from Split-DB (`gitmap-special-repos.db` or `gitmap.db`) in memory without spawning secondary shells.
- **Zero Shell Warnings:** Bypasses `cmdsetup.WarnIfNoWrapper()` completely.
- **Zero Handoff Files:** Bypasses `cmdas.WriteShellHandoff()` completely.
- **Clean Stdout:** Emits purely the requested path string with a trailing newline (or JSON object if `--json` is supplied).

---

## 3. Required Updates to `.agents/skills/gitmap/SKILL.md`

### 3.1 Addition to Non-Negotiable Command Replacement Matrix

The Command Replacement Matrix in `.agents/skills/gitmap/SKILL.md` MUST be updated with the following entry:

| Prohibited Legacy Tool / Pattern | GitMap Mandatory Command | Why GitMap is Required |
| :--- | :--- | :--- |
| ❌ Python subprocess wrappers to query GitMap paths (`gitmap py -c "import subprocess; print(subprocess.check_output(['gitmap', 'cd', 'rc']))"`) | ✅ `gitmap rc --path`, `gitmap path rc`, or `gitmap rs --path` | Zero-subprocess native Go path resolution in $<5\text{ms}$. Eliminates multi-runtime memory overhead, shell wrapper warning noise, and stdout parsing bugs. |

### 3.2 Addition to Essential Command Cheat Sheet

Add a dedicated section under Section 5 (Special Repositories & Fast Path Resolution):

```markdown
### 5. Special Repositories & High-Speed Path Resolution
- `gitmap rc` — Interactive navigation to repo-cache storage.
- `gitmap rc --path` (or `gitmap rc path`) — Native high-speed path query; emits clean repo-cache directory path to stdout in <5ms.
- `gitmap rs` — Interactive navigation to repo-secrets vault.
- `gitmap rs --path` (or `gitmap rs path`) — Native high-speed path query; emits clean repo-secrets directory path to stdout in <5ms.
- `gitmap path <target>` — Universal path resolver for `rc`, `rs`, or any repository by name/alias.
- `gitmap path <target> --json` — Emits structured JSON record containing `name`, `path`, and positive boolean `isValid`.
```

---

## 4. Required Updates to `02-spec/02-coding-guidelines/`

### 4.1 Update: `02-spec/02-coding-guidelines/06-ai-optimization/04-common-ai-mistakes.md`

Add Mistake #16 to the Top AI Mistakes catalog:

```markdown
### Mistake #16: Python Subprocess Wrappers for GitMap CLI Queries

**Frequency:** High
**Rule:** AH-P4

```python
# ❌ AI GENERATES THIS (Subprocess wrapper burning 1000ms+ and triggering shell warnings)
import subprocess
path = subprocess.check_output(["gitmap", "cd", "rc"], text=True).strip()

# ✅ CORRECT (AI executes native GitMap query directly via run_command)
# Execute: gitmap rc --path
# Output: work-dir/repo-cache
```

**Rationale:** GitMap provides native Go CLI query flags (`gitmap rc --path`, `gitmap path rc`, `gitmap rs --path`) that execute in $<5\text{ms}$ with clean stdout. Spawning Python subprocesses introduces runtime latency, consumes process memory, and triggers interactive shell wrapper warnings.
```

### 4.2 Update: `02-spec/02-coding-guidelines/12-python/02-standards.md`

Add Section 5 (GitMap CLI Integration & Subprocess Wrapper Prohibition):

```markdown
## 5. GitMap CLI Integration & Subprocess Wrapper Prohibition

- Python scripts MUST NEVER wrap `gitmap` CLI commands using `subprocess.check_output` or `subprocess.run` to query directory paths, configuration values, or repository statuses.
- Use direct GitMap native CLI commands via shell execution:
  - For repo-cache path: run `gitmap rc --path` or `gitmap path rc`.
  - For repo-secrets path: run `gitmap rs --path` or `gitmap path rs`.
  - For repository path: run `gitmap path <repo-name>`.
- Python scripts that require repository metadata should consume JSON outputs emitted by native commands (e.g., `gitmap status --json`, `gitmap path <target> --json`).
```

### 4.3 Update: `02-spec/02-coding-guidelines/03-golang/`

Document the positive boolean requirements for path query models:
- Struct fields and options MUST use positive prefixes: `isPathOnly`, `hasRepo`, `isValidPath`, `isJSON`.

---

## 5. Enforcement & Automated Verification Checklist

- [ ] Verify that `.agents/skills/gitmap/SKILL.md` includes the Python subprocess ban in its replacement matrix.
- [ ] Verify that `.agents/skills/gitmap/SKILL.md` documents `gitmap rc --path`, `gitmap rs --path`, and `gitmap path`.
- [ ] Verify that `02-spec/02-coding-guidelines/06-ai-optimization/04-common-ai-mistakes.md` indexes Mistake #16.
- [ ] Verify that `02-spec/02-coding-guidelines/12-python/02-standards.md` includes Section 5.
- [ ] Confirm all documentation uses relative paths exclusively (`02-spec/...`, `.agents/...`, `cli/...`).
- [ ] Confirm all boolean variables and options use affirmative `is`/`has` naming.
