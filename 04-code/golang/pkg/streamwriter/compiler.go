package streamwriter

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Compilable represents any object or struct that knows how to compile itself.
type Compilable interface {
	Compile() string
}

// Compiler transpiles generic payloads into strictly ordered, deterministic string outputs.
type Compiler struct {
	maxDepth int
}

// NewCompiler creates a new Compiler with configurable max depth.
func NewCompiler() *Compiler {
	return &Compiler{maxDepth: 32}
}

// DefaultCompiler is the global package compiler instance.
var DefaultCompiler = NewCompiler()

// Compile is the universal generic helper that transpiles any payload T into an ordered string.
func Compile[T any](payload T) string {
	return DefaultCompiler.CompileValue(payload)
}

// CompileValue transpiles any value recursively according to specific order rules.
func (c *Compiler) CompileValue(val any) string {
	return c.compileRecursive(reflect.ValueOf(val), 0, false)
}

func (c *Compiler) compileRecursive(v reflect.Value, depth int, isNested bool) string {
	if depth > c.maxDepth {
		return "..."
	}

	// 1. Nil check
	if !v.IsValid() {
		return "nil"
	}

	// 2. Check if the value or pointer to value implements Compilable interface
	compilable, isComp := extractCompilable(v)
	if isComp {
		return compilable.Compile()
	}

	// 3. Type-specific ordered transpilation
	switch v.Kind() {
	case reflect.String:
		if isNested {
			return strconv.Quote(v.String())
		}
		return v.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)

	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)

	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)

	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"

	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return "nil"
		}
		return c.compileRecursive(v.Elem(), depth+1, isNested)

	case reflect.Slice, reflect.Array:
		length := v.Len()
		if length == 0 {
			return "[]"
		}
		var elements []string
		for i := 0; i < length; i++ {
			elemStr := c.compileRecursive(v.Index(i), depth+1, true)
			elements = append(elements, elemStr)
		}
		return "[" + strings.Join(elements, ", ") + "]"

	case reflect.Map:
		length := v.Len()
		if length == 0 {
			return "{}"
		}

		// Sort keys lexicographically for deterministic order-wise printing
		type keyEntry struct {
			keyStr string
			keyVal reflect.Value
		}

		mapKeys := v.MapKeys()
		entries := make([]keyEntry, 0, len(mapKeys))
		for _, k := range mapKeys {
			kStr := c.compileRecursive(k, depth+1, false)
			entries = append(entries, keyEntry{keyStr: kStr, keyVal: k})
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].keyStr < entries[j].keyStr
		})

		var pairs []string
		for _, entry := range entries {
			valStr := c.compileRecursive(v.MapIndex(entry.keyVal), depth+1, true)
			pairs = append(pairs, entry.keyStr+": "+valStr)
		}
		return "{" + strings.Join(pairs, ", ") + "}"

	case reflect.Struct:
		structType := v.Type()
		numFields := v.NumField()
		if numFields == 0 {
			return "{}"
		}

		var fieldPairs []string
		for i := 0; i < numFields; i++ {
			fieldType := structType.Field(i)

			// Skip unexported fields
			if !fieldType.IsExported() {
				continue
			}

			fieldName, shouldSkip := resolveJSONFieldName(fieldType)
			if shouldSkip {
				continue
			}

			fieldVal := v.Field(i)
			valStr := c.compileRecursive(fieldVal, depth+1, true)
			fieldPairs = append(fieldPairs, fieldName+": "+valStr)
		}

		return "{" + strings.Join(fieldPairs, ", ") + "}"

	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

func extractCompilable(v reflect.Value) (Compilable, bool) {
	comp, ok := extractFromInterface(v)
	if ok {
		return comp, true
	}

	return extractFromAddr(v)
}

func extractFromInterface(v reflect.Value) (Compilable, bool) {
	if !v.CanInterface() {
		return nil, false
	}

	comp, ok := v.Interface().(Compilable)

	return comp, ok
}

func extractFromAddr(v reflect.Value) (Compilable, bool) {
	if v.Kind() == reflect.Ptr || !v.CanAddr() {
		return nil, false
	}

	comp, ok := v.Addr().Interface().(Compilable)

	return comp, ok
}

func resolveJSONFieldName(fieldType reflect.StructField) (string, bool) {
	tag := fieldType.Tag.Get("json")
	if tag == "" {
		return fieldType.Name, false
	}

	parts := strings.Split(tag, ",")
	if parts[0] == "-" {
		return "", true
	}

	if parts[0] != "" {
		return parts[0], false
	}

	return fieldType.Name, false
}
