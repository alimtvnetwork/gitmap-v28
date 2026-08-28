package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func runHead(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: gitmap head <file> [lines]")
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
}

func runTail(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: gitmap tail <file> [lines]")
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
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
}
