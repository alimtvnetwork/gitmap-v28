package cmd

import "strings"

func advanceFlagIndex(arg string, idx, total int, flagsWithValues map[string]bool) int {
	hasVal := flagsWithValues[arg]
	hasMore := idx+1 < total
	if hasVal && hasMore {
		return idx + 1
	}

	return idx
}

func extractNonFlagTokens(args []string) []string {
	flagsWithValues := map[string]bool{
		"--dir": true, "--slug": true, "--description": true, "-d": true,
		"--profile": true, "--org": true,
	}
	var pos []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		isFlag := strings.HasPrefix(arg, "-")
		if !isFlag {
			pos = append(pos, arg)
			continue
		}
		i = advanceFlagIndex(arg, i, len(args), flagsWithValues)
	}

	return pos
}
