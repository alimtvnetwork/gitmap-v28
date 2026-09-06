package fileutil

import (
	"os"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/enum/openfiletype"
)

// Default I/O chunk and buffer size constants (64 KB).
const (
	DefaultChunkSize  = 64 * 1024
	DefaultBufferSize = DefaultChunkSize
)

const (
	FileOpenInvalid               FileOpenModeType = openfiletype.Invalid
	FileOpenReadOnly              FileOpenModeType = openfiletype.ReadOnly
	FileOpenWriteOnly             FileOpenModeType = openfiletype.WriteOnly
	FileOpenReadWrite             FileOpenModeType = openfiletype.ReadWrite
	FileOpenAppend                FileOpenModeType = openfiletype.Append
	FileOpenCreateAppend          FileOpenModeType = openfiletype.CreateAppend
	FileOpenCreateTruncate        FileOpenModeType = openfiletype.CreateTruncate
	FileOpenCreateNew             FileOpenModeType = openfiletype.CreateNew
	FileOpenReadOrCreateOnly      FileOpenModeType = openfiletype.ReadOrCreateOnly
	FileOpenWriteOrCreateOnly     FileOpenModeType = openfiletype.WriteOrCreateOnly
	FileOpenReadWriteOrCreateOnly FileOpenModeType = openfiletype.ReadWriteOrCreateOnly
)

const (
	FileOpReadOnly FileOpType = iota
	FileOpWriteOnly
	FileOpReadWrite
	FileOpAppend
	FileOpCreate
	FileOpCreateAppend
	FileOpCreateTruncate
	FileOpDelete
)

const (
	FilePermNone                 FilePermType = 0000
	FilePermOwnerReadOnly        FilePermType = 0400
	FilePermOwnerWriteOnly       FilePermType = 0200
	FilePermOwnerExecOnly        FilePermType = 0100
	FilePermOwnerReadWrite       FilePermType = 0600
	FilePermPrivate              FilePermType = 0600
	FilePermOwnerAll             FilePermType = 0700
	FilePermOwnerExec            FilePermType = 0700
	FilePermGroupReadOnly        FilePermType = 0440
	FilePermGroupWriteOnly       FilePermType = 0220
	FilePermGroupReadWrite       FilePermType = 0660
	FilePermGroupExec            FilePermType = 0750
	FilePermGroupAll             FilePermType = 0770
	FilePermReadOnly             FilePermType = 0444
	FilePermPublicReadOnly       FilePermType = 0444
	FilePermPublicWriteOnly      FilePermType = 0222
	FilePermStandard             FilePermType = 0644
	FilePermGroupSharedOtherRead FilePermType = 0664
	FilePermPublicReadWrite      FilePermType = 0666
	FilePermExecutable           FilePermType = 0755
	FilePermGroupSharedDir       FilePermType = 0775
	FilePermPublicAll            FilePermType = 0777
	FilePermStickyDir            FilePermType = 01777
	FilePermSetuidExec           FilePermType = 04755
	FilePermSetgidExec           FilePermType = 02755
)

const (
	FileWriteModeDirect   FileWriteModeType = 1
	FileWriteModeAtomic   FileWriteModeType = 2
	FileWriteModeTruncate FileWriteModeType = 3
)

type (
	ChunkHandlerFunc func(chunk []byte) error

	ChunkCallbackFunc func(chunk []byte) *appfault.AppError

	BoundFileActionFunc func(w *BoundFileWriter) *appfault.AppError

	WithLockFunc = BoundFileActionFunc

	FileFilterFunc func(path string, info os.FileInfo) bool
)
