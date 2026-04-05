//go:build !linux || !cgo

package system

import "github.com/N1xev/bubbleMonitor/internal/data"

func fetchAmdGpus() []data.GpuInfo {
	return nil
}
