// Package release — ziparchive.go creates ZIP archives from zip groups
// with maximum compression (Deflate level 9) for release assets.
package release

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/verbose"
)

// ZipGroupArchive holds the result of archiving a zip group.
type ZipGroupArchive struct {
	GroupName   string
	ArchivePath string
	ItemCount   int
}

// BuildZipGroupArchives resolves persistent zip groups from the DB and
// creates max-compression ZIP archives for each. Returns archive paths.
func BuildZipGroupArchives(db *store.DB, groupNames []string, stagingDir string) []string {
	archives := make([]string, 0, len(groupNames))

	for _, name := range groupNames {
		archive, err := buildOneZipGroup(db, name, stagingDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrZGCompress, name, err)

			continue
		}

		archives = append(archives, archive)
	}

	return archives
}

// buildOneZipGroup loads a group's items and compresses them.
func buildOneZipGroup(db *store.DB, name, stagingDir string) (string, error) {
	group, err := db.FindZipGroupByName(name)
	if err != nil {
		return "", fmt.Errorf(constants.ErrZGGroupNotDB, name)
	}

	items, err := db.ListZipGroupItems(name)
	if err != nil {
		return "", err
	}

	if len(items) == 0 {
		fmt.Printf(constants.MsgZGSkipEmpty, name)

		return "", fmt.Errorf("empty group")
	}

	if verbose.IsEnabled() {
		verbose.Get().Log("zip-group %q: %d item(s)", name, len(items))
		for _, item := range items {
			p := item.FullPath
			if len(p) == 0 {
				p = item.Path
			}
			kind := "file"
			if item.IsFolder {
				kind = "folder"
			}
			verbose.Get().Log("  → %s (%s)", p, kind)
		}
	}

	archiveName := resolveArchiveName(group)
	archivePath := filepath.Join(stagingDir, archiveName)

	err = createMaxCompressZip(archivePath, items)
	if err != nil {
		return "", err
	}

	fmt.Printf(constants.MsgZGCompressed, name, archiveName)

	return archivePath, nil
}

// resolveArchiveName returns the archive filename for a group.
func resolveArchiveName(g model.ZipGroup) string {
	if len(g.ArchiveName) > 0 {
		return g.ArchiveName
	}

	return g.Name + ".zip"
}

// BuildAdHocArchive creates a ZIP from ad-hoc paths provided via -Z flags.
// If bundleName is set, all items go into one archive; otherwise each
// gets its own archive.
func BuildAdHocArchive(paths []string, bundleName, stagingDir string) []string {
	if len(bundleName) > 0 {
		return buildAdHocBundle(paths, bundleName, stagingDir)
	}

	return buildAdHocIndividual(paths, stagingDir)
}

// buildAdHocBundle bundles all ad-hoc paths into a single named archive.
func buildAdHocBundle(paths []string, bundleName, stagingDir string) []string {
	items := pathsToItems(paths)
	archivePath := filepath.Join(stagingDir, bundleName)

	err := createMaxCompressZip(archivePath, items)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrZGCompress, bundleName, err)

		return nil
	}

	fmt.Printf(constants.MsgZGCompressed, "ad-hoc", bundleName)

	return []string{archivePath}
}

// buildAdHocIndividual creates one archive per ad-hoc path.
func buildAdHocIndividual(paths []string, stagingDir string) []string {
	archives := make([]string, 0, len(paths))

	for _, p := range paths {
		base := filepath.Base(p)
		archiveName := strings.TrimSuffix(base, filepath.Ext(base)) + ".zip"
		archivePath := filepath.Join(stagingDir, archiveName)

		items := pathsToItems([]string{p})

		err := createMaxCompressZip(archivePath, items)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrZGCompress, p, err)

			continue
		}

		fmt.Printf(constants.MsgZGCompressed, p, archiveName)
		archives = append(archives, archivePath)
	}

	return archives
}

// pathsToItems converts raw paths to ZipGroupItem entries.
func pathsToItems(paths []string) []model.ZipGroupItem {
	items := make([]model.ZipGroupItem, 0, len(paths))

	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			fmt.Printf(constants.MsgZGSkipMissing, p)

			continue
		}

		items = append(items, model.ZipGroupItem{
			FullPath: p,
			Path:     p,
			IsFolder: info.IsDir(),
		})
	}

	return items
}
