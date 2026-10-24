package lazyregex

import (
	"regexp"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func (it *LazyRegexp) IsMatch(s string) bool {
	if it == nil {
		return false
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return false
	}

	return re.MatchString(s)
}

// IsFound is a semantic alias for IsMatch.
func (it *LazyRegexp) IsFound(s string) bool {
	return it.IsMatch(s)
}

// MatchString reports whether the string s contains any match of the regular expression.
func (it *LazyRegexp) MatchString(s string) bool {
	return it.IsMatch(s)
}

// MatchResult matches string s and returns a MatchResult with submatches and named groups.
func (it *LazyRegexp) MatchResult(s string) *MatchResult {
	if it == nil {
		return NewMatchFailure("", s, apperror.NewSimple("nil LazyRegexp cannot match", "E9000"))
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return NewMatchFailure(it.expression, s, FormatMatchFailureError(it.expression, s, err))
	}

	submatches := re.FindStringSubmatch(s)
	if len(submatches) == 0 {
		return NewMatchFailure(it.expression, s, FormatMatchFailureError(it.expression, s, nil))
	}

	named := it.extractNamedGroups(re, submatches)

	return NewMatchSuccess(it.expression, s, submatches, named)
}

// MatchAllResults matches all non-overlapping occurrences of the pattern in s.
func (it *LazyRegexp) MatchAllResults(s string) *MatchAllResult {
	if it == nil {
		return NewMatchAllFailure("", s, apperror.NewSimple("nil LazyRegexp cannot match", "E9000"))
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return NewMatchAllFailure(it.expression, s, FormatMatchFailureError(it.expression, s, err))
	}

	allMatches := re.FindAllStringSubmatch(s, -1)
	if len(allMatches) == 0 {
		return NewMatchAllFailure(it.expression, s, FormatMatchFailureError(it.expression, s, nil))
	}

	groups := it.buildAllMatchGroups(re, allMatches, s)

	return NewMatchAllSuccess(it.expression, s, groups)
}

func (it *LazyRegexp) buildAllMatchGroups(re *regexp.Regexp, allMatches [][]string, s string) []*MatchGroup {
	groups := make([]*MatchGroup, 0, len(allMatches))
	for _, match := range allMatches {
		named := it.extractNamedGroups(re, match)
		groups = append(groups, &MatchGroup{
			Pattern:     it.expression,
			Content:     s,
			Submatches:  match,
			NamedGroups: named,
		})
	}

	return groups
}

func (it *LazyRegexp) extractNamedGroups(re *regexp.Regexp, match []string) GroupMap {
	res := NewGroupMap()
	names := re.SubexpNames()
	for i, name := range names {
		if name == "" || i >= len(match) {
			continue
		}

		res[name] = match[i]
	}

	return res
}

// MatchError returns nil if s matches the regular expression, or an informative AppError on mismatch.
func (it *LazyRegexp) MatchError(s string) error {
	res := it.MatchResult(s)
	if res.IsSuccess() {
		return nil
	}

	return res.AppError()
}

// MatchPattern compiles or retrieves the pattern and matches against comparing string.
func (it *LazyRegexp) Count(s string) int {
	if it == nil {
		return 0
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return 0
	}

	matches := re.FindAllString(s, -1)

	return len(matches)
}

// GroupBy extracts named capture groups (?P<name>...) from the first match into a GroupMap.
func (it *LazyRegexp) GroupBy(s string) GroupMap {
	res := NewGroupMap()
	if it == nil {
		return res
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return res
	}

	match := re.FindStringSubmatch(s)
	if len(match) == 0 {
		return res
	}

	names := re.SubexpNames()
	for i, name := range names {
		if name == "" || i >= len(match) {
			continue
		}

		res[name] = match[i]
	}

	return res
}

// FindGroups is an alias for GroupBy.
func (it *LazyRegexp) FindGroups(s string) GroupMap {
	return it.GroupBy(s)
}

// FindAllGroups extracts named capture groups across all non-overlapping matches in s into a GroupList.
func (it *LazyRegexp) FindAllGroups(s string) GroupList {
	results := NewGroupList()
	if it == nil {
		return results
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return results
	}

	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return results
	}

	names := re.SubexpNames()
	for _, match := range matches {
		groupMap := NewGroupMap()
		for i, name := range names {
			if name == "" || i >= len(match) {
				continue
			}

			groupMap[name] = match[i]
		}

		results = append(results, groupMap)
	}

	return results
}

// FindString returns the leftmost match in s of the regular expression.
func (it *LazyRegexp) FindString(s string) string {
	if it == nil {
		return ""
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return ""
	}

	return re.FindString(s)
}

// FindStringSubmatch returns a slice of strings holding the leftmost submatches in s.
func (it *LazyRegexp) FindStringSubmatch(s string) []string {
	if it == nil {
		return nil
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return nil
	}

	return re.FindStringSubmatch(s)
}

// FindAllString returns a slice of all successive matches of the expression.
func (it *LazyRegexp) FindAllString(s string, n int) []string {
	if it == nil {
		return nil
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return nil
	}

	return re.FindAllString(s, n)
}

// ReplaceAllString returns a copy of src with all matches replaced by repl.
func (it *LazyRegexp) ReplaceAllString(src, repl string) string {
	if it == nil {
		return src
	}

	re, err := it.compiledRegex()
	if err != nil || re == nil {
		return src
	}

	return re.ReplaceAllString(src, repl)
}

// Find returns a slice holding the text of the leftmost match in b of the regular expression.
func (it *LazyRegexp) Find(b []byte) []byte {
	if it == nil {
		return nil
	}

	re, err := it.compiledRegex()
	if err != nil {
		return nil
	}
	if re == nil {
		return nil
	}

	return re.Find(b)
}

// FindBytes is an alias for Find.
func (it *LazyRegexp) FindBytes(b []byte) []byte {
	return it.Find(b)
}

// IsMatchBytes reports whether the byte slice b contains any match of the regular expression.
func (it *LazyRegexp) IsMatchBytes(b []byte) bool {
	if it == nil {
		return false
	}

	re, err := it.compiledRegex()
	if err != nil {
		return false
	}
	if re == nil {
		return false
	}

	return re.Match(b)
}

// CacheLen returns the number of uniquely cached regular expression instances in the global registry.
