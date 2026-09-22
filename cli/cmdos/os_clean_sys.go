package cmdos

import (
	"fmt"
	"runtime"
)

func runOSSystemClean() error {
	if runtime.GOOS != "linux" {
		fmt.Printf("ℹ System package clean is designed for Linux (APT/DNF/Pacman). Current OS is %s.\n", runtime.GOOS)
		return nil
	}
	engine := GetSystemCleanEngine()
	fmt.Println("▶ Cleaning package manager cache...")
	_, _ = engine.CleanSystemPackages()

	fmt.Println("▶ Vacuuming systemd journal logs (older than 3 days)...")
	_ = engine.VacuumJournals()

	fmt.Println("✔ System cleanup complete.")
	return nil
}

func isSystemCleanTarget(arg string) bool {
	return arg == "sys" || arg == "system" || arg == "--system"
}
