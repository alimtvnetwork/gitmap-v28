package templates

import (
	"io/fs"
	"os"
)

// osDirFS returns an fs.FS rooted at dir, falling back to an empty FS when
// dir doesn't exist. Lets List() treat a missing overlay as "no entries"
// without an error.
func osDirFS(dir string) fs.FS {
	if _, err := os.Stat(dir); err != nil {
		return emptyFS{}
	}

	return os.DirFS(dir)
}

type emptyFS struct{}

func (emptyFS) Open(_ string) (fs.File, error) { return nil, fs.ErrNotExist }
