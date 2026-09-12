package cmdmacro

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestParseMacroExportOpts_ParsesFlags(t *testing.T) {
	args := []string{"deploy", "-f", "out.yaml", "--all", "-except", "test1,test2"}
	opts := parseMacroExportOpts(args)

	if opts.TargetName != "deploy" {
		t.Fatalf("expected TargetName=deploy, got %s", opts.TargetName)
	}

	if opts.FilePath != "out.yaml" || opts.Format != constants.OutputYAML {
		t.Fatalf("expected out.yaml and yaml format, got %s, %s", opts.FilePath, opts.Format)
	}

	if !opts.IsAll || len(opts.ExceptList) != 2 {
		t.Fatalf("expected IsAll=true and 2 except tokens, got %+v", opts)
	}
}

func TestParseMacroImportOpts_ParsesFlags(t *testing.T) {
	args := []string{"backup.sqlitedb", "--name", "test-macro", "--force", "--dry-run", "--sqlitedb"}
	opts := parseMacroImportOpts(args)

	if opts.FilePath != "backup.sqlitedb" || opts.TargetName != "test-macro" {
		t.Fatalf("expected FilePath=backup.sqlitedb and TargetName=test-macro, got %s, %s", opts.FilePath, opts.TargetName)
	}

	if !opts.IsForce || !opts.IsDryRun || opts.Format != "sqlite" {
		t.Fatalf("unexpected import opts: %+v", opts)
	}
}

func TestInferMacroExportFormat_InfersFromExtension(t *testing.T) {
	optsYaml := macroExportOpts{FilePath: "archive.yaml"}
	inferMacroExportFormat(&optsYaml)
	if optsYaml.Format != constants.OutputYAML {
		t.Fatalf("expected yaml, got %s", optsYaml.Format)
	}

	optsDb := macroExportOpts{FilePath: "archive.sqlite"}
	inferMacroExportFormat(&optsDb)
	if optsDb.Format != "sqlite" {
		t.Fatalf("expected sqlite, got %s", optsDb.Format)
	}

	optsSqliteDb := macroExportOpts{FilePath: "archive.sqlitedb"}
	inferMacroExportFormat(&optsSqliteDb)
	if optsSqliteDb.Format != "sqlite" {
		t.Fatalf("expected sqlite for sqlitedb, got %s", optsSqliteDb.Format)
	}

	optsZip := macroExportOpts{FilePath: "archive.zip"}
	inferMacroExportFormat(&optsZip)
	if optsZip.Format != "zip" {
		t.Fatalf("expected zip, got %s", optsZip.Format)
	}
}

func TestParseMacroExportOpts_SingleFlag(t *testing.T) {
	args := []string{"single", "deploy", "--format", "json"}
	opts := parseMacroExportOpts(args)

	if !opts.IsSingle || opts.TargetName != "deploy" {
		t.Fatalf("expected IsSingle=true and TargetName=deploy, got %+v", opts)
	}
}

func TestParseMacroExportOpts_PositionalOutputFile(t *testing.T) {
	argsAll := []string{"all", "all_macros.sqlitedb"}
	optsAll := parseMacroExportOpts(argsAll)

	if !optsAll.IsAll || optsAll.FilePath != "all_macros.sqlitedb" || optsAll.Format != "sqlite" {
		t.Fatalf("expected IsAll=true, FilePath=all_macros.sqlitedb, Format=sqlite, got %+v", optsAll)
	}

	argsSingle := []string{"single", "deploy", "deploy.yaml"}
	optsSingle := parseMacroExportOpts(argsSingle)

	if !optsSingle.IsSingle || optsSingle.TargetName != "deploy" || optsSingle.FilePath != "deploy.yaml" || optsSingle.Format != constants.OutputYAML {
		t.Fatalf("expected IsSingle=true, TargetName=deploy, FilePath=deploy.yaml, got %+v", optsSingle)
	}
}

func TestReconcileImportPaths_NormalizesPositionalArgs(t *testing.T) {
	opts := macroImportOpts{TargetName: "all", FilePath: "backup.json"}
	reconcileImportPaths(&opts)

	if opts.FilePath != "backup.json" || opts.TargetName != "" {
		t.Fatalf("expected TargetName cleared for 'all', got %+v", opts)
	}
}

func TestParseMacroImportOpts_FormatsAndModifiers(t *testing.T) {
	args := []string{"backup.yaml", "--yaml", "--as", "new-macro", "--single"}
	opts := parseMacroImportOpts(args)

	if opts.FilePath != "backup.yaml" || opts.Format != constants.OutputYAML {
		t.Fatalf("expected backup.yaml and yaml format, got %+v", opts)
	}

	if opts.RenameAs != "new-macro" || !opts.IsSingle {
		t.Fatalf("expected RenameAs=new-macro and IsSingle=true, got %+v", opts)
	}
}
