package lazyregex

import (
	"encoding/json"
	"sort"
)

// GroupMap is a fluent, nil-safe wrapper around a map of named capture groups.
type GroupMap struct {
	items map[string]string
}

// NewGroupMap creates an empty GroupMap instance.
func NewGroupMap() *GroupMap {
	return &GroupMap{
		items: make(map[string]string),
	}
}

// NewGroupMapFrom creates a GroupMap initialized with the given map.
func NewGroupMapFrom(m map[string]string) *GroupMap {
	gm := NewGroupMap()
	if m == nil {
		return gm
	}

	for k, v := range m {
		gm.items[k] = v
	}

	return gm
}

// Has reports whether the specified key exists in the group map.
func (it *GroupMap) Has(key string) bool {
	if it == nil || it.items == nil {
		return false
	}

	_, exists := it.items[key]
	return exists
}

// HasKey is an alias for Has.
func (it *GroupMap) HasKey(key string) bool {
	return it.Has(key)
}

// Get retrieves the value associated with key, or empty string if not found or nil.
func (it *GroupMap) Get(key string) string {
	if it == nil || it.items == nil {
		return ""
	}

	return it.items[key]
}

// GetOrDefault returns the value for key, or defaultValue if the key is missing or receiver is nil.
func (it *GroupMap) GetOrDefault(key, defaultValue string) string {
	if !it.Has(key) {
		return defaultValue
	}

	return it.Get(key)
}

// Set associates key with value and returns the receiver for chaining.
func (it *GroupMap) Set(key, value string) *GroupMap {
	if it == nil {
		return it
	}

	if it.items == nil {
		it.items = make(map[string]string)
	}

	it.items[key] = value
	return it
}

// Add is a fluent alias for Set.
func (it *GroupMap) Add(key, value string) *GroupMap {
	return it.Set(key, value)
}

// Remove deletes the given key and returns the receiver for chaining.
func (it *GroupMap) Remove(key string) *GroupMap {
	if it == nil || it.items == nil {
		return it
	}

	delete(it.items, key)
	return it
}

// Delete is an alias for Remove.
func (it *GroupMap) Delete(key string) *GroupMap {
	return it.Remove(key)
}

// Keys returns a sorted slice of all keys present in the group map.
func (it *GroupMap) Keys() []string {
	if it == nil || it.items == nil {
		return []string{}
	}

	keys := make([]string, 0, len(it.items))
	for k := range it.items {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

// AllKeys is an alias for Keys.
func (it *GroupMap) AllKeys() []string {
	return it.Keys()
}

// Values returns a slice of all values present in the group map corresponding to sorted keys.
func (it *GroupMap) Values() []string {
	if it == nil || it.items == nil {
		return []string{}
	}

	keys := it.Keys()
	values := make([]string, len(keys))
	for i, k := range keys {
		values[i] = it.items[k]
	}

	return values
}

// AllValues is an alias for Values.
func (it *GroupMap) AllValues() []string {
	return it.Values()
}

// Len returns the count of items in the map.
func (it *GroupMap) Len() int {
	if it == nil || it.items == nil {
		return 0
	}

	return len(it.items)
}

// Count is an alias for Len.
func (it *GroupMap) Count() int {
	return it.Len()
}

// Length is an alias for Len.
func (it *GroupMap) Length() int {
	return it.Len()
}

// IsEmpty reports whether the map is nil or has zero items.
func (it *GroupMap) IsEmpty() bool {
	return it.Len() == 0
}

// HasItems reports whether the map contains at least one item.
func (it *GroupMap) HasItems() bool {
	return it.Len() > 0
}

// Clone returns a deep copy of the GroupMap.
func (it *GroupMap) Clone() *GroupMap {
	if it == nil {
		return NewGroupMap()
	}

	return NewGroupMapFrom(it.items)
}

// Clear removes all keys from the group map.
func (it *GroupMap) Clear() *GroupMap {
	if it == nil {
		return it
	}

	it.items = make(map[string]string)
	return it
}

// ToMap returns a copy of the underlying map[string]string.
func (it *GroupMap) ToMap() map[string]string {
	if it == nil || it.items == nil {
		return make(map[string]string)
	}

	m := make(map[string]string, len(it.items))
	for k, v := range it.items {
		m[k] = v
	}

	return m
}

// Raw is an alias for ToMap.
func (it *GroupMap) Raw() map[string]string {
	return it.ToMap()
}

// ForEach invokes fn for each key-value pair in sorted key order.
func (it *GroupMap) ForEach(fn func(key, val string)) {
	if it == nil || it.items == nil || fn == nil {
		return
	}

	keys := it.Keys()
	for _, k := range keys {
		fn(k, it.items[k])
	}
}

// Filter returns a new GroupMap containing only pairs that satisfy predicate.
func (it *GroupMap) Filter(predicate func(key, val string) bool) *GroupMap {
	result := NewGroupMap()
	if it == nil || it.items == nil || predicate == nil {
		return result
	}

	for k, v := range it.items {
		if predicate(k, v) {
			result.Set(k, v)
		}
	}

	return result
}

// String returns a JSON string representation of the GroupMap.
func (it *GroupMap) String() string {
	b, err := it.MarshalJSON()
	if err != nil {
		return "{}"
	}

	return string(b)
}

// MarshalJSON serializes the GroupMap to JSON bytes.
func (it *GroupMap) MarshalJSON() ([]byte, error) {
	if it == nil || it.items == nil {
		return json.Marshal(map[string]string{})
	}

	return json.Marshal(it.items)
}
