package startup

// LinkInfo sub-section of a Shell Link. Per [MS-SHLLINK] §2.3 this
// block tells Windows where the target file lives. We populate the
// "LocalBasePath" variant only — sufficient for any local .exe and
// the form WScript.Shell.CreateShortcut produces when only the
// TargetPath property is set.
//
// Layout (all little-endian):
//
//   uint32  LinkInfoSize             // total bytes of this block
//   uint32  LinkInfoHeaderSize       // 0x1C (no optional Unicode offsets)
//   uint32  LinkInfoFlags            // 0x01 = VolumeIDAndLocalBasePath
//   uint32  VolumeIDOffset           // 0x1C (right after header)
//   uint32  LocalBasePathOffset      // 0x1C + sizeof(VolumeID)
//   uint32  CommonNetworkRelativeLinkOffset  // 0 (none)
//   uint32  CommonPathSuffixOffset   // points to a single 0x00 terminator
//   <VolumeID block>
//   <LocalBasePath ASCII, NUL-terminated>
//   <CommonPathSuffix = single 0x00>
//
// We embed a minimal VolumeID (drive type Fixed, no label) because
// some shell versions reject LinkInfo without one even when the
// LocalBasePath flag is set.

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// LinkInfo + VolumeID constants from [MS-SHLLINK] §2.3.
const (
	linkInfoHeaderSize       uint32 = 0x0000001C
	linkInfoFlagVolumeID     uint32 = 0x00000001
	volumeIDHeaderSize       uint32 = 0x00000010
	driveTypeFixed           uint32 = 0x00000003
	driveSerialDefault       uint32 = 0x00000000
	volumeLabelOffsetMissing uint32 = 0x00000014 // "no label" sentinel
)

// buildLinkInfo assembles the LinkInfo block as a single byte slice.
// Returns an error only when the target path contains characters
// LocalBasePath cannot represent (NUL bytes — would break the
// NUL-terminated ASCII encoding).
//
// Example:
//
//	raw, err := buildLinkInfo(`C:\Tools\gitmap.exe`)
//	// raw[0:4]  = total LinkInfoSize (little-endian)
//	// raw[4:8]  = 0x1C header size
//	// raw[8:12] = 0x01 VolumeIDAndLocalBasePath flag
//	// raw[0x1C:0x2C]                  = VolumeID block
//	// raw[0x2C:0x2C+len(target)+1]    = target + 0x00
//	// raw[len(raw)-1]                 = 0x00 CommonPathSuffix terminator
func buildLinkInfo(target string) ([]byte, error) {
	if bytes.IndexByte([]byte(target), 0) >= 0 {

		return nil, fmt.Errorf("shortcut target contains NUL byte")
	}
	volumeID := buildVolumeID()
	pathBytes := append([]byte(target), 0x00) // NUL-terminated ASCII
	suffixBytes := []byte{0x00}               // empty CommonPathSuffix
	offsets, err := computeLinkInfoOffsets(volumeID, pathBytes, suffixBytes)
	if err != nil {

		return nil, err
	}

	return assembleLinkInfo(offsets, volumeID, pathBytes, suffixBytes), nil
}

// linkInfoOffsets caches every uint32 offset/size used by the
// LinkInfo header so the assembler is a flat sequence of PutUint32
// calls with no inline arithmetic that could re-trigger G115.
type linkInfoOffsets struct {
	volOff, pathOff, suffixOff, totalSize uint32
}

// computeLinkInfoOffsets converts the variable-width section sizes
// to uint32 once, refusing pathological lengths instead of silently
// wrapping. Returns the precomputed offsets ready to write.
//
// Example: with a 16-byte VolumeID, a 20-byte path ("C:\\foo.exe\x00"
// is 12 bytes; pretend 20 for illustration), and a 1-byte suffix:
//
//	volOff    = 0x1C            (right after 28-byte header)
//	pathOff   = 0x1C + 16 = 0x2C
//	suffixOff = 0x2C + 20 = 0x40
//	totalSize = 0x40 + 1  = 0x41
func computeLinkInfoOffsets(volumeID, pathBytes, suffixBytes []byte) (linkInfoOffsets, error) {
	volSize, err := safeUint32(len(volumeID))
	if err != nil {

		return linkInfoOffsets{}, fmt.Errorf("volumeID size: %w", err)
	}
	pathSize, err := safeUint32(len(pathBytes))
	if err != nil {

		return linkInfoOffsets{}, fmt.Errorf("target path size: %w", err)
	}
	suffixSize, err := safeUint32(len(suffixBytes))
	if err != nil {

		return linkInfoOffsets{}, fmt.Errorf("suffix size: %w", err)
	}
	volOff := linkInfoHeaderSize
	pathOff := volOff + volSize
	suffixOff := pathOff + pathSize

	return linkInfoOffsets{volOff, pathOff, suffixOff, suffixOff + suffixSize}, nil
}

// assembleLinkInfo writes the header + payload sections into a
// single buffer. Pure formatting — no arithmetic, no error path.
//
// Example output for target "C:\\a.exe" (8 bytes + NUL = 9):
//
//	out[0:4]   = 0x39 0x00 0x00 0x00   // totalSize = 0x39
//	out[4:8]   = 0x1C 0x00 0x00 0x00   // header size
//	out[12:16] = 0x1C 0x00 0x00 0x00   // VolumeIDOffset
//	out[16:20] = 0x2C 0x00 0x00 0x00   // LocalBasePathOffset
//	out[0x2C:0x35] = "C:\\a.exe\x00"
//	out[0x38]      = 0x00              // CommonPathSuffix
func assembleLinkInfo(o linkInfoOffsets, volumeID, pathBytes, suffixBytes []byte) []byte {
	le := binary.LittleEndian
	out := make([]byte, o.totalSize)
	le.PutUint32(out[0:4], o.totalSize)
	le.PutUint32(out[4:8], linkInfoHeaderSize)
	le.PutUint32(out[8:12], linkInfoFlagVolumeID)
	le.PutUint32(out[12:16], o.volOff)
	le.PutUint32(out[16:20], o.pathOff)
	le.PutUint32(out[20:24], 0) // no CommonNetworkRelativeLink
	le.PutUint32(out[24:28], o.suffixOff)
	copy(out[o.volOff:], volumeID)
	copy(out[o.pathOff:], pathBytes)
	copy(out[o.suffixOff:], suffixBytes)

	return out
}

// safeUint32 narrows a non-negative int to uint32, refusing values
// that would silently wrap on a 64-bit platform. Centralized so the
// gosec G115 guard is uniform across every offset/size we compute.
func safeUint32(n int) (uint32, error) {
	if n < 0 || n > math.MaxUint32 {

		return 0, fmt.Errorf("value %d out of uint32 range", n)
	}

	return uint32(n), nil
}

// buildVolumeID emits a minimal VolumeID block: 16-byte header,
// fixed drive type, zero serial, "label missing" sentinel. No
// VolumeLabel string follows — the sentinel offset 0x14 is the
// spec-defined "no label" marker and equals VolumeIDSize, which
// Windows treats as end-of-block.
func buildVolumeID() []byte {
	le := binary.LittleEndian
	vol := make([]byte, volumeIDHeaderSize)
	le.PutUint32(vol[0:4], volumeIDHeaderSize)  // VolumeIDSize
	le.PutUint32(vol[4:8], driveTypeFixed)      // DriveType
	le.PutUint32(vol[8:12], driveSerialDefault) // DriveSerialNumber
	le.PutUint32(vol[12:16], volumeLabelOffsetMissing)

	return vol
}
