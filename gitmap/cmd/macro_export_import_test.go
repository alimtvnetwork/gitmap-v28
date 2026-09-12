package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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

	optsZip := macroExportOpts{FilePath: "archive.zip"}
	inferMacroExportFormat(&optsZip)
	if optsZip.Format != "zip" {
		t.Fatalf("expected zip, got %s", optsZip.Format)
	}
}
