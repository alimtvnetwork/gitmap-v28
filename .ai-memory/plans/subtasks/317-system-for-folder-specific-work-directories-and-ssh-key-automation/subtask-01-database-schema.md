# Subtask 01: Database & Schema Implementation (SchemaVersion = 34)

## Status: READY
- **Subtask Slug**: `subtask-01-database-schema`
- **Parent Plan**: `.ai-memory/plans/317-system-for-folder-specific-work-directories-and-ssh-key-automation.md`
- **Architecture Spec**: `02-spec/21-app/317-system-for-folder-specific-work-directories-and-ssh-key-automation/01-database-schema-spec.md`
- **Owner**: Database & Storage Subagent
- **Target Subsystem**: `cli/store`, `cli/model`, `cli/constants`
- **Schema Target**: `SchemaVersionCurrent = 34`
- **Exclusive Owned Files**:
  - `cli/constants/constants_settings.go`
  - `cli/constants/constants_workdir_ssh.go` (new)
  - `cli/model/workdir.go`
  - `cli/model/folder_workdir.go` (new)
  - `cli/model/folder_group.go` (new)
  - `cli/model/ssh_association.go` (new)
  - `cli/model/effective_credentials.go` (new)
  - `cli/store/migrations_v34.go` (new)
  - `cli/store/store.go`
  - `cli/store/workdir_schema.go`
  - `cli/store/workdir_list.go`
  - `cli/store/folder_workdir_repo.go` (new)
  - `cli/store/folder_group_repo.go` (new)
  - `cli/store/ssh_association_repo.go` (new)

---

## 1. Context & Objectives

GitMap currently manages workspaces via a single flat `work_directories` table with a global default flag. It does not provide automated directory-to-workspace mapping, workgroup categorization, or path-bound SSH credential injection. Furthermore, the legacy `work_directories` table uses snake_case naming conventions (`id`, `absolute_path`, `is_default`), violating GitMap's PascalCase standards.

This subtask implements the database layer, schema migration, Go entity models, and store repositories for **Task 317**:
1. **Schema Version Advancement:** Advance `SchemaVersionCurrent` from `33` to `34` in `cli/constants/constants_settings.go`.
2. **PascalCase Table Deployment:** Deploy four PascalCase tables in `~/.gitmap/gitmap.db`:
   - `WorkDirectory` (Master workspace registry)
   - `FolderWorkDir` (Folder path to WorkDirectory binding)
   - `FolderGroup` (Folder grouping and categorization)
   - `SshKeyAssociation` (Automated SSH credential routing and author overrides)
3. **100% Zero-Breakage Backward Compatibility:**
   - Migrate existing rows from legacy `work_directories` table into `WorkDirectory` with zero data loss.
   - Deploy `work_directories` as a compatibility view over `WorkDirectory`.
   - Preserve existing callsites using `model.WorkDir` and `store.ListWorkDirs()`.
