package enums

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// TargetSelectorType selects which cluster nodes a resolution targets.
type TargetSelectorType int

const (
	ServersClients TargetSelectorType = 1
	ClientsOnly    TargetSelectorType = 2
	ServersOnly    TargetSelectorType = 3
)

// IsValid reports whether the selector matches a known target selector.
func (t TargetSelectorType) IsValid() bool {
	switch t {
	case ServersClients, ClientsOnly, ServersOnly:
		return true
	default:
		return false
	}
}

// String returns the name of the target selector.
func (t TargetSelectorType) String() string {
	switch t {
	case ServersClients:
		return "ServersClients"
	case ClientsOnly:
		return "ClientsOnly"
	case ServersOnly:
		return "ServersOnly"
	default:
		return ""
	}
}

// ParseTargetSelectorType parses a raw string (name or number) into a validated TargetSelectorType.
func ParseTargetSelectorType(s string) (TargetSelectorType, *apperror.AppError) {
	cleaned := strings.TrimSpace(s)
	if num, err := strconv.Atoi(cleaned); err == nil {
		val := TargetSelectorType(num)
		if val.IsValid() {
			return val, nil
		}
	}

	for _, val := range AllTargetSelectorTypes() {
		if strings.EqualFold(val.String(), cleaned) {
			return val, nil
		}
	}

	return 0, apperror.NewValidation("ParseTargetSelectorType", "E1000", fmt.Sprintf("invalid cluster target selector: %q", s)).WithContext("value", s)
}

// AllTargetSelectorTypes returns all valid TargetSelectorType variants.
func AllTargetSelectorTypes() []TargetSelectorType {
	return []TargetSelectorType{ServersClients, ClientsOnly, ServersOnly}
}

// IsServersClients reports whether the selector targets servers and clients.
func (t TargetSelectorType) IsServersClients() bool {
	return t == ServersClients
}

// IsClientsOnly reports whether the selector targets clients only.
func (t TargetSelectorType) IsClientsOnly() bool {
	return t == ClientsOnly
}

// IsServersOnly reports whether the selector targets servers only.
func (t TargetSelectorType) IsServersOnly() bool {
	return t == ServersOnly
}
