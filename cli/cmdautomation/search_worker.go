package cmdautomation

import (
	"bytes"
	"os"
	"regexp"
	"sync"
)

func searchWorkerRoutine(jobs <-chan string, results chan<- []SearchMatch, opts SearchOptions, wg *sync.WaitGroup) {
	defer wg.Done()
	reg, _ := resolveWorkerRegex(opts)

	for path := range jobs {
		matches := searchFileContent(path, opts, reg)
		if len(matches) > 0 {
			results <- matches
		}
	}
}

func resolveWorkerRegex(opts SearchOptions) (*regexp.Regexp, error) {
	if !opts.IsRegex {
		return nil, nil
	}
	appErrReg, err := GetRegex(opts.Pattern, opts.IsCaseInsensitive)
	if err != nil {
		return nil, err
	}
	return appErrReg, nil
}

func searchFileContent(path string, opts SearchOptions, reg *regexp.Regexp) []SearchMatch {
	data, err := os.ReadFile(path)
	if err != nil || HasBinaryContent(data) {
		return nil
	}

	var matches []SearchMatch
	lines := bytes.Split(data, []byte("\n"))
	for idx, line := range lines {
		if isLineMatch(line, opts, reg) {
			matches = append(matches, SearchMatch{
				Path:        path,
				LineNumber:  idx + 1,
				LineContent: string(bytes.TrimSpace(line)),
			})
		}
	}
	return matches
}

func isLineMatch(line []byte, opts SearchOptions, reg *regexp.Regexp) bool {
	if opts.IsRegex && reg != nil {
		return reg.Match(line)
	}
	return HasLiteralMatch(line, opts.Pattern, opts.IsCaseInsensitive)
}

func gatherSearchResults(results <-chan []SearchMatch) []SearchMatch {
	var all []SearchMatch
	for chunk := range results {
		all = append(all, chunk...)
	}
	return all
}
