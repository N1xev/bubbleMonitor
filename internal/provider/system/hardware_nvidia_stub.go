//go:build !linux || !cgo

package system

import (
	"sync"
	"sync/atomic"

	"github.com/N1xev/bubbleMonitor/internal/data"
)

var (
	nvmlInitialized   bool
	nvmlInitAttempted atomic.Bool
	nvmlInitOnce      sync.Once
)

func ensureNVML() {
	// NVML not available on this platform
}

func fetchNvidiaGpus() []data.GpuInfo {
	return nil
}
