package release

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func newTransport() *http.Transport {
	t, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return &http.Transport{}
	}

	return t.Clone()
}

var githubTransport = newTransport()

func buildGitHubReleasesURL(owner, repo string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   fmt.Sprintf("/repos/%s/%s/releases", url.PathEscape(owner), url.PathEscape(repo)),
	}
}

func buildGitHubUploadURL(owner, repo string, releaseID int, filename string) *url.URL {
	return &url.URL{
		Scheme:   "https",
		Host:     "uploads.github.com",
		Path:     fmt.Sprintf("/repos/%s/%s/releases/%d/assets", url.PathEscape(owner), url.PathEscape(repo), releaseID),
		RawQuery: "name=" + url.QueryEscape(filename),
	}
}

func newGitHubRequest(method string, u *url.URL, body io.ReadCloser, contentLength int64) *http.Request {
	return &http.Request{
		Method:        method,
		URL:           u,
		Host:          u.Host,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		Body:          body,
		ContentLength: contentLength,
	}
}

func doGitHubRequest(req *http.Request) (*http.Response, error) {
	return githubTransport.RoundTrip(req)
}

func readDirNames(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.Readdirnames(-1)
}
