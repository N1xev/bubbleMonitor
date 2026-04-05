package data

// Layout: reserved rows for header, footer, borders, and chrome.
// Used to calculate available content rows: UI.Height - ReservedContentRows.
const ReservedContentRows = 19

// Settings overlay index boundaries.
// Thresholds (0..ThresholdCount-1), Display (ThresholdCount..+DisplayCount),
// Tabs (+TabCount), Appearance (+AppearanceCount).
const (
	ThresholdCount     = 4 // CPU, Memory, Disk, Temp
	DisplayCount       = 6 // ChartType, ViewType, SortBy, HistoryLen, ProcessCPU, SortDir
	TabCount           = 9 // all available tabs
	AppearanceCount    = 5 // Theme, RefreshRate, BorderType, BorderStyle, Background
	TotalSettingsCount = ThresholdCount + DisplayCount + TabCount + AppearanceCount
)

// AllAvailableTabs is the canonical list of every tab the settings UI can toggle/reorder.
var AllAvailableTabs = []string{
	"Metrics", "Processes", "Disks", "Network",
	"System", "Services", "Connections", "Logs", "Remote",
}

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

	// Kill Dialog
	KillDialogDefaultWidth     = 50
	KillDialogButtonTotalWidth = 16
	KillDialogButtonYOffset    = 4

	// Settings Overlay
	SettingsDefaultWidth  = 120
	SettingsDefaultHeight = 22

	// Context Menu
	ContextMenuDefaultWidth = 30

	// Open Files Overlay
	OpenFilesDefaultWidth  = 80
	OpenFilesDefaultHeight = 20

	SamLabDefaultWidth  = 50
	SamLabDefaultHeight = 14
)
