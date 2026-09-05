package fileutil

import "os"

// FileActionType defines the target operation type for file manipulation.
type FileActionType string

const (
	FileActionTypeReadOnly  FileActionType = "READ_ONLY"
	FileActionTypeRead      FileActionType = "READ"
	FileActionTypeWrite     FileActionType = "WRITE"
	FileActionTypeAppend    FileActionType = "APPEND"
	FileActionTypeCreate    FileActionType = "CREATE"
	FileActionTypeDelete    FileActionType = "DELETE"
	FileActionTypeReadWrite FileActionType = "READ_WRITE"
	FileActionTypeTruncate  FileActionType = "TRUNCATE"
)

// String implements fmt.Stringer.
func (a FileActionType) String() string {
	return string(a)
}

// IsRead reports whether the action involves reading.
func (a FileActionType) IsRead() bool {
	return a == FileActionTypeRead || a == FileActionTypeReadOnly || a == FileActionTypeReadWrite
}

// IsWrite reports whether the action involves writing.
func (a FileActionType) IsWrite() bool {
	return a == FileActionTypeWrite || a == FileActionTypeAppend || a == FileActionTypeTruncate || a == FileActionTypeCreate || a == FileActionTypeReadWrite
}

// IsAppend reports whether the action appends to an existing file.
func (a FileActionType) IsAppend() bool {
	return a == FileActionTypeAppend
}

// IsCreate reports whether the action requires file creation.
func (a FileActionType) IsCreate() bool {
	return a == FileActionTypeCreate
}

// IsDelete reports whether the action deletes the file.
func (a FileActionType) IsDelete() bool {
	return a == FileActionTypeDelete
}

// ToOSFlags maps FileActionType to underlying os.OpenFile integer flags.
func (a FileActionType) ToOSFlags() int {
	switch a {
	case FileActionTypeReadOnly, FileActionTypeRead:
		return os.O_RDONLY
	case FileActionTypeAppend:
		return os.O_CREATE | os.O_WRONLY | os.O_APPEND
	case FileActionTypeTruncate, FileActionTypeWrite:
		return os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	case FileActionTypeCreate:
		return os.O_CREATE | os.O_EXCL | os.O_WRONLY
	case FileActionTypeReadWrite:
		return os.O_CREATE | os.O_RDWR
	default:
		return os.O_RDONLY
	}
}
