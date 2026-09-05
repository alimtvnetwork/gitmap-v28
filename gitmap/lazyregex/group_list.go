package lazyregex

import (
	"encoding/json"
	"sort"
)

// GroupList is a fluent, nil-safe wrapper around a collection of GroupMap instances.
type GroupList struct {
	items []*GroupMap
}

// NewGroupList creates a new GroupList initialized with the optional groups.
func NewGroupList(groups ...*GroupMap) *GroupList {
	gl := &GroupList{
		items: make([]*GroupMap, 0, len(groups)),
	}

	for _, g := range groups {
		if g != nil {
			gl.items = append(gl.items, g)
		}
	}

	return gl
}

// NewGroupListFrom creates a GroupList from a slice of primitive maps.
func NewGroupListFrom(maps []map[string]string) *GroupList {
	gl := &GroupList{
		items: make([]*GroupMap, 0, len(maps)),
	}

	for _, m := range maps {
		gl.items = append(gl.items, NewGroupMapFrom(m))
	}

	return gl
}

// Items returns the slice of GroupMap pointers.
func (it *GroupList) Items() []*GroupMap {
	if it == nil || it.items == nil {
		return []*GroupMap{}
	}

	return it.items
}

// Len returns the number of group maps in the list.
func (it *GroupList) Len() int {
	if it == nil || it.items == nil {
		return 0
	}

	return len(it.items)
}

// Count is an alias for Len.
func (it *GroupList) Count() int {
	return it.Len()
}

// Length is an alias for Len.
func (it *GroupList) Length() int {
	return it.Len()
}

// IsEmpty reports whether the list is nil or has zero items.
func (it *GroupList) IsEmpty() bool {
	return it.Len() == 0
}

// HasItems reports whether the list contains at least one group map.
func (it *GroupList) HasItems() bool {
	return it.Len() > 0
}

// First returns the first GroupMap in the list, or an empty GroupMap if empty or nil.
func (it *GroupList) First() *GroupMap {
	if it.IsEmpty() {
		return NewGroupMap()
	}

	return it.items[0]
}

// Last returns the last GroupMap in the list, or an empty GroupMap if empty or nil.
func (it *GroupList) Last() *GroupMap {
	if it.IsEmpty() {
		return NewGroupMap()
	}

	return it.items[len(it.items)-1]
}

// At returns the GroupMap at index, or an empty GroupMap if out of bounds or nil.
func (it *GroupList) At(index int) *GroupMap {
	if it == nil || it.items == nil {
		return NewGroupMap()
	}

	if index < 0 || index >= len(it.items) {
		return NewGroupMap()
	}

	return it.items[index]
}

// Add appends a GroupMap to the list and returns the receiver for chaining.
func (it *GroupList) Add(group *GroupMap) *GroupList {
	if it == nil {
		return it
	}

	if group == nil {
		return it
	}

	it.items = append(it.items, group)
	return it
}

// AllKeys returns a sorted, deduplicated slice of all keys across all GroupMaps in the list.
func (it *GroupList) AllKeys() []string {
	if it == nil || it.items == nil {
		return []string{}
	}

	keySet := make(map[string]struct{})
	for _, gm := range it.items {
		if gm == nil {
			continue
		}
		for _, k := range gm.Keys() {
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
func (it *GroupList) Keys() []string {
	return it.AllKeys()
}

// ValuesOf returns a slice of all captured values for a given key across all groups.
func (it *GroupList) ValuesOf(key string) []string {
	if it == nil || it.items == nil {
		return []string{}
	}

	values := make([]string, 0, len(it.items))
	for _, gm := range it.items {
		if gm != nil && gm.Has(key) {
			values = append(values, gm.Get(key))
		}
	}

	return values
}

// Find returns the first GroupMap matching predicate, or an empty GroupMap if not found.
func (it *GroupList) Find(predicate func(g *GroupMap) bool) *GroupMap {
	if it == nil || it.items == nil || predicate == nil {
		return NewGroupMap()
	}

	for _, g := range it.items {
		if g != nil && predicate(g) {
			return g
		}
	}

	return NewGroupMap()
}

// Filter returns a new GroupList containing only groups matching predicate.
func (it *GroupList) Filter(predicate func(g *GroupMap) bool) *GroupList {
	filtered := NewGroupList()
	if it == nil || it.items == nil || predicate == nil {
		return filtered
	}

	for _, g := range it.items {
		if g != nil && predicate(g) {
			filtered.Add(g)
		}
	}

	return filtered
}

// ForEach invokes fn for each GroupMap along with its index.
func (it *GroupList) ForEach(fn func(index int, g *GroupMap)) {
	if it == nil || it.items == nil || fn == nil {
		return
	}

	for i, g := range it.items {
		fn(i, g)
	}
}

// Clone returns a deep copy of the GroupList and its GroupMaps.
func (it *GroupList) Clone() *GroupList {
	if it == nil {
		return NewGroupList()
	}

	cloned := NewGroupList()
	for _, g := range it.items {
		if g != nil {
			cloned.Add(g.Clone())
		}
	}

	return cloned
}

// ToMaps returns the collection as a slice of raw maps.
func (it *GroupList) ToMaps() []map[string]string {
	if it == nil || it.items == nil {
		return []map[string]string{}
	}

	raw := make([]map[string]string, 0, len(it.items))
	for _, g := range it.items {
		if g != nil {
			raw = append(raw, g.ToMap())
		}
	}

	return raw
}

// Raw is an alias for ToMaps.
func (it *GroupList) Raw() []map[string]string {
	return it.ToMaps()
}

// String returns a JSON string representation of the GroupList.
func (it *GroupList) String() string {
	b, err := it.MarshalJSON()
	if err != nil {
		return "[]"
	}

	return string(b)
}

// MarshalJSON serializes the GroupList into JSON bytes.
func (it *GroupList) MarshalJSON() ([]byte, error) {
	if it == nil || it.items == nil {
		return json.Marshal([]map[string]string{})
	}

	return json.Marshal(it.ToMaps())
}
