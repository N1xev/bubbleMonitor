package system

import (
	"log"
	"os/exec"
	"runtime"
	"sync/atomic"

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

func detectAmd() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("AMD detection panicked: %v", r)
		}
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
		if r := recover(); r != nil {
			log.Printf("Battery detection panicked: %v", r)
		}
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
		if r := recover(); r != nil {
			log.Printf("Network detection panicked: %v", r)
		}
	}()

	interfaces, err := net.IOCounters(false)
	if err == nil && len(interfaces) > 0 {
		networkDetected.Store(true)
	}
}

func detectDiskIO() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("DiskIO detection panicked: %v", r)
		}
	}()

	_, err := disk.IOCounters()
	if err == nil {
		diskIODetected.Store(true)
	}
}

func detectTemp() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Temperature detection panicked: %v", r)
		}
	}()

	sensors, err := host.SensorsTemperatures()
	if err == nil && len(sensors) > 0 {
		tempDetected.Store(true)
	}
}

func detectServices() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Services detection panicked: %v", r)
		}
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
		if r := recover(); r != nil {
			log.Printf("Container detection panicked: %v", r)
		}
	}()

	dockerDetected.Store(HasDocker())
	k8sDetected.Store(HasKubernetes())
}
