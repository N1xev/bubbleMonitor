package data

import (
	"sync"
	"testing"
)

func TestRingBufferConcurrency(t *testing.T) {
	rb := NewRingBuffer(100)
	var wg sync.WaitGroup
	numGoroutines := 10
	iterations := 1000

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(val float64) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				rb.Push(val)
				_ = rb.Max()
				_ = rb.Avg()
				_ = rb.Len()
				_ = rb.Get(0)
			}
		}(float64(i))
	}

	wg.Wait()

	if rb.Len() > 100 {
		t.Errorf("Buffer overflow: length %d exceeds capacity 100", rb.Len())
	}
}

func TestRingBufferBasicOperations(t *testing.T) {
	rb := NewRingBuffer(5)

	if rb.Len() != 0 {
		t.Errorf("New buffer should have length 0, got %d", rb.Len())
	}

	rb.Push(10.0)
	rb.Push(20.0)
	rb.Push(30.0)

	if rb.Len() != 3 {
		t.Errorf("Expected length 3, got %d", rb.Len())
	}

	if max := rb.Max(); max != 30.0 {
		t.Errorf("Expected max 30.0, got %f", max)
	}

	expectedAvg := 20.0
	if avg := rb.Avg(); avg != expectedAvg {
		t.Errorf("Expected avg %f, got %f", expectedAvg, avg)
	}
}

func TestRingBufferWrapAround(t *testing.T) {
	rb := NewRingBuffer(3)

	rb.Push(1.0)
	rb.Push(2.0)
	rb.Push(3.0)
	rb.Push(4.0)
	rb.Push(5.0)

	if rb.Len() != 3 {
		t.Errorf("Expected length 3 after wraparound, got %d", rb.Len())
	}

	if max := rb.Max(); max != 5.0 {
		t.Errorf("Expected max 5.0 after wraparound, got %f", max)
	}
}

func TestRingBufferEmptyState(t *testing.T) {
	rb := NewRingBuffer(10)

	if rb.Len() != 0 {
		t.Error("New buffer should be empty")
	}

	if rb.Max() != 0 {
		t.Error("Max of empty buffer should be 0")
	}

	if rb.Avg() != 0 {
		t.Error("Avg of empty buffer should be 0")
	}

	if rb.Get(0) != 0 {
		t.Error("Get on empty buffer should return 0")
	}
}

func BenchmarkRingBufferPush(b *testing.B) {
	rb := NewRingBuffer(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Push(float64(i))
	}
}

func BenchmarkRingBufferConcurrentAccess(b *testing.B) {
	rb := NewRingBuffer(1000)
	for i := 0; i < 100; i++ {
		rb.Push(float64(i))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rb.Push(42.0)
			_ = rb.Max()
			_ = rb.Avg()
		}
	})
}
