# Fix Zones and Scroll - Learnings

## Task 1: Fix Zone Misalignment in layout.go (line 437)

### Change Made
- **File**: `internal/ui/layout.go`
- **Line 437**: Changed `listStartY := numTabRows + 4` → `listStartY := topBarH + topGap`

### Technical Details
- **Root Cause**: `numTabRows` is just a count (1, 2, 3...), not actual rendered height
- **Solution**: Use `topBarH` (actual rendered height from lipgloss.Height()) + `topGap` (which equals 1)
- **lipgloss.Height()**: Returns the actual physical height of rendered text block (includes wrapping, borders, padding)

### Code Context
```go
topGap := 1
topBarH := lipgloss.Height(topBar)  // Actual rendered height
// ...
listStartY := topBarH + topGap  // Correct: uses actual height
```

### Notes
- Zone coordinates are absolute (0,0 = top-left)
- Mouse events provide absolute coordinates
- listStartY should equal actual rendered height, not count
- The +2 offset (mentioned in inherited wisdom) accounts for: header line + blank line before first process
- Also removed unused `numTabRows` variable declaration (line 436) since it was only used for the incorrect calculation

### Verification
- Build passes with exit code 0
- Process zones should now align with actual rendered positions

---

## Task 2: Fix Zone Misalignment in processes.go (line 188)

### Change Made
- **File**: `internal/ui/tabs/processes.go`
- **Line 188**: Changed `rowY := listStartY + (i - startIdx)` → `rowY := listStartY + 2 + (i - startIdx)`

### Technical Details
- **Root Cause**: Hover detection in processes.go was not accounting for header offset
- **Solution**: Add +2 to rowY calculation to match layout.go zone creation
- **+2 Offset Explanation**: Accounts for the header line (title) + blank line before first process row

### Code Context
```go
for i := startIdx; i < endIdx; i++ {
    proc := filtered[i]
    isSelected := i == s.SelectedProcess

    rowY := listStartY + 2 + (i - startIdx)  // Added +2 for header offset
    isHovered := !isSelected && mouseY == rowY
    // ...
}
```

### Notes
- Zone coordinates are absolute (0,0 = top-left)
- Mouse events provide absolute coordinates
- layout.go adds +2 for header offset when creating zones
- processes.go must match this offset for hover detection
- This ensures hover highlights align with click zones

### Verification
- Build passes with exit code 0
- Hover highlights now align with click zones

---

## Task 3: Fix pgup/pgdown to Require Focus and Scroll Only Active Block

### Changes Made
- **File**: `internal/app/update.go`
- **pgup (System tab)**: Lines 466-477 - Added focus check, scroll only active block
- **pgdown (System tab)**: Lines 533-541 - Added focus check, scroll only active block

### Technical Details
- **Root Cause**: pgup/pgdown were scrolling ALL system blocks simultaneously, regardless of focus
- **Solution**: 
  - Changed condition from `m.SystemBlockCount > 0` to `m.ActiveScrollBlock >= 0 && m.SystemBlockScrollable[m.ActiveScrollBlock]`
  - Removed for loop iterating over all blocks
  - Now only scrolls the focused block using `m.ActiveScrollBlock` index

### Code Context
```go
// pgup - Before (scrolled all blocks)
} else if currentTab == "System" && m.SystemBlockCount > 0 {
    rows := m.Height - 19
    for i := 0; i < m.SystemBlockCount; i++ {
        m.SystemBlockScrollOffsets[i] -= rows
    }
}

// pgup - After (scrolls only active block)
} else if currentTab == "System" && m.ActiveScrollBlock >= 0 && m.SystemBlockScrollable[m.ActiveScrollBlock] {
    rows := m.Height - 19
    m.SystemBlockScrollOffsets[m.ActiveScrollBlock] -= rows
    if m.SystemBlockScrollOffsets[m.ActiveScrollBlock] < 0 {
        m.SystemBlockScrollOffsets[m.ActiveScrollBlock] = 0
    }
}
```

### Notes
- **Focus Indicator**: Border color change from `b` (default) to `p` (focused) via `[]` key
- **ActiveScrollBlock**: Index of currently focused scrollable block (-1 if none)
- **SystemBlockScrollable**: Boolean slice indicating which blocks are scrollable
- pgup/pgdown now only work AFTER pressing `[]` to focus a block
- This matches expected UX: scroll only the block you're focused on

### Verification
- Build passes with exit code 0
- pgup/pgdown now require block focus to work

---

## Task 4: Initialize ActiveScrollBlock to -1 (No Initial Focus)

### Change Made
- **File**: `internal/app/model.go`
- **Line 134**: Changed `ActiveScrollBlock: 0` → `ActiveScrollBlock: -1`

### Technical Details
- **Root Cause**: ActiveScrollBlock was initialized to 0, which implied first block was focused by default
- **Solution**: Initialize to -1 to indicate no block is focused initially
- **Rationale**: User must explicitly press `[]` to focus a block before arrow keys/scroll work

### Code Context
```go
// Before
ActiveScrollBlock: 0,

// After
ActiveScrollBlock: -1,
```

