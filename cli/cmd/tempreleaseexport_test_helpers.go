package cmd

import (
	"fmt"
	"math"
	"strings"
)

// FormatSeq is an exported wrapper for testing formatSeq.
func FormatSeq(seq, digits int) string {
	return formatSeq(seq, digits)
}

// ParseVersionPatternSafe parses a version pattern without os.Exit.
func ParseVersionPatternSafe(pattern string) (string, int) {
	idx := strings.Index(pattern, "$")
	if idx < 0 {
		return pattern, 0
	}

	prefix := pattern[:idx]
	dollarCount := 0

	for i := idx; i < len(pattern) && pattern[i] == '$'; i++ {
		dollarCount++
	}

	return prefix, dollarCount
}

// ResolveTRBranchExported is an exported wrapper for testing resolveTRBranch.
func ResolveTRBranchExported(version string) string {
	return resolveTRBranch(version)
}

// CheckSequenceRange validates sequence range without os.Exit.
func CheckSequenceRange(start, count, digits int) error {
	maxVal := int(math.Pow(10, float64(digits))) - 1
	endSeq := start + count - 1

	if endSeq > maxVal {
		return fmt.Errorf("sequence %d exceeds %d-digit format (max %d)", endSeq, digits, maxVal)
	}

	return nil
}
