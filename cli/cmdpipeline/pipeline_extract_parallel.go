package cmdpipeline

import (
	"runtime"
	"strings"
	"sync"
)

// PreParsedLogLine holds pre-processed state for a single log line.
type PreParsedLogLine struct {
	Job           string
	Step          string
	Text          string
	Key           string
	CompactedText string
	ContextLimit  int
	IsBlank       bool
	IsIgnored     bool
	IsWarning     bool
	IsWarningCont bool
	IsError       bool
	IsOk          bool
}

// splitLogLines splits raw log text into a line slice, normalizing line endings.
func splitLogLines(rawLogs string) []string {
	if len(rawLogs) == 0 {
		return nil
	}

	normalized := strings.ReplaceAll(rawLogs, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")

	return strings.Split(normalized, "\n")
}

func resolveParallelLogWorkerCount(totalLines int) int {
	cpuCores := runtime.NumCPU()
	if cpuCores < 1 {
		cpuCores = 1
	}
	maxWorkers := cpuCores
	if maxWorkers > 8 {
		maxWorkers = 8
	}
	neededWorkers := (totalLines + 99) / 100
	if neededWorkers < maxWorkers {
		return neededWorkers
	}

	return maxWorkers
}

func dispatchPreParseWorker(
	wg *sync.WaitGroup,
	lines []string,
	parsed []PreParsedLogLine,
	startIdx int,
	endIdx int,
) {
	defer wg.Done()
	for i := startIdx; i < endIdx; i++ {
		parsed[i] = preParseSingleLogLine(lines[i])
	}
}

func preParseSingleLogLine(raw string) PreParsedLogLine {
	if isBlankLogLine(raw) {
		return PreParsedLogLine{IsBlank: true}
	}

	job, step, text, isError := parseLogLine(raw)
	if text == "" || isIgnoredLogLine(text) {
		return PreParsedLogLine{Text: text, IsIgnored: true}
	}

	key := job + "|||" + step
	isWarn := isWarningLine(text)
	isWarnCont := isWarningContinuationLine(text)
	compacted := compactLinkerCommandLine(text)
	ctxLimit := resolveContextLimit(text)
	isOk := isOkLogLine(text)

	return PreParsedLogLine{
		Job:           job,
		Step:          step,
		Text:          text,
		Key:           key,
		CompactedText: compacted,
		ContextLimit:  ctxLimit,
		IsWarning:     isWarn,
		IsWarningCont: isWarnCont,
		IsError:       isError,
		IsOk:          isOk,
	}
}

// ParallelPreParseLogLines pre-parses log lines in parallel using worker goroutines.
func ParallelPreParseLogLines(lines []string) []PreParsedLogLine {
	totalLines := len(lines)
	if totalLines == 0 {
		return nil
	}

	parsed := make([]PreParsedLogLine, totalLines)
	if totalLines <= 100 {
		for i, line := range lines {
			parsed[i] = preParseSingleLogLine(line)
		}

		return parsed
	}

	workerCount := resolveParallelLogWorkerCount(totalLines)
	chunkSize := (totalLines + workerCount - 1) / workerCount

	var wg sync.WaitGroup
	for w := 0; w < workerCount; w++ {
		startIdx := w * chunkSize
		endIdx := resolvePartitionEnd(startIdx+chunkSize, totalLines)
		if startIdx >= endIdx {
			continue
		}
		wg.Add(1)
		go dispatchPreParseWorker(&wg, lines, parsed, startIdx, endIdx)
	}
	wg.Wait()

	return parsed
}

func scanLogLinesParallel(lines []string, jobMap map[string]*FailedJobItem, order *[]string) {
	preParsed := ParallelPreParseLogLines(lines)
	aggregatePreParsedLines(preParsed, jobMap, order)
}

func aggregatePreParsedLines(
	preParsed []PreParsedLogLine,
	jobMap map[string]*FailedJobItem,
	order *[]string,
) {
	contextRemaining := 0
	var lastKey string
	warnBuf := make(map[string][]string)
	warnContRemaining := 0

	for _, pl := range preParsed {
		if pl.IsBlank {
			warnContRemaining = 0

			continue
		}

		if pl.IsIgnored || pl.Text == "" {
			warnContRemaining = 0

			continue
		}

		if pl.IsWarning {
			bufferWarning(warnBuf, pl.Key, pl.Text)
			warnContRemaining = 10

			continue
		}

		if warnContRemaining > 0 && pl.IsWarningCont {
			bufferWarning(warnBuf, pl.Key, pl.Text)
			warnContRemaining--

			continue
		}
		warnContRemaining = 0

		if pl.IsError {
			recordPreParsedErrorLine(jobMap, order, pl, &contextRemaining, &lastKey, warnBuf)

			continue
		}

		appendPreParsedContextLine(jobMap, pl, &contextRemaining, &lastKey)
	}
}

func recordPreParsedErrorLine(
	jobMap map[string]*FailedJobItem,
	order *[]string,
	pl PreParsedLogLine,
	ctxRem *int,
	lastKey *string,
	warnBuf map[string][]string,
) {
	item := getOrCreateJobItem(jobMap, order, pl.Key, pl.Job, pl.Step)
	attachBufferedWarnings(item, warnBuf[pl.Key])
	item.ErrorLines = append(item.ErrorLines, pl.CompactedText)
	updateJobSummary(item, pl.CompactedText)
	*ctxRem = pl.ContextLimit
	*lastKey = pl.Key
}

func appendPreParsedContextLine(
	jobMap map[string]*FailedJobItem,
	pl PreParsedLogLine,
	ctxRem *int,
	lastKey *string,
) {
	if *ctxRem <= 0 || *lastKey != pl.Key || pl.IsOk {
		return
	}

	item := jobMap[pl.Key]
	if item == nil {
		return
	}

	item.ErrorLines = append(item.ErrorLines, "    "+pl.CompactedText)
	if pl.IsWarning {
		item.Warnings = append(item.Warnings, pl.Text)
	}
	*ctxRem--
}

// ParseFailedLogLinesParallel parses raw log lines into failed jobs using chunked parallel workers.
func ParseFailedLogLinesParallel(rawLogs string) []FailedJobItem {
	if strings.TrimSpace(rawLogs) == "" {
		return nil
	}

	jobMap := make(map[string]*FailedJobItem)
	var order []string
	lines := splitLogLines(rawLogs)
	if len(lines) > 100 && runtime.NumCPU() > 1 {
		scanLogLinesParallel(lines, jobMap, &order)
	} else {
		scanLogLinesSequential(lines, jobMap, &order)
	}

	return assembleJobItems(jobMap, order, rawLogs)
}
