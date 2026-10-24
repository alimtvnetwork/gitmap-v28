// Package cmdschedule manages scheduled tasks, privacy masking, and mailbox inspection.
package cmdschedule

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type scheduleEmailOptions struct {
	Name          string `json:"name"`
	To            string `json:"to"`
	Subject       string `json:"subject"`
	Body          string `json:"body,omitempty"`
	File          string `json:"file,omitempty"`
	Every         string `json:"every"`
	ShowFullEmail bool   `json:"showFullEmail"`
	IsMasked      bool   `json:"isMasked"`
	IsJSON        bool   `json:"isJson"`
}

// RunScheduleEmail registers a recurring scheduled email dispatch in scheduler_tasks.
func RunScheduleEmail(args []string) *appfault.AppError {
	if isScheduleHelpRequested(args) {
		renderScheduleEmailHelp()
		return nil
	}

	opts, valErr := parseScheduleEmailFlags(args)
	if valErr != nil {
		return valErr
	}

	task, regErr := registerScheduledEmailTask(opts)
	if regErr != nil {
		return regErr
	}

	if opts.IsJSON {
		return renderScheduleEmailJSON(task, opts)
	}

	renderScheduleEmailSummary(opts)
	return nil
}

func parseScheduleEmailFlags(args []string) (scheduleEmailOptions, *appfault.AppError) {
	var opts scheduleEmailOptions
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case slices.Contains([]string{"--to", "-t"}, arg):
			opts.To, i = parseNextFlagValue(args, i, arg)
		case strings.HasPrefix(arg, "--to="):
			opts.To = strings.TrimPrefix(arg, "--to=")
		case slices.Contains([]string{"--subject", "-s"}, arg):
			opts.Subject, i = parseNextFlagValue(args, i, arg)
		case strings.HasPrefix(arg, "--subject="):
			opts.Subject = strings.TrimPrefix(arg, "--subject=")
		case slices.Contains([]string{"--body", "-b"}, arg):
			opts.Body, i = parseNextFlagValue(args, i, arg)
		case strings.HasPrefix(arg, "--body="):
			opts.Body = strings.TrimPrefix(arg, "--body=")
		case slices.Contains([]string{"--file", "-f"}, arg):
			opts.File, i = parseNextFlagValue(args, i, arg)
		case strings.HasPrefix(arg, "--file="):
			opts.File = strings.TrimPrefix(arg, "--file=")
		case slices.Contains([]string{"--every", "-i", "--interval"}, arg):
			opts.Every, i = parseNextFlagValue(args, i, arg)
		case strings.HasPrefix(arg, "--every="):
			opts.Every = strings.TrimPrefix(arg, "--every=")
		case arg == "--show-full-email":
			opts.ShowFullEmail = true
		case slices.Contains([]string{"--json", "-j"}, arg):
			opts.IsJSON = true
		case !strings.HasPrefix(arg, "-") && opts.Name == "":
			opts.Name = arg
		}
	}

	opts.IsMasked = ResolveIsMasked(opts.ShowFullEmail)
	return opts, validateScheduleEmailOptions(opts)
}

func validateScheduleEmailOptions(opts scheduleEmailOptions) *appfault.AppError {
	if fieldErr := validateRequiredScheduleFields(opts); fieldErr != nil {
		return fieldErr
	}

	return validateScheduleEveryDuration(opts.Every)
}

func validateRequiredScheduleFields(opts scheduleEmailOptions) *appfault.AppError {
	if strings.TrimSpace(opts.To) == "" {
		return appfault.NewValidationError("missing required flag: --to <recipient>")
	}
	if strings.TrimSpace(opts.Subject) == "" {
		return appfault.NewValidationError("missing required flag: --subject <subject>")
	}
	if strings.TrimSpace(opts.Every) == "" {
		return appfault.NewValidationError("missing required flag: --every <interval> (e.g. 1h, 24h)")
	}

	return nil
}

func validateScheduleEveryDuration(every string) *appfault.AppError {
	durRes := ParseScheduleDuration(every)
	if durRes.IsFailed() || durRes.Value <= 0 {
		return appfault.NewValidationError(fmt.Sprintf("invalid schedule interval: %s", every))
	}

	return nil
}

