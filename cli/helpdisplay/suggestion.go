package helpdisplay

import "strings"

// RenderContext carries cached runtime values interpolated into suggestions
// at render time — e.g. detected machine IPs, repo names, versions.
type RenderContext struct {
	values map[string]string
}

// NewRenderContext builds a context pre-seeded with cached values.
//
// Callers seed values lazily: detect once (e.g. machine IPs), cache here,
// then reuse across renders. Suggestions show real values
// ("ssh user@192.168.1.22") instead of placeholders ("ssh user@<ip>").
//
// Example: NewRenderContext(map[string]string{"primary-ip": "192.168.1.22"})
func NewRenderContext(values map[string]string) RenderContext {
	return RenderContext{values: values}
}

// Get returns the cached value for key, or "" when the key is missing.
func (c RenderContext) Get(key string) string {
	value, found := c.values[key]
	if !found {
		return ""
	}
	return value
}

// Suggestion is one actionable hint line; variables interpolate {{name}}
// placeholders in text from a RenderContext at render time.
type Suggestion struct {
	text      string
	variables map[string]string
}

// NewSuggestion builds a suggestion; variables maps placeholder name to
// RenderContext key (e.g. {"ip": "primary-ip"} with text "ssh user@{{ip}}").
func NewSuggestion(text string, variables map[string]string) Suggestion {
	return Suggestion{text: text, variables: variables}
}

// Render replaces every {{name}} in text with ctx.Get(variables[name]).
// Placeholders that cannot resolve (missing variable entry or missing
// context value) are left untouched in the output.
func (s Suggestion) Render(ctx RenderContext) string {
	rendered := s.text
	for name, key := range s.variables {
		rendered = replacePlaceholder(rendered, name, ctx.Get(key))
	}
	return rendered
}

// replacePlaceholder swaps one {{name}} placeholder for value; an empty
// (missing) value leaves the placeholder in place.
func replacePlaceholder(text, name, value string) string {
	if len(value) == 0 {
		return text
	}
	return strings.ReplaceAll(text, "{{"+name+"}}", value)
}
