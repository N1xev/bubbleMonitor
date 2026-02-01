package process

import (
	"sync"

	"github.com/N1xev/bubbleMonitor/internal/data"
)

// --- Slice Pool ---

var procSlicePool = sync.Pool{
	New: func() interface{} {
		// Start with a reasonable capacity, e.g., 500
		s := make([]data.ProcessInfo, 0, 500)
		return &s
	},
}

// GetProcSlice returns a pointer to a slice from the pool.
// The slice is reset to length 0.
func GetProcSlice() *[]data.ProcessInfo {
	slice, ok := procSlicePool.Get().(*[]data.ProcessInfo)
	if !ok {
		s := make([]data.ProcessInfo, 0, 500)
		return &s
	}
	return slice
}

// PutProcSlice returns a slice to the pool.
// It resets the length to 0 to be ready for reuse.
func PutProcSlice(s *[]data.ProcessInfo) {
	if s == nil || *s == nil {
		return
	}
	// Reset length, keep capacity
	*s = (*s)[:0]
	procSlicePool.Put(s)
}

// --- String Interner ---

type Interner struct {
	mu    sync.RWMutex
	cache map[string]string
	count uint64
}

const maxInternerSize = 5000

var globalInterner = &Interner{
	cache: make(map[string]string),
}

// Intern returns a deduplicated string.
// If s is already in the cache, the cached version is returned.
// Otherwise, s is added to the cache and returned.
func Intern(s string) string {
	return globalInterner.Intern(s)
}

func (i *Interner) Intern(s string) string {
	i.mu.RLock()
	v, ok := i.cache[s]
	i.mu.RUnlock()
	if ok {
		return v
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	if v, ok := i.cache[s]; ok {
		return v
	}

	i.cache[s] = s
	i.count++

	if i.count%1000 == 0 && len(i.cache) > maxInternerSize {
		i.cache = make(map[string]string)
		i.count = 0
	}

	return s
}
