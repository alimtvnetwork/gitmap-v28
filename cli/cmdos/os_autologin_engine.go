package cmdos

// defaultAutoLoginEngine holds the active engine implementation.
var defaultAutoLoginEngine AutoLoginEngine = newPlatformAutoLoginEngine()

// GetAutoLoginEngine returns the active auto-login engine.
func GetAutoLoginEngine() AutoLoginEngine {
	return defaultAutoLoginEngine
}

// SetAutoLoginEngine overrides the active engine (used in unit tests).
func SetAutoLoginEngine(engine AutoLoginEngine) {
	defaultAutoLoginEngine = engine
}
