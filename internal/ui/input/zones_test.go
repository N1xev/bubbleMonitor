package input

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestNewZoneManager(t *testing.T) {
	zm := NewZoneManager()
	if zm == nil {
		t.Fatal("NewZoneManager() returned nil")
	}

	zones := zm.GetZones()
	if len(zones) != 0 {
		t.Errorf("Expected 0 zones, got %d", len(zones))
	}
}

func TestRegister(t *testing.T) {
	zm := NewZoneManager()

	zone := Zone{
		ID:     "test-zone",
		Type:   ZoneTypeButton,
		X:      10,
		Y:      20,
		Width:  100,
		Height: 30,
	}

	zm.Register(zone)

	zones := zm.GetZones()
	if len(zones) != 1 {
		t.Errorf("Expected 1 zone, got %d", len(zones))
	}

	if zones[0].ID != "test-zone" {
		t.Errorf("Expected zone ID 'test-zone', got '%s'", zones[0].ID)
	}
}

func TestUnregister(t *testing.T) {
	t.Skip("Unregister method removed")
}

func TestUnregisterNonExistent(t *testing.T) {
	t.Skip("Unregister method removed")
}

func TestUnregisterHoveredZone(t *testing.T) {
	t.Skip("Unregister method removed")
}

func TestFindZoneAt(t *testing.T) {
	zm := NewZoneManager()

	zone := Zone{
		ID:     "test-zone",
		X:      10,
		Y:      10,
		Width:  50,
		Height: 20,
	}
	zm.Register(zone)

	tests := []struct {
		name     string
		expected string
		x, y     int
	}{
		{"Inside zone", "test-zone", 15, 15},
		{"On left edge", "test-zone", 10, 15},
		{"On top edge", "test-zone", 15, 10},
		{"Left of zone", "", 5, 15},
		{"Above zone", "", 15, 5},
		{"Right of zone", "", 65, 15},
		{"Below zone", "", 15, 35},
		{"On right boundary", "", 60, 15},
		{"On bottom boundary", "", 15, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := zm.FindZoneAt(tt.x, tt.y)
			if tt.expected == "" {
				if found != nil {
					t.Errorf("Expected no zone at (%d, %d), found '%s'", tt.x, tt.y, found.ID)
				}
			} else {
				if found == nil {
					t.Errorf("Expected zone '%s' at (%d, %d), found nil", tt.expected, tt.x, tt.y)
				} else if found.ID != tt.expected {
					t.Errorf("Expected zone '%s' at (%d, %d), found '%s'", tt.expected, tt.x, tt.y, found.ID)
				}
			}
		})
	}
}

func TestFindZoneAtZOrder(t *testing.T) {
	zm := NewZoneManager()

	bottomZone := Zone{ID: "bottom", X: 0, Y: 0, Width: 100, Height: 100}
	topZone := Zone{ID: "top", X: 10, Y: 10, Width: 50, Height: 50}

	zm.Register(bottomZone)
	zm.Register(topZone)

	found := zm.FindZoneAt(20, 20)
	if found == nil || found.ID != "top" {
		t.Error("Should return top-most zone at overlapping position")
	}

	foundOutside := zm.FindZoneAt(80, 80)
	if foundOutside == nil || foundOutside.ID != "bottom" {
		t.Error("Should return bottom zone when outside top zone")
	}
}

func TestIsHovered(t *testing.T) {
	zm := NewZoneManager()

	zone := Zone{ID: "test-zone", X: 0, Y: 0, Width: 10, Height: 10}
	zm.Register(zone)

	if zm.IsHovered("test-zone") {
		t.Error("Zone should not be hovered initially")
	}

	zm.UpdateMousePos(5, 5)

	if !zm.IsHovered("test-zone") {
		t.Error("Zone should be hovered after mouse enters")
	}

	zm.UpdateMousePos(20, 20)

	if zm.IsHovered("test-zone") {
		t.Error("Zone should not be hovered after mouse leaves")
	}
}

func TestGetHoveredZone(t *testing.T) {
	zm := NewZoneManager()

	if zm.GetHoveredZone() != nil {
		t.Error("GetHoveredZone should return nil initially")
	}

	zone := Zone{ID: "hovered", X: 0, Y: 0, Width: 10, Height: 10}
	zm.Register(zone)
	zm.UpdateMousePos(5, 5)

	hovered := zm.GetHoveredZone()
	if hovered == nil || hovered.ID != "hovered" {
		t.Fatalf("Expected hovered zone ID 'hovered', got '%v'", hovered)
	}
}

func TestHoverCallback(t *testing.T) {
	t.Skip("OnHover callback removed")
}

func TestUpdateMousePosDoesNotTriggerClick(t *testing.T) {
	zm := NewZoneManager()

	clickCalled := false

	zone := Zone{
		ID:     "click-zone",
		X:      0,
		Y:      0,
		Width:  10,
		Height: 10,
		OnClick: func() tea.Cmd {
			clickCalled = true
			return nil
		},
	}

	zm.Register(zone)
	zm.UpdateMousePos(5, 5)

	if clickCalled {
		t.Error("OnClick callback should NOT be invoked by UpdateMousePos")
	}
}

func TestClear(t *testing.T) {
	zm := NewZoneManager()

	zm.Register(Zone{ID: "zone-1", X: 0, Y: 0, Width: 10, Height: 10})
	zm.Register(Zone{ID: "zone-2", X: 20, Y: 20, Width: 10, Height: 10})
	zm.UpdateMousePos(5, 5)

	zm.Clear()

	if len(zm.GetZones()) != 0 {
		t.Errorf("Expected 0 zones after Clear, got %d", len(zm.GetZones()))
	}
}

func TestGetZonesReturnsCopy(t *testing.T) {
	zm := NewZoneManager()

	zone := Zone{ID: "original", X: 0, Y: 0, Width: 10, Height: 10}
	zm.Register(zone)

	zones := zm.GetZones()
	zones[0].ID = "modified"

	originalZones := zm.GetZones()
	if originalZones[0].ID != "original" {
		t.Error("GetZones should return a copy, not a reference to internal slice")
	}
}

func TestEmptyManager(t *testing.T) {
	zm := NewZoneManager()

	if zm.FindZoneAt(5, 5) != nil {
		t.Error("FindZoneAt should return nil for empty manager")
	}

	if zm.GetHoveredZone() != nil {
		t.Error("GetHoveredZone should return nil for empty manager")
	}

	if zm.IsHovered("any-id") {
		t.Error("IsHovered should return false for empty manager")
	}
}

func TestZoneTypeString(t *testing.T) {
	tests := []struct {
		expected string
		zoneType ZoneType
	}{
		{"generic", ZoneTypeGeneric},
		{"button", ZoneTypeButton},
		{"tab", ZoneTypeTab},
		{"list-item", ZoneTypeListItem},
		{"menu-item", ZoneTypeMenuItem},
		{"link", ZoneTypeLink},
		{"generic", ZoneType(99)},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.zoneType.String()
			if result != tt.expected {
				t.Errorf("ZoneType.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestZoneWithMetadata(t *testing.T) {
	zm := NewZoneManager()

	zone := Zone{
		ID:     "meta-zone",
		X:      0,
		Y:      0,
		Width:  10,
		Height: 10,
		Metadata: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
	}

	zm.Register(zone)
	zones := zm.GetZones()

	if zones[0].Metadata["key1"] != "value1" {
		t.Error("Metadata should be preserved")
	}

	if zones[0].Metadata["key2"] != 42 {
		t.Error("Metadata values should be preserved")
	}
}

func TestZoneWithParentID(t *testing.T) {
	t.Skip("ParentID field removed")
}
