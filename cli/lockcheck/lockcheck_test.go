package lockcheck

import (
	"testing"
)

func TestIsProtectedProcess_IDE(t *testing.T) {
	cases := []struct {
		name string
		pid  int
		want bool
	}{
		{"Antigravity.exe", 1234, true},
		{"antigravity", 5678, true},
		{"code.exe", 1111, true},
		{"electron.exe", 2222, true},
		{"msedgewebview2.exe", 3333, true},
		{"node.exe", 4444, true},
		{"python.exe", 5555, true},
		{"pwsh.exe", 6666, true},
		{"bash", 7777, true},
		{"random_app.exe", 8888, false},
		{"git.exe", 9999, false},
		{"unknown", 0, true},
		{"unknown", -1, true},
		{"unknown", 4, true},
	}

	for _, tc := range cases {
		got := IsProtectedProcess(tc.name, tc.pid)
		if got != tc.want {
			t.Errorf("IsProtectedProcess(%q, %d) = %v; want %v", tc.name, tc.pid, got, tc.want)
		}
	}
}

func TestKillProcess_RefusesProtected(t *testing.T) {
	err := KillProcess(4)
	if err == nil {
		t.Fatal("expected KillProcess(4) to return error refusing to kill System PID")
	}
}

func TestFormatProcessList(t *testing.T) {
	procs := []LockingProcess{
		{PID: 100, Name: "worker"},
		{PID: 200, Name: "daemon"},
	}
	out := FormatProcessList(procs)
	if out == "" {
		t.Fatal("expected non-empty output")
	}

	emptyOut := FormatProcessList(nil)
	if emptyOut != "" {
		t.Fatalf("expected empty string for nil list, got %q", emptyOut)
	}
}
