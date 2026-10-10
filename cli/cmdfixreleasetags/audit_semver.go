// Package cmdfixreleasetags provides semantic version comparison and tag sorting.
package cmdfixreleasetags

import (
	"sort"
	"strings"
)

func splitSemverInts(v string) [3]int {
	norm := NormalizeVersionTag(v)
	var nums [3]int
	parts := strings.SplitN(norm, ".", 3)
	for i, p := range parts {
		if i >= 3 {
			break
		}

		n := 0
		for _, r := range p {
			if r < '0' || r > '9' {
				break
			}
			n = n*10 + int(r-'0')
		}
		nums[i] = n
	}

	return nums
}

// CompareSemver compares two semantic versions. Returns 1 if a > b, -1 if a < b, 0 if equal.
func CompareSemver(a, b string) int {
	pa := splitSemverInts(a)
	pb := splitSemverInts(b)

	for i := 0; i < 3; i++ {
		if pa[i] > pb[i] {
			return 1
		}
		if pa[i] < pb[i] {
			return -1
		}
	}

	return 0
}

// SortTagsDescending sorts tag strings in descending semantic version order.
func SortTagsDescending(tags []string) {
	sort.Slice(tags, func(i, j int) bool {
		cmp := CompareSemver(tags[i], tags[j])
		if cmp != 0 {
			return cmp > 0
		}

		return tags[i] > tags[j]
	})
}

// FindLatestHealthyTag locates the highest semver release that is healthy.
func FindLatestHealthyTag(records []ReleaseTagAuditRecord) string {
	var healthyTags []string
	for _, rec := range records {
		if rec.AuditReason == ReasonHealthy {
			healthyTags = append(healthyTags, rec.Tag)
		}
	}

	if len(healthyTags) == 0 {
		return ""
	}

	SortTagsDescending(healthyTags)

	return healthyTags[0]
}
