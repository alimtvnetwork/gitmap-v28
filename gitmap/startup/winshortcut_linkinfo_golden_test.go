package startup

import (
	"encoding/hex"
	"testing"
)

// TestBuildLinkInfoGoldenBytes pins the exact byte layout of the
// LinkInfo block for a known target. The fixture was derived by
// hand from [MS-SHLLINK] §2.3 — any drift here means the on-disk
// .lnk format we emit changed and Windows shell parsers may break.
func TestBuildLinkInfoGoldenBytes(t *testing.T) {
	const target = `C:\a.exe`
	// Layout for target "C:\a.exe" (9 bytes incl. NUL):
	//   header(28) + VolumeID(16) + path(9) + suffix(1) = 54 bytes (0x36)
	const wantHex = "360000001c000000010000001c0000002c00000000000000350000" +
		"0010000000030000000000000014000000433a5c612e6578650000"
	want, err := hex.DecodeString(wantHex)
	if err != nil {
		t.Fatalf("decode golden hex: %v", err)
	}

	got, err := buildLinkInfo(target)
	if err != nil {
		t.Fatalf("buildLinkInfo: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("length mismatch: got %d, want %d\n got=%s\nwant=%s",
			len(got), len(want), hex.EncodeToString(got), wantHex)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("byte %d (0x%02x): got 0x%02x, want 0x%02x\n got=%s\nwant=%s",
				i, i, got[i], want[i], hex.EncodeToString(got), wantHex)
		}
	}
}

// TestBuildLinkInfoRejectsNULTarget pins the only documented error
// path of buildLinkInfo: a NUL byte inside the target string.
func TestBuildLinkInfoRejectsNULTarget(t *testing.T) {
	if _, err := buildLinkInfo("C:\\bad\x00path.exe"); err == nil {
		t.Fatalf("expected error for NUL-containing target, got nil")
	}
}
