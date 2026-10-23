package cmdscan

import (
	"reflect"
	"testing"
)

func TestParseNumberedSelection_Range(t *testing.T) {
	indices, err := ParseNumberedSelection("1-3", 5)
	if err != nil {
		t.Fatalf("unexpected error parsing range 1-3: %v", err)
	}

	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(indices, expected) {
		t.Errorf("expected %v, got %v", expected, indices)
	}
}

func TestParseNumberedSelection_CommaList(t *testing.T) {
	indices, err := ParseNumberedSelection("2,4", 5)
	if err != nil {
		t.Fatalf("unexpected error parsing 2,4: %v", err)
	}

	expected := []int{2, 4}
	if !reflect.DeepEqual(indices, expected) {
		t.Errorf("expected %v, got %v", expected, indices)
	}
}

func TestParseNumberedSelection_Single(t *testing.T) {
	indices, err := ParseNumberedSelection("5", 5)
	if err != nil {
		t.Fatalf("unexpected error parsing 5: %v", err)
	}

	expected := []int{5}
	if !reflect.DeepEqual(indices, expected) {
		t.Errorf("expected %v, got %v", expected, indices)
	}
}
