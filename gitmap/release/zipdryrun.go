package release

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// DryRunZipGroups prints what zip groups would produce without creating them.
func DryRunZipGroups(db *store.DB, groupNames []string) {
	if len(groupNames) == 0 {
		return
	}

	fmt.Printf(constants.MsgZGDryRunHeader, len(groupNames))

	for _, name := range groupNames {
		items, err := db.ListZipGroupItems(name)
		if err != nil {
			continue
		}

		paths := make([]string, len(items))
		for i, item := range items {
			paths[i] = item.FullPath
			if len(paths[i]) == 0 {
				paths[i] = item.Path
			}
		}

		group, _ := db.FindZipGroupByName(name)
		archiveName := resolveArchiveName(group)

		fmt.Printf(constants.MsgZGDryRunEntry, archiveName, len(items), strings.Join(paths, ", "))
	}
}

// DryRunAdHoc prints what ad-hoc zip items would produce without creating them.
func DryRunAdHoc(paths []string, bundleName string) {
	if len(paths) == 0 {
		return
	}

	if len(bundleName) > 0 {
		fmt.Printf(constants.MsgZGDryRunEntry, bundleName, len(paths), strings.Join(paths, ", "))

		return
	}

	for _, p := range paths {
		base := filepath.Base(p)
		archiveName := strings.TrimSuffix(base, filepath.Ext(base)) + ".zip"

		fmt.Printf(constants.MsgZGDryRunEntry, archiveName, 1, p)
	}
}
