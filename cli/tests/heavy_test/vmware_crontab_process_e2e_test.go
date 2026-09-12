package heavy_test

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvmware"
)

func TestCrontabHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_CRONTAB_HELPER") != "1" {
		return
	}

	args := extractHelperArgs()
	handleHelperExecution(args)
}

func extractHelperArgs() []string {
	args := os.Args
	for i, arg := range args {
		if arg == "--" && i+1 < len(args) {
			return args[i+1:]
		}
	}

	return nil
}

func handleHelperExecution(args []string) {
	if len(args) == 0 {
		os.Exit(2)
	}

	dispatchHelperSubcommand(args[0])
}

func dispatchHelperSubcommand(sub string) {
	if sub == "-l" {
		handleHelperList()

		return
	}

	if sub == "-" {
		handleHelperWrite()

		return
	}

	os.Exit(2)
}

func handleHelperList() {
	statePath := os.Getenv("CRON_STATE_FILE")
	data, err := os.ReadFile(statePath)
	if err != nil {
		_, _ = os.Stderr.WriteString("no crontab for a\n")
		os.Exit(1)
	}

	_, _ = os.Stdout.Write(data)
	os.Exit(0)
}

func handleHelperWrite() {
	statePath := os.Getenv("CRON_STATE_FILE")
	data, _ := io.ReadAll(os.Stdin)
	content := string(data)
	if strings.Contains(content, "no crontab") {
		_, _ = os.Stderr.WriteString("\"-\":0: bad minute\nerrors in crontab file, can't install.\n")
		os.Exit(1)
	}

	_ = os.WriteFile(statePath, data, 0644)
	os.Exit(0)
}

func TestCrontabE2EFullUbuntuLifecycle(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "crontab.txt")
	origCmd := cmdvmware.CrontabCommandFunc
	defer func() { cmdvmware.CrontabCommandFunc = origCmd }()

	cmdvmware.CrontabCommandFunc = makeHelperCrontabCmd(statePath)

	err := cmdvmware.EnsureCrontabPersistence()
	if err != nil {
		t.Fatalf("first-time ensureCrontabPersistence failed: %v", err)
	}

	assertCrontabFileValid(t, statePath)
	assertSecondPersistenceCall(t)
}

func makeHelperCrontabCmd(statePath string) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		testArgs := []string{"-test.run=TestCrontabHelperProcess", "--"}
		testArgs = append(testArgs, args...)
		cmd := exec.Command(os.Args[0], testArgs...)
		cmd.Env = append(os.Environ(), "GO_WANT_CRONTAB_HELPER=1", "CRON_STATE_FILE="+statePath)

		return cmd
	}
}

func assertCrontabFileValid(t *testing.T, statePath string) {
	t.Helper()
	data, err := os.ReadFile(statePath)
	if err != nil {
		t.Fatalf("crontab file missing: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "@reboot") {
		t.Errorf("expected @reboot in saved crontab, got: %q", content)
	}

	if strings.Contains(content, "no crontab") {
		t.Errorf("saved crontab contains 'no crontab': %q", content)
	}
}

func assertSecondPersistenceCall(t *testing.T) {
	t.Helper()
	err := cmdvmware.EnsureCrontabPersistence()
	if err != nil {
		t.Fatalf("idempotent ensureCrontabPersistence failed: %v", err)
	}
}

func TestCrontabE2ERejectsBadMinuteFormat(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "crontab.txt")
	origCmd := cmdvmware.CrontabCommandFunc
	defer func() { cmdvmware.CrontabCommandFunc = origCmd }()

	cmdvmware.CrontabCommandFunc = makeHelperCrontabCmd(statePath)

	badCrontab := "no crontab for a\n@reboot mount\n"
	err := cmdvmware.WriteCrontabFunc(badCrontab)
	if err == nil {
		t.Fatalf("expected writeCrontab to fail on bad minute input")
	}

	assertBadMinuteErrorMessage(t, err.Error())
}

func assertBadMinuteErrorMessage(t *testing.T, errMsg string) {
	t.Helper()
	if !strings.Contains(errMsg, "bad minute") {
		t.Errorf("expected error to mention 'bad minute', got: %q", errMsg)
	}
}
