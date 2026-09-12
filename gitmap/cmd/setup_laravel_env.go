package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const (
	laravelAppKeyBytes  = 32
	laravelAppKeyPrefix = "base64:"
)

// LaravelEnvOptions contains configurable options for Laravel .env synthesis.
type LaravelEnvOptions struct {
	AppName         string
	AppEnv          string
	AppKey          string
	AppDebug        string
	AppUrl          string
	DBConnection    string
	DBHost          string
	DBPort          string
	DBDatabase      string
	DBUsername      string
	DBPassword      string
	CacheStore      string
	SessionDriver   string
	QueueConnection string
}

// GenerateLaravelAppKey generates a cryptographically secure 32-byte Base64 APP_KEY.
func GenerateLaravelAppKey() (string, *apperror.AppError) {
	bytes := make([]byte, laravelAppKeyBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", apperror.WrapSimple(err, "crypto.rand.Read")
	}

	encoded := base64.StdEncoding.EncodeToString(bytes)

	return laravelAppKeyPrefix + encoded, nil
}

func extractEnvKey(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	isComment := strings.HasPrefix(trimmed, "#")
	if isComment {
		return "", false
	}

	idx := strings.Index(trimmed, "=")
	isKV := idx > 0
	if isKV {
		return strings.TrimSpace(trimmed[:idx]), true
	}

	return "", false
}

func formatEnvValue(val string) string {
	hasSpace := strings.Contains(val, " ")
	if hasSpace {
		return fmt.Sprintf("%q", val)
	}

	return val
}

func formatEnvLine(key, val string) string {
	return fmt.Sprintf("%s=%s", key, formatEnvValue(val))
}

func resolveUpdatedEnvLine(key, originalLine string, updates map[string]string, handled map[string]bool) string {
	newVal, hasUpdate := updates[key]
	if hasUpdate {
		handled[key] = true

		return formatEnvLine(key, newVal)
	}

	return originalLine
}

func processEnvLine(line string, updates map[string]string, handled map[string]bool) string {
	key, hasKey := extractEnvKey(line)
	if hasKey {
		return resolveUpdatedEnvLine(key, line, updates, handled)
	}

	return line
}

func appendIfMissing(lines []string, key string, updates map[string]string, handled map[string]bool) []string {
	isHandled := handled[key]
	if isHandled {
		return lines
	}

	val, hasVal := updates[key]
	if hasVal {
		handled[key] = true

		return append(lines, formatEnvLine(key, val))
	}

	return lines
}

// MergeEnvContent updates existing keys in place and appends new keys in order.
func MergeEnvContent(baseContent string, updates map[string]string, orderedKeys []string) string {
	lines := strings.Split(baseContent, constants.NewLineUnix)
	handled := make(map[string]bool)
	out := make([]string, 0, len(lines)+len(updates))
	for _, line := range lines {
		out = append(out, processEnvLine(line, updates, handled))
	}

	for _, key := range orderedKeys {
		out = appendIfMissing(out, key, updates, handled)
	}

	return strings.Join(out, constants.NewLineUnix)
}

func addEnvOption(m map[string]string, keys *[]string, key, val string) {
	hasVal := val != ""
	if hasVal {
		m[key] = val
		*keys = append(*keys, key)
	}
}

// BuildLaravelEnvMap converts options into a key-value map and ordered key list.
func BuildLaravelEnvMap(opts LaravelEnvOptions) (map[string]string, []string) {
	m := make(map[string]string)
	keys := make([]string, 0, 14)
	addEnvOption(m, &keys, "APP_NAME", opts.AppName)
	addEnvOption(m, &keys, "APP_ENV", opts.AppEnv)
	addEnvOption(m, &keys, "APP_KEY", opts.AppKey)
	addEnvOption(m, &keys, "APP_DEBUG", opts.AppDebug)
	addEnvOption(m, &keys, "APP_URL", opts.AppUrl)
	addEnvOption(m, &keys, "DB_CONNECTION", opts.DBConnection)
	addEnvOption(m, &keys, "DB_HOST", opts.DBHost)
	addEnvOption(m, &keys, "DB_PORT", opts.DBPort)
	addEnvOption(m, &keys, "DB_DATABASE", opts.DBDatabase)
	addEnvOption(m, &keys, "DB_USERNAME", opts.DBUsername)
	addEnvOption(m, &keys, "DB_PASSWORD", opts.DBPassword)
	addEnvOption(m, &keys, "CACHE_STORE", opts.CacheStore)
	addEnvOption(m, &keys, "SESSION_DRIVER", opts.SessionDriver)
	addEnvOption(m, &keys, "QUEUE_CONNECTION", opts.QueueConnection)

	return m, keys
}

func resolveLaravelAppKey(currentKey string) (string, *apperror.AppError) {
	hasKey := currentKey != ""
	if hasKey {
		return currentKey, nil
	}

	return GenerateLaravelAppKey()
}

// SynthesizeLaravelEnv synthesizes an ordered Laravel .env file with cryptographic APP_KEY.
func SynthesizeLaravelEnv(baseContent string, opts LaravelEnvOptions) (string, *apperror.AppError) {
	key, keyErr := resolveLaravelAppKey(opts.AppKey)
	if keyErr != nil {
		return "", keyErr
	}

	finalOpts := opts
	finalOpts.AppKey = key
	updates, orderedKeys := BuildLaravelEnvMap(finalOpts)

	return MergeEnvContent(baseContent, updates, orderedKeys), nil
}
