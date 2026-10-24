package enums

// CloneStrategyType represents the clone execution strategy.
type CloneStrategyType string

const (
	StrategyBatch       CloneStrategyType = "batch"
	StrategySequential  CloneStrategyType = "sequential"
	StrategyInteractive CloneStrategyType = "interactive"
	StrategyParallel    CloneStrategyType = "parallel"
)

// IsValid reports whether the clone strategy is recognized.
func (s CloneStrategyType) IsValid() bool {
	switch s {
	case StrategyBatch, StrategySequential, StrategyInteractive, StrategyParallel:
		return true
	default:
		return false
	}
}

// String returns the string representation.
func (s CloneStrategyType) String() string {
	return string(s)
}

// IsBatch reports whether the strategy is batch.
func (s CloneStrategyType) IsBatch() bool {
	return s == StrategyBatch
}

// IsSequential reports whether the strategy is sequential.
func (s CloneStrategyType) IsSequential() bool {
	return s == StrategySequential
}

// IsInteractive reports whether the strategy is interactive.
func (s CloneStrategyType) IsInteractive() bool {
	return s == StrategyInteractive
}

// IsParallel reports whether the strategy is parallel.
func (s CloneStrategyType) IsParallel() bool {
	return s == StrategyParallel
}
