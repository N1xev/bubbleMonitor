package provider

import "time"

// Provider Limits and Timeouts
const (
	// Logging
	MaxLogLines = 50

	// Timeouts
	SSHTimeoutSeconds = 2 * time.Second

	// Process Management
	ProcessListCapacity      = 500
	InternerCleanupFrequency = 1000

	// Network
	NetworkBaseRateMBps = 10.0

	// Protocols
	ProtoTCP = 1
	ProtoUDP = 2
)
