package appfault

import (
	"errors"

	"coding-guidelines/common/pkg/errtype"
)

// AppErrorDataModel is the serializable DTO for AppError with PascalCase properties.
type AppErrorDataModel struct {
	Type       errtype.Variation `json:"Type,omitempty" yaml:"Type,omitempty"`
	Message    string            `json:"Message,omitempty" yaml:"Message,omitempty"`
	Caller     CallerInfo        `json:"Caller,omitempty" yaml:"Caller,omitempty"`
	Stack      StackTrace        `json:"Stack,omitempty" yaml:"Stack,omitempty"`
	Ctx        ContextMap        `json:"Ctx,omitempty" yaml:"Ctx,omitempty"`
	Cause      string            `json:"Cause,omitempty" yaml:"Cause,omitempty"`
	StatusCode int               `json:"StatusCode,omitempty" yaml:"StatusCode,omitempty"`
}

// extractCauseString safely extracts the cause error message string.
func extractCauseString(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}

// populateDataModel populates DTO fields from non-nil AppError.
func (e *AppError) populateDataModel() AppErrorDataModel {
	return AppErrorDataModel{
		Type: e.errType, Message: e.message, Caller: e.caller,
		Stack: e.stack, Ctx: e.ctx.Clone(), Cause: extractCauseString(e.cause),
		StatusCode: e.statusCode,
	}
}

// ToDataModel converts an AppError into its serializable data model.
func (e *AppError) ToDataModel() AppErrorDataModel {
	if e == nil {
		return AppErrorDataModel{}
	}

	return e.populateDataModel()
}

// buildCauseError safely constructs an error if cause is non-empty.
func buildCauseError(cause string) error {
	if len(cause) == 0 {
		return nil
	}

	return errors.New(cause)
}

// ToAppError reconstructs an AppError from the data model.
func (m AppErrorDataModel) ToAppError() *AppError {
	return &AppError{
		errType:    m.Type,
		message:    m.Message,
		caller:     m.Caller,
		stack:      m.Stack,
		ctx:        m.Ctx.Clone(),
		cause:      buildCauseError(m.Cause),
		statusCode: m.StatusCode,
	}
}
