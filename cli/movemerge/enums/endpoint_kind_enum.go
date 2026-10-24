package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// EndpointKindType classifies a positional argument once at command start.
type EndpointKindType int

const (
	// EndpointFolder is a plain on-disk path (relative or absolute).
	EndpointFolder EndpointKindType = 0
	// EndpointURL is an https/http/ssh/git@ remote, optionally :branch.
	EndpointURL EndpointKindType = 1
)

// IsValid reports whether the value matches a known variant.
func (e EndpointKindType) IsValid() bool {
	switch e {
	case EndpointFolder, EndpointURL:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (e EndpointKindType) String() string {
	return strconv.Itoa(int(e))
}

// ParseEndpointKindType parses a raw string into a validated EndpointKindType.
func ParseEndpointKindType(s string) (EndpointKindType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := EndpointKindType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllEndpointKindTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "endpointfolder":
		return EndpointFolder, nil
	case "folder":
		return EndpointFolder, nil
	case "endpointurl":
		return EndpointURL, nil
	case "url":
		return EndpointURL, nil
	}
	return 0, apperror.NewValidation("ParseEndpointKindType", "E1000", fmt.Sprintf("invalid endpointkindtype: %q", s)).WithContext("value", s)
}

// AllEndpointKindTypes returns all valid variants.
func AllEndpointKindTypes() []EndpointKindType {
	return []EndpointKindType{
		EndpointFolder,
		EndpointURL,
	}
}

// IsFolder reports whether the value is EndpointFolder.
func (e EndpointKindType) IsFolder() bool {
	return e == EndpointFolder
}

// IsURL reports whether the value is EndpointURL.
func (e EndpointKindType) IsURL() bool {
	return e == EndpointURL
}

// HasKind reports whether the value is non-empty and valid.
func (e EndpointKindType) HasKind() bool {
	return e.IsValid()
}
