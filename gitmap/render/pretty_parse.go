package render

import "strings"

// blockKind identifies the type of a parsed block.
type blockKind int

const (
	bkParagraph blockKind = iota
	bkHeading
	bkSubtitle
	bkFence
	bkList
	bkBlank
)

// block is one unit of parsed markdown.
type block struct {
	kind  blockKind
	text  string   // for headings: the raw heading line; paragraphs: joined lines
	lines []string // for fence: body lines (between ``` markers)
	depth int      // heading level (1..6); 0 for non-headings
}

// parse turns raw lines into a slice of blocks, applying the contextual
// rules (subtitle detection, fence collapse) at parse time.
func parse(lines []string) []block {
	var out []block
	i := 0
	for i < len(lines) {
		line := lines[i]
		switch {
		case isFence(line):
			body, next := readFence(lines, i)
			out = appendFence(out, body)
			i = next
		case isHeading(line):
			depth := headingDepth(line)
			out = append(out, block{kind: bkHeading, text: line, depth: depth})
			i++
			// Subtitle peek: tolerate a single blank line between the
			// heading and an italic subtitle (typical markdown style).
			j := i
			if j < len(lines) && strings.TrimSpace(lines[j]) == "" {
				j++
			}
			if j < len(lines) && isItalic(lines[j]) {
				out = append(out, block{kind: bkSubtitle, text: stripItalic(lines[j])})
				i = j + 1
			}

		case strings.TrimSpace(line) == "":
			out = append(out, block{kind: bkBlank})
			i++
		case isListItem(line):
			items, next := readList(lines, i)
			out = append(out, block{kind: bkList, lines: items})
			i = next
		case isIndentedCode(line):
			body, next := readIndentedCode(lines, i)
			out = append(out, block{kind: bkFence, lines: body})
			i = next
		default:
			para, next := readParagraph(lines, i)
			out = append(out, block{kind: bkParagraph, text: para})
			i = next
		}
	}

	return out
}

// appendFence either appends the fence as-is or, when the previous non-blank
// block is a paragraph whose normalized text matches the fence body,
// replaces that paragraph with a collapsed yellow arrow line.
func appendFence(out []block, body []string) []block {
	prevIdx := lastNonBlank(out)
	if prevIdx >= 0 && out[prevIdx].kind == bkParagraph &&
		normalize(out[prevIdx].text) == normalize(strings.Join(body, "\n")) {
		out[prevIdx] = block{
			kind: bkParagraph,
			text: TokYellowOpen + collapseArrow + strings.Join(body, " ") + TokYellowClose,
		}

		return out
	}

	return append(out, block{kind: bkFence, lines: body})
}

func lastNonBlank(bs []block) int {
	for i := len(bs) - 1; i >= 0; i-- {
		if bs[i].kind != bkBlank {
			return i
		}
	}

	return -1
}

// readFence returns the body lines between matching ``` markers and the
// index of the line AFTER the closing fence.
func readFence(lines []string, start int) ([]string, int) {
	var body []string
	i := start + 1
	for i < len(lines) && !isFence(lines[i]) {
		body = append(body, lines[i])
		i++
	}
	if i < len(lines) {
		i++ // consume closing fence
	}

	return body, i
}

// readParagraph greedily joins consecutive non-blank, non-special lines.
func readParagraph(lines []string, start int) (string, int) {
	var buf []string
	i := start
	for i < len(lines) {
		line := lines[i]
		if strings.TrimSpace(line) == "" || isHeading(line) || isFence(line) ||
			isListItem(line) || isIndentedCode(line) {
			break
		}
		buf = append(buf, line)
		i++
	}

	return strings.Join(buf, " "), i
}

func isFence(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "```")
}

func isHeading(line string) bool {
	t := strings.TrimLeft(line, " ")

	return strings.HasPrefix(t, "# ") || strings.HasPrefix(t, "## ") ||
		strings.HasPrefix(t, "### ") || strings.HasPrefix(t, "#### ") ||
		strings.HasPrefix(t, "##### ") || strings.HasPrefix(t, "###### ")
}

func headingDepth(line string) int {
	t := strings.TrimLeft(line, " ")
	depth := 0
	for depth < len(t) && t[depth] == '#' {
		depth++
	}

	return depth
}

func isItalic(line string) bool {
	t := strings.TrimSpace(line)
	if len(t) < 3 {
		return false
	}

	return (strings.HasPrefix(t, "*") && strings.HasSuffix(t, "*") &&
		!strings.HasPrefix(t, "**")) ||
		(strings.HasPrefix(t, "_") && strings.HasSuffix(t, "_") &&
			!strings.HasPrefix(t, "__"))
}

func stripItalic(line string) string {
	t := strings.TrimSpace(line)
	t = strings.TrimPrefix(t, "*")
	t = strings.TrimPrefix(t, "_")
	t = strings.TrimSuffix(t, "*")
	t = strings.TrimSuffix(t, "_")

	return t
}

// normalize lowercases, trims, and collapses whitespace for paragraph↔fence
// equivalence checks (rule 1).
func normalize(s string) string {
	return strings.Join(strings.Fields(strings.ToLower(s)), " ")
}

// isListItem matches bullet (- / *) and table (|) rows so they render
// as discrete padded lines instead of being collapsed into a paragraph.
func isListItem(line string) bool {
	t := strings.TrimLeft(line, " ")
	if strings.HasPrefix(t, "- ") || strings.HasPrefix(t, "* ") {
		return true
	}

	return strings.HasPrefix(t, "|")
}

// readList gathers consecutive list / table rows into a single block.
func readList(lines []string, start int) ([]string, int) {
	var out []string
	i := start
	for i < len(lines) && isListItem(lines[i]) {
		out = append(out, lines[i])
		i++
	}

	return out, i
}

// isIndentedCode matches a markdown indented code block line: a non-blank
// line beginning with at least four spaces. Help texts use this to render
// CLI examples; the parser must capture them as fenced output so each
// source line stays on its own row instead of being joined into a single
// paragraph.
func isIndentedCode(line string) bool {
	if !strings.HasPrefix(line, "    ") {
		return false
	}

	return strings.TrimSpace(line) != ""
}

// readIndentedCode gathers consecutive indented (and interleaved blank)
// lines into a code block. Trailing blanks are dropped. The leading four
// spaces are stripped so emitBlock can apply the standard body indent.
func readIndentedCode(lines []string, start int) ([]string, int) {
	var body []string
	i := start
	for i < len(lines) {
		line := lines[i]
		if isIndentedCode(line) {
			body = append(body, strings.TrimPrefix(line, "    "))
			i++

			continue
		}
		if strings.TrimSpace(line) == "" && i+1 < len(lines) && isIndentedCode(lines[i+1]) {
			body = append(body, "")
			i++

			continue
		}
		break
	}

	return body, i
}
