package power

import "fmt"

// Settings represents the operating system power and timeout configuration.
type Settings struct {
	Platform              string `json:"platform"`
	DisplayTimeoutMinutes int    `json:"displayTimeoutMinutes"`
	SleepTimeoutMinutes   int    `json:"sleepTimeoutMinutes"`
	DiskTimeoutMinutes    int    `json:"diskTimeoutMinutes"`
	IsNeverSleep          bool   `json:"isNeverSleep"`
	IsLockDisabled        bool   `json:"isLockDisabled"`
	Source                string `json:"source"`
}

// Summary returns a human-readable description of current power settings.
func (s Settings) Summary() string {
	state := "active"
	if s.IsNeverSleep {
		state = "never-sleep (inhibited)"
	}

	return fmt.Sprintf("display: %dm, sleep: %dm, state: %s",
		s.DisplayTimeoutMinutes, s.SleepTimeoutMinutes, state)
}

// Validate checks whether settings contain valid timeout values.
func (s Settings) Validate() error {
	if s.DisplayTimeoutMinutes < 0 || s.SleepTimeoutMinutes < 0 {
		return fmt.Errorf("timeout minutes cannot be negative")
	}

	return nil
}
