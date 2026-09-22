//go:build linux

package cmdos

import (
	"os/exec"
)

type linuxSystemCleanEngine struct{}

func newPlatformSystemCleanEngine() SystemCleanOperator {
	return &linuxSystemCleanEngine{}
}

func (l *linuxSystemCleanEngine) CleanSystemPackages() (int64, error) {
	if isExecutable("apt-get") {
		_ = exec.Command("apt-get", "clean").Run()
		_ = exec.Command("apt-get", "autoremove", "-y").Run()
		return 0, nil
	}
	if isExecutable("dnf") {
		_ = exec.Command("dnf", "clean", "all").Run()
		return 0, nil
	}
	if isExecutable("pacman") {
		_ = exec.Command("pacman", "-Sc", "--noconfirm").Run()
		return 0, nil
	}
	return 0, nil
}

func (l *linuxSystemCleanEngine) VacuumJournals() error {
	if isExecutable("journalctl") {
		_ = exec.Command("journalctl", "--vacuum-time=3d").Run()
	}
	return nil
}

func isExecutable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
