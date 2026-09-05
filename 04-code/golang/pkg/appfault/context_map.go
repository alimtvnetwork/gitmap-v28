package appfault

import (
	"encoding/json"
	"fmt"
)

// ContextMap provides a rich, typed wrapper around diagnostic metadata.
type ContextMap map[string]any

// NewContextMap creates an empty ContextMap.
func NewContextMap() ContextMap {
	return make(ContextMap)
}

// NewContextMapWithCapacity creates a ContextMap with preallocated capacity.
func NewContextMapWithCapacity(capacity int) ContextMap {
	return make(ContextMap, capacity)
}

// Set inserts or updates a key-value pair.
func (cm ContextMap) Set(key string, val any) ContextMap {
	cm[key] = val

	return cm
}

// Add is an alias for Set.
func (cm ContextMap) Add(key string, val any) ContextMap {
	return cm.Set(key, val)
}

// Get retrieves a value and presence boolean.
func (cm ContextMap) Get(key string) (any, bool) {
	val, exists := cm[key]

	return val, exists
}

// GetString retrieves a value formatted as a string.
func (cm ContextMap) GetString(key string) string {
	val, exists := cm[key]
	if !exists {
		return ""
	}

	return fmt.Sprintf("%v", val)
}

// Has returns true if the key exists.
func (cm ContextMap) Has(key string) bool {
	_, exists := cm[key]

	return exists
}

// Remove deletes a key from the map.
func (cm ContextMap) Remove(key string) ContextMap {
	delete(cm, key)

	return cm
}

// Count returns the number of key-value entries.
func (cm ContextMap) Count() int {
	return len(cm)
}

// ToJson exports the ContextMap as indented JSON bytes.
func (cm ContextMap) ToJson() ([]byte, error) {
	return json.MarshalIndent(cm, "", "  ")
}

// ToJsonString exports the ContextMap as a JSON string.
func (cm ContextMap) ToJsonString() string {
	b, err := cm.ToJson()
	if err != nil {
		return "{}"
	}

	return string(b)
}

// ContextMapFromJson parses ContextMap from JSON bytes.
func ContextMapFromJson(data []byte) (ContextMap, error) {
	var cm ContextMap
	if err := json.Unmarshal(data, &cm); err != nil {
		return nil, err
	}

	return cm, nil
}
