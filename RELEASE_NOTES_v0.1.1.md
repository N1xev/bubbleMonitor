# v0.1.1 - Stability and Bug Fixes

## Overview
This release addresses thread safety, memory management, and error handling issues identified during review.

## Fixes

### 1. RingBuffer Concurrency
- **Problem**: Concurrent access to RingBuffer from multiple goroutines caused potential race conditions.
- **Solution**: Added `sync.RWMutex` to protect all field access. RLock used for reads (Get, Len), Lock for writes (Push).

### 2. String Interner Memory Leak
- **Problem**: Global string interner grew indefinitely.
- **Solution**: Implemented cache pruning when size exceeds 5000 entries.

### 3. CSV Export Errors
- **Problem**: Silent failures during file writes.
- **Solution**: Added proper error checking and user notification via toast messages.

## Testing
- Verified with `go test -race`
- All tests passing.

## Installation

### Binary
Download the latest release for your platform from the [Releases](https://github.com/N1xev/bubbleMonitor/releases) page.

### Source
```bash
git clone https://github.com/N1xev/bubbleMonitor.git
cd bubbleMonitor
go build -o bub ./cmd/bub
```
