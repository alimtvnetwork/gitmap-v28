package result

const (
	DefaultSuccessBanner = "✅ [OK]"
	DefaultFailureBanner = "❌ [FAIL]"
)

type ResultFormatter[T any] func(r Wrap[T]) string

type ResultPredicate[T any] func(r Wrap[T]) bool

type ResultMapper[T any, U any] func(val T) U
