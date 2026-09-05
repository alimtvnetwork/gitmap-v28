package openfiletype

import (
	"os"
	"strings"

	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
)

var variantLabels = [...]string{
	Invalid:               "Invalid",
	ReadOnly:              "ReadOnly",
	WriteOnly:             "WriteOnly",
	ReadWrite:             "ReadWrite",
	Append:                "Append",
	CreateAppend:          "CreateAppend",
	CreateTruncate:        "CreateTruncate",
	CreateNew:             "CreateNew",
	ReadOrCreateOnly:      "ReadOrCreateOnly",
	WriteOrCreateOnly:     "WriteOrCreateOnly",
	ReadWriteOrCreateOnly: "ReadWriteOrCreateOnly",
}

var openFlags = [...]int{
	Invalid:               os.O_RDONLY,
	ReadOnly:              os.O_RDONLY,
	WriteOnly:             os.O_WRONLY,
	ReadWrite:             os.O_RDWR,
	Append:                os.O_WRONLY | os.O_APPEND,
	CreateAppend:          os.O_CREATE | os.O_WRONLY | os.O_APPEND,
	CreateTruncate:        os.O_CREATE | os.O_WRONLY | os.O_TRUNC,
	CreateNew:             os.O_CREATE | os.O_EXCL | os.O_WRONLY,
	ReadOrCreateOnly:      os.O_RDONLY | os.O_CREATE,
	WriteOrCreateOnly:     os.O_WRONLY | os.O_CREATE,
	ReadWriteOrCreateOnly: os.O_RDWR | os.O_CREATE,
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
		return result.WrapFailureWithId[Variant](errtype.Validation, "cannot parse empty string as openfiletype")
	}

	for idx, label := range variantLabels {
		if strings.EqualFold(label, trimmed) {
			return result.WrapSuccess(Variant(idx))
		}
	}

	return result.WrapFailureWithId[Variant](errtype.NotFound, "unknown openfiletype variant: "+s)
}
