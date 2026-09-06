package logleveltype

import (
	"encoding/json"
	"fmt"
	"strings"
)

type (
	Variant byte

	VariantPredicate func(v Variant) bool
)

const (
	Invalid Variant = iota
	Debug
	Info
	Warn
	Error
	Fatal
)

const Unknown = Invalid

func (v Variant) Name() string {
	if int(v) < len(variantLabels) {
		return variantLabels[v]
	}

	return fmt.Sprintf("LogLevel(%d)", byte(v))
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

func (v Variant) IsEnabled(threshold Variant) bool {
	return v >= threshold
}

func (v Variant) IsDebug() bool {
	return v == Debug
}

func (v Variant) IsInfo() bool {
	return v == Info
}

func (v Variant) IsWarn() bool {
	return v == Warn
}

func (v Variant) IsError() bool {
	return v == Error
}

func (v Variant) IsFatal() bool {
	return v == Fatal
}

func (v Variant) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Name())
}

func (v *Variant) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 || trimmed == "null" {
		*v = Invalid

		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		res := Parse(str)
		if res.IsSuccess() {
			*v = res.Data()

			return nil
		}

		return res.Fault()
	}

	var raw byte
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if int(raw) >= len(variantLabels) {
		return fmt.Errorf("invalid logleveltype numeric value %d, supported range: 0..%d", raw, len(variantLabels)-1)
	}

	*v = Variant(raw)

	return nil
}
