# bubbleMonitor

[![Go Version](https://img.shields.io/badge/Go-1.25.5-blue.svg)](https://golang.org/dl/)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey.svg)](https://github.com/N1xev/bubbleMonitor/releases)
[![License](https://img.shields.io/badge/license-AGPLv3-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

A beautiful, cross-platform terminal-based system monitoring tool built with Go and BubbleTea. Real-time system metrics, process management, and customizable themes in an intuitive TUI interface.
**"and it shows you only, what you want to see!😄"**

<!-- ![bubbleMonitor Screenshot](https://via.placeholder.com/800x400/1a1b26/ffffff?text=bubbleMonitor+Screenshot) -->

## Features

### **Real-time System Monitoring**
- **CPU Usage** - Live CPU utilization with per-core monitoring
- **Memory Usage** - RAM consumption with detailed breakdown  
- **Disk Usage** - Storage utilization across all partitions
- **Network Activity** - Real-time bandwidth monitoring (upload/download)
- **Battery Status** - Laptop battery information and health
- **System Temperature** - Hardware temperature monitoring
- **Process Management** - Advanced process control and analysis

### **Customizable Interface**
- **30+ Themes** - Choose from dark, light, and colorful themes
  - Popular themes: Dracula, Nord, Gruvbox, Solarized, Monokai, Catppuccin, Tokyo Night
  - Custom theme support with full color customization
- **Multiple Border Styles** - Single, double, and dashed borders
- **Flexible Layout** - Responsive design that adapts to terminal size
- **Chart Types** - Sparklines, line charts, bar charts, braille charts, and TTY charts

### **Advanced Functionality**
- **Multiple Tabs** - Overview, Metrics, Processes, Disks, Network, System
- **Process Tree View** - Hierarchical process visualization
- **Smart Filtering** - Real-time process filtering and search
- **Process Control** - Kill, suspend, resume, and priority management
- **Alert System** - Configurable thresholds with visual notifications
- **History Tracking** - Configurable metrics history (1m to 1h)
- **Keyboard Shortcuts** - Full keyboard navigation support

### Platform Compatibility

bubbleMonitor provides comprehensive system monitoring across Windows, Linux, and macOS with some platform-specific features and considerations:

| Feature | Windows | Linux | macOS | Notes |
|---------|---------|-------|-------|-------|
| **CPU Monitoring** | ✅ | ✅ | ✅ | Per-core monitoring available on all platforms |
| **Memory Usage** | ✅ | ✅ | ✅ | Full memory statistics and breakdown |
| **Disk Monitoring** | ✅ | ✅ | ✅ | Storage utilization for all partitions |
| **Network Activity** | ✅ | ✅ | ✅ | Real-time bandwidth monitoring |
| **Battery Status** | ✅ | ✅ | ✅ | Laptop battery information and health |
| **System Temperature** | ⚠️ | ✅ | ✅ | Windows may require administrator privileges |
| **Process Management** | ✅ | ✅ | ✅ | List, filter, kill, suspend/resume processes |
| **Process Tree View** | ✅ | ✅ | ✅ | Hierarchical process visualization |
| **Load Averages** | ❌ | ✅ | ❌ | Unix load averages (1m, 5m, 15m) - N/A on Windows/macOS |
| **GPU Monitoring** | ✅ | ✅ | ⚠️ | NVIDIA cards only; macOS has limited support |
| **Process Priority** | ⚠️ | ✅ | ✅ | Windows uses WMIC (requires admin); Unix/macOS use renice |

**Platform-Specific Requirements:**

**Windows:**
- Some system metrics (especially temperature) may require administrator privileges
- Process priority changes use WMIC commands
- Load averages are not available (displays "N/A")

**Linux:**
- Full support for all features including load averages
- Process priority changes use the `renice` command
- GPU monitoring available for NVIDIA cards via `nvidia-smi`

**macOS:**
- Process priority changes use the `renice` command
- Load averages are not available (displays "N/A")
- Limited GPU monitoring support (nvidia-smi not available)

## Quick Start

### Option 1: Download Pre-built Binary

Download the latest release for your platform:

```bash
# Windows
curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.0/bub-windows-amd64-v0.1.0.exe -o bub.exe

# Linux (AMD64)
curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.0/bub-linux-amd64-v0.1.0 -o bub
chmod +x bub

# macOS (Apple Silicon)
curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.0/bub-darwin-arm64-v0.1.0 -o bub
chmod +x bub
```

### Option 2: Build from Source

Ensure you have Go 1.25.5 or later installed:

```bash
# Clone the repository
git clone https://github.com/N1xev/bubbleMonitor.git
cd bubbleMonitor

# Build for your current platform
go build -o bub main.go

# Or build for all platforms (requires Make or similar)
make build-all
```

## Usage

### Basic Usage

Simply run the binary to start monitoring:

```bash
./bub
```

The application will start in your terminal with real-time system monitoring.

### Keyboard Shortcuts Reference

| Action | Shortcut | Description |
|--------|----------|-------------|
| **Navigation** | | |
| Next Tab | `Tab` / `→` / `L` | Move to next tab |
| Previous Tab | `Shift+Tab` / `←` / `h` | Move to previous tab |
| Jump to Tab | `1-6` | Jump directly to specific tab |
| **Controls** | | |
| Pause/Resume | `P` | Pause or resume monitoring |
| Refresh | `R` | Refresh all data |
| Sort Processes | `S` | Cycle through sort options |
| Toggle Help | `?` | Show/hide keyboard shortcuts |
| Quit | `Q` | Exit the application |
| Settings | `.` | Open settings menu |
| History Length | `H` | Cycle history duration |
| Chart Type | `C` | Change chart visualization |
| **Process Management** | | |
| Move Down | `j` / `↓` | Select next process |
| Move Up | `k` / `↑` | Select previous process |
| Top/Bottom | `g` / `G` | Jump to top or bottom |
| Filter | `f` | Filter processes |
| Clear Filter | `c` | Remove active filter |
| Suspend/Resume | `z` / `x` | Suspend or resume process |
| Kill Process | `K` | Terminate selected process |
| Open Files | `o` | View open files for process |
| Tree View | `T` | Toggle hierarchical view |
| Collapse/Expand | `Space` | Collapse or expand tree node |
| Priority | `+` / `-` | Increase or decrease priority |

### Configuration

bubbleMonitor stores configuration in `~/.config/bubble-monitor/config.json`. The application creates this automatically with sensible defaults.

#### Available Themes

bubbleMonitor includes 33 carefully crafted themes:

**Dark Themes:**
- `dark` (default) - Clean dark interface
- `dracula` - Purple and pink color scheme
- `nord` - Arctic blue tones
- `gruvbox` - Warm retro colors
- `tokyonight` - Modern Japanese night theme
- `onedark` - Atom One Dark inspired
- `catppuccin` - Soothing pastel colors
- `material` - Google Material Design colors
- `cyberpunk` - Neon cyberpunk aesthetic

**Light Themes:**
- `light` - Clean light interface
- `github` - GitHub-inspired colors
- `solarized` - Ethan Schoonover's Solarized
- `ocean` - Ocean blue tones

**Special Themes:**
- `custom` - Fully customizable colors
- `tty` - Native ANSI terminal colors

#### Configuration Options

```json
{
  "theme": "dark",
  "refresh_rate": 1000,
  "history_length": 60,
  "chart_type": "sparkline",
  "border_type": "rounded",
  "border_style": "single",
  "background_opaque": false,
  "sort_by": "cpu",
  "view_type": "normal",
  "thresholds": {
    "CPU": 90.0,
    "Memory": 90.0,
    "Disk": 90.0,
    "Temperature": 85.0
  },
  "tabs": ["Overview", "Metrics", "Processes", "Disks", "Network", "System"]
}
```

### Custom Themes

Create your own theme by modifying the `custom` section in the configuration:

```json
{
  "theme": "custom",
  "custom_theme": {
    "primary": "#7D56F4",
    "secondary": "#EE6FF8", 
    "success": "#A1E3AD",
    "warning": "#F5A962",
    "alert": "#F25D94",
    "text": "#F0F0F0",
    "muted": "#A0A0A0",
    "border": "#4A4A4A",
    "background": "#1C1C1C"
  }
}
```

## Development

### Building

Use the included build scripts for multi-platform builds:

```bash
# Windows
build.bat

# Or manually
go build -ldflags="-s -w" -o bub main.go
```

### Project Structure

```
bubbleMonitor/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── build.bat              # Windows build script
├── Makefile               # Unix/Linux/macOS build script
├── src/
│   ├── commands/          # System command implementations
│   │   ├── process/       # Process management
│   │   └── system/        # System information
│   ├── config/            # Configuration management
│   ├── data/              # Data structures and logic
│   ├── messages/          # Inter-component messaging
│   ├── model/             # Application model
│   ├── ui/                # User interface components
│   │   ├── overlays/      # Modal overlays (help, settings)
│   │   ├── tabs/          # Tab-specific views
│   │   └── widgets/       # Reusable UI widgets
│   └── utils/             # Utility functions
└── plans/                 # Development plans
```

### Key Dependencies

- **[BubbleTea v2](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[Lipgloss v2](https://github.com/charmbracelet/lipgloss)** - Style and layout
- **[gopsutil v3](https://github.com/shirou/gopsutil)** - System information
- **[battery](https://github.com/distatus/battery)** - Battery monitoring

## Troubleshooting

### Common Issues

**Application won't start:**
```bash
# Check if binary has execute permissions (Unix/macOS)
chmod +x bub

# Verify Go installation (if building from source)
go version
```

**No system metrics showing:**
- On Windows, some metrics require administrator privileges
- Ensure your terminal supports true color (24-bit) for best experience
- Check that your terminal emulator supports the required Unicode characters

**Performance issues:**
- Increase refresh rate to reduce CPU usage: `H` then select higher interval
- Reduce history length for less memory usage
- Disable background opacity for better performance

**Themes not applying:**
- Verify JSON syntax in configuration file
- Check that theme name matches exactly (case-sensitive)
- Restart the application after changing theme

### Getting Help

- Press `?` in the application for keyboard shortcuts
- Check configuration file: `~/.bubble-monitor/config.json`
- Create an issue on GitHub for bugs or feature requests

## Contributing

We welcome contributions! Please see our contributing guidelines:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Commit changes**: `git commit -m 'Add amazing feature'`
4. **Push to branch**: `git push origin feature/amazing-feature`
5. **Open a Pull Request**

### Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/bubbleMonitor.git
cd bubbleMonitor

# Install dependencies
go mod download

# Run tests
go test ./...

# Build for development
go build -o bub-dev main.go
```

### Code Style

- Follow Go best practices and conventions
- Add tests for new features
- Update documentation for API changes
- Use meaningful commit messages

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- **[Charm Bracelet](https://charm.sh/)** - For BubbleTea and Lipgloss
- **[gopsutil](https://github.com/shirou/gopsutil)** - Cross-platform system utilities
- **Community Contributors** - For feedback, bug reports, and feature requests

## Project Status

- **Version**: 0.1.0
- **Status**: Stable Release
- **Platforms**: Windows, Linux, macOS (AMD64, ARM64)
- **Go Version**: 1.25.5+

---

**Made with ❤️ by the bubbleMonitor team**

For more information, visit our [GitHub repository](https://github.com/N1xev/bubbleMonitor) or check the [documentation](https://github.com/N1xev/bubbleMonitor/wiki).