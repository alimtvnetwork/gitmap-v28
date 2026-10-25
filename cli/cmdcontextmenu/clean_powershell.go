package cmdcontextmenu

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const (
	hkcrPwshLegacyBg  = `HKCR\Directory\Background\shell\PowerShellMenu`
	hkcrPwsh7Bg       = `HKCR\Directory\Background\shell\PowerShell7x64`
	hkcrPwshLegacyDir = `HKCR\Directory\shell\PowerShellMenu`
	hkcrPwsh7Dir      = `HKCR\Directory\shell\PowerShell7x64`
	hkcuPwshLegacyBg  = `HKCU\Software\Classes\Directory\Background\shell\PowerShellMenu`
	hkcuPwshLegacyDir = `HKCU\Software\Classes\Directory\shell\PowerShellMenu`
	valLegacyDisable  = "LegacyDisable"
	valProgAccessOnly = "ProgrammaticAccessOnly"
)

type regRunnerFunc func(args ...string) ([]byte, error)

var (
	targetOS     = runtime.GOOS
	overrideKeys = []string{hkcuPwshLegacyBg, hkcuPwshLegacyDir}
	regRunner    regRunnerFunc = executeRegCommand
)

// PowerShellMenuStatus reports context menu status for PowerShell installations.
type PowerShellMenuStatus struct {
	HasLegacyMenu bool
	HasPwsh7Menu  bool
	IsDuplicate   bool
	IsCleaned     bool
}

// CleanDuplicatePowerShell cleans or restores legacy PowerShell context menu entries.
func CleanDuplicatePowerShell(isRestore, isDryRun bool) error {
	if targetOS != "windows" {
		fmt.Println("PowerShell context menu cleaner is only applicable on Windows.")

		return nil
	}
	if isRestore {
		return handlePowerShellRestore(isDryRun)
	}

	return handlePowerShellClean(isDryRun)
}

// DetectPowerShellStatus checks the presence of PowerShell menu keys in Windows registry.
func DetectPowerShellStatus() PowerShellMenuStatus {
	hasLegacy := isKeyPresent(hkcrPwshLegacyBg) || isKeyPresent(hkcrPwshLegacyDir)
	hasPwsh7 := isKeyPresent(hkcrPwsh7Bg) || isKeyPresent(hkcrPwsh7Dir)
	isCleaned := isKeyPresent(hkcuPwshLegacyBg) || isKeyPresent(hkcuPwshLegacyDir)

	return PowerShellMenuStatus{
		HasLegacyMenu: hasLegacy,
		HasPwsh7Menu:  hasPwsh7,
		IsDuplicate:   hasLegacy && hasPwsh7,
		IsCleaned:     isCleaned,
	}
}

func handlePowerShellClean(isDryRun bool) error {
	if isDryRun {
		printCleanDryRun()

		return nil
	}

	return executePowerShellClean()
}

func handlePowerShellRestore(isDryRun bool) error {
	if isDryRun {
		printRestoreDryRun()

		return nil
	}

	return executePowerShellRestore()
}

func printCleanDryRun() {
	fmt.Println("[dry-run] Would apply HKCU context menu overrides:")
	for _, key := range overrideKeys {
		fmt.Printf("  + %s (%s, %s)\n", key, valLegacyDisable, valProgAccessOnly)
	}
}

func printRestoreDryRun() {
	fmt.Println("[dry-run] Would remove HKCU context menu overrides:")
	for _, key := range overrideKeys {
		fmt.Printf("  - %s\n", key)
	}
}

func executePowerShellClean() error {
	for _, key := range overrideKeys {
		if err := applyOverrideKey(key); err != nil {
			return err
		}
	}
	fmt.Printf("✓ %sLegacy PowerShell context menu disabled.%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Println("  PowerShell 7 remains active in Explorer.")

	return nil
}

func executePowerShellRestore() error {
	for _, key := range overrideKeys {
		if err := removeOverrideKey(key); err != nil {
			return err
		}
	}
	fmt.Printf("✓ %sRestored default PowerShell context menu entries.%s\n", constants.ColorGreen, constants.ColorReset)

	return nil
}

func applyOverrideKey(key string) error {
	if err := writeRegValue(key, valLegacyDisable); err != nil {
		return err
	}

	return writeRegValue(key, valProgAccessOnly)
}

func writeRegValue(key, valName string) error {
	_, err := regRunner("add", key, "/v", valName, "/t", "REG_SZ", "/d", "", "/f")
	if err != nil {
		return apperror.WrapSimple(err, "failed to set registry value: "+valName)
	}

	return nil
}

func removeOverrideKey(key string) error {
	isExist := isKeyPresent(key)
	if !isExist {
		return nil
	}
	_, err := regRunner("delete", key, "/f")
	if err != nil {
		return apperror.WrapSimple(err, "failed to remove registry override: "+key)
	}

	return nil
}

func isKeyPresent(keyPath string) bool {
	_, err := regRunner("query", keyPath)

	return err == nil
}

func executeRegCommand(args ...string) ([]byte, error) {
	cmd := exec.Command("reg", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, err
	}

	return out, nil
}
