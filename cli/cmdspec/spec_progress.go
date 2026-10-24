package cmdspec

import (
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"path/filepath"
	"sort"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

const specProgressHelp = `spec progress — per-workstream completion overview for a spec.

Description:
  Parses the .md files under --spec for ledger checkboxes
  (- [x] done / - [ ] todo) and renders a per-workstream progress bar.
  Read-only.

Usage:
  gitmap spec progress --spec <path>

Flags:
  --spec <path>   spec directory to parse (required)
  --help, -h      show this help

Example:
  gitmap spec progress --spec 02-spec/21-app/gitmap-help-text-and-ui-enhancement
`

// checkboxPatternText is the markdown checkbox pattern, compiled lazily on first use.
const checkboxPatternText = `(?m)^\s*-\s*\[([ xX])\]\s*(.*)$`

var checkboxPattern = lazyregex.New(checkboxPatternText)

// specStream is one parsed markdown file's checkbox rollup.
type specStream struct {
	name     string
	done     int
	total    int
	nextTodo string
}

// RunSpecProgress implements `gitmap spec progress --spec <path>`.
func RunSpecProgress(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(specProgressHelp)
			return nil
		}
	}
	fs := flag.NewFlagSet("spec progress", flag.ContinueOnError)
	specDir := fs.String("spec", "", "spec directory to parse")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(specProgressHelp)
		return apperror.WrapSimple(parseErr, "spec.progress.flags")
	}
	if strings.TrimSpace(*specDir) == "" {
		fmt.Print(specProgressHelp)
		return appfault.NewValidationError("missing required flag: --spec")
	}
	info, statErr := os.Stat(*specDir)
	if statErr != nil || !info.IsDir() {
		return appfault.NewNotFoundError(fmt.Sprintf("spec dir not found: %s", *specDir))
	}

	streams := parseSpecStreams(*specDir)
	if len(streams) == 0 {
		fmt.Printf("spec progress: no markdown files with checkboxes under %s\n", *specDir)
		return nil
	}

	totalDone := 0
	totalAll := 0
	for _, s := range streams {
		totalDone += s.done
		totalAll += s.total
		fmt.Printf("  %-28s %s %d/%d done\n", s.name, progressBar(s.done, s.total), s.done, s.total)
	}
	fmt.Printf("  %-28s %s %d/%d done\n", "TOTAL", progressBar(totalDone, totalAll), totalDone, totalAll)
	if next := nextSpecTodo(streams); next != "" {
		fmt.Printf("next: %s\n", next)
	}

	return nil
}

// parseSpecStreams walks dir for .md files and rolls up their checkboxes.
func parseSpecStreams(dir string) []specStream {
	var streams []specStream
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && path != dir {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		stream := specStream{name: streamName(path)}
		for _, m := range checkboxPattern.Re().FindAllStringSubmatch(string(data), -1) {
			stream.total++
			if m[1] == "x" || m[1] == "X" {
				stream.done++
			} else if stream.nextTodo == "" {
				stream.nextTodo = strings.TrimSpace(m[2])
			}
		}
		if stream.total > 0 {
			streams = append(streams, stream)
		}

		return nil
	})
	sort.Slice(streams, func(i, j int) bool { return streams[i].name < streams[j].name })

	return streams
}

// streamName derives a workstream label from the file name.
func streamName(path string) string {
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	base = strings.TrimLeft(base, "0123456789-")
	if base == "" {
		base = filepath.Base(path)
	}

	return base
}

// progressBar renders an 8-block bar for done/total.
func progressBar(done, total int) string {
	const width = 8
	filled := 0
	if total > 0 {
		filled = done * width / total
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

// nextSpecTodo returns the first todo of the least-complete stream.
func nextSpecTodo(streams []specStream) string {
	best := ""
	bestRatio := 2.0
	for _, s := range streams {
		if s.total == 0 || s.nextTodo == "" {
			continue
		}
		ratio := float64(s.done) / float64(s.total)
		if ratio < bestRatio {
			bestRatio = ratio
			best = s.name + ": " + s.nextTodo
		}
	}

	return best
}
