//go:build linux

package power

import (
	"fmt"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type linuxDriver struct {
	run CmdRunner
}

func init() {
	RegisterDriver("linux", func(r CmdRunner) Manager {
		return &linuxDriver{run: r}
	})
}

func (l *linuxDriver) Platform() string {
	return "linux"
}

func (l *linuxDriver) GetStatus() (Settings, error) {
	displayMins, err := l.queryGnomeDisplay()
	if err != nil {
		displayMins = 0
	}

	sleepMins, err := l.queryGnomeSleep()
	if err != nil {
		sleepMins = 0
	}

	isNever := displayMins == 0 && sleepMins == 0

	return Settings{
		Platform:              "linux",
		DisplayTimeoutMinutes: displayMins,
		SleepTimeoutMinutes:   sleepMins,
		IsNeverSleep:          isNever,
		Source:                "gsettings",
	}, nil
}

func (l *linuxDriver) queryGnomeDisplay() (int, error) {
	out, err := l.run("gsettings", "get", "org.gnome.desktop.session", "idle-delay")
	if err != nil {
		return 0, apperror.WrapSimple(err, "linux.queryGnomeDisplay")
	}

	return ParseGnomeTimeoutMinutes(string(out)), nil
}

func (l *linuxDriver) queryGnomeSleep() (int, error) {
	out, err := l.run("gsettings", "get", "org.gnome.settings-daemon.plugins.power", "sleep-inactive-ac-timeout")
	if err != nil {
		return 0, apperror.WrapSimple(err, "linux.queryGnomeSleep")
	}

	return ParseGnomeTimeoutMinutes(string(out)), nil
}

func (l *linuxDriver) SetNeverSleep() error {
	_ = l.setGnomeNeverSleep()
	_ = l.setXsetNeverSleep()

	return nil
}

func (l *linuxDriver) setGnomeNeverSleep() error {
	commands := [][]string{
		{"set", "org.gnome.desktop.session", "idle-delay", "0"},
		{"set", "org.gnome.settings-daemon.plugins.power", "sleep-inactive-ac-timeout", "0"},
		{"set", "org.gnome.settings-daemon.plugins.power", "sleep-inactive-ac-type", "nothing"},
		{"set", "org.gnome.desktop.screensaver", "lock-enabled", "false"},
	}

	for _, args := range commands {
		_, _ = l.run("gsettings", args...)
	}

	return nil
}

func (l *linuxDriver) setXsetNeverSleep() error {
	_, _ = l.run("xset", "s", "off")
	_, _ = l.run("xset", "-dpms")
	_, _ = l.run("xset", "s", "0", "0")

	return nil
}

func (l *linuxDriver) SetTimeouts(displayMinutes, sleepMinutes int) error {
	displaySec := strconv.Itoa(displayMinutes * 60)
	sleepSec := strconv.Itoa(sleepMinutes * 60)

	commands := [][]string{
		{"set", "org.gnome.desktop.session", "idle-delay", displaySec},
		{"set", "org.gnome.settings-daemon.plugins.power", "sleep-inactive-ac-timeout", sleepSec},
		{"set", "org.gnome.desktop.screensaver", "lock-enabled", "true"},
	}

	for _, args := range commands {
		_, err := l.run("gsettings", args...)
		if err != nil {
			return apperror.WrapSimple(err, fmt.Sprintf("gsettings %v", args))
		}
	}

	return nil
}

func (l *linuxDriver) ApplySettings(s Settings) error {
	if s.IsNeverSleep {
		return l.SetNeverSleep()
	}

	return l.SetTimeouts(s.DisplayTimeoutMinutes, s.SleepTimeoutMinutes)
}
