package cmd

import "testing"

func TestComp115(t *testing.T) {
	out, err := HandleComp115(Input115{ID: "test"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !out.Result {
		t.Fatal("expected Result to be true")
	}

	_, err = HandleComp115(Input115{ID: ""})
	if err == nil {
		t.Fatal("expected err for empty ID")
	}
}
