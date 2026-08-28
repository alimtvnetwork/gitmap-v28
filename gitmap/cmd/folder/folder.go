package folder

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"gopkg.in/yaml.v3"
)

// Run executes the folder command logic.
func Run(args []string) error {
	if len(args) < 1 {
		return apperror.NewSimple("FolderRun", "E_FOLDER_MISSING_DIR")
	}
	dir := args[0]
	outFile := "files.txt"
	if len(args) >= 2 && args[1] != "" {
		outFile = args[1]
	}

	excludePattern := ""
	if len(args) >= 4 && args[2] == "-exclude" {
		excludePattern = args[3]
		// Ignore the 0|1 for now, pattern is sufficient.
	}

	paths, err := walkAndFilter(dir, excludePattern)
	if err != nil {
		return apperror.WrapSimple(err, "folder: failed to walk directory")
	}

	ext := strings.ToLower(filepath.Ext(outFile))
	switch ext {
	case ".json":
		return writeJson(paths, outFile)
	case ".yaml", ".yml":
		return writeYaml(paths, outFile)
	case ".md", ".txt":
		fallthrough
	default:
		return writeText(paths, outFile)
	}
}

func walkAndFilter(root, excludeGlob string) ([]string, error) {
	var results []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return fs.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}
		if excludeGlob != "" {
			matched, _ := filepath.Match(excludeGlob, rel)
			if !matched {
				// Also try matching just the base name or direct path
				matched, _ = filepath.Match(excludeGlob, filepath.Base(rel))
			}
			if matched {
				return nil
			}
		}
		results = append(results, rel)
		return nil
	})
	return results, err
}

func writeJson(paths []string, out string) error {
	b, err := json.MarshalIndent(paths, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(out, b, 0644)
}

func writeYaml(paths []string, out string) error {
	b, err := yaml.Marshal(paths)
	if err != nil {
		return err
	}
	return os.WriteFile(out, b, 0644)
}

func writeText(paths []string, out string) error {
	var sb strings.Builder
	for _, p := range paths {
		sb.WriteString(p + "\n")
	}
	return os.WriteFile(out, []byte(sb.String()), 0644)
}
