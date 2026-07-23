package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestParseCurrentRepoSlug(t *testing.T) {
	const currentSlugVersion = 25
	cases := []struct {
		in       string
		wantBase string
		wantN    int
		wantErr  bool
	}{
		{fmt.Sprintf("gitmap-v%d", currentSlugVersion), "gitmap", currentSlugVersion, false},
		{"gitmap-v100", "gitmap", 100, false},
		{fmt.Sprintf("tool-v%d", 1), "tool", 1, false},
		{"gitmap", "", 0, true},
		{"gitmap-v", "", 0, true},
		{"-v5", "", 0, true},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			base, n, err := parseCurrentRepoSlug(c.in)
			if (err != nil) != c.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, c.wantErr)
			}
			if c.wantErr {
				return
			}
			if base != c.wantBase || n != c.wantN {
				t.Fatalf("got (%q,%d) want (%q,%d)", base, n, c.wantBase, c.wantN)
			}
		})
	}
}

// rewriteTransport rewrites every outbound request's URL host+scheme to
// point at the test server, preserving path so per-slug routing works.
type rewriteTransport struct{ target *url.URL }

func (r *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = r.target.Scheme
	req.URL.Host = r.target.Host
	return http.DefaultTransport.RoundTrip(req)
}

func newTestClient(server *httptest.Server) *http.Client {
	u, _ := url.Parse(server.URL)
	return &http.Client{Transport: &rewriteTransport{target: u}}
}

// hitSet returns a server that returns 200 for the given slug suffixes
// in the path (e.g. "/alimtvnetwork/gitmap-v28"), 404 otherwise.
func hitSet(hits map[string]bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for slug := range hits {
			if strings.HasSuffix(req.URL.Path, "/"+slug) ||
				strings.Contains(req.URL.Path, "/"+slug+"/") {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestProbeSiblings_MaxHitWins(t *testing.T) {
	current := 23
	winner := fmt.Sprintf("gitmap-v%d", current+5)
	other := fmt.Sprintf("gitmap-v%d", current+2)
	srv := hitSet(map[string]bool{winner: true, other: true})
	defer srv.Close()

	slug, ok := probeSiblings(newTestClient(srv), "gitmap", current, 20)
	if !ok {
		t.Fatal("expected hit")
	}
	if slug != winner {
		t.Fatalf("got %q want %q (max-offset hit must win)", slug, winner)
	}
}

func TestProbeSiblings_NoHits(t *testing.T) {
	srv := hitSet(nil)
	defer srv.Close()
	if slug, ok := probeSiblings(newTestClient(srv), "gitmap", 23, 5); ok {
		t.Fatalf("expected no hit, got %q", slug)
	}
}

func TestResolveLatestRepoSlug_FallbackToRelease(t *testing.T) {
	// Only releases API responds 200, siblings + main HEAD return 404.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.URL.Path, "/repos/") && strings.HasSuffix(req.URL.Path, "/releases/latest") {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	slug, source, err := resolveLatestRepoSlug(newTestClient(srv))
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if source != "current-release" {
		t.Fatalf("source=%q want current-release", source)
	}
	if slug == "" {
		t.Fatal("empty slug")
	}
}

func TestResolveLatestRepoSlug_AllFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	if _, _, err := resolveLatestRepoSlug(newTestClient(srv)); err == nil {
		t.Fatal("expected error when all probes fail")
	}
}
