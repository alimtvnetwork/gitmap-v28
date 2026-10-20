package cmdasset

import (
	"strings"
	"testing"
)

func TestExtractSlugFromURL(t *testing.T) {
	cases := []struct {
		url      string
		expected string
	}{
		{"https://prnt.sc/FnF7IMTncnd5", "FnF7IMTncnd5"},
		{"https://prnt.sc/oxgHWyUQkiWW/", "oxgHWyUQkiWW"},
		{"https://image.prntscr.com/2/abc.png?param=1", "abc.png"},
	}

	for _, c := range cases {
		got := extractSlugFromURL(c.url)
		if got != c.expected {
			t.Errorf("extractSlugFromURL(%s) = %s, expected %s", c.url, got, c.expected)
		}
	}
}

func TestParseImageURLFromHTML(t *testing.T) {
	htmlImg := `<div class="image-container"><img class="screenshot__image" src="https://image.prntscr.com/image/test12345.png" /></div>`
	got, err := parseImageURLFromHTML("https://prnt.sc/test1", htmlImg)
	if err != nil || got != "https://image.prntscr.com/image/test12345.png" {
		t.Errorf("expected direct image URL, got %s, err %v", got, err)
	}

	htmlOg := `<html><head><meta property="og:image" content="https://api.prntscr.com/image/og_test.png" /></head></html>`
	gotOg, errOg := parseImageURLFromHTML("https://prnt.sc/test2", htmlOg)
	if errOg != nil || gotOg != "https://api.prntscr.com/image/og_test.png" {
		t.Errorf("expected og:image URL, got %s, err %v", gotOg, errOg)
	}

	htmlImgur := `<html><body>Some text with https://i.imgur.com/exampleImgur.png inside</body></html>`
	gotImgur, errImgur := parseImageURLFromHTML("https://prnt.sc/test3", htmlImgur)
	if errImgur != nil || !strings.Contains(gotImgur, "i.imgur.com") {
		t.Errorf("expected imgur URL, got %s, err %v", gotImgur, errImgur)
	}
}

func TestParsePrntArgs(t *testing.T) {
	args := []string{"https://prnt.sc/1", "https://prnt.sc/2", "-e", "--json"}
	opts := parsePrntArgs(args)

	if len(opts.URLs) != 2 || opts.URLs[0] != "https://prnt.sc/1" || opts.URLs[1] != "https://prnt.sc/2" {
		t.Errorf("unexpected URLs: %v", opts.URLs)
	}
	if !opts.ExtractOnly {
		t.Errorf("expected ExtractOnly = true")
	}
	if !opts.JsonFormat {
		t.Errorf("expected JsonFormat = true")
	}

	singleWithCustom := []string{"https://prnt.sc/1", "my_custom_screenshot.png"}
	optsSingle := parsePrntArgs(singleWithCustom)
	if len(optsSingle.URLs) != 1 || optsSingle.CustomDest != "my_custom_screenshot.png" {
		t.Errorf("expected 1 URL with custom dest, got %v / %s", optsSingle.URLs, optsSingle.CustomDest)
	}
}
