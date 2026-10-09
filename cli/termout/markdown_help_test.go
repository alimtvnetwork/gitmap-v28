package termout

import (
	"testing"
)

const sampleMarkdown = `# gitmap demo

Demo description text.

## Usage

    gitmap demo [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --fast | false | Enable fast mode |

## Tips

- Tip number one
`

func TestFromMarkdown_EmptyInput(t *testing.T) {
	_, err := FromMarkdown(nil)
	if err == nil {
		t.Fatalf("expected error on empty markdown input")
	}
}

func TestFromMarkdown_TitleAndUsage(t *testing.T) {
	menu, err := FromMarkdown([]byte(sampleMarkdown))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if menu.Title != "gitmap demo" {
		t.Errorf("expected title 'gitmap demo', got %q", menu.Title)
	}
	if len(menu.UsageLines) != 1 {
		t.Errorf("expected 1 usage line, got %d", len(menu.UsageLines))
	}
}

func TestFromMarkdown_FlagsAndTips(t *testing.T) {
	menu, err := FromMarkdown([]byte(sampleMarkdown))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(menu.FooterFlags) != 1 {
		t.Errorf("expected 1 flag, got %d", len(menu.FooterFlags))
	}
	if len(menu.Tips) != 1 {
		t.Errorf("expected 1 tip, got %d", len(menu.Tips))
	}
}

func TestFromMarkdown_SectionsAndBullets(t *testing.T) {
	md := "# gitmap sample\n\n## Subcommands\n\n- add: Add item\n- rm: Remove item\n"
	menu, err := FromMarkdown([]byte(md))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(menu.Sections) < 1 {
		t.Fatalf("expected at least 1 section, got %d", len(menu.Sections))
	}
}
