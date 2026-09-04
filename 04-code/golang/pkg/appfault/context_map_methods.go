package appfault

import (
	"fmt"
	"sort"
	"strings"
)

// Clone creates a shallow copy of the ContextMap.
func (cm ContextMap) Clone() ContextMap {
	if cm == nil {
		return make(ContextMap)
	}

	cloned := make(ContextMap, len(cm))
	for k, v := range cm {
		cloned[k] = v
	}

	return cloned
}

// Merge copies all entries from other into cm.
func (cm ContextMap) Merge(other ContextMap) ContextMap {
	for k, v := range other {
		cm[k] = v
	}

	return cm
}

// Keys returns a sorted slice of keys.
func (cm ContextMap) Keys() []string {
	keys := make([]string, 0, len(cm))
	for k := range cm {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

// Format formats the map as a human-readable comma-separated string.
func (cm ContextMap) Format() string {
	if len(cm) == 0 {
		return "{}"
	}

	keys := cm.Keys()
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, cm[k]))
	}

	return "{" + strings.Join(parts, ", ") + "}"
}
