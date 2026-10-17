package cmddb

// TableInfo represents metadata about a SQLite table or view.
type TableInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	ColumnCount int    `json:"columnCount"`
	RowCount    int64  `json:"rowCount"`
}

// ColumnInfo represents a column definition extracted via PRAGMA table_info.
type ColumnInfo struct {
	Cid          int     `json:"cid"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	NotNull      bool    `json:"notNull"`
	DefaultValue *string `json:"defaultValue"`
	IsPrimaryKey bool    `json:"isPrimaryKey"`
}

// SchemaInfo represents DDL definition and associated index definitions for a table.
type SchemaInfo struct {
	TableName string   `json:"tableName"`
	Type      string   `json:"type"`
	SQL       string   `json:"sql"`
	Indexes   []string `json:"indexes,omitempty"`
}

// QueryResult represents tabular query results with column names and dynamic rows.
type QueryResult struct {
	Columns []string `json:"columns"`
	Rows    [][]any  `json:"rows"`
	Count   int      `json:"count"`
}

// CommonDbOptions represents shared flags parsed across db subcommands.
type CommonDbOptions struct {
	DbPath string
	IsJson bool
	Limit  int
}
