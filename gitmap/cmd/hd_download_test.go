package cmd

import "testing"

func TestInitFunc009(t *testing.T) {
	if err := InitFunc009(); err != nil {
		t.Fatal(err)
	}
}

func TestInitFunc010(t *testing.T) {
	if err := InitFunc010(); err != nil {
		t.Fatal(err)
	}
}
