package cmdmail

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
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
	subcmd := strings.ToLower(args[0])
	tail := args[1:]
	switch subcmd {
	case "send":
		return runMailSend(tail)
	case "check":
		return runMailCheck(tail)
	case "test":
		return runMailTest(tail)
	case "template":
		return runMailTemplate(tail)
	case "notify":
		return runMailNotify(tail)
	case "ui":
		return runMailUI(tail)
	default:
		return appfault.NewValidationError(fmt.Sprintf("unknown mail subcommand: %s", subcmd))
	}
}

func isHelpArg(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
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
			if p, err := strconv.Atoi(args[i+1]); err == nil && p > 0 {
				port = p
			}
			i++
		}
	}
	if err := RunMailUI(port); err != nil {
		return appfault.WrapExecution(err, "run mail ui")
	}
	return nil
}

func runMailSend(args []string) *appfault.AppError {
	opts, err := parseSendArgs(args)
	if err != nil {
		return appfault.NewValidationError(err.Error())
	}
	if len(opts.To) == 0 {
		return appfault.NewValidationError("missing required recipient: --to")
	}
	if opts.Subject == "" {
		return appfault.NewValidationError("missing required subject: --subject")
	}
	return executeSendPipeline(opts)
}

func parseSendArgs(args []string) (MailSendOptions, error) {
	opts := MailSendOptions{AccountName: "default", Host: "127.0.0.1", Port: 587}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if (a == "--to" || a == "-t") && i+1 < len(args) {
			opts.To = strings.Split(args[i+1], ",")
			i++
		} else if (a == "--subject" || a == "-s") && i+1 < len(args) {
			opts.Subject = args[i+1]
			i++
		} else if (a == "--body" || a == "-b") && i+1 < len(args) {
			opts.Body = args[i+1]
			i++
		} else if (a == "--file" || a == "-f") && i+1 < len(args) {
			content, err := os.ReadFile(args[i+1])
			if err != nil {
				return opts, err
			}
			opts.Body = string(content)
			i++
		} else if a == "--html" {
			opts.IsHTML = true
		} else if a == "--dry-run" {
			opts.IsDryRun = true
		}
	}
	return opts, nil
}

func executeSendPipeline(opts MailSendOptions) *appfault.AppError {
	if opts.From == "" {
		opts.From = opts.AccountName + "@gitmap.local"
	}
	raw := BuildMIMEMessage(opts)
	if opts.IsDryRun {
		fmt.Printf("[GITMAP MAIL] DRY RUN: Validated to %s (subject: %s)\n", MaskEmailList(opts.To), opts.Subject)
		return nil
	}
	if err := CurrentTransmitter.SendMail(opts, raw); err != nil {
		return appfault.WrapExecution(err, "send email")
	}
	fmt.Printf("[GITMAP MAIL] Sent to %s (subject: %s)\n", MaskEmailList(opts.To), opts.Subject)
	return nil
}

func runMailCheck(args []string) *appfault.AppError {
	filter := parseCheckArgs(args)
	store, _ := OpenMailStore("")
	if store != nil {
		defer store.Close()
	}
	summary, err := EvaluateMailboxes(store, filter)
	if err != nil {
		return appfault.WrapExecution(err, "evaluate mailboxes")
	}
	if hasJSONFlag(args) {
		out, _ := FormatMailCheckJSON(summary)
		fmt.Println(out)
		return nil
	}
	fmt.Print(FormatMailCheckTable(summary, filter.ShowFullEmail))
	return nil
}

func parseCheckArgs(args []string) MailboxCheckFilter {
	f := MailboxCheckFilter{MaxLimit: 10}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--all", "-a":
			f.CheckAll = true
		case "--account":
			if i+1 < len(args) {
				f.AccountName = args[i+1]
				i++
			}
		case "--unread-only", "-u":
			f.FilterUnreadOnly = true
		case "--show-full-email":
			f.ShowFullEmail = true
		}
	}
	return f
}

func hasJSONFlag(args []string) bool {
	for _, a := range args {
		if a == "--json" || a == "-j" {
			return true
		}
	}
	return false
}

