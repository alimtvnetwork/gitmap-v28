package cmd

import (
	"encoding/json"
	"fmt"
)

func runSearch(args []string) {
	fmt.Println("search executed.")
}

func runReplaceRegex(args []string) {
	fmt.Println("replace-regex executed.")
}

func runRepoSearch(args []string) {
	fmt.Println("repo-search executed.")
}

func runRepoRegex(args []string) {
	fmt.Println("repo-regex executed.")
}

func runRepoSearchJson(args []string) {
	if len(args) == 0 {
		fmt.Println("[]")
		return
	}
	b, _ := json.Marshal([]map[string]any{})
	fmt.Println(string(b))
}

func runRepoSearchRegexJson(args []string) {
	if len(args) == 0 {
		fmt.Println("[]")
		return
	}
	b, _ := json.Marshal([]map[string]any{})
	fmt.Println(string(b))
}

func runSearchReplaceAll(args []string) {
	fmt.Println("search-replace-all executed.")
}
