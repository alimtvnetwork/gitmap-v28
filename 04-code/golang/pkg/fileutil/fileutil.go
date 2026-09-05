package fileutil

import (
	"os"
	"path/filepath"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/result"
)

func OpenFile(path string, openMode FileOpenModeType, perm FilePermType) result.Wrap[*os.File] {
	if len(path) == 0 {
		return result.WrapFailure[*os.File](appfault.New(errtype.Validation, "path cannot be empty"))
	}

	flags := openMode.Flags()
	isCreateMode := (flags & os.O_CREATE) != 0
	if isCreateMode {
		dir := filepath.Dir(path)
		if len(dir) > 0 {
			if dir != "." {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return result.WrapFailure[*os.File](appfault.Wrap(errtype.IO, err, "failed to create parent directory: "+dir))
				}
			}
		}
	}

	f, err := os.OpenFile(path, flags, perm.Mode())
	if err != nil {
		if os.IsNotExist(err) {
			return result.WrapFailure[*os.File](appfault.Wrap(errtype.NotFound, err, "file not found: "+path))
		}

		if os.IsPermission(err) {
			return result.WrapFailure[*os.File](appfault.Wrap(errtype.Forbidden, err, "permission denied: "+path))
		}

		return result.WrapFailure[*os.File](appfault.Wrap(errtype.IO, err, "failed to open file: "+path))
	}

	return result.WrapSuccess(f)
}

func Open(path string) result.Wrap[*os.File] {
	return OpenFile(path, FileOpenReadOnly, FilePermStandard)
}

func EnsureDir(path string, perm FilePermType) result.Wrap[bool] {
	if len(path) == 0 {
		return result.WrapFailure[bool](appfault.New(errtype.Validation, "directory path cannot be empty"))
	}

	err := os.MkdirAll(path, perm.Mode())
	if err != nil {
		return result.WrapFailure[bool](appfault.Wrap(errtype.IO, err, "failed to create directory: "+path))
	}

	return result.WrapSuccess(true)
}

func ReadAll(path string) result.Wrap[[]byte] {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return result.WrapFailure[[]byte](appfault.Wrap(errtype.NotFound, err, "file not found: "+path))
		}

		return result.WrapFailure[[]byte](appfault.Wrap(errtype.IO, err, "failed to read file: "+path))
	}

	return result.WrapSuccess(data)
}

func ReadString(path string) result.Wrap[string] {
	res := ReadAll(path)
	if res.IsFailed() {
		return result.WrapFailure[string](res.Fault())
	}

	return result.WrapSuccess(string(res.Data()))
}

func WriteFile(path string, data []byte, perm FilePermType) result.Wrap[bool] {
	wrap := OpenFile(path, FileOpenCreateTruncate, perm)
	if wrap.IsFailed() {
		return result.WrapFailure[bool](wrap.Fault())
	}

	f := wrap.Data()
	defer f.Close()

	_, err := f.Write(data)
	if err != nil {
		return result.WrapFailure[bool](appfault.Wrap(errtype.IO, err, "failed to write data: "+path))
	}

	return result.WrapSuccess(true)
}

func DeleteFile(path string) result.Wrap[bool] {
	if len(path) == 0 {
		return result.WrapFailure[bool](appfault.New(errtype.Validation, "path cannot be empty"))
	}

	err := os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return result.WrapFailure[bool](appfault.Wrap(errtype.NotFound, err, "file not found: "+path))
		}

		if os.IsPermission(err) {
			return result.WrapFailure[bool](appfault.Wrap(errtype.Forbidden, err, "permission denied: "+path))
		}

		return result.WrapFailure[bool](appfault.Wrap(errtype.IO, err, "failed to delete file: "+path))
	}

	return result.WrapSuccess(true)
}

func Remove(path string) result.Wrap[bool] {
	return DeleteFile(path)
}

func RemoveAll(path string) result.Wrap[bool] {
	if len(path) == 0 {
		return result.WrapFailure[bool](appfault.New(errtype.Validation, "path cannot be empty"))
	}

	err := os.RemoveAll(path)
	if err != nil {
		if os.IsPermission(err) {
			return result.WrapFailure[bool](appfault.Wrap(errtype.Forbidden, err, "permission denied: "+path))
		}

		return result.WrapFailure[bool](appfault.Wrap(errtype.IO, err, "failed to remove path: "+path))
	}

	return result.WrapSuccess(true)
}

func ReadFile(path string) result.Wrap[[]byte] {
	return ReadAll(path)
}

func Stat(path string) result.Wrap[os.FileInfo] {
	if len(path) == 0 {
		return result.WrapFailure[os.FileInfo](appfault.New(errtype.Validation, "path cannot be empty"))
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return result.WrapFailure[os.FileInfo](appfault.Wrap(errtype.NotFound, err, "file not found: "+path))
		}

		if os.IsPermission(err) {
			return result.WrapFailure[os.FileInfo](appfault.Wrap(errtype.Forbidden, err, "permission denied: "+path))
		}

		return result.WrapFailure[os.FileInfo](appfault.Wrap(errtype.IO, err, "failed to stat file: "+path))
	}

	return result.WrapSuccess(info)
}

func FileSize(path string) result.Wrap[int64] {
	statRes := Stat(path)
	if statRes.IsFailed() {
		return result.WrapFailure[int64](statRes.Fault())
	}

	return result.WrapSuccess(statRes.Data().Size())
}

func ExecuteOp(
	path string,
	op FileOpType,
	perm FilePermType,
	data []byte,
) result.Wrap[[]byte] {
	if op.IsDelete() {
		delRes := DeleteFile(path)
		if delRes.IsFailed() {
			return result.WrapFailure[[]byte](delRes.Fault())
		}

		return result.WrapSuccess[[]byte](nil)
	}

	if op.IsReadOnly() {
		return ReadAll(path)
	}

	openRes := OpenFile(path, op.OpenMode(), perm)
	if openRes.IsFailed() {
		return result.WrapFailure[[]byte](openRes.Fault())
	}

	file := openRes.Data()
	defer file.Close()

	if len(data) > 0 {
		_, writeErr := file.Write(data)
		if writeErr != nil {
			return result.WrapFailure[[]byte](appfault.Wrap(errtype.IO, writeErr, "failed to write during op: "+path))
		}
	}

	return result.WrapSuccess(data)
}
