package cmdai

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func resolveDestinationFilename(name string, dir string) string {
	slug := sanitizeScriptSlug(name)
	hasNumber := startsWithNumber(name)
	if hasNumber {
		return ensurePyExtension(name)
	}

	nextNum := calculateNextScriptNumber(dir)

	return fmt.Sprintf("%02d-%s.py", nextNum, slug)
}

func startsWithNumber(s string) bool {
	trimmed := strings.TrimSpace(s)
	isEmpty := len(trimmed) == 0
	if isEmpty {
		return false
	}

	return unicode.IsDigit(rune(trimmed[0]))
}

func ensurePyExtension(name string) string {
	hasExt := strings.HasSuffix(name, ".py")
	if hasExt {
		return name
	}

	return name + ".py"
}

func calculateNextScriptNumber(dir string) int {
	maxNum := 0
	entries, err := os.ReadDir(dir)
	hasErr := err != nil
	if hasErr {
		return 45
	}

	for _, entry := range entries {
		num := extractLeadingNumber(entry.Name())
		isGreater := num > maxNum
		if isGreater {
			maxNum = num
		}
	}

	return maxNum + 1
}

func extractLeadingNumber(filename string) int {
	var digits []rune
	for _, r := range filename {
		isDigit := unicode.IsDigit(r)
		if !isDigit {
			break
		}
		digits = append(digits, r)
	}

	isEmpty := len(digits) == 0
	if isEmpty {
		return 0
	}

	val, err := strconv.Atoi(string(digits))
	hasErr := err != nil
	if hasErr {
		return 0
	}

	return val
}
