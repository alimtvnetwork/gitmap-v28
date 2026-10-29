package setup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// CreatePreFlightBackup creates a timestamped copy of targetPath if it exists.
func CreatePreFlightBackup(targetPath string) (BackupResult, *apperror.AppError) {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return BackupResult{OriginalPath: targetPath, IsBackupCreated: false}, nil
	}

	timestamp := time.Now().UTC().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.bak.%s", targetPath, timestamp)
	if err := copyBackupFile(targetPath, backupPath); err != nil {
		return BackupResult{}, err
	}

	return BackupResult{OriginalPath: targetPath, BackupPath: backupPath, IsBackupCreated: true}, nil
}

func copyBackupFile(srcPath, backupPath string) *apperror.AppError {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return apperror.WrapExecution(err, "CreatePreFlightBackup.ReadFile")
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return apperror.WrapExecution(err, "CreatePreFlightBackup.WriteFile")
	}

	return nil
}

// MergeJSONFile reads existing JSON, deep merges updates, and writes back atomically.
func MergeJSONFile(targetPath string, updates map[string]interface{}) *apperror.AppError {
	existingData, err := loadExistingJSON(targetPath)
	if err != nil {
		return err
	}

	DeepMergeMaps(existingData, updates)
	outBytes, marshalErr := json.MarshalIndent(existingData, "", "  ")
	if marshalErr != nil {
		return apperror.WrapExecution(marshalErr, "MergeJSONFile.Marshal")
	}

	return AtomicWriteFile(targetPath, outBytes)
}

func loadExistingJSON(targetPath string) (map[string]interface{}, *apperror.AppError) {
	existingData := make(map[string]interface{})
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return existingData, nil
	}

	readData, readErr := os.ReadFile(targetPath)
	if readErr != nil {
		return nil, apperror.WrapExecution(readErr, "loadExistingJSON.ReadFile")
	}

	_ = json.Unmarshal(readData, &existingData)

	return existingData, nil
}

// DeepMergeMaps recursively overlays keys from src into dest without clobbering existing maps.
func DeepMergeMaps(dest map[string]interface{}, src map[string]interface{}) {
	for key, srcVal := range src {
		mergeSingleKey(dest, key, srcVal)
	}
}

func mergeSingleKey(dest map[string]interface{}, key string, srcVal any) {
	destVal, exists := dest[key]
	if !exists {
		dest[key] = srcVal

		return
	}

	if mergeNestedMaps(dest, key, destVal, srcVal) {
		return
	}

	dest[key] = srcVal
}

func mergeNestedMaps(dest map[string]interface{}, key string, destVal any, srcVal any) bool {
	destMap, isDestMap := destVal.(map[string]interface{})
	srcMap, isSrcMap := srcVal.(map[string]interface{})
	if !isDestMap || !isSrcMap {
		return false
	}

	DeepMergeMaps(destMap, srcMap)

	return true
}

// AtomicWriteFile ensures the parent directory exists, writes to a temp file, and renames atomically.
func AtomicWriteFile(targetPath string, data []byte) *apperror.AppError {
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return apperror.WrapExecution(err, "AtomicWriteFile.MkdirAll")
	}

	tempFile := fmt.Sprintf("%s.tmp.%d", targetPath, os.Getpid())
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return apperror.WrapExecution(err, "AtomicWriteFile.WriteTemp")
	}

	return commitAtomicFile(tempFile, targetPath)
}

func commitAtomicFile(tempFile, targetPath string) *apperror.AppError {
	if err := os.Rename(tempFile, targetPath); err != nil {
		_ = os.Remove(tempFile)

		return apperror.WrapExecution(err, "AtomicWriteFile.Rename")
	}

	return nil
}
