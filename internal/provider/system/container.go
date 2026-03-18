package system

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/msg"
)

var (
	containerChecked   bool
	dockerAvailable    bool
	k8sAvailable       bool
	containerCheckOnce sync.Once
)

func init() {
	containerCheckOnce.Do(checkContainerAvailability)
}

func checkContainerAvailability() {
	if containerChecked {
		return
	}
	dockerAvailable = isDockerAvailable()
	k8sAvailable = isKubernetesAvailable()
	containerChecked = true
}

func isDockerAvailable() bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("docker", "info")
		_ = cmd.Run()
		return cmd.ProcessState.ExitCode() == 0
	}
	conn, err := net.DialTimeout("unix", "/var/run/docker.sock", time.Second)
	if err != nil {
		conn, err = net.DialTimeout("unix", "/run/docker.sock", time.Second)
	}
	if err == nil {
		_ = conn.Close()
		return true
	}
	return false
}

func isKubernetesAvailable() bool {
	cmd := exec.Command("kubectl", "cluster-info")
	err := cmd.Run()
	return err == nil
}

func HasDocker() bool {
	return dockerAvailable
}

func HasKubernetes() bool {
	return k8sAvailable
}

func ContainerCmd() tea.Cmd {
	return func() tea.Msg {
		var containers []data.ContainerInfo
		var pods []data.K8sPodInfo
		var err error

		if dockerAvailable {
			containers, err = fetchDockerContainers()
			if err != nil {
				containers = nil
			}
		}

		if k8sAvailable {
			pods, err = fetchKubernetesPods()
			if err != nil {
				pods = nil
			}
		}

		return msg.ContainerInfoMsg{
			Err:           err,
			Containers:    containers,
			Pods:          pods,
			HasDocker:     dockerAvailable,
			HasKubernetes: k8sAvailable,
		}
	}
}

type dockerStats struct {
	Read     string         `json:"read"`
	CPU      dockerCPUStats `json:"cpu_stats"`
	PreCPU   dockerCPUStats `json:"precpu_stats"`
	Memory   dockerMemStats `json:"memory_stats"`
	Networks map[string]struct {
		RxBytes uint64 `json:"rx_bytes"`
		TxBytes uint64 `json:"tx_bytes"`
	} `json:"networks"`
}

type dockerCPUStats struct {
	CPUUsage struct {
		TotalUsage uint64 `json:"total_usage"`
	} `json:"cpu_usage"`
	SystemCPUUsage uint64  `json:"system_cpu_usage"`
	OnlineCPUs     float64 `json:"online_cpus"`
}

type dockerMemStats struct {
	Usage uint64 `json:"usage"`
	Limit uint64 `json:"limit"`
}

func fetchDockerContainers() ([]data.ContainerInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "ps", "-a", "--format", "{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t{{.State}}\t{{.CreatedAt}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	result := make([]data.ContainerInfo, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 6 {
			continue
		}

		ci := data.ContainerInfo{
			Name:    parts[1],
			ID:      parts[0],
			Status:  parts[3],
			State:   parts[4],
			Image:   parts[2],
			Created: parts[5],
			Type:    "docker",
		}

		stats, err := fetchDockerContainerStats(ctx, ci.ID)
		if err == nil && stats != nil {
			ci.CPUPercent = calculateCPUPercent(stats)
			ci.MemUsage = stats.Memory.Usage
			ci.MemLimit = stats.Memory.Limit
			if stats.Memory.Limit > 0 {
				ci.MemPct = float64(stats.Memory.Usage) / float64(stats.Memory.Limit) * 100
			}
			for _, netStats := range stats.Networks {
				ci.NetRx += netStats.RxBytes
				ci.NetTx += netStats.TxBytes
			}
		}

		result = append(result, ci)
	}

	return result, nil
}

func fetchDockerContainerStats(ctx context.Context, containerID string) (*dockerStats, error) {
	cmd := exec.CommandContext(ctx, "docker", "stats", "--no-stream", "--format", "{{json .}}")
	cmd.Stdin = strings.NewReader(containerID + "\n")

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return nil, fmt.Errorf("empty stats output")
	}

	var rawStats struct {
		CPUPerc  string `json:"CPUPerc"`
		MemUsage string `json:"MemUsage"`
		MemPerc  string `json:"MemPerc"`
		NetRx    string `json:"NetRx"`
		NetTx    string `json:"NetTx"`
	}

	if err := json.Unmarshal([]byte(output), &rawStats); err != nil {
		return nil, err
	}

	stats := &dockerStats{}

	if cpu, err := strconv.ParseFloat(strings.TrimSuffix(rawStats.CPUPerc, "%"), 64); err == nil {
		stats.CPU.OnlineCPUs = 1
		stats.CPU.CPUUsage.TotalUsage = uint64(cpu * 1000000)
	}

	if memParts := strings.Split(rawStats.MemUsage, "/"); len(memParts) == 2 {
		if memUsage, err := parseMemValue(memParts[0]); err == nil {
			stats.Memory.Usage = memUsage
		}
		if memLimit, err := parseMemValue(memParts[1]); err == nil {
			stats.Memory.Limit = memLimit
		}
	}

	stats.Networks = make(map[string]struct {
		RxBytes uint64 `json:"rx_bytes"`
		TxBytes uint64 `json:"tx_bytes"`
	})

	if netRx, err := parseMemValue(rawStats.NetRx); err == nil {
		if stats.Networks == nil {
			stats.Networks = make(map[string]struct {
				RxBytes uint64 `json:"rx_bytes"`
				TxBytes uint64 `json:"tx_bytes"`
			})
		}
		stats.Networks["eth0"] = struct {
			RxBytes uint64 `json:"rx_bytes"`
			TxBytes uint64 `json:"tx_bytes"`
		}{RxBytes: netRx}
	}

	if netTx, err := parseMemValue(rawStats.NetTx); err == nil {
		if stats.Networks == nil {
			stats.Networks = make(map[string]struct {
				RxBytes uint64 `json:"rx_bytes"`
				TxBytes uint64 `json:"tx_bytes"`
			})
		}
		if eth0, ok := stats.Networks["eth0"]; ok {
			eth0.TxBytes = netTx
			stats.Networks["eth0"] = eth0
		}
	}

	return stats, nil
}