4. **Hierarchical Path & Multi-Stage Credential Resolution:** Implement single-query SQLite resolution logic supporting Windows (`\`) and Unix (`/`) path separators with `COLLATE NOCASE` collation.
5. **Strict Coding Standard Compliance:**
   - All booleans start with `is` or `has` only (no negative booleans).
   - Functions strictly <= 15 lines max (<= 8 lines preferred).
   - Zero nested `if` statements (guard clauses only).
   - Blank line after `}` and blank line before `return`.
   - Leaf function error return sovereignty (no exit handler calls).

---

## 2. Target Files Manifest

| Target File | Purpose / Modifications | Action |
| :--- | :--- | :--- |
| `cli/constants/constants_settings.go` | Bump `SchemaVersionCurrent` to `34` | Edit |
| `cli/constants/constants_workdir_ssh.go` | Define SQL DDL, indexes, and CRUD query constants | Create |
| `cli/model/workdir.go` | Add/align `WorkDirectory` model and aliases | Edit |
| `cli/model/folder_workdir.go` | Define `FolderWorkDir` entity struct | Create |
| `cli/model/folder_group.go` | Define `FolderGroup` entity struct | Create |
| `cli/model/ssh_association.go` | Define `SshKeyAssociation` entity struct | Create |
| `cli/model/effective_credentials.go` | Define `EffectiveFolderCredentials` resolution payload | Create |
| `cli/store/migrations_v34.go` | Implement `Migration_AddFolderWorkDirAndSshKeyAutomation` | Create |
| `cli/store/store.go` | Wire v34 migration into `Migrate()` and drop cascade in `Reset()` | Edit |
| `cli/store/workdir_schema.go` | Update constants to reference PascalCase `WorkDirectory` | Edit |
| `cli/store/workdir_list.go` | Refactor to query `WorkDirectory` table directly | Edit |
| `cli/store/folder_workdir_repo.go` | Repository methods for `FolderWorkDir` bindings & ancestor lookup | Create |
| `cli/store/folder_group_repo.go` | Repository methods for `FolderGroup` clustering | Create |
| `cli/store/ssh_association_repo.go` | Repository methods for `SshKeyAssociation` and 5-stage cascading lookup | Create |

---

## 3. Step-by-Step Atomic Checklist

### Phase 1: Constants & SQL DDL Definitions
- [ ] **1.1 Schema Version Bump:** Edit `cli/constants/constants_settings.go` to set `const SchemaVersionCurrent = 34`.
- [ ] **1.2 Create DDL Constants in `cli/constants/constants_workdir_ssh.go`:**
  - Define `SQLCreateWorkDirectoryTable` (`WorkDirectoryId`, `AbsolutePath`, `Label`, `IsDefault`, `HasActiveBinding`, timestamps).
  - Define `SQLCreateWorkDirectoriesCompatView` (`work_directories` view mapping to `WorkDirectory`).
  - Define `SQLCreateFolderWorkDirTable` (`FolderWorkDirId`, `FolderPath`, `WorkDirectoryId`, `IsRecursive`, `Priority`, `IsActive`, timestamps).
  - Define `SQLCreateFolderGroupTable` (`FolderGroupId`, `GroupName`, `FolderPath`, `Description`, `IsDefault`, `IsActive`, timestamps).
  - Define `SQLCreateSshKeyAssociationTable` (`SshKeyAssociationId`, `SshKeyId`, `AssociationType`, `TargetIdentifier`, `FolderPath`, `GroupName`, `WorkDirectoryId`, `RepoId`, `GitAuthorName`, `GitAuthorEmail`, `IsActive`, `IsAutoInjected`, `Priority`, timestamps).
- [ ] **1.3 Define Secondary Indexes & Triggers:**
  - `IdxWorkDirectory_AbsolutePath` (UNIQUE, `COLLATE NOCASE`).
  - `IdxFolderWorkDir_FolderPath` (UNIQUE, `COLLATE NOCASE`).
  - `IdxFolderGroup_NamePath` (UNIQUE, composite `GroupName, FolderPath`).
  - `IdxSshKeyAssoc_TargetKey` (UNIQUE, composite `AssociationType, TargetIdentifier, SshKeyId`).
- [ ] **1.4 Parameterized Query Constants:**
  - Ancestor `FolderWorkDir` lookup with `COLLATE NOCASE` and prefix pattern matching (`LIKE %/%` and `LIKE %\%\`).
  - 5-stage cascading `SshKeyAssociation` lookup query.

### Phase 2: Domain Go Entity Models
- [ ] **2.1 Modernize `cli/model/workdir.go`:**
  - Add `WorkDirectory` struct with `WorkDirectoryId`, `AbsolutePath`, `Label`, `IsDefault`, `HasActiveBinding`, `CreatedAt`, `UpdatedAt`.
  - Maintain `WorkDir` struct for backward compatibility with existing callsites.
- [ ] **2.2 Create `cli/model/folder_workdir.go`:**
  - Define `FolderWorkDir` struct with `FolderWorkDirId`, `FolderPath`, `WorkDirectoryId`, `IsRecursive`, `Priority`, `IsActive`, `CreatedAt`, `UpdatedAt`.
- [ ] **2.3 Create `cli/model/folder_group.go`:**
  - Define `FolderGroup` struct with `FolderGroupId`, `GroupName`, `FolderPath`, `Description`, `IsDefault`, `IsActive`, `CreatedAt`, `UpdatedAt`.
- [ ] **2.4 Create `cli/model/ssh_association.go`:**
  - Define `SshKeyAssociation` struct with `SshKeyAssociationId`, `SshKeyId`, `AssociationType`, `TargetIdentifier`, `FolderPath`, `GroupName`, `WorkDirectoryId`, `RepoId`, `GitAuthorName`, `GitAuthorEmail`, `IsActive`, `IsAutoInjected`, `Priority`, `CreatedAt`, `UpdatedAt`.
- [ ] **2.5 Create `cli/model/effective_credentials.go`:**
  - Define `EffectiveFolderCredentials` containing resolved SSH key, key path, author overrides, and source association type.

### Phase 3: Migration Engine v34 (`cli/store/migrations_v34.go`)
- [ ] **3.1 Check Version Gate:**
  - Guard: If `Setting` table does not exist, return `nil`.
  - Guard: If `schema_version` is already >= 34, return `nil`.
- [ ] **3.2 Execute Table Creation:**
  - Create `WorkDirectory` table and indexes.
- [ ] **3.3 Execute Data Backfill:**
  - Detect if legacy `work_directories` table exists.
  - If exists, execute `INSERT OR IGNORE INTO WorkDirectory ... SELECT ... FROM work_directories`.
- [ ] **3.4 Drop Legacy View / Rebuild View:**
  - If `work_directories` was a table, rename or drop and create `work_directories` as a view pointing to `WorkDirectory`.
  - If `WorkDirectory` was previously a view, drop it before creating the table.
- [ ] **3.5 Deploy Child Tables:**
  - Create `FolderWorkDir`, `FolderGroup`, and `SshKeyAssociation` tables and indexes.
- [ ] **3.6 Bump Setting Version:**
  - Execute `INSERT INTO Setting (Key, Value) VALUES ('schema_version', '34') ON CONFLICT(Key) DO UPDATE SET Value='34'`.
- [ ] **3.7 Store Integration:**
  - Call `Migration_AddFolderWorkDirAndSshKeyAutomation(db.conn)` inside `db.Migrate()` in `cli/store/store.go`.
  - Add child table drop statements to `db.Reset()` in `cli/store/store.go`.

### Phase 4: Store Repositories Implementation
- [ ] **4.1 Modernize `cli/store/workdir_list.go`:**
  - Point `EnsureWorkDirsTable`, `EnsureWorkDir`, `ListWorkDirsWrap`, and `GetWorkDirByPath` to the `WorkDirectory` table.
- [ ] **4.2 Implement `cli/store/folder_workdir_repo.go`:**
  - `BindFolderWorkDir(folderPath string, workDirId int64, isRecursive bool, priority int) (*model.FolderWorkDir, error)`
  - `UnbindFolderWorkDir(folderPath string) error`
  - `ListFolderWorkDirs() ([]model.FolderWorkDir, error)`
  - `ResolveWorkDirectoryForPath(targetPath string) (*model.WorkDirectory, error)`
- [ ] **4.3 Implement `cli/store/folder_group_repo.go`:**
  - `AddFolderToGroup(groupName, folderPath, description string) (*model.FolderGroup, error)`
  - `RemoveFolderFromGroup(groupName, folderPath string) error`
  - `ListFoldersInGroup(groupName string) ([]model.FolderGroup, error)`
  - `ListActiveGroups() ([]model.FolderGroupSummary, error)`
- [ ] **4.4 Implement `cli/store/ssh_association_repo.go`:**
  - `AssociateSSHKeyWithFolder(sshKeyId int64, folderPath, authorName, authorEmail string) (*model.SshKeyAssociation, error)`
  - `AssociateSSHKeyWithGroup(sshKeyId int64, groupName, authorName, authorEmail string) (*model.SshKeyAssociation, error)`
  - `AssociateSSHKeyWithWorkDir(sshKeyId int64, workDirId int64) (*model.SshKeyAssociation, error)`
  - `AssociateSSHKeyWithRepo(sshKeyId int64, repoId int64) (*model.SshKeyAssociation, error)`
  - `RemoveSSHKeyAssociation(assocId int64) error`
  - `ResolveEffectiveCredentialsForPath(repoPath string, repoId int64) (*model.EffectiveFolderCredentials, error)`

### Phase 5: Verification & Zero-Breakage Invariants
- [ ] **5.1 SQLite DDL Verification:**
  - Verify all 4 tables exist with PascalCase column naming.
- [ ] **5.2 Backward-Compatibility Verification:**
  - Verify `SELECT * FROM work_directories` succeeds and mirrors `WorkDirectory`.
- [ ] **5.3 Windows Collation Verification:**
  - Verify case-insensitive path lookups (`C:\Work\Project` matches `c:\work\project`).
- [ ] **5.4 Cascade Deletion Verification:**
  - Verify deleting a `WorkDirectory` cascades to child `FolderWorkDir` and sets null in `SshKeyAssociation`.

---

## 4. Coding Standard Enforcement Checklist

- [ ] **Positive Booleans Only:**
  - All variables, fields, and parameters start with `is` or `has` (`isRecursive`, `isActive`, `isDefault`, `hasActiveBinding`, `isAutoInjected`, `hasCustomKey`).
  - Total ban on negative boolean names (`isNotActive`, `hasNoKey`).
- [ ] **Function Length Discipline:**
  - Every function strictly <= 15 lines max (<= 8 lines preferred).
  - Scanner logic and row mapping split into discrete helpers (`scanFolderWorkDirRows`, `scanSshAssocRows`).
- [ ] **Zero Nested Ifs:**
  - Flattened branching with early guard clauses and returns.
- [ ] **Return Styling:**
  - Exactly one blank line after `}` and one blank line before `return`.
- [ ] **Error Sovereignty:**
  - Errors wrapped via `apperror.WrapSimple(err, ...)` or `apperror.Wrap(err, ...)`.
  - Zero calls to `os.Exit` or `cliexit.HandleError` in store functions.

---

## 5. Acceptance Criteria

1. `SchemaVersionCurrent` equals `34`.
2. Existing databases cleanly upgrade upon CLI command execution without manual intervention or data loss.
3. Legacy `work_directories` reads continue to succeed via the compatibility view.
4. Calling `ResolveWorkDirectoryForPath` correctly matches exact folder paths as well as deep descendant child paths.
5. Calling `ResolveEffectiveCredentialsForPath` resolves the expected SSH key and custom author attributes following the 5-stage cascading priority order.
6. Zero compilation warnings, zero build regressions.
