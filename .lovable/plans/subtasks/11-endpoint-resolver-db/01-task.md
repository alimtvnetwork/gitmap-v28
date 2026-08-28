# Subtask 1: DB & SQL Constants Update

1. Open `gitmap/constants/constants_store.go`.
2. Add `SQLSelectRepoByID = "SELECT RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, Notes, IdentifiedTransport FROM Repo WHERE RepoId = ?"`.
3. Open `gitmap/store/repo.go`.
4. Add `func (db *DB) FindByID(id int64) ([]model.ScanRecord, error) { ... }` executing `SQLSelectRepoByID`.
5. Ensure `store` package compiles without errors.
