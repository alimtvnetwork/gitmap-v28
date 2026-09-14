//go:build !windows && !linux && !darwin

package cmdservice

// ResolveServiceDriver returns a fallback service driver for unsupported platforms.
func ResolveServiceDriver() ServiceDriver {
	return nil
}
