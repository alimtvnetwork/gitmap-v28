package openfiletype

import (
	"encoding/json"
	"fmt"
	"os"
)

type Variant byte

const (
	Invalid Variant = iota
	ReadOnly
	WriteOnly
	ReadWrite
	Append
	CreateAppend
	CreateTruncate
	CreateNew
	ReadOrCreateOnly
	WriteOrCreateOnly
	ReadWriteOrCreateOnly
)

type VariantPredicate func(v Variant) bool

func (v Variant) Flags() int {
	if int(v) < len(openFlags) {
		return openFlags[v]
	}

	return os.O_RDONLY
}

func (v Variant) Name() string {
	if int(v) < len(variantLabels) {
		return variantLabels[v]
	}

	return fmt.Sprintf("OpenFile(%d)", byte(v))
}

func (v Variant) Label() string {
	return v.Name()
}

func (v Variant) String() string {
	return v.Name()
}

func (v Variant) IsValid() bool {
	return v > Invalid && int(v) < len(variantLabels)
}

func (v Variant) IsInvalid() bool {
	return v <= Invalid || int(v) >= len(variantLabels)
}

func (v Variant) IsReadOnly() bool {
	return v == ReadOnly
}

func (v Variant) IsWriteOnly() bool {
	return v == WriteOnly
}

func (v Variant) IsReadWrite() bool {
	return v == ReadWrite
}

func (v Variant) IsAppend() bool {
	return v == Append
}

func (v Variant) IsCreateAppend() bool {
	return v == CreateAppend
}

func (v Variant) IsCreateTruncate() bool {
	return v == CreateTruncate
}

func (v Variant) IsCreateNew() bool {
	return v == CreateNew
}

func (v Variant) IsReadOrCreateOnly() bool {
	return v == ReadOrCreateOnly
}

func (v Variant) IsWriteOrCreateOnly() bool {
	return v == WriteOrCreateOnly
}

func (v Variant) IsReadWriteOrCreateOnly() bool {
	return v == ReadWriteOrCreateOnly
}

func (v Variant) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Name())
}

func (v *Variant) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		res := Parse(str)
		if res.IsSuccess() {
			*v = res.Data()

			return nil
		}
	}

	var raw byte
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*v = Variant(raw)

	return nil
}
