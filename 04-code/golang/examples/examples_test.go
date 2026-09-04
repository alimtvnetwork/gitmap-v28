package examples_test

import (
	"bytes"
	"context"
	"testing"

	"coding-guidelines/common/examples"
	"coding-guidelines/common/pkg/errtype"
	"coding-guidelines/common/pkg/logger"
)

func TestDatabaseQuerySuccess(t *testing.T) {
	repo := examples.NewPluginRepository(nil)
	res := repo.FindById(context.Background(), 1)
	if !res.IsSuccess() || res.IsFailed() {
		t.Fatal("expected successful db query")
	}

	if res.Value.Slug != "seo-optimizer" {
		t.Fatalf("expected slug 'seo-optimizer', got %s", res.Value.Slug)
	}
}

func TestDatabaseQueryNotFound(t *testing.T) {
	repo := examples.NewPluginRepository(nil)
	res := repo.FindById(context.Background(), 404)
	if !res.IsFailed() || !res.HasValidError() {
		t.Fatal("expected error for 404 id")
	}

	if !res.Fault().Is(errtype.NotFound) {
		t.Fatalf("expected errtype.NotFound, got %v", res.Fault().GetType())
	}
}

func newTestWorkflowService(buf *bytes.Buffer) *examples.PluginWorkflowService {
	log := logger.New(logger.DefaultOptions().WithOutput(buf))
	repo := examples.NewPluginRepository(nil)
	client := examples.NewWordPressClient("https://site.example.com")

	return examples.NewPluginWorkflowService(repo, client, log)
}

func TestWorkflowServiceSuccess(t *testing.T) {
	svc := newTestWorkflowService(&bytes.Buffer{})
	res := svc.ActivateWorkflow(context.Background(), 10, 1)
	if !res.IsSuccess() {
		t.Fatal("expected workflow to succeed")
	}

	if res.Value.PluginSummary.Slug != "seo-optimizer" {
		t.Fatalf("expected plugin slug 'seo-optimizer', got %s", res.Value.PluginSummary.Slug)
	}
}

func TestWorkflowServicePropagatesErrorWithoutRewrapping(t *testing.T) {
	svc := newTestWorkflowService(&bytes.Buffer{})
	res := svc.ActivateWorkflow(context.Background(), 10, 404)
	if !res.IsFailed() {
		t.Fatal("expected workflow to fail when plugin is not found in db")
	}

	if !res.Fault().Is(errtype.NotFound) {
		t.Fatalf("expected original errtype.NotFound, got %v", res.Fault().GetType())
	}
}
