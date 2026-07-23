package cmd

import (
	"archive/zip"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// chooseDocsExtractTarget peeks into the zip to decide where to extract.
// If every entry is already prefixed with `docs-site/`, the zip is
// "wrapped" and must extract to binaryDir so the prefix lands as the
// docs-site/ folder. Otherwise (flat zip containing dist/, index.html,
// etc.) extract directly into docsDir so the files end up inside it.
func chooseDocsExtractTarget(zipPath, binaryDir, docsDir string) string {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return docsDir
	}
	defer r.Close()

	prefix := constants.HDDocsDir + "/"
	for _, f := range r.File {
		name := strings.ReplaceAll(f.Name, "\\", "/")
		if name == "" || strings.HasPrefix(name, "__MACOSX/") {
			continue
		}
		if !strings.HasPrefix(name, prefix) && name != constants.HDDocsDir {
			return docsDir
		}
	}

	return binaryDir
}
