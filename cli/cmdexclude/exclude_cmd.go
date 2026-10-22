package cmdexclude

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunExclude routes exclude and repo-exclude subcommands.
func RunExclude(args []string) error {
	if len(args) == 0 {
		return PrintExclusionList()
	}

	sub := strings.ToLower(args[0])
	switch sub {
	case "add":
		return runExcludeAdd(args[1:])
	case "rm", "remove", "delete":
		return runExcludeRemove(args[1:])
	case "clean", "purge":
		return RunCleanRepos(args[1:])
	case "list", "ls", "show":
		return PrintExclusionList()
	default:
		// If first argument is not a subcommand, treat all args as patterns to add
		return runExcludeAdd(args)
	}
}

func runExcludeAdd(patterns []string) error {
	if len(patterns) == 0 {
		fmt.Println("Usage: gitmap exclude add <path-or-pattern> [...]")
		return nil
	}

	existing, _ := LoadCustomExclusions()
	set := make(map[string]bool)
	for _, p := range existing {
		set[strings.ToLower(p)] = true
	}

	addedCount := 0
	for _, p := range patterns {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		lower := strings.ToLower(trimmed)
		if !set[lower] {
			existing = append(existing, trimmed)
			set[lower] = true
			addedCount++
			fmt.Printf("  ✓ Added custom exclusion: %s\n", trimmed)
		} else {
			fmt.Printf("  ◦ Already excluded: %s\n", trimmed)
		}
	}

	if addedCount > 0 {
		if err := SaveCustomExclusions(existing); err != nil {
			return apperror.WrapSimple(err, "failed to save custom exclusions")
		}
		fmt.Printf("Saved %d new exclusion pattern(s).\n", addedCount)
	}

	return nil
}

func runExcludeRemove(patterns []string) error {
	if len(patterns) == 0 {
		fmt.Println("Usage: gitmap exclude rm <path-or-pattern> [...]")
		return nil
	}

	existing, _ := LoadCustomExclusions()
	toRemove := make(map[string]bool)
	for _, p := range patterns {
		toRemove[strings.ToLower(strings.TrimSpace(p))] = true
	}

	var kept []string
	removedCount := 0
	for _, p := range existing {
		if toRemove[strings.ToLower(p)] {
			removedCount++
			fmt.Printf("  ✓ Removed custom exclusion: %s\n", p)
		} else {
			kept = append(kept, p)
		}
	}

	if removedCount > 0 {
		if err := SaveCustomExclusions(kept); err != nil {
			return apperror.WrapSimple(err, "failed to update custom exclusions")
		}
		fmt.Printf("Updated exclusions: %d removed, %d remaining.\n", removedCount, len(kept))
	} else {
		fmt.Println("No matching custom exclusion patterns found to remove.")
	}

	return nil
}

// PrintExclusionList renders active canonical and user-defined exclusion patterns.
func PrintExclusionList() error {
	fmt.Println("╔════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║               GitMap Repository Exclusion Registry                 ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════╝")

	fmt.Println("\n■ Canonical Default Excluded Directories:")
	for i, d := range constants.DefaultScanExcludeDirs {
		fmt.Printf("  • %-18s", d)
		if (i+1)%3 == 0 {
			fmt.Println()
		}
	}
	if len(constants.DefaultScanExcludeDirs)%3 != 0 {
		fmt.Println()
	}

	custom, err := LoadCustomExclusions()
	if err != nil {
		return apperror.WrapSimple(err, "could not read custom exclusions")
	}

	fmt.Printf("\n■ Custom User Exclusions (%d active):\n", len(custom))
	if len(custom) == 0 {
		fmt.Println("  (No custom exclusions set. Use 'gitmap exclude add <pattern>' to add.)")
	} else {
		for _, p := range custom {
			fmt.Printf("  • %s\n", p)
		}
	}

	fmt.Println("\nCommands:")
	fmt.Println("  gitmap exclude add <pattern>      Add directory or pattern to exclusions")
	fmt.Println("  gitmap exclude rm <pattern>       Remove pattern from custom exclusions")
	fmt.Println("  gitmap exclude clean              Purge excluded/cache repos from GitMap DB")

	return nil
}
