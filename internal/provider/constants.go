package provider

import "time"

// Provider Limits and Timeouts
const (
	// Logging
	MaxLogLines = 1000

	// Timeouts
	SSHTimeout = 2 * time.Second

	// Process Management
	ProcessListCapacity = 500
)
