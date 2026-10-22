// Package cmdspace implements the `gitmap space` subcommands that were
// split out of cli/cmd during the file-size enforcement work (spec 243).
package cmdspace

import (
	"errors"
	"strings"
	"testing"
)

func TestBuildTriggerEmailHTMLEscapesAndRows(t *testing.T) {
	info := MachineInfo{Host: "h<ost>", User: "u&ser", IP: "10.0.0.1"}
	outcome := TaskOutcome{
		Repo: "<script>alert(1)</script>", Branch: "main", Operation: "release",
		IsSuccess: true, Summary: "ok & done", StartedAt: "s", FinishedAt: "f",
	}

	html := BuildTriggerEmailHTML(EmailTaskRelease, outcome, info)

	for _, want := range []string{
		"&lt;script&gt;alert(1)&lt;/script&gt;", "h&lt;ost&gt;", "u&amp;ser",
		"10.0.0.1", "main", "release", "ok &amp; done", "s", "f",
		"Host", "User", "Outbound IP", "Repo", "Branch", "Operation", "Outcome", "Started", "Finished",
		bannerOKColor, statusOKText,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("html missing %q", want)
		}
	}

	if strings.Contains(html, "<script>") {
		t.Fatal("html contains unescaped script tag")
	}
}

func TestBuildTriggerEmailHTMLFailureBanner(t *testing.T) {
	info := CollectMachineInfo()
	outcome := InferTaskOutcome(EmailTaskRelease, errors.New("boom"))

	html := BuildTriggerEmailHTML(EmailTaskRelease, outcome, info)

	if !strings.Contains(html, bannerFailColor) || !strings.Contains(html, statusFailText) {
		t.Fatal("expected red FAILED banner")
	}

	if !strings.Contains(html, "boom") {
		t.Fatal("expected outcome summary in html")
	}
}

func TestInferTaskOutcome(t *testing.T) {
	ok := InferTaskOutcome(EmailTaskPull, nil)
	if !ok.IsSuccess || ok.Summary != triggerSuccessSummary || ok.Operation != "pull" {
		t.Fatalf("unexpected success outcome: %+v", ok)
	}

	failed := InferTaskOutcome(EmailTaskClone, errors.New("network down"))
	if failed.IsSuccess || failed.Summary != "network down" {
		t.Fatalf("unexpected failure outcome: %+v", failed)
	}
}

func TestFireEmailTaskTriggers(t *testing.T) {
	path := openTestTriggerDBPath(t)
	stubTriggerStore(t, path)

	saveTestTrigger(t, path, EmailTrigger{Name: "rel-alert", To: "me@example.com", Task: EmailTaskRelease})
	saveTestTrigger(t, path, EmailTrigger{Name: "pull-alert", To: "other@example.com", Task: EmailTaskPull})

	var gotTo, gotSubject, gotHTML string
	var callCount int
	original := sendTriggerEmail
	sendTriggerEmail = func(to, subject, htmlBody string) error {
		callCount++
		gotTo, gotSubject, gotHTML = to, subject, htmlBody

		return nil
	}
	t.Cleanup(func() { sendTriggerEmail = original })

	outcome := TaskOutcome{
		Repo: "demo", Branch: "main", Operation: "release",
		IsSuccess: true, Summary: "v1.2.3 shipped", StartedAt: "s", FinishedAt: "f",
	}
	FireEmailTaskTriggers(EmailTaskRelease, outcome)

	if callCount != 1 {
		t.Fatalf("expected 1 send, got %d", callCount)
	}

	if gotTo != "me@example.com" {
		t.Fatalf("expected to me@example.com, got %q", gotTo)
	}

	if !strings.Contains(gotSubject, "release") || !strings.Contains(gotHTML, "v1.2.3 shipped") {
		t.Fatalf("subject/html missing outcome: %q / %.80s", gotSubject, gotHTML)
	}
}

func TestFireEmailTaskTriggersSendFailureIsSilent(t *testing.T) {
	path := openTestTriggerDBPath(t)
	stubTriggerStore(t, path)

	saveTestTrigger(t, path, EmailTrigger{Name: "x", To: "a@b.c", Task: EmailTaskScan})

	original := sendTriggerEmail
	sendTriggerEmail = func(to, subject, htmlBody string) error { return errors.New("smtp down") }
	t.Cleanup(func() { sendTriggerEmail = original })

	FireEmailTaskTriggers(EmailTaskScan, InferTaskOutcome(EmailTaskScan, nil))
}
