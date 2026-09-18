package cmdssh

import (
	"context"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestFindMultiTargetIndex(t *testing.T) {
	cases := []struct {
		args []string
		want int
	}{
		{[]string{"192.168.1.10,192.168.1.11"}, 0},
		{[]string{"add", "192.168.1.10,192.168.1.11"}, 1},
		{[]string{"--user", "root", "10.0.0.1,10.0.0.2"}, 2},
		{[]string{"192.168.1.10", "single-box"}, -1},
		{[]string{"--help"}, -1},
		{[]string{}, -1},
	}

	for _, tc := range cases {
		got := findMultiTargetIndex(tc.args)
		if got != tc.want {
			t.Errorf("findMultiTargetIndex(%v) = %d, want %d", tc.args, got, tc.want)
		}
	}
}

func TestBuildSingleTargetArgs(t *testing.T) {
	args := []string{"--user", "ubuntu", "10.0.0.1,10.0.0.2", "--auth"}
	single := buildSingleTargetArgs(args, 2, "10.0.0.1")

	if single[2] != "10.0.0.1" || len(single) != 4 {
		t.Errorf("unexpected single target args: %v", single)
	}
	if args[2] != "10.0.0.1,10.0.0.2" {
		t.Error("original args slice was mutated")
	}
}

func TestRunSSHJoinCLI_MultiEnrollment(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		args := []string{"10.0.20.1,10.0.20.2"}
		if err := runSSHJoinCLI(args); err != nil {
			t.Fatalf("runSSHJoinCLI with comma targets failed: %v", err)
		}

		ctx := context.Background()
		h1, err1 := store.GetHostByAlias(ctx, "host-10.0.20.1", db.SQL())
		if err1 != nil || h1.IP != "10.0.20.1" {
			t.Errorf("failed to find first enrolled host: %v, %+v", err1, h1)
		}

		h2, err2 := store.GetHostByAlias(ctx, "host-10.0.20.2", db.SQL())
		if err2 != nil || h2.IP != "10.0.20.2" {
			t.Errorf("failed to find second enrolled host: %v, %+v", err2, h2)
		}
	})
}

func TestIsTargetAddress(t *testing.T) {
	if !isTargetAddress("192.168.1.1") || !isTargetAddress("root@10.0.0.1") || !isTargetAddress("10.0.0.1:2222") {
		t.Error("expected valid IP addresses to return true")
	}
	if isTargetAddress("devbox") || isTargetAddress("my-alias") || isTargetAddress("") {
		t.Error("expected aliases and empty string to return false")
	}
}

func TestExtractTargetPositionsAndFlags(t *testing.T) {
	args := []string{"10.0.30.1", "10.0.30.2", "--auth", "-u", "ubuntu"}
	targets := extractTargetPositions(args)
	flags := extractFlagArgs(args)

	if len(targets) != 2 || targets[0] != "10.0.30.1" || targets[1] != "10.0.30.2" {
		t.Errorf("unexpected extracted targets: %v", targets)
	}
	if len(flags) != 2 || flags[0] != "--auth" || flags[1] != "-u" {
		t.Errorf("unexpected extracted flags: %v", flags)
	}
}

func TestRunSSHJoinCLI_SpaceMultiEnrollment(t *testing.T) {
	withMockSSHDB(t, func(db *store.DB) {
		args := []string{"10.0.40.1", "10.0.40.2"}
		if err := runSSHJoinCLI(args); err != nil {
			t.Fatalf("runSSHJoinCLI with space targets failed: %v", err)
		}

		ctx := context.Background()
		h1, err1 := store.GetHostByAlias(ctx, "host-10.0.40.1", db.SQL())
		if err1 != nil || h1.IP != "10.0.40.1" {
			t.Errorf("failed to find first enrolled host: %v, %+v", err1, h1)
		}

		h2, err2 := store.GetHostByAlias(ctx, "host-10.0.40.2", db.SQL())
		if err2 != nil || h2.IP != "10.0.40.2" {
			t.Errorf("failed to find second enrolled host: %v, %+v", err2, h2)
		}
	})
}

