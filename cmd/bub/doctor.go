package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/spf13/cobra"

	configpkg "github.com/N1xev/bubbleMonitor/internal/config"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Run diagnostic checks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor(cmd)
		},
	}
}

func runDoctor(cmd *cobra.Command) error {
	out := cmd.OutOrStdout()
	failures := 0

	fmt.Fprintf(out, "\n")

	// Config file
	configPath, err := configpkg.GetConfigPath()
	if err != nil {
		fmt.Fprintf(out, "  ✗ Config path: %v\n", err)
		failures++
	} else {
		if _, err := os.Stat(configPath); err != nil {
			fmt.Fprintf(out, "  ⚠ Config file: not found at %s (will be created on first run)\n", configPath)
		} else {
			fmt.Fprintf(out, "  ✓ Config file: %s\n", configPath)
		}
	}

	// gopsutil - CPU metrics
	if pcts, err := cpu.Percent(0, false); err != nil || len(pcts) == 0 {
		fmt.Fprintf(out, "  ✗ CPU metrics: unavailable (%v)\n", err)
		failures++
	} else {
		fmt.Fprintf(out, "  ✓ CPU metrics: available (%.1f%%)\n", pcts[0])
	}

	// Host info
	if hi, err := host.Info(); err != nil {
		fmt.Fprintf(out, "  ✗ Host info: unavailable (%v)\n", err)
		failures++
	} else {
		fmt.Fprintf(out, "  ✓ Host info: %s %s\n", hi.OS, hi.PlatformVersion)
	}

	// GPU checks
	checkGPU(out, "NVIDIA", "nvidia-smi")
	checkGPU(out, "AMD", "rocm-smi")

	// SSH binary
	if path, err := exec.LookPath("ssh"); err != nil {
		fmt.Fprintf(out, "  ✗ SSH: not found in PATH\n")
		failures++
	} else {
		fmt.Fprintf(out, "  ✓ SSH binary: %s\n", path)
	}

	// Remote hosts connectivity
	cfg, _ := loadConfigWithOverrides()
	for _, h := range cfg.RemoteHosts {
		checkRemoteHost(out, h)
	}

	fmt.Fprintf(out, "\n")
	if failures > 0 {
		os.Exit(1)
	}
	return nil
}

func checkGPU(out io.Writer, name, cmd string) {
	if path, err := exec.LookPath(cmd); err != nil {
		fmt.Fprintf(out, "  ⚠ %s GPU: %s not found\n", name, cmd)
	} else {
		fmt.Fprintf(out, "  ✓ %s GPU: %s detected\n", name, path)
	}
}

func checkRemoteHost(out io.Writer, host configpkg.RemoteHostConfig) {
	args := []string{"-o", "ConnectTimeout=5", "-o", "BatchMode=yes", host.Host, "echo", "ok"}
	start := time.Now()
	cmd := exec.Command("ssh", args...)
	err := cmd.Run()
	elapsed := time.Since(start).Truncate(time.Millisecond)

	if err != nil {
		fmt.Fprintf(out, "  ✗ Remote %q (%s): unreachable (%v)\n", host.Name, host.Host, err)
	} else {
		fmt.Fprintf(out, "  ✓ Remote %q (%s): reachable (%v)\n", host.Name, host.Host, elapsed)
	}
}
