// Package dashboard collects Git repository data for the HTML dashboard.
package dashboard

import (
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// authorAcc accumulates per-author commit data during aggregation.
type authorAcc struct {
	name   string
	email  string
	count  int
	first  string
	last   string
	daySet map[string]bool
}

// buildAuthors aggregates commits into per-author statistics.
func buildAuthors(commits []model.CommitInfo) []model.AuthorInfo {
	index := make(map[string]*authorAcc)

	for _, c := range commits {
		acc, exists := index[c.Email]
		if exists {
			acc.count++
			acc.daySet[c.Date[:10]] = true
			updateDateRange(acc, c.Date)

			continue
		}

		index[c.Email] = &authorAcc{
			name:   c.Author,
			email:  c.Email,
			count:  1,
			first:  c.Date,
			last:   c.Date,
			daySet: map[string]bool{c.Date[:10]: true},
		}
	}

	return collectAuthors(index)
}

// collectAuthors converts the accumulator map to a sorted slice.
func collectAuthors(index map[string]*authorAcc) []model.AuthorInfo {
	authors := make([]model.AuthorInfo, 0, len(index))

	for _, acc := range index {
		authors = append(authors, model.AuthorInfo{
			Name:         acc.name,
			Email:        acc.email,
			TotalCommits: acc.count,
			FirstCommit:  acc.first,
			LastCommit:   acc.last,
			ActiveDays:   len(acc.daySet),
		})
	}

	sort.Slice(authors, func(i, j int) bool {
		return authors[i].TotalCommits > authors[j].TotalCommits
	})

	return authors
}

// updateDateRange expands the first/last bounds of an author accumulator.
func updateDateRange(acc *authorAcc, date string) {
	if date < acc.first {
		acc.first = date
	}

	if date > acc.last {
		acc.last = date
	}
}

// buildFrequency buckets commit dates into daily, weekly, and monthly counts.
func buildFrequency(commits []model.CommitInfo) model.FrequencyData {
	daily := make(map[string]int)
	weekly := make(map[string]int)
	monthly := make(map[string]int)

	for _, c := range commits {
		day := c.Date[:10]
		daily[day]++
		weekly[day[:7]+weekSuffix(day)]++
		monthly[day[:7]]++
	}

	return model.FrequencyData{
		Daily:   daily,
		Weekly:  weekly,
		Monthly: monthly,
	}
}

// weekSuffix returns a "-WNN" suffix based on the day of month.
func weekSuffix(day string) string {
	d := day[8:10]
	if d <= "07" {
		return "-W1"
	}
	if d <= "14" {
		return "-W2"
	}
	if d <= "21" {
		return "-W3"
	}

	return "-W4"
}

// attachTagsToCommits maps tag SHAs to commit entries.
func attachTagsToCommits(commits []model.CommitInfo, tags []model.TagInfo) []model.CommitInfo {
	tagMap := make(map[string][]string, len(tags))
	for _, t := range tags {
		tagMap[t.SHA] = append(tagMap[t.SHA], t.Name)
	}

	for i := range commits {
		short := commits[i].ShortSHA
		full := commits[i].SHA
		matched := tagMap[full]
		if len(matched) == 0 {
			matched = tagMap[short]
		}
		if len(matched) > 0 {
			commits[i].Tags = matched
		}
	}

	return commits
}

// isMergeCommit checks whether the parent string contains multiple parents.
func isMergeCommit(parents string) bool {
	return strings.Contains(parents, " ")
}