### Notes
- **Focus Indicator**: Border color change from `b` (default) to `p` (focused) via `[]` key
- **-1 Value**: Indicates "no block focused" state
- **Arrow Keys**: Only work after explicit focus via `[]` key
- **pgup/pgdown**: Already fixed in Task 3 to require focus

### Verification
- Build passes with exit code 0
- System tab starts with no block focused
- User must press `[]` to focus a block first

---

## Task 5: Replace Hardcoded Zone Creation with Dynamic Zone System

### Changes Made
- **File**: `internal/ui/layout.go`
- **Tabs (lines 174-206)**: Replaced MarkAt with Mark using cursor tracking
- **Footer buttons (lines 291-386)**: Replaced MarkAt with Mark using cursor tracking

### Technical Details
- **Root Cause**: Zone positions were calculated separately from text rendering using MarkAt(), which could lead to misalignment
- **Solution**: Use dynamic Mark() with cursor tracking (SetCursor, Mark, TrackText) to ensure zones track exactly where text renders

### Zone Tracker Methods
- `MarkAt(id, x, y, w, h)` - Uses explicit coordinates (potentially misaligned)
- `Mark(id, w, h)` - Uses cursor position (dynamically tracked)
- `SetCursor(x, y)` - Sets cursor position
- `TrackText(text)` - Advances cursor by text width
- `TrackSpaces(count)` - Advances cursor by space count
- `TrackNewline()` - Moves cursor to next line

### Code Context - Tabs
```go
// Before: Used MarkAt with calculated positions
tabsStartX := headerWidth + spacerWidth
currentRowY := 1
currentRowWidth := 0
for i, titleRaw := range s.ActiveTabs {
    // ...
    zt.MarkAt("tab-"+titleRaw, tabsStartX+currentRowWidth, currentRowY, tabWidth, 1)
    // ...
    currentRowWidth += tabWidth
}

// After: Use dynamic Mark with cursor tracking
zt.SetCursor(tabsStartX, currentRowY)
for i, titleRaw := range s.ActiveTabs {
    // ...
    zt.Mark("tab-"+titleRaw, tabWidth, 1)  // Uses cursor position
    // ...
    zt.TrackText(tabText)  // Advance cursor after rendering
}
```

### Code Context - Footer
```go
// Before: Used MarkAt with buttonStart variable
buttonStart := 0
zt.MarkAt("footer-help", buttonStart, footerY, helpWidth, 1)
if isFooterHover && s.MouseX >= buttonStart && s.MouseX < buttonStart+helpWidth {
    // ...
}
buttonStart += helpWidth

// After: Use dynamic Mark with cursor tracking  
zt.SetCursor(0, footerY)
zt.Mark("footer-help", helpWidth, 1)  // Uses cursor position
if isFooterHover && s.MouseX >= zt.cursorX && s.MouseX < zt.cursorX+helpWidth {
    // ...
}
zt.TrackText(helpText)  // Advance cursor after rendering
```

### Notes
- Zone coordinates are absolute (0,0 = top-left)
- Mouse events provide absolute coordinates
- Dynamic tracking ensures zones align with actual rendered button/tab positions
- Hover detection updated to use `zt.cursorX` instead of manual buttonStart tracking

### Verification
- Build passes with exit code 0
- Tab zones and footer button zones now use dynamic cursor tracking
- Hover detection uses cursor position for accurate alignment

---

## Task 6: Verify Top Bar Zone Alignment (Task 1 of fix-top-bar-and-scroll plan)

### Verification Performed
- **File**: `internal/ui/layout.go`
- Searched for all `MarkAt()` calls using ast-grep and grep
- Verified top bar zones (tabs and footer buttons) use dynamic system

### Findings

#### Top Bar Tabs (lines 174-210)
- Line 178: `zt.SetCursor(tabsStartX, currentRowY)` - Sets cursor at tab start
- Line 192: `zt.Mark("tab-"+titleRaw, tabWidth, 1)` - Uses dynamic Mark
- Line 194: Hover detection uses `zt.cursorX` and `zt.cursorY` - Correct!
- Line 209: `zt.TrackText(tabText)` - Advances cursor after rendering

#### Footer Buttons (lines 291-387)
- Line 293: `zt.SetCursor(0, footerY)` - Sets cursor at footer start
- Lines 306, 316, 327, 337, 350, 360, 370, 380: All use `zt.Mark()` - Dynamic!
- Hover detection uses `zt.cursorX` - Correct!
- All buttons use `zt.TrackText()` to advance cursor

#### Remaining MarkAt Calls (NOT top bar)
- Line 490: Process rows (not top bar)
- Lines 593-594: Kill dialog buttons (overlay)
- Line 698: Context menu (overlay)

These are NOT top bar zones and don't need to be changed.

### Verification
- Build passes with exit code 0
- Top bar zones already use dynamic zone tracking system
- No hardcoded zone coordinates in tabs or footer buttons

---

## Task 7: Fix Mouse Scroll to Only Affect Container Under Cursor (Task 2 of fix-top-bar-and-scroll plan)

