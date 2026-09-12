package fileutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDirAndFile(t *testing.T) {
	tmp := t.TempDir()

	dirPath := filepath.Join(tmp, "sub")
	resDir := CreateDir(dirPath, FilePermStandard)
	if resDir.HasError() {
		t.Fatalf("Expected CreateDir to succeed, got %v", resDir.Fault().Error())
	}

	stat, err := os.Stat(dirPath)
	if err != nil || !stat.IsDir() {
		t.Fatalf("Directory was not created")
	}

	filePath := filepath.Join(dirPath, "test.txt")
	resFile := CreateFile(filePath, FilePermStandard)
	if resFile.HasError() {
		t.Fatalf("Expected CreateFile to succeed, got %v", resFile.Fault().Error())
	}

	f := resFile.Data()
	f.Close()

	stat, err = os.Stat(filePath)
	if err != nil || stat.IsDir() {
		t.Fatalf("File was not created")
	}
}
