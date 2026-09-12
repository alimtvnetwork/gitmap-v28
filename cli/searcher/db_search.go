package searcher

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
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

func collectExactMatches(rows *sql.Rows, query string) ([]SearchResult, *apperror.AppError) {
	var allResults []SearchResult
	for rows.Next() {
		rel, abs, content, err := scanRow(rows)
		if err != nil {
			return nil, apperror.WrapSimple(err, "searcher.collectExactMatches.Scan")
		}

		allResults = append(allResults, SearchExact(content, query, abs, rel)...)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "searcher.collectExactMatches.Rows")
	}

	return allResults, nil
}

func collectRegexMatches(rows *sql.Rows, lz *lazyregex.LazyRegexp) ([]SearchResult, *apperror.AppError) {
	var allResults []SearchResult
	for rows.Next() {
		rel, abs, content, err := scanRow(rows)
		if err != nil {
			return nil, apperror.WrapSimple(err, "searcher.collectRegexMatches.Scan")
		}

		allResults = append(allResults, SearchRegex(content, lz, abs, rel)...)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "searcher.collectRegexMatches.Rows")
	}

	return allResults, nil
}

func applyLimit(results []SearchResult, limit int) []SearchResult {
	if limit <= 0 {
		return results
	}

	hasFit := len(results) <= limit
	if hasFit {
		return results
	}

	return results[:limit]
}

func executeSearchExact(ctx context.Context, db *sql.DB, query string) ([]SearchResult, *apperror.AppError) {
	rows, err := db.QueryContext(ctx, sqlSearchExact, "%"+query+"%")
	if err != nil {
		return nil, apperror.WrapSimple(err, "searcher.executeSearchExact.Query")
	}

	defer rows.Close()

	return collectExactMatches(rows, query)
}

// SearchRepoDB searches the RepoDB for exact match, with analytical caching.
func SearchRepoDB(ctx context.Context, db *sql.DB, query string, limit int, useCache bool) ([]SearchResult, *apperror.AppError) {
	res, hasCache, appErr := maybeGetCachedSearchResults(ctx, db, query, limit, useCache)
	if appErr != nil {
		return nil, appErr
	}

	if hasCache {
		return res, nil
	}

	results, err := executeSearchExact(ctx, db, query)
	if err != nil {
		return nil, err
	}

	maybeUpdateCache(ctx, db, query, results, useCache)

	return applyLimit(results, limit), nil
}

func executeSearchRegex(ctx context.Context, db *sql.DB, lz *lazyregex.LazyRegexp) ([]SearchResult, *apperror.AppError) {
	rows, err := db.QueryContext(ctx, sqlSearchRegex)
	if err != nil {
		return nil, apperror.WrapSimple(err, "searcher.executeSearchRegex.Query")
	}

	defer rows.Close()

	return collectRegexMatches(rows, lz)
}

// SearchRepoDBRegex searches the RepoDB using a regex pattern.
func SearchRepoDBRegex(ctx context.Context, db *sql.DB, expr string, limit int, useCache bool) ([]SearchResult, *apperror.AppError) {
	cacheKey := "regex:" + expr
	res, hasCache, appErr := maybeGetCachedSearchResults(ctx, db, cacheKey, limit, useCache)
	if appErr != nil {
		return nil, appErr
	}

	if hasCache {
		return res, nil
	}

	lz := lazyregex.New(expr)
	results, err := executeSearchRegex(ctx, db, lz)
	if err != nil {
		return nil, err
	}

	maybeUpdateCache(ctx, db, cacheKey, results, useCache)

	return applyLimit(results, limit), nil
}

func maybeUpdateCache(ctx context.Context, db *sql.DB, cacheKey string, results []SearchResult, useCache bool) {
	if !useCache {
		return
	}

	if appErr := updateCache(ctx, db, cacheKey, results); appErr != nil {
		appErr.HandleError()
	}
}

func updateCache(ctx context.Context, db *sql.DB, cacheKey string, results []SearchResult) *apperror.AppError {
	b, err := json.Marshal(results)
	if err != nil {
		return apperror.WrapSimple(err, "searcher.updateCache.Marshal")
	}

	_, execErr := db.ExecContext(ctx, sqlUpsertCache, cacheKey, string(b))
	if execErr != nil {
		return apperror.WrapSimple(execErr, "searcher.updateCache.Exec")
	}

	return nil
}

func incrementCacheHits(ctx context.Context, db *sql.DB, key string) *apperror.AppError {
	_, err := db.ExecContext(ctx, "UPDATE SearchCache SET Hits = Hits + 1 WHERE Query = ?", key)
	if err != nil {
		return apperror.WrapSimple(err, "searcher.incrementCacheHits.Exec")
	}

	return nil
}

func queryCachedJson(ctx context.Context, db *sql.DB, query string) (string, bool, *apperror.AppError) {
	var cachedJson string
	err := db.QueryRowContext(ctx, "SELECT ResultJson FROM SearchCache WHERE Query = ?", query).Scan(&cachedJson)
	if err == sql.ErrNoRows {
		return "", false, nil
	}

	if err != nil {
		return "", false, apperror.WrapSimple(err, "searcher.queryCachedJson.Scan")
	}

	hasJson := cachedJson != ""

	return cachedJson, hasJson, nil
}

func getCachedSearchResults(ctx context.Context, db *sql.DB, query string, limit int) ([]SearchResult, bool, *apperror.AppError) {
	cachedJson, hasJson, appErr := queryCachedJson(ctx, db, query)
	if appErr != nil {
		return nil, false, appErr
	}

	if !hasJson {
		return nil, false, nil
	}

	var res []SearchResult
	if err := json.Unmarshal([]byte(cachedJson), &res); err != nil {
		return nil, false, apperror.WrapSimple(err, "searcher.getCachedSearchResults.Unmarshal")
	}

	if hitErr := incrementCacheHits(ctx, db, query); hitErr != nil {
		hitErr.HandleError()
	}

	return applyLimit(res, limit), true, nil
}

func maybeGetCachedSearchResults(
	ctx context.Context,
	db *sql.DB,
	query string,
	limit int,
	useCache bool,
) ([]SearchResult, bool, *apperror.AppError) {
	if !useCache {
		return nil, false, nil
	}

	return getCachedSearchResults(ctx, db, query, limit)
}
