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
