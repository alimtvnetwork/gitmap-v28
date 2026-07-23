// Package constants_test — unit tests for SEO-write constants integrity.
package constants_test

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestSEOSQL_CreateTableContainsColumns verifies table schema (v15: CommitTemplate singular + CommitTemplateId PK).
func TestSEOSQL_CreateTableContainsColumns(t *testing.T) {
	sql := constants.SQLCreateCommitTemplate

	required := []string{"CommitTemplateId", "Kind", "Template", "CreatedAt"}
	for _, col := range required {
		if !strings.Contains(sql, col) {
			t.Errorf("expected SQLCreateCommitTemplate to contain column %q", col)
		}
	}
}

// TestSEOSQL_InsertHasThreePlaceholders verifies insert SQL.
func TestSEOSQL_InsertHasThreePlaceholders(t *testing.T) {
	sql := constants.SQLInsertTemplate
	count := strings.Count(sql, "?")
	if count != 2 {
		t.Errorf("expected 2 placeholders in insert SQL, got %d", count)
	}
}

// TestSEOSQL_SelectByKindHasWhere verifies select SQL filters by Kind.
func TestSEOSQL_SelectByKindHasWhere(t *testing.T) {
	sql := constants.SQLSelectTemplatesByKind
	if !strings.Contains(sql, "WHERE Kind = ?") {
		t.Error("expected SELECT to filter by Kind")
	}
}

// TestSEOSQL_CountReturnsAggregate verifies count SQL.
func TestSEOSQL_CountReturnsAggregate(t *testing.T) {
	sql := constants.SQLCountTemplates
	if !strings.Contains(sql, "COUNT(*)") {
		t.Error("expected COUNT(*) in count SQL")
	}
}

// TestSEOErrorMessages_NonEmpty verifies all error messages are non-empty.
func TestSEOErrorMessages_NonEmpty(t *testing.T) {
	errors := []string{
		constants.ErrSEOURLRequired,
		constants.ErrSEOCSVRead,
		constants.ErrSEOCSVEmpty,
		constants.ErrSEOTemplateRead,
		constants.ErrSEOTemplateEmpty,
		constants.ErrSEOIntervalFmt,
		constants.ErrSEONoFiles,
		constants.ErrSEORotateNotFound,
		constants.ErrSEOGitStage,
		constants.ErrSEOGitCommit,
		constants.ErrSEOGitPush,
		constants.ErrSEOSeedRead,
		constants.ErrSEOCreateWrite,
		constants.ErrSEODBInsert,
	}

	for _, e := range errors {
		if len(e) == 0 {
			t.Error("found empty error message constant")
		}
	}
}

// TestSEOHelpStrings_NonEmpty verifies all help strings are non-empty.
func TestSEOHelpStrings_NonEmpty(t *testing.T) {
	helps := []string{
		constants.HelpSEOWrite,
		constants.HelpSEOWriteFlags,
		constants.HelpSEOCSV,
		constants.HelpSEOURL,
		constants.HelpSEOService,
		constants.HelpSEOArea,
		constants.HelpSEOCompany,
		constants.HelpSEOPhone,
		constants.HelpSEOEmail,
		constants.HelpSEOAddress,
		constants.HelpSEOMaxCommits,
		constants.HelpSEOInterval,
		constants.HelpSEOFilesFlag,
		constants.HelpSEORotate,
		constants.HelpSEODryRunFlag,
		constants.HelpSEOTemplateF,
		constants.HelpSEOCreateTpl,
	}

	for _, h := range helps {
		if len(h) == 0 {
			t.Error("found empty help string constant")
		}
	}
}

// TestSEOMessages_ContainFormatVerbs verifies format messages have placeholders.
func TestSEOMessages_ContainFormatVerbs(t *testing.T) {
	formatted := map[string]string{
		"MsgSEOHeader":   constants.MsgSEOHeader,
		"MsgSEOCommit":   constants.MsgSEOCommit,
		"MsgSEORotation": constants.MsgSEORotation,
		"MsgSEODone":     constants.MsgSEODone,
		"MsgSEODryTitle": constants.MsgSEODryTitle,
		"MsgSEODryDesc":  constants.MsgSEODryDesc,
		"MsgSEOSeeded":   constants.MsgSEOSeeded,
		"MsgSEOWaiting":  constants.MsgSEOWaiting,
	}

	for name, msg := range formatted {
		if !strings.Contains(msg, "%") {
			t.Errorf("expected format verb in %s: %q", name, msg)
		}
	}
}

// TestSEOTableName verifies the table name constant (v15: singular).
func TestSEOTableName(t *testing.T) {
	if constants.TableCommitTemplate != "CommitTemplate" {
		t.Errorf("expected CommitTemplate, got %q", constants.TableCommitTemplate)
	}
}