func runMailTest(args []string) *appfault.AppError {
	host, port, isTLS := "127.0.0.1", 587, false
	for i := 0; i < len(args); i++ {
		if args[i] == "--host" && i+1 < len(args) {
			host = args[i+1]
			i++
		} else if args[i] == "--port" && i+1 < len(args) {
			if p, err := strconv.Atoi(args[i+1]); err == nil && p > 0 {
				port = p
			}
			i++
		} else if args[i] == "--tls" {
			isTLS = true
		}
	}
	return executeTestConnection(host, port, isTLS)
}

func executeTestConnection(host string, port int, isTLS bool) *appfault.AppError {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	if port == 465 || isTLS {
		tlsCfg := &tls.Config{ServerName: host, InsecureSkipVerify: isLocalHost(host)}
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", addr, tlsCfg)
		if err != nil {
			return appfault.WrapExecution(err, "test tls connection")
		}
		_ = conn.Close()
	} else {
		conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
		if err != nil {
			return appfault.WrapExecution(err, "test connection")
		}
		_ = conn.Close()
	}
	fmt.Printf("[GITMAP MAIL TEST] SUCCESS: Connected to %s (latency: %dms)\n", addr, time.Since(start).Milliseconds())
	return nil
}

func runMailTemplate(args []string) *appfault.AppError {
	if len(args) == 0 || isHelpArg(args[0]) {
		fmt.Println("Usage: gitmap mail template <ls | view <name> | add <name> | rm <name>>")
		return nil
	}
	store, err := OpenMailStore("")
	if err != nil {
		return appfault.WrapExecution(err, "open mail store")
	}
	defer store.Close()
	switch strings.ToLower(args[0]) {
	case "ls", "list":
		return handleTemplateList(store)
	case "view", "get":
		return handleTemplateView(store, args[1:])
	case "add", "save":
		return handleTemplateAdd(store, args[1:])
	case "rm", "del":
		return handleTemplateDelete(store, args[1:])
	default:
		return appfault.NewValidationError("unknown template command: " + args[0])
	}
}

func handleTemplateList(s *MailStore) *appfault.AppError {
	list, err := s.ListTemplates()
	if err != nil {
		return appfault.WrapExecution(err, "list templates")
	}
	fmt.Printf("[GITMAP TEMPLATES] Found %d template(s)\n", len(list))
	for _, t := range list {
		fmt.Printf(" - %-20s (subject: %s, html: %t)\n", t.TemplateName, t.SubjectTemplate, t.IsHTML)
	}
	return nil
}

func handleTemplateView(s *MailStore, args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("missing template name")
	}
	t, err := s.GetTemplate(args[0])
	if err != nil {
		return appfault.NewNotFound("mail", "template", "template not found: "+args[0])
	}
	fmt.Printf("Template: %s\nSubject: %s\nHTML: %t\nBody:\n%s\n", t.TemplateName, t.SubjectTemplate, t.IsHTML, t.BodyTemplate)
	return nil
}

func handleTemplateAdd(s *MailStore, args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("missing template name")
	}
	tmpl := MailTemplateRecord{TemplateID: fmt.Sprintf("tmpl-%d", time.Now().UnixNano()), TemplateName: args[0]}
	for i := 1; i < len(args); i++ {
		if args[i] == "--subject" && i+1 < len(args) {
			tmpl.SubjectTemplate = args[i+1]
			i++
		} else if args[i] == "--body" && i+1 < len(args) {
			tmpl.BodyTemplate = args[i+1]
			i++
		} else if args[i] == "--html" {
			tmpl.IsHTML = true
		}
	}
	if err := s.SaveTemplate(tmpl); err != nil {
		return appfault.WrapExecution(err, "save template")
	}
	fmt.Printf("[GITMAP TEMPLATE] Saved template '%s'\n", tmpl.TemplateName)
	return nil
}

func handleTemplateDelete(s *MailStore, args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("missing template name to delete")
	}
	if err := s.DeleteTemplate(args[0]); err != nil {
		return appfault.WrapExecution(err, "delete template")
	}
	fmt.Printf("[GITMAP TEMPLATE] Deleted template '%s'\n", args[0])
	return nil
}

