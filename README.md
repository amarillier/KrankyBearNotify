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

## GUI Fallback System

The app automatically selects the best available GUI method:

### 1. **Fyne GUI** (Primary - Requires OpenGL)
- Modern interface
- Custom icons support
- Auto-close with timeout
- Full feature set

### 2. **WebView GUI** (Optional Fallback - HTML/CSS/JavaScript)
- Used when OpenGL is not available (if compiled with `-tags webview`)
- Web based gradient UI with animations
- Auto-close with countdown timer
- Works in VMs and Remote Desktop
- **Requires**: Build with `-tags webview` flag
- **Runtime**: WebView2 on Windows, WebKit on macOS/Linux

### 3. **Native Dialog** (GUI Fallback)
- Windows: MessageBox API
- macOS/Linux: Basic system dialogs
- Simple text-only notifications
- Manual close only

### 4. **Wall Broadcast** (Linux No-GUI Fallback)
- Used when no GUI is detected on Linux
- Sends message to all logged-in users' terminals
- Perfect for headless servers and SSH sessions
- Formatted text broadcast with timestamps
- Auto-expiry notifications

**Test availability:**
- OpenGL: `./krankybearnotify -check-opengl`
- Wall: `./krankybearnotify -check-wall` (Linux only)

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

**Standard build:**
```bash
git clone https://github.com/amarillier/krankybearnotify.git
cd krankybearnotify
go build
```

**Build with WebView support (optional, for better fallback UI):**

*Note: Requires webkit2gtk on Linux. No extra dependencies on macOS/Windows.*
- macOS: See [BUILD_WEBVIEW_MACOS.md](BUILD_WEBVIEW_MACOS.md)
- Linux: See [BUILD_WEBVIEW_LINUX.md](BUILD_WEBVIEW_LINUX.md)

```bash
# Linux: Install webkit2gtk first
# Ubuntu/Debian: sudo apt install libwebkit2gtk-4.0-dev pkg-config
# Fedora: sudo dnf install webkit2gtk3-devel pkg-config

# Then build with webview tag
go get github.com/webview/webview_go
go mod tidy
go build -tags webview
```

This enables HTML/CSS/JavaScript UI as fallback when OpenGL is not available.

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
| `-width` | Window width in pixels | 400 |
| `-height` | Window height in pixels | 250 |
| `-icon`, `-image` | Path to icon image file (PNG, JPEG, etc.) | "" (no icon) |
| `-check-gui` | Check if GUI mode is available and exit | false |
| `-check-opengl` | Check if OpenGL is available and exit (Windows) | false |
| `-check-wall` | Check if wall broadcast is available (Linux) and exit | false |
| `-force-basic` | Force basic GUI mode (skip OpenGL, use MessageBox/WebView) | false |
| `-force-webview` | Force WebView mode (HTML/CSS/JS UI, requires webview build) | false |
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

On Windows systems, check if OpenGL drivers are actually functional (not just if the DLL exists):

```bash
./krankybearnotify -check-opengl
```

The detection performs a comprehensive test:
- ✅ Checks if `opengl32.dll` exists
- ✅ Verifies WGL functions are present
- ✅ Tests device context availability
- ✅ **Attempts to choose a pixel format** (the critical test)

**Why this matters for VMs:**
In virtualized environments (Proxmox, VirtualBox, VMware) without GPU passthrough, `opengl32.dll` exists but OpenGL doesn't actually work. The improved detection catches this by testing pixel format selection, which fails without real OpenGL drivers.

If OpenGL is not available, the app automatically falls back to:
1. WebView (if compiled with `-tags webview`)
2. Native Windows MessageBox (always available)

**Example in a VM without GPU:**
```bash
$ ./krankybearnotify -check-opengl
OpenGL check: No suitable pixel format found (likely no OpenGL drivers)
OpenGL is not available
Will use native Windows MessageBox as fallback
Exit code: 1

$ ./krankybearnotify -title "Test" -message "This works!"
# Automatically uses MessageBox - no crash!
```

See [OPENGL_DETECTION_IMPROVED.md](OPENGL_DETECTION_IMPROVED.md) for technical details.

### Force Basic GUI Mode (VM Workaround)

If you're running in a VM where OpenGL detection passes but Fyne still fails to initialize, use the `-force-basic` flag to skip OpenGL entirely:

```bash
./krankybearnotify -force-basic -title "VM Test" -message "This bypasses OpenGL"
```

**When to use:**
- Proxmox/KVM VMs where GL context tests pass but Fyne window creation fails
- VMs with partial OpenGL support
- Remote Desktop sessions where OpenGL is unreliable
- Any environment where you want guaranteed MessageBox fallback

**What it does:**
- Skips all OpenGL detection
- Never attempts to initialize Fyne
- Goes directly to MessageBox (Windows) or basic fallback
- Works with `-tags webview` builds for better UI

**Example for scripts:**
```bash
# Always use basic mode in VM
krankybearnotify.exe -force-basic -title "Alert" -message "Server status"
```

