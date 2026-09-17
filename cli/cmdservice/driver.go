package cmdservice

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ServiceDriverResolver defines the function type to resolve an OS service driver.
type ServiceDriverResolver func() ServiceDriver

// DefaultServiceDriverResolver is the active driver resolver, customizable for hermetic testing.
var DefaultServiceDriverResolver ServiceDriverResolver = ResolveServiceDriver

// EnsureServiceDriver returns the active driver or an AppError if unsupported.
func EnsureServiceDriver() (ServiceDriver, *apperror.AppError) {
	d := DefaultServiceDriverResolver()
	if d == nil {
		return nil, apperror.NewSimple("unsupported operating system for service management", "E_SERVICE_UNSUPPORTED_OS")
	}

	return d, nil
}
