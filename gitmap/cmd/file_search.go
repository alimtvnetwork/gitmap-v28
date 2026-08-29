package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/searcher"
)

func runFileSearch(args []string) error {
	if len(args) < 2 {
		return apperror.New("Usage: gitmap file-search <file> <regex> [contextBefore] [contextAfter]\n\nAlternative LLM Examples:\n  Get-ChildItem -Path . -Recurse -File | Select-String \"type SearchResult struct\"\n  Get-ChildItem -Path cmd -Filter *.go | Select-String \"func dispatch[A-Z]\"\n  cat cmd/root.go | Select-String \"func finishCommandAudit\" -Context 0,10", "E9000", nil)
	}

	filePath := args[0]
	pattern := args[1]
	contextBefore, contextAfter := parseContextCounts(args)

	rx, err := regexp.Compile(pattern)
	if err != nil {
		return apperror.WrapSimple(err, "Invalid regex pattern:")
	}

	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, "Error connecting to db:")
	}
	defer mainDB.Close()
	defer db.Close()

	cacheKey := fmt.Sprintf("file-search:%s:%s:%d:%d", filePath, pattern, contextBefore, contextAfter)
	if res, ok := fetchCachedFileSearch(ctx, db, cacheKey); ok {
		printFileSearchResults(res)
		return nil
	}

	content, absPath, errContent := readFileSearchContent(ctx, db, filePath)
	if errContent != nil {
		return errContent
	}

	results := executeLineSearch(content, rx, filePath, absPath, contextBefore, contextAfter)
	printFileSearchResults(results)
	updateFileSearchCache(ctx, db, cacheKey, results)
	return nil
}

func parseContextCounts(args []string) (int, int) {
	before := parsePositionalInt(args, 2)
	after := parsePositionalInt(args, 3)
	return before, after
}

func parsePositionalInt(args []string, index int) int {
	if len(args) <= index {
		return 0
	}
	val, err := strconv.Atoi(args[index])
	if err != nil || val <= 0 {
		return 0
	}
	return val
}

func fetchCachedFileSearch(ctx context.Context, db *sql.DB, cacheKey string) ([]searcher.SearchResult, bool) {
	var cachedJson string
	err := db.QueryRowContext(ctx, "SELECT ResultJson FROM SearchCache WHERE Query = ?", cacheKey).Scan(&cachedJson)
	if err != nil || cachedJson == "" {
		return nil, false
	}
	var res []searcher.SearchResult
	if err := json.Unmarshal([]byte(cachedJson), &res); err != nil {
		return nil, false
	}
	db.ExecContext(ctx, "UPDATE SearchCache SET Hits = Hits + 1 WHERE Query = ?", cacheKey)
	return res, true
}

func readFileSearchContent(ctx context.Context, db *sql.DB, filePath string) (string, string, error) {
	var content, absPath string
	err := db.QueryRowContext(ctx, "SELECT AbsolutePath, Content FROM RepoFile WHERE RelativePath = ?", filePath).Scan(&absPath, &content)
	if err == nil {
		return content, absPath, nil
	}

	b, errDisk := os.ReadFile(filePath)
	if errDisk != nil {
		return "", "", apperror.WrapSimple(errDisk, "Error reading file:")
	}
	return string(b), filePath, nil
}

func executeLineSearch(content string, rx *regexp.Regexp, filePath, absPath string, before, after int) []searcher.SearchResult {
	lines := strings.Split(content, "\n")
	var results []searcher.SearchResult

	for i, line := range lines {
		if !rx.MatchString(line) {
			continue
		}
		matchedText, endIdx := buildContextSnippet(lines, i, before, after)
		results = append(results, searcher.SearchResult{
			MatchedText:   matchedText,
			StartPosition: i + 1,
			EndPosition:   endIdx + 1,
			RelativePath:  filePath,
			FilePath:      absPath,
		})
	}
	return results
}

func buildContextSnippet(lines []string, i, before, after int) (string, int) {
	startIdx := max(0, i-before)
	endIdx := min(len(lines)-1, i+after)

	var matchedContext []string
	for j := startIdx; j <= endIdx; j++ {
		matchedContext = append(matchedContext, fmt.Sprintf("%d: %s", j+1, lines[j]))
	}
	return strings.Join(matchedContext, "\n"), endIdx
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func printFileSearchResults(results []searcher.SearchResult) {
	for _, r := range results {
		fmt.Printf("--- Match at Line %d ---\n%s\n\n", r.StartPosition, r.MatchedText)
	}
}

func updateFileSearchCache(ctx context.Context, db *sql.DB, cacheKey string, results []searcher.SearchResult) {
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
