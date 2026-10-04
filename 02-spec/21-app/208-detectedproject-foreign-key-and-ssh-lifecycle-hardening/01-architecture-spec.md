# Spec 208: DetectedProject Foreign Key RCA, Comprehensive Database Reset, SSH Lifecycle & Cross-OS Firewall Hardening

## 1. Specification Metadata
- **Spec ID:** 208
- **Title:** DetectedProject Foreign Key RCA, Comprehensive Database Reset, SSH Lifecycle & Cross-OS Firewall Hardening
- **Status:** Completed
- **Created Date:** 2026-10-04
- **Related Plans:** [Plan 78](../../.ai-memory/plans/completed/78-detectedproject-foreign-key-and-ssh-lifecycle-hardening.md)
- **Domain:** Scanner Subsystem, Split Database Engine, OpenSSH Fleet Management, Cross-OS Firewall

---

## 2. User Request (Verbatim)

```text
  Note: safe-pull is auto-enabled when existing repos are detected.

gitmap: json: wrote 37 record(s), 0 validation issue(s)
gitmap: json: wrote 37 record(s), 0 validation issue(s)
  • Cache       last-scan.json

■ Database
────────────────────────────────────────────
  ✔ 37 repositories upserted into database
  • Tagged 37 repo(s) with scan folder #1
  ✔ Auto-generated 79 repository alias(es)
  ↪ background probe queued for 37 repo(s) (workers=3)

■ Project Detection
────────────────────────────────────────────
  [nav] Detected 53 project(s) across 35 repo(s)
  - go-projects.json       15 record(s)
  - react-projects.json    31 record(s)
  - node-projects.json     7 record(s)
[QueryWrapper Error]: exec failed: constraint failed: FOREIGN KEY constraint failed (787)
query: INSERT INTO DetectedProject
        (RepoId, ProjectTypeId, ProjectName, AbsolutePath, RepoPath, RelativePath, PrimaryIndicator)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(RepoId, ProjectTypeId, RelativePath) DO UPDATE SET
                ProjectName=excluded.ProjectName,
                AbsolutePath=excluded.AbsolutePath,
                RepoPath=excluded.RepoPath,
                PrimaryIndicator=excluded.PrimaryIndicator,
                DetectedAt=CURRENT_TIMESTAMP
failed to upsert detected project: constraint failed: FOREIGN KEY constraint failed (787)
... (8 errors total)
  [OK] Saved 45 detected project(s) to database
  ✔ projects.json synced: 0 added, 0 updated, 37 unchanged (73 total)
  [stats] Benchmark log: ./.gitmap\output\scan-benchmark.log
  [wait] Waiting for background probes to finish (37 remaining)...

# High Priority Instruction

Okay, good work so far. I want to have a few more changes inside the Gitmap. First of all, we should be able to reset all the database. If we do the reset, that would actually basically going to remove all the database on the database folder, and it would tell, "We have removed," and then also reseed it, so it would mention that. That is something I think is missing when we do the Gitmap reset. Gitmap reset is going to reset everything, every trace of previous information from the system to Gitmap memory. So it would clean it. Remember that. That's the first thing. Second, when we do the Gitmap SSH, it should show the SSH information at the end, not in the middle. The second thing is that if you do the create, Gitmap SSH create, it can give the SSH already there. Then it's going to ask the user, "Are you going to overwrite on the same one?" Because it needs to also pass the prompt Y to when it is running as a wrap. Now, if the user does not provide hyphen Y or hyphen confirm at the beginning of the run, then it would prompt the user, "Do you want to overwrite?" If the user wants to overwrite, then they would do it. But also at the same time, we keep the backup of the previous one so that we can undo and redo and make sure that everything actually enqueued using task. That is a very crucial step. Now, if we do the SSH and enter, it shows commands, it shows SSH correct. We should also do Gitmap SSH space view. It should also only view the SSH stuff, not the commands anymore. Remember that. And it would copy the SSH code as well. So create is going to create the SSH key based on the name or email that is given. The example and the help text in the terminal needs to be updated a bit. And when we view it, there will be a suggestion shown that if you want to recreate or regenerate, you use this command. So this is very crucial. And not only that, the same way, I think we need to update in the scripts fixture project, which is also not regenerating the SSH key. Usually, it should regenerate using... So I'll give you. Also at the same time, there are some issues on another machine that I faced. So I'm just going to give you the error logs so that you can understand if those error logs are important or not, if you are going to fix it or not. Is there any bug or not? Think about that. If there is any bug, try to fix it properly. Now, the next part in the SSH. When we do the SSH, we should have a couple of more commands. We should be able to enable SSH to a specific port. We can enable SSH with firewall. We can disable SSH as well. We can disable certain ports as well. These commands needs to be added, implemented, and that should work regardless of the OS. Try to understand. Also, SSH can be enabled in multiple ports as well. So inside the SSH, this type of enable in ports or disabling ports or change port, these commands needs to be there with help text and also in the UI. And we can also do SSH space UI that would actually open a browser window with a UI that would have all these features that we could try out and see from the browser easily. Now, also, I wanted to have new commands, like if we are in the, let's say, SSH, then we can actually enable something like the... We can enable like SSH to public. We can have commands that actually does this. That means it's going to enable that VM's SSH automatically to public. So make sure that it can go through. That means VM cannot go public, but it can always put it towards the host, so that host can access probably. And from the host, if the Gitmap is accessible, the user can run some commands like Gitmap enable, like in public SSH. Gitmap SSH enable public, that could be another command, and give a port. And then that port, SSH would be, let's say, public available. That means, the VM, if that has a public IP, it can be accessed from outside. And it would automatically handle the firewall situation and everything else. And depending on the OS, we should have a firewall command where we could actually add IP routing, enable, disable. Disable is kind of like deny. So it would automatically work regardless of the Windows, Ubuntu, CentOS or macOS. For Windows, you can test out some of the ports enabling and see whether or not that it works. You can also, enable and disable applications, certain behavior, like it can read from the internet, incoming and outgoing, but cannot do outgoing. So there could be some sort of examples inside the firewall section that we could do, for Windows and also same for, Ubuntu, CentOS and things like that. Now, anytime in Gitmap, if we encounter any help or error, let's say, someone is trying to fail or getting fail into the SSH connect to another VM. So some of the common function or steps we could share. Some of them I have in my mind, but you can be the best to understand what should they follow. You can use your expert opinion and share those suggestions so that user can follow, just in case the SSH is not getting connected. The first of one would be installing the Gitmap into that machine and running the SSH enable, things like that. So these would be first suggestion. Second is checking the machine's IP. Third is checking the machine's username and password. IP and username and password would be step by step. So check those and then try a manual SSH connect first. So this is like suggestion that would be there, and then it would try to use using Gitmap, there would be one liner command that it would try with that IP and machine to connect, using the password and prompting password. And then use the deploy command to deploy the current machine's public key to that machine. So these type of things should be in the suggestion, just in case if a SSH fails. But you can also add your own suggestion that might be more helpful. It is like what I am thinking. I'm just sharing, but there could be better ones. So try to follow these first, then we can discuss the next steps.
```

