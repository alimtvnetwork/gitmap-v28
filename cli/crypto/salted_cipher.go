package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	asciiPrintableStart = 32
	asciiPrintableEnd   = 126
	asciiPrintableRange = 95 // 126 - 32 + 1
)

// GenerateRandomSalt creates a cryptographically random hex salt of the given byte length.
func GenerateRandomSalt(byteLen int) string {
	if byteLen <= 0 {
		byteLen = 8
	}
	buf := make([]byte, byteLen)
	if _, err := rand.Read(buf); err != nil {
		return "a1b2c3d4e5f60718"
	}
	return hex.EncodeToString(buf)
}

// EncryptSalted encrypts plaintext using a salted character rotation cipher.
// Output format: salt:<salt_hex>:<base64_payload>
func EncryptSalted(plain string, salt string) string {
	if plain == "" {
		return ""
	}
	if salt == "" {
		salt = GenerateRandomSalt(8)
	}
	saltBytes := []byte(salt)
	saltLen := len(saltBytes)
	plainRunes := []rune(plain)
	shifted := make([]rune, len(plainRunes))

	for i, r := range plainRunes {
		if r >= asciiPrintableStart && r <= asciiPrintableEnd {
			sByte := int(saltBytes[i%saltLen])
			shift := (sByte + (i+1)*7) % asciiPrintableRange
			val := int(r) - asciiPrintableStart
			newVal := (val + shift) % asciiPrintableRange
			shifted[i] = rune(asciiPrintableStart + newVal)
			continue
		}
		shifted[i] = r
	}

	enc := base64.StdEncoding.EncodeToString([]byte(string(shifted)))
	return fmt.Sprintf("salt:%s:%s", salt, enc)
}

// DecryptSalted decrypts a salted rotation cipher string.
func DecryptSalted(cipherText string) (string, error) {
	if cipherText == "" {
		return "", nil
	}
	if !strings.HasPrefix(cipherText, "salt:") {
		return cipherText, nil
	}
	parts := strings.Split(cipherText, ":")
	if len(parts) != 3 {
		return "", errors.New("invalid salted cipher format, expected salt:<salt>:<payload>")
	}
	salt, payload := parts[1], parts[2]
	if salt == "" || payload == "" {
		return "", errors.New("missing salt or payload in salted cipher")
	}

	shiftedBytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", fmt.Errorf("decode salted cipher payload: %w", err)
	}

	saltBytes := []byte(salt)
	saltLen := len(saltBytes)
	shiftedRunes := []rune(string(shiftedBytes))
	plainRunes := make([]rune, len(shiftedRunes))

	for i, r := range shiftedRunes {
		if r >= asciiPrintableStart && r <= asciiPrintableEnd {
			sByte := int(saltBytes[i%saltLen])
			shift := (sByte + (i+1)*7) % asciiPrintableRange
			val := int(r) - asciiPrintableStart
			origVal := normalizeAsciiVal(val, shift)
			plainRunes[i] = rune(asciiPrintableStart + origVal)
			continue
		}
		plainRunes[i] = r
	}

	return string(plainRunes), nil
}

func normalizeAsciiVal(val, shift int) int {
	origVal := (val - shift) % asciiPrintableRange
	if origVal < 0 {
		return origVal + asciiPrintableRange
	}
	return origVal
}

func resolveSaltByte(saltBytes []byte, saltLen, index int) int {
	if saltLen <= 0 {
		return 0
	}
	return int(saltBytes[index%saltLen])
}

// EncryptCaesar encrypts plaintext using a Caesar cipher with optional salt.
// Format: caesar:<shift>:<payload> or caesar:<shift>:<salt>:<payload>
func EncryptCaesar(plain string, shift int, salt string) string {
	if plain == "" {
		return ""
	}
	normShift := shift % asciiPrintableRange
	if normShift < 0 {
		normShift += asciiPrintableRange
	}
	saltBytes := []byte(salt)
	saltLen := len(saltBytes)
	plainRunes := []rune(plain)
	shifted := make([]rune, len(plainRunes))

	for i, r := range plainRunes {
		if r >= asciiPrintableStart && r <= asciiPrintableEnd {
			sByte := resolveSaltByte(saltBytes, saltLen, i)
			curShift := (normShift + sByte) % asciiPrintableRange
			val := int(r) - asciiPrintableStart
			newVal := (val + curShift) % asciiPrintableRange
			shifted[i] = rune(asciiPrintableStart + newVal)
			continue
		}
		shifted[i] = r
	}

	enc := base64.StdEncoding.EncodeToString([]byte(string(shifted)))
	if salt == "" {
		return fmt.Sprintf("caesar:%d:%s", normShift, enc)
	}
	return fmt.Sprintf("caesar:%d:%s:%s", normShift, salt, enc)
}

// DecryptCaesar decrypts a Caesar cipher string.
func DecryptCaesar(cipherText string) (string, error) {
	if cipherText == "" {
		return "", nil
	}
	if !strings.HasPrefix(cipherText, "caesar:") {
		return cipherText, nil
	}
	parts := strings.Split(cipherText, ":")
	if len(parts) < 3 {
		return "", errors.New("invalid caesar cipher format, expected caesar:<shift>:<payload> or caesar:<shift>:<salt>:<payload>")
	}

	shift, parseErr := strconv.Atoi(parts[1])
	if parseErr != nil {
		return "", fmt.Errorf("invalid caesar shift value: %w", parseErr)
	}
	normShift := shift % asciiPrintableRange
	if normShift < 0 {
		normShift += asciiPrintableRange
	}

	salt := ""
	payload := parts[2]
	if len(parts) == 4 {
		salt = parts[2]
		payload = parts[3]
	}

	shiftedBytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", fmt.Errorf("decode caesar cipher payload: %w", err)
	}

	saltBytes := []byte(salt)
	saltLen := len(saltBytes)
	shiftedRunes := []rune(string(shiftedBytes))
	plainRunes := make([]rune, len(shiftedRunes))

	for i, r := range shiftedRunes {
		if r >= asciiPrintableStart && r <= asciiPrintableEnd {
			sByte := resolveSaltByte(saltBytes, saltLen, i)
			curShift := (normShift + sByte) % asciiPrintableRange
			val := int(r) - asciiPrintableStart
			origVal := normalizeAsciiVal(val, curShift)
			plainRunes[i] = rune(asciiPrintableStart + origVal)
			continue
		}
		plainRunes[i] = r
	}

	return string(plainRunes), nil
}
