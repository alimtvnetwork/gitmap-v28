package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func updateRepoInDB(db *store.DB, repoID int64, newPath, newName string) error {
	query := "UPDATE Repo SET AbsolutePath = ?, RepoName = ? WHERE Id = ?"
	_, err := store.ExecWrapper(db.Conn(), query, newPath, newName, repoID).Destruct()
	if err != nil {
		return err
	}
	aliasQuery := "UPDATE Alias SET AbsolutePath = ? WHERE RepoID = ?"
	_, _ = store.ExecWrapper(db.Conn(), aliasQuery, newPath, repoID).Destruct()
	return nil
}
