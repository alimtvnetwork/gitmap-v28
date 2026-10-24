package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// FormatType is a string tag persisted in ArchiveHistory.ArchiveFormat. It
// reads cleanly in PascalCase logs ("Zip", "TarGz") yet round-trips
// through the canonical extension via FormatFromExt / FormatType.Extension.
type FormatType string

const (
	FormatZip     FormatType = "Zip"
	FormatTar     FormatType = "Tar"
	FormatTarGz   FormatType = "TarGz"
	FormatTarBz2  FormatType = "TarBz2"
	FormatTarXz   FormatType = "TarXz"
	FormatTarZst  FormatType = "TarZst"
	FormatGz      FormatType = "Gz"
	FormatBz2     FormatType = "Bz2"
	FormatXz      FormatType = "Xz"
	FormatZst     FormatType = "Zst"
	Format7z      FormatType = "SevenZip"
	FormatRar     FormatType = "Rar"
	FormatUnknown FormatType = ""
)

// IsValid reports whether the value matches a known variant.
func (f FormatType) IsValid() bool {
	switch f {
	case FormatZip, FormatTar, FormatTarGz, FormatTarBz2, FormatTarXz, FormatTarZst, FormatGz, FormatBz2, FormatXz, FormatZst, Format7z, FormatRar, FormatUnknown:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (f FormatType) String() string {
	return string(f)
}

// ParseFormatType parses a raw string into a validated FormatType.
func ParseFormatType(s string) (FormatType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllFormatTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "formatzip":
		return FormatZip, nil
	case "zip":
		return FormatZip, nil
	case "formattar":
		return FormatTar, nil
	case "tar":
		return FormatTar, nil
	case "formattargz":
		return FormatTarGz, nil
	case "targz":
		return FormatTarGz, nil
	case "formattarbz2":
		return FormatTarBz2, nil
	case "tarbz2":
		return FormatTarBz2, nil
	case "formattarxz":
		return FormatTarXz, nil
	case "tarxz":
		return FormatTarXz, nil
	case "formattarzst":
		return FormatTarZst, nil
	case "tarzst":
		return FormatTarZst, nil
	case "formatgz":
		return FormatGz, nil
	case "gz":
		return FormatGz, nil
	case "formatbz2":
		return FormatBz2, nil
	case "bz2":
		return FormatBz2, nil
	case "formatxz":
		return FormatXz, nil
	case "xz":
		return FormatXz, nil
	case "formatzst":
		return FormatZst, nil
	case "zst":
		return FormatZst, nil
	case "format7z":
		return Format7z, nil
	case "7z":
		return Format7z, nil
	case "formatrar":
		return FormatRar, nil
	case "rar":
		return FormatRar, nil
	case "formatunknown":
		return FormatUnknown, nil
	case "unknown":
		return FormatUnknown, nil
	}
	return "", apperror.NewValidation("ParseFormatType", "E1000", fmt.Sprintf("invalid formattype: %q", s)).WithContext("value", s)
}

// AllFormatTypes returns all valid variants.
func AllFormatTypes() []FormatType {
	return []FormatType{
		FormatZip,
		FormatTar,
		FormatTarGz,
		FormatTarBz2,
		FormatTarXz,
		FormatTarZst,
		FormatGz,
		FormatBz2,
		FormatXz,
		FormatZst,
		Format7z,
		FormatRar,
		FormatUnknown,
	}
}

// IsZip reports whether the value is FormatZip.
func (f FormatType) IsZip() bool {
	return f == FormatZip
}

// IsTar reports whether the value is FormatTar.
func (f FormatType) IsTar() bool {
	return f == FormatTar
}

// IsTarGz reports whether the value is FormatTarGz.
func (f FormatType) IsTarGz() bool {
	return f == FormatTarGz
}

// IsTarBz2 reports whether the value is FormatTarBz2.
func (f FormatType) IsTarBz2() bool {
	return f == FormatTarBz2
}

// IsTarXz reports whether the value is FormatTarXz.
func (f FormatType) IsTarXz() bool {
	return f == FormatTarXz
}

// IsTarZst reports whether the value is FormatTarZst.
func (f FormatType) IsTarZst() bool {
	return f == FormatTarZst
}

// IsGz reports whether the value is FormatGz.
func (f FormatType) IsGz() bool {
	return f == FormatGz
}

// IsBz2 reports whether the value is FormatBz2.
func (f FormatType) IsBz2() bool {
	return f == FormatBz2
}

// IsXz reports whether the value is FormatXz.
func (f FormatType) IsXz() bool {
	return f == FormatXz
}

// IsZst reports whether the value is FormatZst.
func (f FormatType) IsZst() bool {
	return f == FormatZst
}

// IsFormat7z reports whether the value is Format7z.
func (f FormatType) IsFormat7z() bool {
	return f == Format7z
}

// IsRar reports whether the value is FormatRar.
func (f FormatType) IsRar() bool {
	return f == FormatRar
}

// IsUnknown reports whether the value is FormatUnknown.
func (f FormatType) IsUnknown() bool {
	return f == FormatUnknown
}

// HasFormat reports whether the value is non-empty and valid.
func (f FormatType) HasFormat() bool {
	return f.IsValid()
}

// formatExtMap maps each format to its canonical extension (with
// leading dot). Moved with FormatType.Extension from archive.go: it is
// only consumed by Extension, so it lives with the enum.
var formatExtMap = map[FormatType]string{
	FormatZip:    ".zip",
	FormatTar:    ".tar",
	FormatTarGz:  ".tar.gz",
	FormatTarBz2: ".tar.bz2",
	FormatTarXz:  ".tar.xz",
	FormatTarZst: ".tar.zst",
	FormatGz:     ".gz",
	FormatBz2:    ".bz2",
	FormatXz:     ".xz",
	FormatZst:    ".zst",
	Format7z:     ".7z",
	FormatRar:    ".rar",
}

// Extension returns the canonical extension (with leading dot) the
// CreateArchive path uses to construct mholt/archives Format objects.
func (f FormatType) Extension() string {
	return formatExtMap[f]
}