---

## 3. Subsystem Architectural Blueprint

### 3.1 DetectedProject Foreign Key (787) RCA & In-Depth Fix
1. **Schema Constraint:** `FOREIGN KEY (RepoId) REFERENCES Repo(RepoId)` and `FOREIGN KEY (ProjectTypeId) REFERENCES ProjectType(ProjectTypeId)`.
2. **Analysis of the 8 Failures out of 53:**
   - 15 Go projects, 31 React projects, 7 Node projects detected (53 total). 45 saved, 8 failed.
   - When a project is detected in a directory where `RepoId` resolution cannot locate the parent `RepoId` in SQLite (e.g. subdirectories, nested repositories, or case sensitivity differences on Windows), `RepoId` defaults to 0 or null, violating SQLite foreign key constraints.
   - Additionally, if `node` project type is seeded as ID 3 but queried as 0 or mismatched string, or if `ProjectType` table was created without all project types seeded on older database versions, foreign key 787 triggers.
3. **Remediation:**
   - Enforce mandatory validation: verify `ProjectType` table contains all project types (`Go`, `React`, `Node`, `Python`, `Rust`) with IDs 1 to 5.
   - Add pre-insert repository verification: if `RepoId <= 0`, dynamically look up the repository by prefix matching `AbsolutePath` against all known `RepoPath` entries in SQLite.
   - If no parent repository is found, automatically upsert an ad-hoc repo record or safely skip with an actionable diagnostic warning rather than failing the scan.

### 3.2 Database Reset Engine (`gitmap reset` / `gitmap db reset`)
- Comprehensive discovery of all `.db`, `.db-wal`, `.db-shm` files across:
  - `%LOCALAPPDATA%\gitmap-cli\data`
  - Current workspace `.gitmap/data/`
  - Global cache directories
- Reports each unlinked file with human-readable size.
- Reseeds schemas, executes `db.Migrate()`, and runs `db.SeedProjectTypes()`.

### 3.3 SSH Terminal Layout & Dedicated View Command
- Reorder `gitmap ssh` terminal output: command reference rendered first, public key rendered at bottom.
- Implement `gitmap ssh view` (aliases: `v`, `show`, `key`): displays only formatted key card, copies to clipboard, and prints regeneration hint.

### 3.4 SSH Key Creation Overwrite Guard & Journal Task
- In `gitmap ssh create`, check for existing key pair.
- If present, prompt `[y/N]` unless `-y` / `--yes` / `--confirm` / `-f` is passed.
- Generate timestamped backup (`id_rsa.bak.<unix-timestamp>`).
- Enqueue operation in `TaskHistory` for undo/redo.

### 3.5 Scripts-Fixer Project Alignment
- Audit and patch scripts-fixer project to ensure SSH key regeneration supports `-y`/`--force` and handles existing keys cleanly.

### 3.6 Cross-OS SSH Port Management & Firewall Automation
- `gitmap ssh port <ls|add|rm|set>`
- `gitmap ssh enable [--port <p>]` / `gitmap ssh disable`
- `gitmap ssh enable-public <port>`
- Unified driver `cli/firewall/` supporting Windows, Linux (UFW/firewalld/iptables), and macOS (pfctl).

### 3.7 Interactive SSH Web UI & 7-Step Troubleshooting
- Browser-based UI launcher via `gitmap ssh ui`.
- 7-step remediation workflow for failed SSH target connections.
