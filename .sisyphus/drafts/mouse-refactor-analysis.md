# Analysis: Mouse/Click/Hover Refactoring - Single Source of Truth

## Current Architecture Overview

### Entry Point Flow
```
User Input → tea.MouseMsg
            ↓
     internal/app/update.go:43
            ↓
   m.HandleMouse(msg)
            ↓
     internal/app/mouse.go
            ↓
    Zone lookup OR coordinate math
            ↓
    Action execution
```

### Files Involved in Mouse Handling

| File | Responsibility | Lines of Code |
|------|----------------|---------------|
| `internal/app/mouse.go` | Primary mouse handler | 321 |
| `internal/app/update.go` | Entry point routing | 1252 |
| `internal/ui/zones.go` | Zone tracking infrastructure | 105 |
| `internal/ui/layout.go` | Zone registration during render | 761 |
| `internal/ui/tabs/processes.go` | Process row hover | ~200 |
| `internal/ui/overlays/kill.go` | Kill dialog hover | ~50 |

---

## Problems Identified

### 1. DUPLICATE ZONE DETECTION LOGIC
**Location:** Two identical `FindZoneAt` implementations

**File 1:** `internal/app/mouse.go:136-144`
```go
func (m *Model) FindZoneAt(x, y int) *data.Zone {
    for i := len(m.Zones) - 1; i >= 0; i-- {
        zone := m.Zones[i]
        if x >= zone.X && x < zone.X+zone.Width && y >= zone.Y && y < zone.Y+zone.Height {
            return &zone
        }
    }
    return nil
}
```

**File 2:** `internal/ui/zones.go:73-81`
```go
func (z *ZoneTracker) FindZoneAt(x, y int) *Zone {
    for i := len(z.zones) - 1; i >= 0; i-- {
        zone := z.zones[i]
        if x >= zone.X && x < zone.X+zone.Width && y >= zone.Y && y < zone.Y+zone.Height {
            return &zone
        }
    }
    return nil
}
```

**Impact:** Maintenance burden, risk of divergent implementations

### 2. INCONSISTENT HOVER DETECTION (5 Different Patterns!)

**Pattern A: Tab Hover** (`layout.go:192`)
```go
isTabHover := s.MouseY == currentRowY && s.MouseX >= zt.cursorX && 
              s.MouseX < zt.cursorX+lipgloss.Width(renderedTab)
```

**Pattern B: Process Row Hover** (`tabs/processes.go:189`)
```go
rowY := listStartY + 2 + (i - startIdx)
isHovered := !isSelected && mouseY == rowY
```

**Pattern C: Kill Dialog Button Hover** (`overlays/kill.go:35-36`)
```go
isHoverYes := s.MouseY == buttonY && s.MouseX >= buttonsStartX && s.MouseX < buttonsStartX+7
isHoverNo := s.MouseY == buttonY && s.MouseX >= buttonsStartX+10 && s.MouseX < buttonsStartX+16
```

**Pattern D: Footer Button Hover** (`layout.go:428-429`)
```go
isFooterHover := s.MouseY == footerY && s.MouseX >= zt.cursorX && s.MouseX < zt.cursorX+info.width
```

**Pattern E: Context Menu Option Hover** (`layout.go:696`)
```go
isHovered := s.MouseY == optY && s.MouseX >= menuX && s.MouseX < menuX+menuWidth
```

**Impact:** No single way to check hover state, duplicated coordinate math everywhere

### 3. MIXED CLICK HANDLING APPROACHES

**Zone-Based (Good):**
- Tabs: `tab-{name}`
- Footer buttons: `footer-help`, `footer-settings`, etc.
- SamLab link

**Coordinate Math (Bad):**
- Right-click on process rows (`mouse.go:115-127`)
- Kill dialog buttons (no zones)
- Context menu items (no zones)
- Settings overlay controls (no zones)

**Impact:** Some clicks use zones, others use hardcoded coordinates. Inconsistent behavior.

### 4. ZONES ONLY REGISTERED DURING RENDER

Zones are created in `layout.go` during the render cycle:
```go
zt.Mark("tab-"+titleRaw, tabWidth, 1)  // Line 207
zt.Mark("footer-help", width, 1)       // Line 441
zt.Mark("samlab", width, 1)            // Line 385
```

**Impact:** 
- Can't query zone info outside render
- Zones don't exist for hit-testing in overlays
- Race condition risk with global singleton

### 5. GLOBAL SINGLETON ZONE TRACKER

```go
// internal/ui/zones.go:101-104
var globalZones = NewZoneTracker()

func GetGlobalZones() *ZoneTracker {
    return globalZones
}
```

**Impact:** Potential race conditions, hard to test, unclear ownership

### 6. NO ZONE REGISTRATION FOR DYNAMIC ELEMENTS

