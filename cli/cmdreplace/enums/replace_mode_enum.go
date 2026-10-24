package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// ReplaceModeType enumerates the invocation shapes the spec accepts.
type ReplaceModeType int

const (
	ReplaceModeTypeUnknown  ReplaceModeType = 0
	ReplaceModeTypeLiteral  ReplaceModeType = 1
	ReplaceModeTypeVersionN ReplaceModeType = 2
	ReplaceModeTypeAll      ReplaceModeType = 3
	ReplaceModeTypeAudit    ReplaceModeType = 4
)

// IsValid reports whether the value matches a known variant.
func (r ReplaceModeType) IsValid() bool {
	switch r {
	case ReplaceModeTypeUnknown, ReplaceModeTypeLiteral, ReplaceModeTypeVersionN, ReplaceModeTypeAll, ReplaceModeTypeAudit:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (r ReplaceModeType) String() string {
	return strconv.Itoa(int(r))
}

// ParseReplaceModeType parses a raw string into a validated ReplaceModeType.
func ParseReplaceModeType(s string) (ReplaceModeType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := ReplaceModeType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllReplaceModeTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "replacemodetypeunknown":
		return ReplaceModeTypeUnknown, nil
	case "unknown":
		return ReplaceModeTypeUnknown, nil
	case "replacemodetypeliteral":
		return ReplaceModeTypeLiteral, nil
	case "literal":
		return ReplaceModeTypeLiteral, nil
	case "replacemodetypeversionn":
		return ReplaceModeTypeVersionN, nil
	case "versionn":
		return ReplaceModeTypeVersionN, nil
	case "replacemodetypeall":
		return ReplaceModeTypeAll, nil
	case "all":
		return ReplaceModeTypeAll, nil
	case "replacemodetypeaudit":
		return ReplaceModeTypeAudit, nil
	case "audit":
		return ReplaceModeTypeAudit, nil
	}
	return 0, apperror.NewValidation("ParseReplaceModeType", "E1000", fmt.Sprintf("invalid replacemodetype: %q", s)).WithContext("value", s)
}

// AllReplaceModeTypes returns all valid variants.
func AllReplaceModeTypes() []ReplaceModeType {
	return []ReplaceModeType{
		ReplaceModeTypeUnknown,
		ReplaceModeTypeLiteral,
		ReplaceModeTypeVersionN,
		ReplaceModeTypeAll,
		ReplaceModeTypeAudit,
	}
}

// IsUnknown reports whether the value is ReplaceModeTypeUnknown.
func (r ReplaceModeType) IsUnknown() bool {
	return r == ReplaceModeTypeUnknown
}

// IsLiteral reports whether the value is ReplaceModeTypeLiteral.
func (r ReplaceModeType) IsLiteral() bool {
	return r == ReplaceModeTypeLiteral
}

// IsVersionN reports whether the value is ReplaceModeTypeVersionN.
func (r ReplaceModeType) IsVersionN() bool {
	return r == ReplaceModeTypeVersionN
}

// IsAll reports whether the value is ReplaceModeTypeAll.
func (r ReplaceModeType) IsAll() bool {
	return r == ReplaceModeTypeAll
}

// IsAudit reports whether the value is ReplaceModeTypeAudit.
func (r ReplaceModeType) IsAudit() bool {
	return r == ReplaceModeTypeAudit
}

// HasMode reports whether the value is non-empty and valid.
func (r ReplaceModeType) HasMode() bool {
	return r.IsValid()
}
