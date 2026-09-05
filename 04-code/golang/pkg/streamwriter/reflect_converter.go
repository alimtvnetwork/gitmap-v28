package streamwriter

import (
	"encoding/json"
	"fmt"
	"reflect"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// ReflectTypeInfo describes structural metadata inspected from any runtime value.
type ReflectTypeInfo struct {
	Name       string
	PkgPath    string
	Kind       reflect.Kind
	IsPointer  bool
	IsSlice    bool
	IsStruct   bool
	IsNil      bool
	FieldCount int
}

// reflectConverterSingleton provides reflection, dynamic unmarshaling, pointer reduction,
// and interface slice transformation utilities.
type reflectConverterSingleton struct{}

// Reflect is the global singleton instance for reflection and dynamic conversion.
var Reflect = reflectConverterSingleton{}

// UnmarshalTo deserializes raw JSON byte data directly into the provided destination pointer.
// It verifies that targetPtr is a non-nil pointer and returns *appfault.AppError on failure.
func (reflectConverterSingleton) UnmarshalTo(data []byte, targetPtr any) *appfault.AppError {
	if len(data) == 0 {
		return appfault.New(errtype.Precondition, "unmarshal payload bytes cannot be empty")
	}

	if targetPtr == nil {
		return appfault.New(errtype.Precondition, "unmarshal target pointer cannot be nil")
	}

	rv := reflect.ValueOf(targetPtr)
	if rv.Kind() != reflect.Ptr {
		return appfault.New(errtype.Precondition, fmt.Sprintf("unmarshal target must be a pointer, got %s", rv.Kind()))
	}

	if rv.IsNil() {
		return appfault.New(errtype.Precondition, "unmarshal target pointer points to nil")
	}

	if err := json.Unmarshal(data, targetPtr); err != nil {
		return appfault.Wrap(errtype.Serialization, err, "failed to unmarshal JSON into target")
	}

	return nil
}

// UnmarshalToType deserializes JSON bytes into a new instance of generic type T.
func UnmarshalToType[T any](data []byte) (T, *appfault.AppError) {
	var result T
	if len(data) == 0 {
		return result, appfault.New(errtype.Precondition, "unmarshal payload bytes cannot be empty")
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, appfault.Wrap(errtype.Serialization, err, "failed to unmarshal JSON into target type")
	}

	return result, nil
}

// ReducePointer recursively unwraps nested pointers (e.g., ***T -> T) down to the underlying value.
// If any pointer along the chain is nil, it returns nil without panicking.
func (reflectConverterSingleton) ReducePointer(anyItem any) any {
	return Reflect.ReducePointerLevel(anyItem, 10)
}

// ReducePointerLevel unwraps nested pointers up to maxLevel depths.
func (reflectConverterSingleton) ReducePointerLevel(anyItem any, maxLevel int) any {
	if anyItem == nil {
		return nil
	}

	rv := reflect.ValueOf(anyItem)
	level := 0

	for rv.Kind() == reflect.Ptr && level < maxLevel {
		if rv.IsNil() {
			return nil
		}

		rv = rv.Elem()
		level++
	}

	if !rv.IsValid() {
		return nil
	}

	return rv.Interface()
}

// ToInterfaces converts any slice or array into a slice of empty interfaces ([]any).
// If the input is not a slice, array, or is nil, it returns an empty slice.
func (reflectConverterSingleton) ToInterfaces(sliceOrArray any) []any {
	if sliceOrArray == nil {
		return []any{}
	}

	rv := reflect.ValueOf(sliceOrArray)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return []any{}
		}

		rv = rv.Elem()
	}

	k := rv.Kind()
	if k != reflect.Slice && k != reflect.Array {
		return []any{}
	}

	length := rv.Len()
	if length == 0 {
		return []any{}
	}

	result := make([]any, length)
	for i := 0; i < length; i++ {
		item := rv.Index(i)
		if item.Kind() == reflect.Ptr && !item.IsNil() {
			result[i] = item.Interface()
		} else if item.IsValid() && item.CanInterface() {
			result[i] = item.Interface()
		} else {
			result[i] = nil
		}
	}

	return result
}

// ArgsToReflectValues converts a slice of any values into []reflect.Value.
func (reflectConverterSingleton) ArgsToReflectValues(args []any) []reflect.Value {
	if len(args) == 0 {
		return []reflect.Value{}
	}

	list := make([]reflect.Value, len(args))
	for i, arg := range args {
		list[i] = reflect.ValueOf(arg)
	}

	return list
}

// ReflectValuesToInterfaces converts []reflect.Value back into []any.
func (reflectConverterSingleton) ReflectValuesToInterfaces(reflectValues []reflect.Value) []any {
	if len(reflectValues) == 0 {
		return []any{}
	}

	list := make([]any, len(reflectValues))
	for i, rv := range reflectValues {
		list[i] = Reflect.ReflectValueToAnyValue(rv)
	}

	return list
}

// ReflectValueToAnyValue extracts an interface{} from reflect.Value safely.
func (reflectConverterSingleton) ReflectValueToAnyValue(rv reflect.Value) any {
	if !rv.IsValid() {
		return nil
	}

	k := rv.Kind()
	switch k {
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return nil
		}

		return rv.Elem().Interface()

	case reflect.String:
		return rv.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint()

	case reflect.Bool:
		return rv.Bool()

	case reflect.Float32, reflect.Float64:
		return rv.Float()

	default:
		if rv.CanInterface() {
			return rv.Interface()
		}

		return nil
	}
}

// InspectType examines the runtime type metadata of any value.
func (reflectConverterSingleton) InspectType(anyItem any) ReflectTypeInfo {
	if anyItem == nil {
		return ReflectTypeInfo{
			Name:  "nil",
			Kind:  reflect.Invalid,
			IsNil: true,
		}
	}

	t := reflect.TypeOf(anyItem)
	v := reflect.ValueOf(anyItem)

	isPtr := t.Kind() == reflect.Ptr
	isNilVal := false

	if isPtr {
		isNilVal = v.IsNil()
		if !isNilVal {
			t = t.Elem()
			v = v.Elem()
		}
	}

	isSlice := t.Kind() == reflect.Slice || t.Kind() == reflect.Array
	isStruct := t.Kind() == reflect.Struct
	fieldCount := 0

	if isStruct {
		fieldCount = t.NumField()
	}

	return ReflectTypeInfo{
		Name:       t.Name(),
		PkgPath:    t.PkgPath(),
		Kind:       t.Kind(),
		IsPointer:  isPtr,
		IsSlice:    isSlice,
		IsStruct:   isStruct,
		IsNil:      isNilVal,
		FieldCount: fieldCount,
	}
}
