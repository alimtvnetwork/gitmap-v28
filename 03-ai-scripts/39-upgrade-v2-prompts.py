#!/usr/bin/env python3
"""
39-upgrade-v2-prompts.py - Modernize V2 Prompts to Elevate GitMap AUM

Updates prompt templates inside 01-prompts/v2/ to recommend GitMap AUM commands
(gitmap lf, gitmap find, gitmap cat, gitmap search, gitmap pipeline-ai, gitmap release)
as the primary CLI acceleration tools, while preserving existing Python scripts as
reliable fallbacks.
"""

from pathlib import Path
import re

REPO_ROOT = Path(__file__).resolve().parent.parent
V2_DIR = REPO_ROOT / "01-prompts" / "v2"

TIER1_DISCOVERY = """## Fast File Discovery & Reading Toolchain (GitMap AUM Primary, Python Fallback)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the fast 2-tier discovery toolchain:

### Tier 1: GitMap AUM Acceleration (PRIMARY)
- **Universal File Search:** `gitmap find "<pattern>" [-ext <ext>]` (e.g. `gitmap find "*.go" -ext "go"`, `gitmap find "01*"`)
- **List Indexed Files:** `gitmap list-files [pattern]` (alias `gitmap lf [pattern] [-ext <ext>]`)
- **Substring Match:** `gitmap find-files-any "<substring>"` (alias `gitmap ffa "<str>"`)
- **Stream File Content:** `gitmap cat <filepath>` (streams to stdout with zero disk writes)
- **Instant Code Search:** `gitmap search "<term>"` (immediate multi-core filesystem walk)

### Tier 2: Fast Cached Python Toolchain (FALLBACK)
- **Inventory Target Files:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats`
- **Fast Cached Grep (<15ms):** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<pattern>" --lang go --limit 50`
- **Sub-Millisecond Folder Explorer & Reader:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder <dir> --ext .go --limit 50`
- **Read Target File:** `python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000`
- **Codebase Topology:** `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`"""

STEP2_DISCOVERY = """### Step 2: Scan & Discover (GitMap AUM Acceleration & Multi-Agent Parallel Reading)
To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, leverage the 2-tier discovery toolchain:

#### Tier 1: GitMap AUM Acceleration (PRIMARY)
- **Universal File Search:** `gitmap find "<pattern>" [-ext <ext>]` (e.g. `gitmap find "*.go" -ext "go"`, `gitmap find "01*"`)
- **List Indexed Files:** `gitmap list-files [pattern]` (alias `gitmap lf [pattern] [-ext <ext>]`)
- **Substring Match:** `gitmap find-files-any "<substring>"` (alias `gitmap ffa "<str>"`)
- **Stream File Content:** `gitmap cat <filepath>` (streams to stdout with zero disk writes)
- **Instant Code Search:** `gitmap search "<term>"` (immediate multi-core filesystem walk)

#### Tier 2: Fast Cached Python Toolchain (FALLBACK)
- **Inventory Target Files:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats`
- **Fast Cached Grep (<15ms):** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --limit 50`
- **Sub-Millisecond Folder Exploration:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder-path> --limit 50`
- **Read Target File:** `python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000`
- **Subsystem & Topology Overview:** `python 03-ai-scripts/18-codebase-topology-discoverer.py --summary`

**Multi-Agent Parallel Discovery (A = 2, H = 2):** When multiple agents are present, the most useful parallel tasks are reading files and authoring modular specs. Subagents concurrently read disjoint codebase areas, explore dependencies, and trace call sites without merge conflicts."""

