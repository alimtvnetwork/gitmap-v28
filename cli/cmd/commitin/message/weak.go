package message

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)

const weakPunctuation = ".,:;!?"

func firstWordLower(msg string) string {
	title := msg
	if i := strings.IndexByte(msg, '\n'); i >= 0 {
		title = msg[:i]
	}

	title = strings.TrimSpace(title)
	if title == "" {
		return ""
	}

	word := title
	if i := strings.IndexAny(title, " \t"); i >= 0 {
		word = title[:i]
	}

	word = strings.TrimRight(word, weakPunctuation)

	return strings.ToLower(word)
}

func matchesWeak(msg string, weakWords []string) bool {
	if len(weakWords) == 0 {
		return false
	}

	w := firstWordLower(msg)
	if w == "" {
		return false
	}

	for _, ww := range weakWords {
		if strutil.EqualFoldAnyTrim(ww, w) {
			return true
		}
	}

	return false
}
