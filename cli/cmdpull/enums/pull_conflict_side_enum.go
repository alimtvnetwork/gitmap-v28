package enums

// ConflictSideType represents which side of a merge conflict to checkout.
type ConflictSideType string

const (
	ConflictSideTheirs ConflictSideType = "theirs"
	ConflictSideOurs   ConflictSideType = "ours"
)

// IsValid reports whether the conflict side is either theirs or ours.
func (c ConflictSideType) IsValid() bool {
	return c == ConflictSideTheirs || c == ConflictSideOurs
}

// String returns the string value of the conflict side.
func (c ConflictSideType) String() string {
	return string(c)
}
