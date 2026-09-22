package cmdos

// defaultAutoLoginEngine holds the active engine implementation.
var defaultAutoLoginEngine AutoLoginOperator = newPlatformAutoLoginEngine()

// GetAutoLoginEngine returns the active auto-login engine.
func GetAutoLoginEngine() AutoLoginOperator {
	return defaultAutoLoginEngine
}

// SetAutoLoginEngine overrides the active engine (used in unit tests).
func SetAutoLoginEngine(engine AutoLoginOperator) {
	defaultAutoLoginEngine = engine
}
