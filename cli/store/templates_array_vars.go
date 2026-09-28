package store

import (
	"encoding/json"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	varArrayRefPattern = regexp.MustCompile(`^([A-Za-z0-9_.-]+)\[(\d+)\]$`)
	seededRand         = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// ResolveTemplateVarValue resolves a placeholder key against vars supporting arrays, indexing, and PascalCase normalization.
func ResolveTemplateVarValue(key string, vars map[string]string) (string, bool) {
	if len(vars) == 0 || key == "" {
		return "", false
	}

	// 1. Check indexed array reference: e.g. Var[0]
	if match := varArrayRefPattern.FindStringSubmatch(key); len(match) == 3 {
		return tryResolveIndexedArray(match[1], match[2], vars)
	}

	// 2. Standard or random array reference: e.g. Var
	rawVal, isFound := lookupKeyWithNormalization(key, vars)
	if !isFound {
		return "", false
	}

	elements, isArray := parseArraySegments(rawVal)
	if isArray && len(elements) > 0 {
		return pickRandomArrayElement(elements), true
	}

	return rawVal, true
}

func tryResolveIndexedArray(baseKey, idxStr string, vars map[string]string) (string, bool) {
	idx, err := strconv.Atoi(idxStr)
	if err != nil {
		return "", false
	}
	return resolveIndexedArrayValue(baseKey, idx, vars)
}

func resolveIndexedArrayValue(baseKey string, idx int, vars map[string]string) (string, bool) {

	rawVal, isFound := lookupKeyWithNormalization(baseKey, vars)
	if !isFound {
		return "", false
	}
	elements, isArray := parseArraySegments(rawVal)
	if !isArray || len(elements) == 0 {
		return rawVal, true
	}
	if idx < 0 {
		return elements[0], true
	}
	if idx >= len(elements) {
		return elements[len(elements)-1], true
	}
	return elements[idx], true
}

func lookupKeyWithNormalization(key string, vars map[string]string) (string, bool) {
	if val, ok := vars[key]; ok {
		return val, true
	}
	for k, v := range vars {
		if strings.EqualFold(k, key) {
			return v, true
		}
	}
	normKey := normalizeVarKey(key)
	for k, v := range vars {
		if normalizeVarKey(k) == normKey {
			return v, true
		}
	}
	return "", false
}

func normalizeVarKey(k string) string {
	clean := strings.ToLower(k)
	clean = strings.ReplaceAll(clean, "_", "")
	clean = strings.ReplaceAll(clean, "-", "")
	clean = strings.ReplaceAll(clean, ".", "")
	return clean
}

func parseArraySegments(raw string) ([]string, bool) {
	trimmed := strings.TrimSpace(raw)
	if !strings.HasPrefix(trimmed, "[") || !strings.HasSuffix(trimmed, "]") {
		return nil, false
	}
	var strSlice []string
	if err := json.Unmarshal([]byte(trimmed), &strSlice); err == nil {
		return strSlice, true
	}
	var anySlice []any
	if err := json.Unmarshal([]byte(trimmed), &anySlice); err == nil {
		res := make([]string, 0, len(anySlice))
		for _, item := range anySlice {
			res = append(res, strings.TrimSpace(itemString(item)))
		}
		return res, true
	}
	return nil, false
}

func itemString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func pickRandomArrayElement(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	if len(arr) == 1 {
		return arr[0]
	}
	n := seededRand.Intn(len(arr))
	return arr[n]
}
