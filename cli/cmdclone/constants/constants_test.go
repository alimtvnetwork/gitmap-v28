package constants

import "testing"

func TestConstants(t *testing.T) {
	if FlagMaxConcurrency != "max-concurrency" {
		t.Fatalf("unexpected flag: %s", FlagMaxConcurrency)
	}
	if MsgMultiCloneNoURLs == "" {
		t.Fatal("empty MsgMultiCloneNoURLs")
	}
	if FlagTargetDir != "target-dir" {
		t.Fatalf("unexpected flag: %s", FlagTargetDir)
	}
}
