package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// ErrNoPrompt is returned when --no-prompt is set but a value is needed.
// Callers map this to exit code CommitInExitMissingAnswer.
var ErrNoPrompt = errors.New("no-prompt set; cannot ask")

// Asker abstracts terminal I/O for hermetic tests.
type Asker struct {
	NoPrompt bool
	In       io.Reader
	Out      io.Writer
	reader   *bufio.Reader
}

// New builds an Asker bound to stdin/stderr (stderr so prompts never
// pollute machine-readable stdout).
func New(noPrompt bool) *Asker {
	return &Asker{NoPrompt: noPrompt, In: os.Stdin, Out: os.Stderr}
}

// AskString returns the user's answer or ErrNoPrompt under --no-prompt.
// `field` names the missing setting; it is interpolated into the
// standardized error format so the user knows what was unset.
func (a *Asker) AskString(field, question, fallback string) (string, error) {
	if a.NoPrompt {
		fmt.Fprintf(a.Out, constants.CommitInErrMissingAnswer, field)
		return "", ErrNoPrompt
	}
	if a.reader == nil {
		a.reader = bufio.NewReader(a.In)
	}
	if fallback != "" {
		fmt.Fprintf(a.Out, "%s [%s]: ", question, fallback)
	} else {
		fmt.Fprintf(a.Out, "%s: ", question)
	}
	line, err := a.reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("read answer: %w", err)
	}
	line = strings.TrimRight(line, "\r\n")
	if line == "" {
		return fallback, nil
	}
	return line, nil
}

// AskEnum loops until the answer is in `valid` (case-sensitive). Empty
// answer with non-empty fallback returns the fallback unchecked.
func (a *Asker) AskEnum(field, question string, valid []string, fallback string) (string, error) {
	for {
		ans, err := a.AskString(field, question+" ("+strings.Join(valid, "|")+")", fallback)
		if err != nil {
			return "", err
		}
		if isMember(ans, valid) {
			return ans, nil
		}
		fmt.Fprintf(a.Out, "  invalid: must be one of %s\n", strings.Join(valid, "|"))
	}
}

func isMember(v string, set []string) bool {
	for _, s := range set {
		if s == v {
			return true
		}
	}
	return false
}
