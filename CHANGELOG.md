# Changelog

## [v0.1.1] - 2026-02-01

### Bug Fixes

#### Critical
- **Fixed potential race condition in RingBuffer** (#1)
  - Added `sync.RWMutex` to protect concurrent access to history buffers
  - While the current BubbleTea implementation is sequential, this future-proofs against background tasks
  - Includes comprehensive tests with race detector

- **Prevented memory leak in string interner** (#2)
  - Global interner was accumulating strings indefinitely, causing unbounded memory growth
  - Now prunes cache when it exceeds 5000 entries
  - Prevents memory issues during long monitoring sessions

- **Fixed unsafe type assertion in slice pool** (#2)
  - `GetProcSlice()` now uses comma-ok idiom instead of blind type assertion
  - Returns fresh slice if pool gets corrupted instead of panicking

#### High Priority
- **Added proper error handling for CSV exports** (#3)
  - Previously, write operations silently ignored errors
  - Users now get clear feedback if snapshot export fails
  - Prevents corrupted files being reported as successful

- **Removed duplicate code (copy-paste error)** (#3)
  - Fixed triple assignment to `m.LastError` in help overlay handler

### Testing
- Added comprehensive test suite for `RingBuffer` with concurrency tests
- Added test suite for `Interner` including pruning and concurrency tests
- Added tests for slice pool type safety and reusability
- All tests pass with `go test -race`

### Notes
- All features and functionality remain unchanged
- This is a stability and reliability release
- No breaking changes

---

## [v0.1.0] - 2024-12-24

Initial release with full TUI system monitoring functionality.
