package cmd

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestResolverByPathAndSlug(t *testing.T) {
	repos := []model.ScanRecord{
		{ID: 1, Slug: "prompt-architect", RepoName: "prompt-architect", AbsolutePath: `D:/work/prompt-architect`},
		{ID: 2, Slug: "errorwrapper", RepoName: "errorwrapper", AbsolutePath: `D:/work/03-aukgo/errorwrapper`},
	}

	// 1. Relative Windows path with dot-slash
	hit := resolveByPath(`.\prompt-architect`, repos)
	if hit == nil || hit.ID != 1 {
		t.Fatalf("resolveByPath .\\prompt-architect failed: got %+v", hit)
	}

	// 2. Trailing slash
	hit = resolveByPath(`D:/work/prompt-architect/`, repos)
	if hit == nil || hit.ID != 1 {
		t.Fatalf("resolveByPath trailing slash failed: got %+v", hit)
	}

	// 3. Basename match
	hit = resolveByPath(`prompt-architect`, repos)
	if hit == nil || hit.ID != 1 {
		t.Fatalf("resolveByPath basename failed: got %+v", hit)
	}

	// 4. Case-insensitive slug match
	hit = resolveBySlug(`PROMPT-ARCHITECT`, repos)
	if hit == nil || hit.ID != 1 {
		t.Fatalf("resolveBySlug case-insensitive failed: got %+v", hit)
	}

	// 5. Glob match
	hits := resolveByGlob(`*architect`, repos)
	if len(hits) != 1 || hits[0].ID != 1 {
		t.Fatalf("resolveByGlob *architect failed: got %+v", hits)
	}

	_ = filepath.Clean("")
}
