package appwriter

import (
	"context"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/fileutil"
	"coding-guidelines/common/pkg/payloadconv"
)

// FileWriterOptions specifies configuration for file-backed writers.
type FileWriterOptions struct {
	Name     string
	FilePath string
	OpenMode fileutil.FileOpenModeType
	PermMode fileutil.FilePermType
	IsLocked bool
}

// NewFileWriter creates a file writer using fileutil enums and wrap constructors.
func NewFileWriter(opts FileWriterOptions) BaseWriterWrap {
	if len(opts.FilePath) == 0 {
		return WrapWriter.FailureWithId(errtype.Validation, "file path cannot be empty")
	}

	openMode := opts.OpenMode
	if openMode == 0 {
		openMode = fileutil.FileOpenCreateAppend
	}

	permMode := opts.PermMode
	if permMode == 0 {
		permMode = fileutil.FilePermStandard
	}

	fileWrap := fileutil.OpenFile(opts.FilePath, openMode, permMode)
	if fileWrap.IsFailed() {
		return WrapWriter.FailureFromWrap(fileWrap)
	}

	name := opts.Name
	if len(name) == 0 {
		name = opts.FilePath
	}

	writer := NewBaseWriter(name, fileWrap.Data(), opts.IsLocked, fileWriteFunc)

	return WrapWriter.Success(writer)
}

func fileWriteFunc(ctx context.Context, self Writer, payload any) *appfault.AppError {
	data := payloadconv.ToBytesMust(payload)

	_, err := self.Destination().Write(data)
	if err != nil {
		return appfault.Wrap(errtype.IO, err, "failed to write payload to file destination")
	}

	return nil
}
