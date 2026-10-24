package enums

// WhereClauseType identifies the kind of a query builder where clause.
type WhereClauseType int

const (
	WhereOp WhereClauseType = iota
	WhereLocate
	WhereColumnOp
)

// IsValid reports whether the where clause type is a recognized enumeration value.
func (w WhereClauseType) IsValid() bool {
	switch w {
	case WhereOp, WhereLocate, WhereColumnOp:
		return true
	default:
		return false
	}
}

// String returns the string representation of the where clause type.
func (w WhereClauseType) String() string {
	switch w {
	case WhereOp:
		return "op"
	case WhereLocate:
		return "locate"
	case WhereColumnOp:
		return "column_op"
	default:
		return "unknown"
	}
}

// AllWhereClauseTypes returns all valid WhereClauseType variants.
func AllWhereClauseTypes() []WhereClauseType {
	return []WhereClauseType{
		WhereOp,
		WhereLocate,
		WhereColumnOp,
	}
}

// IsOp reports whether the clause is a plain operator comparison.
func (w WhereClauseType) IsOp() bool {
	return w == WhereOp
}

// IsLocate reports whether the clause is a substring locate filter.
func (w WhereClauseType) IsLocate() bool {
	return w == WhereLocate
}

// IsColumnOp reports whether the clause is a column-to-column comparison.
func (w WhereClauseType) IsColumnOp() bool {
	return w == WhereColumnOp
}
