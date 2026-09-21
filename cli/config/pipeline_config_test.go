package config

import (
	"os"
	"testing"
	"time"
)

func TestResolvePipelineMaxDbBytes_Default(t *testing.T) {
	_ = os.Unsetenv("GITMAP_PIPELINE_MAX_DB_MB")
	maxBytes := ResolvePipelineMaxDbBytes()
	if maxBytes != 10*1024*1024 {
		t.Fatalf("expected 10MB default, got %d", maxBytes)
	}
}

func TestResolvePipelineMaxDbBytes_EnvOverride(t *testing.T) {
	_ = os.Setenv("GITMAP_PIPELINE_MAX_DB_MB", "25")
	defer func() { _ = os.Unsetenv("GITMAP_PIPELINE_MAX_DB_MB") }()

	maxBytes := ResolvePipelineMaxDbBytes()
	if maxBytes != 25*1024*1024 {
		t.Fatalf("expected 25MB, got %d", maxBytes)
	}
}

func TestResolvePipelineCacheTTL_Default(t *testing.T) {
	_ = os.Unsetenv("GITMAP_PIPELINE_CACHE_TTL_SEC")
	ttl := ResolvePipelineCacheTTL()
	if ttl != 5*time.Second {
		t.Fatalf("expected 5s default, got %v", ttl)
	}
}

func TestResolvePipelineCacheTTL_EnvOverride(t *testing.T) {
	_ = os.Setenv("GITMAP_PIPELINE_CACHE_TTL_SEC", "12")
	defer func() { _ = os.Unsetenv("GITMAP_PIPELINE_CACHE_TTL_SEC") }()

	ttl := ResolvePipelineCacheTTL()
	if ttl != 12*time.Second {
		t.Fatalf("expected 12s, got %v", ttl)
	}
}
