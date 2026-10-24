package cmddb

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

func handleReadInspect(flags readFlags, dbArg string) *appfault.AppError {
	start := time.Now()
	conn, dbPath, err := openReadOnlyResolvedDb(dbArg)
	if err != nil {
		return failRead(os.Stdout, flags.isJson, "db-read:inspect", dbArg, dbArg, "ERR_DB_NOT_FOUND", 20, err.Error(), "Verify database alias or run 'gitmap scan' to initialize.")
	}
	defer conn.Close()

	tables, colMap, inspectErr := fetchOrCacheMetadata(conn, dbPath)
	elapsed := time.Since(start)
	if inspectErr != nil {
		return failRead(os.Stdout, flags.isJson, "db-read:inspect", dbPath, dbArg, "ERR_EXECUTION", 23, inspectErr.Error(), "Check database schema integrity.")
	}

	return outputInspectResult(flags, dbPath, dbArg, tables, colMap, elapsed)
}

func outputInspectResult(flags readFlags, dbPath, dbArg string, tables []TableInfo, colMap map[string][]ColumnInfo, elapsed time.Duration) *appfault.AppError {
	if flags.isJson {
		data := map[string]any{"tables": tables, "columns": colMap, "totalTables": len(tables)}
		outputReadSuccess(os.Stdout, "db-read:inspect", dbPath, dbArg, data, elapsed)

		return nil
	}

	renderInspectTerminal(os.Stdout, dbPath, tables, colMap, flags.isShowColumns)
	fmt.Println(formatCompletionTime(elapsed))

	return nil
}

func fetchOrCacheMetadata(conn *sql.DB, dbPath string) ([]TableInfo, map[string][]ColumnInfo, error) {
	cache := GetGlobalTableMetadataCache()
	tables, colMap, hasCached := cache.Get(dbPath)
	if hasCached {
		return tables, colMap, nil
	}

	tables, err := listTablesAndViews(conn)
	if err != nil {
		return nil, nil, err
	}

	colMap = populateColumnsMap(conn, tables)
	cache.Set(dbPath, tables, colMap)

	return tables, colMap, nil
}

func populateColumnsMap(conn *sql.DB, tables []TableInfo) map[string][]ColumnInfo {
	colMap := make(map[string][]ColumnInfo, len(tables))
	for _, t := range tables {
		cols, err := fetchColumnInfos(conn, t.Name)
		if err == nil {
			colMap[t.Name] = cols
		}
	}

	return colMap
}

func renderInspectTerminal(out io.Writer, dbPath string, tables []TableInfo, colMap map[string][]ColumnInfo, isShowColumns bool) {
	renderTablesTerminal(out, dbPath, tables)
	if !isShowColumns {
		return
	}

	for _, t := range tables {
		cols := colMap[t.Name]
		if len(cols) > 0 {
			fmt.Println()
			renderColumnsTerminal(out, dbPath, t.Name, cols)
		}
	}
}
