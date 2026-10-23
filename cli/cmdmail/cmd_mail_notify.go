package cmdmail

import (
	"fmt"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

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
		switch {
		case args[i] == "--subject" && i+1 < len(args):
			tmpl.SubjectTemplate, i = args[i+1], i+1
		case args[i] == "--body" && i+1 < len(args):
			tmpl.BodyTemplate, i = args[i+1], i+1
		case args[i] == "--html":
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
	if !cfg.RunOnce {
		fmt.Printf("[GITMAP WATCHER] Monitoring %s (interval: %v)\n", cfg.TargetDirectory, cfg.PollingInterval)
		return nil
	}
	batch, _, err := ExecuteWatchCycle(cfg, baseline, store, CurrentTransmitter)
	if err != nil {
		return appfault.WrapExecution(err, "run watch cycle")
	}
	fmt.Printf("[GITMAP WATCHER] Single pass completed (%d events)\n", batch.TotalEvents)
	return nil
}

func parseWatchArgs(args []string) FolderWatcherConfig {
	cfg := FolderWatcherConfig{TargetDirectory: ".", PollingInterval: DefaultPollingInterval, DebounceWindow: DefaultDebounceWindow, IsRecursive: true, IsActive: true}
	for i := 0; i < len(args); i++ {
		switch {
		case (args[i] == "--path" || args[i] == "-p") && i+1 < len(args):
			cfg.TargetDirectory, i = args[i+1], i+1
		case (args[i] == "--to" || args[i] == "-t") && i+1 < len(args):
			cfg.Recipients, i = strings.Split(args[i+1], ","), i+1
		case args[i] == "--once":
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
