package appfaults

import "coding-guidelines/common/pkg/appfault"

// HasError returns true if collection contains at least one active error.
func (c *Collection) HasError() bool {
	return c != nil && len(c.items) > 0
}

// HasNoError returns true if collection contains zero errors or is nil.
func (c *Collection) HasNoError() bool {
	return c == nil || len(c.items) == 0
}

// HasNullError returns true if collection is nil or empty.
func (c *Collection) HasNullError() bool {
	return c.HasNoError()
}

// IsNull returns true if collection pointer is nil.
func (c *Collection) IsNull() bool {
	return c == nil
}

// IsSuccess returns true if collection is nil or contains zero errors.
func (c *Collection) IsSuccess() bool {
	return c.HasNoError()
}

// IsEmpty returns true if collection contains no errors.
func (c *Collection) IsEmpty() bool {
	return c.HasNoError()
}

// HasZero returns true if collection is nil or contains zero errors.
func (c *Collection) HasZero() bool {
	return c.HasNoError()
}

// IsZero returns true if collection is nil or contains zero errors.
func (c *Collection) IsZero() bool {
	return c.HasNoError()
}

// HasNull returns true if collection is nil or contains zero errors.
func (c *Collection) HasNull() bool {
	return c.HasNoError()
}

// IsValid returns true if collection is in a healthy valid state (no errors).
func (c *Collection) IsValid() bool {
	return c.HasNoError()
}

// IsInvalid returns true if collection has active errors.
func (c *Collection) IsInvalid() bool {
	return c.HasError()
}

// IsFailed returns true if collection has active errors.
func (c *Collection) IsFailed() bool {
	return c.HasError()
}

// Count returns the number of active errors in the collection.
func (c *Collection) Count() int {
	if c == nil {
		return 0
	}

	return len(c.items)
}

// Items returns a defensive copy of the underlying slice.
func (c *Collection) Items() []*appfault.AppError {
	if c.IsEmpty() {
		return []*appfault.AppError{}
	}

	copied := make([]*appfault.AppError, len(c.items))
	copy(copied, c.items)

	return copied
}

// First returns the first AppError or nil if empty.
func (c *Collection) First() *appfault.AppError {
	if c.IsEmpty() {
		return nil
	}

	return c.items[0]
}

// Last returns the last AppError or nil if empty.
func (c *Collection) Last() *appfault.AppError {
	if c.IsEmpty() {
		return nil
	}

	return c.items[len(c.items)-1]
}
