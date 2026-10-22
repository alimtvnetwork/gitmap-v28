// Package fsutil — file_classify.go classifies files into read, edit, and created operations.
package fsutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// FileItem represents a classified file target with its anticipated operation.
type FileItem struct {
	Path   string `json:"path"`
	Action string `json:"action"` // "read", "edit", or "created"
}

// ClassifiedFiles represents files organized by category and itemized details.
type ClassifiedFiles struct {
	Read    []string   `json:"read"`
	Edit    []string   `json:"edit"`
	Created []string   `json:"created"`
	Items   []FileItem `json:"items"`
	Total   int        `json:"total"`
}

// NewClassifiedFiles initializes empty non-nil slices.
func NewClassifiedFiles() ClassifiedFiles {
	return ClassifiedFiles{
		Read:    []string{},
		Edit:    []string{},
		Created: []string{},
		Items:   []FileItem{},
		Total:   0,
	}
}

// ClassifyFileTargets parses raw file inputs and classifies them into read, edit, and created.
func ClassifyFileTargets(raw, defaultAction, rootDir string) ClassifiedFiles {
	result := NewClassifiedFiles()
	rawItems := extractRawFileList(raw)

	for _, rawItem := range rawItems {
		item := classifySingleFile(rawItem, defaultAction, rootDir)
		appendClassifiedItem(&result, item)
	}
	result.Total = len(result.Items)

	return result
}

func appendClassifiedItem(cf *ClassifiedFiles, item FileItem) {
	if item.Path == "" {
		return
	}
	cf.Items = append(cf.Items, item)

	switch item.Action {
	case "read":
		cf.Read = append(cf.Read, item.Path)
	case "created":
		cf.Created = append(cf.Created, item.Path)
	default:
		cf.Edit = append(cf.Edit, item.Path)
	}
}

func extractRawFileList(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed == "[]" {
		return nil
	}
	var jsonList []any
	isJsonList := json.Unmarshal([]byte(trimmed), &jsonList) == nil
	if isJsonList {
		return parseJsonListItems(jsonList)
	}

	return parseDelimitedItems(trimmed)
}

func parseJsonListItems(items []any) []string {
	var out []string
	for _, it := range items {
		str, isStr := it.(string)
		if isStr && strings.TrimSpace(str) != "" {
			out = append(out, strings.TrimSpace(str))
		}
	}

	return out
}

func parseDelimitedItems(raw string) []string {
	clean := strings.ReplaceAll(raw, ";", ",")
	clean = strings.ReplaceAll(clean, "\n", ",")
	parts := strings.Split(clean, ",")
	var out []string

	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}

	return out
}

func classifySingleFile(entry, defaultAction, rootDir string) FileItem {
	trimmed := strings.TrimSpace(entry)
	action, cleanPath, hasPrefix := extractPrefixAction(trimmed)
	if hasPrefix {
		return FileItem{Path: NormalizeToForwardSlashes(cleanPath), Action: action}
	}
	normalizedPath := NormalizeToForwardSlashes(trimmed)
	resolvedAction := resolveFallbackAction(normalizedPath, defaultAction, rootDir)

	return FileItem{Path: normalizedPath, Action: resolvedAction}
}

func extractPrefixAction(entry string) (string, string, bool) {
	prefixes := []struct {
		prefix string
		action string
	}{
		{"read:", "read"},
		{"view:", "read"},
		{"inspect:", "read"},
		{"create:", "created"},
		{"created:", "created"},
		{"new:", "created"},
		{"add:", "created"},
		{"edit:", "edit"},
		{"update:", "edit"},
		{"modify:", "edit"},
	}

	lower := strings.ToLower(entry)
	for _, p := range prefixes {
		hasPrefix := strings.HasPrefix(lower, p.prefix)
		if hasPrefix {
			path := strings.TrimSpace(entry[len(p.prefix):])

			return p.action, path, true
		}
	}

	return "", entry, false
}

func resolveFallbackAction(path, defaultAction, rootDir string) string {
	normalizedDefault := normalizeActionKeyword(defaultAction)
	if normalizedDefault != "" {
		return normalizedDefault
	}
	isFilePresent := checkFileExists(path, rootDir)
	if isFilePresent {
		return "edit"
	}

	return "created"
}

var actionKeywordCanonical = map[string]string{
	"read": "read", "view": "read", "inspect": "read",
	"create": "created", "created": "created", "new": "created", "add": "created",
	"edit": "edit", "update": "edit", "modify": "edit", "touch": "edit",
}

func normalizeActionKeyword(act string) string {
	return actionKeywordCanonical[strings.ToLower(strings.TrimSpace(act))]
}

func checkFileExists(path, rootDir string) bool {
	checkPath := path
	if rootDir != "" && !filepath.IsAbs(path) {
		checkPath = filepath.Join(rootDir, path)
	}
	info, err := os.Stat(checkPath)

	return err == nil && info != nil
}