### Changes Made
- **File**: `internal/ui/layout.go`
  - Added zone creation for System blocks (lines 501-526)
- **File**: `internal/app/mouse.go`
  - **handleScrollUp (lines 199-209)**: Use zone-based cursor detection instead of ActiveScrollBlock
  - **handleScrollDown (lines 274-281)**: Use zone-based cursor detection instead of ActiveScrollBlock

### Technical Details
- **Root Cause**: Mouse scroll was using `m.ActiveScrollBlock` (set by pressing `[]`) instead of detecting which block is under the cursor
- **Solution**: 
  - Added zones for System blocks in layout.go using `zt.MarkAt("systemblock-"+fmt.Sprint(i), blockX, blockY, blockWidth, blockHeight)`
  - Use `m.FindZoneAt(m.MouseX, m.MouseY)` to find which "systemblock-N" zone is under the cursor
  - Extract block index from zone ID using `fmt.Sscanf(zone.ID, "systemblock-%d", &blockIdx)`

### Code Context - Zone Creation (layout.go)
```go
// Create zones for System blocks
if s.SystemBlockCount > 0 {
    cols := 1
    if s.Width >= 80 {
        cols = 2
    }
    if s.Width >= 120 {
        cols = 3
    }
    numRows := (s.SystemBlockCount + cols - 1) / cols
    rowHeight := availHeight / numRows
    if rowHeight < 1 {
        rowHeight = 1
    }

    for i := 0; i < s.SystemBlockCount; i++ {
        row := i / cols
        col := i % cols
        blockY := listStartY + row*rowHeight
        blockHeight := rowHeight
        blockX := col * (s.Width / cols)
        blockWidth := s.Width / cols
        if col == cols-1 {
            blockWidth = s.Width - blockX
        }
        zt.MarkAt("systemblock-"+fmt.Sprint(i), blockX, blockY, blockWidth, blockHeight)
    }
}
```

### Code Context - Scroll Handler (mouse.go)
```go
// Before (used ActiveScrollBlock)
activeBlock := m.ActiveScrollBlock
if activeBlock >= 0 && m.SystemBlockScrollable[activeBlock] {
    m.SystemBlockScrollOffsets[activeBlock] -= rows
}

// After (uses cursor position)
blockIdx := -1
zone := m.FindZoneAt(m.MouseX, m.MouseY)
if zone != nil && strings.HasPrefix(zone.ID, "systemblock-") {
    fmt.Sscanf(zone.ID, "systemblock-%d", &blockIdx)
}
if blockIdx >= 0 && m.SystemBlockScrollable[blockIdx] {
    m.SystemBlockScrollOffsets[blockIdx] -= rows
}
```

### Notes
- Zone coordinates are absolute (0,0 = top-left)
- Mouse events provide absolute coordinates
- FindZoneAt() returns the topmost zone at the cursor position
- If cursor is not over a system block, blockIdx stays -1 and no scroll occurs
- This allows scrolling ANY scrollable block just by hovering over it (no need to focus with `[]`)

### Verification
- Build passes with exit code 0
- Mouse scroll now only affects container under cursor on System tab
- Works for both scroll up and scroll down

---

## Task 8: Fix Arrow Key Bug - Move System Handling Outside Logs Block

### Change Made
- **File**: `internal/app/update.go`
- **Lines 700-713**: Moved System block handling from nested inside Logs block to separate else if branch

### Technical Details
- **Root Cause**: Arrow key (j/down) handling had System block scroll nested inside Logs block's else-if
- Specifically: `if m.LogsScrollOffset < maxScroll { ... } else if currentTab == "System" ...`
- This meant System scrolling only executed when Logs was at max scroll (wrong!)
- **Solution**: Move System to separate else if branch at same level as Logs

### Code Context
```go
// Before (BUGGY - System nested inside Logs else-if)
} else if currentTab == "Logs" && len(m.SystemLogs) > 0 {
    // ... Logs scrolling ...
    if m.LogsScrollOffset < maxScroll {
        m.LogsScrollOffset++
    } else if currentTab == "System" && ... {  // BUG: Only runs when Logs at max!
        m.SystemBlockScrollOffsets[m.ActiveScrollBlock]++
    }
}

// After (FIXED - System as separate else if)
} else if currentTab == "Logs" && len(m.SystemLogs) > 0 {
    // ... Logs scrolling ...
    if m.LogsScrollOffset < maxScroll {
        m.LogsScrollOffset++
    }
} else if currentTab == "System" && m.SystemBlockCount > 0 && m.ActiveScrollBlock >= 0 && m.SystemBlockScrollable[m.ActiveScrollBlock] {
    m.SystemBlockScrollOffsets[m.ActiveScrollBlock]++
}
```

### Notes
- The "up" arrow key (k/up) at lines 715-733 already had the correct structure
- Only the "down" arrow key (j/down) at lines 700-713 had this bug
- This fix enables arrow keys to work properly for System tab

### Verification
- Build passes with exit code 0
- Arrow keys now work for System tab (after pressing [] to focus a block)

