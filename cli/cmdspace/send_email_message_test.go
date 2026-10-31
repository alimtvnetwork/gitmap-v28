package cmdspace

import (
	"testing"
)

// Split from send_email_test.go on 31-Oct-26 so send_email_test.go stays
// under the 300-line file-size gate. Same package, same helpers
// (writeTempEmailFile lives in send_email_test.go).

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
