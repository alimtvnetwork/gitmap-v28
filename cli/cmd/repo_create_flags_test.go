package cmd

import "testing"

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
