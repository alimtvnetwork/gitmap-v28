package appfault

import (
	"fmt"
	"strings"
)

// appendYamlLine appends key and value if val is non-empty.
func appendYamlLine(b *strings.Builder, key, val string) {
	if len(val) > 0 {
		b.WriteString(fmt.Sprintf("%s: \"%s\"\n", key, val))
	}
}

// appendBasicYamlFields appends basic header fields to YAML builder.
func appendBasicYamlFields(b *strings.Builder, m AppErrorDataModel) {
	b.WriteString(fmt.Sprintf("Type: %d\n", m.Type.Code()))
	appendYamlLine(b, "TypeName", m.Type.Name())
	appendYamlLine(b, "Message", m.Message)
}

// appendDetailYamlFields appends caller, stack, cause, and status to YAML.
func appendDetailYamlFields(b *strings.Builder, m AppErrorDataModel) {
	appendYamlLine(b, "Caller", m.Caller.String())
	appendYamlLine(b, "Cause", m.Cause)
	if m.StatusCode > 0 {
		b.WriteString(fmt.Sprintf("StatusCode: %d\n", m.StatusCode))
	}
}

// appendCtxYamlFields appends Ctx entries to YAML builder.
func appendCtxYamlFields(b *strings.Builder, ctx ContextMap) {
	if len(ctx) == 0 {
		return
	}

	b.WriteString("Ctx:\n")
	for _, k := range ctx.Keys() {
		b.WriteString(fmt.Sprintf("  %s: %v\n", k, ctx[k]))
	}
}

// ToYAMLString serializes the AppError to clean YAML representation.
func (e *AppError) ToYAMLString() string {
	if e == nil {
		return "{}\n"
	}

	m := e.ToDataModel()
	var b strings.Builder
	appendBasicYamlFields(&b, m)
	appendDetailYamlFields(&b, m)
	appendCtxYamlFields(&b, m.Ctx)

	return b.String()
}

// ToYAML exports the AppError as YAML bytes.
func (e *AppError) ToYAML() ([]byte, error) {
	return []byte(e.ToYAMLString()), nil
}
