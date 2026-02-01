package remote

import (
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"github.com/N1xev/bubbleMonitor/internal/msg"
)

func CheckRemoteCmd(host string) tea.Cmd {
	return func() tea.Msg {
		// Runs: ssh -o ConnectTimeout=2 host uptime
		cmd := exec.Command("ssh", "-o", "ConnectTimeout=2", "-o", "BatchMode=yes", host, "uptime")
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
