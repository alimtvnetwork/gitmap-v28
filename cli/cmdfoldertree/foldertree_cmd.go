package cmdfoldertree

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunFolderTree is the main entry point for gitmap folder-tree / ft.
func RunFolderTree(args []string) error {
	opts, remaining := parseFolderTreeArgs(args)
	if opts.Subcommand == "help" {
		PrintFolderTreeHelp()
		return nil
	}
	switch opts.Subcommand {
	case "export", "exp", "export-folder":
		return handleExport(opts, remaining)
	case "import", "imp", "import-folder":
		return handleImport(opts, remaining)
	default:
		return handleLS(opts, remaining)
	}
}

func handleLS(opts FolderTreeOptions, remaining []string) error {
	if len(remaining) > 0 && isFolderArg(remaining[0]) {
		opts.DirsOnly = true
		remaining = remaining[1:]
	}
	target := resolveTargetPath(remaining, opts.TargetDir)
	root, err := ScanFolderTree(target, opts)
	if err != nil {
		return apperror.WrapSimple(err, "scan folder tree "+target)
	}
	if opts.IsPreview {
		fmt.Print(RenderPreview(root, opts))
		return nil
	}
	fmt.Print(RenderTree(root, opts))
	return nil
}

func handleExport(opts FolderTreeOptions, remaining []string) error {
	if len(remaining) > 0 && isFolderArg(remaining[0]) {
		opts.DirsOnly = true
		remaining = remaining[1:]
	}
	if opts.Subcommand == "export-folder" {
		opts.DirsOnly = true
	}
	target := resolveTargetPath(remaining, opts.TargetDir)
	root, err := ScanFolderTree(target, opts)
	if err != nil {
		return apperror.WrapSimple(err, "scan folder tree "+target)
	}
	content, err := ExportFolderTree(root, opts)
	if err != nil {
		return apperror.WrapSimple(err, "export folder tree")
	}
	if opts.OutputFile == "" {
		fmt.Println(content)
	} else {
		fmt.Printf("✔ Exported folder tree to %s\n", opts.OutputFile)
	}
	return nil
}

func handleImport(opts FolderTreeOptions, remaining []string) error {
	if len(remaining) > 0 && isFolderArg(remaining[0]) {
		opts.DirsOnly = true
		remaining = remaining[1:]
	}
	if opts.Subcommand == "import-folder" {
		opts.DirsOnly = true
	}
	if len(remaining) == 0 && opts.InputFile == "" {
		return apperror.NewValidationError("usage: gitmap folder-tree import <file.json|file.yaml|file.txt> [target-dir]")
	}
	inputFile := opts.InputFile
	targetDir := "."
	if len(remaining) > 0 {
		inputFile = remaining[0]
	}
	if len(remaining) > 1 {
		targetDir = remaining[1]
	}
	summary, err := ImportFolderTree(inputFile, targetDir, opts)
	if err != nil {
		return apperror.WrapSimple(err, "import folder tree")
	}
	printImportSummary(summary, opts.DryRun)
	return nil
}
