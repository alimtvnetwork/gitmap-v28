package enums

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// VHostSiteType represents the supported site types for virtual hosts.
type VHostSiteType string

const (
	VHostSiteTypeWordpress VHostSiteType = "wordpress"
	VHostSiteTypeLaravel   VHostSiteType = "laravel"
	VHostSiteTypePHP       VHostSiteType = "php"
	VHostSiteTypeStatic    VHostSiteType = "static"
)

// IsValid reports whether the value matches a known variant.
func (v VHostSiteType) IsValid() bool {
	switch v {
	case VHostSiteTypeWordpress, VHostSiteTypeLaravel, VHostSiteTypePHP, VHostSiteTypeStatic:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (v VHostSiteType) String() string {
	return string(v)
}

func ParseVHostSiteType(raw string) (VHostSiteType, *apperror.AppError) {
	norm := strings.ToLower(strings.TrimSpace(raw))
	switch norm {
	case "wordpress", "wp":

		return VHostSiteTypeWordpress, nil
	case "laravel", "art", "artisan":

		return VHostSiteTypeLaravel, nil
	case "php":

		return VHostSiteTypePHP, nil
	case "static", "html":

		return VHostSiteTypeStatic, nil
	default:

		return "", apperror.NewValidationError("unrecognized site type: " + raw)
	}
}

// AllVHostSiteTypes returns all valid variants.
func AllVHostSiteTypes() []VHostSiteType {
	return []VHostSiteType{
		VHostSiteTypeWordpress,
		VHostSiteTypeLaravel,
		VHostSiteTypePHP,
		VHostSiteTypeStatic,
	}
}

// IsWordpress reports whether the value is VHostSiteTypeWordpress.
func (v VHostSiteType) IsWordpress() bool {
	return v == VHostSiteTypeWordpress
}

// IsLaravel reports whether the value is VHostSiteTypeLaravel.
func (v VHostSiteType) IsLaravel() bool {
	return v == VHostSiteTypeLaravel
}

// IsPHP reports whether the value is VHostSiteTypePHP.
func (v VHostSiteType) IsPHP() bool {
	return v == VHostSiteTypePHP
}

// IsStatic reports whether the value is VHostSiteTypeStatic.
func (v VHostSiteType) IsStatic() bool {
	return v == VHostSiteTypeStatic
}

// HasHostSite reports whether the value is non-empty and valid.
func (v VHostSiteType) HasHostSite() bool {
	return v.IsValid()
}
