package cmddb

// DatabaseAlias represents a registered SQLite database alias mapping.
type DatabaseAlias struct {
	DatabaseAliasId    int64  `json:"databaseAliasId"`
	AliasName          string `json:"aliasName"`
	DatabasePath       string `json:"databasePath"`
	DatabaseType       string `json:"databaseType"`
	IsActive           bool   `json:"isActive"`
	FileSizeBytes      int64  `json:"fileSizeBytes"`
	SearchHistoryCount int64  `json:"searchHistoryCount"`
	LastAccessedAt     string `json:"lastAccessedAt"`
	Description        string `json:"description"`
	Notes              string `json:"notes"`
}

// DatabaseAliasSummary represents aggregated database listing statistics or an enriched listing entry.
type DatabaseAliasSummary struct {
	Alias          string `json:"alias"`
	Path           string `json:"path"`
	Type           string `json:"type,omitempty"`
	Size           string `json:"size"`
	FileSizeBytes  int64  `json:"fileSizeBytes"`
	TableCount     int    `json:"tableCount"`
	SearchCount    int64  `json:"searchCount"`
	IsActive       bool   `json:"isActive"`
	Description    string `json:"description,omitempty"`
	LastAccessedAt string `json:"lastAccessedAt,omitempty"`
}

// DatabaseListingPayload represents the JSON envelope payload for db-listing.
type DatabaseListingPayload struct {
	Databases  []DatabaseAliasSummary `json:"databases"`
	TotalCount int                    `json:"totalCount"`
	TotalBytes int64                  `json:"totalBytes"`
	TotalSize  string                 `json:"totalSize"`
}
