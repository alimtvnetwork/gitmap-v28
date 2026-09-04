package appfaults

import (
	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// Concat creates a new Collection containing current items plus err.
func (c *Collection) Concat(err *appfault.AppError) *Collection {
	cloned := NewWithCapacity(c.Count() + 1)
	if c != nil {
		cloned.AddAll(c.items...)
	}

	cloned.Add(err)

	return cloned
}

// ConcatMultiple creates a new Collection containing current items plus faults.
func (c *Collection) ConcatMultiple(faults ...*appfault.AppError) *Collection {
	cloned := NewWithCapacity(c.Count() + len(faults))
	if c != nil {
		cloned.AddAll(c.items...)
	}

	cloned.AddAll(faults...)

	return cloned
}

// countTotalFaults computes total error count across collections.
func (c *Collection) countTotalFaults(others []*Collection) int {
	total := c.Count()
	for _, o := range others {
		total += o.Count()
	}

	return total
}

// ConcatNew creates a new Collection combining c and others without mutating any.
func (c *Collection) ConcatNew(others ...*Collection) *Collection {
	result := NewWithCapacity(c.countTotalFaults(others))
	if c != nil {
		result.AddAll(c.items...)
	}

	for _, o := range others {
		if o != nil {
			result.AddAll(o.items...)
		}
	}

	return result
}

// ConcatErrors wraps multiple causes and returns a new combined Collection.
func (c *Collection) ConcatErrors(errType errtype.Variation, causes ...error) *Collection {
	cloned := NewWithCapacity(c.Count() + len(causes))
	if c != nil {
		cloned.AddAll(c.items...)
	}

	for _, cause := range causes {
		cloned.AddError(errType, cause)
	}

	return cloned
}
