package funcintel

import (
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

var detectors = map[string]Detector{}

var extToLang = map[string]string{
	".go":   constants.CommitInLanguageGo,
	".js":   constants.CommitInLanguageJavaScript,
	".mjs":  constants.CommitInLanguageJavaScript,
	".cjs":  constants.CommitInLanguageJavaScript,
	".ts":   constants.CommitInLanguageTypeScript,
	".tsx":  constants.CommitInLanguageTypeScript,
	".rs":   constants.CommitInLanguageRust,
	".py":   constants.CommitInLanguagePython,
	".php":  constants.CommitInLanguagePhp,
	".java": constants.CommitInLanguageJava,
	".cs":   constants.CommitInLanguageCSharp,
}

// LanguageForPath returns the FunctionIntelLanguage enum token for
// the given relative path, or "" if no detector covers the extension.
func LanguageForPath(rel string) string {
	return extToLang[strings.ToLower(filepath.Ext(rel))]
}

// Get returns the Detector registered for `lang`, or nil.
func Get(lang string) Detector { return detectors[lang] }

func register(lang string, d Detector) { detectors[lang] = d }

// EnabledLanguages filters to languages that have a registered detector.
func EnabledLanguages(wanted []string) []string {
	out := make([]string, 0, len(wanted))
	for _, l := range wanted {
		if _, ok := detectors[l]; ok {
			out = append(out, l)
		}
	}
	return out
}
