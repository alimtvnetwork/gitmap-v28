package cmd

import "testing"

func TestInitPullArrayUIPart1(t *testing.T) {
	if err := InitPullArrayUIPart1(); err != nil {
		t.Fatal(err)
	}
}
