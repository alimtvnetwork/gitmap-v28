//go:build windows

package power

import (
	"fmt"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type windowsDriver struct {
	run CmdRunner
}

func init() {
	RegisterDriver("windows", func(r CmdRunner) Manager {
		return &windowsDriver{run: r}
	})
}

func (w *windowsDriver) Platform() string {
	return "windows"
}

func (w *windowsDriver) GetStatus() (Settings, error) {
	displayMins, err := w.queryDisplayTimeout()
	if err != nil {
		displayMins = 0
	}

	sleepMins, err := w.querySleepTimeout()
	if err != nil {
		sleepMins = 0
	}

	isNever := displayMins == 0 && sleepMins == 0

	return Settings{
		Platform:              "windows",
		DisplayTimeoutMinutes: displayMins,
		SleepTimeoutMinutes:   sleepMins,
		IsNeverSleep:          isNever,
		Source:                "powercfg",
	}, nil
}

func (w *windowsDriver) queryDisplayTimeout() (int, error) {
	out, err := w.run("powercfg", "/q", "SCHEME_CURRENT", "SUB_VIDEO")
	if err != nil {
		return 0, apperror.WrapSimple(err, "power.queryDisplay")
	}

	sec, err := ParsePowercfgSettingIndex(string(out), "VIDEOIDLE")
	if err != nil {
		return 0, err
	}

	return ConvertSecondsToMinutes(sec), nil
}

func (w *windowsDriver) querySleepTimeout() (int, error) {
	out, err := w.run("powercfg", "/q", "SCHEME_CURRENT", "SUB_SLEEP")
	if err != nil {
		return 0, apperror.WrapSimple(err, "power.querySleep")
	}

	sec, err := ParsePowercfgSettingIndex(string(out), "STANDBYIDLE")
	if err != nil {
		return 0, err
	}

	return ConvertSecondsToMinutes(sec), nil
}

func (w *windowsDriver) SetNeverSleep() error {
	commands := [][]string{
		{"/change", "monitor-timeout-ac", "0"},
		{"/change", "monitor-timeout-dc", "0"},
		{"/change", "standby-timeout-ac", "0"},
		{"/change", "standby-timeout-dc", "0"},
		{"/change", "hibernate-timeout-ac", "0"},
		{"/change", "hibernate-timeout-dc", "0"},
	}

	return w.executePowercfgBatch(commands)
}

func (w *windowsDriver) SetTimeouts(displayMinutes, sleepMinutes int) error {
	commands := [][]string{
		{"/change", "monitor-timeout-ac", strconv.Itoa(displayMinutes)},
		{"/change", "monitor-timeout-dc", strconv.Itoa(displayMinutes)},
		{"/change", "standby-timeout-ac", strconv.Itoa(sleepMinutes)},
		{"/change", "standby-timeout-dc", strconv.Itoa(sleepMinutes)},
	}

	return w.executePowercfgBatch(commands)
}

func (w *windowsDriver) executePowercfgBatch(cmdArgs [][]string) error {
	for _, args := range cmdArgs {
		_, err := w.run("powercfg", args...)
		if err != nil {
			return apperror.WrapSimple(err, fmt.Sprintf("powercfg %v", args))
		}
	}

	return nil
}

func (w *windowsDriver) ApplySettings(s Settings) error {
	if s.IsNeverSleep {
		return w.SetNeverSleep()
	}

	return w.SetTimeouts(s.DisplayTimeoutMinutes, s.SleepTimeoutMinutes)
}
