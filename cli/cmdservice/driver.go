package cmdservice

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// EnsureServiceDriver returns the active driver or an AppError if unsupported.
func EnsureServiceDriver() (ServiceDriver, *apperror.AppError) {
	d := ResolveServiceDriver()
	if d == nil {
		return nil, apperror.NewSimple("unsupported operating system for service management", "E_SERVICE_UNSUPPORTED_OS")
	}

	return d, nil
}
