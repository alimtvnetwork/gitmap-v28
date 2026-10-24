package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// ZshThemeType represents supported ZSH prompt theme names.
type ZshThemeType string

const (
	ZshThemeRobbyrussell ZshThemeType = "robbyrussell"
	ZshThemeAgnoster     ZshThemeType = "agnoster"
	ZshThemeFletcherm    ZshThemeType = "fletcherm"
	ZshThemeAfMagic      ZshThemeType = "af-magic"
	ZshThemeAfowler      ZshThemeType = "afowler"
	ZshThemeBira         ZshThemeType = "bira"
	ZshThemeCandy        ZshThemeType = "candy"
	ZshThemeClean        ZshThemeType = "clean"
	ZshThemeCloud        ZshThemeType = "cloud"
	ZshThemeGnzh         ZshThemeType = "gnzh"
	ZshThemeHalfLife     ZshThemeType = "half-life"
)

// IsValid reports whether the value matches a known variant.
func (z ZshThemeType) IsValid() bool {
	switch z {
	case ZshThemeRobbyrussell, ZshThemeAgnoster, ZshThemeFletcherm, ZshThemeAfMagic, ZshThemeAfowler, ZshThemeBira, ZshThemeCandy, ZshThemeClean, ZshThemeCloud, ZshThemeGnzh, ZshThemeHalfLife:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (z ZshThemeType) String() string {
	return string(z)
}

// ParseZshThemeType parses a raw string into a validated ZshThemeType.
func ParseZshThemeType(s string) (ZshThemeType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllZshThemeTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "zshthemerobbyrussell":
		return ZshThemeRobbyrussell, nil
	case "robbyrussell":
		return ZshThemeRobbyrussell, nil
	case "zshthemeagnoster":
		return ZshThemeAgnoster, nil
	case "agnoster":
		return ZshThemeAgnoster, nil
	case "zshthemefletcherm":
		return ZshThemeFletcherm, nil
	case "fletcherm":
		return ZshThemeFletcherm, nil
	case "zshthemeafmagic":
		return ZshThemeAfMagic, nil
	case "afmagic":
		return ZshThemeAfMagic, nil
	case "zshthemeafowler":
		return ZshThemeAfowler, nil
	case "afowler":
		return ZshThemeAfowler, nil
	case "zshthemebira":
		return ZshThemeBira, nil
	case "bira":
		return ZshThemeBira, nil
	case "zshthemecandy":
		return ZshThemeCandy, nil
	case "candy":
		return ZshThemeCandy, nil
	case "zshthemeclean":
		return ZshThemeClean, nil
	case "clean":
		return ZshThemeClean, nil
	case "zshthemecloud":
		return ZshThemeCloud, nil
	case "cloud":
		return ZshThemeCloud, nil
	case "zshthemegnzh":
		return ZshThemeGnzh, nil
	case "gnzh":
		return ZshThemeGnzh, nil
	case "zshthemehalflife":
		return ZshThemeHalfLife, nil
	case "halflife":
		return ZshThemeHalfLife, nil
	}
	return "", apperror.NewValidation("ParseZshThemeType", "E1000", fmt.Sprintf("invalid zshthemetype: %q", s)).WithContext("value", s)
}

// AllZshThemeTypes returns all valid variants.
func AllZshThemeTypes() []ZshThemeType {
	return []ZshThemeType{
		ZshThemeRobbyrussell,
		ZshThemeAgnoster,
		ZshThemeFletcherm,
		ZshThemeAfMagic,
		ZshThemeAfowler,
		ZshThemeBira,
		ZshThemeCandy,
		ZshThemeClean,
		ZshThemeCloud,
		ZshThemeGnzh,
		ZshThemeHalfLife,
	}
}

// IsRobbyrussell reports whether the value is ZshThemeRobbyrussell.
func (z ZshThemeType) IsRobbyrussell() bool {
	return z == ZshThemeRobbyrussell
}

// IsAgnoster reports whether the value is ZshThemeAgnoster.
func (z ZshThemeType) IsAgnoster() bool {
	return z == ZshThemeAgnoster
}

// IsFletcherm reports whether the value is ZshThemeFletcherm.
func (z ZshThemeType) IsFletcherm() bool {
	return z == ZshThemeFletcherm
}

// IsAfMagic reports whether the value is ZshThemeAfMagic.
func (z ZshThemeType) IsAfMagic() bool {
	return z == ZshThemeAfMagic
}

// IsAfowler reports whether the value is ZshThemeAfowler.
func (z ZshThemeType) IsAfowler() bool {
	return z == ZshThemeAfowler
}

// IsBira reports whether the value is ZshThemeBira.
func (z ZshThemeType) IsBira() bool {
	return z == ZshThemeBira
}

// IsCandy reports whether the value is ZshThemeCandy.
func (z ZshThemeType) IsCandy() bool {
	return z == ZshThemeCandy
}

// IsClean reports whether the value is ZshThemeClean.
func (z ZshThemeType) IsClean() bool {
	return z == ZshThemeClean
}

// IsCloud reports whether the value is ZshThemeCloud.
func (z ZshThemeType) IsCloud() bool {
	return z == ZshThemeCloud
}

// IsGnzh reports whether the value is ZshThemeGnzh.
func (z ZshThemeType) IsGnzh() bool {
	return z == ZshThemeGnzh
}

// IsHalfLife reports whether the value is ZshThemeHalfLife.
func (z ZshThemeType) IsHalfLife() bool {
	return z == ZshThemeHalfLife
}

// HasTheme reports whether the value is non-empty and valid.
func (z ZshThemeType) HasTheme() bool {
	return z.IsValid()
}
