package templates

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/templates/enums"

// SourceType describes where a resolved template came from.
// SourceType aliases the strongly-typed enum from the enums package.
type SourceType = enums.SourceType

const (
	SourceNone  = enums.SourceNone
	SourceUser  = enums.SourceUser
	SourceEmbed = enums.SourceEmbed
)

// Resolved is a single template resolution result.
type Resolved struct {
	Kind    string
	Lang    string
	Path    string // overlay absolute path, or embedded virtual path
	Source  SourceType
	Content []byte
}

// Resolve looks up a (kind, lang) template, preferring the user overlay
// over the embedded asset.
func Resolve(kind, lang string) (Resolved, error) {
	rel := relPath(kind, lang)

	if r, ok, err := resolveUser(kind, lang, rel); err != nil {
		return Resolved{}, err
	} else if ok {
		return r, nil
	}

	return resolveEmbed(kind, lang, rel)
}

// relPath returns "<kind>/<lang><ext>" for a given (kind, lang).
func relPath(kind, lang string) string {
	return filepath.ToSlash(filepath.Join(kind, lang+extFor(kind)))
}

// extFor returns the file extension for a given template kind.
func extFor(kind string) string {
	if kind == kindAttributes || kind == kindLFS {
		return templateExtAttributes
	}

	return templateExtIgnore
}

// resolveUser checks the overlay directory.
func resolveUser(kind, lang, rel string) (Resolved, bool, error) {
	dir, err := UserDir()
	if err != nil {
		return Resolved{}, false, err
	}

	full := filepath.Join(dir, rel)
	data, err := os.ReadFile(full)
	if errors.Is(err, fs.ErrNotExist) {
		return Resolved{}, false, nil
	}

	if err != nil {
		return Resolved{}, false, apperror.WrapSimple(err, "failed to read template "+full)
	}

	return Resolved{Kind: kind, Lang: lang, Path: full, Source: SourceUser, Content: data}, true, nil
}

// resolveEmbed reads from the embedded FS.
func resolveEmbed(kind, lang, rel string) (Resolved, error) {
	full := filepath.ToSlash(filepath.Join(embedAssetsRoot, rel))
	data, err := FS.ReadFile(full)
	if errors.Is(err, fs.ErrNotExist) {
		return Resolved{}, apperror.NewNotFoundError(fmt.Sprintf(errTemplateNotFound, kind, lang))
	}

	if err != nil {
		return Resolved{}, apperror.WrapSimple(err, "failed to read template "+full)
	}

	return Resolved{Kind: kind, Lang: lang, Path: full, Source: SourceEmbed, Content: data}, nil
}
