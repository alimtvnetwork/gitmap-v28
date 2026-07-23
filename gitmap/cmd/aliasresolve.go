package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// aliasContext holds the resolved alias information for the current command.
var aliasContext *resolvedAlias

// resolvedAlias stores the result of resolving a -A/--alias flag.
type resolvedAlias struct {
	Alias        string
	AbsolutePath string
	Slug         string
}

// extractAliasFlag scans args for -A or --alias and returns the alias
// name and the remaining args with the flag removed.
func extractAliasFlag(args []string) (string, []string) {
	for i, arg := range args {
		if arg == "-A" || arg == "--alias" {
			if i+1 < len(args) {
				return args[i+1], removeElements(args, i, 2)
			}

			fmt.Fprintln(os.Stderr, constants.ErrAliasEmpty)
			os.Exit(1)
		}

		if hasAliasPrefix(arg, "-A=") {
			return arg[3:], removeElements(args, i, 1)
		}
		if hasAliasPrefix(arg, "--alias=") {
			return arg[8:], removeElements(args, i, 1)
		}
	}

	return "", args
}

// hasAliasPrefix checks if an arg starts with the given prefix.
func hasAliasPrefix(arg, prefix string) bool {
	return len(arg) > len(prefix) && arg[:len(prefix)] == prefix
}

// removeElements removes count elements starting at index from a slice.
func removeElements(args []string, index, count int) []string {
	result := make([]string, 0, len(args)-count)
	result = append(result, args[:index]...)
	result = append(result, args[index+count:]...)

	return result
}

// resolveAliasContext looks up an alias and sets the global context.
// Returns the resolved repo path for use in commands.
func resolveAliasContext(aliasName string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	resolved, err := db.ResolveAlias(aliasName)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		exitWith(1)
	}

	aliasContext = &resolvedAlias{
		Alias:        resolved.Alias.Alias,
		AbsolutePath: resolved.AbsolutePath,
		Slug:         resolved.Slug,
	}

	fmt.Fprintf(os.Stderr, constants.MsgAliasResolved,
		resolved.Alias, resolved.AbsolutePath, resolved.Slug)
}

// GetAliasPath returns the resolved alias path if set, or empty string.
func GetAliasPath() string {
	if aliasContext == nil {
		return ""
	}

	return aliasContext.AbsolutePath
}

// GetAliasSlug returns the resolved alias slug if set, or empty string.
func GetAliasSlug() string {
	if aliasContext == nil {
		return ""
	}

	return aliasContext.Slug
}

// HasAlias returns true if a -A flag was resolved.
func HasAlias() bool {
	return aliasContext != nil
}

// AliasAsRecords returns the alias as a single-element ScanRecord slice.
// Useful for commands that operate on a list of repos.
func AliasAsRecords() []store.AliasWithRepo {
	if aliasContext == nil {
		return nil
	}

	return []store.AliasWithRepo{{
		AbsolutePath: aliasContext.AbsolutePath,
		Slug:         aliasContext.Slug,
	}}
}
