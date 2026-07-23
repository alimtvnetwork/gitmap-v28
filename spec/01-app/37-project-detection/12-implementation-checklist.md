# Project Type Detection — Implementation Checklist

## How to Use

Each section maps a spec file to the Go source files it drives.
Check off items as they are implemented. Dependencies are noted
where order matters.

---

## Phase 1: Foundation (no dependencies)

These files have no inter-dependencies and can be built in parallel.

### From `11-constants-project.md`

- [ ] **CREATE** `constants/constants_project.go` — IDs, keys, table
      names, JSON filenames, indicators, exclusion dirs, command
      names, aliases, help text, flag names, messages, errors
- [ ] **CREATE** `constants/constants_project_sql.go` — All SQL:
      create tables, seed, upsert, query, stale cleanup, drop

### From `09-package-structure.md` (model structs)

- [ ] **CREATE** `model/projecttype.go` — `ProjectType` struct
- [ ] **CREATE** `model/project.go` — `DetectedProject` struct
- [ ] **CREATE** `model/gometadata.go` — `GoProjectMetadata` +
      `GoRunnableFile` structs
- [ ] **CREATE** `model/csharpmetadata.go` — `CSharpProjectMetadata` +
      `CSharpProjectFile` + `CSharpKeyFile` structs

---

## Phase 2: Storage (depends on Phase 1)

### From `03-data-model.md`

- [ ] **MODIFY** `store/store.go` — Add 7 new create-table statements
      to `Migrate()`, add seed call, add 7 drop statements to `Reset()`
- [ ] **CREATE** `store/projecttype.go` — Seed `ProjectTypes`, query
      by key

### From `03-data-model.md` + `04-go-metadata.md`

- [ ] **CREATE** `store/project.go` — `UpsertDetectedProject`,
      `SelectProjectsByTypeKey`, `CountProjectsByTypeKey`,
      `DeleteStaleProjects`
- [ ] **CREATE** `store/gometadata.go` — `UpsertGoMetadata`,
      `UpsertGoRunnable`, `SelectGoMetadata`, `SelectGoRunnables`,
      `DeleteStaleGoRunnables`

### From `03-data-model.md` + `05-csharp-metadata.md`

- [ ] **CREATE** `store/csharpmetadata.go` — `UpsertCSharpMetadata`,
      `UpsertCSharpProjectFile`, `UpsertCSharpKeyFile`,
      `SelectCSharpMetadata`, `SelectCSharpProjectFiles`,
      `SelectCSharpKeyFiles`, `DeleteStaleCSharpFiles`,
      `DeleteStaleCSharpKeyFiles`

---

## Phase 3: Detection (depends on Phase 1)

### From `02-detection-rules.md`

- [ ] **CREATE** `detector/rules.go` — Per-type indicator matching,
      exclusion dir check, `cmake-build-*` prefix match, React dep
      list check, C# `.sln` precedence logic

### From `02-detection-rules.md` + `08-scan-integration.md`

- [ ] **CREATE** `detector/detector.go` — `DetectProjects(repoPath)`
      entry point, tree walker, collect `[]DetectedProject` + metadata

### From `04-go-metadata.md`

- [ ] **CREATE** `detector/goparser.go` — Parse `go.mod` for module
      name + Go version, locate `go.sum`, scan `cmd/` for runnables,
      check root `main.go`

### From `02-detection-rules.md` (Node/React)

- [ ] **CREATE** `detector/parser.go` — Parse `package.json` name
      field, check dependencies for React indicators, parse
      `CMakeLists.txt` `project()` directive for C++ project name

### From `05-csharp-metadata.md`

- [ ] **CREATE** `detector/csharpparser.go` — Parse `.csproj` XML
      (`TargetFramework`, `OutputType`, `Sdk`), find `.sln`, parse
      `global.json` SDK version, collect key files

---

## Phase 4: Scan Integration (depends on Phases 2 + 3)

### From `08-scan-integration.md`

- [ ] **MODIFY** `cmd/scan.go` — After `BuildRecords`, call
      `detector.DetectProjects` per repo, then upsert projects +
      metadata to DB, then cleanup stale records

### From `06-json-output.md`

- [ ] **MODIFY** `cmd/scanoutput.go` — Write per-type JSON files
      (`go-projects.json`, etc.) to output dir; skip empty types;
      sort by `repoName` then `relativePath`

---

## Phase 5: Query Commands (depends on Phase 2)

### From `07-commands.md`

- [ ] **CREATE** `cmd/projectrepos.go` — `runProjectRepos` handler,
      dispatch by project type key, `--json` and `--count` flags,
      read from DB, error if no DB
- [ ] **CREATE** `cmd/projectreposoutput.go` — Terminal format
      (type + name + path + indicator), JSON format
- [ ] **MODIFY** `cmd/root.go` — Register `go-repos`, `node-repos`,
      `react-repos`, `cpp-repos`, `csharp-repos` + aliases in
      command dispatch

---

## Summary

| Action  | Count | Files                                              |
|---------|-------|----------------------------------------------------|
| CREATE  | 17    | 2 constants, 4 model, 4 store, 5 detector, 2 cmd   |
| MODIFY  | 4     | `store/store.go`, `cmd/scan.go`, `cmd/scanoutput.go`, `cmd/root.go` |
| TOTAL   | 21    |                                                    |

---

## Dependency Graph

```
Phase 1 (constants + models)
    │
    ├──→ Phase 2 (store)
    │        │
    │        ├──→ Phase 4 (scan integration)
    │        │
    │        └──→ Phase 5 (query commands)
    │
    └──→ Phase 3 (detector)
             │
             └──→ Phase 4 (scan integration)
```
