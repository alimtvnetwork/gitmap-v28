# Execution Plan: GitMap Repo-Cache Scan, Clone, and LS Ecosystem

- **Task Slug:** `gitmap-repo-cache-scan-clone-and-ls-ecosystem`
- **Specification:** [02-spec/21-app/gitmap-repo-cache-scan-clone-and-ls-ecosystem/01-architecture-spec.md](file:///d:/work/gitmap/02-spec/21-app/gitmap-repo-cache-scan-clone-and-ls-ecosystem/01-architecture-spec.md)
- **Status:** APPROVED & IN PROGRESS
- **Budget:** 300 steps
- **Concurrency (A=2, H=2):** 2 Subagents, 2 Hands (4 concurrent operations capacity)

---

## 1. User Request (Verbatim)

```text
gitmap scan . --rc # keep the json also in the repo cache inside 01-gitmap/gitmap.json if already exist try to merge 
gitmap scan . --rc --seprate/sep # keep the json also in the repo cache inside 01-gitmap/gitmap.json if already exist then create next one 02-gitmap.json

gitmap scan --rc --seprate/sep # keep the json also in the repo cache inside 01-gitmap/gitmap.json if already exist then create next one 02-gitmap.json

gitmap scan --rc --s # keep the json also in the repo cache inside 01-gitmap/gitmap.json if already exist then create next one 02-gitmap.json

gitmap clone/cfr/cfrp --rc/repo-cache # automatically scan gitmap.json or type of json available the frist try would be in the 01-gitmap/gitmap.json , then 02-, 03- json file or also try to look into other json files to understand format if fesible also keep a cache of it in the split db so that doens't need to parse all the time and shows a list of repo wants to clone first , user says -y or select a number from the json list each json will also say a short list of repos incliding the user name but not the full url as shorrt as possibl;e as treeview and user can choose they want to clone all, specific , or rotate one json at a time

gitmap clone/cfr/cfrp rc  # automatically scan gitmap.json or type of json available the frist try would be in the 01-gitmap/gitmap.json , then 02-, 03- json file or also try to look into other json files to understand format if fesible also keep a cache of it in the split db so that doens't need to parse all the time and shows a list of repo wants to clone first , user says -y or select a number from the json list each json will also say a short list of repos incliding the user name but not the full url as shorrt as possibl;e as treeview and user can choose they want to clone all, specific , or rotate one json at a time

gitmap ls rc # shows the list of json, and which json contains which repo and in which format public or ssh format, show a tree view summary to pick fromt he list and along with command to import it, clear??

now write detailed spec first for this in spec folder, please
```

---

## 2. Discrete Subtasks & File Ownership Matrix

| Subtask ID | Title | Owner | Owned Files | Status |
|:---|:---|:---|:---|:---:|
| **Task-01** | Scan Engine: `--rc`, `--separate` (`--sep`, `--s`), Auto-Merge & Sequential Allocator | Worker 01 | `cli/cmdscan/flags.go`, `cli/cmdscan/scan.go`, `cli/cmdscan/scan_rc_export.go`, `cli/cmdscan/scan_rc_merge.go`, `cli/cmdscan/scan_rc_test.go` | PENDING |
| **Task-02** | Clone & LS Engine: `clone/cfr/cfrp --rc/rc`, Split-DB Cache, Short TreeView & `ls rc` | Worker 02 | `cli/cmdclone/clone_rc.go`, `cli/cmdclone/clone_rc_db.go`, `cli/cmdclone/clone_rc_tree.go`, `cli/cmdclone/clone.go`, `cli/cmdclone/clonefixrepo.go`, `cli/cmdlist/list_rc.go`, `cli/cmdlist/list.go`, `cli/cmdclone/clone_rc_test.go` | PENDING |

---

## 3. Execution Phases

### Phase 1: Planning & Specification (Complete)
- Authored `01-architecture-spec.md` and `02-component-and-cli-spec.md`.
- Enqueued master plan and subtask plans.
- Initialized SQLite Task DB via `gitmap agent-ai create-task`.

### Phase 2: Parallel Worker Dispatch (A=2, H=2)
- Worker 01 executes Subtask 01 in disjoint file box `cli/cmdscan/`.
- Worker 02 executes Subtask 02 in disjoint file boxes `cli/cmdclone/` and `cli/cmdlist/`.
- Workers strictly enforce Rule R-PRETOUCH via `gitmap agent-ai editing` and declare file boxes via `gitmap agent subtask claim-files`.

### Phase 3: Verification & Atomic Push
- Run targeted tests: `go test ./cmdscan -run TestScanRC` and `go test ./cmdclone ./cmdlist -run TestCloneRC`.
- Run linters: `check-relative-paths.py -c` and `check-forbidden-strings.py`.
- Atomic commit & push: `gitmap cpf "feat - add repo-cache scan merge, sequential allocator, clone rc treeview and ls rc inspector"`.
