package process

import (
	tea "charm.land/bubbletea/v2"
	"github.com/shirou/gopsutil/v3/process"

	"github.com/N1xev/bubbleMonitor/internal/msg"
)

// FetchOpenFilesCmd retrieves open files for a process
func FetchOpenFilesCmd(pid int32) tea.Cmd {
	return func() tea.Msg {
		proc, err := process.NewProcess(pid)
		if err != nil {
			return msg.OpenFilesMsg{Pid: pid, Err: err}
		}
		files, err := proc.OpenFiles()
		if err != nil {
			return msg.OpenFilesMsg{Pid: pid, Err: err}
		}
		return msg.OpenFilesMsg{Pid: pid, Files: files}
	}
}
