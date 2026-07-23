package cmd

// JSON encoder for `gitmap bookmark list --json`.
//
// Migrated off json.MarshalIndent onto gitmap/stablejson so key order
// becomes a compile-time decision rather than a reflection accident.
// Schema: spec/08-json-schemas/bookmark-list.schema.json.

import (
	"io"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/stablejson"
)

// bookmark-list wire keys. Names + order are the contract.
const (
	bookmarkListKeyID        = "id"
	bookmarkListKeyName      = "name"
	bookmarkListKeyCommand   = "command"
	bookmarkListKeyArgs      = "args"
	bookmarkListKeyFlags     = "flags"
	bookmarkListKeyCreatedAt = "createdAt"
)

// encodeBookmarkListJSON writes records as a stablejson 2-space-indented
// array. Empty input emits `[]\n`.
func encodeBookmarkListJSON(w io.Writer, records []model.BookmarkRecord) error {
	return stablejson.WriteArray(w, buildBookmarkListJSONItems(records))
}

// buildBookmarkListJSONItems is the single source of (field name,
// field order, value) for bookmark-list.
func buildBookmarkListJSONItems(records []model.BookmarkRecord) [][]stablejson.Field {
	items := make([][]stablejson.Field, 0, len(records))
	for _, r := range records {
		items = append(items, []stablejson.Field{
			{Key: bookmarkListKeyID, Value: r.ID},
			{Key: bookmarkListKeyName, Value: r.Name},
			{Key: bookmarkListKeyCommand, Value: r.Command},
			{Key: bookmarkListKeyArgs, Value: r.Args},
			{Key: bookmarkListKeyFlags, Value: r.Flags},
			{Key: bookmarkListKeyCreatedAt, Value: r.CreatedAt},
		})
	}

	return items
}
