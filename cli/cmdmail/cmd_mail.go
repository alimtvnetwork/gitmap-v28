package cmdmail

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunMail is the top-level CLI router for the gitmap mail command family.
func RunMail(args []string) *appfault.AppError {
	if len(args) == 0 || isHelpArg(args[0]) {
		printMailHelp()
		return nil
	}
	switch strings.ToLower(args[0]) {
	case "send":
		return runMailSend(args[1:])
	case "check":
		return runMailCheck(args[1:])
	case "test":
		return runMailTest(args[1:])
	case "template":
		return runMailTemplate(args[1:])
	case "notify":
		return runMailNotify(args[1:])
	case "ui":
		return runMailUI(args[1:])
	default:
		return appfault.NewValidationError(fmt.Sprintf("unknown mail subcommand: %s", args[0]))
	}
}

func isHelpArg(arg string) bool {
	return slices.Contains([]string{"-h", "--help", "help"}, arg)
}

// RunMailUI boots the embedded mail composer web server.
func RunMailUI(port int) error {
	server, listener, boundPort, err := StartMailUIServer(MailUIOptions{Port: port})
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Printf("[GITMAP MAIL UI] Server running at http://127.0.0.1:%d/ui\n", boundPort)
	return server.Serve(listener)
}

func runMailUI(args []string) *appfault.AppError {
	port := DefaultMailUIPort
	for i := 0; i < len(args); i++ {
		if (args[i] == "--port" || args[i] == "-p") && i+1 < len(args) {
			port, i = parsePortFlag(args[i+1], port), i+1
		}
	}
	if err := RunMailUI(port); err != nil {
		return appfault.WrapExecution(err, "run mail ui")
	}
	return nil
}

func parsePortFlag(arg string, fallback int) int {
	if p, err := strconv.Atoi(arg); err == nil && p > 0 {
		return p
	}
	return fallback
}

func runMailSend(args []string) *appfault.AppError {
	opts, err := parseSendArgs(args)
	if err != nil {
		return appfault.NewValidationError(err.Error())
	}
	if len(opts.To) == 0 || opts.Subject == "" {
		return appfault.NewValidationError("missing required arguments: --to and --subject")
	}
	if opts.From == "" {
		opts.From = opts.AccountName + "@gitmap.local"
	}
	if opts.IsDryRun {
		fmt.Printf("[GITMAP MAIL] DRY RUN: Validated to %s (subject: %s)\n", MaskEmailList(opts.To), opts.Subject)
		return nil
	}
	if err := CurrentTransmitter.SendMail(opts, BuildMIMEMessage(opts)); err != nil {
		return appfault.WrapExecution(err, "send email")
	}
	fmt.Printf("[GITMAP MAIL] Sent to %s (subject: %s)\n", MaskEmailList(opts.To), opts.Subject)
	return nil
}

func parseSendArgs(args []string) (MailSendOptions, error) {
	opts := MailSendOptions{AccountName: "default", Host: "127.0.0.1", Port: 587}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case (a == "--to" || a == "-t") && i+1 < len(args):
			opts.To, i = strings.Split(args[i+1], ","), i+1
		case (a == "--subject" || a == "-s") && i+1 < len(args):
			opts.Subject, i = args[i+1], i+1
		case (a == "--body" || a == "-b") && i+1 < len(args):
			opts.Body, i = args[i+1], i+1
		case (a == "--file" || a == "-f") && i+1 < len(args):
			body, err := readBodyFile(args[i+1])
			if err != nil {
				return opts, err
			}
			opts.Body, i = body, i+1
		case a == "--html":
			opts.IsHTML = true
		case (a == "--account" || a == "-a") && i+1 < len(args):
			opts.AccountName, i = args[i+1], i+1
		case a == "--dry-run":
			opts.IsDryRun = true
		}
	}
	return opts, nil
}

func readBodyFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", appfault.WrapSimple(err, "read body file")
	}
	return string(content), nil
}

