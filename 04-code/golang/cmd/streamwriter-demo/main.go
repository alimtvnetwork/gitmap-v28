package main

import (
	"fmt"
	"os"

	"coding-guidelines/common/examples"
)

func main() {
	fmt.Println("================================================================")
	fmt.Println("        STREAMWRITER DEMO: LOGGER, JSON, AND STREAMER          ")
	fmt.Println("================================================================")

	// 1. Run Logger Example
	fmt.Println("\n>>> [1] RUNNING LOGGER DEMONSTRATION")
	fmt.Println("----------------------------------------------------------------")
	if appErr := examples.RunLoggerExample(os.Stdout); appErr != nil {
		fmt.Fprintf(os.Stderr, "Logger example failed: %v\n", appErr)
		os.Exit(1)
	}

	// 2. Run Json Example
	fmt.Println("\n>>> [2] RUNNING JSON DEMONSTRATION")
	fmt.Println("----------------------------------------------------------------")
	if appErr := examples.RunJsonExample(os.Stdout); appErr != nil {
		fmt.Fprintf(os.Stderr, "Json example failed: %v\n", appErr)
		os.Exit(1)
	}

	// 3. Run Streamer Example
	fmt.Println("\n>>> [3] RUNNING STREAMER & WRITER DEMONSTRATION")
	fmt.Println("----------------------------------------------------------------")
	if appErr := examples.RunStreamerExample(os.Stdout); appErr != nil {
		fmt.Fprintf(os.Stderr, "Streamer example failed: %v\n", appErr)
		os.Exit(1)
	}

	fmt.Println("\n================================================================")
	fmt.Println("              ALL DEMONSTRATIONS COMPLETED SUCCESSFULLY         ")
	fmt.Println("================================================================")
}
