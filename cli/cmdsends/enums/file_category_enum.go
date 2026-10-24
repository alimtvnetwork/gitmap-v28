package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// FileCategoryType defines semantic classification categories for changed files.
type FileCategoryType string

const (
	// CategoryStyle classifies CSS and styling assets.
	CategoryStyle FileCategoryType = "style"
	// CategoryCore classifies core programming code logic.
	CategoryCore FileCategoryType = "core"
	// CategoryBoolean classifies boolean conventions and predicate files.
	CategoryBoolean FileCategoryType = "boolean"
	// CategoryAsset classifies static visual assets, icons, and logos.
	CategoryAsset FileCategoryType = "asset"
	// CategorySpec classifies documentation and architectural specifications.
	CategorySpec FileCategoryType = "spec"
	// CategoryOther classifies auxiliary configuration and general files.
	CategoryOther FileCategoryType = "other"
)

// IsValid reports whether the value matches a known variant.
func (f FileCategoryType) IsValid() bool {
	switch f {
	case CategoryStyle, CategoryCore, CategoryBoolean, CategoryAsset, CategorySpec, CategoryOther:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (f FileCategoryType) String() string {
	return string(f)
}

// ParseFileCategoryType parses a raw string into a validated FileCategoryType.
func ParseFileCategoryType(s string) (FileCategoryType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllFileCategoryTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "categorystyle":
		return CategoryStyle, nil
	case "style":
		return CategoryStyle, nil
	case "categorycore":
		return CategoryCore, nil
	case "core":
		return CategoryCore, nil
	case "categoryboolean":
		return CategoryBoolean, nil
	case "boolean":
		return CategoryBoolean, nil
	case "categoryasset":
		return CategoryAsset, nil
	case "asset":
		return CategoryAsset, nil
	case "categoryspec":
		return CategorySpec, nil
	case "spec":
		return CategorySpec, nil
	case "categoryother":
		return CategoryOther, nil
	case "other":
		return CategoryOther, nil
	}
	return "", apperror.NewValidation("ParseFileCategoryType", "E1000", fmt.Sprintf("invalid filecategorytype: %q", s)).WithContext("value", s)
}

// AllFileCategoryTypes returns all valid variants.
func AllFileCategoryTypes() []FileCategoryType {
	return []FileCategoryType{
		CategoryStyle,
		CategoryCore,
		CategoryBoolean,
		CategoryAsset,
		CategorySpec,
		CategoryOther,
	}
}

// IsStyle reports whether the value is CategoryStyle.
func (f FileCategoryType) IsStyle() bool {
	return f == CategoryStyle
}

// IsCore reports whether the value is CategoryCore.
func (f FileCategoryType) IsCore() bool {
	return f == CategoryCore
}

// IsBoolean reports whether the value is CategoryBoolean.
func (f FileCategoryType) IsBoolean() bool {
	return f == CategoryBoolean
}

// IsAsset reports whether the value is CategoryAsset.
func (f FileCategoryType) IsAsset() bool {
	return f == CategoryAsset
}

// IsSpec reports whether the value is CategorySpec.
func (f FileCategoryType) IsSpec() bool {
	return f == CategorySpec
}

// IsOther reports whether the value is CategoryOther.
func (f FileCategoryType) IsOther() bool {
	return f == CategoryOther
}

// HasCategory reports whether the value is non-empty and valid.
func (f FileCategoryType) HasCategory() bool {
	return f.IsValid()
}
