package folder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// OutputFormat defines the rendering strategy for folder output.
type OutputFormatType string

type OutputFormat = OutputFormatType

const (
	FormatTree OutputFormatType = "tree"
	FormatMd   OutputFormatType = "md"
	FormatJson OutputFormatType = "json"
	FormatYaml OutputFormatType = "yaml"
	FormatFlat OutputFormatType = "flat"
)

// Options holds CLI configurations and flags for folder scanning.
type Options struct {
	TargetDir  string
	OutFile    string
	Format     OutputFormat
	IsDetailed bool
	Filter     FilterConfig
}

// Run parses arguments, scans directory files, and renders the requested format.
func Run(args []string) error {
	opts, err := ParseArgs(args)
	if err != nil {
		return err
	}

	files, err := ScanDirectory(opts.TargetDir, opts.Filter)
	if err != nil {
		return apperror.Wrap(err, fmt.Sprintf("scan directory %s", opts.TargetDir), nil)
	}

	content, err := RenderOutput(opts, files)
	if err != nil {
		return err
	}

	return writeOrPrintOutput(content, opts.OutFile)
}

// ScanDirectory walks the filesystem hierarchy and extracts metadata for qualifying files.
func ScanDirectory(root string, filter FilterConfig) ([]*FileMeta, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	var results []*FileMeta
	errWalk := filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, errIn error) error {
		if errIn != nil {
			return errIn
		}

		rel, errRel := filepath.Rel(absRoot, path)
		if errRel != nil {
			rel = path
		}

		if d.IsDir() {
			return handleDirSkip(rel, filter)
		}

		if filter.IsPathExcluded(rel, false) {
			return nil
		}

		meta, errMeta := ExtractMetadata(path, rel)
		if errMeta != nil {
			return nil
		}

		if filter.IsMetaAllowed(meta) {
			results = append(results, meta)
		}

		return nil
	})

	return results, errWalk
}

func handleDirSkip(rel string, filter FilterConfig) error {
	if rel != "." && filter.IsPathExcluded(rel, true) {
		return fs.SkipDir
	}

	if filter.MaxDepth > 0 && calculateDepth(rel) > filter.MaxDepth {
		return fs.SkipDir
	}

	return nil
}

func calculateDepth(rel string) int {
	norm := filepath.ToSlash(rel)
	if norm == "." || norm == "" {
		return 0
	}

	return strings.Count(norm, "/") + 1
}

// RenderOutput delegates to specific format renderers.
func RenderOutput(opts Options, files []*FileMeta) (string, error) {
	rootName := filepath.Base(opts.TargetDir)
	if rootName == "." || rootName == "/" || rootName == "\\" {
		rootName = "root"
	}

	switch opts.Format {
	case FormatTree:
		treeRoot := BuildTree(rootName, files)

		return RenderTree(treeRoot, opts.IsDetailed), nil
	case FormatMd:
		treeRoot := BuildTree(rootName, files)

		return RenderMarkdown(treeRoot, opts.IsDetailed), nil
	case FormatJson:
		report := BuildReport(opts.TargetDir, files)

		return RenderJson(report)
	case FormatYaml:
		report := BuildReport(opts.TargetDir, files)

		return RenderYaml(report)
	case FormatFlat:
		return RenderFlat(files, opts.IsDetailed), nil
	default:
		treeRoot := BuildTree(rootName, files)

		return RenderTree(treeRoot, opts.IsDetailed), nil
	}
}

func writeOrPrintOutput(content, outFile string) error {
	if outFile != "" {
		return os.WriteFile(outFile, []byte(content), 0644)
	}

	fmt.Print(content)

	return nil
}
