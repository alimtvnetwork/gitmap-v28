package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func printMoveDryRun(src, dest string) {
	fmt.Println("  (dry-run) Move operation preview:")
	fmt.Printf("     Source:      %s\n", src)
	fmt.Printf("     Destination: %s\n", dest)
	fmt.Println("     Actions:     Relocate directory, update SQLite Repo path, sync VS Code & GitHub Desktop")
}

func confirmMovePrompt(slug, src, dest string) bool {
	fmt.Printf("Move %s\n  From: %s\n  To:   %s\nProceed? [y/N] ", slug, src, dest)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(line))
	return ans == "y" || ans == "yes"
}
