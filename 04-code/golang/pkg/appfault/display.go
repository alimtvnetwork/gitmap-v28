package appfault

import (
	"fmt"
	"strings"
)

// formatBasicError returns formatted type name, code, and message.
func formatBasicError(e *AppError) string {
	return fmt.Sprintf("[%s:%d] %s", e.errType.Name(), e.errType.Code(), e.message)
}

// appendCallerAndCause appends caller site and cause to formatted string.
func appendCallerAndCause(base, caller string, cause error) string {
	if len(caller) > 0 {
		base += fmt.Sprintf(" (at=%s)", caller)
	}

	if cause != nil {
		base += fmt.Sprintf(" (cause=%v)", cause)
	}

	return base
}

// Error implements the standard Go error interface.
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	return appendCallerAndCause(formatBasicError(e), e.caller.String(), e.cause)
}

// appendHeader writes diagnostic header info.
func appendHeader(b *strings.Builder, e *AppError) {
	b.WriteString(fmt.Sprintf("ERROR: [%s:%d] %s\n", e.errType.Name(), e.errType.Code(), e.message))
	if !e.caller.IsEmpty() {
		b.WriteString(fmt.Sprintf("CALLER: %s\n", e.caller.String()))
	}

	if e.cause != nil {
		b.WriteString(fmt.Sprintf("CAUSE: %v\n", e.cause))
	}
}

// appendContextAndStack writes context map and stack trace.
func appendContextAndStack(b *strings.Builder, ctx ContextMap, stack StackTrace) {
	if len(ctx) > 0 {
		b.WriteString(fmt.Sprintf("CONTEXT: %s\n", ctx.Format()))
	}

	if len(stack) > 0 {
		b.WriteString("STACK TRACE:\n" + stack.String())
	}
}

// FullString returns a comprehensive diagnostic dump of the AppError.
func (e *AppError) FullString() string {
	if e == nil {
		return ""
	}

	var b strings.Builder
	appendHeader(&b, e)
	appendContextAndStack(&b, e.ctx, e.stack)

	return b.String()
}

// appendMarkdownCauseAndStack writes cause and codeblock stack trace.
func appendMarkdownCauseAndStack(b *strings.Builder, cause error, stack StackTrace) {
	if cause != nil {
		b.WriteString(fmt.Sprintf("- **Cause:** `%v`\n", cause))
	}

	if len(stack) > 0 {
		b.WriteString("\n```\n" + stack.String() + "```\n")
	}
}

// ToClipboard returns a Markdown formatted report for AI analysis.
func (e *AppError) ToClipboard() string {
	if e == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("### Error Report\n\n- **Type:** `%s (%d)`\n- **Message:** %s\n", e.errType.Name(), e.errType.Code(), e.message))
	appendMarkdownCauseAndStack(&b, e.cause, e.stack)

	return b.String()
}

// DisplayError prints a terminal banner representation.
func (e *AppError) DisplayError() {
	if e != nil {
		fmt.Printf("❌ [%s:%d] %s (at %s)\n", e.errType.Name(), e.errType.Code(), e.message, e.caller.String())
	}
}

// DefaultFaultFormatter formats the error into a clean, human-readable terminal line.
func DefaultFaultFormatter(e *AppError) string {
	if e == nil {
		return ""
	}

	callerInfo := ""
	if !e.caller.IsEmpty() {
		callerInfo = fmt.Sprintf(" (at %s)", e.caller.String())
	}

	return fmt.Sprintf("❌ [%s:%d] %s%s", e.errType.Name(), e.errType.Code(), e.message, callerInfo)
}

// Print outputs the default formatted fault representation to standard output.
func (e *AppError) Print() {
	if e != nil {
		fmt.Println(e.Format(DefaultFaultFormatter))
	}
}

// Format formats the AppError using a specified or default formatter.
func (e *AppError) Format(formatter FaultFormatter) string {
	if e == nil {
		return ""
	}

	if formatter != nil {
		return formatter(e)
	}

	return DefaultFaultFormatter(e)
}

// FormatStdout produces a rich terminal banner with icon, HTTP status, and caller info.
func FormatStdout(e *AppError) string {
	if e == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("❌ ERROR [%s:%d] %s (HTTP %d)", e.errType.Name(), e.errType.Code(), e.message, e.StatusCode()))
	if !e.caller.IsEmpty() {
		b.WriteString(fmt.Sprintf("\n   Caller:  %s", e.caller.String()))
	}

	if e.cause != nil {
		b.WriteString(fmt.Sprintf("\n   Cause:   %v", e.cause))
	}

	if len(e.ctx) > 0 {
		b.WriteString(fmt.Sprintf("\n   Context: %s", e.ctx.Format()))
	}

	return b.String()
}

// FormatJson produces a structured, indented JSON representation of the error.
func FormatJson(e *AppError) string {
	if e == nil {
		return "{}"
	}

	return e.ToJsonString()
}

// FormatJSON is an alias for FormatJson.
func FormatJSON(e *AppError) string {
	return FormatJson(e)
}

// FormatTextLog produces a single-line structured log string suitable for file logging.
func FormatTextLog(e *AppError) string {
	if e == nil {
		return ""
	}

	callerStr := "unknown"
	if !e.caller.IsEmpty() {
		callerStr = e.caller.String()
	}

	logLine := fmt.Sprintf("[ERROR] [%s:%d] status=%d caller=%q msg=%q",
		e.errType.Name(), e.errType.Code(), e.StatusCode(), callerStr, e.message)

	if e.cause != nil {
		logLine += fmt.Sprintf(" cause=%q", e.cause.Error())
	}

	if len(e.ctx) > 0 {
		logLine += fmt.Sprintf(" ctx=%q", e.ctx.Format())
	}

	return logLine
}

// FormatStdout returns the rich terminal banner formatting.
func (e *AppError) FormatStdout() string {
	return FormatStdout(e)
}

// FormatJson returns the formatted JSON string.
func (e *AppError) FormatJson() string {
	return FormatJson(e)
}

// FormatJSON is an alias for FormatJson.
func (e *AppError) FormatJSON() string {
	return e.FormatJson()
}

// FormatTextLog returns the structured single-line log format.
func (e *AppError) FormatTextLog() string {
	return FormatTextLog(e)
}

// PrintStdout prints the error formatted for console stdout.
func (e *AppError) PrintStdout() {
	if e != nil {
		fmt.Println(e.FormatStdout())
	}
}

// PrintJson prints the error formatted as JSON.
func (e *AppError) PrintJson() {
	if e != nil {
		fmt.Println(e.FormatJson())
	}
}

// PrintJSON is an alias for PrintJson.
func (e *AppError) PrintJSON() {
	e.PrintJson()
}

// PrintLog prints the error formatted as a structured log line.
func (e *AppError) PrintLog() {
	if e != nil {
		fmt.Println(e.FormatTextLog())
	}
}

// PrintWith outputs the error using a customized formatter.
func (e *AppError) PrintWith(formatter FaultFormatter) {
	if e != nil {
		fmt.Println(e.Format(formatter))
	}
}
