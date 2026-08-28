package cmd

import (
	"context"
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

	contextBefore := 0
	contextAfter := 0

	if len(args) > 2 {
		if val, err := strconv.Atoi(args[2]); err == nil && val > 0 {
			contextBefore = val
		}
	}
	if len(args) > 3 {
		if val, err := strconv.Atoi(args[3]); err == nil && val > 0 {
			contextAfter = val
		}
	}

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

	// Check cache
	cacheKey := fmt.Sprintf("file-search:%s:%s:%d:%d", filePath, pattern, contextBefore, contextAfter)
	var cachedJson string
	err = db.QueryRowContext(ctx, "SELECT ResultJson FROM SearchCache WHERE Query = ?", cacheKey).Scan(&cachedJson)
	if err == nil && cachedJson != "" {
		var res []searcher.SearchResult
		if err := json.Unmarshal([]byte(cachedJson), &res); err == nil {
			db.ExecContext(ctx, "UPDATE SearchCache SET Hits = Hits + 1 WHERE Query = ?", cacheKey)
			printFileSearchResults(res)
			return nil
		}
	}

	// Fetch file content
	var content string
	var absPath string
	err = db.QueryRowContext(ctx, "SELECT AbsolutePath, Content FROM RepoFile WHERE RelativePath = ?", filePath).Scan(&absPath, &content)
	if err != nil {
		// Fallback to reading from disk
		b, errDisk := os.ReadFile(filePath)
		if errDisk != nil {
			return apperror.WrapSimple(errDisk, "Error reading file:")
		}
		content = string(b)
		absPath = filePath
	}

	lines := strings.Split(content, "\n")
	var results []searcher.SearchResult

	for i, line := range lines {
		if rx.MatchString(line) {
			startIdx := i - contextBefore
			if startIdx < 0 {
				startIdx = 0
			}
			endIdx := i + contextAfter
			if endIdx >= len(lines) {
				endIdx = len(lines) - 1
			}

			var matchedContext []string
			for j := startIdx; j <= endIdx; j++ {
				matchedContext = append(matchedContext, fmt.Sprintf("%d: %s", j+1, lines[j]))
			}

			results = append(results, searcher.SearchResult{
				MatchedText:   strings.Join(matchedContext, "\n"),
				StartPosition: i + 1, // Line number
				EndPosition:   endIdx + 1,
				FilePath:      absPath,
				RelativePath:  filePath,
			})
		}
	}

	if len(results) > 0 {
		b, _ := json.Marshal(results)
		db.ExecContext(ctx, "INSERT INTO SearchCache (Query, ResultJson, Timestamp, Hits) VALUES (?, ?, CURRENT_TIMESTAMP, 1)", cacheKey, string(b))
	}

	printFileSearchResults(results)
	return nil
}

func printFileSearchResults(res []searcher.SearchResult) {
	for _, r := range res {
		fmt.Printf("Match at line %d:\n", r.StartPosition)
		fmt.Println(r.MatchedText)
		fmt.Println(strings.Repeat("-", 40))
	}
}
