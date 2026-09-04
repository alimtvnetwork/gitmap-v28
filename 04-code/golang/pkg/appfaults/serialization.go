package appfaults

import (
	"encoding/json"
	"fmt"
	"strings"

	"coding-guidelines/common/pkg/appfault"
)

// String implements fmt.Stringer returning multi-line formatted summary.
func (c *Collection) String() string {
	return c.Format()
}

// Format formats the collection as numbered list of error messages.
func (c *Collection) Format() string {
	if !c.HasError() {
		return "No faults recorded."
	}

	lines := make([]string, 0, len(c.items))
	for idx, item := range c.items {
		lines = append(lines, fmt.Sprintf("[%d] %s", idx+1, item.Error()))
	}

	return strings.Join(lines, "\n")
}

// MarshalJSON serializes the collection as a JSON array of AppErrors.
func (c *Collection) MarshalJSON() ([]byte, error) {
	if !c.HasError() {
		return []byte("[]"), nil
	}

	return json.Marshal(c.items)
}

// UnmarshalJSON unmarshals a JSON array of AppErrors into the collection.
func (c *Collection) UnmarshalJSON(data []byte) error {
	var items []*appfault.AppError
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	c.items = make([]*appfault.AppError, 0, len(items))
	c.AddAll(items...)

	return nil
}
