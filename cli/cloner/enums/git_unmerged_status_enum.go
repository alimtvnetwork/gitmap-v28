package enums

// GitUnmergedStatusType represents two-letter git status porcelain codes for unmerged files.
type GitUnmergedStatusType string

const (
	StatusBothModified  GitUnmergedStatusType = "UU"
	StatusBothAdded     GitUnmergedStatusType = "AA"
	StatusUpdatedByUs   GitUnmergedStatusType = "UD"
	StatusUpdatedByThem GitUnmergedStatusType = "DU"
	StatusBothDeleted   GitUnmergedStatusType = "DD"
)

// IsValid reports whether the status is a recognized two-letter unmerged conflict status.
func (s GitUnmergedStatusType) IsValid() bool {
	switch s {
	case StatusBothModified, StatusBothAdded, StatusUpdatedByUs, StatusUpdatedByThem, StatusBothDeleted:
		return true
	default:
		return false
	}
}

// String returns the string representation.
func (s GitUnmergedStatusType) String() string {
	return string(s)
}
