package termout

// AlignmentType defines column text alignment for terminal tables.
type AlignmentType string

const (
	// AlignLeft positions text flush with the left margin.
	AlignLeft AlignmentType = "left"
	// AlignRight positions text flush with the right margin (ideal for numeric metrics).
	AlignRight AlignmentType = "right"
	// AlignCenter centers text symmetrically within the column width.
	AlignCenter AlignmentType = "center"
)

// AlignType is maintained as a type alias for AlignmentType for backward compatibility.
type AlignType = AlignmentType

// BorderStyleType defines the visual framing style for rendered terminal tables.
type BorderStyleType string

const (
	// BorderNone indicates no border framing (standard spaced table).
	BorderNone BorderStyleType = "none"
	// BorderRounded uses curved UTF-8 box-drawing characters (╭, ╮, ╯, ╰).
	BorderRounded BorderStyleType = "rounded"
	// BorderSharp uses classic rectangular box-drawing characters (┌, ┐, ┘, └).
	BorderSharp BorderStyleType = "sharp"
	// BorderMinimal uses understated separator lines without outer vertical bars.
	BorderMinimal BorderStyleType = "minimal"
)

// ColumnConfig defines configuration options for an individual table column.
type ColumnConfig struct {
	Title     string
	MaxWidth  int
	MinWidth  int
	Alignment AlignmentType
	Align     AlignmentType // Backward-compatible field alias
}

// Column is maintained as a type alias for ColumnConfig for backward compatibility.
type Column = ColumnConfig

// Row defines a single row containing column cell data.
type Row struct {
	Cells []string
	Color string
}

// TableConfig defines full configuration for table rendering.
type TableConfig struct {
	Columns      []ColumnConfig
	Rows         []Row
	HeaderColor  string
	BorderColor  string
	HasBorders   bool
	BorderStyle  BorderStyleType
	EllipsisText string
}
