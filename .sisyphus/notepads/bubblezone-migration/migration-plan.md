# BubbleZone v2 Migration Plan for bubbleMonitor

## Executive Summary

This document outlines a comprehensive migration plan to replace the current manual zone tracking system with bubblezone v2 library in the bubbleMonitor TUI application. The current implementation uses a custom `ZoneTracker` struct with manual coordinate calculations, while bubblezone v2 provides element-linked zones with automatic position calculation and `InBounds()` detection.

## Current Implementation Analysis

### ZoneTracker Structure
```go
// Current implementation in internal/ui/zones.go
type ZoneTracker struct {
    zones   []Zone
    cursorX int
    cursorY int
    enabled bool
}

// Zone definition in internal/data/types.go
type Zone struct {
    ID     string
    X      int
    Y      int
    Width  int
    Height int
}
```

### Current Zone Usage Patterns

#### 1. Tab Navigation (layout.go:192)
**Current:**
```go
zt.Mark("tab-"+titleRaw, tabWidth, 1)
```

**BubbleZone v2:**
```go
zone.Mark("tab-"+titleRaw, tabWidth, 1)
```

#### 2. Footer Buttons (layout.go:306, 316, 327, etc.)
**Current:**
```go
zt.Mark("footer-help", helpWidth, 1)
zt.Mark("footer-settings", settingsWidth, 1)
zt.Mark("footer-quit", quitWidth, 1)
```

**BubbleZone v2:**
```go
zone.Mark("footer-help", helpWidth, 1)
zone.Mark("footer-settings", settingsWidth, 1)
zone.Mark("footer-quit", quitWidth, 1)
```

#### 3. Process Rows (layout.go:490)
**Current:**
```go
zt.MarkAt("process-"+fmt.Sprint(proc.Pid), 0, rowY, s.Width, 1)
```

**BubbleZone v2:**
```go
zone.Mark("process-"+fmt.Sprint(proc.Pid), s.Width, 1)
```

#### 4. System Blocks (layout.go:526)
**Current:**
```go
zt.MarkAt("systemblock-"+fmt.Sprint(i), blockX, blockY, blockWidth, blockHeight)
```

**BubbleZone v2:**
```go
zone.Mark("systemblock-"+fmt.Sprint(i), blockWidth, blockHeight)
```

#### 5. Dialog Buttons (layout.go:622, 623)
**Current:**
```go
zt.MarkAt("kill-yes", dialogX+3+(dialogWidth-6-16)/2, dialogY+dialogHeight-3, 7, 1)
zt.MarkAt("kill-no", dialogX+3+(dialogWidth-6-16)/2+10, dialogY+dialogHeight-3, 6, 1)
```

**BubbleZone v2:**
```go
zone.Mark("kill-yes", 7, 1)
zone.Mark("kill-no", 6, 1)
```

#### 6. Context Menu Options (layout.go:727)
**Current:**
```go
zt.MarkAt("context-menu-"+fmt.Sprint(i), menuX, optY, menuWidth, 1)
```

**BubbleZone v2:**
```go
zone.Mark("context-menu-"+fmt.Sprint(i), menuWidth, 1)
```

#### 7. Mouse Event Detection (mouse.go:136, 200, 275)
**Current:**
```go
zone := m.FindZoneAt(x, y)
```

**BubbleZone v2:**
```go
if zone.Get("zone-id").InBounds(msg) {
    // Handle click
}
```

## Migration Strategy

### Phase 1: Setup and Dependencies

#### 1.1 Add bubblezone v2 Dependency
```bash
go get github.com/lrstanley/bubblezone/v2@latest
```

#### 1.2 Initialize Global Zone Manager
**Current:**
```go
var globalZones = NewZoneTracker()
func GetGlobalZones() *ZoneTracker {
    return globalZones
}
```

**BubbleZone v2:**
```go
import "github.com/lrstanley/bubblezone/v2"

func init() {
    zone.NewGlobal()
}
```

### Phase 2: Replace ZoneTracker with bubblezone

#### 2.1 Remove Current ZoneTracker Implementation
**Files to modify:**
- `internal/ui/zones.go` - Remove entire ZoneTracker struct and methods
- `internal/data/types.go` - Remove Zone struct definition

#### 2.2 Update Layout View Function
**Current:**
```go
zt := GetGlobalZones()
zt.Reset()
zt.SetCursor(0, 0)
// ... various Mark/MarkAt calls
```

