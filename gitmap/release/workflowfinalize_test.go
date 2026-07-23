package release

import (
	"reflect"
	"testing"
)

func TestCollectZipGroupNames_Empty(t *testing.T) {
	opts := Options{}
	got := collectZipGroupNames(opts)

	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestCollectZipGroupNames_PersistentOnly(t *testing.T) {
	opts := Options{ZipGroups: []string{"docs", "configs"}}
	got := collectZipGroupNames(opts)
	want := []string{"docs", "configs"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCollectZipGroupNames_BundleOnly(t *testing.T) {
	opts := Options{BundleName: "extras.zip"}
	got := collectZipGroupNames(opts)
	want := []string{"extras.zip"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCollectZipGroupNames_Merged(t *testing.T) {
	opts := Options{
		ZipGroups:  []string{"docs", "scripts"},
		BundleName: "bundle.zip",
	}
	got := collectZipGroupNames(opts)
	want := []string{"docs", "scripts", "bundle.zip"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCollectZipGroupNames_EmptyBundleName(t *testing.T) {
	opts := Options{
		ZipGroups:  []string{"assets"},
		BundleName: "",
	}
	got := collectZipGroupNames(opts)
	want := []string{"assets"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
