package cmdfoldertree

import "strconv"

func handleValueFlags(a string, args []string, i int, opts *FolderTreeOptions) (bool, int) {
	if a == "--depth" || a == "-L" {
		return extractDepthFlag(args, i, opts)
	}
	if a == "--format" || a == "-fmt" || (a == "-f" && hasFormatValue(args, i)) {
		return extractFormatFlag(args, i, opts)
	}
	if a == "--output" || a == "-o" || a == "--file" {
		return extractOutputFlag(args, i, opts)
	}
	return false, i
}

func hasFormatValue(args []string, i int) bool {
	if i+1 >= len(args) {
		return false
	}
	switch args[i+1] {
	case "json", "yaml", "yml", "tree", "txt", "text", "preview", "gap", "folder", "folders", "dirs", "files":
		return true
	default:
		return false
	}
}

func extractDepthFlag(args []string, i int, opts *FolderTreeOptions) (bool, int) {
	if i+1 >= len(args) {
		return false, i
	}
	if d, err := strconv.Atoi(args[i+1]); err == nil {
		opts.MaxDepth = d
	}
	return true, i + 1
}

func extractFormatFlag(args []string, i int, opts *FolderTreeOptions) (bool, int) {
	if i+1 < len(args) {
		opts.Format = args[i+1]
		return true, i + 1
	}
	return false, i
}

func extractOutputFlag(args []string, i int, opts *FolderTreeOptions) (bool, int) {
	if i+1 < len(args) {
		opts.OutputFile = args[i+1]
		return true, i + 1
	}
	return false, i
}
