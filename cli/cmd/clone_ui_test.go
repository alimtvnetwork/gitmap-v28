package cmd

import "testing"

func TestInitCloneUIPart1(t *testing.T) {
	if err := InitCloneUIPart1(); err != nil {
		t.Fatal(err)
	}
}
