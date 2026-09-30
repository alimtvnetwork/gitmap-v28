package cmdssh

import (
	"testing"
)

func TestDeduplicatePublicKeys_FilterDuplicatesAndComments(t *testing.T) {
	rawKeys := []string{
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIG12345 user@host1",
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIG12345 user@host2", // duplicate key
		"# this is a comment",
		"",
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC98765 user@host3",
		"invalid-key-format",
	}

	unique := deduplicatePublicKeys(rawKeys)
	if len(unique) != 2 {
		t.Fatalf("expected 2 unique keys, got %d: %+v", len(unique), unique)
	}
	if unique[0] != "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIG12345 user@host1" {
		t.Errorf("unexpected first key: %s", unique[0])
	}
	if unique[1] != "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC98765 user@host3" {
		t.Errorf("unexpected second key: %s", unique[1])
	}
}

func TestFilterMissingPublicKeys(t *testing.T) {
	existingAuth := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIG12345 user@host1\n"
	existingMap := buildExistingSignaturesMap(existingAuth)

	keys := []string{
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIG12345 user@host1",
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIG99999 user@host2",
	}

	missing := filterMissingPublicKeys(keys, existingMap)
	if len(missing) != 1 {
		t.Fatalf("expected 1 missing key, got %d", len(missing))
	}
	if missing[0] != keys[1] {
		t.Errorf("expected missing key %s, got %s", keys[1], missing[0])
	}
}

func TestParseDeployKeysFlags(t *testing.T) {
	args := []string{"all", "--except", "worker-1,node-2", "--dry-run", "--json"}
	except, isDryRun, isJSON := parseDeployKeysFlags(args)

	if except != "worker-1,node-2" {
		t.Errorf("unexpected except: %s", except)
	}
	if !isDryRun {
		t.Errorf("expected isDryRun true")
	}
	if !isJSON {
		t.Errorf("expected isJSON true")
	}
}

func TestParseDeployKeysFlags_DirectVariants(t *testing.T) {
	testCases := [][]string{
		{"deploy-keys-all", "--except=worker-1"},
		{"keys-all", "--dry-run"},
		{"--all", "-n"},
		{"deploy-keys", "all"},
		{"all-keys", "--dry-run"},
		{"all", "keys", "-n"},
		{"keys", "all", "-n"},
		{"key", "all", "-n"},
	}

	for _, tc := range testCases {
		except, isDryRun, _ := parseDeployKeysFlags(tc)
		if len(tc) > 1 && tc[1] == "--except=worker-1" && except != "worker-1" {
			t.Errorf("expected except worker-1, got %s", except)
		}
		if len(tc) > 1 && (tc[len(tc)-1] == "--dry-run" || tc[len(tc)-1] == "-n") && !isDryRun {
			t.Errorf("expected isDryRun true for args %v", tc)
		}
	}
}

func TestIsIgnoredDeployKeysToken(t *testing.T) {
	validTokens := []string{
		"keys", "key", "k", "deploy", "ssh", "all", "--all", "-all",
		"keys-all", "deploy-keys", "deploy-keys-all", "keys-hyphen-all", "hyphen-all",
		"all-keys", "all-key", "allkeys",
	}
	for _, tok := range validTokens {
		if !isIgnoredDeployKeysToken(tok) {
			t.Errorf("expected %q to be recognized as ignored token", tok)
		}
	}
	if isIgnoredDeployKeysToken("worker-1") {
		t.Errorf("expected node identifier not to be ignored")
	}
}

func TestIsKnownLegacyDeployWord_DeployKeysVariants(t *testing.T) {
	variants := []string{"all-keys", "all-key", "allkeys", "all", "keys", "key", "keys-all", "deploy-keys", "deploy-keys-all"}
	for _, v := range variants {
		if !isKnownLegacyDeployWord(v) {
			t.Errorf("expected isKnownLegacyDeployWord(%q) to be true", v)
		}
	}
}
