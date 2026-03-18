package system

import (
	"os/exec"
	"runtime"
	"sync/atomic"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/net"

	"github.com/N1xev/bubbleMonitor/internal/data"
)

var (
	detectionDone    atomic.Bool
	nvidiaDetected   atomic.Bool
	amdDetected      atomic.Bool
	intelDetected    atomic.Bool
	batteryDetected  atomic.Bool
	networkDetected  atomic.Bool
	diskIODetected   atomic.Bool
	servicesDetected atomic.Bool
	tempDetected     atomic.Bool
	dockerDetected   atomic.Bool
	k8sDetected      atomic.Bool
)

func DetectHardware() *data.HardwareCapabilities {
	if detectionDone.Load() {
		return getCapabilities()
	}

	detectNvidia()
	detectAmd()
	detectBattery()
	detectNetwork()
	detectDiskIO()
	detectTemp()
	detectServices()
	detectContainers()

	detectionDone.Store(true)

	return getCapabilities()
}

func getCapabilities() *data.HardwareCapabilities {
	return &data.HardwareCapabilities{
		HasNvidiaGPU:         nvidiaDetected.Load(),
		HasAmdGPU:            amdDetected.Load(),
		HasIntelGPU:          intelDetected.Load(),
		HasBattery:           batteryDetected.Load(),
		HasNetworkInterfaces: networkDetected.Load(),
		HasDiskIO:            diskIODetected.Load(),
		HasServices:          servicesDetected.Load(),
		HasTempSensors:       tempDetected.Load(),
		HasDocker:            dockerDetected.Load(),
		HasKubernetes:        k8sDetected.Load(),
	}
}

func detectNvidia() {
	defer func() {
		_ = recover()
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

func detectAmd() {
	defer func() {
		_ = recover()
	}()

	if runtime.GOOS != "linux" {
		return
	}

	_, err := exec.LookPath("rocm-smi")
	if err != nil {
		return
	}

	amdDetected.Store(true)
}

func detectBattery() {
	defer func() {
		_ = recover()
	}()

	_, err := exec.LookPath("upower")
	if err != nil {
		if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
			batteryDetected.Store(true)
		}
		return
	}

	batteryDetected.Store(true)
}

func detectNetwork() {
	defer func() {
		_ = recover()
	}()

	interfaces, err := net.IOCounters(false)
	if err == nil && len(interfaces) > 0 {
		networkDetected.Store(true)
	}
}

func detectDiskIO() {
	defer func() {
		_ = recover()
	}()

	_, err := disk.IOCounters()
	if err == nil {
		diskIODetected.Store(true)
	}
}

func detectTemp() {
	defer func() {
		_ = recover()
	}()

	sensors, err := host.SensorsTemperatures()
	if err == nil && len(sensors) > 0 {
		tempDetected.Store(true)
	}
}

func detectServices() {
	defer func() {
		_ = recover()
	}()

	if runtime.GOOS == "linux" {
		_, err := exec.LookPath("systemctl")
		if err == nil {
			servicesDetected.Store(true)
		}
	} else if runtime.GOOS == "windows" {
		servicesDetected.Store(true)
	}
}

func detectContainers() {
	defer func() {
		_ = recover()
	}()

	dockerDetected.Store(HasDocker())
	k8sDetected.Store(HasKubernetes())
}
