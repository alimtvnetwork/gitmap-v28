package cmdsearch

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	searchconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdsearch/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/repodb"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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

		return nil, nil, apperror.WrapSimple(err, fmt.Sprintf(searchconstants.FmtErrFindRepoByPath, cwd))
	}
	if len(repos) == 0 {
		mainDB.Close()

		return nil, nil, apperror.NewNotFoundError(searchconstants.ErrRepoNotFound)
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
		arg := args[i]
		if isLimitFlag(arg) && i+1 < len(args) {
			limit, _ = strconv.Atoi(args[i+1])
			i++
			continue
		}

		if isLimitFlag(arg) {
			continue
		}

		if strings.HasPrefix(arg, searchconstants.FlagArgLimitPrefix) {
			limit, _ = strconv.Atoi(strings.TrimPrefix(arg, searchconstants.FlagArgLimitPrefix))
			continue
		}

		cleanArgs = append(cleanArgs, arg)
	}

	return limit, cleanArgs
}

func isLimitFlag(arg string) bool {
	return arg == searchconstants.FlagArgLimit || arg == searchconstants.FlagArgLimitShort
}
