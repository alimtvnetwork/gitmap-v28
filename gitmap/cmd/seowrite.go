// Package cmd — seowrite.go handles flag parsing and orchestration for seo-write.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// seoWriteFlags holds parsed flags for the seo-write command.
type seoWriteFlags struct {
	csv            string
	url            string
	service        string
	area           string
	company        string
	phone          string
	email          string
	address        string
	maxCommits     int
	interval       string
	files          string
	rotateFile     string
	dryRun         bool
	templatePath   string
	createTemplate bool
	authorName     string
	authorEmail    string
}

// runSEOWrite is the entry point for the seo-write command.
func runSEOWrite(args []string) {
	checkHelp("seo-write", args)
	if isCreateTemplateShorthand(args) {
		createTemplateFile()

		return
	}

	flags := parseSEOWriteFlags(args)
	if flags.createTemplate {
		createTemplateFile()

		return
	}

	executeSEOWrite(flags)
}

// isCreateTemplateShorthand checks if the first arg is the "ct" alias.
func isCreateTemplateShorthand(args []string) bool {
	if len(args) < 1 {
		return false
	}

	return args[0] == constants.CmdCreateTemplate
}

// executeSEOWrite runs the main seo-write workflow.
func executeSEOWrite(flags seoWriteFlags) {
	messages := resolveMessages(flags)
	if len(messages) == 0 {
		return
	}

	if flags.dryRun {
		printDryRun(messages, flags)

		return
	}

	intervalMin, intervalMax := parseInterval(flags.interval)
	runCommitLoop(flags, messages, intervalMin, intervalMax)
}

// resolveMessages loads commit messages from CSV or templates.
func resolveMessages(flags seoWriteFlags) []commitMessage {
	if flags.csv != "" {
		return loadCSVMessages(flags.csv)
	}

	if flags.url == "" {
		fmt.Fprint(os.Stderr, constants.ErrSEOURLRequired)
		os.Exit(1)
	}

	return loadTemplateMessages(flags)
}

// printDryRun outputs all planned commit messages without executing.
func printDryRun(messages []commitMessage, flags seoWriteFlags) {
	if flags.authorName != "" || flags.authorEmail != "" {
		author := resolveAuthorFlag(flags.authorName, flags.authorEmail)
		fmt.Printf(constants.MsgSEODryAuthor, author)
	}

	for i, m := range messages {
		fmt.Printf(constants.MsgSEODryTitle, i+1, m.title)
		fmt.Printf(constants.MsgSEODryDesc, m.description)
	}
}

// parseSEOWriteFlags parses command-line flags for seo-write.
func parseSEOWriteFlags(args []string) seoWriteFlags {
	fs := flag.NewFlagSet(constants.CmdSEOWrite, flag.ExitOnError)
	var f seoWriteFlags

	fs.StringVar(&f.csv, constants.FlagSEOCSV, "", constants.FlagDescSEOCSV)
	fs.StringVar(&f.url, constants.FlagSEOURL, "", constants.FlagDescSEOURL)
	fs.StringVar(&f.service, constants.FlagSEOService, "", constants.FlagDescSEOService)
	fs.StringVar(&f.area, constants.FlagSEOArea, "", constants.FlagDescSEOArea)
	fs.StringVar(&f.company, constants.FlagSEOCompany, "", constants.FlagDescSEOCompany)
	fs.StringVar(&f.phone, constants.FlagSEOPhone, "", constants.FlagDescSEOPhone)
	fs.StringVar(&f.email, constants.FlagSEOEmail, "", constants.FlagDescSEOEmail)
	fs.StringVar(&f.address, constants.FlagSEOAddress, "", constants.FlagDescSEOAddress)
	fs.IntVar(&f.maxCommits, constants.FlagSEOMaxCommits, 0, constants.FlagDescSEOMaxCommits)
	fs.StringVar(&f.interval, constants.FlagSEOInterval, constants.SEODefaultInterval, constants.FlagDescSEOInterval)
	fs.StringVar(&f.files, constants.FlagSEOFiles, "", constants.FlagDescSEOFiles)
	fs.StringVar(&f.rotateFile, constants.FlagSEORotateFile, "", constants.FlagDescSEORotateFile)
	fs.BoolVar(&f.dryRun, constants.FlagSEODryRun, false, constants.FlagDescSEODryRun)
	fs.StringVar(&f.templatePath, constants.FlagSEOTemplate, "", constants.FlagDescSEOTemplate)
	fs.BoolVar(&f.createTemplate, constants.FlagSEOCreateTemplate, false, constants.FlagDescSEOCreateTemplate)
	fs.StringVar(&f.authorName, constants.FlagSEOAuthorName, "", constants.FlagDescSEOAuthorName)
	fs.StringVar(&f.authorEmail, constants.FlagSEOAuthorEmail, "", constants.FlagDescSEOAuthorEmail)

	_ = fs.Parse(args)

	return f
}
