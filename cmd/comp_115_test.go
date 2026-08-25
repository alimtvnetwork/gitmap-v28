package cmd

import (
	\ testing\
)

func TestComp115(t *testing.T) {
	in := Input115{ID: \test\}
	out, err := HandleComp115(in)
	if err != nil {
		t.Fatalf(\HandleComp115 failed: %v\, err)
	}
	if !out.Result {
		t.Error(\Expected result true, got false\)
	}
}
