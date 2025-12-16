# bubbleMonitor

[![Go Version](https://img.shields.io/badge/Go-1.25.5-blue.svg)](https://golang.org/dl/)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey.svg)](https://github.com/N1xev/bubbleMonitor/releases)
[![License](https://img.shields.io/badge/license-AGPLv3-blue.svg)](LICENSE)

A beautiful terminal-based system monitor built with Go and BubbleTea. Track your system metrics in real-time with a slick TUI interface.

**"shows you only what you want to see! 😄"**

![bubbleMonitor Screenshot](https://github.com/user-attachments/assets/8929a57d-5160-4ef8-9169-69e1e42af11f)

## What's Inside

Monitor everything that matters: CPU usage (per-core!), memory consumption, disk space, network activity, battery status, and system temperatures. Manage processes with ease—filter, sort, kill, suspend, or check what files they're using.

Choose from 30+ gorgeous themes (Dracula, Nord, Gruvbox, Tokyo Night, and more), customize borders, and pick your favorite chart style. Works beautifully on Windows, Linux, and macOS.

## Getting Started

Grab the latest binary for your system:

```bash
# Windows
curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.0/bub-windows-amd64-v0.1.0.exe -o bub.exe

# Linux
curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.0/bub-linux-amd64-v0.1.0 -o bub
chmod +x bub

# macOS (Apple Silicon)
curl -L https://github.com/N1xev/bubbleMonitor/releases/download/v0.1.0/bub-darwin-arm64-v0.1.0 -o bub
chmod +x bub
```

Or build from source if you're feeling adventurous:

```bash
git clone https://github.com/N1xev/bubbleMonitor.git
cd bubbleMonitor
go build -o bub main.go
```

Then just run `./bub` and you're good to go!

## Keyboard Shortcuts

- `Tab` / `1-6` - Navigate between tabs
- `P` - Pause/resume monitoring
- `S` - Sort processes
- `f` - Filter processes
- `K` - Kill selected process
- `z` / `x` - Suspend/resume process
- `.` - Open settings
- `?` - Show all shortcuts
- `Q` - Quit

## Configuration

bubbleMonitor creates a config file at `~/.config/bubble-monitor/config.json` with sensible defaults. Tweak the refresh rate, history length, theme, or set custom alert thresholds for CPU, memory, disk, and temperature.

Want your own colors? Switch to the `custom` theme and define your palette:

```json
{
  "theme": "custom",
  "custom_theme": {
    "primary": "#7D56F4",
    "secondary": "#EE6FF8",
    "success": "#A1E3AD",
    "warning": "#F5A962",
    "alert": "#F25D94"
  }
}
```

## Platform Notes

Most features work everywhere, but there are a few quirks:

- **Windows**: Temperature monitoring might need admin privileges. Load averages aren't available.
- **Linux**: Full support for everything, including load averages.
- **macOS**: Load averages show as "N/A", and GPU monitoring is limited.

GPU monitoring works best with NVIDIA cards across all platforms.

## Contributing

Found a bug or have an idea? Open an issue or submit a pull request! Fork the repo, create a branch, make your changes, and send it over.

```bash
git checkout -b feature/cool-new-thing
git commit -m 'Add cool new thing'
git push origin feature/cool-new-thing
```

## Built With

- [BubbleTea](https://github.com/charmbracelet/bubbletea) - The amazing TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - For making things pretty
- [gopsutil](https://github.com/shirou/gopsutil) - Cross-platform system info
- [battery](https://github.com/distatus/battery) - Battery monitoring

## License

GNU Affero General Public License v3.0 - see [LICENSE](LICENSE) for details.

---

**Made with ❤️ by [Alaa Elsamouly](https://github.com/N1xev)**