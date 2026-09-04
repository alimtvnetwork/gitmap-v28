package appfault

import (
	"errors"
	"fmt"
	"strings"
)

// writeCompileDetails writes message, caller and context to builder.
func (e *AppError) writeCompileDetails(b *strings.Builder) {
	b.WriteString(fmt.Sprintf("%sMessage: %s%s", IndentTab, e.message, Newline))
	if !e.caller.IsEmpty() {
		b.WriteString(fmt.Sprintf("%sCaller:  %s%s", IndentTab, e.caller.String(), Newline))
	}

	if len(e.ctx) > 0 {
		b.WriteString(fmt.Sprintf("%sContext: %s%s", IndentTab, e.ctx.Format(), Newline))
	}
}

// Compile builds a formatted diagnostic summary using layout constants.
func (e *AppError) Compile() string {
	if e.HasNullError() {
		return ""
	}

	var b strings.Builder
	b.WriteString(HeaderPrefix)
	b.WriteString(fmt.Sprintf("%s (Code: %d)%s", e.errType.Name(), e.errType.Code(), Newline))
	e.writeCompileDetails(&b)

	return b.String()
}

// CompileWithStack builds a complete formatted report including stack trace.
func (e *AppError) CompileWithStack() string {
	if e.HasNullError() {
		return ""
	}

	compiled := e.Compile()
	if len(e.stack) > 0 {
		return fmt.Sprintf("%s%sStack Trace:%s%s%s", compiled, SectionPrefix, Newline, e.stack.String(), Newline)
	}

	return compiled
}

// CompiledError returns a standard error with the compiled string representation.
func (e *AppError) CompiledError() error {
	if e.HasNullError() {
		return nil
	}

	return errors.New(e.Compile())
}
