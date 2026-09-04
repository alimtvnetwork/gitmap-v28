package appfaults

import (
	"sync"

	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/errtype"
)

// MutexCollection provides concurrency-safe operations over a Collection.
type MutexCollection struct {
	sync.RWMutex
	inner *Collection
}

// NewMutexCollection creates a thread-safe error collection.
func NewMutexCollection() *MutexCollection {
	return &MutexCollection{
		inner: New(),
	}
}

// Add safely appends an error to the collection.
func (mc *MutexCollection) Add(err *appfault.AppError) *MutexCollection {
	mc.Lock()
	defer mc.Unlock()
	mc.inner.Add(err)

	return mc
}

// AddType safely adds an error by type.
func (mc *MutexCollection) AddType(errType errtype.Variation) *MutexCollection {
	mc.Lock()
	defer mc.Unlock()
	mc.inner.AddType(errType)

	return mc
}

// AddTypeMsg safely adds an error by type and message.
func (mc *MutexCollection) AddTypeMsg(errType errtype.Variation, msg string) *MutexCollection {
	mc.Lock()
	defer mc.Unlock()
	mc.inner.AddTypeMsg(errType, msg)

	return mc
}

// AddError safely wraps and appends a standard error with explicit type.
func (mc *MutexCollection) AddError(errType errtype.Variation, cause error) *MutexCollection {
	mc.Lock()
	defer mc.Unlock()
	mc.inner.AddError(errType, cause)

	return mc
}

// HasError safely checks if any error is stored.
func (mc *MutexCollection) HasError() bool {
	mc.RLock()
	defer mc.RUnlock()

	return mc.inner.HasError()
}

// IsSuccess safely checks if collection is empty.
func (mc *MutexCollection) IsSuccess() bool {
	return !mc.HasError()
}

// Count safely returns the item count.
func (mc *MutexCollection) Count() int {
	mc.RLock()
	defer mc.RUnlock()

	return mc.inner.Count()
}

// Snapshot safely returns a cloned Collection.
func (mc *MutexCollection) Snapshot() *Collection {
	mc.RLock()
	defer mc.RUnlock()

	return NewFromFaults(mc.inner.items...)
}
