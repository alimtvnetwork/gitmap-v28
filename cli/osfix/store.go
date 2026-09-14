package osfix

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// ListFixes loads all registered fixes from the persistent store.
func ListFixes() FixSliceResult {
	path := resolveFixStorePath()
	if !hasFile(path) {
		return result.OkSlice([]FixItem{})
	}
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return result.FailSlice[FixItem](apperror.WrapSimple(readErr, "read fixes store"))
	}
	var items []FixItem
	if err := json.Unmarshal(content, &items); err != nil {
		return result.FailSlice[FixItem](apperror.WrapSimple(err, "unmarshal fixes store"))
	}
	return result.OkSlice(items)
}

// GetFix retrieves a fix by name.
func GetFix(name string) FixResult {
	res := ListFixes()
	if res.IsFailure() {
		return result.Fail[FixItem](res.AppError())
	}
	for _, item := range res.Value {
		if item.Name == name {
			return result.Ok(item)
		}
	}
	return result.Fail[FixItem](apperror.NewSimple("fix not found: "+name, "E_FIX_NOT_FOUND"))
}

// SaveFix adds or edits a registered fix.
func SaveFix(item FixItem) FixResult {
	res := ListFixes()
	if res.IsFailure() {
		return result.Fail[FixItem](res.AppError())
	}
	items := res.Value
	isUpdated := false
	for i, existing := range items {
		if existing.Name == item.Name {
			item.CreatedAt = existing.CreatedAt
			item.UpdatedAt = time.Now()
			items[i] = item
			isUpdated = true
			break
		}
	}
	if !isUpdated {
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()
		items = append(items, item)
	}
	return persistFixItems(items, item)
}

func persistFixItems(items []FixItem, target FixItem) FixResult {
	path := resolveFixStorePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return result.Fail[FixItem](apperror.WrapSimple(err, "create fix dir"))
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return result.Fail[FixItem](apperror.WrapSimple(err, "marshal fixes"))
	}
	if writeErr := os.WriteFile(path, data, 0o644); writeErr != nil {
		return result.Fail[FixItem](apperror.WrapSimple(writeErr, "write fixes store"))
	}
	return result.Ok(target)
}

// DeleteFix removes a fix by name.
func DeleteFix(name string) FixBoolResult {
	res := ListFixes()
	if res.IsFailure() {
		return result.Fail[bool](res.AppError())
	}
	var kept []FixItem
	isFound := false
	for _, item := range res.Value {
		if item.Name == name {
			isFound = true
			continue
		}
		kept = append(kept, item)
	}
	if !isFound {
		return result.Fail[bool](apperror.NewSimple("fix not found: "+name, "E_FIX_NOT_FOUND"))
	}
	return writeKeptFixes(kept)
}

func writeKeptFixes(items []FixItem) FixBoolResult {
	path := resolveFixStorePath()
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return result.Fail[bool](apperror.WrapSimple(err, "marshal fixes"))
	}
	if writeErr := os.WriteFile(path, data, 0o644); writeErr != nil {
		return result.Fail[bool](apperror.WrapSimple(writeErr, "write fixes store"))
	}
	return result.Ok(true)
}

func resolveFixStorePath() string {
	home := os.Getenv("USERPROFILE")
	if home == "" {
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, ".gitmap", "os_fixes.json")
}

func hasFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
