package constants

// Compress and checksum messages.
const (
	MsgCompressArchive   = "  ✓ Compressed %s → %s\n"
	ErrCompressFailed    = "  ✗ Failed to compress %s: %v\n"
	FlagDescCompress     = "Wrap release assets in .zip (Windows) or .tar.gz (Linux/macOS)"
	MsgChecksumGenerated = "  ✓ Generated %s (SHA256)\n"
	ErrChecksumFailed    = "  ✗ Failed to hash %s: %v\n"
	FlagDescChecksums    = "Generate SHA256 checksums.txt for release assets"
	ChecksumsFile        = "checksums.txt"
)
