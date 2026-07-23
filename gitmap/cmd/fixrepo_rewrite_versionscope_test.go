package cmd

// v5.39.0: regression matrix locking down the bare-base rewrite scope
// across current versions v1..v4. Bare `{base}` substitution must ONLY
// fire on the v1→v2 transition (current==2 with v1 in targets). Every
// other (current, targets) shape must preserve standalone `{base}`.

import "testing"

func TestApplyAllTargets_VersionScopeMatrix(t *testing.T) {
	// IMPORTANT: use a synthetic base ("testpkg") that does NOT match this
	// repo's own base name. If we used "gitmap" here, the fix-repo digit-
	// capture rewriter would smash literal `gitmap-v28` fixtures to the
	// current version on every bump (see mem: fix-repo digit-capture rule).
	const base = "testpkg"
	type tc struct {
		name    string
		current int
		targets []int
		in      string
		want    string
	}
	cases := []tc{
		{
			// current=1 is a no-op floor: nothing to bump to.
			name:    "v1_no_rewrite",
			current: 1,
			targets: []int{},
			in:      "testpkg and otherpkg-v9 stay put",
			want:    "testpkg and otherpkg-v9 stay put",
		},
		{
			// v1→v2: the ONLY case where bare base is rewritten.
			name:    "v2_bare_base_rewritten",
			current: 2,
			targets: []int{1},
			in:      "url=https://github.com/x/testpkg plus otherpkg-v9 token",
			want:    "url=https://github.com/x/testpkg-v2 plus otherpkg-v9 token",
		},
		{
			// v3: bare base preserved even with v1 in targets.
			name:    "v3_bare_base_preserved",
			current: 3,
			targets: []int{1, 2},
			in:      "testpkg binary and otherpkg-v9 and otherpkg-v9",
			want:    "testpkg binary and otherpkg-v9 and otherpkg-v9",
		},
		{
			// v4: bare base preserved across full target sweep.
			name:    "v4_bare_base_preserved",
			current: 4,
			targets: []int{1, 2, 3},
			in:      "testpkg binary, otherpkg-v9, otherpkg-v9, otherpkg-v9",
			want:    "testpkg binary, otherpkg-v9, otherpkg-v9, otherpkg-v9",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _ := applyAllTargets(c.in, base, c.current, c.targets)
			if got != c.want {
				t.Fatalf("scope mismatch (current=%d).\n got:  %q\n want: %q", c.current, got, c.want)
			}
		})
	}
}

