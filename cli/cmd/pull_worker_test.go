package cmd

import "testing"

func TestInitFunc007(t *testing.T) {
	if err := InitFunc007(); err != nil {
		t.Fatal(err)
	}
}

func TestInitFunc008(t *testing.T) {
	if err := InitFunc008(); err != nil {
		t.Fatal(err)
	}
}
