package cmdos

// AutoLoginConfig encapsulates credentials and configuration for OS auto-login.
type AutoLoginConfig struct {
	Username  string
	Domain    string
	Password  string
	IsEnabled bool
}

// AutoLoginStatus represents current auto-login state inspection.
type AutoLoginStatus struct {
	IsEnabled      bool
	Username       string
	Domain         string
	HasPassword    bool
	DisplayManager string
}

// AutoLoginEngine abstracts OS-specific auto-login read/write operations.
type AutoLoginEngine interface {
	Configure(cfg AutoLoginConfig) error
	Disable() error
	Status() (AutoLoginStatus, error)
}
