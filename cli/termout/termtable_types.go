package termout

// AlignType defines column text alignment.
type AlignType string

const (
	AlignLeft   AlignType = "left"
	AlignRight  AlignType = "right"
	AlignCenter AlignType = "center"
)

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
	BorderStyle  BorderStyleType
	EllipsisText string
}
