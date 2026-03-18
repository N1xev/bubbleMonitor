package testutil

type MockRemoteProvider struct {
	UptimeFunc      func(host, user, password string) (string, error)
	ProcessListFunc func(host, user, password string) (string, error)
	DiskUsageFunc   func(host, user, password string) (string, error)
	MemoryUsageFunc func(host, user, password string) (string, error)
}

func (m *MockRemoteProvider) Uptime(host, user, password string) (string, error) {
	if m.UptimeFunc != nil {
		return m.UptimeFunc(host, user, password)
	}
	return "uptime: 5 days", nil
}

func (m *MockRemoteProvider) ProcessList(host, user, password string) (string, error) {
	if m.ProcessListFunc != nil {
		return m.ProcessListFunc(host, user, password)
	}
	return "PID COMMAND\n1 init\n2 systemd", nil
}

func (m *MockRemoteProvider) DiskUsage(host, user, password string) (string, error) {
	if m.DiskUsageFunc != nil {
		return m.DiskUsageFunc(host, user, password)
	}
	return "Filesystem      Size  Used Avail Use% Mounted on\n/dev/sda1       100G   50G   50G  50% /", nil
}

func (m *MockRemoteProvider) MemoryUsage(host, user, password string) (string, error) {
	if m.MemoryUsageFunc != nil {
		return m.MemoryUsageFunc(host, user, password)
	}
	return "Mem:         16384     8192     8192  50%\nSwap:         8192      1024     7168  12%", nil
}
