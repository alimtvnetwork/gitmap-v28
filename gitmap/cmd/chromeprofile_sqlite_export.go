package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// exportChromeSQLite extracts the portable SQLite databases from a Chrome profile
// into a destination directory. It returns the number of files copied and any error.
//
//nolint:unused
func exportChromeSQLite(srcProfile, dstDir string) (int, error) {
	if err := os.MkdirAll(dstDir, constants.DirPermission); err != nil {
		return 0, fmt.Errorf("mkdir %s: %w", dstDir, err)
	}

	total := 0
	for _, name := range constants.ChromeProfileSQLiteEntries {
		src := filepath.Join(srcProfile, name)
		dst := filepath.Join(dstDir, name)
		n, err := copyEntry(src, dst)
		if err != nil {
			// locked files or missing files are often skipped gracefully by copyEntry,
			// but if it's a hard error, we return it.
			return total, fmt.Errorf("copy %s: %w", name, err)
		}
		total += n
	}
	return total, nil
}
