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

// ToJSON exports the AppError as formatted JSON bytes.
func (e *AppError) ToJSON() ([]byte, error) {
	return json.MarshalIndent(e, "", "  ")
}

// ToJSONString exports the AppError as a JSON string.
func (e *AppError) ToJSONString() string {
	b, err := e.ToJSON()
	if err != nil {
		return "{}"
	}

	return string(b)
}

// FromJSON parses an AppError from JSON byte slice.
func FromJSON(data []byte) (*AppError, error) {
	var appErr AppError
	if err := json.Unmarshal(data, &appErr); err != nil {
		return nil, err
	}

	return &appErr, nil
}
