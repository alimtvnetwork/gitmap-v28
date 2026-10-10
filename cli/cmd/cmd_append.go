package cmd

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"path/filepath"
)

func RunAppend(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: gitmap append <file> <content>")

		return nil
	}

	filePath := args[0]
	content := args[1]

	err := doAppendFile(filePath, content)
	if err != nil {
		fmt.Println("Error appending to file:", err)

		return apperror.WrapSimple(err, "append file")
	}

	return nil
}

func RunWrite(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: gitmap write <file> <content>")

		return nil
	}

	filePath := args[0]
	content := args[1]

	err := doWriteFile(filePath, content)
	if err != nil {
		fmt.Println("Error writing to file:", err)

		return apperror.WrapSimple(err, "write file")
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
