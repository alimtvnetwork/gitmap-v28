package cmdfoldertree

import (
	"fmt"
)

// RunFolderTree is the main entry point for gitmap folder-tree / ft.
func RunFolderTree(args []string) error {
	opts, remaining := parseFolderTreeArgs(args)
	if opts.Subcommand == "help" {
		PrintFolderTreeHelp()
		return nil
	}
	switch opts.Subcommand {
	case "export", "exp":
		return handleExport(opts, remaining)
	case "import", "imp":
		return handleImport(opts, remaining)
	default:
		return handleLS(opts, remaining)
	}
}

func handleLS(opts FolderTreeOptions, remaining []string) error {
	target := resolveTargetPath(remaining, opts.TargetDir)
	root, err := ScanFolderTree(target, opts)
	if err != nil {
		return fmt.Errorf("scan folder tree %s: %w", target, err)
	}
	if opts.IsPreview {
		fmt.Print(RenderPreview(root, opts))
		return nil
	}
	fmt.Print(RenderTree(root, opts))
	return nil
}

func handleExport(opts FolderTreeOptions, remaining []string) error {
	target := resolveTargetPath(remaining, opts.TargetDir)
	root, err := ScanFolderTree(target, opts)
	if err != nil {
		return fmt.Errorf("scan folder tree %s: %w", target, err)
	}
	content, err := ExportFolderTree(root, opts)
	if err != nil {
		return fmt.Errorf("export folder tree: %w", err)
	}
	if opts.OutputFile == "" {
		fmt.Println(content)
	} else {
		fmt.Printf("✔ Exported folder tree to %s\n", opts.OutputFile)
	}
	return nil
}

func handleImport(opts FolderTreeOptions, remaining []string) error {
	if len(remaining) == 0 && opts.InputFile == "" {
		return fmt.Errorf("usage: gitmap folder-tree import <file.json|file.yaml> [target-dir]")
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
		return fmt.Errorf("import folder tree: %w", err)
	}
	printImportSummary(summary, opts.DryRun)
	return nil
}

func printImportSummary(s *ImportSummary, dryRun bool) {
	prefix := "✔ Created"
	if dryRun {
		prefix = "✔ [DRY RUN] Would create"
	}
	fmt.Printf("%s %d directories and %d files in %s\n",
		prefix, s.DirsCreated, s.FilesCreated, s.TargetDir)
}

func resolveTargetPath(remaining []string, targetDir string) string {
	if len(remaining) > 0 && remaining[0] != "" {
		return remaining[0]
	}
	if targetDir != "" {
		return targetDir
	}
	return "."
}
