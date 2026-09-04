package appfaults

import (
	"errors"
	"fmt"
	"strings"

	"coding-guidelines/common/pkg/appfault"
)

// Compile builds a multi-line formatted summary using standard constants.
func (c *Collection) Compile() string {
	if c.IsEmpty() {
		return "No faults recorded."
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%sFault Collection (%d items):%s", appfault.HeaderPrefix, c.Count(), appfault.Newline))
	for idx, item := range c.items {
		b.WriteString(fmt.Sprintf("%s[%d] [%s:%d] %s%s", appfault.IndentTab, idx+1, item.Type().Name(), item.Type().Code(), item.Message(), appfault.Newline))
	}

	return b.String()
}

// appendFaultTraces appends all stack traces from faults.
func (c *Collection) appendFaultTraces(b *strings.Builder) {
	for idx, item := range c.items {
		b.WriteString(fmt.Sprintf("%s--- Fault #%d ---%s", appfault.Newline, idx+1, appfault.Newline))
		b.WriteString(item.CompileWithStack())
	}
}

// CompileWithStack builds a multi-line report including all stack traces.
func (c *Collection) CompileWithStack() string {
	if c.IsEmpty() {
		return "No faults recorded."
	}

	var b strings.Builder
	b.WriteString(c.Compile())
	b.WriteString(fmt.Sprintf("%s%sDetailed Fault Diagnostics:%s", appfault.Newline, appfault.SectionPrefix, appfault.Newline))
	c.appendFaultTraces(&b)

	return b.String()
}

// CompiledError returns a combined Go error.
func (c *Collection) CompiledError() error {
	if c.IsEmpty() {
		return nil
	}

	return errors.New(c.Compile())
}
