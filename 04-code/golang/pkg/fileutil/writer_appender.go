package fileutil

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/streamwriter"
)

type FileWriteModeType uint8

var writeModeNames = map[FileWriteModeType]string{
	FileWriteModeDirect:   "Direct",
	FileWriteModeAtomic:   "Atomic",
	FileWriteModeTruncate: "Truncate",
}

func (m FileWriteModeType) Name() string {
	if name, ok := writeModeNames[m]; ok {
		return name
	}

	return fmt.Sprintf("FileWriteMode(%d)", uint8(m))
}

func (m FileWriteModeType) String() string {
	return m.Name()
}

func (m FileWriteModeType) IsValid() bool {
	_, ok := writeModeNames[m]

	return ok
}

type (
	FileWriterOptions struct {
		Path        string
		Mode        FileWriteModeType
		Perm        FilePermType
		SyncOnWrite bool
	}

	FileWriter struct {
		mu          sync.RWMutex
		path        string
		mode        FileWriteModeType
		perm        FilePermType
		syncOnWrite bool
		file        *os.File
	}
)

func NewFileWriterEngine(path string) *FileWriter {
	return &FileWriter{
		path:        path,
		mode:        FileWriteModeDirect,
		perm:        FilePermStandard,
		syncOnWrite: false,
	}
}

func NewFileWriterWithOptions(opts FileWriterOptions) *FileWriter {
	perm := opts.Perm
	if perm == 0 {
		perm = FilePermStandard
	}

	mode := opts.Mode
	if mode == 0 {
		mode = FileWriteModeDirect
	}

	return &FileWriter{
		path:        opts.Path,
		mode:        mode,
		perm:        perm,
		syncOnWrite: opts.SyncOnWrite,
	}
}

func (w *FileWriter) Name() string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return fmt.Sprintf("file-writer[%s]", filepath.Base(w.path))
}

func (w *FileWriter) Path() string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.path
}

func (w *FileWriter) Mode() FileWriteModeType {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.mode
}

func (w *FileWriter) SetMode(mode FileWriteModeType) *FileWriter {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.mode = mode

	return w
}

func (w *FileWriter) SetPerm(perm FilePermType) *FileWriter {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.perm = perm

	return w
}

func (w *FileWriter) SetSyncOnWrite(isSync bool) *FileWriter {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.syncOnWrite = isSync

	return w
}

func (w *FileWriter) Write(ctx context.Context, payload []byte) *appfault.AppError {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.path == "" {
		return appfault.New(errtype.Precondition, "file path cannot be empty")
	}

	switch w.mode {
	case FileWriteModeAtomic:
		res := WriteAtomic(w.path, payload, w.perm)
		if res.IsFailed() {
			return res.Fault()
		}

		return nil

	case FileWriteModeTruncate:
		return w.writeDirect(payload, os.O_CREATE|os.O_TRUNC|os.O_WRONLY)

	default:
		return w.writeDirect(payload, os.O_CREATE|os.O_WRONLY)
	}
}

func (w *FileWriter) writeDirect(payload []byte, flags int) *appfault.AppError {
	dir := filepath.Dir(w.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return appfault.Wrap(errtype.IO, err, "failed to create parent directories")
	}

	f, err := os.OpenFile(w.path, flags, w.perm.Mode())
	if err != nil {
		return appfault.Wrap(errtype.IO, err, "failed to open target file")
	}

	defer f.Close()

	if _, err := f.Write(payload); err != nil {
		return appfault.Wrap(errtype.IO, err, "failed to write payload to file")
	}

	if w.syncOnWrite {
		if err := f.Sync(); err != nil {
			return appfault.Wrap(errtype.IO, err, "failed to sync file to storage")
		}
	}

	return nil
}

func (w *FileWriter) WriteString(ctx context.Context, text string) *appfault.AppError {
	return w.Write(ctx, []byte(text))
}

func (w *FileWriter) WriteStd(p []byte) (n int, err error) {
	appErr := w.Write(context.Background(), p)
	if appErr != nil {
		return 0, appErr
	}

	return len(p), nil
}

func (w *FileWriter) StdWriter() io.WriteCloser {
	return &fileWriterStdAdapter{writer: w}
}

type fileWriterStdAdapter struct {
	writer *FileWriter
}

func (s *fileWriterStdAdapter) Write(p []byte) (n int, err error) {
	return s.writer.WriteStd(p)
}

func (s *fileWriterStdAdapter) Close() error {
	appErr := s.writer.Close()
	if appErr != nil {
		return appErr
	}

	return nil
}

