package cmd

import (
	"reflect"
	"testing"
)

func TestMergeStringPaths(t *testing.T) {
	tests := []struct {
		name     string
		existing []string
		incoming []string
		want     []string
	}{
		{
			name:     "empty lists",
			existing: []string{},
			incoming: []string{},
			want:     []string{},
		},
		{
			name:     "only incoming",
			existing: []string{},
			incoming: []string{"/a/b", "/c/d"},
			want:     []string{"/a/b", "/c/d"},
		},
		{
			name:     "only existing",
			existing: []string{"/a/b", "/c/d"},
			incoming: []string{},
			want:     []string{"/a/b", "/c/d"},
		},
		{
			name:     "no duplicates",
			existing: []string{"/a/b"},
			incoming: []string{"/c/d"},
			want:     []string{"/a/b", "/c/d"},
		},
		{
			name:     "with duplicates",
			existing: []string{"/a/b", "/c/d"},
			incoming: []string{"/c/d", "/e/f", "/a/b"},
			want:     []string{"/a/b", "/c/d", "/e/f"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeStringPaths(tt.existing, tt.incoming)
			// Handle empty slice comparison nil vs []
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeStringPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}
