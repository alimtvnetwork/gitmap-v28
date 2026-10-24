package constants

import (
	"strings"
	"testing"
)

func TestConstantsInvariants(t *testing.T) {
	testFlagConstants(t)
	testBinaryConstants(t)
	testPathConstants(t)
	testUIConstants(t)
	testDraculaColors(t)
	testErrorCode(t)
}

func testFlagConstants(t *testing.T) {
	if FlagLongInstall != "--install" || FlagShortInstall != "-i" {
		t.Error("install flag mismatch")
	}
	if FlagLongHelp != "--help" || FlagShortHelp != "-h" {
		t.Error("help flag mismatch")
	}
	if FlagLongYes != "--yes" || FlagShortYes != "-y" {
		t.Error("yes flag mismatch")
	}
	if FlagLongVersion != "--version" || FlagShortVersion != "-V" {
		t.Error("version flag mismatch")
	}
	if FlagLongRelease != "--release" || FlagShortRelease != "-r" {
		t.Error("release flag mismatch")
	}
	if FlagLongCheck != "--check" || FlagDoubleDash != "--" {
		t.Error("check or separator flag mismatch")
	}
}

func testBinaryConstants(t *testing.T) {
	if BinCargo != "cargo" || BinRustc != "rustc" || BinRustup != "rustup" || BinClang != "clang" {
		t.Error("binary name mismatch")
	}
}

func testPathConstants(t *testing.T) {
	if DirCargo != ".cargo" || DirBin != "bin" || ExtExe != ".exe" {
		t.Error("path token mismatch")
	}
	if PathUsrBin != "/usr/bin" || PathUsrLocalBin != "/usr/local/bin" {
		t.Error("system path token mismatch")
	}
}

func testUIConstants(t *testing.T) {
	if DisplayNameCargo != "Cargo" || DisplayNameRustc != "Rustc" || DisplayNameRustup != "Rustup" {
		t.Error("display name mismatch")
	}
	if StatusHeaderTitle != "Cargo & Rust Toolchain Status:" {
		t.Error("status header title mismatch")
	}
	if StatusInstalled != "installed" || StatusNotInstalled != "Not installed" {
		t.Error("status label mismatch")
	}
	if !strings.Contains(WingetInstallRust, "winget") || !strings.Contains(BrewInstallRust, "brew") {
		t.Error("suggestion text mismatch")
	}
}

func testDraculaColors(t *testing.T) {
	colors := []string{
		ColorDraculaPurple,
		ColorDraculaWhite,
		ColorDraculaGreen,
		ColorDraculaRed,
		ColorDraculaGray,
		ColorDraculaCyan,
	}

	for _, c := range colors {
		if !strings.HasPrefix(c, "#") {
			t.Errorf("color %q must be a hex string", c)
		}
	}
}

func testErrorCode(t *testing.T) {
	if ErrCodeMissingCargo != "E7100" {
		t.Errorf("expected ErrCodeMissingCargo to be E7100, got %q", ErrCodeMissingCargo)
	}
	if len(MissingCargoErrorMsg) == 0 {
		t.Error("expected non-empty MissingCargoErrorMsg")
	}
}
