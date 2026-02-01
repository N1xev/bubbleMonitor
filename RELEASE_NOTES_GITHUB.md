# v0.1.1 - Stability and Bug Fixes

This release focuses on improving stability and fixing critical bugs discovered during code review.

## 🐛 Bug Fixes

### Critical
- **Fixed potential race condition in RingBuffer** 
  - Added `sync.RWMutex` to protect concurrent access to history buffers
  - Future-proofs against background tasks accessing metrics history
  - Verified with race detector tests

- **Prevented memory leak in string interner**
  - Global interner was accumulating strings indefinitely
  - Now prunes cache when it exceeds 5000 entries  
  - Prevents memory issues during long monitoring sessions

- **Fixed unsafe type assertion in slice pool**
  - `GetProcSlice()` now uses comma-ok idiom
  - Returns fresh slice if pool gets corrupted instead of panicking

### High Priority
- **Added proper error handling for CSV exports**
  - Previously, write operations silently ignored errors
  - Users now get clear feedback if snapshot export fails
  - Prevents corrupted files being reported as successful

- **Removed duplicate code (copy-paste error)**
  - Fixed triple assignment to `m.LastError` in help overlay handler

## 📦 Installation

### Download Pre-built Binaries

Choose the binary for your platform:

#### Windows
```bash
curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.1/bub-windows-amd64-v0.1.1.exe -o bub.exe
./bub.exe
```

#### Linux
```bash
curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.1/bub-linux-amd64-v0.1.1 -o bub
chmod +x bub
./bub
```

#### macOS (Intel)
```bash
curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.1/bub-darwin-amd64-v0.1.1 -o bub
chmod +x bub
./bub
```

#### macOS (Apple Silicon)
```bash
curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.1/bub-darwin-arm64-v0.1.1 -o bub
chmod +x bub
./bub
```

### Build from Source

```bash
git clone https://github.com/N1xev/bubbleMonitor.git
cd bubbleMonitor
go build -o bub ./cmd/bub
./bub
```

## 🧪 Testing
- Comprehensive test suite with race detector
- All tests passing on Linux/macOS/Windows
- No memory leaks detected

## 📝 Notes
- All features and functionality remain unchanged
- Full backward compatibility maintained
- No breaking changes

---

## 📊 What's Changed
* Refactor project structure to follow Go conventions by @N1xev in 2f18905
* Fix potential race condition in RingBuffer by @N1xev in 5ff87e2
* Prevent memory leak in string interner and fix type assertion by @N1xev in 2e25728
* Add proper error handling for CSV writes by @N1xev in e11a028

**Full Changelog**: https://github.com/N1xev/bubbleMonitor/compare/v0.1.0...v0.1.1
