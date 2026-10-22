package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// annotation is one curated catalog entry from nc_annotations.json.
// Name, Version, Category, Description, and Example are the editorial
// copy; Alias is the display alias used when the name is a multi-word
// form ("te all") or the base command has no dispatcher registration.
type annotation struct {
	Name        string `json:"name"`
	Alias       string `json:"alias"`
	Version     string `json:"version"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Example     string `json:"example"`
}

// loadAnnotations reads one JSON object per line and validates that every
// entry carries the fields the catalog tests require.
func loadAnnotations(path string) ([]annotation, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read annotations")
	}

	var out []annotation

	dec := json.NewDecoder(strings.NewReader(string(data)))
	for dec.More() {
		var a annotation
		if err := dec.Decode(&a); err != nil {
			return nil, apperror.WrapSimple(err, "decode annotations")
		}

		if err := validateAnnotation(a); err != nil {
			return nil, err
		}

		out = append(out, a)
	}

	if len(out) == 0 {
		return nil, apperror.NewNotFoundError(fmt.Sprintf("no annotations found in %s", path))
	}

	return out, nil
}

func validateAnnotation(a annotation) error {
	missing := []string{}

	if a.Name == "" {
		missing = append(missing, "name")
	}

	if a.Version == "" {
		missing = append(missing, "version")
	}

	if a.Category == "" {
		missing = append(missing, "category")
	}

	if a.Description == "" {
		missing = append(missing, "description")
	}

	if a.Example == "" {
		missing = append(missing, "example")
	}

	if len(missing) > 0 {
		return apperror.NewValidationError(fmt.Sprintf("annotation %q missing fields: %s", a.Name, strings.Join(missing, ", ")))
	}

	return nil
}
