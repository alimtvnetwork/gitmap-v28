// Package uipref centralizes runtime UX preferences (quiet output and
// no-color) so individual commands can consult a single source of
// truth instead of re-implementing env-var parsing.
//
// Honored environment variables:
//
//	GITMAP_QUIET=1     - suppress spinners and decorative banners
//	GITMAP_NO_COLOR=1  - strip / skip ANSI color escapes
//	NO_COLOR (any val) - de-facto cross-tool convention (https://no-color.org)
//
// Item #13 of the post-v6.53.0 suggestions list. Wired into the
// shared clone runner (clonespinner) first; future call sites
// (chrome profile copy, history, backup, ssh status) should consult
// IsQuiet() / IsNoColor() before emitting decorative output.
package uipref

import (
	"os"
	"path/filepath"
	"strings"
)

// Env var names — kept here (not constants/) because this package
// is the only legitimate reader.
const (
	EnvQuiet      = "GITMAP_QUIET"
	EnvNoColor    = "GITMAP_NO_COLOR"
	EnvNoColorStd = "NO_COLOR"
	EnvMacroPwd   = "GITMAP_MACRO_PWD"
)

var macroPwdOverride *bool

// IsMacroPwdVisible reports whether the current working directory banner
// should be displayed above each step prompt in the interactive macro builder.
func IsMacroPwdVisible() bool {
	if macroPwdOverride != nil {
		return *macroPwdOverride
	}

	if val, hasEnv := os.LookupEnv(EnvMacroPwd); hasEnv {
		return isTruthyMacroVal(val)
	}

	return isStoredMacroPwdEnabled()
}

// SetMacroPwdOverride sets a runtime override for Macro PWD visibility.
func SetMacroPwdOverride(isVisible bool) {
	val := isVisible
	macroPwdOverride = &val
}

// ResetMacroPwdOverride clears any runtime override for Macro PWD visibility.
func ResetMacroPwdOverride() {
	macroPwdOverride = nil
}

func isTruthyMacroVal(val string) bool {
	return val != "0" && val != "false" && val != "FALSE" && val != "off"
}

func getMacroPrefPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".gitmap", "macro_pwd.pref"), nil
}

func isStoredMacroPwdEnabled() bool {
	path, err := getMacroPrefPath()
	if err != nil {
		return true
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return true
	}

	return isTruthyMacroVal(strings.TrimSpace(string(data)))
}

// SetMacroShowPwd persists the user preference for displaying PWD in macro builder.
func SetMacroShowPwd(isVisible bool) error {
	SetMacroPwdOverride(isVisible)

	path, err := getMacroPrefPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	content := "1"
	if !isVisible {
		content = "0"
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// IsQuiet reports whether decorative / progress output should be
// suppressed. Treats any non-empty value (other than "0"/"false")
// as enabled.
func IsQuiet() bool { return truthyEnv(EnvQuiet) }

// IsNoColor reports whether ANSI color codes should be skipped.
// NO_COLOR is honored per the cross-tool convention: any value
// (including empty string set explicitly) disables color, but we
// also accept the project-scoped GITMAP_NO_COLOR for parity with
// GITMAP_QUIET.
func IsNoColor() bool {
	if _, set := os.LookupEnv(EnvNoColorStd); set {
		return true
	}

	return truthyEnv(EnvNoColor)
}

// truthyEnv returns true when v is set and not "", "0", "false".
func truthyEnv(name string) bool {
	v := os.Getenv(name)
	if v == "" || v == "0" || v == "false" || v == "FALSE" {
		return false
	}

	return true
}
