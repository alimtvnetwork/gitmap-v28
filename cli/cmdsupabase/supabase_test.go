package cmdsupabase

import (
	"database/sql"
	"os"
	"testing"
)

func TestMaskSecret(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "none"},
		{"short", "****"},
		{"12345678", "****"},
		{"1234567890", "1234...****"},
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", "eyJhbG...****"},
	}

	for _, tc := range tests {
		actual := MaskSecret(tc.input)
		if actual != tc.expected {
			t.Errorf("MaskSecret(%q) = %q; want %q", tc.input, actual, tc.expected)
		}
	}
}

func TestEncryptDecryptSecret_Roundtrip(t *testing.T) {
	original := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.super-secret-service-role-key"
	encrypted, err := EncryptSecret(original)
	if err != nil {
		t.Fatalf("EncryptSecret failed: %v", err)
	}

	if encrypted == original {
		t.Fatalf("Ciphertext must not match original plaintext!")
	}

	decrypted, err := DecryptSecret(encrypted)
	if err != nil {
		t.Fatalf("DecryptSecret failed: %v", err)
	}

	if decrypted != original {
		t.Fatalf("Decrypted = %q; want %q", decrypted, original)
	}
}

func TestSupabaseDb_CRUD(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer db.Close()

	if err := initSupabaseSchema(db); err != nil {
		t.Fatalf("failed to init schema: %v", err)
	}

	rec := SupabaseProjectRecord{
		Alias:         "test-proj",
		ProjectRef:    "test-ref-123",
		ApiUrl:        "https://test.supabase.co",
		AnonKeyEnc:    "enc_anon_key",
		ServiceKeyEnc: "enc_service_key",
		DbUrlEnc:      "enc_db_url",
		Status:        StatusActive,
		Description:   "Test project",
	}

	if err := executeInsertProject(db, rec); err != nil {
		t.Fatalf("executeInsertProject failed: %v", err)
	}

	fetched, err := queryProjectByAlias(db, "test-proj")
	if err != nil {
		t.Fatalf("queryProjectByAlias failed: %v", err)
	}

	if fetched.Alias != "test-proj" || fetched.ApiUrl != "https://test.supabase.co" {
		t.Fatalf("fetched record mismatch: %+v", fetched)
	}

	all, err := queryAllProjects(db)
	if err != nil {
		t.Fatalf("queryAllProjects failed: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 project, got %d", len(all))
	}
}

func TestIsValidSupabaseUrl(t *testing.T) {
	if !isValidSupabaseUrl("https://xyz.supabase.co") {
		t.Errorf("expected https url to be valid")
	}
	if !isValidSupabaseUrl("http://localhost:54321") {
		t.Errorf("expected http url to be valid")
	}
	if isValidSupabaseUrl("ftp://localhost") {
		t.Errorf("expected ftp url to be invalid")
	}
}

func TestParseAddFlags(t *testing.T) {
	args := []string{
		"my-proj", "https://xyz.supabase.co", "anon", "service", "postgres://db",
		"--desc", "My Description", "--project-ref", "ref123",
	}

	pos, flags := parseAddFlags(args)
	if len(pos) != 5 {
		t.Fatalf("expected 5 positional args, got %d", len(pos))
	}
	if flags["desc"] != "My Description" {
		t.Errorf("expected desc 'My Description', got %q", flags["desc"])
	}
	if flags["project-ref"] != "ref123" {
		t.Errorf("expected project-ref 'ref123', got %q", flags["project-ref"])
	}
}

func TestRunSupabaseCLI_Help(t *testing.T) {
	if err := RunSupabaseCLI([]string{"help"}); err != nil {
		t.Errorf("expected nil error for help, got: %v", err)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
