package constants

// CleanCategoryType defines categories for cleanup operations.
type CleanCategoryType string

// Backwards-compatible type alias.
type CleanCategory = CleanCategoryType

const (
	CleanCategoryInvalid  CleanCategoryType = ""
	CleanCategoryTerminal CleanCategoryType = "terminal"
	CleanCategoryDev      CleanCategoryType = "dev"
	CleanCategoryCache    CleanCategoryType = "cache"
)

// String returns the canonical serialized string representation.
func (c CleanCategoryType) String() string {
	return string(c)
}

// Label returns the human-readable display label.
func (c CleanCategoryType) Label() string {
	return c.String()
}

// IsValid reports whether the clean category is recognized.
func (c CleanCategoryType) IsValid() bool {
	hasCategory := c.IsTerminal() || c.IsDev() || c.IsCache()

	return hasCategory
}

// IsTerminal reports whether the category targets terminal history.
func (c CleanCategoryType) IsTerminal() bool {
	return c == CleanCategoryTerminal
}

// IsDev reports whether the category targets developer caches.
func (c CleanCategoryType) IsDev() bool {
	return c == CleanCategoryDev
}

// IsCache reports whether the category targets general cache.
func (c CleanCategoryType) IsCache() bool {
	return c == CleanCategoryCache
}
