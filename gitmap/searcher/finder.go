package searcher

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/lazyregex"
)

type FileFindResult struct {
	RelativePath string `json:"relative_path"`
	AbsolutePath string `json:"absolute_path"`
}

// FindFile searches RepoDB for file names matching the query or pattern.
func FindFile(ctx context.Context, db *sql.DB, query string, limit int, useCache bool) ([]FileFindResult, error) {
	if useCache {
		var cached string
		err := db.QueryRowContext(ctx, "SELECT ResultJson FROM SearchCache WHERE Query = ?", "find:"+query).Scan(&cached)
		if err == nil && cached != "" {
			var res []FileFindResult
			if err := json.Unmarshal([]byte(cached), &res); err == nil {
				db.ExecContext(ctx, "UPDATE SearchCache SET Hits = Hits + 1 WHERE Query = ?", "find:"+query)
				if limit > 0 && len(res) > limit {
					res = res[:limit]
				}
				return res, nil
			}
		}
	}

	sqlQuery := "SELECT RelativePath, AbsolutePath FROM RepoFile WHERE RelativePath LIKE ?"
	likePattern := "%" + query + "%"
	
	if limit > 0 {
		sqlQuery += " LIMIT ?"
	}

	var rows *sql.Rows
	var err error
	if limit > 0 {
		rows, err = db.QueryContext(ctx, sqlQuery, likePattern, limit)
	} else {
		rows, err = db.QueryContext(ctx, sqlQuery, likePattern)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []FileFindResult
	for rows.Next() {
		var r FileFindResult
		if err := rows.Scan(&r.RelativePath, &r.AbsolutePath); err == nil {
			results = append(results, r)
		}
	}

	if useCache {
		updateCache(ctx, db, "find:"+query, castFileFindResults(results))
	}

	return results, nil
}

// FindFileRegex searches RepoDB using regex on the filename.
func FindFileRegex(ctx context.Context, db *sql.DB, expr string, limit int, useCache bool) ([]FileFindResult, error) {
	if useCache {
		var cached string
		err := db.QueryRowContext(ctx, "SELECT ResultJson FROM SearchCache WHERE Query = ?", "find_regex:"+expr).Scan(&cached)
		if err == nil && cached != "" {
			var res []FileFindResult
			if err := json.Unmarshal([]byte(cached), &res); err == nil {
				db.ExecContext(ctx, "UPDATE SearchCache SET Hits = Hits + 1 WHERE Query = ?", "find_regex:"+expr)
				if limit > 0 && len(res) > limit {
					res = res[:limit]
				}
				return res, nil
			}
		}
	}

	lz := lazyregex.New(expr)
	rows, err := db.QueryContext(ctx, "SELECT RelativePath, AbsolutePath FROM RepoFile")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []FileFindResult
	for rows.Next() {
		if limit > 0 && len(results) >= limit {
			break
		}
		var r FileFindResult
		if err := rows.Scan(&r.RelativePath, &r.AbsolutePath); err == nil {
			if lz.Re().MatchString(r.RelativePath) {
				results = append(results, r)
			}
		}
	}

	if useCache {
		updateCache(ctx, db, "find_regex:"+expr, castFileFindResults(results))
	}

	return results, nil
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
func FindAndRead(ctx context.Context, db *sql.DB, query string, isRegex bool, limit int, useCache bool) ([]FileReadResult, error) {
	var files []FileFindResult
	var err error

	if isRegex {
		files, err = FindFileRegex(ctx, db, query, limit, useCache)
	} else {
		files, err = FindFile(ctx, db, query, limit, useCache)
	}

	if err != nil {
		return nil, err
	}

	var results []FileReadResult
	for _, f := range files {
		var isBig int
		var content string
		err := db.QueryRowContext(ctx, "SELECT IsBig, Content FROM RepoFile WHERE RelativePath = ?", f.RelativePath).Scan(&isBig, &content)
		if err != nil {
			continue
		}
		
		if isBig == 1 {
			// Stub for big file read
			content = "[BIG_FILE_CONTENT]"
		}

		results = append(results, FileReadResult{
			RelativePath: f.RelativePath,
			AbsolutePath: f.AbsolutePath,
			Content:      content,
		})
	}

	return results, nil
}
