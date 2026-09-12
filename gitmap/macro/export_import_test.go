package macro

import (
	"path/filepath"
	"testing"
	"time"
)

func sampleMacroFixture(name string) Macro {
	return Macro{
		Name:        name,
		Description: "test macro description for " + name,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		TotalSteps:  2,
		Tags:        "test,ci",
		Steps: []MacroStep{
			{StepNum: 1, CommandLine: "echo 1", WorkingDir: "/tmp", ContinueOnError: true, TimeoutSeconds: 10},
			{StepNum: 2, CommandLine: "echo 2", WorkingDir: "/tmp", ContinueOnError: false, TimeoutSeconds: 20},
		},
	}
}

func TestExportMacros_SingleAndAllJSON(t *testing.T) {
	m1 := sampleMacroFixture("json-macro-1")
	m2 := sampleMacroFixture("json-macro-2")
	payload, err := ExportToJSON([]Macro{m1, m2})
	if err != nil {
		t.Fatalf("ExportToJSON error: %v", err)
	}

	parsedList, err := ParseImportJSON(payload)
	if err != nil || len(parsedList) != 2 {
		t.Fatalf("ParseImportJSON slice failed: count=%d, err=%v", len(parsedList), err)
	}

	singlePayload, err := ExportToJSON(m1)
	if err != nil {
		t.Fatalf("ExportToJSON single error: %v", err)
	}

	singleParsed, err := ParseImportJSON(singlePayload)
	if err != nil || len(singleParsed) != 1 {
		t.Fatalf("ParseImportJSON single failed: count=%d, err=%v", len(singleParsed), err)
	}
}

func TestExportMacros_SingleAndAllYAML(t *testing.T) {
	m1 := sampleMacroFixture("yaml-macro-1")
	m2 := sampleMacroFixture("yaml-macro-2")
	payload, err := ExportToYAML([]Macro{m1, m2})
	if err != nil {
		t.Fatalf("ExportToYAML error: %v", err)
	}

	parsedList, err := ParseImportYAML(payload)
	if err != nil || len(parsedList) != 2 {
		t.Fatalf("ParseImportYAML slice failed: count=%d, err=%v", len(parsedList), err)
	}

	singlePayload, err := ExportToYAML(m1)
	if err != nil {
		t.Fatalf("ExportToYAML single error: %v", err)
	}

	singleParsed, err := ParseImportYAML(singlePayload)
	if err != nil || len(singleParsed) != 1 {
		t.Fatalf("ParseImportYAML single failed: count=%d, err=%v", len(singleParsed), err)
	}
}

func TestExportMacros_SQLiteDatabaseRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_macros.db")
	m1 := sampleMacroFixture("sql-macro-1")
	m2 := sampleMacroFixture("sql-macro-2")
	if err := ExportMacrosToSQLite([]Macro{m1, m2}, dbPath); err != nil {
		t.Fatalf("ExportMacrosToSQLite error: %v", err)
	}

	imported, err := ParseImportSQLite(dbPath)
	if err != nil {
		t.Fatalf("ParseImportSQLite error: %v", err)
	}

	if len(imported) != 2 {
		t.Fatalf("expected 2 imported macros, got %d", len(imported))
	}

	if len(imported[0].Steps) != 2 || imported[0].Steps[0].CommandLine != "echo 1" {
		t.Fatalf("imported macro steps mismatch: %+v", imported[0].Steps)
	}
}

func TestExportMacros_ZIPArchiveRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	zipPath := filepath.Join(tmpDir, "test_macros.zip")
	m1 := sampleMacroFixture("zip-macro-1")
	m2 := sampleMacroFixture("zip-macro-2")
	if err := ExportToZIP([]Macro{m1, m2}, zipPath); err != nil {
		t.Fatalf("ExportToZIP error: %v", err)
	}

	imported, err := ParseImportZIP(zipPath)
	if err != nil {
		t.Fatalf("ParseImportZIP error: %v", err)
	}

	if len(imported) != 2 {
		t.Fatalf("expected 2 imported macros from zip, got %d", len(imported))
	}
}

func TestValidateMacro_RejectsInvalidNameAndEmptySteps(t *testing.T) {
	emptyName := Macro{Name: "", Steps: []MacroStep{{CommandLine: "echo"}}}
	if err := ValidateMacro(&emptyName); err == nil {
		t.Fatal("expected error for empty macro name")
	}

	traversalName := Macro{Name: "../badname", Steps: []MacroStep{{CommandLine: "echo"}}}
	if err := ValidateMacro(&traversalName); err == nil {
		t.Fatal("expected error for path traversal name")
	}

	slashName := Macro{Name: "sub/folder", Steps: []MacroStep{{CommandLine: "echo"}}}
	if err := ValidateMacro(&slashName); err == nil {
		t.Fatal("expected error for slash in name")
	}

	noSteps := Macro{Name: "validname", Steps: nil}
	if err := ValidateMacro(&noSteps); err == nil {
		t.Fatal("expected error for macro with zero steps")
	}

	emptyCommand := Macro{Name: "validname", Steps: []MacroStep{{CommandLine: "   "}}}
	if err := ValidateMacro(&emptyCommand); err == nil {
		t.Fatal("expected error for macro step with empty command")
	}
}

func TestSerializeSingleOrAll_DifferentiatesSingleAndList(t *testing.T) {
	m1 := sampleMacroFixture("single-macro")
	singleBytes, err := SerializeSingleOrAll([]Macro{m1}, true, "json")
	if err != nil || len(singleBytes) == 0 {
		t.Fatalf("SerializeSingleOrAll single failed: %v", err)
	}

	if singleBytes[0] == '[' {
		t.Fatalf("expected JSON object starting with '{', got array: %s", string(singleBytes))
	}

	listBytes, err := SerializeSingleOrAll([]Macro{m1}, false, "json")
	if err != nil || len(listBytes) == 0 {
		t.Fatalf("SerializeSingleOrAll list failed: %v", err)
	}

	if listBytes[0] != '[' {
		t.Fatalf("expected JSON array starting with '[', got: %s", string(listBytes))
	}
}

func TestImportMacros_HandlesDryRunAndExceptFilter(t *testing.T) {
	m1 := sampleMacroFixture("dry-m1")
	m2 := sampleMacroFixture("dry-m2")
	res, err := ImportMacros([]Macro{m1, m2}, ImportOptions{
		IsDryRun:   true,
		ExceptList: []string{"dry-m2"},
	})
	if err != nil {
		t.Fatalf("ImportMacros error: %v", err)
	}

	if res.Imported != 1 || res.TotalFound != 2 {
		t.Fatalf("unexpected dry run result: %+v", res)
	}

	filteredRes, err := ImportMacros([]Macro{m1, m2}, ImportOptions{
		IsDryRun:   true,
		TargetName: "dry-m2",
	})
	if err != nil {
		t.Fatalf("ImportMacros with TargetName error: %v", err)
	}

	if filteredRes.Imported != 1 || filteredRes.Names[0] != "dry-m2" {
		t.Fatalf("expected only dry-m2 imported, got %+v", filteredRes)
	}
}
