package cmdos

// defaultTweakEngine holds the active tweak implementation.
var defaultTweakEngine TweakOperator = newPlatformTweakEngine()

// GetTweakEngine returns the active tweak engine.
func GetTweakEngine() TweakOperator {
	return defaultTweakEngine
}

// SetTweakEngine overrides the active tweak engine (for tests).
func SetTweakEngine(engine TweakOperator) {
	defaultTweakEngine = engine
}