func runMailCheck(args []string) *appfault.AppError {
	filter, isJSON := parseCheckArgs(args)
	store, _ := OpenMailStore("")
	if store != nil {
		defer store.Close()
	}
	summary, err := EvaluateMailboxes(store, filter)
	if err != nil {
		return appfault.WrapExecution(err, "evaluate mailboxes")
	}
	if isJSON {
		return outputMailCheckJSON(summary)
	}
	fmt.Print(FormatMailCheckTable(summary, filter.ShowFullEmail))
	return nil
}

func outputMailCheckJSON(summary *MailCheckSummary) *appfault.AppError {
	jsonStr, err := FormatMailCheckJSON(summary)
	if err != nil {
		return appfault.WrapExecution(err, "format json")
	}
	fmt.Println(jsonStr)
	return nil
}

func parseCheckArgs(args []string) (MailboxCheckFilter, bool) {
	filter := MailboxCheckFilter{CheckAll: true}
	isJSON := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case (a == "--account" || a == "-a") && i+1 < len(args):
			filter.AccountName, filter.CheckAll, i = args[i+1], false, i+1
		case a == "--all":
			filter.CheckAll = true
		case a == "--unread-only" || a == "-u":
			filter.FilterUnreadOnly = true
		case a == "--show-full-email":
			filter.ShowFullEmail = true
		case a == "--json" || a == "-j":
			isJSON = true
		}
	}
	return filter, isJSON
}

func runMailTest(args []string) *appfault.AppError {
	host := "127.0.0.1"
	port := 587
	isTLS := false
	for i := 0; i < len(args); i++ {
		switch {
		case (args[i] == "--host" || args[i] == "-h") && i+1 < len(args):
			host, i = args[i+1], i+1
		case (args[i] == "--port" || args[i] == "-p") && i+1 < len(args):
			port, i = parsePortFlag(args[i+1], port), i+1
		case args[i] == "--tls":
			isTLS = true
		}
	}
	return executeTestConnection(host, port, isTLS)
}

func executeTestConnection(host string, port int, isTLS bool) *appfault.AppError {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	start := time.Now()
	if port == 465 || isTLS {
		return testTLSConnection(host, addr, start)
	}
	return testPlainConnection(addr, start)
}

func testTLSConnection(host, addr string, start time.Time) *appfault.AppError {
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", addr, &tls.Config{ServerName: host, InsecureSkipVerify: isLocalHost(host)})
	if err != nil {
		return appfault.WrapExecution(err, "test tls connection")
	}
	_ = conn.Close()
	fmt.Printf("[GITMAP MAIL TEST] SUCCESS: Connected to %s (%dms)\n", addr, time.Since(start).Milliseconds())
	return nil
}

func testPlainConnection(addr string, start time.Time) *appfault.AppError {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return appfault.WrapExecution(err, "test connection")
	}
	_ = conn.Close()
	fmt.Printf("[GITMAP MAIL TEST] SUCCESS: Connected to %s (%dms)\n", addr, time.Since(start).Milliseconds())
	return nil
}

func printMailHelp() {
	const helpText = `Usage: gitmap mail <command> [flags]

Commands:
  send        Send an email directly via SMTP
  check       Check configured mailboxes and unread counts
  test        Test SMTP connectivity and TLS handshakes
  template    Manage email templates (ls, view, add, rm)
  notify      Folder watcher notifier (watch, ls, trigger)
  ui          Start the embedded mailing composer web UI

Flags:
  --to, -t             Recipient addresses (comma-separated)
  --subject, -s        Subject line
  --body, -b           Body content
  --file, -f           Path to message body file
  --html               Send as HTML format
  --unread-only, -u    Filter mailboxes to unread messages only
  --show-full-email    Display unmasked email addresses
  --json, -j           Emit output as Universal JSON Envelope`
	fmt.Println(helpText)
}
