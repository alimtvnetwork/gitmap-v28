package streamwriter

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

// ReentrantMutex allows the same goroutine to acquire the lock multiple times without deadlocking.
type ReentrantMutex struct {
	mu        sync.Mutex
	owner     atomic.Int64
	recursion int32
}

// Lock acquires the mutex or increments the recursion count if already owned by current goroutine.
func (m *ReentrantMutex) Lock() {
	gid := getGoroutineId()
	if m.owner.Load() == gid {
		m.recursion++

		return
	}

	m.mu.Lock()
	m.owner.Store(gid)
	m.recursion = 1
}

// Unlock decrements the recursion count or releases the mutex when count reaches zero.
func (m *ReentrantMutex) Unlock() {
	gid := getGoroutineId()
	if m.owner.Load() != gid {
		return
	}

	m.recursion--
	if m.recursion <= 0 {
		m.recursion = 0
		m.owner.Store(0)
		m.mu.Unlock()
	}
}

func getGoroutineId() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	b := buf[:n]
	prefix := []byte("goroutine ")
	if !bytes.HasPrefix(b, prefix) {
		return 0
	}

	b = b[len(prefix):]
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		return 0
	}

	id, err := strconv.ParseInt(string(b[:i]), 10, 64)
	if err != nil {
		return 0
	}

	return id
}
