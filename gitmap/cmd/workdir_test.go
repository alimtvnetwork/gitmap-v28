package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func TestWorkDirCLIParse(t *testing.T) {
	opts := parseWorkDirFlags([]string{"add", "/home/user/work", "--label", "my-work"})
	if opts.Action != "add" || opts.Target != "/home/user/work" || opts.Label != "my-work" {
		t.Fatalf("unexpected parsed options: %+v", opts)
	}

	optsEmpty := parseWorkDirFlags([]string{})
	if optsEmpty.Action != "help" {
		t.Fatalf("expected default action help, got %s", optsEmpty.Action)
	}

	optsHelp := parseWorkDirFlags([]string{"help"})
	if optsHelp.Action != "help" {
		t.Fatalf("expected action help, got %s", optsHelp.Action)
	}

	optsDef := parseWorkDirFlags([]string{"default"})
	if optsDef.Action != "default" {
		t.Fatalf("expected action default, got %s", optsDef.Action)
	}
}

func TestWorkDirHelpAndUsage(t *testing.T) {
	optsHelp := workDirOptions{Action: "help"}
	err := dispatchWorkDirAction(optsHelp)
	if err != nil {
		t.Fatalf("dispatchWorkDirAction help returned error: %v", err)
	}

	optsDashH := workDirOptions{Action: "-h"}
	errDashH := dispatchWorkDirAction(optsDashH)
	if errDashH != nil {
		t.Fatalf("dispatchWorkDirAction -h returned error: %v", errDashH)
	}
}

func TestCDWorkDirKeyword(t *testing.T) {
	if !isWorkDirKeyword("work") {
		t.Error("expected 'work' to be a workdir keyword")
	}

	if !isWorkDirKeyword("default") {
		t.Error("expected 'default' to be a workdir keyword")
	}

	if !isWorkDirKeyword("workdir") {
		t.Error("expected 'workdir' to be a workdir keyword")
	}

	if isWorkDirKeyword("unknown-repo") {
		t.Error("did not expect 'unknown-repo' to be a workdir keyword")
	}
}

func TestAutoRegisterFirstWorkDir(t *testing.T) {
	tempDir := t.TempDir()
	db, errDB := store.OpenDefault()
	if errDB != nil {
		t.Skip("sqlite unavailable")
	}
	defer db.Close()
	ensureWorkDirsTableExists()

	_, _ = db.SQL().Exec("DELETE FROM work_directories")

	isRegistered := autoRegisterFirstWorkDir(tempDir, true)
	if !isRegistered {
		t.Fatalf("expected autoRegisterFirstWorkDir to return true on empty store")
	}

	def, errDef := db.GetDefaultWorkDir()
	if errDef != nil || def == nil {
		t.Fatalf("expected default workdir to be set, got err: %v", errDef)
	}
	if def.AbsolutePath != tempDir {
		t.Fatalf("expected default workdir %s, got %s", tempDir, def.AbsolutePath)
	}

	isSecondRegistered := autoRegisterFirstWorkDir(tempDir, true)
	if isSecondRegistered {
		t.Fatalf("expected second autoRegisterFirstWorkDir to return false")
	}
}

func TestWorkDirAddDefaultsAndDuplicates(t *testing.T) {
	tempDir := t.TempDir()
	db, errDB := store.OpenDefault()
	if errDB != nil {
		t.Skip("sqlite unavailable")
	}
	defer db.Close()
	ensureWorkDirsTableExists()

	_, _ = db.SQL().Exec("DELETE FROM work_directories")

	errAdd := runWorkDirAdd(tempDir, "temp-label")
	if errAdd != nil {
		t.Fatalf("runWorkDirAdd failed: %v", errAdd)
	}

	errDuplicate := runWorkDirAdd(tempDir, "temp-label")
	if errDuplicate != nil {
		t.Fatalf("expected duplicate runWorkDirAdd to succeed with notice, got: %v", errDuplicate)
	}

	errEmpty := runWorkDirAdd("", "empty-target")
	if errEmpty != nil {
		t.Fatalf("expected empty target to succeed, got: %v", errEmpty)
	}
}
