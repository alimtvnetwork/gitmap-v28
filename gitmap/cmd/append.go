package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func runAppend(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: gitmap append <file> <content>")
		return nil
	}
	filePath := args[0]
	content := args[1]

	err := doAppendFile(filePath, content)
	if err != nil {
		fmt.Println("Error appending to file:", err)
		return apperror.New("fatal error", "E9000", nil)
	}
	return nil
}

func runWrite(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: gitmap write <file> <content>")
		return nil
	}
	filePath := args[0]
	content := args[1]

	err := doWriteFile(filePath, content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return apperror.New("fatal error", "E9000", nil)
	}
	return nil
}

func doAppendFile(filePath string, content string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Ensure a trailing newline just like PowerShell's Add-Content does
	_, err = f.WriteString(content + "\n")
	return err
}

func doWriteFile(filePath string, content string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content + "\n")
	return err
}
