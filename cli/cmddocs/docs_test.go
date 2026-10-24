package cmddocs

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddocs/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddocs/enums"
)

func TestResolveViewerMode(t *testing.T) {
	mode := ResolveViewerMode()

	if mode != enums.DocsViewerModeBrowser || !mode.IsBrowser() {
		t.Errorf("expected default viewer mode to be browser, got %s", mode)
	}
}

func TestResolveBrowserCmd(t *testing.T) {
	testResolveBrowserCmdWindows(t)
	testResolveBrowserCmdDarwin(t)
	testResolveBrowserCmdLinux(t)
	testResolveBrowserCmdFallback(t)
}

func testResolveBrowserCmdWindows(t *testing.T) {
	url := "https://example.com"
	cmd := resolveBrowserCmd(constants.OSWindows, url)

	if len(cmd.Args) != 4 {
		t.Fatalf("expected 4 args for windows, got %d", len(cmd.Args))
	}
	if cmd.Args[0] != constants.CmdWindowsShell || cmd.Args[1] != constants.CmdArgSlashC {
		t.Errorf("unexpected windows command prefix: %v", cmd.Args)
	}
	if cmd.Args[2] != constants.CmdArgStart || cmd.Args[3] != url {
		t.Errorf("unexpected windows command suffix: %v", cmd.Args)
	}
}

func testResolveBrowserCmdDarwin(t *testing.T) {
	url := "https://example.com"
	cmd := resolveBrowserCmd(constants.OSDarwin, url)

	if len(cmd.Args) != 2 {
		t.Fatalf("expected 2 args for darwin, got %d", len(cmd.Args))
	}
	if cmd.Args[0] != constants.CmdOpen || cmd.Args[1] != url {
		t.Errorf("unexpected darwin command args: %v", cmd.Args)
	}
}

func testResolveBrowserCmdLinux(t *testing.T) {
	url := "https://example.com"
	cmd := resolveBrowserCmd(constants.OSLinux, url)

	if len(cmd.Args) != 2 {
		t.Fatalf("expected 2 args for linux, got %d", len(cmd.Args))
	}
	if cmd.Args[0] != constants.CmdXdgOpen || cmd.Args[1] != url {
		t.Errorf("unexpected linux command args: %v", cmd.Args)
	}
}

func testResolveBrowserCmdFallback(t *testing.T) {
	url := "https://example.com"
	cmd := resolveBrowserCmd("freebsd", url)

	if len(cmd.Args) != 2 {
		t.Fatalf("expected 2 args for fallback, got %d", len(cmd.Args))
	}
	if cmd.Args[0] != constants.CmdXdgOpen || cmd.Args[1] != url {
		t.Errorf("unexpected fallback command args: %v", cmd.Args)
	}
}

func TestRunDocsSuccess(t *testing.T) {
	origStart := startCmdFn
	defer func() { startCmdFn = origStart }()

	startCmdFn = func(_ *exec.Cmd) error {
		return nil
	}

	err := RunDocs([]string{"docs"})

	if err != nil {
		t.Errorf("expected RunDocs to succeed, got %v", err)
	}
}

func TestRunDocsFailure(t *testing.T) {
	origStart := startCmdFn
	defer func() { startCmdFn = origStart }()

	startCmdFn = func(_ *exec.Cmd) error {
		return errors.New("launch error")
	}

	err := RunDocs([]string{"docs"})

	if err == nil {
		t.Error("expected RunDocs to fail when startCmdFn fails")
	}
}

func TestRunDocsArgsHandling(t *testing.T) {
	origStart := startCmdFn
	origHelp := CheckHelpFn
	defer func() {
		startCmdFn = origStart
		CheckHelpFn = origHelp
	}()

	startCmdFn = func(_ *exec.Cmd) error { return nil }

	var capturedCmd string
	var capturedArgs []string
	CheckHelpFn = func(cmd string, args []string) {
		capturedCmd = cmd
		capturedArgs = args
	}

	args := []string{"--help", "--browser"}
	err := RunDocs(args)

	if err != nil {
		t.Errorf("expected RunDocs to succeed, got %v", err)
	}
	if capturedCmd != constants.CmdDocs || len(capturedArgs) != 2 {
		t.Errorf("unexpected checkHelp capture: cmd=%s, args=%v", capturedCmd, capturedArgs)
	}
}

func TestCheckHelpNilShim(t *testing.T) {
	origHelp := CheckHelpFn
	defer func() { CheckHelpFn = origHelp }()

	CheckHelpFn = nil
	checkHelp("docs", []string{"--help"})
}
