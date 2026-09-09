package cmd

import (
	"bytes"
	"database/sql"
	"path/filepath"
	"testing"
)

func TestChromeTokensDoubleBase64Roundtrip(t *testing.T) {
	testPayloads := [][]byte{
		[]byte("1//0eXYZ12345_sample_refresh_token_string"),
		[]byte("v20\xcc]\x88\x12\x00\x00\x00binary\x00data\xff\xfe"),
		[]byte(""),
		[]byte("short"),
	}

	for _, original := range testPayloads {
		encoded := EncodeDoubleBase64(original)
		decoded, err := DecodeDoubleBase64(encoded)
		if err != nil {
			t.Fatalf("DecodeDoubleBase64 failed: %v", err)
		}
		if !bytes.Equal(original, decoded) {
			t.Errorf("expected %q, got %q", original, decoded)
		}
	}
}

func TestChromeTokensCaesarCipherRoundtrip(t *testing.T) {
	testStrings := []string{
		"MTIvMHhYWVoxMjM0NV9zYW1wbGVfcmVmcmVzaF90b2tlbl9zdHJpbmc=",
		"v20_Binary_Token_Data_1234567890",
		"Hello World 2026!",
	}

	shifts := []int{1, 5, 13, 25, 26, 42}
	for _, shift := range shifts {
		for _, s := range testStrings {
			encrypted := EncodeCaesarCipher(s, shift)
			decrypted := DecodeCaesarCipher(encrypted, shift)
			if s != decrypted {
				t.Errorf("shift %d: expected %s, got %s", shift, s, decrypted)
			}
		}
	}
}

func TestChromeTokensCaesarByteShiftRoundtrip(t *testing.T) {
	data := []byte("v20\xcc]\x88\x12\x00\x00\x00arbitrary_token_payload_with_high_bytes\xfe\xff")
	shifts := []byte{1, 7, 13, 255}

	for _, shift := range shifts {
		encoded := EncodeCaesarByteShift(data, shift)
		decoded, err := DecodeCaesarByteShift(encoded, shift)
		if err != nil {
			t.Fatalf("DecodeCaesarByteShift failed: %v", err)
		}
		if !bytes.Equal(data, decoded) {
			t.Errorf("shift %d: byte data mismatch", shift)
		}
	}
}

func TestChromeTokenVaultSQLiteRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "Web Data")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE token_service (service VARCHAR PRIMARY KEY NOT NULL, encrypted_token BLOB NOT NULL)")
	if err != nil {
		t.Fatalf("create table failed: %v", err)
	}

	tokenData := []byte("v20_sample_encrypted_refresh_token_payload")
	_, err = db.Exec("INSERT INTO token_service (service, encrypted_token) VALUES (?, ?)", "AccountId-118122973631983074723", tokenData)
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	db.Close()

	vault, err := readChromeTokenService(tmpDir)
	if err != nil {
		t.Fatalf("readChromeTokenService failed: %v", err)
	}
	if vault == nil {
		t.Fatalf("expected non-nil vault")
	}
	if vault.Count != 1 {
		t.Errorf("expected count 1, got %d", vault.Count)
	}

	entry := vault.Tokens[0]
	if entry.Service != "AccountId-118122973631983074723" {
		t.Errorf("unexpected service: %s", entry.Service)
	}
	if _, ok := entry.Variations["doubleBase64"]; !ok {
		t.Errorf("missing doubleBase64 variation")
	}
	if _, ok := entry.Variations["caesarCipher"]; !ok {
		t.Errorf("missing caesarCipher variation")
	}
	if _, ok := entry.Variations["caesarByteShift"]; !ok {
		t.Errorf("missing caesarByteShift variation")
	}

	// Test restoring into a new destination profile directory
	dstDir := t.TempDir()
	restoreErr := restoreChromeTokenService(dstDir, vault)
	if restoreErr != nil {
		t.Fatalf("restoreChromeTokenService failed: %v", restoreErr)
	}

	dstDB, err := sql.Open("sqlite", filepath.Join(dstDir, "Web Data"))
	if err != nil {
		t.Fatalf("open dst db failed: %v", err)
	}
	defer dstDB.Close()

	var restoredService string
	var restoredBytes []byte
	err = dstDB.QueryRow("SELECT service, encrypted_token FROM token_service WHERE service=?", entry.Service).Scan(&restoredService, &restoredBytes)
	if err != nil {
		t.Fatalf("query restored token failed: %v", err)
	}
	if !bytes.Equal(tokenData, restoredBytes) {
		t.Errorf("restored token mismatch: expected %q, got %q", tokenData, restoredBytes)
	}
}

func TestBuildExportFromDiskIncludesTokenVault(t *testing.T) {
	tmpDir := t.TempDir()
	initDummyTokenWebData(t, tmpDir)
	exp := buildExportFromDisk("Profile 1", tmpDir)
	if exp.TokenVault == nil {
		t.Fatalf("expected TokenVault to be populated")
	}
	if exp.TokenVault.Count != 1 {
		t.Errorf("expected 1 token in vault, got %d", exp.TokenVault.Count)
	}
}

func initDummyTokenWebData(t *testing.T, dir string) {
	dbPath := filepath.Join(dir, "Web Data")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()
	_, _ = db.Exec("CREATE TABLE token_service (service VARCHAR PRIMARY KEY NOT NULL, encrypted_token BLOB NOT NULL)")
	_, _ = db.Exec("INSERT INTO token_service (service, encrypted_token) VALUES (?, ?)", "service-1", []byte("secret-token"))
}
