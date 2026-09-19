package cmdautomation

import (
	"encoding/json"
	"strings"
	"unicode/utf16"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	// EncodingUTF8 represents standard UTF-8 stream encoding.
	EncodingUTF8 = "utf-8"

	// EncodingUTF16LE represents UTF-16 Little Endian stream encoding.
	EncodingUTF16LE = "utf-16le"
)

// IsUTF16Encoding reports whether the requested encoding is UTF-16 / UTF-16LE.
func IsUTF16Encoding(encoding string) bool {
	lower := strings.ToLower(strings.TrimSpace(encoding))
	return lower == "utf-16" || lower == "utf16" || lower == "utf-16le" || lower == "utf16le"
}

// NormalizeEncoding normalizes the requested encoding string to canonical form.
func NormalizeEncoding(encoding string) string {
	if IsUTF16Encoding(encoding) {
		return EncodingUTF16LE
	}
	return EncodingUTF8
}

// NegotiateEncoding negotiates the stream encoding based on request and runtime.
func NegotiateEncoding(requested, runtimeName string) string {
	if requested != "" {
		return NormalizeEncoding(requested)
	}
	return EncodingUTF8
}

// EncodeToStream encodes a FileContext into JSON bytes with the negotiated encoding.
func EncodeToStream(ctx FileContext, encoding string) ([]byte, *apperror.AppError) {
	jsonData, err := json.Marshal(ctx)
	if err != nil {
		return nil, apperror.WrapSimple(err, "EncodeToStream")
	}
	if IsUTF16Encoding(encoding) {
		return encodeUTF16LE(jsonData), nil
	}
	return append(jsonData, '\n'), nil
}

func encodeUTF16LE(data []byte) []byte {
	runes := []rune(string(data))
	u16s := utf16.Encode(runes)
	buf := make([]byte, 2+len(u16s)*2+2)
	buf[0], buf[1] = 0xFF, 0xFE
	idx := 2
	for _, v := range u16s {
		buf[idx] = byte(v)
		buf[idx+1] = byte(v >> 8)
		idx += 2
	}
	buf[idx], buf[idx+1] = 0x0A, 0x00
	return buf
}

// DecodeFromStream decodes raw stream bytes into a string, auto-detecting BOMs.
func DecodeFromStream(data []byte, encoding string) string {
	if len(data) == 0 {
		return ""
	}
	if hasUTF16LEBOM(data) {
		return decodeUTF16LE(data[2:])
	}
	if hasUTF8BOM(data) {
		return string(data[3:])
	}
	if IsUTF16Encoding(encoding) {
		return decodeUTF16LE(data)
	}
	return string(data)
}

func hasUTF16LEBOM(data []byte) bool {
	return len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE
}

func hasUTF8BOM(data []byte) bool {
	return len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF
}

func decodeUTF16LE(data []byte) string {
	pairCount := len(data) / 2
	u16s := make([]uint16, pairCount)
	for i := 0; i < pairCount; i++ {
		u16s[i] = uint16(data[i*2]) | (uint16(data[i*2+1]) << 8)
	}
	return string(utf16.Decode(u16s))
}