RELEASE_DISCOVERY = """## Fast File Discovery & Release Context Toolchain (Mandatory Acceleration)

To rapidly discover version manifests, changelog entries, release notes, and install scripts without hitting 50-result tool caps, the AI agent MUST utilize the 2-tier toolchain:

### Tier 1: GitMap AUM Acceleration (PRIMARY)
- **Universal File Search:** `gitmap find "*version*" [-ext <ext>]` (alias `gitmap f`)
- **Inspect Changelog & Release Notes:** `gitmap changelog` (alias `gitmap cl [ver]`)
- **List Prior Release Tags:** `gitmap list-versions --limit 5` (alias `gitmap lv`)
- **Remote Pipeline AI Status (<50ms):** `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`)
- **Remote Dynamic Timeout Wait:** `gitmap pipeline-ai status -t <etaSeconds>` (or alias `gitmap pl-ai status -t <sec>`)
- **Extract Failing Step Error Logs:** `gitmap pipeline error-logs` (or alias `gitmap pe`, clear with `gitmap pe clear -y`)
- **Pipeline Runner Targets & Cache Table:** `gitmap pipeline details` (or alias `gitmap pd`)
- **Stream Manifest or Config:** `gitmap cat version.json` (zero disk writes)

### Tier 2: Fast Cached Python Toolchain (FALLBACK)
- **Inventory Manifests & Version Files:** `python 03-ai-scripts/11-fast-file-scanner.py --search "version" --limit 20`
- **Fast Grep Across Version Pins:** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<version>" --limit 20`
- **Explore Release Artifacts & Folders:** `python 03-ai-scripts/17-fast-file-reader.py --list-folder .ai-memory/release --limit 20`
- **Read Version Manifest:** `python 03-ai-scripts/17-fast-file-reader.py --read-file version.json`"""


