package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RustToolchainBinaryType represents standard Rust toolchain executable binaries.
type RustToolchainBinaryType string

const (
	RustBinaryCargo  RustToolchainBinaryType = "cargo"
	RustBinaryRustc  RustToolchainBinaryType = "rustc"
	RustBinaryRustup RustToolchainBinaryType = "rustup"
	RustBinaryClang  RustToolchainBinaryType = "clang"
)

// IsValid reports whether the binary matches a known toolchain executable.
func (r RustToolchainBinaryType) IsValid() bool {
	switch r {
	case RustBinaryCargo, RustBinaryRustc, RustBinaryRustup, RustBinaryClang:
		return true
	default:
		return false
	}
}

// String returns the string representation of the binary.
func (r RustToolchainBinaryType) String() string {
	return string(r)
}

// ParseRustToolchainBinaryType parses a raw string into a validated RustToolchainBinaryType.
func ParseRustToolchainBinaryType(s string) (RustToolchainBinaryType, *apperror.AppError) {
	val := RustToolchainBinaryType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseRustToolchainBinaryType", "E1000", fmt.Sprintf("invalid rust toolchain binary: %q", s)).WithContext("value", s)
}

// AllRustToolchainBinaryTypes returns all valid RustToolchainBinaryType variants.
func AllRustToolchainBinaryTypes() []RustToolchainBinaryType {
	return []RustToolchainBinaryType{
		RustBinaryCargo,
		RustBinaryRustc,
		RustBinaryRustup,
		RustBinaryClang,
	}
}

// IsCargo reports whether the binary is cargo.
func (r RustToolchainBinaryType) IsCargo() bool {
	return r == RustBinaryCargo
}

// IsRustc reports whether the binary is rustc.
func (r RustToolchainBinaryType) IsRustc() bool {
	return r == RustBinaryRustc
}

// IsRustup reports whether the binary is rustup.
func (r RustToolchainBinaryType) IsRustup() bool {
	return r == RustBinaryRustup
}

// IsClang reports whether the binary is clang.
func (r RustToolchainBinaryType) IsClang() bool {
	return r == RustBinaryClang
}

// HasBinary reports whether the binary is non-empty and valid.
func (r RustToolchainBinaryType) HasBinary() bool {
	return r.IsValid()
}