**BubbleZone v2:**
```go
// No need to reset or set cursor - bubblezone handles this automatically
// ... various zone.Mark calls
```

#### 2.3 Update Mouse Event Handling
**Current:**
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

**BubbleZone v2:**
```go
// Remove FindZoneAt method entirely
// Use zone.Get("zone-id").InBounds(msg) directly in event handlers
```

### Phase 3: Update Event Handlers

#### 3.1 Tab Click Handling
**Current:**
```go
if strings.HasPrefix(zone.ID, "tab-") {
    tabName := strings.TrimPrefix(zone.ID, "tab-")
    for i, tab := range m.ActiveTabs {
        if tab == tabName {
            m.SelectedTab = i
            break
        }
    }
}
```

**BubbleZone v2:**
```go
if zone.Get("tab-"+tabName).InBounds(msg) {
    for i, tab := range m.ActiveTabs {
        if tab == tabName {
            m.SelectedTab = i
            break
        }
    }
}
```

#### 3.2 Footer Button Handling
**Current:**
```go
switch zone.ID {
case "footer-help":
    m.ShowHelp = !m.ShowHelp
case "footer-settings":
    m.ShowSettings = !m.ShowSettings
// ... other cases
}
```

**BubbleZone v2:**
```go
if zone.Get("footer-help").InBounds(msg) {
    m.ShowHelp = !m.ShowHelp
} else if zone.Get("footer-settings").InBounds(msg) {
    m.ShowSettings = !m.ShowSettings
// ... other conditions
}
```

#### 3.3 Process Row Selection
**Current:**
```go
if strings.HasPrefix(zone.ID, "process-") {
    pidStr := strings.TrimPrefix(zone.ID, "process-")
    var pid int32
    fmt.Sscanf(pidStr, "%d", &pid)
    // Find process and select
}
```

**BubbleZone v2:**
```go
if zone.Get("process-"+pidStr).InBounds(msg) {
    // Process is already selected via mouse coordinates
    // No need to parse PID from ID
}
```

### Phase 4: Update View Functions

#### 4.1 Wrap Views with zone.Scan
**Current:**
```go
ui := lipgloss.JoinVertical(lipgloss.Left,
    topBar,
    topPad,
    content,
)
```

**BubbleZone v2:**
```go
ui := zone.Scan(lipgloss.JoinVertical(lipgloss.Left,
    topBar,
    topPad,
    content,
))
```

#### 4.2 Remove Manual Cursor Tracking
**Current:**
```go
zt.SetCursor(tabsStartX, currentRowY)
zt.TrackSpaces(0)
zt.TrackText(tabText)
```

**BubbleZone v2:**
```go
// No manual cursor tracking needed - bubblezone handles positioning automatically
```

## Performance Considerations

### 1. Zone Registration Overhead
- **Current:** Manual coordinate calculations on every render
- **BubbleZone v2:** Automatic zone registration with ANSI sequence injection
- **Impact:** Minimal overhead, optimized for performance

### 2. Memory Usage
- **Current:** Stores all zones in memory with manual tracking
- **BubbleZone v2:** Uses efficient data structures with automatic cleanup
- **Impact:** Similar or slightly better memory usage

### 3. Render Performance
- **Current:** Additional cursor tracking operations
- **BubbleZone v2:** Streamlined zone registration
- **Impact:** Potential slight improvement in render performance

## Testing Strategy

### 1. Unit Tests
- Test zone registration and detection
- Verify mouse event handling
- Test edge cases (overlapping zones, boundary conditions)

### 2. Integration Tests
- Test complete UI interactions
- Verify tab navigation
- Test footer button functionality
- Test process selection and context menus

### 3. Performance Tests
- Measure render performance before/after migration
- Test memory usage
- Verify responsiveness with large numbers of zones

### 4. Visual Regression Tests
- Ensure UI appearance remains consistent
- Verify no visual artifacts from zone markers

## Dependencies and Prerequisites

### Required Dependencies
```go
go get github.com/lrstanley/bubblezone/v2@latest
go get charm.land/bubbletea/v2@latest
go get charm.land/lipgloss/v2@latest
```

### Version Compatibility
- Go 1.23.0 or later
- bubbletea v2.0.0 or later
- lipgloss v2.0.0 or later

## Risk Assessment

### Low Risk
- Zone registration API is similar
- Mouse event handling patterns are compatible
- No breaking changes to UI appearance

### Medium Risk
- Requires updating all event handlers
- Need to ensure proper zone cleanup
- Potential performance impacts during transition

