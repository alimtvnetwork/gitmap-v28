package config

import (
	"os"
	"testing"
)

func TestVariablesStoreAndExpansion(t *testing.T) {
	// Set global variable
	if err := SetVariable("global", "MY_TEST_VAR", "hello_world"); err != nil {
		t.Fatalf("SetVariable failed: %v", err)
	}

	val, found := GetVariable("global", "MY_TEST_VAR")
	if !found || val != "hello_world" {
		t.Fatalf("expected hello_world, got %s (found: %v)", val, found)
	}

	// Set scoped variable
	if err := SetVariable("agy", "TARGET_DIR", "d:/work/custom"); err != nil {
		t.Fatalf("SetVariable scoped failed: %v", err)
	}

	scopedVal, scopedFound := GetVariable("agy", "TARGET_DIR")
	if !scopedFound || scopedVal != "d:/work/custom" {
		t.Fatalf("expected d:/work/custom, got %s", scopedVal)
	}

	// Test expansion
	expanded := ExpandVariables("folder is $MY_TEST_VAR and path is ${TARGET_DIR}", "agy")
	expected := "folder is hello_world and path is d:/work/custom"
	if expanded != expected {
		t.Fatalf("expected %q, got %q", expected, expanded)
	}

	// Test fallback to OS env
	t.Setenv("SYSTEM_PORT", "9090")
	expandedEnv := ExpandVariables("port=$SYSTEM_PORT", "global")
	if expandedEnv != "port=9090" {
		t.Fatalf("expected port=9090, got %s", expandedEnv)
	}

	// Test delete
	_ = DeleteVariable("global", "MY_TEST_VAR")
	_ = DeleteVariable("agy", "TARGET_DIR")
}

func TestExportVariables(t *testing.T) {
	_ = SetVariable("global", "EXP_VAR_1", "v1")
	_ = SetVariable("testsub", "EXP_VAR_2", "v2")

	if err := ExportVariablesToEnv(); err != nil {
		t.Fatalf("ExportVariablesToEnv failed: %v", err)
	}

	if os.Getenv("EXP_VAR_1") != "v1" {
		t.Errorf("expected EXP_VAR_1=v1, got %s", os.Getenv("EXP_VAR_1"))
	}
	if os.Getenv("TESTSUB_EXP_VAR_2") != "v2" {
		t.Errorf("expected TESTSUB_EXP_VAR_2=v2, got %s", os.Getenv("TESTSUB_EXP_VAR_2"))
	}
}
