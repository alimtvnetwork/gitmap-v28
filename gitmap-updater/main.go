// Package main is the entry point for gitmap-updater.
package main

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap-updater/cmd"
)

func main() {
	if len(os.Args) < 2 {
		cmd.PrintUsage()
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "check":
		cmd.RunCheck()
	case "run":
		cmd.RunUpdate()
	case "update-worker":
		cmd.RunWorker()
	case "version", "v":
		fmt.Println(cmd.Version)
	case "help", "--help", "-h":
		cmd.PrintUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		cmd.PrintUsage()
		os.Exit(1)
	}
}
