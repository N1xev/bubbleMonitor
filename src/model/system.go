package model

// DiskPartition holds information about a disk partition
type DiskPartition struct {
	Mountpoint string
	Device     string
	Fstype     string
	Total      uint64
	Used       uint64
	UsedPct    float64
}

// GpuInfo holds information about a GPU
type GpuInfo struct {
	Name        string
	Driver      string
	MemoryTotal string
	MemoryUsed  string
}
