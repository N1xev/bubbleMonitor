package system

import (
	tea "charm.land/bubbletea/v2"
	"github.com/shirou/gopsutil/v3/net"

	"github.com/N1xev/bubbleMonitor/src/messages"
)

// NetworkInterfacesCmd fetches network interface stats
// NetworkInterfacesCmd fetches network interface stats
func NetworkInterfacesCmd() tea.Cmd {
	return func() tea.Msg {
		ioCounters, _ := net.IOCounters(true)
		return messages.NetworkInterfacesMsg(ioCounters)
	}
}
