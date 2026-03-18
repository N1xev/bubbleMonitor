package process

import (
	"sync"

	"github.com/N1xev/bubbleMonitor/internal/data"
)

var procSlicePool = sync.Pool{
	New: func() interface{} {
		s := make([]data.ProcessInfo, 0, ProcessListCapacity)
		return &s
	},
}

func GetProcSlice() *[]data.ProcessInfo {
	slice, ok := procSlicePool.Get().(*[]data.ProcessInfo)
	if !ok {
		s := make([]data.ProcessInfo, 0, ProcessListCapacity)
		return &s
	}
	return slice
}

func PutProcSlice(s *[]data.ProcessInfo) {
	if s == nil || *s == nil {
		return
	}
	*s = (*s)[:0]
	procSlicePool.Put(s)
}

type Interner struct {
	cache      map[string]*internerEntry
	mu         sync.RWMutex
	count      uint64
	useCounter uint64 // Incremented on each access for LRU tracking
}

const maxInternerSize = 5000
const internerEvictBatch = 500 // Evict 500 oldest entries at a time

type internerEntry struct {
	value string
	used  uint64 // Last use counter for LRU
}

var globalInterner = &Interner{
	cache: make(map[string]*internerEntry),
}

func Intern(s string) string {
	return globalInterner.Intern(s)
}

func (i *Interner) Intern(s string) string {
	i.mu.Lock()
	defer i.mu.Unlock()

	if entry, ok := i.cache[s]; ok {
		i.useCounter++
		entry.used = i.useCounter
		return entry.value
	}

	i.useCounter++
	i.cache[s] = &internerEntry{value: s, used: i.useCounter}
	i.count++

	if len(i.cache) > maxInternerSize {
		i.evictOldest(internerEvictBatch)
	}

	return s
}

// evictOldest removes the n oldest entries from the cache
func (i *Interner) evictOldest(n int) {
	if len(i.cache) <= n {
		return
	}

	// Find n entries with lowest use counter
	type item struct {
		key  string
		used uint64
	}
	items := make([]item, 0, len(i.cache))
	for k, v := range i.cache {
		items = append(items, item{key: k, used: v.used})
	}

	// Simple selection of n oldest (could use heap for large n, but n is small)
	if len(items) > n {
		// Partial sort to find n oldest
		for j := 0; j < n; j++ {
			minIdx := j
			for k := j + 1; k < len(items); k++ {
				if items[k].used < items[minIdx].used {
					minIdx = k
				}
			}
			items[j], items[minIdx] = items[minIdx], items[j]
		}

		// Delete the n oldest
		for j := 0; j < n; j++ {
			delete(i.cache, items[j].key)
		}
	}
}
