package data

type RingBuffer struct {
	data   []float64
	curr   int
	length int
	full   bool

	cachedMax float64
	cachedSum float64
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
	return r.length
}

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

	idx := (r.curr + i) % len(r.data)
	return r.data[idx]
}

func (r *RingBuffer) Max() float64 {
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
