package startup

// Cross-platform structural tests for the in-process Shell Link
// writer. No build tag — these run on Linux CI and validate the
// .lnk byte structure against [MS-SHLLINK] §2.1 (ShellLinkHeader)
// and §2.3 (LinkInfo) without needing Windows UI execution.
//
// The tests assert:
//   1. Header magic + size                  (4-byte signature)
//   2. CLSID byte layout                    (canonical Shell Link)
//   3. LinkFlags = HasLinkInfo | IsUnicode  (no other flags set)
//   4. FileAttributes = FILE_ATTRIBUTE_NORMAL
//   5. ShowCommand = SW_SHOWNORMAL
//   6. LinkInfo offsets are internally consistent
//   7. Embedded LocalBasePath round-trips bytewise
//   8. Total length matches the spec formula
//
// What we don't assert: actual Windows execution. That's covered
// by the Windows-tagged round-trip test in windows_test.go.

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
)

// TestBuildShortcutBytes_HeaderShape locks in the 76-byte header
// layout. Any change to LinkFlags / ShowCommand would shift Windows
// shell behavior and must be a deliberate decision.
func TestBuildShortcutBytes_HeaderShape(t *testing.T) {
	data := mustBuildShortcut(t, `C:\Windows\System32\notepad.exe`)
	le := binary.LittleEndian

	if got := le.Uint32(data[0:4]); got != 0x0000004C {
		t.Errorf("HeaderSize = %#x, want 0x4C", got)
	}
	if !bytes.Equal(data[4:20], linkCLSID[:]) {
		t.Errorf("CLSID = %x, want %x", data[4:20], linkCLSID[:])
	}
	wantFlags := uint32(0x02 | 0x80) // HasLinkInfo | IsUnicode
	if got := le.Uint32(data[20:24]); got != wantFlags {
		t.Errorf("LinkFlags = %#x, want %#x", got, wantFlags)
	}
	if got := le.Uint32(data[24:28]); got != 0x80 {
		t.Errorf("FileAttributes = %#x, want FILE_ATTRIBUTE_NORMAL", got)
	}
	if got := le.Uint32(data[68:72]); got != 0x01 {
		t.Errorf("ShowCommand = %#x, want SW_SHOWNORMAL", got)
	}
}

// TestBuildShortcutBytes_LinkInfoOffsets verifies the LinkInfo
// internal offset chain (every offset points inside the block and
// each section is reachable). A broken offset would cause Windows
// shell to fall back to "this shortcut is empty" behavior at
// dispatch time.
func TestBuildShortcutBytes_LinkInfoOffsets(t *testing.T) {
	target := `C:\Program Files\App\app.exe`
	data := mustBuildShortcut(t, target)
	li := data[76:] // LinkInfo starts right after header
	le := binary.LittleEndian

	totalSize := le.Uint32(li[0:4])
	headerSize := le.Uint32(li[4:8])
	flags := le.Uint32(li[8:12])
	volOff := le.Uint32(li[12:16])
	pathOff := le.Uint32(li[16:20])
	suffixOff := le.Uint32(li[24:28])

	if int(totalSize) != len(li) {
		t.Errorf("LinkInfoSize = %d, want %d", totalSize, len(li))
	}
	if headerSize != 0x1C {
		t.Errorf("LinkInfoHeaderSize = %#x, want 0x1C", headerSize)
	}
	if flags&0x01 == 0 {
		t.Errorf("LinkInfoFlags missing VolumeIDAndLocalBasePath bit")
	}
	if volOff != 0x1C || pathOff <= volOff || suffixOff <= pathOff {
		t.Errorf("offsets non-monotonic: vol=%d path=%d suffix=%d",
			volOff, pathOff, suffixOff)
	}
}

// TestBuildShortcutBytes_EmbeddedPath confirms the LocalBasePath
// is written verbatim as NUL-terminated ASCII at the offset the
// LinkInfo header advertises. This is the byte the Windows shell
// reads to find what to launch.
func TestBuildShortcutBytes_EmbeddedPath(t *testing.T) {
	target := `D:\tools\watcher.exe`
	data := mustBuildShortcut(t, target)
	li := data[76:]
	pathOff := binary.LittleEndian.Uint32(li[16:20])

	embedded := readCString(li[pathOff:])
	if embedded != target {
		t.Errorf("embedded path = %q, want %q", embedded, target)
	}
}

// TestBuildShortcutBytes_RejectsEmpty confirms the writer fails
// loudly on an empty target. A silent zero-byte LocalBasePath
// would produce a .lnk Windows treats as "broken shortcut" with
// a confusing dialog at next login.
func TestBuildShortcutBytes_RejectsEmpty(t *testing.T) {
	if _, err := buildShortcutBytes(""); err == nil {
		t.Fatal("expected error for empty target, got nil")
	}
}

// TestBuildShortcutBytes_RejectsNUL confirms NUL bytes inside the
// target are caught before they corrupt the NUL-terminated ASCII
// LocalBasePath encoding.
func TestBuildShortcutBytes_RejectsNUL(t *testing.T) {
	if _, err := buildShortcutBytes("C:\\bad\x00path.exe"); err == nil {
		t.Fatal("expected error for NUL byte in target, got nil")
	}
}

// TestBuildShortcutBytes_LengthFormula asserts the total file size
// matches the spec formula: 76 (header) + 28 (LinkInfo header) +
// 16 (VolumeID) + len(target)+1 + 1 (CommonPathSuffix terminator).
// Catches any future drift in the LinkInfo assembly that doesn't
// show up in the offset test.
func TestBuildShortcutBytes_LengthFormula(t *testing.T) {
	target := `C:\x.exe`
	data := mustBuildShortcut(t, target)
	want := 76 + 28 + 16 + len(target) + 1 + 1
	if len(data) != want {
		t.Errorf("total length = %d, want %d", len(data), want)
	}
}

// mustBuildShortcut is a tiny test helper that builds bytes and
// fails the test on error. Keeps each test case to one assertion
// concept.
func mustBuildShortcut(t *testing.T, target string) []byte {
	t.Helper()
	data, err := buildShortcutBytes(target)
	if err != nil {
		t.Fatalf("buildShortcutBytes(%q): %v", target, err)
	}

	return data
}

// readCString reads bytes up to the first NUL. Used to extract the
// embedded LocalBasePath without pulling in golang.org/x/text or
// guessing termination from offsets.
func readCString(b []byte) string {
	idx := bytes.IndexByte(b, 0)
	if idx < 0 {
		return string(b)
	}

	return strings.Clone(string(b[:idx]))
}
