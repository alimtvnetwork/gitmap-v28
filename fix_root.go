package main

import (
	"os"
	"strings"
)

func main() {
	b, _ := os.ReadFile("gitmap/cmd/root.go")
	s := string(b)
	s = strings.ReplaceAll(s, "cliexit.Reportf(command, \"execute\", \"\", err)\n\t\t\tpanic(\"fatal error\")", "handleGlobalError(command, err)")
	s = strings.ReplaceAll(s, "cliexit.Reportf(command, \"execute\", \"\", err)\r\n\t\t\tpanic(\"fatal error\")", "handleGlobalError(command, err)")

	if !strings.Contains(s, "\"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror\"") {
		s = strings.Replace(s, "\"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit\"", "\"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror\"\n\t\"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit\"\n\t\"github.com/alimtvnetwork/gitmap-v28/gitmap/config\"", 1)
	}

	os.WriteFile("gitmap/cmd/root.go", []byte(s), 0644)
}
