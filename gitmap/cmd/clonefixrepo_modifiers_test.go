package cmd

import (
	"reflect"
	"testing"
)

// TestParseCfrModifiers pins the order-independent contract for the `cg`
// and `p` pre-URL modifier tokens accepted by `cfr` / `cfrp`. Any change
// that alters recognition, ordering, or the passthrough boundary must be
// deliberate and update this table.
func TestParseCfrModifiers(t *testing.T) {
	t.Parallel()

	type want struct {
		cg   bool
		pub  bool
		rest []string
	}
	cases := []struct {
		name string
		in   []string
		want want
	}{
		{
			name: "no_modifiers_passthrough",
			in:   []string{"https://x/y"},
			want: want{rest: []string{"https://x/y"}},
		},
		{
			name: "cg_only",
			in:   []string{"cg", "https://x/y"},
			want: want{cg: true, rest: []string{"https://x/y"}},
		},
		{
			name: "p_only",
			in:   []string{"p", "https://x/y"},
			want: want{pub: true, rest: []string{"https://x/y"}},
		},
		{
			name: "cg_then_p",
			in:   []string{"cg", "p", "https://x/y"},
			want: want{cg: true, pub: true, rest: []string{"https://x/y"}},
		},
		{
			name: "p_then_cg_order_independent",
			in:   []string{"p", "cg", "https://x/y"},
			want: want{cg: true, pub: true, rest: []string{"https://x/y"}},
		},
		{
			name: "duplicate_tokens_idempotent",
			in:   []string{"cg", "cg", "p", "https://x/y"},
			want: want{cg: true, pub: true, rest: []string{"https://x/y"}},
		},
		{
			name: "flag_stops_scan",
			in:   []string{"--ssh", "cg", "https://x/y"},
			want: want{rest: []string{"--ssh", "cg", "https://x/y"}},
		},
		{
			name: "url_stops_scan_trailing_tokens_ignored",
			in:   []string{"https://x/y", "cg"},
			want: want{rest: []string{"https://x/y", "cg"}},
		},
		{
			name: "unknown_token_stops_scan",
			in:   []string{"xx", "cg", "https://x/y"},
			want: want{rest: []string{"xx", "cg", "https://x/y"}},
		},
		{
			name: "empty_argv",
			in:   nil,
			want: want{},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, rest := ParseCfrModifiers(tc.in)
			if got.InstallCodingGuidelines != tc.want.cg || got.PromotePublic != tc.want.pub {
				t.Fatalf("flags: got {cg=%v p=%v}, want {cg=%v p=%v}",
					got.InstallCodingGuidelines, got.PromotePublic, tc.want.cg, tc.want.pub)
			}
			if !reflect.DeepEqual(rest, tc.want.rest) && !(len(rest) == 0 && len(tc.want.rest) == 0) {
				t.Fatalf("rest: got %#v, want %#v", rest, tc.want.rest)
			}
		})
	}
}
