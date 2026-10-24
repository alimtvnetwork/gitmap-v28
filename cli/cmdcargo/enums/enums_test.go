package enums

import (
	"testing"
)

func TestCargoSubcommandType(t *testing.T) {
	for _, sc := range AllCargoSubcommandTypes() {
		if !sc.IsValid() {
			t.Errorf("expected valid for %s", sc)
		}
		if sc.String() != string(sc) {
			t.Errorf("expected string representation to match %s", sc)
		}
		if !sc.HasSubcommand() {
			t.Errorf("expected HasSubcommand to be true for %s", sc)
		}
	}

	if CargoSubcommandType("invalid").IsValid() {
		t.Error("expected invalid subcommand to not be valid")
	}

	parsed, err := ParseCargoSubcommandType("  STATUS ")
	if err != nil || parsed != CargoSubcommandStatus {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	_, err = ParseCargoSubcommandType("unknown")
	if err == nil {
		t.Error("expected error for unknown subcommand")
	}

	if !CargoSubcommandStatus.IsStatus() || !CargoSubcommandSt.IsStatus() || !CargoSubcommandInfo.IsStatus() {
		t.Error("expected status variants to return true for IsStatus")
	}
	if CargoSubcommandInstall.IsStatus() {
		t.Error("install should not be status")
	}

	if !CargoSubcommandInstall.IsInstall() || !CargoSubcommandIn.IsInstall() {
		t.Error("expected install variants to return true for IsInstall")
	}
	if CargoSubcommandStatus.IsInstall() {
		t.Error("status should not be install")
	}
}

func TestCargoCommandType(t *testing.T) {
	all := AllCargoCommandTypes()
	if len(all) != 13 {
		t.Errorf("expected 13 cargo commands, got %d", len(all))
	}

	for _, cmd := range all {
		if !cmd.IsValid() {
			t.Errorf("expected valid for %s", cmd)
		}
		if cmd.String() != string(cmd) {
			t.Errorf("expected string representation to match %s", cmd)
		}
		if !cmd.HasCommand() {
			t.Errorf("expected HasCommand to be true for %s", cmd)
		}
	}

	if CargoCommandType("nope").IsValid() {
		t.Error("expected invalid command to not be valid")
	}

	parsed, err := ParseCargoCommandType("  CLIPPY ")
	if err != nil || parsed != CargoCmdClippy {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	_, err = ParseCargoCommandType("unknown")
	if err == nil {
		t.Error("expected error parsing unknown command")
	}

	testCommandFlags(t)
}

func testCommandFlags(t *testing.T) {
	if !CargoCmdBuild.IsBuild() || !CargoCmdCheck.IsCheck() || !CargoCmdTest.IsTest() {
		t.Error("command predicate failure for build/check/test")
	}
	if !CargoCmdClippy.IsClippy() || !CargoCmdRun.IsRun() || !CargoCmdClean.IsClean() {
		t.Error("command predicate failure for clippy/run/clean")
	}
	if !CargoCmdFmt.IsFmt() || !CargoCmdDoc.IsDoc() || !CargoCmdUpdate.IsUpdate() {
		t.Error("command predicate failure for fmt/doc/update")
	}
	if !CargoCmdInit.IsInit() || !CargoCmdNew.IsNew() || !CargoCmdBench.IsBench() || !CargoCmdPublish.IsPublish() {
		t.Error("command predicate failure for init/new/bench/publish")
	}
}

func TestRustToolchainBinaryType(t *testing.T) {
	all := AllRustToolchainBinaryTypes()
	if len(all) != 4 {
		t.Errorf("expected 4 rust binaries, got %d", len(all))
	}

	for _, bin := range all {
		if !bin.IsValid() {
			t.Errorf("expected valid for %s", bin)
		}
		if bin.String() != string(bin) {
			t.Errorf("expected string representation to match %s", bin)
		}
		if !bin.HasBinary() {
			t.Errorf("expected HasBinary to be true for %s", bin)
		}
	}

	if RustToolchainBinaryType("gcc").IsValid() {
		t.Error("expected invalid binary to not be valid")
	}

	parsed, err := ParseRustToolchainBinaryType("  RUSTUP ")
	if err != nil || parsed != RustBinaryRustup {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	_, err = ParseRustToolchainBinaryType("unknown")
	if err == nil {
		t.Error("expected error parsing unknown binary")
	}

	if !RustBinaryCargo.IsCargo() || !RustBinaryRustc.IsRustc() || !RustBinaryRustup.IsRustup() || !RustBinaryClang.IsClang() {
		t.Error("binary predicate failure")
	}
}

func TestToolchainEnvVarType(t *testing.T) {
	all := AllToolchainEnvVarTypes()
	if len(all) != 9 {
		t.Errorf("expected 9 env vars, got %d", len(all))
	}

	for _, ev := range all {
		if !ev.IsValid() {
			t.Errorf("expected valid for %s", ev)
		}
		if ev.String() != string(ev) {
			t.Errorf("expected string representation to match %s", ev)
		}
		if !ev.HasEnvVar() {
			t.Errorf("expected HasEnvVar to be true for %s", ev)
		}
	}

	if ToolchainEnvVarType("HOME").IsValid() {
		t.Error("expected invalid env var to not be valid")
	}

	parsed, err := ParseToolchainEnvVarType("  cargo_home ")
	if err != nil || parsed != ToolchainEnvCargoHome {
		t.Errorf("unexpected parse result: %v, err: %v", parsed, err)
	}

	_, err = ParseToolchainEnvVarType("invalid_env")
	if err == nil {
		t.Error("expected error parsing unknown env var")
	}

	testEnvVarFlags(t)
}

func testEnvVarFlags(t *testing.T) {
	if !ToolchainEnvLibclangPath.IsLibclangPath() || !ToolchainEnvPath.IsPath() {
		t.Error("env predicate failure for libclang_path/path")
	}
	if !ToolchainEnvCargoHome.IsCargoHome() || !ToolchainEnvRustupHome.IsRustupHome() {
		t.Error("env predicate failure for cargo_home/rustup_home")
	}
	if !ToolchainEnvSystemRoot.IsSystemRoot() || !ToolchainEnvComSpec.IsComSpec() {
		t.Error("env predicate failure for systemroot/comspec")
	}
	if !ToolchainEnvTmp.IsTmp() || !ToolchainEnvTemp.IsTemp() || !ToolchainEnvUserProfile.IsUserProfile() {
		t.Error("env predicate failure for tmp/temp/userprofile")
	}
}
