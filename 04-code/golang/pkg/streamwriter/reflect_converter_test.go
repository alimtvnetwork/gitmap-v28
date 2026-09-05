package streamwriter_test

import (
	"reflect"
	"testing"

	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/streamwriter"
)

type reflectSampleUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func TestReflect_UnmarshalTo(t *testing.T) {
	// 1. Success case
	rawJSON := []byte(`{"id":101,"name":"Alice","role":"Admin"}`)
	var user reflectSampleUser

	err := streamwriter.Reflect.UnmarshalTo(rawJSON, &user)
	if err != nil {
		t.Fatalf("expected successful unmarshal, got: %v", err)
	}

	if user.ID != 101 || user.Name != "Alice" || user.Role != "Admin" {
		t.Fatalf("unexpected unmarshaled user: %+v", user)
	}

	// 2. Empty payload bytes error
	err = streamwriter.Reflect.UnmarshalTo([]byte{}, &user)
	if err == nil || err.Type() != errtype.Precondition {
		t.Fatalf("expected Precondition error on empty bytes, got: %v", err)
	}

	// 3. Nil target pointer error
	err = streamwriter.Reflect.UnmarshalTo(rawJSON, nil)
	if err == nil || err.Type() != errtype.Precondition {
		t.Fatalf("expected Precondition error on nil target, got: %v", err)
	}

	// 4. Non-pointer target error
	err = streamwriter.Reflect.UnmarshalTo(rawJSON, user)
	if err == nil || err.Type() != errtype.Precondition {
		t.Fatalf("expected Precondition error on non-pointer target, got: %v", err)
	}

	// 5. Invalid JSON syntax error
	err = streamwriter.Reflect.UnmarshalTo([]byte(`{invalid-json`), &user)
	if err == nil || err.Type() != errtype.Serialization {
		t.Fatalf("expected Serialization error on invalid JSON, got: %v", err)
	}
}

func TestUnmarshalToType(t *testing.T) {
	rawJSON := []byte(`{"id":202,"name":"Bob","role":"Developer"}`)
	user, err := streamwriter.UnmarshalToType[reflectSampleUser](rawJSON)
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	if user.ID != 202 || user.Name != "Bob" {
		t.Fatalf("unexpected user: %+v", user)
	}

	// Empty bytes
	_, err = streamwriter.UnmarshalToType[reflectSampleUser](nil)
	if err == nil {
		t.Fatalf("expected error on nil bytes")
	}
}

func TestReflect_ReducePointer(t *testing.T) {
	val := "deep-value"
	ptr1 := &val
	ptr2 := &ptr1
	ptr3 := &ptr2

	reduced := streamwriter.Reflect.ReducePointer(ptr3)
	str, ok := reduced.(string)
	if !ok || str != "deep-value" {
		t.Fatalf("expected 'deep-value', got: %v", reduced)
	}

	// Non pointer
	num := 42
	if streamwriter.Reflect.ReducePointer(num) != 42 {
		t.Fatalf("expected 42")
	}

	// Nil pointer chain
	var nilPtr *string
	nilChain := &nilPtr
	if streamwriter.Reflect.ReducePointer(nilChain) != nil {
		t.Fatalf("expected nil for broken pointer chain")
	}
}

func TestReflect_ToInterfaces(t *testing.T) {
	nums := []int{10, 20, 30}
	infs := streamwriter.Reflect.ToInterfaces(nums)

	if len(infs) != 3 {
		t.Fatalf("expected length 3, got: %d", len(infs))
	}

	if infs[0] != 10 || infs[1] != 20 || infs[2] != 30 {
		t.Fatalf("unexpected slice elements: %v", infs)
	}

	// Pointer to slice
	ptrSlice := &nums
	infs2 := streamwriter.Reflect.ToInterfaces(ptrSlice)
	if len(infs2) != 3 {
		t.Fatalf("expected length 3 from pointer to slice")
	}

	// Non-slice input
	empty := streamwriter.Reflect.ToInterfaces("not-a-slice")
	if len(empty) != 0 {
		t.Fatalf("expected empty slice for non-slice input")
	}
}

func TestReflect_InspectType(t *testing.T) {
	user := reflectSampleUser{ID: 1, Name: "Admin", Role: "Super"}
	info := streamwriter.Reflect.InspectType(&user)

	if !info.IsPointer {
		t.Fatalf("expected IsPointer to be true")
	}

	if !info.IsStruct {
		t.Fatalf("expected IsStruct to be true")
	}

	if info.FieldCount != 3 {
		t.Fatalf("expected FieldCount 3, got: %d", info.FieldCount)
	}

	if info.Kind != reflect.Struct {
		t.Fatalf("expected Kind Struct, got: %v", info.Kind)
	}

	// Nil inspect
	nilInfo := streamwriter.Reflect.InspectType(nil)
	if !nilInfo.IsNil {
		t.Fatalf("expected IsNil true for nil value")
	}
}
