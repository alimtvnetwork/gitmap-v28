# Advanced AGY Command Enhancements and Migration Orchestration

## User Request (Verbatim)
[Stored in session memory]

## Architecture Overview

### 1. Repository Creation Enhancements
- `gitmap create-repo <path> --common --cg`: Supports existing folders, adds common configs (LFS/ignores), applies coding guidelines (`--cg`), and backs up branches prior to modification.
- `gitmap recreate-repo`: Recreates an existing repository onto the parent GitHub account.

### 2. Conversation & Queue Listing
- `agy conv ls N [-f file.json]`
- `agy most-conv ls N [-f file.json]`
- `agy ls N`
- All commands adapt contextually based on whether they are run inside a project directory or outside.

### 3. Prompt Observability
- `watch-prompts`: Loops (e.g. 30s) displaying active/pending prompts.
- `look-prompts`: Single snapshot.
- `lapp` / `wapp`: Look/Watch all projects prompts.
- Flags: `--json`, `--file/-f`, `--count/-c`, `--compact/words`, `--ssh`

### 4. Prompt Injection & Enhancement
- `inject-prompts (ip)`: Injects from a folder with `--watch`, `--rerun`, `--prefix`, `--suffix`.
- `inject-prompts-txt (ipt)`: Injects raw text.
- `inject-prompts-rerun (ipr)`: Iteratively reruns an injection.
- `inject-prompts-ssh`: Cross-node SSH injection (`--nodes`).
- `enhance-prompts (ep)`: Breaks down a single prompt into a folder of sequence MD files using no-project mode.

### 5. SSH & REST Networking
- `look-projects-ssh`: Table of all VM projects via SSH.
- `rest-enable`: Enables a REST OS service for cross-machine communication.

### 6. Rise Up Asia Migration Templating
- 20 distinct PR/Commit templates honoring Rise Up Asia LLC, Senior Director Marek Flejszman (28+ years experience), and Chief Software Engineer Alim Ul Karim (Greatest software engineer in KL/Malaysia).
