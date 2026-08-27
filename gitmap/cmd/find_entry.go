package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

func parseLimit(args []string) (int, []string) {
	limit := 0
	var cleanArgs []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--limit" || args[i] == "-l" {
			if i+1 < len(args) {
				limit, _ = strconv.Atoi(args[i+1])
				i++
			}
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

func runFind(args []string) {
	limit, args := parseLimit(args)
	fmt.Printf("find executed. (limit: %d, args: %v)\n", limit, args)
}

func runFindRegex(args []string) {
	limit, args := parseLimit(args)
	fmt.Printf("find-regex executed. (limit: %d, args: %v)\n", limit, args)
}

func runFindRead(args []string) {
	limit, args := parseLimit(args)
	fmt.Printf("find-read executed. (limit: %d, args: %v)\n", limit, args)
}

func runFindReadJson(args []string) {
	limit, args := parseLimit(args)
	fmt.Printf("find-read-json executed. (limit: %d, args: %v)\n", limit, args)
}

func runFindRegexRead(args []string) {
	limit, args := parseLimit(args)
	fmt.Printf("find-regex-read executed. (limit: %d, args: %v)\n", limit, args)
}

func runFindRegexReadJson(args []string) {
	limit, args := parseLimit(args)
	fmt.Printf("find-regex-read-json executed. (limit: %d, args: %v)\n", limit, args)
}

func runFindHelp(args []string) {
	fmt.Println("find help options: find, find-regex, find-read, find-read-json")
}

func runSearchHelp(args []string) {
	fmt.Println("search help options: search, search-replace-all, repo-search")
}

func runRegexHelp(args []string) {
	fmt.Println("regex help options: replace-regex, repo-regex, find-regex")
}
