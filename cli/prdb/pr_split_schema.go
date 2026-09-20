package prdb

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func prSchemaStatements() []string {
	return []string{
		constants.SQLCreatePullRequest,
		constants.SQLCreatePullRequestIndexes,
		constants.SQLCreatePrRelease,
		constants.SQLCreatePrReleaseIndexes,
		constants.SQLCreatePrBranch,
		constants.SQLCreatePrBranchIndexes,
		constants.SQLCreatePrViews,
	}
}

func execPrSchemaStatements(conn *sql.DB, stmts []string) *apperror.AppError {
	for _, stmt := range stmts {
		if _, err := conn.Exec(stmt); err != nil {
			return apperror.WrapSimple(err, "prdb.execPrSchemaStatements")
		}
	}

	return nil
}

// InitPrSchema initializes the PullRequest, PrRelease, and PrBranch tables, indexes, and views.
func InitPrSchema(conn *sql.DB) *apperror.AppError {
	return execPrSchemaStatements(conn, prSchemaStatements())
}
