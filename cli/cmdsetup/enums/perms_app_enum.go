package enums

// PermsAppType defines the web application profile for permission auditing.
type PermsAppType string

const (
	PermsAppTypeGeneric   PermsAppType = "generic"
	PermsAppTypeWordPress PermsAppType = "wordpress"
	PermsAppTypeLaravel   PermsAppType = "laravel"
)

// IsValid reports whether the value matches a known variant.
func (p PermsAppType) IsValid() bool {
	switch p {
	case PermsAppTypeGeneric, PermsAppTypeWordPress, PermsAppTypeLaravel:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (p PermsAppType) String() string {
	return string(p)
}

// AllPermsAppTypes returns all valid variants.
func AllPermsAppTypes() []PermsAppType {
	return []PermsAppType{
		PermsAppTypeGeneric,
		PermsAppTypeWordPress,
		PermsAppTypeLaravel,
	}
}

// IsGeneric reports whether the value is PermsAppTypeGeneric.
func (p PermsAppType) IsGeneric() bool {
	return p == PermsAppTypeGeneric
}

// IsWordPress reports whether the value is PermsAppTypeWordPress.
func (p PermsAppType) IsWordPress() bool {
	return p == PermsAppTypeWordPress
}

// IsLaravel reports whether the value is PermsAppTypeLaravel.
func (p PermsAppType) IsLaravel() bool {
	return p == PermsAppTypeLaravel
}

// HasApp reports whether the value is non-empty and valid.
func (p PermsAppType) HasApp() bool {
	return p.IsValid()
}
