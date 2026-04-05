//go:build linux && cgo

package system

import (
	"log"
	"os/exec"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

func detectNvidia() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("NVIDIA detection panicked: %v", r)
		}
	}()

	ensureNVML()
	if nvmlInitialized {
		count, ret := nvml.DeviceGetCount()
		if ret == nvml.SUCCESS && count > 0 {
			nvidiaDetected.Store(true)
			return
		}
	}

	cmd := exec.Command("nvidia-smi", "--query-gpu=index", "--format=csv,noheader")
	if out, err := cmd.Output(); err == nil && len(out) > 0 {
		nvidiaDetected.Store(true)
	}
}
