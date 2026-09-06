package appwriter

import (
	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
)

type writerWrapConstructor struct{}

var (
	WrapWriter = writerWrapConstructor{}
	WriterWrap = WrapWriter
	Wrap       = WrapWriter
)

func (writerWrapConstructor) Success(w *BaseWriter) BaseWriterWrap {
	return result.WrapSuccess(w)
}

func (writerWrapConstructor) Failure(err *appfault.AppError) BaseWriterWrap {
	return result.WrapFailure[*BaseWriter](err)
}

func (writerWrapConstructor) FailureFromError(err *appfault.AppError) BaseWriterWrap {
	return result.WrapFailure[*BaseWriter](err)
}

func (writerWrapConstructor) FailureWithId(errType errtype.Variation, msg string) BaseWriterWrap {
	return result.WrapFailureWithId[*BaseWriter](errType, msg)
}

func (writerWrapConstructor) FailureWithCause(errType errtype.Variation, cause error, msg string) BaseWriterWrap {
	return result.WrapFailureWithCause[*BaseWriter](errType, cause, msg)
}

func (writerWrapConstructor) FailureFromWrap(failed any) BaseWriterWrap {
	if f, ok := failed.(interface{ Fault() *appfault.AppError }); ok {
		return result.WrapFailure[*BaseWriter](f.Fault())
	}

	if a, ok := failed.(interface{ AppError() *appfault.AppError }); ok {
		return result.WrapFailure[*BaseWriter](a.AppError())
	}

	return result.WrapFailure[*BaseWriter](appfault.New(errtype.Generic, "unknown wrap failure"))
}

func WrapWriterSuccess(w *BaseWriter) BaseWriterWrap {
	return WrapWriter.Success(w)
}

func WrapWriterFailure(err *appfault.AppError) BaseWriterWrap {
	return WrapWriter.Failure(err)
}

func WrapWriterFailureFromError(err *appfault.AppError) BaseWriterWrap {
	return WrapWriter.FailureFromError(err)
}

func WrapWriterFailureWithId(errType errtype.Variation, msg string) BaseWriterWrap {
	return WrapWriter.FailureWithId(errType, msg)
}

func WrapWriterFailureWithCause(errType errtype.Variation, cause error, msg string) BaseWriterWrap {
	return WrapWriter.FailureWithCause(errType, cause, msg)
}

func WrapWriterFailureFromWrap[U any](failed result.Wrap[U]) BaseWriterWrap {
	return WrapWriter.Failure(failed.Fault())
}
