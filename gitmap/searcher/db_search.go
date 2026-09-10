package searcher

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/lazyregex"
)

const (
	sqlSearchExact = "SELECT RelativePath, AbsolutePath, Content FROM RepoFile WHERE IsBig = 0 AND Content LIKE ?"
	sqlSearchRegex = "SELECT RelativePath, AbsolutePath, Content FROM RepoFile WHERE IsBig = 0"
	sqlUpsertCache = `
		INSERT INTO SearchCache (Query, Hits, ResultJson, CreatedAt, UpdatedAt)
		VALUES (?, 1, ?, 0, 0)
		ON CONFLICT(Query) DO UPDATE SET ResultJson=excluded.ResultJson, Hits=Hits+1;`
)

func scanRow(rows *sql.Rows) (string, string, string, error) {
	var rel, abs, content string
	err := rows.Scan(&rel, &abs, &content)
	return rel, abs, content, err
}

func collectExactMatches(rows *sql.Rows, query string) ([]SearchResult, error) {
	var allResults []SearchResult
	for rows.Next() {
		rel, abs, content, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, SearchExact(content, query, abs, rel)...)
	}
	return allResults, rows.Err()
}

func collectRegexMatches(rows *sql.Rows, lz *lazyregex.LazyRegexp) ([]SearchResult, error) {
	var allResults []SearchResult
	for rows.Next() {
		rel, abs, content, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, SearchRegex(content, lz, abs, rel)...)
	}
	return allResults, rows.Err()
}

func applyLimit(results []SearchResult, limit int) []SearchResult {
	if limit <= 0 {
		return results
	}
	if len(results) <= limit {
		return results
	}
	return results[:limit]
}

func executeSearchExact(ctx context.Context, db *sql.DB, query string) ([]SearchResult, error) {
	rows, err := db.QueryContext(ctx, sqlSearchExact, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectExactMatches(rows, query)
}

// SearchRepoDB searches the RepoDB for exact match, with analytical caching.
func SearchRepoDB(ctx context.Context, db *sql.DB, query string, limit int, useCache bool) ([]SearchResult, error) {
	if res, hasCache := maybeGetCachedSearchResults(ctx, db, query, limit, useCache); hasCache {
		return res, nil
	}
	results, err := executeSearchExact(ctx, db, query)
	if err != nil {
		return nil, err
	}
	if useCache {
		updateCache(ctx, db, query, results)
	}
	return applyLimit(results, limit), nil
}

func executeSearchRegex(ctx context.Context, db *sql.DB, lz *lazyregex.LazyRegexp) ([]SearchResult, error) {
	rows, err := db.QueryContext(ctx, sqlSearchRegex)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectRegexMatches(rows, lz)
}

// SearchRepoDBRegex searches the RepoDB using a regex pattern.
func SearchRepoDBRegex(ctx context.Context, db *sql.DB, expr string, limit int, useCache bool) ([]SearchResult, error) {
	cacheKey := "regex:" + expr
	if res, hasCache := maybeGetCachedSearchResults(ctx, db, cacheKey, limit, useCache); hasCache {
		return res, nil
	}
	lz := lazyregex.New(expr)
	results, err := executeSearchRegex(ctx, db, lz)
	if err != nil {
		return nil, err
	}
	if useCache {
		updateCache(ctx, db, cacheKey, results)
	}
	return applyLimit(results, limit), nil
}

func updateCache(ctx context.Context, db *sql.DB, cacheKey string, results []SearchResult) {
	b, err := json.Marshal(results)
	if err != nil {
		return
	}
	if _, execErr := db.ExecContext(ctx, sqlUpsertCache, cacheKey, string(b)); execErr != nil {
		return
	}
}

func queryCachedJson(ctx context.Context, db *sql.DB, query string) (string, bool) {
	var cachedJson string
	err := db.QueryRowContext(ctx, "SELECT ResultJson FROM SearchCache WHERE Query = ?", query).Scan(&cachedJson)
	if err != nil {
		return "", false
	}
	return cachedJson, cachedJson != ""
}

func getCachedSearchResults(ctx context.Context, db *sql.DB, query string, limit int) ([]SearchResult, bool) {
	cachedJson, hasJson := queryCachedJson(ctx, db, query)
	if !hasJson {
		return nil, false
	}
	var res []SearchResult
	if err := json.Unmarshal([]byte(cachedJson), &res); err != nil {
		return nil, false
	}
	if _, err := db.ExecContext(ctx, "UPDATE SearchCache SET Hits = Hits + 1 WHERE Query = ?", query); err != nil {
		return nil, false
	}
	return applyLimit(res, limit), true
}

func maybeGetCachedSearchResults(ctx context.Context, db *sql.DB, query string, limit int, useCache bool) ([]SearchResult, bool) {
	if !useCache {
		return nil, false
	}
	return getCachedSearchResults(ctx, db, query, limit)
}
