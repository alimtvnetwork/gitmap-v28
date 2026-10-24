package templates

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// Entry describes one discoverable template.
type Entry struct {
	Kind   string     // ignore | attributes | lfs
	Lang   string     // common | go | node | ...
	Source SourceType // SourceUser or SourceEmbed
	Path   string     // absolute (overlay) or virtual embed path
}

// ListWrap returns every available template, with the user-overlay copy
// shadowing the embedded one when both exist. Sorted by (kind, lang), wrapped in EntrySliceWrap.
func ListWrap() EntrySliceWrap {
	merged := map[string]Entry{}

	if err := collectEmbed(merged); err != nil {
		return result.FailSlice[Entry](apperror.WrapSimple(err, "templates.ListWrap"))
	}

	if err := collectUser(merged); err != nil {
		return result.FailSlice[Entry](apperror.WrapSimple(err, "templates.ListWrap"))
	}

	out := make([]Entry, 0, len(merged))
	for _, e := range merged {
		out = append(out, e)
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Kind != out[j].Kind {
			return kindRank(out[i].Kind) < kindRank(out[j].Kind)
		}

		return out[i].Lang < out[j].Lang
	})

	return result.OkSlice(out)
}

// List returns every available template, with the user-overlay copy
// shadowing the embedded one when both exist. Sorted by (kind, lang).
func List() ([]Entry, error) {
	res := ListWrap()
	if res.HasError() {
		return nil, res.Err
	}

	return res.Items(), nil
}

// kindRank gives a stable display order: ignore, attributes, lfs.
func kindRank(kind string) int {
	switch kind {
	case kindIgnore:
		return 0
	case kindAttributes:
		return 1
	case kindLFS:
		return 2
	}

	return 3
}

// collectEmbed walks the embedded FS and adds every template.
func collectEmbed(out map[string]Entry) error {
	return fs.WalkDir(FS, embedAssetsRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if filepath.Base(p) == constants.ReadmeFile {
			return nil
		}

		kind, lang, ok := parseEmbedPath(p)
		if !ok {
			return nil
		}

		out[kind+"/"+lang] = Entry{Kind: kind, Lang: lang, Source: SourceEmbed, Path: p}

		return nil
	})
}

// collectUser walks the user-overlay directory (if it exists) and overrides
// any embed entries with overlay entries.
func collectUser(out map[string]Entry) error {
	dir, err := UserDir()
	if err != nil {
		return err
	}

	walkErr := fs.WalkDir(osDirFS(dir), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // overlay dir may not exist; skip silently
		}

		if d.IsDir() || p == "." {
			return nil
		}

		kind, lang, ok := parseRelTemplatePath(p)
		if !ok {
			return nil
		}

		out[kind+"/"+lang] = Entry{Kind: kind, Lang: lang, Source: SourceUser, Path: filepath.Join(dir, p)}

		return nil
	})

	return walkErr
}

// parseEmbedPath turns "assets/ignore/go.gitignore" into ("ignore", "go").
func parseEmbedPath(p string) (string, string, bool) {
	rel := strings.TrimPrefix(p, embedAssetsRoot+"/")

	return parseRelTemplatePath(rel)
}

// parseRelTemplatePath turns "ignore/go.gitignore" into ("ignore", "go").
func parseRelTemplatePath(rel string) (string, string, bool) {
	rel = filepath.ToSlash(rel)
	parts := strings.SplitN(rel, "/", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	kind := parts[0]
	base := parts[1]
	switch kind {
	case kindIgnore:
		if !strings.HasSuffix(base, templateExtIgnore) {
			return "", "", false
		}

		return kind, strings.TrimSuffix(base, templateExtIgnore), true
	case kindAttributes, kindLFS:
		if !strings.HasSuffix(base, templateExtAttributes) {
			return "", "", false
		}

		return kind, strings.TrimSuffix(base, templateExtAttributes), true
	}

	return "", "", false
}
