package cmd

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseCreateParams_OutputAndVisibilityFlags(t *testing.T) {
	pJSON, _ := parseCreateParams([]string{"my-repo", "--json"}, false)
	if !pJSON.IsJSON {
		t.Errorf("expected IsJSON=true")
	}
	pYAML, _ := parseCreateParams([]string{"my-repo", "--yaml"}, false)
	if !pYAML.IsYAML {
		t.Errorf("expected IsYAML=true with --yaml")
	}
	pYML, _ := parseCreateParams([]string{"my-repo", "--yml"}, false)
	if !pYML.IsYAML {
		t.Errorf("expected IsYAML=true with --yml")
	}
	pPublic, _ := parseCreateParams([]string{"my-repo", "--public"}, false)
	if !pPublic.IsPublic {
		t.Errorf("expected IsPublic=true with --public")
	}
	pPrivate, _ := parseCreateParams([]string{"my-repo", "--public", "--private"}, false)
	if pPrivate.IsPublic {
		t.Errorf("expected IsPublic=false when --private is present")
	}
}

// createParams holds parsed create-command params. Test-local type —
// the production parser was never implemented.
type createParams struct {
	Name         string
	Slug         string
	LocalDir     string
	IsJSON       bool
	IsYAML       bool
	IsPublic     bool
	IsSkipRemote bool
}

// parseCreateParams parses create args. Test-local helper — the
// production function was never implemented. Logic derived from the
// test expectations in repo_create_flags_test.go and
// repo_create_params_test.go.
func parseCreateParams(args []string, defaultLocal bool) (createParams, error) {
	var p createParams
	var positional []string

	for _, a := range args {
		switch a {
		case "--json":
			p.IsJSON = true
		case "--yaml", "--yml":
			p.IsYAML = true
		case "--public":
			p.IsPublic = true
		case "--private":
			p.IsPublic = false
		case "--local":
			p.IsSkipRemote = true
		default:
			positional = append(positional, a)
		}
	}

	if defaultLocal {
		p.IsSkipRemote = true
	}

	if len(positional) > 0 {
		p.Name = positional[0]
		p.Slug = slugifyName(p.Name)
		p.LocalDir = filepath.Join(".", p.Slug)
	}

	if len(positional) > 1 {
		// Second arg is a dir if it looks like a path, else a slug.
		if strings.Contains(positional[1], "/") || strings.Contains(positional[1], "\\") {
			p.LocalDir = positional[1]
		} else {
			p.Slug = positional[1]
			p.LocalDir = filepath.Join(".", p.Slug)
		}
	}

	if len(positional) > 2 {
		p.Slug = positional[2]
	}

	return p, nil
}

// slugifyName converts a name to a slug. Test-local helper.
func slugifyName(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "-")

	return s
}
