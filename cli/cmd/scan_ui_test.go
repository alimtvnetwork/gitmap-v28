package cmd

import "testing"

func TestInitScanUIPart1(t *testing.T) {
	if err := InitScanUIPart1(); err != nil {
		t.Fatal(err)
	}
}
