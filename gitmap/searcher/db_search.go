package searcher

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/lazyregex"
)

// SearchRepoDB searches the RepoDB for exact match, with analytical caching.
func SearchRepoDB(ctx context.Context, db *sql.DB, query string, useCache bool) ([]SearchResult, error) {
	if useCache {
		var cachedJson string
		err := db.QueryRowContext(ctx, "SELECT ResultJson FROM SearchCache WHERE Query = ?", query).Scan(&cachedJson)
		if err == nil && cachedJson != "" {
			// Found in cache
			var res []SearchResult
			if err := json.Unmarshal([]byte(cachedJson), &res); err == nil {
				// Bump hit count
				db.ExecContext(ctx, "UPDATE SearchCache SET Hits = Hits + 1 WHERE Query = ?", query)
				return res, nil
			}
		}
	}

	// Not in cache or cache disabled, perform full search
	rows, err := db.QueryContext(ctx, "SELECT RelativePath, AbsolutePath, Content FROM RepoFile WHERE IsBig = 0 AND Content LIKE ?", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allResults []SearchResult
	for rows.Next() {
		var rel, abs, content string
		if err := rows.Scan(&rel, &abs, &content); err != nil {
			continue
		}
		fileResults := SearchExact(content, query, abs, rel)
		allResults = append(allResults, fileResults...)
	}

	// Update cache
	if useCache {
		updateCache(ctx, db, query, allResults)
	}

	return allResults, nil
}

// SearchRepoDBRegex searches the RepoDB using a regex pattern.
func SearchRepoDBRegex(ctx context.Context, db *sql.DB, expr string, useCache bool) ([]SearchResult, error) {
	if useCache {
		var cachedJson string
		err := db.QueryRowContext(ctx, "SELECT ResultJson FROM SearchCache WHERE Query = ?", "regex:"+expr).Scan(&cachedJson)
		if err == nil && cachedJson != "" {
			var res []SearchResult
			if err := json.Unmarshal([]byte(cachedJson), &res); err == nil {
				db.ExecContext(ctx, "UPDATE SearchCache SET Hits = Hits + 1 WHERE Query = ?", "regex:"+expr)
				return res, nil
			}
		}
	}

	lz := lazyregex.New(expr)

	rows, err := db.QueryContext(ctx, "SELECT RelativePath, AbsolutePath, Content FROM RepoFile WHERE IsBig = 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allResults []SearchResult
	for rows.Next() {
		var rel, abs, content string
		if err := rows.Scan(&rel, &abs, &content); err != nil {
			continue
		}
		fileResults := SearchRegex(content, lz, abs, rel)
		allResults = append(allResults, fileResults...)
	}

	if useCache {
		updateCache(ctx, db, "regex:"+expr, allResults)
	}

	return allResults, nil
}

func updateCache(ctx context.Context, db *sql.DB, cacheKey string, results []SearchResult) {
	// A naive caching implementation: always insert/update
	b, err := json.Marshal(results)
	if err != nil {
		return
	}
	
	query := `
		INSERT INTO SearchCache (Query, Hits, ResultJson, CreatedAt, UpdatedAt)
		VALUES (?, 1, ?, 0, 0)
		ON CONFLICT(Query) DO UPDATE SET ResultJson=excluded.ResultJson, Hits=Hits+1;
	`
	db.ExecContext(ctx, query, cacheKey, string(b))
}
