package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/repodb"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func getRepoDB(ctx context.Context) (*store.DB, *sql.DB, error) {
	mainDB, err := store.OpenDefault()
	if err != nil {
		return nil, nil, err
	}
	
	cwd, err := os.Getwd()
	if err != nil {
		mainDB.Close()
		return nil, nil, err
	}
	
	repos, err := mainDB.FindByPath(cwd)
	if err != nil || len(repos) == 0 {
		mainDB.Close()
		return nil, nil, fmt.Errorf("current directory is not a tracked gitmap repository. run 'gitmap scan' first")
	}
	
	repoDB, err := repodb.OpenRepoDB(ctx, constants.DefaultOutputDir, repos[0].AbsolutePath, repos[0].ID)
	if err != nil {
		mainDB.Close()
		return nil, nil, err
	}
	
	return mainDB, repoDB, nil
}
