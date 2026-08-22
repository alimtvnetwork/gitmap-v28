package result

type Result[T any] struct {
	IsSuccess bool
	IsFailed  bool
	Data      T
	AppError  error
}

func NewSuccess[T any](data T) Result[T] {
	return Result[T]{
		IsSuccess: true,
		IsFailed:  false,
		Data:      data,
	}
}

func NewFailure[T any](err error) Result[T] {
	return Result[T]{
		IsSuccess: false,
		IsFailed:  true,
		AppError:  err,
	}
}
