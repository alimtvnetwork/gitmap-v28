package logleveltype

import (
	"strings"

	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
)

var variantLabels = [...]string{
	Invalid: "Unknown",
	Debug:   "Debug",
	Info:    "Info",
	Warn:    "Warn",
	Error:   "Error",
	Fatal:   "Fatal",
}

func All() []Variant {
	items := make([]Variant, 0, len(variantLabels)-1)
	for i := 1; i < len(variantLabels); i++ {
		items = append(items, Variant(i))
	}

	return items
}

func Values() []string {
	names := make([]string, 0, len(variantLabels)-1)
	for _, label := range variantLabels[1:] {
		names = append(names, label)
	}

	return names
}

func Parse(s string) result.Wrap[Variant] {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) == 0 {
		return result.WrapFailureWithId[Variant](errtype.Validation, "cannot parse empty string as logleveltype")
	}

	if strings.EqualFold(trimmed, "unknown") || strings.EqualFold(trimmed, "invalid") {
		return result.WrapSuccess(Invalid)
	}

	for idx, label := range variantLabels {
		if strings.EqualFold(label, trimmed) {
			return result.WrapSuccess(Variant(idx))
		}
	}

	return result.WrapFailureWithId[Variant](errtype.NotFound, "unknown logleveltype variant: "+s)
}
