# Draft: Zone Fixes - Process Hover/Click & SamLab Transparency

## Issues Identified

1. **SamLab overlay transparency** - `Background(bg)` always applied in samlab.go line 38
2. **Process zones off-by-one** - rowPadding=2 added to zone Y, shifting zones down by 1
3. **Process hover not working** - Will work once zone Y is fixed

## Root Cause Analysis

### SamLab Issue
In `overlays/samlab.go` line 38:
```go
boxStyle := lipgloss.NewStyle().
    Border(border).
    BorderForeground(b).
    Background(bg).  // Always applied - should be conditional!
```
Same pattern as context menu fix - needs to check `s.BackgroundOpaque`.

### Process Zone Issue
In `tabs/processes.go` line 199:
```go
rowY := listStartY + rowPadding + (i - startIdx)
```
Where:
- `listStartY = topBarH + topGap` (layout.go line 606)
- `rowPadding = 2` (processes.go line 20)
- `topGap = 1`

The zones are shifted by (rowPadding - topGap) = 1 row down. Should be just `listStartY + (i - startIdx)`.

## Fixes to Apply

1. **samlab.go**: Make background conditional on s.BackgroundOpaque
2. **processes.go**: Remove rowPadding from zone Y calculation
