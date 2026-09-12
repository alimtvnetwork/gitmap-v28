package lazyregex

import (
	"encoding/json"
	"sort"
)

// GroupList is a collection of GroupMap instances with fluent helper methods.
type GroupList []GroupMap

// NewGroupList creates a new GroupList initialized with the optional groups.
func NewGroupList(groups ...GroupMap) GroupList {
	gl := make(GroupList, 0, len(groups))
	for _, g := range groups {
		if g != nil {
			gl = append(gl, g)
		}
	}

	return gl
}

// NewGroupListFrom creates a GroupList from a slice of primitive maps.
func NewGroupListFrom(maps []map[string]string) GroupList {
	gl := make(GroupList, 0, len(maps))
	for _, m := range maps {
		gl = append(gl, NewGroupMapFrom(m))
	}

	return gl
}

// Items returns the slice of GroupMap instances.
func (it GroupList) Items() []GroupMap {
	if it == nil {
		return []GroupMap{}
	}

	return it
}

// AllItems is an alias for Items.
func (it GroupList) AllItems() []GroupMap {
	return it.Items()
}

// Len returns the number of group maps in the list.
func (it GroupList) Len() int {
	return len(it)
}

// Count is an alias for Len.
func (it GroupList) Count() int {
	return len(it)
}

// Length is an alias for Len.
func (it GroupList) Length() int {
	return len(it)
}

// Size is an alias for Len.
func (it GroupList) Size() int {
	return len(it)
}

// IsEmpty reports whether the list is nil or has zero items.
func (it GroupList) IsEmpty() bool {
	return len(it) == 0
}

// HasItems reports whether the list contains at least one group map.
func (it GroupList) HasItems() bool {
	return len(it) > 0
}

// HasAnyItem is an alias for HasItems.
func (it GroupList) HasAnyItem() bool {
	return len(it) > 0
}

// First returns the first GroupMap in the list, or an empty GroupMap if empty or nil.
func (it GroupList) First() GroupMap {
	if len(it) == 0 {
		return NewGroupMap()
	}

	return it[0]
}

// Last returns the last GroupMap in the list, or an empty GroupMap if empty or nil.
func (it GroupList) Last() GroupMap {
	if len(it) == 0 {
		return NewGroupMap()
	}

	return it[len(it)-1]
}

// At returns the GroupMap at index, or an empty GroupMap if out of bounds or nil.
func (it GroupList) At(index int) GroupMap {
	if index < 0 || index >= len(it) {
		return NewGroupMap()
	}

	return it[index]
}

// Add appends a GroupMap to the list and returns the receiver pointer for chaining.
func (it *GroupList) Add(group GroupMap) *GroupList {
	if it == nil {
		return it
	}

	if group == nil {
		return it
	}

	*it = append(*it, group)

	return it
}

// Append is an alias for Add.
func (it *GroupList) Append(group GroupMap) *GroupList {
	return it.Add(group)
}

// RemoveAt removes the GroupMap at index and returns the receiver pointer for chaining.
func (it *GroupList) RemoveAt(index int) *GroupList {
	if it == nil || index < 0 || index >= len(*it) {
		return it
	}

	*it = append((*it)[:index], (*it)[index+1:]...)

	return it
}

// AllKeys returns a sorted, deduplicated slice of all keys across all GroupMaps in the list.
func (it GroupList) AllKeys() []string {
	if len(it) == 0 {
		return []string{}
	}

	keySet := make(map[string]struct{})
	for _, gm := range it {
		for k := range gm {
			keySet[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

// Keys is an alias for AllKeys.
func (it GroupList) Keys() []string {
	return it.AllKeys()
}

// KeyList is an alias for AllKeys returning a sorted list of all keys.
func (it GroupList) KeyList() []string {
	return it.AllKeys()
}

// ValuesOf returns a slice of all captured values for a given key across all groups.
func (it GroupList) ValuesOf(key string) []string {
	if len(it) == 0 {
		return []string{}
	}

	values := make([]string, 0, len(it))
	for _, gm := range it {
		if gm.Has(key) {
			values = append(values, gm[key])
		}
	}

	return values
}

// Find returns the first GroupMap matching predicate, or an empty GroupMap if not found.
func (it GroupList) Find(predicate func(g GroupMap) bool) GroupMap {
	if predicate == nil {
		return NewGroupMap()
	}

	for _, g := range it {
		if predicate(g) {
			return g
		}
	}

	return NewGroupMap()
}

// Filter returns a new GroupList containing only groups matching predicate.
func (it GroupList) Filter(predicate func(g GroupMap) bool) GroupList {
	filtered := NewGroupList()
	if predicate == nil {
		return filtered
	}

	for _, g := range it {
		if predicate(g) {
			filtered = append(filtered, g)
		}
	}

	return filtered
}

// ForEach invokes fn for each GroupMap along with its index.
func (it GroupList) ForEach(fn func(index int, g GroupMap)) {
	if fn == nil {
		return
	}

	for i, g := range it {
		fn(i, g)
	}
}

// Clone returns a deep copy of the GroupList and its GroupMaps.
func (it GroupList) Clone() GroupList {
	cloned := make(GroupList, 0, len(it))
	for _, g := range it {
		cloned = append(cloned, g.Clone())
	}

	return cloned
}

// ToMaps returns the collection as a slice of raw maps.
func (it GroupList) ToMaps() []map[string]string {
	if len(it) == 0 {
		return []map[string]string{}
	}

	raw := make([]map[string]string, 0, len(it))
	for _, g := range it {
		raw = append(raw, g.ToMap())
	}

	return raw
}

// Raw is an alias for ToMaps.
func (it GroupList) Raw() []map[string]string {
	return it.ToMaps()
}

// String returns a JSON string representation of the GroupList.
func (it GroupList) String() string {
	b, err := it.MarshalJSON()
	if err != nil {
		return "[]"
	}

	return string(b)
}

// MarshalJSON serializes the GroupList into JSON bytes.
func (it GroupList) MarshalJSON() ([]byte, error) {
	if len(it) == 0 {
		return json.Marshal([]map[string]string{})
	}

	return json.Marshal(it.ToMaps())
}