See [IMMEDIATE_WORKAROUND.md](IMMEDIATE_WORKAROUND.md) for more details.

### Force WebView Mode (Better UI Option)

If you want a better UI than MessageBox with animations and auto-close support:

**First, build with WebView on your Windows machine:**
```cmd
REM On Windows VM
go get github.com/webview/webview_go
go mod tidy
go build -tags webview -o krankybearnotify.exe
```

**Then use force-webview flag:**
```cmd
krankybearnotify.exe -force-webview -title "Modern Alert" -message "Beautiful HTML/CSS UI!"
```

**Benefits over MessageBox:**
- ✅ Beautiful gradient UI with animations
- ✅ Auto-close with countdown timer  
- ✅ Modern styling (rounded corners, shadows)
- ✅ Professional appearance
- ✅ No OpenGL required

**Recommended for VMs:**
```cmd
REM Best option for your Proxmox VM
krankybearnotify.exe -force-webview -title "Alert" -message "Much better than MessageBox!"
```

See [BUILD_WEBVIEW_WINDOWS.md](BUILD_WEBVIEW_WINDOWS.md) for detailed build instructions.

**Note:** WebView requires WebView2 runtime (pre-installed on Windows 10/11). Cross-compilation from macOS is complex; build directly on Windows for best results.

### Check Wall Broadcast (Linux)

On Linux systems, check if the `wall` command is available for broadcasting to all logged-in users:

```bash
./krankybearnotify -check-wall
```

This is useful for headless servers or SSH sessions where no GUI is available. The app will automatically use wall broadcast as a fallback when no GUI is detected.

**Example wall broadcast output:**
```
================================================================
  IMPORTANT ALERT
================================================================

Your server backup has completed successfully!

[This notification will be displayed for 30 seconds]
================================================================
Sent: 2025-10-08 14:30:45
```

All logged-in users will see this message in their terminal.

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

# Use your own custom icon (can also use -image as an alias for -icon)
./krankybearnotify -title "Alert" -message "Custom notification" -image "/path/to/your/icon.png"
```

### Custom Window Dimensions

```bash
# Make a larger window for more content
./krankybearnotify -title "Large Notification" -message "This is a larger window with more space for text" -width 600 -height 350

# Make a smaller, compact notification
./krankybearnotify -title "Compact" -message "Small alert" -width 300 -height 150

# Custom size with icon and long message
./krankybearnotify -title "Custom Layout" -message "A notification with custom dimensions and an icon. The text will wrap properly to fit the window width." -icon "./KrankyBearBeret.png" -width 500 -height 300
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

If you're running this on a Linux server without a GUI, the application will automatically use wall broadcast to notify all logged-in users:

```bash
# Check GUI status
$ ./krankybearnotify -check-gui
GUI mode is not available
$ echo $?
1

# Check wall broadcast status
$ ./krankybearnotify -check-wall
Wall broadcast is available
Can send notifications to all logged-in users
$ echo $?
0

# Send a notification (will use wall broadcast automatically)
$ ./krankybearnotify -title "Server Alert" -message "Backup completed"
GUI not available, using wall broadcast

Broadcast Message from root@server
        (somewhere) at 14:30 ...

================================================================
  SERVER ALERT
================================================================

Backup completed

[This notification will be displayed for 10 seconds]
================================================================
Sent: 2025-10-08 14:30:45
```

All logged-in users will see this message in their terminals.

### macOS

On macOS, the application checks if the WindowServer process is running. This is the standard way to detect if the GUI is available.

### Windows

On Windows, the application checks if the process has access to a window station, which indicates GUI availability.

**GUI Fallback System:**

The app tries multiple notification methods in order:

1. **Fyne (OpenGL)** - Full-featured modern UI
2. **WebView (HTML/JS)** - Web-based UI when OpenGL unavailable (optional build)
3. **Native Dialog** - Basic text notifications (Windows MessageBox, etc.)
4. **Wall Broadcast** - Terminal broadcast to all users (Linux only, when no GUI)

**WebView Requirements (optional):**
- Windows: WebView2 Runtime (usually pre-installed on Windows 10/11)
- macOS: Uses native WebKit
- Linux: Requires webkit2gtk

**Wall Broadcast (Linux):**
- Used automatically when no GUI is detected
- Sends notifications to all logged-in users via terminal
- Perfect for headless servers and SSH sessions
- Requires `wall` command (usually pre-installed)

Test methods availability:
```bash
# Test OpenGL
krankybearnotify -check-opengl

# Test wall broadcast (Linux)
krankybearnotify -check-wall
```

**Fallback Comparison:**

| Feature | Fyne | WebView | Native Dialog |
|---------|------|---------|---------------|
| Custom Icons | ✅ | ❌ | ❌ |
| Auto-close | ✅ | ✅ | ❌ |
| Modern UI | ✅ | ✅ | ❌ |
| VMs/RDP | ❌ | ✅ | ✅ |
| GPU Required | Yes | No | No |

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
