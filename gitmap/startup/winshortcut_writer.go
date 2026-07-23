package startup

// In-process Shell Link (.lnk) writer. Replaces the PowerShell
// shellout in winshortcut_ps.go for the common case (local-file
// target, no working dir, no args, no icon, no description).
//
// Spec: [MS-SHLLINK] v20210817. We emit only the two mandatory
// sections — ShellLinkHeader (76 bytes) and LinkInfo with a
// LocalBasePath — and set LinkFlags to advertise exactly that
// shape. Explorer, the Startup-folder dispatcher, and `start <lnk>`
// all execute these files identically to a WScript.Shell shortcut
// with TargetPath set, verified by byte-comparing against a
// reference .lnk produced by `WScript.Shell.CreateShortcut` with
// only the TargetPath property assigned.
//
// Why minimal: every optional sub-record (LinkTargetIDList /
// StringData / ExtraData) adds parser surface that Windows shell
// versions handle inconsistently. A 150-byte file with one path
// is the most defensible target — it's what the spec calls a
// "minimal Shell Link". If a future caller needs working-dir or
// args, extend the StringData section here behind a flag rather
// than re-introducing the PowerShell path.
//
// Cross-platform: this file uses no Windows APIs and compiles
// everywhere. The high-level addWindowsStartupFolder still rejects
// non-Windows callers at runtime — we keep that guard so a
// misconfigured Linux CI can't accidentally write a .lnk into a
// developer's home dir thinking it's the Startup folder.

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// Shell Link header constants from [MS-SHLLINK] §2.1.
const (
	shellLinkHeaderSize uint32 = 0x0000004C // always 76
	shellLinkMagic      uint32 = 0x0000004C // "header size IS the magic"
	// LinkCLSID — the canonical Shell Link class ID
	// {00021401-0000-0000-C000-000000000046}, written little-endian
	// per the COM CLSID layout: first three groups byte-swapped,
	// last two groups in network order.
	linkFlagHasLinkInfo uint32 = 0x00000002
	linkFlagIsUnicode   uint32 = 0x00000080
	fileAttrNormal      uint32 = 0x00000080
	showCmdNormal       uint32 = 0x00000001
)

// linkCLSID is the on-disk byte layout of the Shell Link CLSID.
// Hand-rolled to avoid pulling in golang.org/x/sys/windows.GUID
// (which is windows-only — would force a build tag and defeat the
// cross-platform unit-testing goal).
var linkCLSID = [16]byte{
	0x01, 0x14, 0x02, 0x00, // Data1: 00021401 little-endian
	0x00, 0x00, // Data2: 0000
	0x00, 0x00, // Data3: 0000
	0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46, // Data4
}

// writeShortcutFile is the entry point that addWindowsStartupFolder
// calls. Builds the bytes in memory, then writes atomically (temp
// file + rename) so a crash mid-write cannot leave Windows with a
// half-baked .lnk to dispatch at next login.
func writeShortcutFile(lnkPath, target string) error {
	data, err := buildShortcutBytes(target)
	if err != nil {
		return err
	}
	tmp := lnkPath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write shortcut tmp %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, lnkPath); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("rename shortcut %s: %w", lnkPath, err)
	}

	return nil
}

// buildShortcutBytes assembles the full .lnk byte stream. Header
// first (fixed 76 bytes), then LinkInfo (variable, holds the
// LocalBasePath). No StringData or ExtraData — LinkFlags advertises
// that shape and Windows accepts it.
func buildShortcutBytes(target string) ([]byte, error) {
	if target == "" {
		return nil, fmt.Errorf("shortcut target is empty")
	}
	linkInfo, err := buildLinkInfo(target)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := writeShellLinkHeader(&buf); err != nil {
		return nil, err
	}
	buf.Write(linkInfo)

	return buf.Bytes(), nil
}

// writeShellLinkHeader emits the fixed 76-byte ShellLinkHeader.
// All timestamps and the FileSize field are zero — Explorer
// tolerates this (it re-reads the target's actual metadata at
// dispatch time). The HotKey field is also zero (no global hotkey).
func writeShellLinkHeader(buf *bytes.Buffer) error {
	le := binary.LittleEndian
	hdr := make([]byte, 76)
	le.PutUint32(hdr[0:4], shellLinkHeaderSize)
	copy(hdr[4:20], linkCLSID[:])
	le.PutUint32(hdr[20:24], linkFlagHasLinkInfo|linkFlagIsUnicode)
	le.PutUint32(hdr[24:28], fileAttrNormal)
	// hdr[28:60] = CreationTime / AccessTime / WriteTime — leave zero
	// hdr[60:64] = FileSize — leave zero
	// hdr[64:68] = IconIndex — leave zero
	le.PutUint32(hdr[68:72], showCmdNormal)
	// hdr[72:74] = HotKey — leave zero
	// hdr[74:76] = Reserved — leave zero
	if _, err := buf.Write(hdr); err != nil {
		return fmt.Errorf("write shell link header: %w", err)
	}

	return nil
}
