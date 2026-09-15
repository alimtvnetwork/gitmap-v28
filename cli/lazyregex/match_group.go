package lazyregex

// Items returns all captured submatches (index 0 is full match, 1+ are groups).
func (it *MatchGroup) Items() []string {
	if it == nil || len(it.Submatches) == 0 {
		return []string{}
	}

	return it.Submatches
}

// Map returns the named capture groups map.
func (it *MatchGroup) Map() GroupMap {
	if it == nil || it.NamedGroups == nil {
		return NewGroupMap()
	}

	return it.NamedGroups
}

// First returns the first submatch (the full match) or empty string.
func (it *MatchGroup) First() string {
	if it == nil || len(it.Submatches) == 0 {
		return ""
	}

	return it.Submatches[0]
}

// Last returns the last captured submatch or empty string.
func (it *MatchGroup) Last() string {
	if it == nil || len(it.Submatches) == 0 {
		return ""
	}

	return it.Submatches[len(it.Submatches)-1]
}

// FirstOrDefault returns the first submatch or defaultValue if empty.
func (it *MatchGroup) FirstOrDefault(defaultValue string) string {
	first := it.First()
	if first == "" {
		return defaultValue
	}

	return first
}

// At returns the submatch at the specified index or empty string if out of bounds.
func (it *MatchGroup) At(index int) string {
	if it == nil || index < 0 || index >= len(it.Submatches) {
		return ""
	}

	return it.Submatches[index]
}

// Count returns the number of captured submatches.
func (it *MatchGroup) Count() int {
	if it == nil {
		return 0
	}

	return len(it.Submatches)
}

// Len is an alias for Count.
func (it *MatchGroup) Len() int {
	return it.Count()
}

// HasNamed reports whether a named capture group exists.
func (it *MatchGroup) HasNamed(name string) bool {
	if it == nil || it.NamedGroups == nil {
		return false
	}

	return it.NamedGroups.Has(name)
}

// GetNamed retrieves the value of a named capture group or empty string.
func (it *MatchGroup) GetNamed(name string) string {
	if it == nil || it.NamedGroups == nil {
		return ""
	}

	return it.NamedGroups.Get(name)
}

// GetNamedOrDefault retrieves a named group value or defaultValue if absent.
func (it *MatchGroup) GetNamedOrDefault(name, defaultValue string) string {
	if it == nil || it.NamedGroups == nil {
		return defaultValue
	}

	return it.NamedGroups.GetOrDefault(name, defaultValue)
}
