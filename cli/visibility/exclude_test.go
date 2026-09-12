package visibility

import (
	"reflect"
	"testing"
)

func TestParseExclusionListNone(t *testing.T) {
	idx, isAll, err := ParseExclusionList("none", 10)
	if err != nil || isAll || len(idx) != 0 {
		t.Fatalf("expected none: %v %v %v", idx, isAll, err)
	}

	idx2, isAll2, err2 := ParseExclusionList("", 10)
	if err2 != nil || isAll2 || len(idx2) != 0 {
		t.Fatalf("expected empty: %v %v %v", idx2, isAll2, err2)
	}
}

func TestParseExclusionListAll(t *testing.T) {
	idx, isAll, err := ParseExclusionList("all", 10)
	if err != nil || !isAll || idx != nil {
		t.Fatalf("expected all: %v %v %v", idx, isAll, err)
	}
}

func TestParseExclusionListRanges(t *testing.T) {
	idx, isAll, err := ParseExclusionList("1, 3-5, 2", 10)
	if err != nil || isAll {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(idx, want) {
		t.Fatalf("got %v, want %v", idx, want)
	}
}

func TestParseExclusionListErrors(t *testing.T) {
	if _, _, err := ParseExclusionList("11", 10); err == nil {
		t.Fatal("expected error for index > totalCount")
	}

	if _, _, err := ParseExclusionList("0", 10); err == nil {
		t.Fatal("expected error for 0 index")
	}

	if _, _, err := ParseExclusionList("5-3", 10); err == nil {
		t.Fatal("expected error for backward range")
	}

	if _, _, err := ParseExclusionList("foo", 10); err == nil {
		t.Fatal("expected error for non-numeric token")
	}
}
