package process

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/N1xev/bubbleMonitor/src/messages"
	"github.com/shirou/gopsutil/v3/process"
)

// ReniceProcessCmd changes the priority of a process by a delta
// delta < 0 increases priority (on Unix), delta > 0 decreases priority
func ReniceProcessCmd(pid int32, delta int) tea.Cmd {
	return func() tea.Msg {

		// WINDOWS LOGIC
		// Priorities: Idle(64) < Below(16384) < Normal(32) < Above(32768) < High(128) < Realtime(256)
		// We use a predefined list to step through.
		winPriorities := []int32{64, 16384, 32, 32768, 128, 256}

		proc, err := process.NewProcess(pid)
		if err != nil {
			return messages.PriorityChangeMsg{Pid: pid, Err: err}
		}

		currentNice, err := proc.Nice()
		if err != nil {
			return messages.PriorityChangeMsg{Pid: pid, Err: err}
		}

		var newPrio int32

		if currentNice > 20 || currentNice < -20 {
			// Windows Logic
			currIdx := -1
			for i, p := range winPriorities {
				if p == currentNice {
					currIdx = i
					break
				}
			}

			// If not found, maybe default to Normal(32)?
			if currIdx == -1 {
				currIdx = 2 // Normal
			}

			newIdx := currIdx
			if delta < 0 {
				// Increase Priority (Step Up in Windows List)
				newIdx++
			} else {
				// Decrease Priority (Step Down)
				newIdx--
			}

			// Clamp
			if newIdx < 0 {
				newIdx = 0
			}
			if newIdx >= len(winPriorities) {
				newIdx = len(winPriorities) - 1
			}

			newPrio = winPriorities[newIdx]

			// Execute wmic
			// wmic process where processid=PID call setpriority VAL
			cmd := exec.Command("wmic", "process", "where", "processid="+strings.TrimSpace(fmt.Sprint(pid)), "call", "setpriority", fmt.Sprint(newPrio))
			// Only for Windows
			if err := cmd.Run(); err != nil {
				return messages.PriorityChangeMsg{Pid: pid, Priority: newPrio, Err: err}
			}
			// Success
			return messages.PriorityChangeMsg{Pid: pid, Priority: newPrio, Err: nil}

		} else {
			// Unix Logic
			// New = Current + Delta
			// If delta < 0 (e.g. -1), New < Current (Higher Priority)
			newPrio = currentNice + int32(delta)

			// Use generic Renice if gopsutil supported it, but it doesn't.
			// Use "renice" command
			cmd := exec.Command("renice", "-n", fmt.Sprint(newPrio), "-p", fmt.Sprint(pid))
			if err := cmd.Run(); err != nil {
				return messages.PriorityChangeMsg{Pid: pid, Priority: newPrio, Err: err}
			}
			return messages.PriorityChangeMsg{Pid: pid, Priority: newPrio, Err: nil}
		}
	}
}

// ReniceProcessCmdSafe changes the priority of a process by a delta
// delta < 0 increases priority, delta > 0 decreases priority
func ReniceProcessCmdSafe(pid int32, delta int) tea.Cmd {
	return func() tea.Msg {
		proc, err := process.NewProcess(pid)
		if err != nil {
			return messages.PriorityChangeMsg{Pid: pid, Err: err}
		}

		currentNice, err := proc.Nice()
		if err != nil {
			return messages.PriorityChangeMsg{Pid: pid, Err: err}
		}

		if runtime.GOOS == "windows" {
			// Window Priority Classes (wmic Expects):
			// 64=Idle, 16384=Below, 32=Normal, 32768=Above, 128=High, 256=Realtime
			winPriorities := []int32{64, 16384, 32, 32768, 128, 256}

			// Gopsutil returns Base Priority (Level), not Class ID.
			// Mapping: 4->Idle(0), 6->Below(1), 8->Normal(2), 10->Above(3), 13->High(4), 24->Realtime(5)

			var currIdx int
			switch currentNice {
			case 4:
				currIdx = 0
			case 6:
				currIdx = 1
			case 8:
				currIdx = 2
			case 10:
				currIdx = 3
			case 13:
				currIdx = 4
			case 24:
				currIdx = 5
			default:
				// Fallback: check if it matched a Class ID directly?
				currIdx = -1
				for i, p := range winPriorities {
					if p == currentNice {
						currIdx = i
						break
					}
				}
				if currIdx == -1 {
					currIdx = 2 // Default Normal
				}
			}

			newIdx := currIdx
			if delta < 0 {
				newIdx++ // Increase Priority (Step Up)
			} else {
				newIdx-- // Decrease Priority (Step Down)
			}

			if newIdx < 0 {
				newIdx = 0
			}
			if newIdx >= len(winPriorities) {
				newIdx = len(winPriorities) - 1
			}

			newPrio := winPriorities[newIdx]

			cmd := exec.Command("wmic", "process", "where", "processid="+strings.TrimSpace(fmt.Sprint(pid)), "call", "setpriority", fmt.Sprint(newPrio))
			if err := cmd.Run(); err != nil {
				return messages.PriorityChangeMsg{Pid: pid, Priority: newPrio, Err: err}
			}
			return messages.PriorityChangeMsg{Pid: pid, Priority: newPrio, Err: nil}
		}

		newPrio := currentNice + int32(delta)
		cmd := exec.Command("renice", "-n", fmt.Sprint(newPrio), "-p", fmt.Sprint(pid))
		if err := cmd.Run(); err != nil {
			return messages.PriorityChangeMsg{Pid: pid, Priority: newPrio, Err: err}
		}
		return messages.PriorityChangeMsg{Pid: pid, Priority: newPrio, Err: nil}
	}
}

// SuspendProcessCmd suspends a process
func SuspendProcessCmd(pid int32) tea.Cmd {
	return func() tea.Msg {
		proc, err := process.NewProcess(pid)
		if err != nil {
			return messages.ProcessControlMsg{Pid: pid, Action: "suspend", Err: err}
		}
		err = proc.Suspend()
		return messages.ProcessControlMsg{Pid: pid, Action: "suspend", Err: err}
	}
}

// ResumeProcessCmd resumes a process
func ResumeProcessCmd(pid int32) tea.Cmd {
	return func() tea.Msg {
		proc, err := process.NewProcess(pid)
		if err != nil {
			return messages.ProcessControlMsg{Pid: pid, Action: "resume", Err: err}
		}
		err = proc.Resume()
		return messages.ProcessControlMsg{Pid: pid, Action: "resume", Err: err}
	}
}
