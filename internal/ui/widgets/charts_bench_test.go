package widgets

import (
	"testing"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/N1xev/bubbleMonitor/internal/data"
)

func BenchmarkRenderSparkline(b *testing.B) {
	ring := data.NewRingBuffer(100)
	for i := 0; i < 100; i++ {
		ring.Push(float64(i % 100))
	}

	c1 := compat.AdaptiveColor{Light: lipgloss.Color("#00FF00"), Dark: lipgloss.Color("#00FF00")}
	c2 := compat.AdaptiveColor{Light: lipgloss.Color("#FF0000"), Dark: lipgloss.Color("#FF0000")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderSparkline(ring, 80, 1, c1, c2, 100.0)
	}
}

func BenchmarkRenderLineChart(b *testing.B) {
	ring := data.NewRingBuffer(100)
	for i := 0; i < 100; i++ {
		ring.Push(float64(i % 100))
	}

	c1 := compat.AdaptiveColor{Light: lipgloss.Color("#00FF00"), Dark: lipgloss.Color("#00FF00")}
	c2 := compat.AdaptiveColor{Light: lipgloss.Color("#FF0000"), Dark: lipgloss.Color("#FF0000")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderLineChart(ring, 80, 10, c1, c2, 100.0)
	}
}

func BenchmarkRenderBrailleChart(b *testing.B) {
	ring := data.NewRingBuffer(100)
	for i := 0; i < 100; i++ {
		ring.Push(float64(i % 100))
	}

	c1 := compat.AdaptiveColor{Light: lipgloss.Color("#00FF00"), Dark: lipgloss.Color("#00FF00")}
	c2 := compat.AdaptiveColor{Light: lipgloss.Color("#FF0000"), Dark: lipgloss.Color("#FF0000")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderBrailleChart(ring, 80, 10, c1, c2, 100.0)
	}
}
