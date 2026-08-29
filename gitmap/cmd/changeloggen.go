package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// runChangelogGen handles the 'changelog-generate' command.
func runChangelogGen(args []string) error {
	checkHelp("changelog-generate", args)

	from, to, write := parseChangelogGenFlags(args)
	fromTag, toRef, err := release.ResolveTagRange(from, to)

	if err != nil {
		return apperror.WrapSimple(err, "✗")
	}

	commits, err := release.GenerateChangelog(fromTag, toRef)
	if err != nil {
		return apperror.WrapSimple(err, "✗")
	}

	if len(commits) == 0 {
		fmt.Printf(constants.MsgChangelogGenEmpty, fromTag, toRef)

		return nil
	}

	version := resolveGenVersion(toRef)
	section := release.FormatChangelogSection(version, commits)

	if !write {
		printChangelogPreview(section)
		return nil
	}
	return writeChangelogSection(section)
}

// parseChangelogGenFlags parses flags for changelog-generate.
func parseChangelogGenFlags(args []string) (from, to string, write bool) {
	fs := flag.NewFlagSet(constants.CmdChangelogGen, flag.ExitOnError)
	fromFlag := fs.String("from", "", constants.FlagDescFrom)
	toFlag := fs.String("to", "", constants.FlagDescTo)
	writeFlag := fs.Bool("write", false, constants.FlagDescWrite)
	_ = fs.Parse(args)

	return *fromFlag, *toFlag, *writeFlag
}

// resolveGenVersion returns the version label for the generated section.
func resolveGenVersion(toRef string) string {
	if toRef == constants.GitHEAD {
		return "Unreleased"
	}

	return release.NormalizeVersion(toRef)
}

// printChangelogPreview prints the generated changelog to stdout.
func printChangelogPreview(section string) {
	fmt.Print(constants.MsgChangelogGenPreview)
	fmt.Print(section)
}

// writeChangelogSection prepends the section to changelog.md.
func writeChangelogSection(section string) *apperror.AppError {
	existing, err := os.ReadFile(constants.ChangelogFile)
	if err != nil && !os.IsNotExist(err) {
		return apperror.WrapSimple(err, constants.ErrChangelogGenRead)
	}

	content := section + "\n" + string(existing)

	err = os.WriteFile(constants.ChangelogFile, []byte(content), 0o644)
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrChangelogGenWrite)
	}

	fmt.Printf(constants.MsgChangelogGenWritten, constants.ChangelogFile)
	return nil
}
