package fileutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
	"coding-guidelines/common/pkg/streamwriter"
)

var chunkBufferPool = sync.Pool{
	New: func() any {
		b := make([]byte, DefaultBufferSize)

		return &b
	},
}

func WriteAtomic(path string, data []byte, perm FilePermType) result.Wrap[bool] {
	if len(path) == 0 {
		return result.WrapFailure[bool](appfault.New(errtype.Validation, "path cannot be empty"))
	}

	dir := filepath.Dir(path)
	ensureRes := EnsureDir(dir, FilePermStandard)
	if ensureRes.IsFailed() {
		return result.WrapFailure[bool](ensureRes.Fault())
	}

	tmpPattern := fmt.Sprintf(".%s.tmp.%d.%d", filepath.Base(path), os.Getpid(), time.Now().UnixNano())
	tmpPath := filepath.Join(dir, tmpPattern)

	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, perm.Mode())
	if err != nil {
		return result.WrapFailure[bool](appfault.Wrap(errtype.IO, err, "failed to create atomic temp file: "+tmpPath))
	}

	writeFailed := false
	if len(data) > 0 {
		if _, writeErr := tmpFile.Write(data); writeErr != nil {
			writeFailed = true
			_ = tmpFile.Close()
			_ = os.Remove(tmpPath)

			return result.WrapFailure[bool](appfault.Wrap(errtype.IO, writeErr, "failed to write atomic data: "+tmpPath))
		}
	}

	if writeFailed {
		return result.WrapFailure[bool](appfault.New(errtype.IO, "atomic write aborted"))
	}

	if syncErr := tmpFile.Sync(); syncErr != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)

		return result.WrapFailure[bool](appfault.Wrap(errtype.IO, syncErr, "failed to sync atomic file: "+tmpPath))
	}

	if closeErr := tmpFile.Close(); closeErr != nil {
		_ = os.Remove(tmpPath)

		return result.WrapFailure[bool](appfault.Wrap(errtype.IO, closeErr, "failed to close atomic file: "+tmpPath))
	}

	// On Windows, rename fails if destination exists, so we remove destination first
	_ = os.Remove(path)

	if renameErr := os.Rename(tmpPath, path); renameErr != nil {
		_ = os.Remove(tmpPath)

		return result.WrapFailure[bool](appfault.Wrap(errtype.IO, renameErr, "failed to atomically rename: "+path))
	}

	return result.WrapSuccess(true)
}

func ReadChunked(path string, chunkSize int, onChunk ChunkCallbackFunc) result.Wrap[int64] {
	if len(path) == 0 {
		return result.WrapFailure[int64](appfault.New(errtype.Validation, "path cannot be empty"))
	}

	effectiveChunkSize := chunkSize
	if effectiveChunkSize <= 0 {
		effectiveChunkSize = DefaultBufferSize
	}

	openRes := Open(path)
	if openRes.IsFailed() {
		return result.WrapFailure[int64](openRes.Fault())
	}

	f := openRes.Data()
	defer f.Close()

	bufPtr := chunkBufferPool.Get().(*[]byte)
	buf := *bufPtr
	if len(buf) < effectiveChunkSize {
		buf = make([]byte, effectiveChunkSize)
	}

	defer chunkBufferPool.Put(&buf)

	var totalBytes int64
	for {
		n, readErr := f.Read(buf[:effectiveChunkSize])
		if n > 0 {
			totalBytes += int64(n)
			if onChunk != nil {
				chunkFault := onChunk(buf[:n])
				if chunkFault != nil {
					return result.WrapFailure[int64](chunkFault)
				}
			}
		}

		if readErr != nil {
			if readErr == io.EOF {
				break
			}

			return result.WrapFailure[int64](appfault.Wrap(errtype.IO, readErr, "error reading chunk from: "+path))
		}
	}

	return result.WrapSuccess(totalBytes)
}

func WriteChunked(path string, perm FilePermType, reader io.Reader, bufferSize int) result.Wrap[int64] {
	if len(path) == 0 {
		return result.WrapFailure[int64](appfault.New(errtype.Validation, "path cannot be empty"))
	}

	if reader == nil {
		return result.WrapFailure[int64](appfault.New(errtype.Validation, "reader cannot be nil"))
	}

	effectiveBufSize := bufferSize
	if effectiveBufSize <= 0 {
		effectiveBufSize = DefaultBufferSize
	}

	openRes := OpenFile(path, FileOpenCreateTruncate, perm)
	if openRes.IsFailed() {
		return result.WrapFailure[int64](openRes.Fault())
	}

	f := openRes.Data()
	defer f.Close()

	bufPtr := chunkBufferPool.Get().(*[]byte)
	buf := *bufPtr
	if len(buf) < effectiveBufSize {
		buf = make([]byte, effectiveBufSize)
	}

	defer chunkBufferPool.Put(&buf)

	var totalWritten int64
	for {
		n, readErr := reader.Read(buf[:effectiveBufSize])
		if n > 0 {
			written, writeErr := f.Write(buf[:n])
			totalWritten += int64(written)
			if writeErr != nil {
				return result.WrapFailure[int64](appfault.Wrap(errtype.IO, writeErr, "failed writing chunk to: "+path))
			}
		}

		if readErr != nil {
			if readErr == io.EOF {
				break
			}

			return result.WrapFailure[int64](appfault.Wrap(errtype.IO, readErr, "error reading source for write to: "+path))
		}
	}

	if syncErr := f.Sync(); syncErr != nil {
		return result.WrapFailure[int64](appfault.Wrap(errtype.IO, syncErr, "failed syncing file: "+path))
	}

	return result.WrapSuccess(totalWritten)
}

func NewFileWriter(path string, openMode FileOpenModeType, perm FilePermType) result.Wrap[*streamwriter.PluggableWriter[any]] {
	openRes := OpenFile(path, openMode, perm)
	if openRes.IsFailed() {
		return result.WrapFailure[*streamwriter.PluggableWriter[any]](openRes.Fault())
	}

	file := openRes.Data()
	writer := streamwriter.NewAnyWriter(streamwriter.WriterOptions[any]{
		Name:        "file-writer:" + filepath.Base(path),
		Destination: file,
	})

	return result.WrapSuccess(writer)
}
