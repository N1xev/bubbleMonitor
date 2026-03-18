package data

import (
	"testing"
)

func BenchmarkBuildProcessTree(b *testing.B) {
	state := &AppState{
		Process: ProcessState{
			CollapsedPids: make(map[int32]bool),
		},
	}

	procs := make([]ProcessInfo, 500)
	for i := 0; i < 500; i++ {
		procs[i] = ProcessInfo{
			Pid:    int32(i + 1),
			Ppid:   int32((i / 10) * 10),
			Name:   "test_process",
			Cpu:    50.0,
			Memory: 30.0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = state.buildProcessTree(procs)
	}
}

func BenchmarkGetFilteredProcesses(b *testing.B) {
	state := &AppState{
		Process: ProcessState{
			ProcessFilter: "test",
		},
	}

	state.Process.Processes = make([]ProcessInfo, 500)
	for i := 0; i < 500; i++ {
		state.Process.Processes[i] = ProcessInfo{
			Pid:           int32(i + 1),
			Name:          "Test_Process",
			Username:      "User_Name",
			Cmdline:       "/usr/bin/Test_Command",
			NameLower:     "test_process",
			UsernameLower: "user_name",
			CmdlineLower:  "/usr/bin/test_command",
			Cpu:           50.0,
			Memory:        30.0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = state.GetFilteredProcesses()
	}
}
