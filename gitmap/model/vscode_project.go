package model

// VSCodeProject is one row in the VSCodeProject table — the gitmap-side
// source of truth for entries synced into VS Code Project Manager's
// projects.json.
//
// `tags` is not stored on purpose: it lives only inside projects.json and
// is preserved across syncs.
//
// `Paths` (multi-root extras, schema v20+) IS stored, JSON-encoded as a
// TEXT column in SQLite, and surfaced here as a decoded []string. The DB
// list is UNIONed with any user-added paths from the VS Code UI on every
// sync — gitmap never silently removes a user-added entry.
type VSCodeProject struct {
	ID         int64    `json:"id"`
	RootPath   string   `json:"rootPath"`
	Name       string   `json:"name"`
	Paths      []string `json:"paths"`
	Enabled    bool     `json:"enabled"`
	Profile    string   `json:"profile"`
	LastSeenAt string   `json:"lastSeenAt"`
	CreatedAt  string   `json:"createdAt"`
	UpdatedAt  string   `json:"updatedAt"`
}
