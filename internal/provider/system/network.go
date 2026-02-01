package system

import (
	tea "charm.land/bubbletea/v2"
	"github.com/shirou/gopsutil/v3/net"

	"github.com/N1xev/bubbleMonitor/internal/msg"
)

// NetworkInterfacesCmd fetches network interface stats
func NetworkInterfacesCmd() tea.Cmd {
	return func() tea.Msg {
		ioCounters, _ := net.IOCounters(true)
		return msg.NetworkInterfacesMsg{
			Interfaces: ioCounters,
			Err:        nil,
		}
	}
}
