// Package helptextemitter renders HelpDisplay structs to markdown help
// topics. It is the Go half of the DRY help generator (spec 243.3):
// 03-ai-scripts/51-helptext-generator.py drives it via `go run` on an
// ephemeral runner, so no stray package mains live in the repository.
package helpdoc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdspace"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdisplay"
)

// topic binds a helptext/*.md topic name to its HelpDisplay builder and
// the builder's Go symbol (cited in the generated banner).
type topic struct {
	name    string
	builder string
	build   func() helpdisplay.Displayer
}

// topics is the registry of generated helptext/*.md topics. Add a row
// here when a command migrates to the displayer; the generator picks
// it up with no other changes.
var topics = []topic{
	{name: "space", builder: "cmdspace.SpaceHelpDisplay", build: cmdspace.SpaceHelpDisplay},
	{name: "scan", builder: "cmdscan.ScanHelpDisplay", build: cmdscan.ScanHelpDisplay},
	{name: "pull-all", builder: "cmdpull.PullAllHelpDisplay", build: cmdpull.PullAllHelpDisplay},
}

var (
	errMissingOutDir = errors.New("helptextemitter: --out directory is required")
	errUnknownTopic  = errors.New("helptextemitter: unknown topic")
)

// ansiSequence matches SGR color sequences so generated markdown stays plain.
var ansiSequence = regexp.MustCompile("\x1b\\[[0-9;]*m")

// TopicNames returns the sorted names of all registered topics.
func TopicNames() []string {
	names := make([]string, 0, len(topics))
	for _, t := range topics {
		names = append(names, t.name)
	}
	sort.Strings(names)

	return names
}

// findTopic returns the topic for name, or a typed unknown-topic error.
func findTopic(name string) (topic, error) {
	for _, t := range topics {
		if t.name == name {
			return t, nil
		}
	}

	return topic{}, fmt.Errorf("%w: %q", errUnknownTopic, name)
}

// RenderTopic renders one topic's HelpDisplay to its markdown document.
// ANSI color sequences are stripped so the committed file stays plain.
func RenderTopic(name string) (string, error) {
	t, err := findTopic(name)
	if err != nil {
		return "", err
	}
	plain := ansiSequence.ReplaceAllString(t.build().Render(helpdisplay.NewRenderContext(nil)), "")

	return markdownDocument(t, plain), nil
}

// markdownDocument wraps the plain render in the generated-topic banner.
func markdownDocument(t topic, plain string) string {
	var out strings.Builder
	out.WriteString("# gitmap " + t.name + "\n\n")
	out.WriteString("> **Generated — do not hand-edit.** This topic is rendered from the\n")
	out.WriteString("> `" + t.builder + "` HelpDisplay struct. Regenerate with:\n>\n")
	out.WriteString("> `gitmap py 03-ai-scripts/51-helptext-generator.py --topics " + t.name + "`\n\n")
	out.WriteString("```text\n")
	out.WriteString(strings.TrimRight(plain, "\n") + "\n")
	out.WriteString("```\n")

	return out.String()
}

// Generate renders topics (comma-separated; "" means all) into outDir as
// <topic>.md files. It returns typed errors; nothing is swallowed.
func Generate(outDir, topicsCSV string) error {
	if len(outDir) == 0 {
		return errMissingOutDir
	}
	names, err := resolveNames(topicsCSV)
	if err != nil {
		return err
	}
	for _, name := range names {
		if err := writeTopic(outDir, name); err != nil {
			return err
		}
	}

	return nil
}

// resolveNames expands "" to all topics, else validates the CSV list.
func resolveNames(topicsCSV string) ([]string, error) {
	if len(topicsCSV) == 0 {
		return TopicNames(), nil
	}
	names := make([]string, 0)
	for _, raw := range strings.Split(topicsCSV, ",") {
		name := strings.TrimSpace(raw)
		if _, err := findTopic(name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, nil
}

// writeTopic renders one topic and writes <outDir>/<name>.md.
func writeTopic(outDir, name string) error {
	doc, err := RenderTopic(name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("helptextemitter: create out dir %q: %w", outDir, err)
	}
	path := filepath.Join(outDir, name+".md")
	if err := os.WriteFile(path, []byte(doc), 0o644); err != nil {
		return fmt.Errorf("helptextemitter: write %q: %w", path, err)
	}

	return nil
}
