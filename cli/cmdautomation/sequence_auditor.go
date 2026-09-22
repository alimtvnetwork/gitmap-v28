package cmdautomation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var (
	numPrefixRegex = regexp.MustCompile(`^(\d+)[-_]`)
	h1HeaderRegex  = regexp.MustCompile(`(?m)^(#\s*)(\d+)([\s—\-]+.*)$`)
)

// RunSequenceAuditor scans directories for numbered markdown files and verifies sequence and H1 title consistency.
func RunSequenceAuditor(opts SequenceAuditorOptions) (SequenceAuditorResult, *apperror.AppError) {
	start := time.Now()
	res := SequenceAuditorResult{}
	targetDir := resolveGuardTargetDir(opts.Dir)

	walkErr := filepath.Walk(targetDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && isExcludedDir(info.Name()) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			auditDirectorySequence(p, opts, &res)
		}
		return nil
	})

	if walkErr != nil {
		ctx := map[string]any{"dir": targetDir, "err": walkErr.Error()}
		return res, apperror.Wrap(walkErr, "sequence_auditor", ctx)
	}

	res.Duration = time.Since(start)
	return res, nil
}

func auditDirectorySequence(dir string, opts SequenceAuditorOptions, res *SequenceAuditorResult) {
	numbered := collectNumberedMarkdownFiles(dir)
	if len(numbered) == 0 {
		return
	}

	res.ScannedDirs++
	res.TotalNumberedFiles += len(numbered)

	dirRes := SequenceDirResult{
		Dir:                filepath.ToSlash(dir),
		NumberedFilesCount: len(numbered),
	}

	checkSequenceGaps(numbered, &dirRes)
	checkAndFixTitles(dir, numbered, opts, &dirRes)

	dirRes.IsClean = len(dirRes.SequenceGaps) == 0 && len(dirRes.TitleMismatches) == 0
	res.TotalGaps += len(dirRes.SequenceGaps)
	res.TotalMismatches += len(dirRes.TitleMismatches)
	res.TotalFixed += len(dirRes.FixedTitles)
	res.DirResults = append(res.DirResults, dirRes)
}

type numberedFileInfo struct {
	Number   int
	Filename string
	Path     string
}

func collectNumberedMarkdownFiles(dir string) []numberedFileInfo {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var list []numberedFileInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		m := numPrefixRegex.FindStringSubmatch(entry.Name())
		if len(m) > 1 {
			num, _ := strconv.Atoi(m[1])
			list = append(list, numberedFileInfo{
				Number:   num,
				Filename: entry.Name(),
				Path:     filepath.Join(dir, entry.Name()),
			})
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Number < list[j].Number })
	return list
}

func checkSequenceGaps(list []numberedFileInfo, dirRes *SequenceDirResult) {
	if len(list) == 0 {
		return
	}
	first := list[0].Number
	if first != 0 && first != 1 {
		return
	}

	expected := first
	for _, item := range list {
		if item.Number != expected {
			gap := filepath.ToSlash(filepath.Join(dirRes.Dir, item.Filename)) +
				" (found " + strconv.Itoa(item.Number) + ", expected " + strconv.Itoa(expected) + ")"
			dirRes.SequenceGaps = append(dirRes.SequenceGaps, gap)
		}
		expected++
	}
}

func checkAndFixTitles(dir string, list []numberedFileInfo, opts SequenceAuditorOptions, dirRes *SequenceDirResult) {
	for _, item := range list {
		data, err := os.ReadFile(item.Path)
		if err != nil {
			continue
		}
		inspectAndFixSingleTitle(item, data, opts, dirRes)
	}
}

func inspectAndFixSingleTitle(item numberedFileInfo, data []byte, opts SequenceAuditorOptions, dirRes *SequenceDirResult) {
	loc := h1HeaderRegex.FindSubmatchIndex(data)
	if len(loc) < 8 {
		return
	}
	h1NumStr := string(data[loc[4]:loc[5]])
	h1Num, _ := strconv.Atoi(h1NumStr)

	if h1Num == item.Number {
		return
	}

	rel := filepath.ToSlash(filepath.Join(dirRes.Dir, item.Filename))
	msg := rel + " (prefix " + strconv.Itoa(item.Number) + " != H1 header " + h1NumStr + ")"
	dirRes.TitleMismatches = append(dirRes.TitleMismatches, msg)

	if opts.IsFixMode {
		fixTitleHeaderOnDisk(item.Path, data, item.Number, dirRes)
	}
}

func fixTitleHeaderOnDisk(path string, data []byte, correctNum int, dirRes *SequenceDirResult) {
	formattedNum := fmt.Sprintf("%02d", correctNum)
	replacement := "${1}" + formattedNum + "${3}"
	fixed := h1HeaderRegex.ReplaceAll(data, []byte(replacement))
	if err := os.WriteFile(path, fixed, 0644); err == nil {
		dirRes.FixedTitles = append(dirRes.FixedTitles, filepath.ToSlash(path))
	}
}
