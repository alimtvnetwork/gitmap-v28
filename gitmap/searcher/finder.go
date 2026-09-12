package searcher

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/lazyregex"
)

type FileFindResult struct {
	RelativePath string `json:"relative_path"`
	AbsolutePath string `json:"absolute_path"`
}

// FindFile searches RepoDB for file names matching the query or pattern.
func FindFile(
	ctx context.Context,
	db *sql.DB,
	query string,
	limit int,
	useCache bool,
) ([]FileFindResult, *apperror.AppError) {
	res, hasCache, appErr := maybeGetCachedFindResults(ctx, db, "find:"+query, limit, useCache)
	if appErr != nil {
		return nil, appErr
	}

	if hasCache {
		return res, nil
	}

	rows, queryErr := queryFindFileRows(ctx, db, query, limit)
	if queryErr != nil {
		return nil, queryErr
	}

	defer rows.Close()

	return processFindFileResults(ctx, db, query, rows, useCache)
}

func processFindFileResults(
	ctx context.Context,
	db *sql.DB,
	query string,
	rows *sql.Rows,
	useCache bool,
) ([]FileFindResult, *apperror.AppError) {
	results, scanErr := scanFindFileRows(rows)
	if scanErr != nil {
		return nil, scanErr
	}

	maybeUpdateCache(ctx, db, "find:"+query, castFileFindResults(results), useCache)

	return results, nil
}

func queryFindFileRows(
	ctx context.Context,
	db *sql.DB,
	query string,
	limit int,
) (*sql.Rows, *apperror.AppError) {
	if limit > 0 {
		return queryFindFileWithLimit(ctx, db, query, limit)
	}

	return queryFindFileAll(ctx, db, query)
}

func queryFindFileWithLimit(
	ctx context.Context,
	db *sql.DB,
	query string,
	limit int,
) (*sql.Rows, *apperror.AppError) {
	sqlQuery := "SELECT RelativePath, AbsolutePath FROM RepoFile WHERE RelativePath LIKE ? LIMIT ?"
	likePattern := "%" + query + "%"
	rows, err := db.QueryContext(ctx, sqlQuery, likePattern, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "searcher.queryFindFileWithLimit.Query")
	}

	return rows, nil
}

func queryFindFileAll(
	ctx context.Context,
	db *sql.DB,
	query string,
) (*sql.Rows, *apperror.AppError) {
	sqlQuery := "SELECT RelativePath, AbsolutePath FROM RepoFile WHERE RelativePath LIKE ?"
	likePattern := "%" + query + "%"
	rows, err := db.QueryContext(ctx, sqlQuery, likePattern)
	if err != nil {
		return nil, apperror.WrapSimple(err, "searcher.queryFindFileAll.Query")
	}

	return rows, nil
}

func scanFindFileRows(rows *sql.Rows) ([]FileFindResult, *apperror.AppError) {
	var results []FileFindResult
	for rows.Next() {
		var r FileFindResult
		err := rows.Scan(&r.RelativePath, &r.AbsolutePath)
		if err != nil {
			return nil, apperror.WrapSimple(err, "searcher.scanFindFileRows.Scan")
		}

		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "searcher.scanFindFileRows.Rows")
	}

	return results, nil
}

// FindFileRegex searches RepoDB using regex on the filename.
func FindFileRegex(
	ctx context.Context,
	db *sql.DB,
	expr string,
	limit int,
	useCache bool,
) ([]FileFindResult, *apperror.AppError) {
	cacheKey := "find_regex:" + expr
	res, hasCache, appErr := maybeGetCachedFindResults(ctx, db, cacheKey, limit, useCache)
	if appErr != nil {
		return nil, appErr
	}

	if hasCache {
		return res, nil
	}

	rows, queryErr := db.QueryContext(ctx, "SELECT RelativePath, AbsolutePath FROM RepoFile")
	if queryErr != nil {
		return nil, apperror.WrapSimple(queryErr, "searcher.FindFileRegex.Query")
	}

	defer rows.Close()

	return processFindRegexResults(ctx, db, cacheKey, rows, expr, limit, useCache)
}

func processFindRegexResults(
	ctx context.Context,
	db *sql.DB,
	cacheKey string,
	rows *sql.Rows,
	expr string,
	limit int,
	useCache bool,
) ([]FileFindResult, *apperror.AppError) {
	lz := lazyregex.New(expr)
	results, scanErr := scanFindFileRegexRows(rows, lz, limit)
	if scanErr != nil {
		return nil, scanErr
	}

	maybeUpdateCache(ctx, db, cacheKey, castFileFindResults(results), useCache)

	return results, nil
}

