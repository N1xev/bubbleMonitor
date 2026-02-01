package system

import (
	"strconv"

	tea "charm.land/bubbletea/v2"
	"github.com/N1xev/bubbleMonitor/internal/data"
	"github.com/N1xev/bubbleMonitor/internal/msg"
	"github.com/shirou/gopsutil/v3/net"
)

func ConnectionsCmd() tea.Cmd {
	return func() tea.Msg {
		conns, err := net.Connections("all")
		if err != nil {
			return msg.ConnectionsMsg{}
		}

		var result []data.ConnectionInfo
		for _, c := range conns {
			result = append(result, data.ConnectionInfo{
				LocalAddr:  c.Laddr.IP + ":" + strconv.Itoa(int(c.Laddr.Port)),
				RemoteAddr: c.Raddr.IP + ":" + strconv.Itoa(int(c.Raddr.Port)),
				State:      c.Status,
				Pid:        c.Pid,
				Protocol:   getProto(c.Type),
			})
		}

		return msg.ConnectionsMsg(result)
	}
}

func getProto(t uint32) string {
	switch t {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	default:
		return "UNK"
	}
}
