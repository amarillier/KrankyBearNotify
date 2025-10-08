# KrankyBear Notify

A cross-platform notification application built with Go and the Fyne GUI library. This application displays notification windows for any currently logged-in user on macOS, Windows, and Linux with proper GUI mode detection.

## Features

- **Cross-Platform Support**: Works on macOS, Windows, and Linux
  - **Linux**: Works on all desktop environments (GNOME, KDE, XFCE, Cinnamon, MATE, etc.) with X11 or Wayland
  - **macOS**: macOS 10.13 (High Sierra) or later
  - **Windows**: Windows 10 or later
- **GUI Mode Detection**: Automatically detects if a graphical environment is available
  - **Linux**: Checks `DISPLAY`, `WAYLAND_DISPLAY` environment variables and systemd `graphical.target`
  - **macOS**: Checks if WindowServer is running
  - **Windows**: Checks if the process has access to a window station
- **Customizable Notifications**: Configure title, message, timeout, and custom icons
- **Custom Icons**: Display your own images alongside notifications
- **Auto-dismiss**: Optional timeout to automatically close notifications
- **Modern UI**: Clean, simple interface built with Fyne

# Features & Known Issues

## OpenGL Fallback (Windows)
- Automatically detects if OpenGL drivers are available
- Falls back to native Windows MessageBox if OpenGL is not available
- Useful for Windows VMs, Remote Desktop, or systems without GPU acceleration
- Test OpenGL availability: `./krankybearnotify -check-opengl`

## Installation
- This is a standalone / portable app

### Prerequisites

- Go 1.21 or later (for building from source)
- For Linux: X11 or Wayland display server
  - **Works on all major desktop environments**: GNOME, KDE Plasma, XFCE, Cinnamon, MATE, LXQt, Budgie, Deepin, and more
  - **Works on all major distributions**: Ubuntu, Fedora, Debian, Arch, openSUSE, Manjaro, Linux Mint, Pop!_OS, etc.
- For macOS: macOS 10.13 or later
- For Windows: Windows 10 or later

### Build from Source

```bash
git clone https://github.com/amarillier/krankybearnotify.git
cd krankybearnotify
go build
```

### Cross-Platform Building

**Building for Windows from macOS:**

```bash
# Install mingw-w64 for Windows cross-compilation
brew install mingw-w64

# Build for Windows
make build-windows
```

**Building for Linux:**

Cross-compiling Fyne apps to Linux from macOS requires Docker-based `fyne-cross`:

```bash
# Install fyne-cross
go install github.com/fyne-io/fyne-cross@latest

# Build for Linux (amd64 and arm64)
fyne-cross linux -arch=amd64,arm64

# Output will be in fyne-cross/dist/linux-*/ directories
```

Or build directly on Linux:
```bash
go build -ldflags="-w -s" -o krankybearnotify
```

**Note:** Simple `GOOS=linux go build` doesn't work for Fyne apps because they require CGO and Linux-specific libraries.

### Install Dependencies

The project uses Go modules, so dependencies will be automatically downloaded:

```bash
go mod download
```

## Usage

### Getting Help

Show help message with all options and examples:

```bash
./krankybearnotify           # No arguments shows help
./krankybearnotify -h        # Also shows help
./krankybearnotify -help     # Also shows help
```

Show version information:

```bash
./krankybearnotify -version
```

Check for updates:

```bash
./krankybearnotify -checkupdate
# or use the short alias
./krankybearnotify -cu
```

**Note:** Update check data is saved to `latestcheck.json` in the same directory as the executable.

### Basic Usage

Show a notification with default settings:

```bash
./krankybearnotify
```

### Custom Notification

```bash
./krankybearnotify -title "Important Alert" -message "Your task is complete!" -timeout 5
```

### Notification with Icon

```bash
./krankybearnotify -title "Build Complete" -message "Your project built successfully!" -icon "./KrankyBearBeret.png" -timeout 10
```