func scanFindFileRegexRows(
	rows *sql.Rows,
	lz *lazyregex.LazyRegexp,
	limit int,
) ([]FileFindResult, *apperror.AppError) {
	var results []FileFindResult
	for rows.Next() {
		hasHitLimit := limit > 0 && len(results) >= limit
		if hasHitLimit {
			break
		}

		var r FileFindResult
		if err := rows.Scan(&r.RelativePath, &r.AbsolutePath); err != nil {
			return nil, apperror.WrapSimple(err, "searcher.scanFindFileRegexRows.Scan")
		}

		if lz.Re().MatchString(r.RelativePath) {
			results = append(results, r)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.WrapSimple(err, "searcher.scanFindFileRegexRows.Rows")
	}

	return results, nil
}

func getCachedFindResults(
	ctx context.Context,
	db *sql.DB,
	key string,
	limit int,
) ([]FileFindResult, bool, *apperror.AppError) {
	var cached string
	err := db.QueryRowContext(ctx, "SELECT ResultJson FROM SearchCache WHERE Query = ?", key).Scan(&cached)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, apperror.WrapSimple(err, "searcher.getCachedFindResults.Scan")
	}

	if cached == "" {
		return nil, false, nil
	}

	return parseCachedFindResults(ctx, db, key, cached, limit)
}

func parseCachedFindResults(
	ctx context.Context,
	db *sql.DB,
	key string,
	cached string,
	limit int,
) ([]FileFindResult, bool, *apperror.AppError) {
	var res []FileFindResult
	if err := json.Unmarshal([]byte(cached), &res); err != nil {
		return nil, false, apperror.WrapSimple(err, "searcher.parseCachedFindResults.Unmarshal")
	}

	if appErr := incrementCacheHits(ctx, db, key); appErr != nil {
		appErr.HandleError()
	}

	hasTruncate := limit > 0 && len(res) > limit
	if hasTruncate {
		res = res[:limit]
	}

	return res, true, nil
}

// helper to reuse cache tables
func castFileFindResults(in []FileFindResult) []SearchResult {
	var out []SearchResult
	for _, i := range in {
		out = append(out, SearchResult{
			RelativePath: i.RelativePath,
			FilePath:     i.AbsolutePath,
		})
	}

	return out
}

type FileReadResult struct {
	RelativePath string `json:"relative_path"`
	AbsolutePath string `json:"absolute_path"`
	Content      string `json:"content"`
}

// FindAndRead searches RepoDB for file names, and then reads their content.
func FindAndRead(
	ctx context.Context,
	db *sql.DB,
	query string,
	isRegex bool,
	limit int,
	useCache bool,
) ([]FileReadResult, *apperror.AppError) {
	files, err := resolveFindFiles(ctx, db, query, isRegex, limit, useCache)
	if err != nil {
		return nil, err
	}

	var results []FileReadResult
	for _, f := range files {
		item, readErr := readFileFindContent(ctx, db, f)
		if readErr != nil {
			return nil, readErr
		}

		results = append(results, item)
	}

	return results, nil
}

func resolveFindFiles(
	ctx context.Context,
	db *sql.DB,
	query string,
	isRegex bool,
	limit int,
	useCache bool,
) ([]FileFindResult, *apperror.AppError) {
	if isRegex {
		return FindFileRegex(ctx, db, query, limit, useCache)
	}

	return FindFile(ctx, db, query, limit, useCache)
}

func readFileFindContent(
	ctx context.Context,
	db *sql.DB,
	f FileFindResult,
) (FileReadResult, *apperror.AppError) {
	var isBig int
	var content string
	err := db.QueryRowContext(ctx, "SELECT IsBig, Content FROM RepoFile WHERE RelativePath = ?", f.RelativePath).Scan(&isBig, &content)
	if err != nil {
		return FileReadResult{}, apperror.WrapSimple(err, "searcher.readFileFindContent.Query")
	}

	hasBigContent := isBig == 1
	if hasBigContent {
		content = "[BIG_FILE_CONTENT]"
	}

	return FileReadResult{
		RelativePath: f.RelativePath,
		AbsolutePath: f.AbsolutePath,
		Content:      content,
	}, nil
}

func maybeGetCachedFindResults(
	ctx context.Context,
	db *sql.DB,
	cacheKey string,
	limit int,
	useCache bool,
) ([]FileFindResult, bool, *apperror.AppError) {
	if !useCache {
		return nil, false, nil
	}

	return getCachedFindResults(ctx, db, cacheKey, limit)
}
