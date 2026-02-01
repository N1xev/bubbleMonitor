package data

// Health Scoring Constants
const (
	MaxHealthScore = 100

	// Thresholds for resource usage (Percentage)
	HealthThresholdHealthy  = 90.0 // Below this is fine, above starts penalty
	HealthThresholdWarning  = 70.0 // Warning level
	HealthThresholdCritical = 95.0 // Critical level

	// Penalties deducted from health score
	HealthDeductionCPUCritical    = 30
	HealthDeductionCPUHigh        = 10
	HealthDeductionMemoryCritical = 30
	HealthDeductionMemoryHigh     = 10
	HealthDeductionDiskCritical   = 20
	HealthDeductionTempCritical   = 30
	HealthDeductionTempHigh       = 10

	// Process Tracking
	TopProcessesTrackCount = 5
)
