package input

import tea "charm.land/bubbletea/v2"

// ZoneType categorizes different types of interactive zones
type ZoneType int

const (
	ZoneTypeGeneric ZoneType = iota
	ZoneTypeButton
	ZoneTypeTab
	ZoneTypeListItem
	ZoneTypeMenuItem
	ZoneTypeLink
)

// String returns the string representation of a ZoneType
func (z ZoneType) String() string {
	switch z {
	case ZoneTypeButton:
		return "button"
	case ZoneTypeTab:
		return "tab"
	case ZoneTypeListItem:
		return "list-item"
	case ZoneTypeMenuItem:
		return "menu-item"
	case ZoneTypeLink:
		return "link"
	default:
		return "generic"
	}
}

// Zone represents an interactive clickable/hoverable area in the UI
type Zone struct {
	Metadata map[string]interface{}
	OnClick  func() tea.Cmd
	ID       string
	Type     ZoneType
	X        int
	Y        int
	Width    int
	Height   int
}

// ZoneManager defines the interface for managing interactive zones
type ZoneManager interface {
	Register(zone Zone)
	FindZoneAt(x, y int) *Zone
	GetHoveredZone() *Zone
	IsHovered(zoneID string) bool
	UpdateMousePos(x, y int) tea.Cmd
	Clear()
	GetZones() []Zone
}

// zoneManager implements the ZoneManager interface
type zoneManager struct {
	hoveredZoneID string
	zones         []Zone
	lastMouseX    int
	lastMouseY    int
}

// Ensure zoneManager implements ZoneManager interface
var _ ZoneManager = (*zoneManager)(nil)
