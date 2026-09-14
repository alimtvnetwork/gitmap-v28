package cmdssh

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	_ "modernc.org/sqlite"
)

func writeTempFile(t *testing.T, filename, content string) string {
	path := filepath.Join(t.TempDir(), filename)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func sampleClusterConfigContent() string {
	return `{
  "control": {
    "master": "192.168.0.20"
  },
  "nodes": {
    "worker-1": "192.168.0.21",
    "worker-2": "192.168.0.22"
  },
  "user": {
    "name": "kubeadmin",
    "password": "SecretPassword123!"
  }
}`
}

func TestLoadClusterConfigFile_Valid(t *testing.T) {
	path := writeTempFile(t, "valid-config.json", sampleClusterConfigContent())
	cfg, err := LoadClusterConfigFile(path)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if cfg.User.Name != "kubeadmin" || cfg.User.Password != "SecretPassword123!" {
		t.Fatalf("user mismatch: %+v", cfg.User)
	}
	if cfg.Control["master"] != "192.168.0.20" || len(cfg.Nodes) != 2 {
		t.Fatalf("topology mismatch: %+v", cfg)
	}
}

func TestLoadClusterConfigFile_Errors(t *testing.T) {
	_, errNotFound := LoadClusterConfigFile(filepath.Join(t.TempDir(), "missing.json"))
	if errNotFound == nil {
		t.Fatal("expected not found error")
	}
	invalidPath := writeTempFile(t, "bad.json", "{invalid-json")
	_, errInvalid := LoadClusterConfigFile(invalidPath)
	if errInvalid == nil {
		t.Fatal("expected json parse error")
	}
}

func TestValidateClusterUser(t *testing.T) {
	validUser := ClusterUserJSON{Name: "admin", Password: "pwd"}
	if err := ValidateClusterUser(validUser); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	emptyName := ClusterUserJSON{Name: "", Password: "pwd"}
	if err := ValidateClusterUser(emptyName); err == nil {
		t.Fatal("expected error for empty name")
	}
	emptyPass := ClusterUserJSON{Name: "admin", Password: ""}
	if err := ValidateClusterUser(emptyPass); err == nil {
		t.Fatal("expected error for empty password")
	}
}

func generateTestRSAPrivateKeyFile(t *testing.T) string {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate rsa key: %v", err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})
	return writeTempFile(t, "id_rsa", string(pemBytes))
}

func withMockSSHRSAKey(t *testing.T, fn func(keyPath string)) {
	keyPath := generateTestRSAPrivateKeyFile(t)
	prevLocator := defaultSSHKeyLocator
	defaultSSHKeyLocator = func() (string, bool) { return keyPath, true }
	defer func() { defaultSSHKeyLocator = prevLocator }()
	fn(keyPath)
}

func assertRSACiphertext(t *testing.T, encPass string) {
	isRSA := isRSACiphertext(encPass)
	if !isRSA {
		t.Fatalf("expected rsa ciphertext prefix, got: %s", encPass)
	}
	status := formatEncryptionStatus(encPass)
	if status != "rsa:sha256 [encrypted]" {
		t.Fatalf("unexpected status: %s", status)
	}
}

func TestClusterImport_RSAEncryptionValidation(t *testing.T) {
	withMockSSHRSAKey(t, func(keyPath string) {
		encPass, err := EncryptSSHPassword("KubeSecret@2026")
		if err != nil {
			t.Fatalf("EncryptSSHPassword failed: %v", err)
		}
		assertRSACiphertext(t, encPass)
	})
}

func setupClusterTestDB(t *testing.T) *store.DB {
	dbPath := filepath.Join(t.TempDir(), "cluster_test.db")
	dbConn, err := store.OpenAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open cluster test db: %v", err)
	}
	_ = dbConn.Migrate()
	_ = store.EnsureSSHHostsTable(dbConn.SQL())
	_ = store.EnsureSSHHistoryTable(dbConn.SQL())
	return dbConn
}

func withClusterTestContext(t *testing.T, fn func(db *store.DB)) {
	testDB := setupClusterTestDB(t)
	defer testDB.Close()
	prevOpener := openSSHDBFunc
	openSSHDBFunc = func() (*store.DB, error) { return testDB, nil }
	defer func() { openSSHDBFunc = prevOpener }()
	fn(testDB)
}

func assertHostEnrolled(t *testing.T, db *sql.DB, alias, ip, role string) {
	var user, clusterRole, encPass string
	var port int
	query := "SELECT username, port, encrypted_password, cluster_role FROM ssh_hosts WHERE alias = ?"
	err := db.QueryRow(query, alias).Scan(&user, &port, &encPass, &clusterRole)
	if err != nil {
		t.Fatalf("query host failed for alias %s: %v", alias, err)
	}
	if user != "kubeadmin" || port != 22 || clusterRole != role || encPass == "" {
		t.Fatalf("host data mismatch for %s: user=%s, port=%d, role=%s", alias, user, port, clusterRole)
	}
}

func assertTableOutput(t *testing.T, output string) {
	tokens := []string{"ROLE", "ALIAS", "IP", "USERNAME", "STATUS", "control", "master", "worker-1", "worker-2", "rsa:sha256 [encrypted]"}
	for _, token := range tokens {
		isContained := strings.Contains(output, token)
		if !isContained {
			t.Fatalf("missing expected token %q in output:\n%s", token, output)
		}
	}
}

func captureClusterImportOutput(t *testing.T, path string) string {
	var buf bytes.Buffer
	prevOut := clusterImportOut
	clusterImportOut = &buf
	defer func() { clusterImportOut = prevOut }()
	if err := RunClusterImportCLI([]string{path}); err != nil {
		t.Fatalf("RunClusterImportCLI failed: %v", err)
	}
	return buf.String()
}

func TestRunClusterImportCLI_Success(t *testing.T) {
	withClusterTestContext(t, func(testDB *store.DB) {
		withMockSSHRSAKey(t, func(keyPath string) {
			path := writeTempFile(t, "01-config.json", sampleClusterConfigContent())
			out := captureClusterImportOutput(t, path)
			assertHostEnrolled(t, testDB.SQL(), "master", "192.168.0.20", "control")
			assertHostEnrolled(t, testDB.SQL(), "worker-1", "192.168.0.21", "worker")
			assertHostEnrolled(t, testDB.SQL(), "worker-2", "192.168.0.22", "worker")
			assertTableOutput(t, out)
		})
	})
}

func TestSJClusterImportSubcommandAliases(t *testing.T) {
	aliases := []string{"import-cluster", "import", "import-config"}
	for _, alias := range aliases {
		isMatch := isSJClusterImportSubcommand(alias)
		if !isMatch {
			t.Fatalf("expected alias %s to match", alias)
		}
		isSub := isSJSubcommand(alias)
		if !isSub {
			t.Fatalf("expected alias %s to be recognized as sj subcommand", alias)
		}
	}
}

func TestResolveImportConfigPath(t *testing.T) {
	defaultPath := resolveImportConfigPath([]string{})
	if defaultPath != "./01-config.json" {
		t.Fatalf("expected ./01-config.json, got %s", defaultPath)
	}
	customPath := resolveImportConfigPath([]string{"custom.json"})
	if customPath != "custom.json" {
		t.Fatalf("expected custom.json, got %s", customPath)
	}
}