### Command-Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-title` | Notification title | "Notification" |
| `-message` | Notification message | "This is a notification message" |
| `-timeout` | Auto-close timeout in seconds (0 for no timeout) | 10 |
| `-icon` | Path to icon image file (PNG, JPEG, etc.) | "" (no icon) |
| `-check-gui` | Check if GUI mode is available and exit | false |
| `-check-opengl` | Check if OpenGL is available and exit (Windows) | false |
| `-version` | Show version information and exit | false |
| `-checkupdate`, `-cu` | Check for updates and exit | false |
| `-h`, `-help` | Show help message with examples | - |

### Check GUI Availability

To check if a GUI environment is available without showing a notification:

```bash
./krankybearnotify -check-gui
```

This will exit with code 0 if GUI is available, or code 1 if not.

### Check OpenGL Availability (Windows)

On Windows systems, check if OpenGL drivers are available for Fyne:

```bash
./krankybearnotify -check-opengl
```

If OpenGL is not available, the app will automatically fall back to native Windows MessageBox.

## Examples

### Simple Notification

```bash
./krankybearnotify -title "Hello" -message "World!"
```

### Long-Running Notification

```bash
./krankybearnotify -title "Manual Close" -message "Click OK to close" -timeout 0
```

### Quick Alert

```bash
./krankybearnotify -title "Quick Alert" -message "This will close in 3 seconds" -timeout 3
```

### Custom Icon Notification

```bash
# Use one of the bundled KrankyBear icons
./krankybearnotify -title "Success!" -message "Operation completed" -icon "./KrankyBearFedoraRed.png"

# Use your own custom icon
./krankybearnotify -title "Alert" -message "Custom notification" -icon "/path/to/your/icon.png"
```

### Script Integration

```bash
#!/bin/bash

# Check if GUI is available before showing notification
if ./krankybearnotify -check-gui; then
    ./krankybearnotify -title "Backup Complete" -message "Your backup finished successfully"
else
    echo "GUI not available, skipping notification"
fi
```

## Platform-Specific Notes

### Linux

On Linux, the application checks for GUI availability in the following order:

1. `DISPLAY` environment variable (X11)
2. `WAYLAND_DISPLAY` environment variable (Wayland)
3. systemd `graphical.target` status

**Desktop Environment Compatibility:**
The application works on **all major Linux desktop environments**, including:
- ✅ GNOME (both X11 and Wayland sessions)
- ✅ KDE Plasma (both X11 and Wayland sessions)
- ✅ XFCE
- ✅ Cinnamon
- ✅ MATE
- ✅ LXQt / LXDE
- ✅ Budgie
- ✅ Deepin
- ✅ Any other desktop environment using X11 or Wayland

The detection is desktop-agnostic and only checks for the presence of a display server, not specific desktop APIs. For the `graphical.target` check to work, systemd must be installed and running.

#### Running on Linux Server

If you're running this on a Linux server without a GUI, the application will detect this and exit gracefully:

```bash
$ ./krankybearnotify -check-gui
GUI mode is not available
$ echo $?
1
```

### macOS

On macOS, the application checks if the WindowServer process is running. This is the standard way to detect if the GUI is available.

### Windows

On Windows, the application checks if the process has access to a window station, which indicates GUI availability.

**OpenGL Fallback:**
The app automatically detects if OpenGL drivers are available. If not (common in VMs or Remote Desktop sessions), it falls back to native Windows MessageBox API for notifications. This ensures notifications work even without GPU acceleration.

Test OpenGL availability:
```bash
krankybearnotify -check-opengl
```

**Note:** Fallback mode limitations:
- No custom icons support
- Auto-close timeout not supported (manual close only)
- Simpler UI compared to Fyne

## Development

### Running Tests

Run all tests:

```bash
go test -v
```

Run tests with coverage:

```bash
go test -v -cover
```

Run benchmarks:

```bash
go test -bench=.
```

### Project Structure

