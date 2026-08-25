package cmd

import "testing"

func TestComp117(t *testing.T) {
	out, err := HandleComp117(Input117{ID: "test"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !out.Result {
		t.Fatal("expected Result to be true")
	}

	_, err = HandleComp117(Input117{ID: "error"})
	if err == nil {
		t.Fatal("expected err for ID='error'")
	}
}