func (w *FileWriter) Sync() *appfault.AppError {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.file != nil {
		if err := w.file.Sync(); err != nil {
			return appfault.Wrap(errtype.IO, err, "sync failed")
		}
	}

	return nil
}

func (w *FileWriter) Close() *appfault.AppError {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		if err != nil {
			return appfault.Wrap(errtype.IO, err, "close failed")
		}
	}

	return nil
}

func (w *FileWriter) Lock() {
	w.mu.Lock()
}

func (w *FileWriter) Unlock() {
	w.mu.Unlock()
}

func (w *FileWriter) AsWriter() streamwriter.Writer[[]byte] {
	return w
}

type FileAppender struct {
	mu            sync.Mutex
	path          string
	perm          FilePermType
	autoSync      bool
	file          *os.File
	bytesAppended atomic.Int64
}

func NewFileAppender(path string, perm FilePermType) *FileAppender {
	if perm == 0 {
		perm = FilePermStandard
	}

	return &FileAppender{
		path:     path,
		perm:     perm,
		autoSync: false,
	}
}

func (a *FileAppender) Path() string {
	return a.path
}

func (a *FileAppender) SetAutoSync(isAuto bool) *FileAppender {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.autoSync = isAuto

	return a
}

func (a *FileAppender) ensureOpen() error {
	if a.file != nil {
		return nil
	}

	dir := filepath.Dir(a.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(a.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, a.perm.Mode())
	if err != nil {
		return err
	}

	a.file = f

	return nil
}

func (a *FileAppender) Append(ctx context.Context, data []byte) *appfault.AppError {
	a.mu.Lock()
	defer a.mu.Unlock()

	if err := a.ensureOpen(); err != nil {
		return appfault.Wrap(errtype.IO, err, "failed to open appender target file")
	}

	n, err := a.file.Write(data)
	if err != nil {
		return appfault.Wrap(errtype.IO, err, "failed to append bytes to file")
	}

	a.bytesAppended.Add(int64(n))

	if a.autoSync {
		if err := a.file.Sync(); err != nil {
			return appfault.Wrap(errtype.IO, err, "failed to sync appender file")
		}
	}

	return nil
}

func (a *FileAppender) AppendString(ctx context.Context, text string) *appfault.AppError {
	return a.Append(ctx, []byte(text))
}

func (a *FileAppender) WriteStd(p []byte) (n int, err error) {
	appErr := a.Append(context.Background(), p)
	if appErr != nil {
		return 0, appErr
	}

	return len(p), nil
}

func (a *FileAppender) StdWriter() io.WriteCloser {
	return &appenderStdAdapter{appender: a}
}

type appenderStdAdapter struct {
	appender *FileAppender
}

func (s *appenderStdAdapter) Write(p []byte) (n int, err error) {
	return s.appender.WriteStd(p)
}

func (s *appenderStdAdapter) Close() error {
	appErr := s.appender.Close()
	if appErr != nil {
		return appErr
	}

	return nil
}

func (a *FileAppender) BytesAppended() int64 {
	return a.bytesAppended.Load()
}

func (a *FileAppender) Sync() *appfault.AppError {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.file != nil {
		if err := a.file.Sync(); err != nil {
			return appfault.Wrap(errtype.IO, err, "appender sync failed")
		}
	}

	return nil
}

func (a *FileAppender) Close() *appfault.AppError {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.file != nil {
		err := a.file.Close()
		a.file = nil
		if err != nil {
			return appfault.Wrap(errtype.IO, err, "appender close failed")
		}
	}

	return nil
}

func (a *FileAppender) Lock() {
	a.mu.Lock()
}

func (a *FileAppender) Unlock() {
	a.mu.Unlock()
}

func (a *FileAppender) Name() string {
	return fmt.Sprintf("file-appender[%s]", filepath.Base(a.path))
}

func (a *FileAppender) Write(ctx context.Context, payload []byte) *appfault.AppError {
	return a.Append(ctx, payload)
}

func (a *FileAppender) AsWriter() streamwriter.Writer[[]byte] {
	return a
}

var _ streamwriter.Writer[[]byte] = (*FileWriter)(nil)
var _ sync.Locker = (*FileWriter)(nil)
var _ streamwriter.Writer[[]byte] = (*FileAppender)(nil)
var _ sync.Locker = (*FileAppender)(nil)
var _ io.WriteCloser = (*fileWriterStdAdapter)(nil)
var _ io.WriteCloser = (*appenderStdAdapter)(nil)
