package db

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type NodePathAlias struct {
	NodePathAliasId int64
	NodeId          string
	Alias           string
	AbsolutePath    string
}

func UpsertPathAlias(db *sql.DB, nodeId, alias, absPath string) error { return nil }
func GetPathAlias(db *sql.DB, nodeId, alias string) (string, error)   { return "", nil }
func SetDefaultPath(db *sql.DB, nodeId, path string) error            { return nil }
func GetDefaultPath(db *sql.DB, nodeId string) (string, error)        { return "", nil }
func ListPathAliases(db *sql.DB, nodeId string) result.ResultSlice[NodePathAlias] {
	return result.OkSlice([]NodePathAlias{})
}
