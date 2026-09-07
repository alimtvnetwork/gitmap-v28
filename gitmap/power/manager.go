package power

import (
	"os/exec"
	"runtime"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// Manager abstracts OS-specific power and screen timeout operations.
type Manager interface {
	Platform() string
	GetStatus() (Settings, error)
	SetNeverSleep() error
	SetTimeouts(displayMinutes, sleepMinutes int) error
	ApplySettings(s Settings) error
}

// CmdRunner represents a function that executes an external process.
type CmdRunner func(name string, args ...string) ([]byte, error)

var (
	registryMu sync.RWMutex
	drivers              = make(map[string]func(runner CmdRunner) Manager)
	runner     CmdRunner = defaultExecRunner
)

func defaultExecRunner(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)

	return cmd.CombinedOutput()
}

// SetRunnerForTesting sets a custom process runner for tests.
func SetRunnerForTesting(r CmdRunner) func() {
	prev := runner
	runner = r

	return func() {
		runner = prev
	}
}

// RegisterDriver registers an OS power manager factory.
func RegisterDriver(platform string, factory func(runner CmdRunner) Manager) {
	registryMu.Lock()
	defer registryMu.Unlock()
	drivers[platform] = factory
}

// NewManager creates a power Manager for the current host operating system.
func NewManager() (Manager, error) {
	registryMu.RLock()
	factory, exists := drivers[runtime.GOOS]
	registryMu.RUnlock()

	if !exists {
		return nil, apperror.NewSimple("unsupported power management platform: "+runtime.GOOS, "E_UNSUPPORTED_OS")
	}

	return factory(runner), nil
}