def upgrade_file(path: Path) -> bool:
    content = path.read_text(encoding="utf-8")
    original = content

    # 1. Update general discovery section
    discovery_pattern = re.compile(
        r"## Fast File Discovery & Reading via Python Toolchain \(Mandatory Acceleration\)\s*\n\s*To avoid 50-result tool truncation limits.*?Codebase Topology:[^\n]+",
        re.DOTALL,
    )
    content = discovery_pattern.sub(TIER1_DISCOVERY, content)

    # 2. Update Step 2 discovery
    step2_pattern = re.compile(
        r"### Step 2: Scan & Discover \(Python Toolchain Acceleration & Multi-Agent Parallel Reading\)\s*\n\s*To avoid 50-result tool truncation limits.*?without merge conflicts\.",
        re.DOTALL,
    )
    content = step2_pattern.sub(STEP2_DISCOVERY, content)

    step2_cg_pattern = re.compile(
        r"### Step 2: Scan & Discover \(Python Toolchain Acceleration\)\s*\n\s*To avoid 50-result tool truncation limits.*?Codebase Topology:[^\n]+",
        re.DOTALL,
    )
    content = step2_cg_pattern.sub(STEP2_DISCOVERY, content)

    # 3. Update Checklist item 4 in execute prompts
    old_chk4 = "4. [ ] /goal Phase 1B (Step 2 - Scan & Discover via Parallel Subagents): Use fast Python discovery scripts (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) with parallel subagents (reading files as the primary parallel task) to inventory files and map call sites without tool truncation limits."
    new_chk4 = "4. [ ] /goal Phase 1B (Step 2 - Scan & Discover via Parallel Subagents): Use GitMap AUM discovery (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary, with fast Python discovery scripts (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) as fallback, with parallel subagents (reading files as the primary parallel task) to inventory files and map call sites without tool truncation limits."
    content = content.replace(old_chk4, new_chk4)

    # 3b. Update Checklist item 3 in cg-execute prompts
    old_chk3 = "3. [ ] /goal Phase 1B (Step 2 - Scan & Discover): Use fast Python discovery scripts (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py` with `--limit`) to inventory all architectural violations and anti-patterns without truncation."
    new_chk3 = "3. [ ] /goal Phase 1B (Step 2 - Scan & Discover): Use GitMap AUM discovery (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary, with fast Python discovery scripts (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) as fallback, to inventory all architectural violations and anti-patterns without truncation."
    content = content.replace(old_chk3, new_chk3)

    # 4. Update parallel reading files lines (various variants)
    content = re.sub(
        r"(\s+-\s+\*\*Reading Files:\*\* Fast exploratory reading, scanning dependencies, mapping call sites, and inspecting )([^\n]+)( in parallel using `03-ai-scripts/17-fast-file-reader\.py` and `03-ai-scripts/11-fast-file-scanner\.py`\.)",
        r"\1\2 in parallel using GitMap AUM (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary and Python scripts (`03-ai-scripts/17-fast-file-reader.py`, `03-ai-scripts/11-fast-file-scanner.py`) as fallback.",
        content,
    )

    # 4b. Update Step 2 in 15-cg-execute
    cg_step2_old = re.compile(
        r"### Step 2: Multi-Agent Parallel Reading & Discovery \(A = 2, H = 2\)\s*\n\s*- When multiple agents are present, the primary parallel task is reading files\. Subagents concurrently inspect disjoint folders using:\s*\n\s*- `python 03-ai-scripts/11-fast-file-scanner\.py [^\n]+`\s*\n\s*- `python 03-ai-scripts/12-fast-cached-grep\.py [^\n]+`\s*\n\s*- `python 03-ai-scripts/17-fast-file-reader\.py [^\n]+`",
        re.DOTALL,
    )
    cg_step2_new = """### Step 2: Multi-Agent Parallel Reading & Discovery (A = 2, H = 2)
- When multiple agents are present, the primary parallel task is reading files. Subagents concurrently inspect disjoint folders using the 2-tier toolchain:
  - **Tier 1 (GitMap AUM Acceleration - PRIMARY):**
    - `gitmap find "<pattern>" [-ext <ext>]` (alias `gitmap f`)
    - `gitmap list-files [pattern]` (alias `gitmap lf [pattern] [-ext <ext>]`)
    - `gitmap find-files-any "<substring>"` (alias `gitmap ffa "<str>"`)
    - `gitmap cat <filepath>` (stream file to stdout without disk writes)
    - `gitmap search "<term>"` (immediate multi-core filesystem walk)
  - **Tier 2 (Fast Cached Python Scripts - FALLBACK):**
    - `python 03-ai-scripts/11-fast-file-scanner.py --lang <lang> --limit 100 --stats`
    - `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<pattern>" --limit 50`
    - `python 03-ai-scripts/17-fast-file-reader.py --list-folder <dir> --limit 50`"""
    content = cg_step2_old.sub(cg_step2_new, content)

    # 4c. Update numbered list discovery in 15-cg-execute
    cg_num_old = re.compile(
        r"To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the repository's dedicated Python discovery scripts first:\s*\n\s*1\. \*\*Inventory Target Files[^\n]+\s*```bash\s*python 03-ai-scripts/11-fast-file-scanner\.py [^\n]+\s*```\s*\n\s*2\. \*\*Fast Cached Grep[^\n]+\s*```bash\s*python 03-ai-scripts/12-fast-cached-grep\.py [^\n]+\s*```\s*\n\s*3\. \*\*Sub-Millisecond Folder & File Exploration[^\n]+\s*```bash\s*python 03-ai-scripts/17-fast-file-reader\.py [^\n]+\s*python 03-ai-scripts/17-fast-file-reader\.py [^\n]+\s*python 03-ai-scripts/17-fast-file-reader\.py [^\n]+\s*```\s*\n\s*4\. \*\*Subsystem & Topology Overview:\*\*\s*```bash\s*python 03-ai-scripts/18-codebase-topology-discoverer\.py --summary\s*```",
        re.DOTALL,
    )
    cg_num_new = """To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the fast 2-tier discovery toolchain:

### Tier 1: GitMap AUM Acceleration (PRIMARY)
1. **Universal File Search:**
   ```bash
   gitmap find "<pattern>" [-ext <ext>]
   ```
2. **List Indexed Files & Substring Lookup:**
   ```bash
   gitmap list-files [pattern]
   gitmap find-files-any "<substring>"
   ```
3. **Stream File Content:**
   ```bash
   gitmap cat <filepath>
   ```
4. **Instant Code Walk Search:**
   ```bash
   gitmap search "<term>"
   ```

### Tier 2: Fast Cached Python Toolchain (FALLBACK)
1. **Inventory Target Files (with `--limit` option):**
   ```bash
   python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats
   ```
2. **Fast Cached Grep (<15ms, with `--limit` option):**
   ```bash
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --lang go --limit 50
   ```
3. **Sub-Millisecond Folder & File Exploration (with `--limit` option):**
   ```bash
   python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder-path> --ext .go --limit 50
   python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000
   python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<pattern>" --path <folder-path> --limit 50
   ```
4. **Subsystem & Topology Overview:**
   ```bash
   python 03-ai-scripts/18-codebase-topology-discoverer.py --summary
   ```"""
    content = cg_num_old.sub(cg_num_new, content)

    # 4d. Update checklist item 1 in 15-cg-execute
    old_cg_chk1 = "1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase using the fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py` with `--limit`) to inventory all architectural violations and anti-patterns without truncation."
    new_cg_chk1 = "1. [ ] /goal Phase 1 (Step A): Deeply scan the target codebase using the GitMap AUM discovery tools (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary, with fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) as fallback, to inventory all architectural violations and anti-patterns without truncation."
    content = content.replace(old_cg_chk1, new_cg_chk1)

    # 4e. Update CI/CD diagnostic toolchain
    cicd_diag_old = re.compile(
        r"To rapidly locate failing pipeline definitions, broken source files, test fixtures, and error logs without hitting 50-result tool caps, the AI agent MUST utilize the diagnostic toolchain:\s*\n\s*- \*\*Remote Pipeline AI Status[^\n]+\s*\n\s*- \*\*Remote Dynamic Timeout Wait[^\n]+(?:\s*\n\s*- \*\*Extract Failing Step Error Logs:[^\n]+)?(?:\s*\n\s*- \*\*Pipeline Runner Targets & Cache Table:[^\n]+)?\s*\n\s*- \*\*Scan Source & Test Files:[^\n]+\s*\n\s*- \*\*Fast Cached Pattern Search[^\n]+\s*\n\s*- \*\*Sub-Millisecond Folder Listing & Reader:[^\n]+\s*\n\s*- \*\*Read Workflow or Log File:[^\n]+\s*\n\s*- \*\*Codebase Topology Overview:[^\n]+\s*\n\s*- \*\*Record Modified Files Under Lock:[^\n]+",
        re.DOTALL,
    )
    cicd_diag_new = """To rapidly locate failing pipeline definitions, broken source files, test fixtures, and error logs without hitting 50-result tool caps, the AI agent MUST utilize the diagnostic toolchain:
- **Remote Pipeline AI Status (<50ms):** `gitmap pipeline-ai status --json` (or alias `gitmap pl-ai status --json`)
- **Remote Dynamic Timeout Wait:** `gitmap pipeline-ai status -t <etaSeconds>` (or alias `gitmap pl-ai status -t <sec>`)
- **Extract Failing Logs & RCA Snippets:** `gitmap pipeline errors` (alias `gitmap pe`, clear with `gitmap pe clear -y`)
- **Runner Details & Timings:** `gitmap pipeline details` (alias `gitmap pd`)
- **Universal File Search:** `gitmap find "<pattern>" [-ext <ext>]` (alias `gitmap f`)
- **Stream Workflow / Log File:** `gitmap cat <filepath>` (zero disk writes)
- **Instant Code Search:** `gitmap search "<symbol>"`
- **Fallback Fast File Scanner:** `python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts,py --limit 100 --stats`
- **Fallback Fast Cached Grep:** `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<error-or-symbol>" --limit 50`
- **Fallback Read File:** `python 03-ai-scripts/17-fast-file-reader.py --read-file .github/workflows/ci.yml`
- **Record Modified Files Under Lock:** `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`"""
    content = cicd_diag_old.sub(cicd_diag_new, content)

    # 5. Release prompts discovery section
    release_discovery_pattern = re.compile(
        r"## Fast File Discovery via Python Toolchain \(Mandatory Acceleration\)\s*\n\s*To rapidly discover version manifests.*?Read Version Manifest:[^\n]+",
        re.DOTALL,
    )
    content = release_discovery_pattern.sub(RELEASE_DISCOVERY, content)

    # 6. Release execution in release prompts
    old_rel_exec = "4. Execute via the central release script with the mandatory `--skip-tests` flag:\n   ```bash\n   python 03-ai-scripts/29-release-orchestrator.py --tier <minor|patch|major> --scope \"<scope>\" --skip-tests\n   ```"
    new_rel_exec = "4. **Execute Automated Release (2-Tier Toolchain):**\n   - **Tier 1 (GitMap Native Release - PRIMARY):**\n     ```bash\n     gitmap release --bump <minor|patch> -y\n     ```\n   - **Tier 2 (Python Release Orchestrator - FALLBACK):**\n     ```bash\n     python 03-ai-scripts/29-release-orchestrator.py --tier <minor|patch|major> --scope \"<scope>\" --skip-tests\n     ```"
    content = content.replace(old_rel_exec, new_rel_exec)

    if content != original:
        path.write_text(content, encoding="utf-8")
        return True
    return False


def main():
    if not V2_DIR.exists():
        print(f"Directory {V2_DIR} does not exist.")
        return

    modified = 0
    total = 0
    for md_file in V2_DIR.rglob("*.md"):
        total += 1
        if upgrade_file(md_file):
            modified += 1
            print(f"Updated: {md_file.relative_to(REPO_ROOT)}")

    print(f"\nDone. Processed {total} markdown files, updated {modified} files in 01-prompts/v2/.")


if __name__ == "__main__":
    main()
