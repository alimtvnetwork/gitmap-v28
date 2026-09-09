package cmd

import (
	"strings"
	"testing"
)

func TestCrontabE2EFirstTimeSetupNoCrontab(t *testing.T) {
	origRead := readCrontabFunc
	origWrite := writeCrontabFunc
	defer func() {
		readCrontabFunc = origRead
		writeCrontabFunc = origWrite
	}()

	readCrontabFunc = func() string {
		return ""
	}

	var writtenCrontab string
	writeCrontabFunc = func(content string) error {
		writtenCrontab = content

		return nil
	}

	err := ensureCrontabPersistence()
	if err != nil {
		t.Fatalf("ensureCrontabPersistence failed: %v", err)
	}

	if strings.Contains(writtenCrontab, "no crontab") {
		t.Fatalf("written crontab must not contain 'no crontab': %q", writtenCrontab)
	}

	if !strings.HasPrefix(writtenCrontab, "@reboot") {
		t.Fatalf("expected written crontab to start with @reboot: %q", writtenCrontab)
	}

	if !strings.HasSuffix(writtenCrontab, "\n") {
		t.Fatalf("expected written crontab to end with newline: %q", writtenCrontab)
	}
}

func TestCrontabE2EIdempotencyAndPreservation(t *testing.T) {
	origRead := readCrontabFunc
	origWrite := writeCrontabFunc
	defer func() {
		readCrontabFunc = origRead
		writeCrontabFunc = origWrite
	}()

	existingJob := "0 12 * * * /usr/local/bin/backup.sh"
	alreadyPersisted := existingJob + "\n" + crontabRebootLine + "\n"

	readCrontabFunc = func() string {
		return alreadyPersisted
	}

	wasWritten := false
	writeCrontabFunc = func(content string) error {
		wasWritten = true

		return nil
	}

	err := ensureCrontabPersistence()
	if err != nil {
		t.Fatalf("expected nil error on already persisted crontab, got: %v", err)
	}

	if wasWritten {
		t.Errorf("writeCrontabFunc should not be called when entry already exists")
	}
}

func TestCrontabE2ENoBadMinuteError(t *testing.T) {
	errTextFromUbuntu := "no crontab for a\n"
	cleaned := cleanCrontabOutput(errTextFromUbuntu)
	if cleaned != "" {
		t.Fatalf("cleanCrontabOutput failed to clear error output: %q", cleaned)
	}

	finalCrontab := buildUpdatedCrontab(cleaned, crontabRebootLine)
	firstWord := strings.Fields(finalCrontab)[0]
	if firstWord != "@reboot" {
		t.Fatalf("first word of crontab is %q (would cause 'bad minute' if not @reboot)", firstWord)
	}
}

func TestCrontabE2EPreserveExistingJobs(t *testing.T) {
	origRead := readCrontabFunc
	origWrite := writeCrontabFunc
	defer func() {
		readCrontabFunc = origRead
		writeCrontabFunc = origWrite
	}()

	userJob := "30 2 * * * /opt/scripts/nightly-backup.sh"
	readCrontabFunc = func() string {
		return userJob
	}

	var writtenContent string
	writeCrontabFunc = func(content string) error {
		writtenContent = content

		return nil
	}

	err := ensureCrontabPersistence()
	if err != nil {
		t.Fatalf("ensureCrontabPersistence failed: %v", err)
	}

	if !strings.Contains(writtenContent, userJob) {
		t.Errorf("expected original user job %q to be preserved in %q", userJob, writtenContent)
	}

	if !strings.Contains(writtenContent, crontabRebootLine) {
		t.Errorf("expected reboot line in %q", writtenContent)
	}
}
