//go:build darwin

package power

import "strconv"

type darwinDriver struct {
	run CmdRunner
}

func init() {
	RegisterDriver("darwin", func(r CmdRunner) Manager {
		return &darwinDriver{run: r}
	})
}

func (d *darwinDriver) Platform() string {
	return "darwin"
}

func (d *darwinDriver) GetStatus() (Settings, error) {
	out, err := d.run("pmset", "-g")
	if err != nil {
		return Settings{Platform: "darwin", Source: "pmset"}, err
	}

	return ParsePmsetSettings(string(out)), nil
}

func (d *darwinDriver) SetNeverSleep() error {
	_, err := d.run("pmset", "-a", "displaysleep", "0", "sleep", "0")

	return err
}

func (d *darwinDriver) SetTimeouts(displayMinutes, sleepMinutes int) error {
	dStr := strconv.Itoa(displayMinutes)
	sStr := strconv.Itoa(sleepMinutes)
	_, err := d.run("pmset", "-a", "displaysleep", dStr, "sleep", sStr)

	return err
}

func (d *darwinDriver) ApplySettings(s Settings) error {
	if s.IsNeverSleep {
		return d.SetNeverSleep()
	}

	return d.SetTimeouts(s.DisplayTimeoutMinutes, s.SleepTimeoutMinutes)
}
