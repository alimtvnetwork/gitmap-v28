package cmdspace

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/linkexpander"
)

// send-email flag names.
const (
	flagSendEmailTo               = "--to"
	flagSendEmailSubject          = "--subject"
	flagSendEmailBody             = "--body"
	flagSendEmailFile             = "--file"
	flagSendEmailHTML             = "--html"
	flagSendEmailExpand           = "--expand-links"
	flagSendEmailExpandShort      = "--expand"
	flagSendEmailInstruction      = "--instruction"
	flagSendEmailInstructionShort = "-i"
	flagSendEmailJSON             = "--json"
	flagSendEmailHelp             = "--help"
	flagSendEmailHelpShort        = "-h"
)

// sendEmailOptions holds the parsed `space send-email` flags.
type sendEmailOptions struct {
	To              string
	Subject         string
	Body            string
	File            string
	InstructionFile string
	IsHTML          bool
	IsExpandLinks   bool
	IsJSON          bool
	HasInstruction  bool
}

// RunSendEmail implements `gitmap space send-email --to X [--subject S]
// [--body B] [--file F] [--html]`: it parses an email file when given,
// merges explicit flags over it, then delegates to the go-email-reader
// binary when installed or sends through the built-in SMTP sender.
func RunSendEmail(args []string) error {
	if hasHelpFlag(args) {
		printSendEmailHelp()

		return nil
	}
	opts, parseErr := parseSendEmailArgs(args)
	if parseErr != nil {
		return parseErr
	}
	if opts.HasInstruction {
		return runSendEmailInstruction(opts)
	}

	return runStandardSendEmail(opts)
}

func runStandardSendEmail(opts sendEmailOptions) error {
	msg, buildErr := buildEmailMessage(opts)
	if buildErr != nil {
		return buildErr
	}
	if bin := resolveGoEmailReaderBin(); len(bin) > 0 {
		return delegateToGoEmailReader(bin, opts, msg)
	}

	return sendEmailWithJSONSupport(opts, msg)
}

func runSendEmailInstruction(opts sendEmailOptions) error {
	cfg, parseErr := parseInstructionFileConfig(opts.InstructionFile)
	if parseErr != nil {
		return parseErr
	}
	applyInstructionAmbientOverrides(&cfg, opts)
	res, execErr := linkexpander.ExecuteInstruction(cfg)
	if execErr != nil {
		return execErr
	}

	return outputInstructionResult(opts, res)
}

func parseInstructionFileConfig(path string) (linkexpander.InstructionConfig, error) {
	cfg := linkexpander.DefaultInstructionConfig()
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return cfg, apperror.WrapSimple(readErr, "read instruction file:")
	}
	if unmarshalErr := json.Unmarshal(data, &cfg); unmarshalErr != nil {
		return cfg, apperror.WrapSimple(unmarshalErr, "parse instruction JSON:")
	}

	return cfg, nil
}

func applyInstructionAmbientOverrides(cfg *linkexpander.InstructionConfig, opts sendEmailOptions) {
	for i := range cfg.Links {
		if len(cfg.Links[i].Recipients) == 0 && len(opts.To) > 0 {
			cfg.Links[i].Recipients = splitAddresses(opts.To)
		}
		if len(cfg.Links[i].Subject) == 0 && len(opts.Subject) > 0 {
			cfg.Links[i].Subject = opts.Subject
		}
	}
}

func outputInstructionResult(opts sendEmailOptions, res linkexpander.InstructionResult) error {
	if opts.IsJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(res)
	}

	return printInstructionSummaryText(res)
}

func printInstructionSummaryText(res linkexpander.InstructionResult) error {
	fmt.Printf("Instruction executed: %d total links, %d expanded, %d sent\n", res.TotalLinks, res.ExpandedCount, res.SentCount)
	for _, item := range res.Results {
		printInstructionItemSummary(item)
	}

	return nil
}

func printInstructionItemSummary(item linkexpander.InstructionLinkReport) {
	if item.IsExpanded {
		fmt.Printf("  %s -> %s (hops: %d)\n", item.OriginalURL, item.ExpandedURL, item.RedirectHops)
		return
	}
	fmt.Printf("  %s (unchanged)\n", item.OriginalURL)
}

func sendEmailWithJSONSupport(opts sendEmailOptions, msg EmailMessage) error {
	sendErr := sendEmailViaSMTP(msg)
	if !opts.IsJSON {
		return sendErr
	}

	return printSendEmailJSON(opts, msg, sendErr)
}

