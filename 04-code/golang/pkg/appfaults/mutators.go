package appfaults

import (
	"fmt"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// Add appends an AppError if it is non-nil and has an active error.
func (c *Collection) Add(err *appfault.AppError) *Collection {
	if c == nil || err == nil || !err.HasError() {
		return c
	}

	c.items = append(c.items, err)

	return c
}

// AddType creates an AppError from an errtype and appends it.
func (c *Collection) AddType(errType errtype.Variation) *Collection {
	return c.Add(appfault.NewType(errType))
}

// AddTypeMsg creates an AppError with message and appends it.
func (c *Collection) AddTypeMsg(errType errtype.Variation, msg string) *Collection {
	return c.Add(appfault.New(errType, msg))
}

// AddTypeMsgf creates an AppError with formatted message and appends it.
func (c *Collection) AddTypeMsgf(errType errtype.Variation, format string, args ...any) *Collection {
	return c.Add(appfault.New(errType, fmt.Sprintf(format, args...)))
}

// AddError wraps a cause error with an explicit errtype and appends it.
func (c *Collection) AddError(errType errtype.Variation, cause error) *Collection {
	return c.Add(appfault.WrapType(errType, cause))
}

// AddErrorMsg wraps a cause error with an explicit errtype, custom message and appends it.
func (c *Collection) AddErrorMsg(errType errtype.Variation, cause error, msg string) *Collection {
	return c.Add(appfault.Wrap(errType, cause, msg))
}

// AddWithContext creates an AppError with context map and appends it.
func (c *Collection) AddWithContext(errType errtype.Variation, msg string, ctx map[string]any) *Collection {
	return c.Add(appfault.NewWithContext(errType, msg, ctx))
}

// AddAll appends multiple AppErrors in order.
func (c *Collection) AddAll(faults ...*appfault.AppError) *Collection {
	for _, f := range faults {
		c.Add(f)
	}

	return c
}

// Merge copies all items from another collection into c.
func (c *Collection) Merge(other *Collection) *Collection {
	if c == nil || other == nil {
		return c
	}

	return c.AddAll(other.items...)
}

// Clear removes all items from the collection.
func (c *Collection) Clear() *Collection {
	if c != nil {
		c.items = c.items[:0]
	}

	return c
}
