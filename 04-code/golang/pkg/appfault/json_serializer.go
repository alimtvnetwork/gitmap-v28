package appfault

import "encoding/json"

// MarshalJSON implements custom JSON marshaling with stringified Cause.
func (e *AppError) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}

	return json.Marshal(e.ToDataModel())
}

// UnmarshalJSON implements custom JSON unmarshaling into AppError.
func (e *AppError) UnmarshalJSON(data []byte) error {
	var model AppErrorDataModel
	if err := json.Unmarshal(data, &model); err != nil {
		return err
	}

	*e = *model.ToAppError()

	return nil
}

// ToJson exports the AppError as formatted JSON bytes.
func (e *AppError) ToJson() ([]byte, error) {
	return json.MarshalIndent(e, "", "  ")
}

// ToJSON is an alias for ToJson.
func (e *AppError) ToJSON() ([]byte, error) {
	return e.ToJson()
}

// ToJsonString exports the AppError as a JSON string.
func (e *AppError) ToJsonString() string {
	b, err := e.ToJson()
	if err != nil {
		return "{}"
	}

	return string(b)
}

// ToJSONString is an alias for ToJsonString.
func (e *AppError) ToJSONString() string {
	return e.ToJsonString()
}

// FromJson parses an AppError from JSON byte slice.
func FromJson(data []byte) (*AppError, error) {
	var appErr AppError
	if err := json.Unmarshal(data, &appErr); err != nil {
		return nil, err
	}

	return &appErr, nil
}

// FromJSON is an alias for FromJson.
func FromJSON(data []byte) (*AppError, error) {
	return FromJson(data)
}
