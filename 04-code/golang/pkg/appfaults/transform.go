package appfaults

import (
	"fmt"
	"strings"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// Filter returns a new Collection containing only matching items.
func (c *Collection) Filter(predicate func(*appfault.AppError) bool) *Collection {
	filtered := New()
	if c == nil || predicate == nil {
		return filtered
	}

	for _, item := range c.items {
		if predicate(item) {
			filtered.Add(item)
		}
	}

	return filtered
}

// FilterByType returns a new Collection with errors matching target variation.
func (c *Collection) FilterByType(target errtype.Variation) *Collection {
	return c.Filter(func(e *appfault.AppError) bool {
		return e.Is(target)
	})
}

// ToAppError merges collection errors into a single composite AppError.
func (c *Collection) ToAppError(errType errtype.Variation, title string) *appfault.AppError {
	if c.IsEmpty() {
		return nil
	}

	compositeMsg := fmt.Sprintf("%s (%d faults):\n%s", title, c.Count(), c.Format())
	first := c.First()
	if first != nil && first.Cause() != nil {
		return appfault.Wrap(errType, first.Cause(), compositeMsg)
	}

	return appfault.New(errType, compositeMsg)
}

// Errors converts all items into a standard []error slice.
func (c *Collection) Errors() []error {
	if c.IsEmpty() {
		return []error{}
	}

	errs := make([]error, 0, len(c.items))
	for _, item := range c.items {
		errs = append(errs, item)
	}

	return errs
}

// ErrorString compiles errors into a joined string.
func (c *Collection) ErrorString() string {
	if c.IsEmpty() {
		return ""
	}

	msgs := make([]string, 0, len(c.items))
	for idx, item := range c.items {
		msgs = append(msgs, fmt.Sprintf("%d. %s", idx+1, item.Error()))
	}

	return strings.Join(msgs, "\n")
}
