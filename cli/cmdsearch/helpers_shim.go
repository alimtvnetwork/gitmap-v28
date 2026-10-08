package cmdsearch

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/repodb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"context"
	"fmt"
	"os"
	"database/sql"
	"strconv"
	"strings"
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
	if err != nil {
		mainDB.Close()

		return nil, nil, fmt.Errorf("find repo by path %s failed: %w", cwd, err)
	}
	if len(repos) == 0 {
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

func parseLimit(args []string) (int, []string) {
	limit := 0
	var cleanArgs []string
	for i := 0; i < len(args); i++ {
		if (args[i] == "--limit" || args[i] == "-l") && i+1 < len(args) {
			limit, _ = strconv.Atoi(args[i+1])
			i++
			continue
		}

		if args[i] == "--limit" || args[i] == "-l" {
			continue
		}

		if strings.HasPrefix(args[i], "--limit=") {
			limit, _ = strconv.Atoi(strings.TrimPrefix(args[i], "--limit="))
			continue
		}

		cleanArgs = append(cleanArgs, args[i])
	}

	return limit, cleanArgs
}
