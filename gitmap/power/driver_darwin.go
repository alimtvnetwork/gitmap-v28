//go:build darwin

package power

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
	return Settings{
		Platform:              "darwin",
		DisplayTimeoutMinutes: 0,
		SleepTimeoutMinutes:   0,
		IsNeverSleep:          false,
		Source:                "pmset",
	}, nil
}

func (d *darwinDriver) SetNeverSleep() error {
	_, err := d.run("pmset", "-a", "displaysleep", "0", "sleep", "0")

	return err
}

func (d *darwinDriver) SetTimeouts(displayMinutes, sleepMinutes int) error {
	_, err := d.run("pmset", "-a", "displaysleep", string(rune(displayMinutes)), "sleep", string(rune(sleepMinutes)))

	return err
}

func (d *darwinDriver) ApplySettings(s Settings) error {
	if s.IsNeverSleep {
		return d.SetNeverSleep()
	}

	return d.SetTimeouts(s.DisplayTimeoutMinutes, s.SleepTimeoutMinutes)
}
