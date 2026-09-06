package appfaults

import (
	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

type (
	Collection struct {
		items []*appfault.AppError
	}

	AppFaults = Collection
)

// New creates an empty, non-nil error collection.
func New() *Collection {
	return &Collection{
		items: make([]*appfault.AppError, 0),
	}
}

// NewWithCapacity preallocates backing slice capacity.
func NewWithCapacity(capacity int) *Collection {
	return &Collection{
		items: make([]*appfault.AppError, 0, capacity),
	}
}

// NewFromFaults constructs a collection filtering out nil errors.
func NewFromFaults(faults ...*appfault.AppError) *Collection {
	c := NewWithCapacity(len(faults))
	c.AddAll(faults...)

	return c
}

// NewFromErrors wraps raw Go errors into AppError instances with an explicit error type.
func NewFromErrors(errType errtype.Variation, errs ...error) *Collection {
	c := NewWithCapacity(len(errs))
	for _, err := range errs {
		c.AddError(errType, err)
	}

	return c
}

// Clone creates a deep copy of the Collection. If receiver is nil, it safely returns an empty Collection.
func (c *Collection) Clone() *Collection {
	if c == nil {
		return New()
	}

	cloned := NewWithCapacity(len(c.items))
	for _, item := range c.items {
		if item != nil {
			cloned.Add(item.Clone())
		}
	}

	return cloned
}
