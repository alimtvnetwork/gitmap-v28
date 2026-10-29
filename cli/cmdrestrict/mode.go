package cmdrestrict

import (
	"fmt"
	"strings"
)

const (
	ModeAuto        = "auto"
	ModeFlag        = "flag"
	settingsModeKey = "restriction.enforcement"
)

// SetStoreDir points restriction reads and writes at dir. An empty dir uses the user store.
func SetStoreDir(dir string) {
	restrictionsDirOverride = dir
}

// CurrentMode reports auto or flag. A missing file means auto.
func CurrentMode() string {
	file, err := loadRestrictions()
	if err != nil {
		return ModeAuto
	}

	return normalizeStoredMode(file.EnforcementMode)
}

// IsAutoEnforced reports whether every repository on this machine enforces the window.
func IsAutoEnforced() bool {
	return CurrentMode() == ModeAuto
}

// ActiveTimezone returns the timezone of the active restriction.
func ActiveTimezone() string {
	file, err := loadRestrictions()
	if err != nil {
		return defaultTimezone
	}
	item := activeRestriction(file)
	if item.Timezone == "" {
		return defaultTimezone
	}

	return item.Timezone
}

// SetMode stores auto or flag for every repository on this machine.
func SetMode(mode string) error {
	normalized, err := normalizeMode(mode)
	if err != nil {
		return err
	}
	file, errLoad := loadRestrictions()
	if errLoad != nil {
		return errLoad
	}
	file.EnforcementMode = normalized

	return saveRestrictions(file)
}

// ExposeMode copies the machine mode into settings attributes so the UI shows it.
func ExposeMode(attributes map[string]string) {
	if attributes == nil {
		return
	}
	attributes[settingsModeKey] = CurrentMode()
}

// ApplySettingsMode writes the settings choice into the machine restriction file.
func ApplySettingsMode(attributes map[string]string) error {
	if attributes == nil {
		return nil
	}
	mode := strings.TrimSpace(attributes[settingsModeKey])
	if mode == "" {
		return nil
	}

	return SetMode(mode)
}

func normalizeMode(mode string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case ModeAuto:
		return ModeAuto, nil
	case ModeFlag:
		return ModeFlag, nil
	default:
		return "", fmt.Errorf("enforcement mode must be auto or flag")
	}
}

func normalizeStoredMode(mode string) string {
	if strings.EqualFold(strings.TrimSpace(mode), ModeFlag) {
		return ModeFlag
	}

	return ModeAuto
}

func runMode(args []string) error {
	if len(args) == 0 {
		fmt.Printf("Enforcement: %s\n", describeMode(CurrentMode()))

		return nil
	}
	if err := SetMode(args[0]); err != nil {
		return err
	}
	fmt.Printf("Enforcement: %s\n", describeMode(CurrentMode()))

	return nil
}

func describeMode(mode string) string {
	if mode == ModeFlag {
		return "flag (only when --schedule is passed)"
	}

	return "auto (every repository on this machine)"
}
