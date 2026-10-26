package cmdsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	searchconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdsearch/constants"
	searchenums "github.com/alimtvnetwork/gitmap-v28/cli/cmdsearch/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/searcher"
)

// RunReplaceRegex executes regex substitution across files.
func RunReplaceRegex(args []string) error {
	checkHelp(searchconstants.CmdReplaceRegex, args)
	if len(args) < 2 {
		printReplaceRegexUsage()

		return nil
	}

	fmt.Println(searchconstants.MsgReplaceRegexExec)

	return nil
}

func printReplaceRegexUsage() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println(searchconstants.UsageReplaceRegex)
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println(searchconstants.ExampleReplaceRegex)
}

// RunRepoSearch searches indexed repositories by keyword.
func RunRepoSearch(args []string) error {
	checkHelp(searchconstants.CmdRepoSearch, args)
	isUI, nonUIArgs := parseSearchUIFlag(args)
	if isUI {
		return launchSearchUI()
	}

	limit, cleanArgs := parseLimit(nonUIArgs)
	if len(cleanArgs) == 0 {
		printRepoSearchUsage()

		return nil
	}

	query := cleanArgs[0]
	action, _ := searchenums.ParseSearchAction(query)
	if action.IsHistoryAction() {
		return searcher.RenderAUMSearchHistoryTable(limit)
	}
	if action.IsClear() {
		fmt.Println(searchconstants.MsgRepoSearchClear)

		return nil
	}

	return executeRepoSearch(query, limit)
}

func printRepoSearchUsage() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println(searchconstants.UsageRepoSearch)
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Description:" + constants.ColorReset)
	fmt.Println("  " + searchconstants.DescRepoSearchMetadata)
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println(searchconstants.ExampleRepoSearch1)
	fmt.Println(searchconstants.ExampleRepoSearch2)
}

func executeRepoSearch(query string, limit int) error {
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, searchconstants.ErrDefaultError)
	}

	defer mainDB.Close()
	defer db.Close()

	res, searchErr := searcher.SearchRepoDB(ctx, db, query, limit, true)
	if searchErr != nil {
		pterm.Error.Println(searchErr)

		return nil
	}

	renderRepoSearchResults(res)

	return nil
}

func renderRepoSearchResults(res []searcher.SearchResult) {
	for _, r := range res {
		pterm.DefaultHeader.WithFullWidth().Printf(searchconstants.FmtFoundHeader, r.RelativePath, r.StartPosition)
		fmt.Println(r.MatchedText)
		fmt.Println()
	}
}

// RunRepoRegex searches indexed repositories using regular expressions.
func RunRepoRegex(args []string) error {
	checkHelp(searchconstants.CmdRepoRegex, args)
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println(searchconstants.UsageRepoRegex)

		return nil
	}

	query := cleanArgs[0]
	action, _ := searchenums.ParseSearchAction(query)
	if action.IsHistory() {
		fmt.Println(searchconstants.MsgRepoRegexHistoryEmpty)

		return nil
	}
	if action.IsClear() {
		fmt.Println(searchconstants.MsgRepoRegexClear)

		return nil
	}

	return executeRepoRegex(query, limit)
}

func executeRepoRegex(query string, limit int) error {
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, searchconstants.ErrDefaultError)
	}

	defer mainDB.Close()
	defer db.Close()

	res, searchErr := searcher.SearchRepoDBRegex(ctx, db, query, limit, true)
	if searchErr != nil {
		pterm.Error.Println(searchErr)

		return nil
	}

	renderRepoSearchResults(res)

	return nil
}

// RunRepoSearchJson searches repositories and emits results as JSON.
func RunRepoSearchJson(args []string) error {
	checkHelp(searchconstants.CmdRepoSearchJSON, args)
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println(searchconstants.EmptyJSONArray)

		return nil
	}

	return executeRepoSearchJSON(cleanArgs[0], limit)
}

func executeRepoSearchJSON(query string, limit int) error {
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		fmt.Println(searchconstants.EmptyJSONArray)

		return apperror.WrapSimple(err, searchconstants.ErrRepoSearchFatal)
	}

	defer mainDB.Close()
	defer db.Close()

	res, searchErr := searcher.SearchRepoDB(ctx, db, query, limit, true)
	if searchErr != nil {
		fmt.Println(searchconstants.EmptyJSONArray)

		return nil
	}

	b, _ := json.Marshal(res)
	fmt.Println(string(b))

	return nil
}

// RunRepoSearchRegexJson performs regex repo search and outputs JSON.
func RunRepoSearchRegexJson(args []string) error {
	checkHelp(searchconstants.CmdRepoSearchRegexJSON, args)
	limit, cleanArgs := parseLimit(args)
	if len(cleanArgs) == 0 {
		fmt.Println(searchconstants.EmptyJSONArray)

		return nil
	}

	return executeRepoSearchRegexJSON(cleanArgs[0], limit)
}

func executeRepoSearchRegexJSON(query string, limit int) error {
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		fmt.Println(searchconstants.EmptyJSONArray)

		return apperror.WrapSimple(err, searchconstants.ErrRepoRegexSearchFatal)
	}

	defer mainDB.Close()
	defer db.Close()

	res, searchErr := searcher.SearchRepoDBRegex(ctx, db, query, limit, true)
	if searchErr != nil {
		fmt.Println(searchconstants.EmptyJSONArray)

		return nil
	}

	b, _ := json.Marshal(res)
	fmt.Println(string(b))

	return nil
}

// RunSearchReplaceAll replaces search indexes and outputs.
func RunSearchReplaceAll(args []string) error {
	checkHelp(searchconstants.CmdSearchReplaceAll, args)
	if len(args) == 0 {
		fmt.Println(searchconstants.MsgSearchReplaceAllExec)

		return nil
	}

	action, _ := searchenums.ParseSearchAction(args[0])
	if action.IsResetAction() {
		os.RemoveAll(searchconstants.PathRepoSearchOutput)
		fmt.Println(searchconstants.MsgSearchReplaceAllReset)

		return nil
	}

	fmt.Println(searchconstants.MsgSearchReplaceAllExec)

	return nil
}
