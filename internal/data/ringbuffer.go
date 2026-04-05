package data

import "sync"

type RingBuffer struct {
	data      []float64
	mu        sync.RWMutex // Protects all fields below
	curr      int
	length    int
	cachedMax float64
	cachedSum float64
	full      bool
	maxDirty  bool
	sumDirty  bool
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data:     make([]float64, size),
		maxDirty: true,
		sumDirty: true,
	}
}

func (r *RingBuffer) Push(val float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	oldVal := r.data[r.curr]
	r.data[r.curr] = val

	if r.full {
		r.cachedSum = r.cachedSum - oldVal + val
		if oldVal == r.cachedMax || val > r.cachedMax {
			r.maxDirty = true
		}
	} else {
		r.cachedSum += val
		if r.length == 0 || val > r.cachedMax {
			r.cachedMax = val
		}
		r.sumDirty = false
	}

	r.curr++
	if r.curr >= len(r.data) {
		r.curr = 0
		r.full = true
	}
	if !r.full {
		r.length = r.curr
	} else {
		r.length = len(r.data)
	}
}

func (r *RingBuffer) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.length
}

func (r *RingBuffer) Get(i int) float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.length == 0 {
		return 0
	}
	if !r.full {
		if i >= r.length {
			return 0
		}
		return r.data[i]
	}

	idx := (r.curr + i) % len(r.data)
	return r.data[idx]
}

func (r *RingBuffer) Max() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.length == 0 {
		return 0
	}

	if !r.maxDirty {
		return r.cachedMax
	}

	max := r.data[0]
	if r.full {
		for _, v := range r.data {
			if v > max {
				max = v
			}
		}
	} else {
		for i := 0; i < r.curr; i++ {
			if r.data[i] > max {
				max = r.data[i]
			}
		}
	}

	r.cachedMax = max
	r.maxDirty = false
	return max
}

func (r *RingBuffer) Avg() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.length == 0 {
		return 0
	}

	if !r.sumDirty {
		return r.cachedSum / float64(r.length)
	}

	sum := 0.0
	if r.full {
		for _, v := range r.data {
			sum += v
		}
	} else {
		for i := 0; i < r.curr; i++ {
			sum += r.data[i]
		}
	}

	r.cachedSum = sum
	r.sumDirty = false
	return sum / float64(r.length)
}

type Accessor interface {
	Len() int
	Get(i int) float64
	Max() float64
}

// Resize creates a new buffer with the given size, preserving as much
// existing data as possible (newest data is kept).
func (r *RingBuffer) Resize(newSize int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if newSize == len(r.data) {
		return
	}

	newData := make([]float64, newSize)

	// Copy existing data, keeping the newest entries
	count := r.length
	if count > newSize {
		count = newSize
	}

	// Read from oldest to newest
	for i := 0; i < count; i++ {
		if r.full {
			idx := (r.curr + i) % len(r.data)
			newData[i] = r.data[idx]
		} else {
			newData[i] = r.data[i]
		}
	}

	r.data = newData
	r.curr = count % newSize
	r.length = count
	r.full = count == newSize
	r.maxDirty = true
	r.sumDirty = true
}
