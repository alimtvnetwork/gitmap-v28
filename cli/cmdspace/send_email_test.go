package cmdspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeTempEmailFile writes content to dir/name and returns the path.
func writeTempEmailFile(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	if writeErr := os.WriteFile(path, []byte(content), 0600); writeErr != nil {
		t.Fatal(writeErr)
	}

	return path
}

func TestParseSendEmailArgsTable(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		want    sendEmailOptions
		wantErr string
	}{
		{
			name: "space separated",
			args: []string{"--to", "a@x.com", "--subject", "Hi", "--body", "Hello"},
			want: sendEmailOptions{To: "a@x.com", Subject: "Hi", Body: "Hello"},
		},
		{
			name: "equals form",
			args: []string{"--to=a@x.com", "--subject=Hi"},
			want: sendEmailOptions{To: "a@x.com", Subject: "Hi"},
		},
		{
			name: "html bool",
			args: []string{"--to", "a@x.com", "--subject", "Hi", "--html"},
			want: sendEmailOptions{To: "a@x.com", Subject: "Hi", IsHTML: true},
		},
		{
			name:    "html rejects value",
			args:    []string{"--to", "a@x.com", "--html=x"},
			wantErr: "takes no value",
		},
		{
			name:    "unknown flag first",
			args:    []string{"--bogus", "--to", "a@x.com"},
			wantErr: "unknown flag: --bogus",
		},
		{
			name:    "trailing unknown flag rejected before value resolution",
			args:    []string{"--to", "a@x.com", "--subject", "Hi", "--bogus"},
			wantErr: "unknown flag: --bogus",
		},
		{
			name:    "equals unknown flag",
			args:    []string{"--to", "a@x.com", "--bogus=x"},
			wantErr: "unknown flag: --bogus",
		},
		{
			name:    "positional rejected",
			args:    []string{"--to", "a@x.com", "oops"},
			wantErr: "unexpected positional argument",
		},
		{
			name:    "missing value",
			args:    []string{"--to"},
			wantErr: "requires a value",
		},
		{
			name: "file flag",
			args: []string{"--file", "note.yaml", "--to", "b@x.com"},
			want: sendEmailOptions{File: "note.yaml", To: "b@x.com"},
		},
		{
			name: "expand links and json bool flags",
			args: []string{"--to", "a@x.com", "--subject", "Hi", "--expand-links", "--json"},
			want: sendEmailOptions{To: "a@x.com", Subject: "Hi", IsExpandLinks: true, IsJSON: true},
		},
		{
			name: "expand short flag",
			args: []string{"--to", "a@x.com", "--subject", "Hi", "--expand"},
			want: sendEmailOptions{To: "a@x.com", Subject: "Hi", IsExpandLinks: true},
		},
		{
			name: "instruction flag",
			args: []string{"--instruction", "job.json", "--json"},
			want: sendEmailOptions{InstructionFile: "job.json", HasInstruction: true, IsJSON: true},
		},
		{
			name: "instruction short flag",
			args: []string{"-i", "job.json"},
			want: sendEmailOptions{InstructionFile: "job.json", HasInstruction: true},
		},
		{
			name: "instruction equals form",
			args: []string{"--instruction=job.json"},
			want: sendEmailOptions{InstructionFile: "job.json", HasInstruction: true},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts, err := parseSendEmailArgs(tc.args)
			if len(tc.wantErr) > 0 {
				assertParseEmailArgsError(t, tc.wantErr, err)
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if opts != tc.want {
				t.Fatalf("got %+v, want %+v", opts, tc.want)
			}
		})
	}
}

func assertParseEmailArgsError(t *testing.T, wantErr string, err error) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), wantErr) {
		t.Fatalf("expected error containing %q, got %v", wantErr, err)
	}
}

func TestRunSendEmailMissingTo(t *testing.T) {
	if err := RunSendEmail([]string{"--subject", "Hi"}); err == nil ||
		!strings.Contains(err.Error(), "missing --to") {
		t.Fatalf("expected missing --to error, got %v", err)
	}
}

func TestRunSendEmailMissingSubject(t *testing.T) {
	if err := RunSendEmail([]string{"--to", "a@x.com"}); err == nil ||
		!strings.Contains(err.Error(), "missing subject") {
		t.Fatalf("expected missing subject error, got %v", err)
	}
}

func TestRunSendEmailUnknownFlag(t *testing.T) {
	err := RunSendEmail([]string{"--to", "a@x.com", "--subject", "Hi", "--nope"})
	if err == nil || !strings.Contains(err.Error(), "unknown flag: --nope") {
		t.Fatalf("expected unknown flag error, got %v", err)
	}
}

func TestParseEmailFileYAML(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.yaml",
		"to: a@x.com, b@x.com\ncc:\n  - c@x.com\nsubject:  Hello \nbody: YAML body\nhtml: true\n")

	msg, err := parseEmailFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(msg.To) != 2 || msg.To[0] != "a@x.com" || msg.To[1] != "b@x.com" {
		t.Fatalf("to wrong: %v", msg.To)
	}
	if len(msg.Cc) != 1 || msg.Cc[0] != "c@x.com" {
		t.Fatalf("cc wrong: %v", msg.Cc)
	}
	if msg.Subject != "Hello" || msg.Body != "YAML body" || !msg.IsHTML {
		t.Fatalf("content wrong: %+v", msg)
	}
}

