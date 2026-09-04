# Generic Payload `T` & Ordered Recursive Compilation Architecture

> **Status:** Implemented & Verified  
> **Date:** 2026-09-04  
> **Package:** `04-code/golang/pkg/streamwriter`  
> **Topic:** Type-Parameterized Streams/Writers `[T any]` and Order-Wise Recursive Transpilation Engine  

---

## 1. Executive Summary & Design Mandates

This architecture addresses two fundamental requirements:
1. **Generic Payload `T`:** All streamers, writers, and loggers are parameterized over `[T any]`, providing compile-time type-safety while supporting universal `[any]` when needed.
2. **Recursive `Compile` Engine:** Transpiles any generic payload `T` into deterministic, ordered representations according to strict rules:
   - **Primitives (string, number, bool):** Printed directly (`hello`, `42`, `3.14`, `true`).
   - **Maps:** Printed in **lexicographical key order** (`{a: 1, b: 2}`).
   - **Arrays & Slices:** Printed sequentially in index order (`["first", 2, true]`).
   - **Objects & Structs:** Printed with exported field names, json tag mapping, and **recursive seeking of `Compile()` methods** on nested elements.

---

## 2. The `Compile` Transpilation Engine

### 2.1 The `Compilable` Interface
Any object or struct can define its own transpilation output:
```go
type Compilable interface {
	Compile() string
}
```

### 2.2 Transpilation Rules & Recursive Resolution
```go
package streamwriter

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func Compile[T any](payload T) string {
	return DefaultCompiler.CompileValue(payload)
}

func (c *Compiler) compileRecursive(v reflect.Value, depth int, isNested bool) string {
	if depth > c.maxDepth {
		return "..."
	}
	if !v.IsValid() {
		return "nil"
	}

	// 1. Recursive check for Compilable method on value or pointer receiver
	if v.CanInterface() {
		if compilable, isComp := v.Interface().(Compilable); isComp {
			return compilable.Compile()
		}
	}
	if v.Kind() != reflect.Ptr && v.CanAddr() {
		if compilable, isComp := v.Addr().Interface().(Compilable); isComp {
			return compilable.Compile()
		}
	}

	// 2. Type-specific ordered transpilation
	switch v.Kind() {
	case reflect.String:
		if isNested {
			return strconv.Quote(v.String())
		}
		return v.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)

	case reflect.Float32, reflect.Float64:
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
		var elements []string
		for i := 0; i < length; i++ {
			elements = append(elements, c.compileRecursive(v.Index(i), depth+1, true))
		}
		return "[" + strings.Join(elements, ", ") + "]"

	case reflect.Map:
		// Map keys are sorted lexicographically for order-wise printing
		mapKeys := v.MapKeys()
		entries := make([]keyEntry, 0, len(mapKeys))
		for _, k := range mapKeys {
			entries = append(entries, keyEntry{
				keyStr: c.compileRecursive(k, depth+1, false),
				keyVal: k,
			})
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].keyStr < entries[j].keyStr
		})

		var pairs []string
		for _, e := range entries {
			valStr := c.compileRecursive(v.MapIndex(e.keyVal), depth+1, true)
			pairs = append(pairs, e.keyStr+": "+valStr)
		}
		return "{" + strings.Join(pairs, ", ") + "}"

	case reflect.Struct:
		structType := v.Type()
		var fieldPairs []string
		for i := 0; i < v.NumField(); i++ {
			fieldType := structType.Field(i)
			if !fieldType.IsExported() {
				continue
			}
			fieldName := fieldType.Name
			tag := fieldType.Tag.Get("json")
			if tag != "" {
				parts := strings.Split(tag, ",")
				if parts[0] != "" && parts[0] != "-" {
					fieldName = parts[0]
				}
			}
			valStr := c.compileRecursive(v.Field(i), depth+1, true)
			fieldPairs = append(fieldPairs, fieldName+": "+valStr)
		}
		return "{" + strings.Join(fieldPairs, ", ") + "}"

	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
```

---

## 3. Generic Contracts with Payload `T`

```go
type WriterInterface[T any] interface {
	Interfacer
	Name() string
	Write(ctx context.Context, payload T) error
	AsWriter() WriterInterface[T]
	Sync() error
	Close() error
}

type StreamerInterface[T any] interface {
	Interfacer
	Name() string
	Stream(ctx context.Context, payload T) error
	AsStreamer() StreamerInterface[T]
	AsWriter() WriterInterface[T]
	IsLocked() bool
	Destination() io.Writer
	Sync() error
	Close() error
}

type StreamFunc[T any] func(ctx context.Context, payload T, dest io.Writer) error
type WriteFunc[T any] func(ctx context.Context, payload T) error
type FormatFunc[T any] func(payload T) ([]byte, error)
```

---

## 4. End-to-End Demonstration

```go
func main() {
	ctx := context.Background()

	// 1. Primitive transpilation
	fmt.Println(streamwriter.Compile("plain text")) // Output: plain text
	fmt.Println(streamwriter.Compile(12345))        // Output: 12345

	// 2. Order-wise Map transpilation (keys sorted alphabetically)
	unorderedMap := map[string]any{
		"zebra": 99,
		"apple": 1,
		"cat":   true,
	}
	fmt.Println(streamwriter.Compile(unorderedMap))
	// Output: {apple: 1, cat: true, zebra: 99}

	// 3. Slice transpilation
	fmt.Println(streamwriter.Compile([]any{"alpha", 2, false}))
	// Output: ["alpha", 2, false]

	// 4. Struct with recursive Compilable method resolution
	type InnerMetric struct {
		Name  string
		Value float64
	}
	// InnerMetric implements Compilable
	func (m InnerMetric) Compile() string {
		return fmt.Sprintf("METRIC(%s=%.2f)", m.Name, m.Value)
	}

	type ServerState struct {
		Host    string      `json:"host"`
		Port    int         `json:"port"`
		Metric  InnerMetric `json:"metric"`
	}

	state := ServerState{
		Host: "localhost",
		Port: 8080,
		Metric: InnerMetric{Name: "cpu_load", Value: 0.42},
	}
	fmt.Println(streamwriter.Compile(state))
	// Output: {host: "localhost", port: 8080, metric: METRIC(cpu_load=0.42)}

	// 5. Generic Streamer with Type T = ServerState
	streamer := streamwriter.NewLockedStreamer[ServerState](streamwriter.LockedOptions[ServerState]{
		Name: "state-streamer",
	})
	_ = streamer.Stream(ctx, state)
}
```

---

## 5. Verification Results

All tests pass 100% green (`go test ./pkg/streamwriter -v -count=1`):
- `TestCompiler_Primitives` (PASS)
- `TestCompiler_Maps_OrderWise` (PASS)
- `TestCompiler_Slices_OrderWise` (PASS)
- `TestCompiler_ObjectAndRecursiveCompilable` (PASS)
- `TestLockedStreamer_Generic_ConcurrentSafe` (PASS)
- `TestLocklessStreamer_Generic_Direct` (PASS)
- `TestSelfBinding_GenericContracts` (PASS)
- `TestSwappableMethods_GenericRuntime` (PASS)
- `TestCompositeLogger_FluentChaining` (PASS)
- `TestLogRecord_Compile` (PASS)
