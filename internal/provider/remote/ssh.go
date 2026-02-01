package remote

import (
	"fmt"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"github.com/N1xev/bubbleMonitor/internal/msg"
	"github.com/N1xev/bubbleMonitor/internal/provider"
)

func CheckRemoteCmd(host string) tea.Cmd {
	return func() tea.Msg {
		// Runs: ssh -o ConnectTimeout=2 host uptime
		timeoutOpt := fmt.Sprintf("ConnectTimeout=%d", int(provider.SSHTimeoutSeconds.Seconds()))
		cmd := exec.Command("ssh", "-o", timeoutOpt, "-o", "BatchMode=yes", host, "uptime")
		out, err := cmd.Output()
		result := string(out)
		if err != nil {
			result = "Error: " + err.Error()
		}
		if result == "" {
			result = "No output"
		}
		return msg.RemoteMsg{Host: host, Output: result}
	}
}
