package templates

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Source describes where a resolved template came from.
type Source int

const (
	SourceNone Source = iota
	SourceUser
	SourceEmbed
)

// Resolved is a single template resolution result.
type Resolved struct {
	Kind    string
	Lang    string
	Path    string // overlay absolute path, or embedded virtual path
	Source  Source
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
		return Resolved{}, false, fmt.Errorf(errTemplateRead, full, err)
	}

	return Resolved{Kind: kind, Lang: lang, Path: full, Source: SourceUser, Content: data}, true, nil
}

// resolveEmbed reads from the embedded FS.
func resolveEmbed(kind, lang, rel string) (Resolved, error) {
	full := filepath.ToSlash(filepath.Join(embedAssetsRoot, rel))
	data, err := FS.ReadFile(full)
	if errors.Is(err, fs.ErrNotExist) {
		return Resolved{}, fmt.Errorf(errTemplateNotFound, kind, lang)
	}
	if err != nil {
		return Resolved{}, fmt.Errorf(errTemplateRead, full, err)
	}

	return Resolved{Kind: kind, Lang: lang, Path: full, Source: SourceEmbed, Content: data}, nil
}
