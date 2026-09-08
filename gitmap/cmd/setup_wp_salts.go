package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// WpSaltLength defines the standard 64-character length for WordPress salts.
const (
	WpSaltLength  = 64
	WpSaltCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_ []{}<>~+=;:?"
)

// WpSalts holds the 8 cryptographically generated WordPress security keys and salts.
type WpSalts struct {
	AuthKey        string
	SecureAuthKey  string
	LoggedInKey    string
	NonceKey       string
	AuthSalt       string
	SecureAuthSalt string
	LoggedInSalt   string
	NonceSalt      string
}

var (
	wpSaltsMarkerRegex  = regexp.MustCompile(`(?s)// >>> gitmap:wp-salts >>>.*?// <<< gitmap:wp-salts <<<`)
	wpSaltsDefaultRegex = regexp.MustCompile(`(?m)(^define\s*\(\s*['"](?:AUTH|SECURE_AUTH|LOGGED_IN|NONCE)_(?:KEY|SALT)['"].*?\n)+`)
)

func fillRandomChars(buf []byte) *apperror.AppError {
	maxVal := big.NewInt(int64(len(WpSaltCharset)))
	for i := range buf {
		num, err := rand.Int(rand.Reader, maxVal)
		if err != nil {
			return apperror.WrapSimple(err, "crypto.rand.Int")
		}
		buf[i] = WpSaltCharset[num.Int64()]
	}

	return nil
}

// GenerateRandomSalt creates a cryptographically random salt string of the given length.
func GenerateRandomSalt(length int) (string, *apperror.AppError) {
	isInvalidLength := length <= 0
	if isInvalidLength {
		return "", apperror.NewValidationError("salt length must be positive")
	}
	buf := make([]byte, length)
	err := fillRandomChars(buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func buildWpSaltsStruct(s []string) *WpSalts {
	return &WpSalts{
		AuthKey:        s[0],
		SecureAuthKey:  s[1],
		LoggedInKey:    s[2],
		NonceKey:       s[3],
		AuthSalt:       s[4],
		SecureAuthSalt: s[5],
		LoggedInSalt:   s[6],
		NonceSalt:      s[7],
	}
}

// GenerateWpSalts generates 8 unique 64-character salts using pure Go crypto/rand.
func GenerateWpSalts() (*WpSalts, *apperror.AppError) {
	salts := make([]string, 8)
	for i := 0; i < 8; i++ {
		val, err := GenerateRandomSalt(WpSaltLength)
		if err != nil {
			return nil, err
		}
		salts[i] = val
	}

	return buildWpSaltsStruct(salts), nil
}

func formatSaltDefine(key, value string) string {
	return fmt.Sprintf("define('%s', '%s');%s", key, value, constants.NewLineUnix)
}

// FormatWpSaltsPHP formats the 8 salts as PHP define statements.
func FormatWpSaltsPHP(salts *WpSalts) string {
	var b strings.Builder
	b.WriteString(formatSaltDefine("AUTH_KEY", salts.AuthKey))
	b.WriteString(formatSaltDefine("SECURE_AUTH_KEY", salts.SecureAuthKey))
	b.WriteString(formatSaltDefine("LOGGED_IN_KEY", salts.LoggedInKey))
	b.WriteString(formatSaltDefine("NONCE_KEY", salts.NonceKey))
	b.WriteString(formatSaltDefine("AUTH_SALT", salts.AuthSalt))
	b.WriteString(formatSaltDefine("SECURE_AUTH_SALT", salts.SecureAuthSalt))
	b.WriteString(formatSaltDefine("LOGGED_IN_SALT", salts.LoggedInSalt))
	b.WriteString(formatSaltDefine("NONCE_SALT", salts.NonceSalt))

	return b.String()
}

// WrapWpSaltsMarker wraps WordPress salts in GitMap marker comment lines.
func WrapWpSaltsMarker(content string) string {
	trimmed := strings.TrimSpace(content)

	return fmt.Sprintf("// >>> gitmap:wp-salts >>>%s%s%s// <<< gitmap:wp-salts <<<%s",
		constants.NewLineUnix, trimmed, constants.NewLineUnix, constants.NewLineUnix)
}

func appendWpSalts(originalContent string, newBlock string) string {
	target := "/* That's all, stop editing!"
	idx := strings.Index(originalContent, target)
	isFound := idx >= 0
	if isFound {
		before := originalContent[:idx]
		after := originalContent[idx:]
		return before + newBlock + constants.NewLineUnix + after
	}

	return originalContent + constants.NewLineUnix + newBlock
}

func injectFreshSalts(originalContent string, newBlock string) string {
	hasDefaultSalts := wpSaltsDefaultRegex.MatchString(originalContent)
	if hasDefaultSalts {
		return wpSaltsDefaultRegex.ReplaceAllLiteralString(originalContent, strings.TrimSpace(newBlock)+constants.NewLineUnix)
	}

	return appendWpSalts(originalContent, newBlock)
}

// InjectWpSalts injects or replaces WordPress security salts idempotently.
func InjectWpSalts(originalContent string, salts *WpSalts) string {
	newBlock := WrapWpSaltsMarker(FormatWpSaltsPHP(salts))
	hasMarker := wpSaltsMarkerRegex.MatchString(originalContent)
	if hasMarker {
		return wpSaltsMarkerRegex.ReplaceAllLiteralString(originalContent, strings.TrimSpace(newBlock))
	}

	return injectFreshSalts(originalContent, newBlock)
}
