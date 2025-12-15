package data

// RingBuffer is a fixed-size circular buffer for float64 data
type RingBuffer struct {
	data   []float64
	curr   int
	length int
	full   bool
}

// NewRingBuffer creates a new ring buffer of the given size
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]float64, size),
	}
}

// Push adds a value to the buffer
func (r *RingBuffer) Push(val float64) {
	r.data[r.curr] = val
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

// Len returns the current number of elements
func (r *RingBuffer) Len() int {
	return r.length
}

// Get returns the logical element at index i (0 is oldest)
func (r *RingBuffer) Get(i int) float64 {
	if r.length == 0 {
		return 0
	}
	if !r.full {
		if i >= r.length {
			return 0
		}
		return r.data[i]
	}

	// If full, logical 0 is at curr (oldest overwritten)
	// Actually, if full, the oldest data is at r.curr.
	// Valid range is 0 to length-1.
	// Physical index = (r.curr + i) % size
	idx := (r.curr + i) % len(r.data)
	return r.data[idx]
}

// Max returns the maximum value in the buffer
func (r *RingBuffer) Max() float64 {
	if r.length == 0 {
		return 0
	}
	max := -1.7976931348623157e+308 // float64 min

	if r.full {
		max = r.data[0]
		for _, v := range r.data {
			if v > max {
				max = v
			}
		}
	} else {
		max = r.data[0]
		for i := 0; i < r.curr; i++ {
			if r.data[i] > max {
				max = r.data[i]
			}
		}
	}
	return max
}

// Avg returns the average value in the buffer
func (r *RingBuffer) Avg() float64 {
	if r.length == 0 {
		return 0
	}
	sum := 0.0
	// Optimization: iterate only existing data
	count := r.length

	if r.full {
		for _, v := range r.data {
			sum += v
		}
	} else {
		for i := 0; i < r.curr; i++ {
			sum += r.data[i]
		}
	}
	return sum / float64(count)
}

// Accessor interface for charts
type Accessor interface {
	Len() int
	Get(i int) float64
}
