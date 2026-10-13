package cmdinstall

import (
	"bytes"
	"os"
)

func isCompressedArchiveMagic(h []byte) bool {
	if len(h) < 4 {
		return false
	}
	isGzip := h[0] == 0x1f && h[1] == 0x8b
	isZip := h[0] == 0x50 && h[1] == 0x4b
	isBz2 := h[0] == 0x42 && h[1] == 0x5a && h[2] == 0x68
	isZst := bytes.Equal(h[:4], []byte{0x28, 0xb5, 0x2f, 0xfd})
	return isGzip || isZip || isBz2 || isZst
}

func isContainerOrExecutableMagic(h []byte) bool {
	if len(h) < 6 {
		return false
	}
	isXz := bytes.Equal(h[:6], []byte{0xfd, 0x37, 0x7a, 0x58, 0x5a, 0x00})
	is7z := bytes.Equal(h[:6], []byte{0x37, 0x7a, 0xbc, 0xaf, 0x27, 0x1c})
	isExe := h[0] == 'M' && h[1] == 'Z'
	isElf := bytes.Equal(h[:4], []byte{0x7f, 'E', 'L', 'F'})
	return isXz || is7z || isExe || isElf
}

func isDmgOrAppleMagic(h []byte) bool {
	if len(h) < 4 {
		return false
	}
	isMachO := (h[0] == 0xfe && h[1] == 0xed && h[2] == 0xfa && (h[3] == 0xce || h[3] == 0xcf)) ||
		(h[0] == 0xcf && h[1] == 0xfa && h[2] == 0xed && h[3] == 0xfe)
	isAppleXar := bytes.Equal(h[:4], []byte{0x78, 0x61, 0x72, 0x21})
	isZlib := h[0] == 0x78 && (h[1] == 0x01 || h[1] == 0x9c || h[1] == 0xda)
	return isMachO || isAppleXar || isZlib
}

func isTarArchiveMagic(h []byte) bool {
	if len(h) < 262 {
		return false
	}
	return bytes.HasPrefix(h[257:], []byte("ustar"))
}

func isSupportedArchiveHeader(h []byte) bool {
	return isCompressedArchiveMagic(h) || isContainerOrExecutableMagic(h) || isTarArchiveMagic(h) || isDmgOrAppleMagic(h)
}

func readArchiveHeaderBytes(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "archive.open")
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, readErr := io.ReadFull(f, buf)
	if readErr != nil && readErr != io.ErrUnexpectedEOF && readErr != io.EOF {
		return nil, apperror.WrapSimple(readErr, "archive.readHeader")
	}
	return buf[:n], nil
}

func isDmgKolyTrailer(filePath string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil || stat.Size() < 512 {
		return false
	}
	trailer := make([]byte, 512)
	_, readErr := f.ReadAt(trailer, stat.Size()-512)
	return readErr == nil && bytes.Equal(trailer[:4], []byte("koly"))
}

func VerifyArchiveHeaderMagic(filePath string) bool {
	header, err := readArchiveHeaderBytes(filePath)
	if err != nil || len(header) < 4 {
		return false
	}
	if isSupportedArchiveHeader(header) {
		return true
	}
	return isDmgKolyTrailer(filePath)
}

