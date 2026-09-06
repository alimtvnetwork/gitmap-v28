package openfiletype

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
)

var (
	variantLabels = [...]string{
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

	openFlags = [...]int{
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

	variantMap = compileVariantMap()
)

func compileVariantMap() map[string]Variant {
	m := make(map[string]Variant, (len(variantLabels)*4)+4)
	for i, label := range variantLabels {
		v := Variant(i)
		m[label] = v
		m[strings.ToLower(label)] = v
		m[strings.ToUpper(label)] = v
		m[strconv.Itoa(i)] = v
	}

	m["unknown"] = Invalid
	m["invalid"] = Invalid
	m["UNKNOWN"] = Invalid
	m["INVALID"] = Invalid

	return m
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

	if v, ok := variantMap[strings.ToLower(trimmed)]; ok {
		return result.WrapSuccess(v)
	}

	return result.WrapFailureWithId[Variant](
		errtype.NotFound,
		fmt.Sprintf("unknown openfiletype variant %q, supported variants: [%s]", s, strings.Join(Values(), ", ")),
	)
}
