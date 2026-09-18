package termtable

// AlignType defines column text alignment.
type AlignType string

const (
	AlignLeft   AlignType = "left"
	AlignRight  AlignType = "right"
	AlignCenter AlignType = "center"
)

// Column defines a table column specification.
type Column struct {
	Title    string
	MaxWidth int
	MinWidth int
	Align    AlignType
}

// Row defines a single row containing column cell data.
type Row struct {
	Cells []string
	Color string
}

// TableConfig defines full configuration for table rendering.
type TableConfig struct {
	Columns      []Column
	Rows         []Row
	HeaderColor  string
	BorderColor  string
	HasBorders   bool
	EllipsisText string
}
