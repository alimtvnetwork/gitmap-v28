package streamwriter

import (
	"encoding/json"
	"fmt"
	"reflect"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// PayloadKind identifies the classification of an incoming generic payload.
type PayloadKind byte

var payloadKindNames = [...]string{
	"Nil",
	"Bytes",
	"String",
	"Error",
	"Map",
	"Struct",
	"Primitive",
}

// String implements fmt.Stringer for PayloadKind.
func (k PayloadKind) String() string {
	idx := int(k)
	if idx < len(payloadKindNames) {
		return payloadKindNames[idx]
	}

	return fmt.Sprintf("PayloadKind(%d)", idx)
}

// InspectPayload determines the classification of any incoming payload.
func InspectPayload(payload any) PayloadKind {
	if payload == nil {
		return PayloadNil
	}

	switch payload.(type) {
	case []byte:
		return PayloadBytes
	case string:
		return PayloadString
	case *appfault.AppError:
		return PayloadError
	case error:
		return PayloadError
	case map[string]any:
		return PayloadMap
	}

	val := reflect.ValueOf(payload)
	kind := val.Kind()
	if kind == reflect.Ptr {
		if val.IsNil() {
			return PayloadNil
		}

		kind = val.Elem().Kind()
	}

	switch kind {
	case reflect.Struct:
		return PayloadStruct
	case reflect.Map:
		return PayloadMap
	case reflect.Slice, reflect.Array:
		return PayloadBytes
	default:
		return PayloadPrimitive
	}
}

// ExtractBytes converts an arbitrary payload into raw bytes without redundant serialization.
// CRITICAL: If the payload is already []byte, it is returned directly WITHOUT json.Marshal,
// preventing Go's standard json encoder from Base64-encoding raw byte slices.
func ExtractBytes(payload any) []byte {
	if payload == nil {
		return nil
	}

	switch casted := payload.(type) {
	case []byte:
		return casted
	case string:
		return []byte(casted)
	case *appfault.AppError:
		return []byte(casted.Error())
	case error:
		return []byte(casted.Error())
	default:
		serialized, err := json.Marshal(payload)
		if err != nil {
			return []byte(fmt.Sprintf("%v", payload))
		}

		return serialized
	}
}

// ExtractJSONBytes converts a payload into valid JSON bytes.
// If payload is already valid JSON []byte, it returns it directly to avoid double-encoding.
func ExtractJSONBytes(payload any) ([]byte, *appfault.AppError) {
	if payload == nil {
		return []byte("null"), nil
	}

	switch casted := payload.(type) {
	case []byte:
		if json.Valid(casted) {
			return casted, nil
		}

		marshaled, err := json.Marshal(string(casted))
		if err != nil {
			return nil, appfault.Wrap(errtype.Validation, err, "failed to serialize raw bytes to json")
		}

		return marshaled, nil

	case string:
		if json.Valid([]byte(casted)) {
			return []byte(casted), nil
		}

		marshaled, err := json.Marshal(casted)
		if err != nil {
			return nil, appfault.Wrap(errtype.Validation, err, "failed to serialize string to json")
		}

		return marshaled, nil

	case *appfault.AppError:
		jsonBytes, err := casted.ToJSON()
		if err != nil {
			return nil, appfault.Wrap(errtype.Execution, err, "failed to serialize AppError to json")
		}

		return jsonBytes, nil

	default:
		marshaled, err := json.Marshal(payload)
		if err != nil {
			return nil, appfault.Wrap(errtype.Execution, err, "failed to marshal payload to json")
		}

		return marshaled, nil
	}
}
