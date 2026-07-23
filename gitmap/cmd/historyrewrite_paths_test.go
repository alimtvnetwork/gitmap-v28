package cmd

import (
	"reflect"
	"testing"
)

// TestParseHistoryPathsNormalizesAllInputForms pins the multi-form
// path parser against the three shapes the spec promises to accept:
// separate args, comma-joined, and comma-space-joined. Order must be
// preserved and duplicates dropped.
func TestParseHistoryPathsNormalizesAllInputForms(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "separate args",
			args: []string{"a/x.go", "b/y.go", "c/z.go"},
			want: []string{"a/x.go", "b/y.go", "c/z.go"},
		},
		{
			name: "single comma-joined arg",
			args: []string{"a/x.go,b/y.go,c/z.go"},
			want: []string{"a/x.go", "b/y.go", "c/z.go"},
		},
		{
			name: "single comma-space-joined arg",
			args: []string{"a/x.go, b/y.go, c/z.go"},
			want: []string{"a/x.go", "b/y.go", "c/z.go"},
		},
		{
			name: "mixed comma and separate args",
			args: []string{"a/x.go,b/y.go", "c/z.go"},
			want: []string{"a/x.go", "b/y.go", "c/z.go"},
		},
		{
			name: "trailing comma and empty pieces dropped",
			args: []string{"a/x.go,", ",b/y.go,,", "c/z.go"},
			want: []string{"a/x.go", "b/y.go", "c/z.go"},
		},
		{
			name: "duplicates removed, order preserved",
			args: []string{"a/x.go", "b/y.go", "a/x.go", "c/z.go", "b/y.go"},
			want: []string{"a/x.go", "b/y.go", "c/z.go"},
		},
		{
			name: "folders and files coexist",
			args: []string{"docs/", "src/main.go,assets/"},
			want: []string{"docs/", "src/main.go", "assets/"},
		},
		{
			name: "empty input",
			args: []string{},
			want: []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseHistoryPaths(tc.args)
			if len(got) == 0 && len(tc.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("parseHistoryPaths(%#v) = %#v, want %#v", tc.args, got, tc.want)
			}
		})
	}
}
