package cmd

import (
	"os"
	"reflect"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestStripVSCodeSyncDisabledFlag covers the three shapes the global
// kill switch can take on argv (long, short, absent) and verifies
// both the cleaned slice and the GITMAP_VSCODE_SYNC_DISABLED side
// effect. Each case restores the env var so cases stay independent.
func TestStripVSCodeSyncDisabledFlag(t *testing.T) {
	cases := []struct {
		name      string
		in        []string
		want      []string
		wantEnvOn bool
	}{
		{
			name:      "absent leaves args and env untouched",
			in:        []string{"clone", "https://example.com/r.git"},
			want:      []string{"clone", "https://example.com/r.git"},
			wantEnvOn: false,
		},
		{
			name:      "long form is stripped and env flipped",
			in:        []string{"--vscode-sync-disabled", "clone", "x"},
			want:      []string{"clone", "x"},
			wantEnvOn: true,
		},
		{
			name:      "short form mid-args is stripped and env flipped",
			in:        []string{"clone", "-vscode-sync-disabled", "x"},
			want:      []string{"clone", "x"},
			wantEnvOn: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			os.Unsetenv(constants.EnvVSCodeSyncDisabled)
			defer os.Unsetenv(constants.EnvVSCodeSyncDisabled)

			got := stripVSCodeSyncDisabledFlag(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("args = %#v, want %#v", got, tc.want)
			}

			gotOn := os.Getenv(constants.EnvVSCodeSyncDisabled) == constants.EnvVSCodeSyncDisabledOn
			if gotOn != tc.wantEnvOn {
				t.Errorf("env on = %v, want %v", gotOn, tc.wantEnvOn)
			}
		})
	}
}
