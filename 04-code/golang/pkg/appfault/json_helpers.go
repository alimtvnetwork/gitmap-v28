package appfault

import (
	"encoding/json"

	"coding-guidelines/common/pkg/errtype"
)

// SerializeToJSON serializes any value to JSON bytes returning *AppError on failure.
func SerializeToJSON(val any) ([]byte, *AppError) {
	data, err := json.Marshal(val)
	if err != nil {
		return nil, Wrap(errtype.Execution, err, "failed to marshal value to json")
	}

	return data, nil
}

// SerializeToJSONString serializes any value to a JSON string returning *AppError on failure.
func SerializeToJSONString(val any) (string, *AppError) {
	data, appErr := SerializeToJSON(val)
	if appErr != nil {
		return "", appErr
	}

	return string(data), nil
}

// DeserializeFromJSON parses JSON bytes into target value returning *AppError on failure.
func DeserializeFromJSON[T any](data []byte) (T, *AppError) {
	var target T
	if err := json.Unmarshal(data, &target); err != nil {
		return target, Wrap(errtype.Validation, err, "failed to unmarshal json data")
	}

	return target, nil
}

// DeserializeFromJSONString parses JSON string into target value returning *AppError on failure.
func DeserializeFromJSONString[T any](str string) (T, *AppError) {
	return DeserializeFromJSON[T]([]byte(str))
}
