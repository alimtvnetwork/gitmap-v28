package secrets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveRepoSecretsRoot(t *testing.T) {
	root := ResolveRepoSecretsRoot()
	if root == "" {
		t.Fatalf("expected non-empty root")
	}
}

func TestResolveRepoSecretsNodesPath_DefaultAndToken(t *testing.T) {
	tmpDir := t.TempDir()
	secretsDir := filepath.Join(tmpDir, "repo-secrets")
	w1Dir := filepath.Join(secretsDir, "04-w1-machine")
	_ = os.MkdirAll(w1Dir, 0755)
	dummyNodes := filepath.Join(w1Dir, "gitmap-ssh-nodes.json")
	_ = os.WriteFile(dummyNodes, []byte("{}"), 0644)

	// Direct path
	resolvedDirect := ResolveRepoSecretsNodesPath(dummyNodes)
	if resolvedDirect != dummyNodes {
		t.Errorf("expected direct path %s, got %s", dummyNodes, resolvedDirect)
	}

	// Folder path token
	folderRes := resolveFolderByToken(secretsDir, "w1")
	if folderRes != w1Dir {
		t.Errorf("expected %s, got %s", w1Dir, folderRes)
	}

	nodesPath := findNodesFileInDir(w1Dir)
	if nodesPath != dummyNodes {
		t.Errorf("expected %s, got %s", dummyNodes, nodesPath)
	}
}

func TestResolveRepoSecretsManifest_Token(t *testing.T) {
	tmpDir := t.TempDir()
	secretsDir := filepath.Join(tmpDir, "repo-secrets")
	w2Dir := filepath.Join(secretsDir, "05-w2-machine")
	_ = os.MkdirAll(w2Dir, 0755)
	dummyManifest := filepath.Join(w2Dir, "gitmap.json")
	_ = os.WriteFile(dummyManifest, []byte("[]"), 0644)

	folderRes := resolveFolderByToken(secretsDir, "w2")
	if folderRes != w2Dir {
		t.Errorf("expected %s, got %s", w2Dir, folderRes)
	}

	manifest := findAnyManifestInDir(w2Dir)
	if manifest != dummyManifest {
		t.Errorf("expected %s, got %s", dummyManifest, manifest)
	}
}