func TestParseEmailFileYAMLListTo(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.yml", "to:\n  - a@x.com\n  - b@x.com\nsubject: S\n")

	msg, err := parseEmailFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(msg.To) != 2 {
		t.Fatalf("to wrong: %v", msg.To)
	}
}

func TestParseEmailFileJSON(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.json",
		`{"to": ["a@x.com"], "bcc": "d@x.com", "subject": "J", "body": "json body"}`)

	msg, err := parseEmailFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(msg.To) != 1 || len(msg.Bcc) != 1 || msg.Subject != "J" || msg.Body != "json body" {
		t.Fatalf("wrong: %+v", msg)
	}
}

func TestParseEmailFileMarkdown(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.md", "# Meeting notes\n\nFirst line.\nSecond line.\n")

	msg, err := parseEmailFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Subject != "Meeting notes" {
		t.Fatalf("subject wrong: %q", msg.Subject)
	}
	if msg.Body != "First line.\nSecond line.\n" {
		t.Fatalf("body wrong: %q", msg.Body)
	}
}

func TestParseEmailFileMarkdownNoHeading(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.md", "no heading here\n")

	if _, err := parseEmailFile(path); err == nil ||
		!strings.Contains(err.Error(), "E_SEND_EMAIL_NO_SUBJECT") {
		t.Fatalf("expected no-subject error, got %v", err)
	}
}

func TestParseEmailFileText(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.txt", "\n\nSubject:  Greetings \n\nBody keeps\n\nblank lines.\n")

	msg, err := parseEmailFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Subject != "Greetings" {
		t.Fatalf("subject wrong: %q", msg.Subject)
	}
	if msg.Body != "\nBody keeps\n\nblank lines.\n" {
		t.Fatalf("body wrong: %q", msg.Body)
	}
}

func TestParseEmailFileErrors(t *testing.T) {
	dir := t.TempDir()

	if _, err := parseEmailFile(filepath.Join(dir, "missing.yaml")); err == nil {
		t.Fatal("expected error for missing file")
	}

	empty := writeTempEmailFile(t, dir, "empty.txt", "   \n  \n")
	if _, err := parseEmailFile(empty); err == nil ||
		!strings.Contains(err.Error(), "E_SEND_EMAIL_EMPTY_FILE") {
		t.Fatalf("expected empty-file error, got %v", err)
	}

	badYAML := writeTempEmailFile(t, dir, "bad.yaml", "to: [unclosed\n")
	if _, err := parseEmailFile(badYAML); err == nil {
		t.Fatal("expected error for bad YAML")
	}

	badJSON := writeTempEmailFile(t, dir, "bad.json", `{"to": `)
	if _, err := parseEmailFile(badJSON); err == nil {
		t.Fatal("expected error for bad JSON")
	}
}

func TestBuildEmailMessageOverrides(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.yaml",
		"to: file@x.com\nsubject: FileSubj\nbody: FileBody\n")

	msg, err := buildEmailMessage(sendEmailOptions{File: path, Subject: "CliSubj", To: "cli@x.com"})
	if err != nil {
		t.Fatal(err)
	}
	if len(msg.To) != 1 || msg.To[0] != "cli@x.com" {
		t.Fatalf("to override wrong: %v", msg.To)
	}
	if msg.Subject != "CliSubj" {
		t.Fatalf("subject override wrong: %q", msg.Subject)
	}
	if msg.Body != "FileBody" {
		t.Fatalf("body should come from file: %q", msg.Body)
	}
}

func TestBuildEmailMessageFileMissingSubject(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "mail.json", `{"to": "a@x.com", "body": "no subject key"}`)

	if _, err := buildEmailMessage(sendEmailOptions{File: path, To: "a@x.com"}); err == nil {
		t.Fatal("expected error when neither file nor flag provides a subject")
	}
}

func TestBuildEmailMessageLinkExpansion(t *testing.T) {
	opts := sendEmailOptions{
		To:            "user@example.com",
		Subject:       "Links",
		Body:          "Activation: https://serviceactivation.in/subscription/new/QuTYNtSw8PTD",
		IsExpandLinks: true,
	}

	msg, err := buildEmailMessage(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !msg.IsExpanded {
		t.Fatal("expected msg.IsExpanded to be true")
	}
	if msg.ExpandedCount != 1 {
		t.Fatalf("expected ExpandedCount 1, got %d", msg.ExpandedCount)
	}
	expected := "Activation: https://serviceactivation.google.com/subscription/new/QuTYNtSw8PTD"
	if msg.Body != expected {
		t.Fatalf("expected expanded body %q, got %q", expected, msg.Body)
	}
}

func TestRunSendEmailWithInstruction_RelaxesToSubject(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "instruction.json", `{
		"links": [
			{
				"url": "https://serviceactivation.in/subscription/new/ALPHA200",
				"recipients": ["job@example.com"],
				"subject": "Subj {{original}}"
			}
		],
		"autoSend": false
	}`)

	args := []string{"--instruction", path, "--json"}
	if err := RunSendEmail(args); err != nil {
		t.Fatalf("expected success with --instruction without --to and --subject, got %v", err)
	}
}

func TestRunSpaceEmailExpand_WithInstruction(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEmailFile(t, dir, "expand_instruction.json", `{
		"links": [
			{"url": "https://serviceactivation.in/subscription/new/BETA300"}
		]
	}`)

	args := []string{"--instruction", path, "--json"}
	if err := runSpaceEmailExpand(args); err != nil {
		t.Fatalf("expected success with runSpaceEmailExpand --instruction, got %v", err)
	}
}