func parseMemValue(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	multiplier := uint64(1)

	if strings.HasSuffix(s, "GIB") || strings.HasSuffix(s, "GB") {
		s = strings.TrimSuffix(s, "GIB")
		s = strings.TrimSuffix(s, "GB")
		multiplier = 1024 * 1024 * 1024
	} else if strings.HasSuffix(s, "MIB") || strings.HasSuffix(s, "MB") {
		s = strings.TrimSuffix(s, "MIB")
		s = strings.TrimSuffix(s, "MB")
		multiplier = 1024 * 1024
	} else if strings.HasSuffix(s, "KIB") || strings.HasSuffix(s, "KB") {
		s = strings.TrimSuffix(s, "KIB")
		s = strings.TrimSuffix(s, "KB")
		multiplier = 1024
	} else if strings.HasSuffix(s, "B") {
		s = strings.TrimSuffix(s, "B")
	}

	s = strings.TrimSpace(s)
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return uint64(num * float64(multiplier)), nil
}

func calculateCPUPercent(stats *dockerStats) float64 {
	cpuDelta := float64(stats.CPU.CPUUsage.TotalUsage) - float64(stats.PreCPU.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPU.SystemCPUUsage) - float64(stats.PreCPU.SystemCPUUsage)

	if systemDelta > 0 && cpuDelta > 0 {
		onlineCPUs := stats.CPU.OnlineCPUs
		if onlineCPUs == 0 {
			onlineCPUs = 1
		}
		return (cpuDelta / systemDelta) * onlineCPUs * 100
	}

	return 0
}

func fetchKubernetesPods() ([]data.K8sPodInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "kubectl", "get", "pods", "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var k8sResult struct {
		Items []struct {
			Metadata struct {
				Name              string `json:"name"`
				Namespace         string `json:"namespace"`
				CreationTimestamp string `json:"creationTimestamp"`
			} `json:"metadata"`
			Status struct {
				Phase      string `json:"phase"`
				PodIP      string `json:"podIP"`
				HostIP     string `json:"hostIP"`
				Conditions []struct {
					Type   string `json:"type"`
					Status string `json:"status"`
				} `json:"conditions"`
				ContainerStatuses []struct {
					Name         string `json:"name"`
					Ready        bool   `json:"ready"`
					RestartCount int    `json:"restartCount"`
					State        struct {
						Running struct{} `json:"running"`
						Waiting struct {
							Reason string `json:"reason"`
						} `json:"waiting"`
						Terminated struct {
							Reason string `json:"reason"`
						} `json:"terminated"`
					} `json:"state"`
				} `json:"containerStatuses"`
			} `json:"status"`
			Spec struct {
				NodeName   string `json:"nodeName"`
				Containers []struct {
					Name      string `json:"name"`
					Resources struct {
						Requests struct {
							CPU    string `json:"cpu"`
							Memory string `json:"memory"`
						} `json:"requests"`
						Limits struct {
							CPU    string `json:"cpu"`
							Memory string `json:"memory"`
						} `json:"limits"`
					} `json:"resources"`
				} `json:"containers"`
			} `json:"spec"`
		} `json:"items"`
	}

	if err := json.Unmarshal(out, &k8sResult); err != nil {
		return nil, err
	}

	pods := make([]data.K8sPodInfo, 0, len(k8sResult.Items))

	for _, item := range k8sResult.Items {
		readyCount := 0
		totalCount := len(item.Status.ContainerStatuses)
		for _, cs := range item.Status.ContainerStatuses {
			if cs.Ready {
				readyCount++
			}
		}

		status := item.Status.Phase
		if len(item.Status.ContainerStatuses) > 0 {
			cs := item.Status.ContainerStatuses[0]
			if cs.State.Waiting.Reason != "" {
				status = cs.State.Waiting.Reason
			} else if cs.State.Terminated.Reason != "" {
				status = cs.State.Terminated.Reason
			}
		}

		cpuReq := ""
		memReq := ""
		cpuLim := ""
		memLim := ""
		if len(item.Spec.Containers) > 0 {
			c := item.Spec.Containers[0]
			cpuReq = c.Resources.Requests.CPU
			memReq = c.Resources.Requests.Memory
			cpuLim = c.Resources.Limits.CPU
			memLim = c.Resources.Limits.Memory
		}

		pod := data.K8sPodInfo{
			Name:      item.Metadata.Name,
			Namespace: item.Metadata.Namespace,
			Status:    status,
			Ready:     fmt.Sprintf("%d/%d", readyCount, totalCount),
			Node:      item.Spec.NodeName,
			CPUReq:    cpuReq,
			MemReq:    memReq,
			CPULim:    cpuLim,
			MemLim:    memLim,
			Age:       item.Metadata.CreationTimestamp,
		}

		if totalCount > 0 {
			pod.Restarts = item.Status.ContainerStatuses[0].RestartCount
		}

		pods = append(pods, pod)
	}

	return pods, nil
}
