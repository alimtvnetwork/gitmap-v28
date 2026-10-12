// Package cmdrelease implements release management commands.
package cmdrelease

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/release"
)

// PinnedTagsPayload represents the serialized schema for .gitmap/release/pinned-tags.json.
type PinnedTagsPayload struct {
	PinnedTags []string `json:"pinnedTags"`
	UpdatedAt  string   `json:"updatedAt"`
}

// PinnedTagsFileName is the canonical persistent pin file name.
const PinnedTagsFileName = "pinned-tags.json"

// PinnedTagsFilePath returns the resolved path to the persistent pin registry.
func PinnedTagsFilePath(workDir string) string {
	if workDir == "" {
		workDir = "."
	}
	return filepath.Join(workDir, constants.DefaultReleaseDir, PinnedTagsFileName)
}

// LoadPinnedTags reads and returns all persistent pinned tags.
func LoadPinnedTags(workDir string) ([]string, error) {
	path := PinnedTagsFilePath(workDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, apperror.WrapSimple(err, "E9010")
	}

	var payload PinnedTagsPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, apperror.WrapSimple(err, "E9011")
	}

	return payload.PinnedTags, nil
}

// SavePinnedTags deduplicates, sorts, and persists the pinned tags list.
func SavePinnedTags(workDir string, tags []string) error {
	path := PinnedTagsFilePath(workDir)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return apperror.WrapSimple(err, "E9012")
	}

	deduped := deduplicateAndSortTags(tags)
	payload := PinnedTagsPayload{
		PinnedTags: deduped,
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	bytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "E9013")
	}

	if err := os.WriteFile(path, append(bytes, '\n'), 0o644); err != nil {
		return apperror.WrapSimple(err, "E9014")
	}

	return nil
}

// IsTagPinned returns true if candidate matches any pinned tag.
func IsTagPinned(candidate string, pins []string) bool {
	normCandidate := normalizeTagForComparison(candidate)
	for _, pin := range pins {
		if normalizeTagForComparison(pin) == normCandidate {
			return true
		}
	}
	return false
}

func normalizeTagForComparison(tag string) string {
	t := strings.TrimSpace(tag)
	t = strings.TrimPrefix(t, "v")
	v, err := release.Parse(t)
	if err == nil {
		return v.CoreString()
	}
	return t
}

func deduplicateAndSortTags(tags []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, t := range tags {
		cleaned := strings.TrimSpace(t)
		if cleaned == "" {
			continue
		}
		if !seen[cleaned] {
			seen[cleaned] = true
			result = append(result, cleaned)
		}
	}
	sort.Strings(result)
	return result
}

// RunReleasePinNotToRemove implements `gitmap release pin-not-to-remove` (alias `pnr`).
func RunReleasePinNotToRemove(args []string) error {
	checkHelp(constants.CmdReleasePinNotToRemove, args)

	listFlag, removeFlag, clearFlag, rawTokens := parsePinFlags(args)
	workDir := "."

	existing, err := LoadPinnedTags(workDir)
	if err != nil {
		return err
	}

	if clearFlag {
		return handleClearPins(workDir)
	}

	if removeFlag != "" {
		return handleRemovePin(workDir, existing, removeFlag)
	}

	if listFlag || len(rawTokens) == 0 {
		printPinnedTagsList(existing)
		return nil
	}

	return handleAddPins(workDir, existing, rawTokens)
}

func parsePinFlags(args []string) (bool, string, bool, []string) {
	fs := flag.NewFlagSet(constants.CmdReleasePinNotToRemove, flag.ExitOnError)
	list := fs.Bool("list", false, "List all currently pinned tags")
	remove := fs.String("remove", "", "Remove a tag from the pin registry")
	unpin := fs.String("unpin", "", "Alias for --remove")
	clear := fs.Bool("clear", false, "Clear all pinned tags")

	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	targetRemove := *remove
	if targetRemove == "" {
		targetRemove = *unpin
	}

	var tokens []string
	for _, a := range fs.Args() {
		parts := strings.Split(a, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				tokens = append(tokens, trimmed)
			}
		}
	}

	return *list, targetRemove, *clear, tokens
}

func handleClearPins(workDir string) error {
	if err := SavePinnedTags(workDir, []string{}); err != nil {
		return err
	}
	fmt.Println("\x1b[1;32m✅ Cleared all pinned release tags.\x1b[0m")
	return nil
}

func handleRemovePin(workDir string, existing []string, toRemove string) error {
	norm := normalizeTagForComparison(toRemove)
	var kept []string
	found := false
	for _, pin := range existing {
		if normalizeTagForComparison(pin) == norm {
			found = true
			continue
		}
		kept = append(kept, pin)
	}

	if !found {
		fmt.Printf("\x1b[33m⚠ Tag %q was not found in the pinned tags registry.\x1b[0m\n", toRemove)
		return nil
	}

	if err := SavePinnedTags(workDir, kept); err != nil {
		return err
	}
	fmt.Printf("\x1b[1;32m✅ Unpinned tag %s (remaining: %d pins).\x1b[0m\n", toRemove, len(kept))
	return nil
}

func handleAddPins(workDir string, existing []string, tokens []string) error {
	merged := append(existing, tokens...)
	if err := SavePinnedTags(workDir, merged); err != nil {
		return err
	}

	fmt.Printf("\x1b[1;36m📌 release pin-not-to-remove\x1b[0m: Pinned %d tag(s)\n", len(tokens))
	for _, t := range tokens {
		fmt.Printf("   • \x1b[1;32m%s\x1b[0m (immune to remove-between and auto-pruning)\n", t)
	}

	updated, _ := LoadPinnedTags(workDir)
	fmt.Printf("   Total active pinned tags: %d\n", len(updated))
	return nil
}

func printPinnedTagsList(pins []string) {
	fmt.Println("\x1b[1;36m📌 Pinned Release Tags (Protected from removal):\x1b[0m")
	if len(pins) == 0 {
		fmt.Println("   (No tags currently pinned. Use `gitmap pnr <tag1>, <tag2>` to pin milestone releases)")
		return
	}
	for i, pin := range pins {
		fmt.Printf("   %2d. \x1b[1;32m%s\x1b[0m\n", i+1, pin)
	}
	fmt.Println("\x1b[2mUse `gitmap pnr --remove <tag>` to unpin a release tag.\x1b[0m")
}
