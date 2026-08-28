package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func runHead(args []string) error {
	if len(args) == 0 {
		return apperror.New("Usage: gitmap head <file> [lines]", "E9000", nil)
	}

	filePath := args[0]
	lines := 10
	if len(args) > 1 {
		if val, err := strconv.Atoi(args[1]); err == nil && val > 0 {
			lines = val
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return apperror.Wrap(err, "Error opening file:", nil)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		count++
		if count >= lines {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return apperror.Wrap(err, "Error reading file:", nil)
	}
	return nil
}

func runTail(args []string) error {
	if len(args) == 0 {
		return apperror.New("Usage: gitmap tail <file> [lines]", "E9000", nil)
	}

	filePath := args[0]
	lines := 10
	if len(args) > 1 {
		if val, err := strconv.Atoi(args[1]); err == nil && val > 0 {
			lines = val
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return apperror.Wrap(err, "Error opening file:", nil)
	}
	defer file.Close()

	// A simple circular buffer approach
	buffer := make([]string, lines)
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		buffer[count%lines] = scanner.Text()
		count++
	}

	if err := scanner.Err(); err != nil {
		return apperror.Wrap(err, "Error reading file:", nil)
	}

	start := 0
	if count > lines {
		start = count % lines
	} else {
		lines = count
	}

	for i := 0; i < lines; i++ {
		fmt.Println(buffer[(start+i)%len(buffer)])
	}
	return nil
}
