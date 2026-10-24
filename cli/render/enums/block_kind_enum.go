package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// BlockKindType identifies the type of a parsed block.
type BlockKindType int

const (
	BlockKindTypeParagraph BlockKindType = 0
	BlockKindTypeHeading   BlockKindType = 1
	BlockKindTypeSubtitle  BlockKindType = 2
	BlockKindTypeFence     BlockKindType = 3
	BlockKindTypeList      BlockKindType = 4
	BlockKindTypeBlank     BlockKindType = 5
)

// IsValid reports whether the value matches a known variant.
func (b BlockKindType) IsValid() bool {
	switch b {
	case BlockKindTypeParagraph, BlockKindTypeHeading, BlockKindTypeSubtitle, BlockKindTypeFence, BlockKindTypeList, BlockKindTypeBlank:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (b BlockKindType) String() string {
	return strconv.Itoa(int(b))
}

// ParseBlockKindType parses a raw string into a validated BlockKindType.
func ParseBlockKindType(s string) (BlockKindType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := BlockKindType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllBlockKindTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "blockkindtypeparagraph":
		return BlockKindTypeParagraph, nil
	case "paragraph":
		return BlockKindTypeParagraph, nil
	case "blockkindtypeheading":
		return BlockKindTypeHeading, nil
	case "heading":
		return BlockKindTypeHeading, nil
	case "blockkindtypesubtitle":
		return BlockKindTypeSubtitle, nil
	case "subtitle":
		return BlockKindTypeSubtitle, nil
	case "blockkindtypefence":
		return BlockKindTypeFence, nil
	case "fence":
		return BlockKindTypeFence, nil
	case "blockkindtypelist":
		return BlockKindTypeList, nil
	case "list":
		return BlockKindTypeList, nil
	case "blockkindtypeblank":
		return BlockKindTypeBlank, nil
	case "blank":
		return BlockKindTypeBlank, nil
	}
	return 0, apperror.NewValidation("ParseBlockKindType", "E1000", fmt.Sprintf("invalid blockkindtype: %q", s)).WithContext("value", s)
}

// AllBlockKindTypes returns all valid variants.
func AllBlockKindTypes() []BlockKindType {
	return []BlockKindType{
		BlockKindTypeParagraph,
		BlockKindTypeHeading,
		BlockKindTypeSubtitle,
		BlockKindTypeFence,
		BlockKindTypeList,
		BlockKindTypeBlank,
	}
}

// IsParagraph reports whether the value is BlockKindTypeParagraph.
func (b BlockKindType) IsParagraph() bool {
	return b == BlockKindTypeParagraph
}

// IsHeading reports whether the value is BlockKindTypeHeading.
func (b BlockKindType) IsHeading() bool {
	return b == BlockKindTypeHeading
}

// IsSubtitle reports whether the value is BlockKindTypeSubtitle.
func (b BlockKindType) IsSubtitle() bool {
	return b == BlockKindTypeSubtitle
}

// IsFence reports whether the value is BlockKindTypeFence.
func (b BlockKindType) IsFence() bool {
	return b == BlockKindTypeFence
}

// IsList reports whether the value is BlockKindTypeList.
func (b BlockKindType) IsList() bool {
	return b == BlockKindTypeList
}

// IsBlank reports whether the value is BlockKindTypeBlank.
func (b BlockKindType) IsBlank() bool {
	return b == BlockKindTypeBlank
}

// HasKind reports whether the value is non-empty and valid.
func (b BlockKindType) HasKind() bool {
	return b.IsValid()
}
