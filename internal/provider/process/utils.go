package process

import (
	"sync"

	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/provider"
)

var procSlicePool = sync.Pool{
	New: func() interface{} {
		s := make([]data.ProcessInfo, 0, provider.ProcessListCapacity)
		return &s
	},
}

func GetProcSlice() *[]data.ProcessInfo {
	slice, ok := procSlicePool.Get().(*[]data.ProcessInfo)
	if !ok {
		s := make([]data.ProcessInfo, 0, provider.ProcessListCapacity)
		return &s
	}
	return slice
}

func PutProcSlice(s *[]data.ProcessInfo) {
	if s == nil || *s == nil {
		return
	}
	// Reset length, keep capacity
	*s = (*s)[:0]
	procSlicePool.Put(s)
}

type Interner struct {
	mu    sync.RWMutex
	cache map[string]string
	count uint64
}

const maxInternerSize = 5000

var globalInterner = &Interner{
	cache: make(map[string]string),
}

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

	if i.count%provider.InternerCleanupFrequency == 0 && len(i.cache) > maxInternerSize {
		i.cache = make(map[string]string)
		i.count = 0
	}

	return s
}
