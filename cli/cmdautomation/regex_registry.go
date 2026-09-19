package cmdautomation

import (
	"bytes"
	"regexp"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RegexRegistry provides thread-safe lazy regex compilation and caching.
type RegexRegistry struct {
	mu    sync.RWMutex
	cache map[string]*regexp.Regexp
}

var globalRegexRegistry = &RegexRegistry{
	cache: make(map[string]*regexp.Regexp),
}

// GetRegex lazily retrieves or compiles a regular expression.
func GetRegex(pattern string, isCaseInsensitive bool) (*regexp.Regexp, *apperror.AppError) {
	return globalRegexRegistry.Get(pattern, isCaseInsensitive)
}

func (r *RegexRegistry) Get(pattern string, isCaseInsensitive bool) (*regexp.Regexp, *apperror.AppError) {
	key := buildRegexKey(pattern, isCaseInsensitive)
	r.mu.RLock()
	cached, hasPattern := r.cache[key]
	r.mu.RUnlock()

	if hasPattern {
		return cached, nil
	}

	return r.compileAndCache(key, pattern, isCaseInsensitive)
}

func (r *RegexRegistry) compileAndCache(key, pattern string, isCaseInsensitive bool) (*regexp.Regexp, *apperror.AppError) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if cached, hasPattern := r.cache[key]; hasPattern {
		return cached, nil
	}

	raw := buildRawPattern(pattern, isCaseInsensitive)
	compiled, err := regexp.Compile(raw)
	if err != nil {
		ctx := map[string]any{"pattern": pattern, "err": err.Error()}
		return nil, apperror.New("regex_compile", "E_REGEX_INVALID", ctx)
	}

	r.cache[key] = compiled
	return compiled, nil
}

func buildRegexKey(pattern string, isCaseInsensitive bool) string {
	if isCaseInsensitive {
		return "ci:" + pattern
	}
	return "cs:" + pattern
}

func buildRawPattern(pattern string, isCaseInsensitive bool) string {
	if isCaseInsensitive {
		return "(?i)" + pattern
	}
	return pattern
}

// HasLiteralMatch executes high-speed Boyer-Moore substring search without regex overhead.
func HasLiteralMatch(content []byte, substr string, isCaseInsensitive bool) bool {
	if isCaseInsensitive {
		lowerContent := bytes.ToLower(content)
		lowerSubstr := strings.ToLower(substr)
		return bytes.Contains(lowerContent, []byte(lowerSubstr))
	}
	return bytes.Contains(content, []byte(substr))
}