func registerScheduledEmailTask(opts scheduleEmailOptions) (store.SchedulerTask, *appfault.AppError) {
	db, err := store.OpenDefault()
	if err != nil {
		return store.SchedulerTask{}, appfault.WrapExecution(err, "open scheduler database")
	}
	defer db.Close()

	taskName := opts.Name
	if taskName == "" {
		taskName = fmt.Sprintf("email-%d", time.Now().Unix())
	}

	slug := store.ScheduleSlug(taskName)
	cmdLine := buildEmailCommandLine(opts)
	task := store.SchedulerTask{
		Name:        taskName,
		Slug:        slug,
		DBPath:      store.ScheduleDBPath(slug),
		CommandLine: cmdLine,
		IntervalVal: opts.Every,
		IsEnabled:   true,
		IsScheduled: true,
	}

	if insertErr := db.InsertSchedule(task); insertErr != nil {
		return store.SchedulerTask{}, appfault.WrapExecution(insertErr, "insert scheduled email task")
	}

	_ = syncScheduleSplitDBConfig(task)
	return task, nil
}

func buildEmailCommandLine(opts scheduleEmailOptions) string {
	parts := []string{"gitmap mail send", fmt.Sprintf("--to %q", opts.To), fmt.Sprintf("--subject %q", opts.Subject)}
	if opts.Body != "" {
		parts = append(parts, fmt.Sprintf("--body %q", opts.Body))
	}
	if opts.File != "" {
		parts = append(parts, fmt.Sprintf("--file %q", opts.File))
	}

	return strings.Join(parts, " ")
}

func renderScheduleEmailSummary(opts scheduleEmailOptions) {
	notifierDisplay := resolveRecipientDisplay(opts.To, opts.IsMasked)
	taskName := opts.Name
	if taskName == "" {
		taskName = "(auto-generated)"
	}

	fmt.Printf("  ✓ Scheduled email task %q registered\n", taskName)
	fmt.Printf("    Notifier: %s\n", notifierDisplay)
	fmt.Printf("    Interval: %s\n", opts.Every)
	fmt.Printf("    Subject:  %s\n", opts.Subject)
}

func resolveRecipientDisplay(rawTo string, isMasked bool) string {
	if !isMasked {
		return rawTo
	}

	if strings.Contains(rawTo, ",") {
		recipients := strings.Split(rawTo, ",")
		trimmed := make([]string, len(recipients))
		for i, r := range recipients {
			trimmed[i] = strings.TrimSpace(r)
		}
		return strings.Join(MaskEmailList(trimmed), ", ")
	}

	return MaskEmail(rawTo)
}

func renderScheduleEmailJSON(task store.SchedulerTask, opts scheduleEmailOptions) *appfault.AppError {
	payload := map[string]any{
		"Status": "success",
		"Attributes": map[string]any{
			"Command":    "schedule-email",
			"ExecutedAt": time.Now().UTC().Format(time.RFC3339),
			"IsMasked":   opts.IsMasked,
		},
		"Task": map[string]any{
			"Name":        task.Name,
			"Slug":        task.Slug,
			"Notifier":    resolveRecipientDisplay(opts.To, opts.IsMasked),
			"Interval":    task.IntervalVal,
			"CommandLine": task.CommandLine,
			"IsEnabled":   task.IsEnabled,
		},
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return appfault.WrapExecution(err, "marshal schedule email json")
	}

	fmt.Println(string(data))
	return nil
}

func renderScheduleEmailHelp() {
	fmt.Println("Usage: gitmap schedule email [name] [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --to, -t string         Recipient email address (required)")
	fmt.Println("  --subject, -s string    Email subject line (required)")
	fmt.Println("  --every, -i string      Execution interval, e.g. 1h, 24h, 30m (required)")
	fmt.Println("  --body, -b string       Inline body text")
	fmt.Println("  --file, -f string       Path to email body or template file")
	fmt.Println("  --show-full-email       Display unmasked recipient PII in summary")
	fmt.Println("  --json, -j              Emit output in Universal JSON format")
}
