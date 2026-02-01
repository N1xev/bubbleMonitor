package system

import (
	"os/exec"
	"runtime"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/msg"
)

func ServicesCmd() tea.Cmd {
	return func() tea.Msg {
		var services []data.ServiceInfo

		if runtime.GOOS == "linux" {
			// systemctl list-units --type=service --no-legend
			cmd := exec.Command("systemctl", "list-units", "--type=service", "--no-legend", "--all")
			out, err := cmd.Output()
			if err == nil {
				lines := strings.Split(string(out), "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) == "" {
						continue
					}
					parts := strings.Fields(line)
					if len(parts) >= 4 {
						// Format: unit load active sub description...
						services = append(services, data.ServiceInfo{
							Name:        parts[0],
							Status:      parts[3], // sub-state: running, dead, exited
							Description: strings.Join(parts[4:], " "),
						})
					}
				}
			}
		}

		return msg.ServicesMsg(services)
	}
}