```
.
├── main.go                 # Main application entry point
├── gui_check_linux.go      # Linux-specific GUI detection
├── gui_check_darwin.go     # macOS-specific GUI detection
├── gui_check_windows.go    # Windows-specific GUI detection
├── gui_check_other.go      # Stub for other platforms
├── gui_check_test.go       # Tests for GUI detection
├── main_test.go            # Tests for main functionality
├── go.mod                  # Go module definition
└── README.md               # This file
```

### Build Tags

The project uses Go build tags to compile platform-specific code:

- `//go:build linux` - Linux-specific code
- `//go:build darwin` - macOS-specific code
- `//go:build windows` - Windows-specific code
- `//go:build !linux && !darwin && !windows` - Other platforms

## Testing GUI Mode Detection

### Linux

Test the `graphical.target` detection on any desktop environment:

```bash
# Check if graphical.target is active
systemctl is-active graphical.target

# Test with environment variables (X11 - GNOME, KDE, XFCE, etc.)
DISPLAY=:0 ./krankybearnotify -check-gui

# Test with Wayland (GNOME Wayland, KDE Wayland, Sway, etc.)
WAYLAND_DISPLAY=wayland-0 ./krankybearnotify -check-gui

# Quick test notification on your current desktop
./krankybearnotify -title "Desktop Test" -message "If you see this, it works on your DE!" -timeout 5
```

**Testing on Different Desktop Environments:**
```bash
# Works the same on all these:
# - GNOME: Ubuntu, Fedora, Debian (X11 or Wayland)
# - KDE Plasma: Kubuntu, KDE neon, openSUSE (X11 or Wayland)
# - XFCE: Xubuntu, Manjaro XFCE
# - Cinnamon: Linux Mint
# - MATE: Ubuntu MATE
# - And any other X11/Wayland desktop
```

### macOS

Test WindowServer detection:

```bash
# Check if WindowServer is running
pgrep -x WindowServer

# Run the GUI check
./krankybearnotify -check-gui
```

### Windows

Test window station detection:

```powershell
# Run the GUI check
.\krankybearnotify.exe -check-gui
```

## Multi-User Support

This application can be used to send notifications to currently logged-in users. Here are some approaches:

### Linux (systemd)

Create a systemd user service:

```ini
# ~/.config/systemd/user/krankybear-notify.service
[Unit]
Description=KrankyBear Notification Service

[Service]
Type=oneshot
ExecStart=/usr/local/bin/krankybearnotify -title "System Alert" -message "Your notification here"

[Install]
WantedBy=default.target
```

### macOS (launchd)

Create a launch agent:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.krankybear.notify</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/krankybearnotify</string>
        <string>-title</string>
        <string>System Alert</string>
        <string>-message</string>
        <string>Your notification here</string>
    </array>
</dict>
</plist>
```

### Windows (Task Scheduler)

Use Windows Task Scheduler to run the notification for all logged-in users.

## Troubleshooting

### "GUI mode is not available" Error

This error occurs when the application cannot detect a graphical environment. Check:

1. **Linux**: Ensure `DISPLAY` or `WAYLAND_DISPLAY` is set, or that `graphical.target` is active
2. **macOS**: Ensure you're running in a GUI session (not SSH without X11 forwarding)
3. **Windows**: Ensure you're running in a desktop session (not a service or background task)

### Notification Doesn't Appear

1. Check if GUI is available: `./krankybearnotify -check-gui`
2. Verify the timeout isn't too short
3. Check if the window is appearing behind other windows

### Build Errors

If you encounter build errors, ensure you have the required dependencies:

```bash
go mod tidy
go build
```

## Dependencies

- [Fyne](https://fyne.io/) v2.6.3 - Cross-platform GUI toolkit
- Go standard library

## License

This project is provided as-is, free for personal, educational and commercial use, under GNU GPL-3.0

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Author

Allan Marillier

## Acknowledgments

- Built with [Fyne](https://fyne.io/) - An easy-to-use GUI toolkit for Go
- Inspired by the need for simple, cross-platform notification systems
