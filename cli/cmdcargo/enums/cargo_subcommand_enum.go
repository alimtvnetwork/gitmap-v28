package enums

import (
	"fmt"
	"strings"
)

// CargoSubcommandType represents recognized GitMap cargo subcommands.
type CargoSubcommandType string

const (
	CargoSubcommandStatus  CargoSubcommandType = "status"
	CargoSubcommandSt      CargoSubcommandType = "st"
	CargoSubcommandInfo    CargoSubcommandType = "info"
	CargoSubcommandInstall CargoSubcommandType = "install"
	CargoSubcommandIn      CargoSubcommandType = "in"
)

// IsValid reports whether the subcommand matches a known variant.
func (c CargoSubcommandType) IsValid() bool {
	switch c {
	case CargoSubcommandStatus, CargoSubcommandSt, CargoSubcommandInfo,
		CargoSubcommandInstall, CargoSubcommandIn:
		return true
	default:
		return false
	}
}

// String returns the string representation of the subcommand.
func (c CargoSubcommandType) String() string {
	return string(c)
}

// ParseCargoSubcommandType parses a raw string into a validated CargoSubcommandType.
func ParseCargoSubcommandType(s string) (CargoSubcommandType, error) {
	val := CargoSubcommandType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid cargo subcommand: %q", s)
}

// AllCargoSubcommandTypes returns all valid CargoSubcommandType variants.
func AllCargoSubcommandTypes() []CargoSubcommandType {
	return []CargoSubcommandType{
		CargoSubcommandStatus,
		CargoSubcommandSt,
		CargoSubcommandInfo,
		CargoSubcommandInstall,
		CargoSubcommandIn,
	}
}

// IsStatus reports whether the subcommand is a status inquiry variant.
func (c CargoSubcommandType) IsStatus() bool {
	return c == CargoSubcommandStatus || c == CargoSubcommandSt || c == CargoSubcommandInfo
}

// IsInstall reports whether the subcommand is an installation variant.
func (c CargoSubcommandType) IsInstall() bool {
	return c == CargoSubcommandInstall || c == CargoSubcommandIn
}

// HasSubcommand reports whether the subcommand is non-empty and valid.
func (c CargoSubcommandType) HasSubcommand() bool {
	return c.IsValid()
}