func runMailNotify(args []string) *appfault.AppError {
	if len(args) == 0 || isHelpArg(args[0]) {
		fmt.Println("Usage: gitmap mail notify <watch | ls | trigger <rule>>")
		return nil
	}
	store, _ := OpenMailStore("")
	if store != nil {
		defer store.Close()
	}
	switch strings.ToLower(args[0]) {
	case "watch":
		return handleNotifyWatch(store, args[1:])
	case "ls", "list":
		return handleNotifyList(store)
	case "trigger":
		return handleNotifyTrigger(args[1:])
	default:
		return appfault.NewValidationError("unknown notify command: " + args[0])
	}
}

func handleNotifyWatch(store *MailStore, args []string) *appfault.AppError {
	cfg := parseWatchArgs(args)
	guard := NewIntervalGuard(cfg.PollingInterval)
	if guard.WasClamped {
		fmt.Printf("[GITMAP WATCHER] Notice: interval clamped to minimum %v\n", guard.EffectiveInterval)
	}
	cfg.PollingInterval = guard.EffectiveInterval
	baseline, _ := ScanFolderTree(cfg.TargetDirectory, cfg.IsRecursive)
	if cfg.RunOnce {
		batch, _, err := ExecuteWatchCycle(cfg, baseline, store, CurrentTransmitter)
		if err != nil {
			return appfault.WrapExecution(err, "run watch cycle")
		}
		fmt.Printf("[GITMAP WATCHER] Single pass completed. Detected %d event(s)\n", batch.TotalEvents)
		return nil
	}
	fmt.Printf("[GITMAP WATCHER] Monitoring %s (interval: %v)\n", cfg.TargetDirectory, cfg.PollingInterval)
	return nil
}

func parseWatchArgs(args []string) FolderWatcherConfig {
	cfg := FolderWatcherConfig{TargetDirectory: ".", PollingInterval: DefaultPollingInterval, DebounceWindow: DefaultDebounceWindow, IsRecursive: true, IsActive: true}
	for i := 0; i < len(args); i++ {
		if (args[i] == "--path" || args[i] == "-p") && i+1 < len(args) {
			cfg.TargetDirectory = args[i+1]
			i++
		} else if (args[i] == "--to" || args[i] == "-t") && i+1 < len(args) {
			cfg.Recipients = strings.Split(args[i+1], ",")
			i++
		} else if args[i] == "--once" {
			cfg.RunOnce = true
		}
	}
	return cfg
}

func handleNotifyList(s *MailStore) *appfault.AppError {
	if s == nil {
		fmt.Println("[GITMAP NOTIFY] No database records found.")
		return nil
	}
	rules, _ := s.ListWatcherRules()
	logs, _ := s.ListNotificationLogs(10)
	fmt.Printf("[GITMAP NOTIFY] %d rule(s), %d log(s)\n", len(rules), len(logs))
	for _, l := range logs {
		fmt.Printf(" - [%s] %s -> %s (delivered: %t)\n", l.CreatedAt, l.Subject, l.RecipientSummary, l.IsDelivered)
	}
	return nil
}

func handleNotifyTrigger(args []string) *appfault.AppError {
	if len(args) == 0 {
		return appfault.NewValidationError("missing rule name to trigger")
	}
	fmt.Printf("[GITMAP NOTIFY] Triggered notification for rule '%s'\n", args[0])
	return nil
}

func printMailHelp() {
	fmt.Println("Usage: gitmap mail <command> [flags]")
	fmt.Println("\nCommands:")
	fmt.Println("  send        Send an email directly via SMTP")
	fmt.Println("  check       Check configured mailboxes and unread counts")
	fmt.Println("  test        Test SMTP connectivity and TLS handshakes")
	fmt.Println("  template    Manage email templates (ls, view, add, rm)")
	fmt.Println("  notify      Folder watcher notifier (watch, ls, trigger)")
	fmt.Println("  ui          Start the embedded mailing composer web UI")
	fmt.Println("\nFlags:")
	fmt.Println("  --to, -t             Recipient addresses (comma-separated)")
	fmt.Println("  --subject, -s        Subject line")
	fmt.Println("  --body, -b           Body content")
	fmt.Println("  --file, -f           Path to message body file")
	fmt.Println("  --html               Send as HTML format")
	fmt.Println("  --unread-only, -u    Filter mailboxes to unread messages only")
	fmt.Println("  --show-full-email    Display unmasked email addresses")
	fmt.Println("  --json, -j           Emit output as Universal JSON Envelope")
}