### High Risk
- Breaking changes if dependencies are not updated
- Complex interactions may require careful testing

## Migration Timeline

### Week 1: Preparation
- Add bubblezone v2 dependency
- Set up testing infrastructure
- Create migration plan documentation

### Week 2: Core Migration
- Replace ZoneTracker implementation
- Update layout view function
- Migrate basic zone usage patterns

### Week 3: Event Handler Migration
- Update mouse event handling
- Migrate tab navigation
- Update footer button interactions

### Week 4: Testing and Optimization
- Comprehensive testing
- Performance optimization
- Documentation updates

## Before/After Code Examples

### Example 1: Tab Navigation

**Before:**
```go
// layout.go
zt.SetCursor(tabsStartX, currentRowY)
zt.TrackSpaces(0)

for i, titleRaw := range s.ActiveTabs {
    title := strings.ToUpper(titleRaw)
    tabText := " " + title + " "
    tabWidth := lipgloss.Width(tabText)
    
    if zt.cursorX-tabsStartX+tabWidth >= availableWidth && zt.cursorX-tabsStartX > 0 {
        zt.TrackNewline()
        zt.SetCursor(tabsStartX, zt.cursorY)
    }
    
    zt.Mark("tab-"+titleRaw, tabWidth, 1)
    zt.TrackText(tabText)
}
```

**After:**
```go
// layout.go
for i, titleRaw := range s.ActiveTabs {
    title := strings.ToUpper(titleRaw)
    tabText := " " + title + " "
    tabWidth := lipgloss.Width(tabText)
    
    zone.Mark("tab-"+titleRaw, tabWidth, 1)
}
```

### Example 2: Process Row Selection

**Before:**
```go
// layout.go
for i := startIdx; i < endIdx; i++ {
    proc := visibleProcs[i]
    rowY := listStartY + 2 + (i - startIdx)
    zt.MarkAt("process-"+fmt.Sprint(proc.Pid), 0, rowY, s.Width, 1)
}

// mouse.go
if strings.HasPrefix(zone.ID, "process-") {
    pidStr := strings.TrimPrefix(zone.ID, "process-")
    var pid int32
    fmt.Sscanf(pidStr, "%d", &pid)
    
    procs := m.GetFilteredProcesses()
    for i, p := range procs {
        if p.Pid == pid {
            m.SelectedProcess = i
            return nil
        }
    }
}
```

**After:**
```go
// layout.go
for i := startIdx; i < endIdx; i++ {
    proc := visibleProcs[i]
    zone.Mark("process-"+fmt.Sprint(proc.Pid), s.Width, 1)
}

// mouse.go
if zone.Get("process-"+pidStr).InBounds(msg) {
    // Process is already selected via mouse coordinates
    return nil
}
```

### Example 3: Footer Button Handling

**Before:**
```go
// layout.go
zt.Mark("footer-help", helpWidth, 1)
zt.Mark("footer-settings", settingsWidth, 1)
zt.Mark("footer-quit", quitWidth, 1)

// mouse.go
switch zone.ID {
case "footer-help":
    m.ShowHelp = !m.ShowHelp
case "footer-settings":
    m.ShowSettings = !m.ShowSettings
case "footer-quit":
    return tea.Quit
}
```

**After:**
```go
// layout.go
zone.Mark("footer-help", helpWidth, 1)
zone.Mark("footer-settings", settingsWidth, 1)
zone.Mark("footer-quit", quitWidth, 1)

// mouse.go
if zone.Get("footer-help").InBounds(msg) {
    m.ShowHelp = !m.ShowHelp
} else if zone.Get("footer-settings").InBounds(msg) {
    m.ShowSettings = !m.ShowSettings
} else if zone.Get("footer-quit").InBounds(msg) {
    return tea.Quit
}
```

## Conclusion

This migration plan provides a comprehensive approach to replacing the current manual zone tracking system with bubblezone v2. The migration offers several benefits:

1. **Simplified Code:** Eliminates manual coordinate calculations and cursor tracking
2. **Better Performance:** Optimized zone registration and detection
3. **Improved Maintainability:** Cleaner API and reduced boilerplate
4. **Enhanced Features:** Automatic zone management and better integration with BubbleTea

The migration is well-structured with clear phases, comprehensive testing strategy, and minimal risk to existing functionality. With proper execution, this migration will significantly improve the codebase quality and maintainability of bubbleMonitor.