package input

import tea "charm.land/bubbletea/v2"

// NewZoneManager creates a new ZoneManager instance
func NewZoneManager() ZoneManager {
	return &zoneManager{
		zones:         []Zone{},
		hoveredZoneID: "",
		lastMouseX:    -1,
		lastMouseY:    -1,
	}
}

// Register adds a new zone to the manager
func (zm *zoneManager) Register(zone Zone) {
	zm.zones = append(zm.zones, zone)
}

// FindZoneAt returns the top-most zone at the given coordinates
func (zm *zoneManager) FindZoneAt(x, y int) *Zone {
	for i := len(zm.zones) - 1; i >= 0; i-- {
		zone := zm.zones[i]
		if x >= zone.X && x < zone.X+zone.Width && y >= zone.Y && y < zone.Y+zone.Height {
			return &zone
		}
	}
	return nil
}

// GetHoveredZone returns the currently hovered zone
func (zm *zoneManager) GetHoveredZone() *Zone {
	if zm.hoveredZoneID == "" {
		return nil
	}
	for i := range zm.zones {
		if zm.zones[i].ID == zm.hoveredZoneID {
			return &zm.zones[i]
		}
	}
	return nil
}

// IsHovered returns true if the zone with the given ID is currently hovered
func (zm *zoneManager) IsHovered(zoneID string) bool {
	if zm.hoveredZoneID == "" {
		return false
	}
	return zm.hoveredZoneID == zoneID
}

// UpdateMousePos updates the mouse position and stores the hovered zone ID
func (zm *zoneManager) UpdateMousePos(x, y int) tea.Cmd {
	zm.lastMouseX = x
	zm.lastMouseY = y

	newHovered := zm.FindZoneAt(x, y)
	var newHoveredID string
	if newHovered != nil {
		newHoveredID = newHovered.ID
	}

	if newHoveredID != zm.hoveredZoneID {
		zm.hoveredZoneID = newHoveredID
	}

	return nil
}

// Clear removes all zones from the manager
func (zm *zoneManager) Clear() {
	zm.zones = []Zone{}
}

// GetZones returns a copy of all registered zones
func (zm *zoneManager) GetZones() []Zone {
	result := make([]Zone, len(zm.zones))
	copy(result, zm.zones)
	return result
}
