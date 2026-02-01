package system

import (
	"os/exec"
	"runtime"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/N1xev/bubbleMonitor/internal/msg"
)

func SystemLogsCmd() tea.Cmd {
	return func() tea.Msg {
		var logs []string

		if runtime.GOOS == "linux" {
			// journalctl -n 50 --no-pager -o short-iso
			cmd := exec.Command("journalctl", "-n", "50", "--no-pager", "-o", "short-iso")
			out, err := cmd.Output()
			if err == nil {
				lines := strings.Split(string(out), "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						logs = append(logs, line)
					}
				}
			} else {
				// Fallback to reading /var/log/syslog if permission denied or no journalctl?
				// Just return error or empty
				logs = append(logs, "Error reading logs or no permission: "+err.Error())
			}
		} else {
			logs = append(logs, "Logs not implemented for "+runtime.GOOS)
		}

		return msg.SysLogMsg(logs)
	}
}
