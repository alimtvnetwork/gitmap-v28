package commitin

import (
	"flag"
	"io"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// csvHolder collects every CSV-shaped flag so finalizeFlagFanout can
// run the comma split + per-flag validation in one pass. Keeps the
// flag-registration func short.
type csvHolder struct {
	exclude          string
	messageExclude   string
	messagePrefix    string
	messageSuffix    string
	overrideMessages string
	weakWords        string
	languages        string
}

// newFlagSet registers every spec §2.5 flag in one place. Returns the
// flag set, the destination RawArgs, and the CSV holder for the
// post-parse fan-out step. Help output is silenced — Parse is meant
// to be called by an outer command that owns the help dispatch.
func newFlagSet() (*flag.FlagSet, *RawArgs, *csvHolder) {
	fs := flag.NewFlagSet(constants.CmdCommitIn, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	raw := &RawArgs{}
	csv := &csvHolder{}
	registerBoolFlags(fs, raw)
	registerStringFlags(fs, raw)
	registerCsvFlags(fs, csv)
	return fs, raw, csv
}

func registerBoolFlags(fs *flag.FlagSet, raw *RawArgs) {
	fs.BoolVar(&raw.UseDefaultProfile, constants.CommitInFlagDefault, false, constants.CommitInDescDefault)
	fs.BoolVar(&raw.UseDefaultProfile, constants.CommitInFlagDefaultShort, false, constants.CommitInDescDefault)
	fs.BoolVar(&raw.SaveProfileOverwrite, constants.CommitInFlagSaveProfileOverwrite, false, constants.CommitInDescSaveProfileOverwrite)
	fs.BoolVar(&raw.SetDefault, constants.CommitInFlagSetDefault, false, constants.CommitInDescSetDefault)
	fs.BoolVar(&raw.OverrideOnlyWeak, constants.CommitInFlagOverrideOnlyWeak, false, constants.CommitInDescOverrideOnlyWeak)
	fs.BoolVar(&raw.IsNoPrompt, constants.CommitInFlagNoPrompt, false, constants.CommitInDescNoPrompt)
	fs.BoolVar(&raw.IsDryRun, constants.CommitInFlagDryRun, false, constants.CommitInDescDryRun)
	fs.BoolVar(&raw.IsKeepTemp, constants.CommitInFlagKeepTemp, false, constants.CommitInDescKeepTemp)
	fs.BoolVar(&raw.IsNoReleaseBranch, constants.CommitInFlagNoReleaseBranch, false, constants.CommitInDescNoReleaseBranch)
}

func registerStringFlags(fs *flag.FlagSet, raw *RawArgs) {
	fs.StringVar(&raw.ProfileName, constants.CommitInFlagProfile, "", constants.CommitInDescProfile)
	fs.StringVar(&raw.SaveProfileName, constants.CommitInFlagSaveProfile, "", constants.CommitInDescSaveProfile)
	fs.StringVar(&raw.AuthorName, constants.CommitInFlagAuthorName, "", constants.CommitInDescAuthorName)
	fs.StringVar(&raw.AuthorEmail, constants.CommitInFlagAuthorEmail, "", constants.CommitInDescAuthorEmail)
	fs.StringVar(&raw.ConflictMode, constants.CommitInFlagConflict, "", constants.CommitInDescConflict)
	fs.StringVar(&raw.TitlePrefix, constants.CommitInFlagTitlePrefix, "", constants.CommitInDescTitlePrefix)
	fs.StringVar(&raw.TitleSuffix, constants.CommitInFlagTitleSuffix, "", constants.CommitInDescTitleSuffix)
	fs.StringVar(&raw.FunctionIntel, constants.CommitInFlagFunctionIntel, "", constants.CommitInDescFunctionIntel)
}

func registerCsvFlags(fs *flag.FlagSet, csv *csvHolder) {
	fs.StringVar(&csv.exclude, constants.CommitInFlagExclude, "", constants.CommitInDescExclude)
	fs.StringVar(&csv.messageExclude, constants.CommitInFlagMessageExclude, "", constants.CommitInDescMessageExclude)
	fs.StringVar(&csv.messagePrefix, constants.CommitInFlagMessagePrefix, "", constants.CommitInDescMessagePrefix)
	fs.StringVar(&csv.messageSuffix, constants.CommitInFlagMessageSuffix, "", constants.CommitInDescMessageSuffix)
	fs.StringVar(&csv.overrideMessages, constants.CommitInFlagOverrideMessages, "", constants.CommitInDescOverrideMessages)
	fs.StringVar(&csv.weakWords, constants.CommitInFlagWeakWords, "", constants.CommitInDescWeakWords)
	fs.StringVar(&csv.languages, constants.CommitInFlagLanguages, "", constants.CommitInDescLanguages)
}

// boolFlagSet enumerates every flag registered as bool above. Kept
// adjacent to registerBoolFlags so adding a bool flag can never desync.
func boolFlagSet() map[string]bool {
	return map[string]bool{
		constants.CommitInFlagDefault:              true,
		constants.CommitInFlagDefaultShort:         true,
		constants.CommitInFlagSaveProfileOverwrite: true,
		constants.CommitInFlagSetDefault:           true,
		constants.CommitInFlagOverrideOnlyWeak:     true,
		constants.CommitInFlagNoPrompt:             true,
		constants.CommitInFlagDryRun:               true,
		constants.CommitInFlagKeepTemp:             true,
		constants.CommitInFlagNoReleaseBranch:      true,
	}
}
