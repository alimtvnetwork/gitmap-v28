package fileutil

import "coding-guidelines/common/pkg/appfault"

var defaultWrapper = NewDefault()

// ReadFile reads a file using the global default wrapper.
func ReadFile(filePath string, action FileActionType, mode FileModeType) ([]byte, *appfault.AppError) {
	return defaultWrapper.Read(filePath, action, mode)
}

// ReadFileString reads file contents as string using default wrapper.
func ReadFileString(filePath string, action FileActionType, mode FileModeType) (string, *appfault.AppError) {
	return defaultWrapper.ReadString(filePath, action, mode)
}

// WriteFile writes data to a file using the global default wrapper.
func WriteFile(filePath string, data []byte, action FileActionType, mode FileModeType) *appfault.AppError {
	return defaultWrapper.Write(filePath, data, action, mode)
}

// AppendFile appends data to a file using default wrapper.
func AppendFile(filePath string, data []byte, mode FileModeType) *appfault.AppError {
	return defaultWrapper.Append(filePath, data, mode)
}

// DeleteFile deletes a file using default wrapper.
func DeleteFile(filePath string) *appfault.AppError {
	return defaultWrapper.Delete(filePath)
}
