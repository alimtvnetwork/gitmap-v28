package cmdssh

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestParseInstallTargetAndPackageZeroArgs(t *testing.T) {
	pkg, target := parseInstallTargetAndPackage(nil, []string{})
	if pkg != "gitmap" || target != "all" {
		t.Fatalf("expected gitmap/all, got %s/%s", pkg, target)
	}
}

func TestParseInstallTargetAndPackageSingleArgGitmap(t *testing.T) {
	pkg, target := parseInstallTargetAndPackage(nil, []string{"gitmap"})
	if pkg != "gitmap" || target != "all" {
		t.Fatalf("expected gitmap/all, got %s/%s", pkg, target)
	}
}

func TestParseInstallTargetAndPackageSingleArgTarget(t *testing.T) {
	conns := []db.SSHConnection{{Alias: "node1", IPAddress: "10.0.0.1"}}
	pkg, target := parseInstallTargetAndPackage(conns, []string{"node1"})
	if pkg != "gitmap" || target != "node1" {
		t.Fatalf("expected gitmap/node1, got %s/%s", pkg, target)
	}
}

func TestParseInstallTargetAndPackageSingleArgPackage(t *testing.T) {
	conns := []db.SSHConnection{{Alias: "node1", IPAddress: "10.0.0.1"}}
	pkg, target := parseInstallTargetAndPackage(conns, []string{"devbox"})
	if pkg != "devbox" || target != "all" {
		t.Fatalf("expected devbox/all, got %s/%s", pkg, target)
	}
}

func TestParseInstallTargetAndPackageTwoArgs(t *testing.T) {
	conns := []db.SSHConnection{{Alias: "node1", IPAddress: "10.0.0.1"}}
	pkg, target := parseInstallTargetAndPackage(conns, []string{"agy", "node1"})
	if pkg != "agy" || target != "node1" {
		t.Fatalf("expected agy/node1, got %s/%s", pkg, target)
	}
}

func TestResolveRemoteShell(t *testing.T) {
	if resolveRemoteShell("windows") != "ps" {
		t.Fatal("expected ps for windows")
	}
	if resolveRemoteShell("linux") != "bash" {
		t.Fatal("expected bash for linux")
	}
}

func TestParseInstallTargetAndPackage_SingleAllFlag(t *testing.T) {
	pkg1, target1 := parseInstallTargetAndPackage(nil, []string{"--all"})
	if pkg1 != "gitmap" || target1 != "all" {
		t.Fatalf("expected gitmap/all, got %s/%s", pkg1, target1)
	}
	pkg2, target2 := parseInstallTargetAndPackage(nil, []string{"-a"})
	if pkg2 != "gitmap" || target2 != "all" {
		t.Fatalf("expected gitmap/all, got %s/%s", pkg2, target2)
	}
}

func TestParseInstallTargetAndPackage_TwoArgsAllFlag(t *testing.T) {
	pkg1, target1 := parseInstallTargetAndPackage(nil, []string{"agy", "--all"})
	if pkg1 != "agy" || target1 != "all" {
		t.Fatalf("expected agy/all, got %s/%s", pkg1, target1)
	}
	pkg2, target2 := parseInstallTargetAndPackage(nil, []string{"--all", "agy"})
	if pkg2 != "agy" || target2 != "all" {
		t.Fatalf("expected agy/all, got %s/%s", pkg2, target2)
	}
}
