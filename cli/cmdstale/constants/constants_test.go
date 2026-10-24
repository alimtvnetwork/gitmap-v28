package constants

import (
	"testing"
)

func TestConstantsInvariants(t *testing.T) {
	testCommandConstants(t)
	testFlagConstants(t)
	testFlagArgConstants(t)
	testThresholdConstants(t)
	testDirectoryAndPathConstants(t)
	testFormatAndHeaderConstants(t)
	testGitTokensConstants(t)
	testMessageAndErrorConstants(t)
}

func testCommandConstants(t *testing.T) {
	if CmdStale != "stale" || CmdStaleAlias != "sta" {
		t.Error("command name and alias mismatch")
	}
	if SubCmdHelp != "help" {
		t.Error("help subcommand mismatch")
	}
}

func testFlagConstants(t *testing.T) {
	if FlagDays != "days" || FlagRoot != "root" {
		t.Error("flag name mismatch for days or root")
	}
	if FlagArchive != "archive" || FlagDryRun != "dry-run" {
		t.Error("flag name mismatch for archive or dry-run")
	}
	if FlagFormat != "format" || FlagHelp != "help" || FlagHelpShort != "h" {
		t.Error("flag name mismatch for format or help")
	}
	if FlagDescDays == "" || FlagDescRoot == "" || FlagDescArchive == "" {
		t.Error("flag descriptions must not be empty")
	}
}

func testFlagArgConstants(t *testing.T) {
	if FlagArgDays != "--days" || FlagArgRoot != "--root" {
		t.Error("flag arg string mismatch")
	}
	if FlagArgArchive != "--archive" || FlagArgDryRun != "--dry-run" {
		t.Error("flag arg string mismatch")
	}
	if FlagArgFormat != "--format" || FlagArgHelp != "--help" || FlagArgHelpShort != "-h" {
		t.Error("flag arg string mismatch")
	}
	if FlagHelpEntryHelp != "-h, --help" {
		t.Error("help entry mismatch")
	}
}

func testThresholdConstants(t *testing.T) {
	if DefaultStaleDays != 90 || DefaultScanRoot != "." {
		t.Error("default threshold mismatch")
	}
	if DefaultOutputFormat != "table" || HoursPerDay != 24 {
		t.Error("format or hours default mismatch")
	}
	if MinHygieneWorkers != 2 || MaxHygieneWorkers != 8 {
		t.Error("worker bounds mismatch")
	}
	if DefaultArchiveDirPerm != 0o755 {
		t.Error("archive directory perm mismatch")
	}
}

func testDirectoryAndPathConstants(t *testing.T) {
	if DirDotGit != ".git" || DirGitmap != ".gitmap" || DirArchive != "archive" {
		t.Error("directory name constants mismatch")
	}
}

func testFormatAndHeaderConstants(t *testing.T) {
	if DateFormatStale != "2006-01-02" || TimestampFormatArchive != "20060102T150405Z" {
		t.Error("date/time layout constants mismatch")
	}
	if JSONIndent != "  " || JSONPrefix != "" {
		t.Error("JSON layout constants mismatch")
	}
	if HeaderPath != "path" || HeaderLastCommitUTC != "last_commit_utc" || HeaderAgeDays != "age_days" {
		t.Error("CSV header constants mismatch")
	}
}

func testGitTokensConstants(t *testing.T) {
	if GitBin != "git" || GitFlagDir != "-C" || GitSubCmdLog != "log" {
		t.Error("git command tokens mismatch")
	}
	if GitArgLimitOne != "-1" || GitFormatCommitEpoch != "--format=%ct" {
		t.Error("git argument tokens mismatch")
	}
	if ColorCyan != "\033[36m" || ColorYellow != "\033[33m" || ColorReset != "\033[0m" {
		t.Error("ANSI color tokens mismatch")
	}
}

func testMessageAndErrorConstants(t *testing.T) {
	if MsgNoStaleReposFmt == "" || MsgStaleHeaderFmt == "" || MsgStaleRowFmt == "" {
		t.Error("message formats must not be empty")
	}
	if MsgWouldArchiveFmt == "" || MsgArchivedFmt == "" || MsgFailedFmt == "" {
		t.Error("archive message formats must not be empty")
	}
	if ErrJSONEncodeFmt == "" || ErrArchiveMkdirFmt == "" || ErrInvalidFormatFmt == "" {
		t.Error("error formats must not be empty")
	}
	if AgeDaysIntFmt != "%d" || SecScanFmt != "%d" || Newline != "\n" {
		t.Error("sentinel constants mismatch")
	}
}
