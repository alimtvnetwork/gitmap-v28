package termout

import "github.com/alimtvnetwork/gitmap-v28/cli/termout/enums"

// AlignmentType defines column text alignment for terminal tables.
// AlignmentType aliases the strongly-typed enum from the enums package.
type AlignmentType = enums.AlignmentType

const (
	AlignLeft   = enums.AlignLeft
	AlignRight  = enums.AlignRight
	AlignCenter = enums.AlignCenter
)

// AlignType is maintained as a type alias for AlignmentType for backward compatibility.
type AlignType = enums.AlignmentType

// BorderStyleType defines the visual framing style for rendered terminal tables.
// BorderStyleType aliases the strongly-typed enum from the enums package.
type BorderStyleType = enums.BorderStyleType

const (
	BorderNone    = enums.BorderNone
	BorderRounded = enums.BorderRounded
	BorderSharp   = enums.BorderSharp
	BorderMinimal = enums.BorderMinimal
)

// WrapModeType defines cell text wrapping behavior for terminal tables.
// WrapModeType aliases the strongly-typed enum from the enums package.
type WrapModeType = enums.WrapModeType

const (
	WrapClip     = enums.WrapClip
	WrapEllipsis = enums.WrapEllipsis
	WrapWord     = enums.WrapWord
)

// ColumnConfig defines configuration options for an individual table column.
type ColumnConfig struct {
	Title     string
	MaxWidth  int
	MinWidth  int
	Alignment AlignmentType
	Align     AlignmentType // Backward-compatible field alias
	WrapMode  WrapModeType
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
	StripeColor  string
	HasBorders   bool
	RowStripe    bool
	BorderStyle  BorderStyleType
	WrapMode     WrapModeType
	EllipsisText string
}
