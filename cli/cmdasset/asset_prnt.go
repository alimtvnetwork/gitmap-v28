package cmdasset

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunAssetPrnt dispatches the screenshot downloader or URL extractor.
func RunAssetPrnt(args []string) error {
	opts := parsePrntArgs(args)
	if len(opts.URLs) == 0 {
		printAssetPrntHelp()
		return nil
	}

	results := executePrntBatch(opts)
	if opts.JsonFormat {
		return printPrntJSON(opts, results)
	}

	printPrntTerminalSummary(opts, results)
	return nil
}

func parsePrntArgs(args []string) PrntOptions {
	opts := PrntOptions{Concurrency: 4}
	var nonFlags []string

	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			return PrntOptions{}
		case a == "-e" || a == "--extract-only" || a == "-p" || a == "--print-only":
			opts.ExtractOnly = true
		case a == "-j" || a == "--json":
			opts.JsonFormat = true
		case a == "-o" || a == "--output-dir" || a == "--output":
			if i+1 < len(args) {
				opts.OutputDir = args[i+1]
				i++
			}
		case strings.HasPrefix(a, "--output-dir="):
			opts.OutputDir = strings.TrimPrefix(a, "--output-dir=")
		default:
			nonFlags = append(nonFlags, a)
		}
	}

	if len(nonFlags) == 2 && !strings.HasPrefix(nonFlags[1], "http://") && !strings.HasPrefix(nonFlags[1], "https://") {
		opts.URLs = []string{nonFlags[0]}
		opts.CustomDest = nonFlags[1]
		return opts
	}

	opts.URLs = nonFlags
	return opts
}

func printAssetPrntHelp() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap prnt <url1> [url2...] [flags]")
	fmt.Println("  gitmap download-prnt <url1> [url2...] [flags]")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Flags:" + constants.ColorReset)
	fmt.Println("  -e, --extract-only   Extract direct image URLs without downloading")
	fmt.Println("  -j, --json           Output structured JSON envelope")
	fmt.Println("  -o, --output-dir     Target directory for downloads (default: assets/screenshots)")
	fmt.Println("  -h, --help           Show this help text")
	fmt.Println()
	fmt.Println(constants.ColorCyan + "Examples:" + constants.ColorReset)
	fmt.Println("  gitmap prnt https://prnt.sc/FnF7IMTncnd5")
	fmt.Println("  gitmap prnt https://prnt.sc/url1 https://prnt.sc/url2")
	fmt.Println("  gitmap prnt https://prnt.sc/url1 https://prnt.sc/url2 -e")
	fmt.Println("  gitmap prnt https://prnt.sc/url1 --json")
}

func executePrntBatch(opts PrntOptions) []PrntResultItem {
	results := make([]PrntResultItem, len(opts.URLs))
	sem := make(chan struct{}, opts.Concurrency)
	var wg sync.WaitGroup

	for i, u := range opts.URLs {
		wg.Add(1)
		go func(idx int, targetURL string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			results[idx] = processSingleURL(targetURL, opts)
		}(i, u)
	}

	wg.Wait()
	return results
}

func processSingleURL(rawURL string, opts PrntOptions) PrntResultItem {
	item := PrntResultItem{SourceURL: rawURL}
	imgURL, err := resolveImageURL(rawURL)
	if err != nil {
		item.IsSuccess = false
		item.Error = err.Error()
		return item
	}

	item.DirectURL = imgURL
	item.Slug = extractSlugFromURL(rawURL)

	if opts.ExtractOnly {
		item.IsSuccess = true
		return item
	}

	data, err := fetchImageBytes(imgURL)
	if err != nil {
		item.IsSuccess = false
		item.Error = fmt.Sprintf("failed to download from %s: %v", imgURL, err)
		return item
	}

	finalPath := resolveDestinationPath(item.Slug, opts.CustomDest, opts.OutputDir)
	_ = os.MkdirAll(filepath.Dir(finalPath), 0755)

	if writeErr := os.WriteFile(finalPath, data, 0644); writeErr != nil {
		item.IsSuccess = false
		item.Error = fmt.Sprintf("failed to save to %s: %v", finalPath, writeErr)
		return item
	}

	item.SavedPath = finalPath
	item.FileSizeBytes = len(data)
	item.BrainSyncPath = maybeSyncToActiveBrain(item.Slug, data)
	item.IsSuccess = true

	return item
}

func printPrntTerminalSummary(opts PrntOptions, results []PrntResultItem) {
	succCount := 0
	for _, r := range results {
		if r.IsSuccess {
			succCount++
		}
	}

	if opts.ExtractOnly {
		fmt.Printf("\n%sExtracted %d direct image URL(s):%s\n", constants.ColorCyan, succCount, constants.ColorReset)
		for _, r := range results {
			printPrntExtractItem(r)
		}
		return
	}

	fmt.Printf("\n%s✔ Processed %d screenshot asset(s) (%d succeeded, %d failed):%s\n",
		constants.ColorGreen, len(results), succCount, len(results)-succCount, constants.ColorReset)

	for _, r := range results {
		printPrntDownloadItem(r)
	}
}

func printPrntExtractItem(r PrntResultItem) {
	if r.IsSuccess {
		fmt.Printf("  • %s -> %s%s%s\n", r.SourceURL, constants.ColorGreen, r.DirectURL, constants.ColorReset)
		return
	}
	fmt.Printf("  • %s -> %sError: %s%s\n", r.SourceURL, constants.ColorRed, r.Error, constants.ColorReset)
}

func printPrntDownloadItem(r PrntResultItem) {
	if !r.IsSuccess {
		fmt.Printf("  • Source: %s\n", r.SourceURL)
		fmt.Printf("    Error:  %s%s%s\n", constants.ColorRed, r.Error, constants.ColorReset)
		return
	}

	fmt.Printf("  • Source: %s\n", r.SourceURL)
	fmt.Printf("    Direct: %s\n", r.DirectURL)
	fmt.Printf("    Saved:  %s%s%s (%d bytes)\n", constants.ColorGreen, r.SavedPath, constants.ColorReset, r.FileSizeBytes)
	printPrntBrainSync(r.BrainSyncPath)
}

func printPrntBrainSync(brainPath string) {
	if brainPath == "" {
		return
	}
	fmt.Printf("    Brain:  %s%s%s\n", constants.ColorCyan, brainPath, constants.ColorReset)
}

func printPrntJSON(opts PrntOptions, results []PrntResultItem) error {
	succCount := 0
	var errs []string
	for _, r := range results {
		if r.IsSuccess {
			succCount++
		} else if r.Error != "" {
			errs = append(errs, fmt.Sprintf("%s: %s", r.SourceURL, r.Error))
		}
	}

	code := 200
	msg := "OK"
	if succCount == 0 && len(results) > 0 {
		code = 500
		msg = "All extractions failed"
	}

	env := PrntResponseEnvelope{
		Status: PrntStatus{
			IsSuccess: succCount > 0,
			Code:      code,
			Message:   msg,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		Attributes: PrntAttributes{
			TotalRequested: len(results),
			TotalSucceeded: succCount,
			TotalFailed:    len(results) - succCount,
			ExtractOnly:    opts.ExtractOnly,
		},
		Results: results,
		Errors:  errs,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(env)
}
