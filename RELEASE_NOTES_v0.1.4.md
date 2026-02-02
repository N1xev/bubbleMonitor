# bubbleMonitor v0.1.3 - Enhanced CPU Monitoring & Optimization

This release addresses user feedback regarding high resource usage and confusion about CPU metrics >100%.

## 🚀 New Features

### Normalized CPU Usage (Btop Style)
By default, Linux reports CPU usage as the sum of all cores (e.g., 200% = 2 full cores). This can be confusing.
We've added a **Normalized Mode** option!

- **Raw Mode** (Default): Shows total usage (e.g., 400% on 4 cores). Accurate for capacity planning.
- **Normalized Mode**: Shows 0-100% scale (divided by core count). Better for "how busy is this process relative to the whole machine".

**How to use:**
- Press `n` to toggle instantly.
- Or change "Process CPU" in the Settings menu (`.`).

## ⚡ Performance Improvements

We addressed reports of high CPU usage by `bub` itself (17.5%):

- **Optimized Rendering**: The process list rendering engine was rewritten to pre-allocate styles, eliminating thousands of memory allocations per frame.
- **Smarter Updates**: Background analysis logic was split to avoid redundant calculations, reducing CPU overhead during monitoring cycles.
- **Memory Fixes**: Fixed a potential memory leak in process history tracking.

## 🛠️ Fixes

- Fixed stale data in history graphs.
- Improved responsiveness of the settings menu.

---

### Installation

**Go Install**
```bash
go install github.com/N1xev/bubbleMonitor/cmd/bub@v0.1.3
```

**Binaries**
Download the pre-compiled binary for your platform from the Assets section.
