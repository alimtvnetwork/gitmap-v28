package regexnew

import (
	"encoding/json"
	"sort"
)

// GroupMap is a map of named capture groups with fluent helper methods.
type GroupMap map[string]string

// NewGroupMap creates an empty GroupMap instance.
func NewGroupMap() GroupMap {
	return make(GroupMap)
}

// NewGroupMapFrom creates a GroupMap initialized with the given map.
func NewGroupMapFrom(m map[string]string) GroupMap {
	gm := make(GroupMap, len(m))
	if m == nil {
		return gm
	}

	for k, v := range m {
		gm[k] = v
	}

	return gm
}

// Has reports whether the specified key exists in the group map.
func (it GroupMap) Has(key string) bool {
	if it == nil {
		return false
	}

	_, exists := it[key]
	return exists
}

// HasKey is an alias for Has.
func (it GroupMap) HasKey(key string) bool {
	return it.Has(key)
}

// ContainsKey is an alias for Has.
func (it GroupMap) ContainsKey(key string) bool {
	return it.Has(key)
}

// Get retrieves the value associated with key, or empty string if not found or nil.
func (it GroupMap) Get(key string) string {
	if it == nil {
		return ""
	}

	return it[key]
}

// GetOrDefault returns the value for key, or defaultValue if the key is missing or receiver is nil.
func (it GroupMap) GetOrDefault(key, defaultValue string) string {
	if !it.Has(key) {
		return defaultValue
	}

	return it[key]
}

// Set associates key with value and returns the receiver for chaining.
func (it GroupMap) Set(key, value string) GroupMap {
	if it == nil {
		it = make(GroupMap)
	}

	it[key] = value
	return it
}

// Add is a fluent alias for Set.
func (it GroupMap) Add(key, value string) GroupMap {
	return it.Set(key, value)
}

// Put is a fluent alias for Set.
func (it GroupMap) Put(key, value string) GroupMap {
	return it.Set(key, value)
}

// Remove deletes the given key and returns the receiver for chaining.
func (it GroupMap) Remove(key string) GroupMap {
	if it != nil {
		delete(it, key)
	}
	return it
}

// Delete is an alias for Remove.
func (it GroupMap) Delete(key string) GroupMap {
	return it.Remove(key)
}

// Keys returns a sorted slice of all keys present in the group map.
func (it GroupMap) Keys() []string {
	if it == nil {
		return []string{}
	}

	keys := make([]string, 0, len(it))
	for k := range it {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

// AllKeys is an alias for Keys.
func (it GroupMap) AllKeys() []string {
	return it.Keys()
}

// KeyList is an alias for Keys returning a sorted list of all keys.
func (it GroupMap) KeyList() []string {
	return it.Keys()
}

// Values returns a slice of all values present in the group map corresponding to sorted keys.
func (it GroupMap) Values() []string {
	if it == nil {
		return []string{}
	}

	keys := it.Keys()
	values := make([]string, len(keys))
	for i, k := range keys {
		values[i] = it[k]
	}

	return values
}

// AllValues is an alias for Values.
func (it GroupMap) AllValues() []string {
	return it.Values()
}

// ValueList is an alias for Values.
func (it GroupMap) ValueList() []string {
	return it.Values()
}

// Len returns the count of items in the map.
func (it GroupMap) Len() int {
	return len(it)
}

// Count is an alias for Len.
func (it GroupMap) Count() int {
	return len(it)
}

// Length is an alias for Len.
func (it GroupMap) Length() int {
	return len(it)
}

// Size is an alias for Len.
func (it GroupMap) Size() int {
	return len(it)
}

// IsEmpty reports whether the map is nil or has zero items.
func (it GroupMap) IsEmpty() bool {
	return len(it) == 0
}

// HasItems reports whether the map contains at least one item.
func (it GroupMap) HasItems() bool {
	return len(it) > 0
}

// HasAnyItem is an alias for HasItems.
func (it GroupMap) HasAnyItem() bool {
	return len(it) > 0
}

// Clone returns a deep copy of the GroupMap.
func (it GroupMap) Clone() GroupMap {
	if it == nil {
		return NewGroupMap()
	}

	cloned := make(GroupMap, len(it))
	for k, v := range it {
		cloned[k] = v
	}

	return cloned
}

// Clear removes all keys from the group map.
func (it GroupMap) Clear() GroupMap {
	for k := range it {
		delete(it, k)
	}
	return it
}

// ToMap returns a copy of the underlying map[string]string.
func (it GroupMap) ToMap() map[string]string {
	if it == nil {
		return make(map[string]string)
	}

	m := make(map[string]string, len(it))
	for k, v := range it {
		m[k] = v
	}

	return m
}

// Raw is an alias for ToMap.
func (it GroupMap) Raw() map[string]string {
	return it.ToMap()
}

// Items returns a copy of the underlying map[string]string (alias for ToMap).
func (it GroupMap) Items() map[string]string {
	return it.ToMap()
}

// ForEach invokes fn for each key-value pair in sorted key order.
func (it GroupMap) ForEach(fn func(key, val string)) {
	if it == nil || fn == nil {
		return
	}

	keys := it.Keys()
	for _, k := range keys {
		fn(k, it[k])
	}
}

// Filter returns a new GroupMap containing only pairs that satisfy predicate.
func (it GroupMap) Filter(predicate func(key, val string) bool) GroupMap {
	result := NewGroupMap()
	if it == nil || predicate == nil {
		return result
	}

	for k, v := range it {
		if predicate(k, v) {
			result[k] = v
		}
	}

	return result
}

// String returns a JSON string representation of the GroupMap.
func (it GroupMap) String() string {
	b, err := it.MarshalJSON()
	if err != nil {
		return "{}"
	}

	return string(b)
}

// MarshalJSON serializes the GroupMap to JSON bytes.
func (it GroupMap) MarshalJSON() ([]byte, error) {
	if it == nil {
		return json.Marshal(map[string]string{})
	}

	return json.Marshal(map[string]string(it))
}
