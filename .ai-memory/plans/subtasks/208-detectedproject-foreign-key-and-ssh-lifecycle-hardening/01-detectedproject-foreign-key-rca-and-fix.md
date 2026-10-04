# Subtask 78.1: DetectedProject FOREIGN KEY Constraint (787) In-Depth RCA & Fix

## 1. Context & Objective
In multi-repo scans across 37 repositories with 53 detected projects (15 Go, 31 React, 7 Node), 45 projects upserted successfully while exactly 8 failed with:
`[QueryWrapper Error]: exec failed: constraint failed: FOREIGN KEY constraint failed (787)`
`query: INSERT INTO DetectedProject (RepoId, ProjectTypeId, ProjectName, AbsolutePath, RepoPath, RelativePath, PrimaryIndicator) VALUES (?, ?, ?, ?, ?, ?, ?)`

The objective is to identify why these 8 failed, ensure complete ProjectType seeding, fix `RepoId` resolution across nested directories, and prevent any foreign key failures during scan.

---

## 2. Technical Investigation
- Examine `cli/store/project.go` and `cli/cmdscan/scanprojects.go`.
- Check `ProjectType` table contents and seeding in `cli/constants/constants_project_sql.go`.
- Check how `ProjectTypeId` is resolved for Node, React, and Go projects.
- Verify how `RepoId` is resolved when `rec.ID <= 0` or when project path is inside a nested repo.

---

## 3. Remediation Checklist
- [ ] Ensure all project types (Go=1, React=2, Node=3, Python=4, Rust=5) are always present in the database.
- [ ] In `scanprojects.go`, if `RepoId <= 0`, dynamically query `Repo` table matching the closest parent path.
- [ ] If no parent repository is found in database, fail safely or log diagnostic warning rather than executing an invalid SQL insert.
- [ ] Run full project scan across `./` or local repos and verify 0 foreign key errors.