Missing zones for:
- Process list rows
- Kill dialog Yes/No buttons
- Context menu items
- Settings controls
- Any overlay interactive elements

---

## Target Architecture: Single Source of Truth

### Design Principles

1. **Unified ZoneManager** - One system handles all hit-testing
2. **Event Delegation** - Zones can have callbacks attached
3. **Separate Hit-Testing from Rendering** - Zones exist independently of render
4. **Consistent Hover State** - Single way to check if cursor is over a zone
5. **Hierarchical Zones** - Parent/child relationships for complex UI

### Proposed New Types

```go
// internal/ui/input/zones.go

type ZoneType int

const (
    ZoneTypeGeneric ZoneType = iota
    ZoneTypeButton
    ZoneTypeTab
    ZoneTypeListItem
    ZoneTypeMenuItem
)

type Zone struct {
    ID       string
    X, Y     int
    Width    int
    Height   int
    Type     ZoneType
    ParentID string          // For hierarchical relationships
    Metadata map[string]interface{} // Process ID, etc.
    
    // Optional callbacks for event delegation
    OnClick  func() tea.Cmd
    OnHover  func(hovered bool)
}

type ZoneManager struct {
    zones       []Zone
    hoveredZone *Zone      // Currently hovered zone
    lastMouseX  int
    lastMouseY  int
}

func (zm *ZoneManager) Register(zone Zone)
func (zm *ZoneManager) Unregister(id string)
func (zm *zm *ZoneManager) FindZoneAt(x, y int) *Zone
func (zm *ZoneManager) GetHoveredZone() *Zone
func (zm *ZoneManager) IsHovered(zoneID string) bool
func (zm *ZoneManager) UpdateMousePos(x, y int) tea.Cmd // Returns commands from callbacks
func (zm *ZoneManager) Clear()
```

---

## Refactoring Strategy

### Phase 1: Foundation (Core Zone System)
- Create new `internal/ui/input/zones.go` with ZoneManager
- Move Zone struct from data/state.go to input package
- Add comprehensive tests for ZoneManager
- Deprecate ZoneTracker gradually

### Phase 2: Migrate Existing Zone Usage
- Replace ZoneTracker with ZoneManager in layout.go
- Update tab zone registration
- Update footer button zone registration
- Update SamLab zone registration

### Phase 3: Consolidate Click Handling
- Move all click logic from mouse.go to zone callbacks
- Make handleLeftClick delegate to ZoneManager
- Remove coordinate-based click detection from handleRightClick

### Phase 4: Add Missing Zones
- Register process rows as zones
- Register kill dialog buttons as zones
- Register context menu items as zones
- Register settings controls as zones

### Phase 5: Unify Hover Detection
- Replace all 5 hover patterns with ZoneManager.IsHovered()
- Add hover state tracking to ZoneManager
- Update process row rendering
- Update kill dialog rendering
- Update context menu rendering

### Phase 6: Cleanup
- Remove deprecated ZoneTracker
- Remove duplicate FindZoneAt from mouse.go
- Remove coordinate-based detection from mouse.go
- Update documentation

---

## Risk Assessment

| Risk | Severity | Mitigation |
|------|----------|------------|
| Breaking existing mouse functionality | High | Comprehensive tests, feature flags |
| Performance regression | Medium | Benchmark hit-testing, optimize if needed |
| Zone registration overhead | Low | Measure render time, lazy registration |
| Regression in hover effects | Medium | Visual QA for all hover states |

---

## Success Metrics

- ✅ Zero duplicate FindZoneAt implementations
- ✅ Single pattern for hover detection (ZoneManager.IsHovered)
- ✅ All clickable elements registered as zones
- ✅ No coordinate-based click detection outside ZoneManager
- ✅ All existing tests pass
- ✅ Performance maintained (±10% render time)

---

## Files to Modify (Estimated)

| File | Changes | Risk |
|------|---------|------|
| `internal/ui/input/zones.go` | NEW FILE - Core zone system | Low |
| `internal/data/state.go` | Remove Zone, update Zones field | Medium |
| `internal/ui/zones.go` | DEPRECATE - Keep for backward compat | Low |
| `internal/app/mouse.go` | Refactor to use ZoneManager | High |
| `internal/ui/layout.go` | Update zone registration | Medium |
| `internal/ui/tabs/processes.go` | Add zone registration | Medium |
| `internal/ui/overlays/kill.go` | Add zone registration | Low |
| `internal/ui/overlays/settings.go` | Add zone registration | Low |
| `internal/ui/overlays/files.go` | Add zone registration | Low |
| `internal/ui/overlays/help.go` | Add zone registration | Low |
| `internal/ui/overlays/samlab.go` | Add zone registration | Low |

Total: ~10 files, estimated 400-600 lines of changes
