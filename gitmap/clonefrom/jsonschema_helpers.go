package clonefrom

// JSON-Schema construction helpers used by jsonschema.go. Kept in
// a sibling file so EmitReportSchema / EmitInputSchema stay readable
// and the per-file line cap is preserved.
//
// Schemas are built as map[string]any trees because (a) encoding/json
// sorts map keys alphabetically on marshal — giving us byte-stable
// output without an ordered-map dependency — and (b) the tree shape
// per JSON-Schema draft 2020-12 maps cleanly to dynamic objects.

import "github.com/alimtvnetwork/gitmap-v27/gitmap/constants"

// propKV pairs a property name with its sub-schema. Used purely as
// an argument-passing convenience for orderedProps; the resulting
// map is unordered (encoding/json re-sorts on emit).
type propKV struct {
	name   string
	schema map[string]any
}

func kv(name string, schema map[string]any) propKV {
	return propKV{name: name, schema: schema}
}

// orderedProps folds a slice of propKV into the canonical
// `properties` object expected by JSON Schema. Name "ordered" is
// historical — encoding/json alphabetizes map keys regardless.
func orderedProps(items ...propKV) map[string]any {
	out := make(map[string]any, len(items))
	for _, it := range items {
		out[it.name] = it.schema
	}

	return out
}

// rootSchema builds the top-level object schema with $schema, $id,
// title, and (optionally) properties + required. Caller may delete
// or override fields for non-object root cases (e.g. arrays).
func rootSchema(id, title string, properties map[string]any, required []string) map[string]any {
	out := map[string]any{
		"$schema":              constants.JSONSchemaDialect2020_12,
		"$id":                  id,
		"title":                title,
		"type":                 "object",
		"additionalProperties": false,
	}
	if properties != nil {
		out["properties"] = properties
	}
	if len(required) > 0 {
		out["required"] = required
	}

	return out
}

func objectSchema(properties map[string]any, required []string, description string) map[string]any {
	out := map[string]any{
		"type":                 "object",
		"properties":           properties,
		"additionalProperties": false,
		"description":          description,
	}
	if len(required) > 0 {
		out["required"] = required
	}

	return out
}

func arraySchema(items map[string]any, description string) map[string]any {
	return map[string]any{
		"type":        "array",
		"items":       items,
		"description": description,
	}
}

func strSchema(description string) map[string]any {
	return map[string]any{"type": "string", "description": description}
}

func intSchema(description string) map[string]any {
	return map[string]any{"type": "integer", "minimum": 0, "description": description}
}

func numSchema(description string) map[string]any {
	return map[string]any{"type": "number", "minimum": 0, "description": description}
}

func enumSchema(description string, values []string) map[string]any {
	asAny := make([]any, 0, len(values))
	for _, v := range values {
		asAny = append(asAny, v)
	}

	return map[string]any{
		"type":        "string",
		"enum":        asAny,
		"description": description,
	}
}

func constIntSchema(value int, description string) map[string]any {
	return map[string]any{
		"type":        "integer",
		"const":       value,
		"description": description,
	}
}

// scanFieldSchema returns a permissive sub-schema for one accepted
// clone-now input field. We deliberately do NOT pin types per field
// here (numbers vs strings) because the parser tolerates either form
// for `id` / `depth` (string-encoded integers from CSV are accepted).
// `["string","integer","null"]` covers every observed shape without
// rejecting valid manifests.
func scanFieldSchema(name string) map[string]any {
	return map[string]any{
		"type":        []string{"string", "integer", "number", "null"},
		"description": "Accepted scan-record field: " + name,
	}
}

// anyOfRequired builds an `anyOf` clause where each branch requires
// exactly one of the named fields. Used to express "at least one of
// httpsUrl/sshUrl must be present" without forbidding either alone.
func anyOfRequired(fields []string) []any {
	out := make([]any, 0, len(fields))
	for _, f := range fields {
		out = append(out, map[string]any{"required": []string{f}})
	}

	return out
}