func printSendEmailJSON(opts sendEmailOptions, msg EmailMessage, sendErr error) error {
	result := linkexpander.SendEmailJSONResult{
		IsSuccess:     sendErr == nil,
		To:            msg.To,
		Subject:       msg.Subject,
		IsHTML:        msg.IsHTML,
		IsExpanded:    msg.IsExpanded,
		ExpandedCount: msg.ExpandedCount,
		Links:         msg.ExpandedLinks,
	}
	if sendErr != nil {
		result.Error = sendErr.Error()
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(result)
}

// parseSendEmailArgs parses the send-email flags; unknown dash flags and
// positional arguments are errors.
func parseSendEmailArgs(args []string) (sendEmailOptions, error) {
	var opts sendEmailOptions

	for i := 0; i < len(args); i++ {
		consumed, stepErr := applyNextSendEmailArg(args, i, &opts)
		if stepErr != nil {
			return opts, stepErr
		}
		i += consumed
	}

	return opts, nil
}

// applyNextSendEmailArg applies one argument: unknown dash flags are
// rejected before any value resolution, positional args are errors.
func applyNextSendEmailArg(args []string, i int, opts *sendEmailOptions) (int, error) {
	arg := args[i]
	if !strings.HasPrefix(arg, "-") {
		return 0, apperror.NewValidationError(fmt.Sprintf("unexpected positional argument %q: usage: gitmap space send-email --to <addr> [--subject <s>] [--body <b>] [--file <f>] [--html]", arg))
	}

	name, _, _ := splitFlagEquals(arg)
	if !isKnownSendEmailFlag(name) {
		return 0, apperror.NewValidationError(fmt.Sprintf("unknown flag: %s", name))
	}

	return applySendEmailFlag(args, i, opts)
}

// isKnownSendEmailFlag reports whether name is a recognized send-email flag.
func isKnownSendEmailFlag(name string) bool {
	switch name {
	case flagSendEmailTo, flagSendEmailSubject, flagSendEmailBody, flagSendEmailFile, flagSendEmailHTML,
		flagSendEmailExpand, flagSendEmailExpandShort, flagSendEmailInstruction, flagSendEmailInstructionShort,
		flagSendEmailJSON, flagSendEmailHelp, flagSendEmailHelpShort:
		return true
	}

	return false
}

// applySendEmailFlag applies one dash-flag; it returns how many extra args
// the flag consumed (1 for `--flag value`, 0 for `--flag=value` or bools).
func applySendEmailFlag(args []string, i int, opts *sendEmailOptions) (int, error) {
	name, inlineVal, hasInline := splitFlagEquals(args[i])

	if name == flagSendEmailHTML {
		return applyHTMLFlag(hasInline, opts)
	}

	if name == flagSendEmailExpand || name == flagSendEmailExpandShort {
		return applyExpandLinksFlag(hasInline, opts)
	}

	if name == flagSendEmailJSON {
		return applyJSONFlag(hasInline, opts)
	}

	val, extra, valErr := resolveFlagValue(args, i, inlineVal, hasInline, name)
	if valErr != nil {
		return 0, valErr
	}

	if assignErr := assignSendEmailFlag(name, val, opts); assignErr != nil {
		return 0, assignErr
	}

	return extra, nil
}

// splitFlagEquals splits --name=value into its name and value parts.
func splitFlagEquals(a string) (name, val string, hasVal bool) {
	name, val, hasVal = strings.Cut(a, "=")

	return name, val, hasVal
}

// applyHTMLFlag sets the --html bool; it takes no value.
func applyHTMLFlag(hasInline bool, opts *sendEmailOptions) (int, error) {
	if hasInline {
		return 0, apperror.NewValidationError("flag --html takes no value")
	}
	opts.IsHTML = true

	return 0, nil
}

func applyExpandLinksFlag(hasInline bool, opts *sendEmailOptions) (int, error) {
	if hasInline {
		return 0, apperror.NewValidationError("flag --expand-links takes no value")
	}
	opts.IsExpandLinks = true

	return 0, nil
}

func applyJSONFlag(hasInline bool, opts *sendEmailOptions) (int, error) {
	if hasInline {
		return 0, apperror.NewValidationError("flag --json takes no value")
	}
	opts.IsJSON = true

	return 0, nil
}

// resolveFlagValue returns the value for a value-flag from --name=value
// or from the next argument.
func resolveFlagValue(args []string, i int, inlineVal string, hasInline bool, name string) (string, int, error) {
	if hasInline {
		return inlineVal, 0, nil
	}

	if i+1 >= len(args) {
		return "", 0, apperror.NewValidationError(fmt.Sprintf("flag %s requires a value", name))
	}

	return args[i+1], 1, nil
}

// assignSendEmailFlag stores a value flag; unknown names are errors.
func assignSendEmailFlag(name, val string, opts *sendEmailOptions) error {
	switch name {
	case flagSendEmailTo:
		opts.To = val
	case flagSendEmailSubject:
		opts.Subject = val
	case flagSendEmailBody:
		opts.Body = val
	case flagSendEmailFile:
		opts.File = val
	case flagSendEmailInstruction, flagSendEmailInstructionShort:
		opts.InstructionFile = val
		opts.HasInstruction = len(strings.TrimSpace(val)) > 0
	default:
		return apperror.NewValidationError(fmt.Sprintf("unknown flag: %s", name))
	}

	return nil
}

// buildEmailMessage parses --file when given, then lets explicit CLI flags
// override the file values, and validates the merged message.
func buildEmailMessage(opts sendEmailOptions) (EmailMessage, error) {
	msg, parseErr := parseOptionalEmailFile(opts.File)
	if parseErr != nil {
		return msg, parseErr
	}
	applySendEmailOverrides(opts, &msg)
	applyLinkExpansion(&msg, opts.IsExpandLinks)

	if validateErr := validateEmailMessage(msg); validateErr != nil {
		return msg, validateErr
	}

	return msg, nil
}

func parseOptionalEmailFile(filePath string) (EmailMessage, error) {
	if len(filePath) == 0 {
		return EmailMessage{}, nil
	}
	return parseEmailFile(filePath)
}

// applySendEmailOverrides lets explicit CLI flags win over --file values.
func applySendEmailOverrides(opts sendEmailOptions, msg *EmailMessage) {
	if len(opts.To) > 0 {
		msg.To = splitAddresses(opts.To)
	}
	if len(opts.Subject) > 0 {
		msg.Subject = opts.Subject
	}
	if len(opts.Body) > 0 {
		msg.Body = opts.Body
	}
	if opts.IsHTML {
		msg.IsHTML = true
	}
}

func applyLinkExpansion(msg *EmailMessage, isExpandLinks bool) {
	if !isExpandLinks {
		return
	}
	msg.Subject, _ = linkexpander.ExpandText(msg.Subject)
	expandedBody, links := linkexpander.ExpandText(msg.Body)
	msg.Body = expandedBody
	msg.IsExpanded = true
	msg.ExpandedCount = len(links)
	msg.ExpandedLinks = links
}

// validateEmailMessage requires at least one recipient and a subject.
func validateEmailMessage(msg EmailMessage) error {
	if len(msg.To) == 0 {
		return apperror.NewSimple("missing --to: usage: gitmap space send-email --to <addr[,addr]> --subject <s> --body <b>", "E_SEND_EMAIL_NO_TO")
	}

	if len(strings.TrimSpace(msg.Subject)) == 0 {
		return apperror.NewSimple("missing subject: pass --subject or provide it in --file", "E_SEND_EMAIL_NO_SUBJECT")
	}

	return nil
}

// delegateToGoEmailReader execs `go-email-reader send` with inherited stdio.
// Passwords are never passed: go-email-reader reads its own encrypted store.
func delegateToGoEmailReader(bin string, opts sendEmailOptions, msg EmailMessage) error {
	cmd := exec.Command(bin, delegateSendEmailArgv(opts, msg)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if runErr := cmd.Run(); runErr != nil {
		return apperror.WrapSimple(runErr, "go-email-reader send:")
	}

	return nil
}

// delegateSendEmailArgv builds the `go-email-reader send` argv.
func delegateSendEmailArgv(opts sendEmailOptions, msg EmailMessage) []string {
	argv := []string{"send", "--to", strings.Join(msg.To, ",")}
	argv = appendValueArg(argv, "--subject", msg.Subject)
	argv = appendValueArg(argv, "--body", msg.Body)
	argv = appendValueArg(argv, "--file", opts.File)
	if opts.IsHTML {
		argv = append(argv, "--html")
	}
	if opts.IsExpandLinks {
		argv = append(argv, "--expand-links")
	}
	if opts.IsJSON {
		argv = append(argv, "--json")
	}

	return argv
}

// appendValueArg appends --name value when the value is non-empty.
func appendValueArg(argv []string, name, val string) []string {
	if len(val) == 0 {
		return argv
	}

	return append(argv, name, val)
}
