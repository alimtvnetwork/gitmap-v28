package cmd

import (
	"reflect"
	"testing"
)

type clusterFlagTestCase struct {
	name       string
	args       []string
	wantExcept string
	wantIPs    []string
	wantIDs    []int
	wantAuto   bool
	wantForce  bool
	wantNoPref bool
	wantPos    []string
	wantErr    bool
}

func clusterFilterCases() []clusterFlagTestCase {
	return []clusterFlagTestCase{
		{name: "except integers", args: []string{"--except", "2,4", "ps", "echo hello"}, wantExcept: "2,4", wantPos: []string{"ps", "echo hello"}},
		{name: "except mixed IP and int", args: []string{"--except", "192.168.1.24,151"}, wantExcept: "192.168.1.24,151", wantPos: []string{}},
		{name: "except range", args: []string{"--except", "2-5"}, wantExcept: "2-5", wantPos: []string{}},
		{name: "include IPs", args: []string{"--ip", "10.0.0.1,10.0.0.2", "--ip", "10.0.0.3"}, wantIPs: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}, wantPos: []string{}},
		{name: "include IDs", args: []string{"--id", "1,2", "--id", "3"}, wantIDs: []int{1, 2, 3}, wantPos: []string{}},
	}
}

func clusterControlAndErrorCases() []clusterFlagTestCase {
	return []clusterFlagTestCase{
		{name: "auto confirm short", args: []string{"-Y", "ps", "pwd"}, wantAuto: true, wantPos: []string{"ps", "pwd"}},
		{name: "auto confirm long", args: []string{"--yes"}, wantAuto: true, wantPos: []string{}},
		{name: "force lifecycle and no preflight", args: []string{"--force-lifecycle", "--no-preflight"}, wantForce: true, wantNoPref: true, wantPos: []string{}},
		{name: "mutual exclusivity error", args: []string{"--except", "2", "--ip", "10.0.0.1"}, wantErr: true},
		{name: "mutual exclusivity error with ID", args: []string{"--except", "2", "--id", "1"}, wantErr: true},
	}
}

func assertClusterScalars(t *testing.T, got ClusterFlags, tt clusterFlagTestCase) {
	if got.ExceptClause != tt.wantExcept {
		t.Errorf("ExceptClause = %v, want %v", got.ExceptClause, tt.wantExcept)
	}
	if got.AutoConfirm != tt.wantAuto || got.ForceLifecycle != tt.wantForce || got.NoPreflight != tt.wantNoPref {
		t.Errorf("flags mismatch: auto=%v/%v, force=%v/%v, noPref=%v/%v",
			got.AutoConfirm, tt.wantAuto, got.ForceLifecycle, tt.wantForce, got.NoPreflight, tt.wantNoPref)
	}
}

func assertClusterSlices(t *testing.T, got ClusterFlags, pos []string, tt clusterFlagTestCase) {
	if (len(got.OnlyIPs) > 0 || len(tt.wantIPs) > 0) && !reflect.DeepEqual(got.OnlyIPs, tt.wantIPs) {
		t.Errorf("OnlyIPs = %v, want %v", got.OnlyIPs, tt.wantIPs)
	}
	if (len(got.OnlyIDs) > 0 || len(tt.wantIDs) > 0) && !reflect.DeepEqual(got.OnlyIDs, tt.wantIDs) {
		t.Errorf("OnlyIDs = %v, want %v", got.OnlyIDs, tt.wantIDs)
	}
	if (len(pos) > 0 || len(tt.wantPos) > 0) && !reflect.DeepEqual(pos, tt.wantPos) {
		t.Errorf("pos = %v, want %v", pos, tt.wantPos)
	}
}

func runClusterFlagTest(t *testing.T, tt clusterFlagTestCase) {
	t.Run(tt.name, func(t *testing.T) {
		got, pos, err := ParseClusterFlags(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseClusterFlags() error = %v, wantErr %v", err, tt.wantErr)
			return nil
		}
		if !tt.wantErr {
			assertClusterScalars(t, got, tt)
			assertClusterSlices(t, got, pos, tt)
		}
	})
}

func TestParseClusterFlags(t *testing.T) {
	tests := append(clusterFilterCases(), clusterControlAndErrorCases()...)
	for _, tt := range tests {
		runClusterFlagTest(t, tt)
	}
}
