package cmdspace

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// send-email flag names.
const (
	flagSendEmailTo        = "--to"
	flagSendEmailSubject   = "--subject"
	flagSendEmailBody      = "--body"
	flagSendEmailFile      = "--file"
	flagSendEmailHTML      = "--html"
	flagSendEmailHelp      = "--help"
	flagSendEmailHelpShort = "-h"
)

// sendEmailOptions holds the parsed `space send-email` flags.
type sendEmailOptions struct {
	To      string
	Subject string
	Body    string
	File    string
	IsHTML  bool
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

	msg, buildErr := buildEmailMessage(opts)
	if buildErr != nil {
		return buildErr
	}

	if bin := resolveGoEmailReaderBin(); len(bin) > 0 {
		return delegateToGoEmailReader(bin, opts, msg)
	}

	return sendEmailViaSMTP(msg)
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
		return 0, fmt.Errorf("unexpected positional argument %q: usage: gitmap space send-email --to <addr> [--subject <s>] [--body <b>] [--file <f>] [--html]", arg)
	}

	name, _, _ := splitFlagEquals(arg)
	if !isKnownSendEmailFlag(name) {
		return 0, fmt.Errorf("unknown flag: %s", name)
	}

	return applySendEmailFlag(args, i, opts)
}

// isKnownSendEmailFlag reports whether name is a recognized send-email flag.
func isKnownSendEmailFlag(name string) bool {
	switch name {
	case flagSendEmailTo, flagSendEmailSubject, flagSendEmailBody, flagSendEmailFile, flagSendEmailHTML, flagSendEmailHelp, flagSendEmailHelpShort:
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
		return 0, fmt.Errorf("flag --html takes no value")
	}
	opts.IsHTML = true

	return 0, nil
}

// resolveFlagValue returns the value for a value-flag from --name=value
// or from the next argument.
func resolveFlagValue(args []string, i int, inlineVal string, hasInline bool, name string) (string, int, error) {
	if hasInline {
		return inlineVal, 0, nil
	}

	if i+1 >= len(args) {
		return "", 0, fmt.Errorf("flag %s requires a value", name)
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
	default:
		return fmt.Errorf("unknown flag: %s", name)
	}

	return nil
}

// buildEmailMessage parses --file when given, then lets explicit CLI flags
// override the file values, and validates the merged message.
func buildEmailMessage(opts sendEmailOptions) (EmailMessage, error) {
	msg := EmailMessage{}

	if len(opts.File) > 0 {
		parsed, parseErr := parseEmailFile(opts.File)
		if parseErr != nil {
			return msg, parseErr
		}
		msg = parsed
	}
	applySendEmailOverrides(opts, &msg)

	if validateErr := validateEmailMessage(msg); validateErr != nil {
		return msg, validateErr
	}

	return msg, nil
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

// resolveGoEmailReaderBin locates the go-email-reader binary: first on
// PATH, then in ~/go/bin, then via the go_email_reader_bin gitmap
// setting. Empty when no usable binary is found.
func resolveGoEmailReaderBin() string {
	if path, lookErr := exec.LookPath("go-email-reader"); lookErr == nil {
		return path
	}

	if home, homeErr := os.UserHomeDir(); homeErr == nil {
		candidate := filepath.Join(home, "go", "bin", "go-email-reader")
		if isExecutableFile(candidate) {
			return candidate
		}
	}

	setting, _ := config.GetVariable("global", constants.SettingGoEmailReaderBin)
	if isExecutableFile(setting) {
		return setting
	}

	return ""
}

// isExecutableFile reports whether path is an executable regular file.
func isExecutableFile(path string) bool {
	if len(path) == 0 {
		return false
	}

	info, statErr := os.Stat(path)
	if statErr != nil || info.IsDir() {
		return false
	}

	return info.Mode().Perm()&0111 != 0
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
	argv = appendValueArg(argv, "--subject", opts.Subject)
	argv = appendValueArg(argv, "--body", opts.Body)
	argv = appendValueArg(argv, "--file", opts.File)
	if opts.IsHTML {
		argv = append(argv, "--html")
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
