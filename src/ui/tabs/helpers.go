package tabs

import (
	"github.com/N1xev/bubbleMonitor/src/data"
)

// CalcNetPercent calculates network usage percentage
func CalcNetPercent(s *data.AppState) float64 {
	total := s.NetSentRate + s.NetRecvRate
	netP := (total / 10) * 100
	if netP > 100 {
		netP = 100
	}
	return netP
}
