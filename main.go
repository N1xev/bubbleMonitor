package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/N1xev/bubbleMonitor/src/model"
)

func main() {
	p := tea.NewProgram(model.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
