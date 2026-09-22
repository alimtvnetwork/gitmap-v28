package cmdautomation

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"sync"
)

type workerSearchContext struct {
	opts       SearchOptions
	reg        *regexp.Regexp
	patBytes   []byte
	lowerBytes []byte
}

func searchWorkerRoutine(jobs <-chan string, results chan<- []SearchMatch, opts SearchOptions, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx := buildWorkerContext(opts)

	for path := range jobs {
		matches := searchFileContent(path, ctx)
		if len(matches) > 0 {
			results <- matches
		}
	}
}

func buildWorkerContext(opts SearchOptions) workerSearchContext {
	reg, _ := resolveWorkerRegex(opts)
	var patBytes, lowerBytes []byte
	if !opts.IsRegex {
		patBytes = []byte(opts.Pattern)
		lowerBytes = []byte(strings.ToLower(opts.Pattern))
	}

	return workerSearchContext{
		opts:       opts,
		reg:        reg,
		patBytes:   patBytes,
		lowerBytes: lowerBytes,
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

func searchFileContent(path string, ctx workerSearchContext) []SearchMatch {
	data, err := os.ReadFile(path)
	if err != nil || HasBinaryContent(data) {
		return nil
	}
	if !hasFileMatch(data, ctx) {
		return nil
	}

	return scanMatchingLines(path, data, ctx)
}

func hasFileMatch(data []byte, ctx workerSearchContext) bool {
	if ctx.opts.IsRegex && ctx.reg != nil {
		return ctx.reg.Match(data)
	}
	if ctx.opts.IsCaseInsensitive {
		return bytes.Contains(bytes.ToLower(data), ctx.lowerBytes)
	}

	return bytes.Contains(data, ctx.patBytes)
}

func scanMatchingLines(path string, data []byte, ctx workerSearchContext) []SearchMatch {
	var matches []SearchMatch
	lineNum := 1
	remaining := data

	for len(remaining) > 0 {
		idx := bytes.IndexByte(remaining, '\n')
		var line []byte
		if idx >= 0 {
			line = remaining[:idx]
			remaining = remaining[idx+1:]
		} else {
			line = remaining
			remaining = nil
		}
		if isLineMatch(line, ctx) {
			matches = append(matches, SearchMatch{
				Path:        path,
				LineNumber:  lineNum,
				LineContent: string(bytes.TrimSpace(line)),
			})
		}
		lineNum++
	}

	return matches
}

func isLineMatch(line []byte, ctx workerSearchContext) bool {
	if ctx.opts.IsRegex && ctx.reg != nil {
		return ctx.reg.Match(line)
	}
	if ctx.opts.IsCaseInsensitive {
		return bytes.Contains(bytes.ToLower(line), ctx.lowerBytes)
	}

	return bytes.Contains(line, ctx.patBytes)
}

func gatherSearchResults(results <-chan []SearchMatch) []SearchMatch {
	var all []SearchMatch
	for chunk := range results {
		all = append(all, chunk...)
	}

	return all
}
