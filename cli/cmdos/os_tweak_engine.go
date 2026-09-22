package cmdos

// defaultTweakEngine holds the active tweak implementation.
var defaultTweakEngine TweakEngine = newPlatformTweakEngine()

// GetTweakEngine returns the active tweak engine.
func GetTweakEngine() TweakEngine {
	return defaultTweakEngine
}

// SetTweakEngine overrides the active tweak engine (for tests).
func SetTweakEngine(engine TweakEngine) {
	defaultTweakEngine = engine
}
