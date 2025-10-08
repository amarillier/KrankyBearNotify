# Custom Icon Feature - Quick Guide

## Overview

Version 1.1.0 adds support for displaying custom images in your notifications! Icons appear on the left side of the notification content, making your alerts more visually distinctive and professional.

## Basic Usage

```bash
# Simple notification with icon
./krankybearnotify -title "Success!" -message "Build completed" -icon "./KrankyBearBeret.png"

# With timeout
./krankybearnotify -title "Alert" -message "Deployment finished" -icon "./icon.png" -timeout 5
```

## Icon Specifications

- **Recommended Size**: 64x64 pixels
- **Supported Formats**: PNG, JPEG, GIF, and other common image formats
- **Path Types**: Both relative and absolute paths are supported
- **Aspect Ratio**: Automatically preserved (images are scaled proportionally)

## Icon Behavior

### Graceful Fallback
If an icon file is not found or fails to load, the notification will still display correctly without the icon. A warning is logged but the notification proceeds normally.

```bash
# If icon.png doesn't exist, notification shows without icon
./krankybearnotify -title "Test" -message "Hello" -icon "./nonexistent.png"
# Output: Warning: Icon file not found: ./nonexistent.png
# Notification still displays with text only
```

### Layout
- Icon appears on the **left side** of the notification
- Vertical separator between icon and content
- Window width automatically adjusted to 500px (increased from 400px)
- Icon is centered vertically within its container

## Examples with Bundled Icons

The project includes three KrankyBear icons you can use:

```bash
# Beret icon - casual/informational
./krankybearnotify -title "Info" -message "Task started" -icon "./KrankyBearBeret.png"

# Fedora icon - success/completion
./krankybearnotify -title "Success" -message "All tests passed!" -icon "./KrankyBearFedoraRed.png"

# Hard hat icon - warning/construction
./krankybearnotify -title "Warning" -message "Maintenance required" -icon "./KrankyBearHardHat.png"
```

## Integration Examples

### Shell Script Function
```bash
notify_with_icon() {
    local title="$1"
    local message="$2"
    local icon="${3:-}"
    local timeout="${4:-10}"
    
    if [ -n "$icon" ] && [ -f "$icon" ]; then
        krankybearnotify -title "$title" -message "$message" -icon "$icon" -timeout "$timeout"
    else
        krankybearnotify -title "$title" -message "$message" -timeout "$timeout"
    fi
}

# Usage
notify_with_icon "Build Status" "Compilation successful" "./success.png" 5
```

### Python
```python
def notify(title, message, icon=None, timeout=10):
    cmd = ['krankybearnotify', '-title', title, '-message', message, '-timeout', str(timeout)]
    if icon and os.path.exists(icon):
        cmd.extend(['-icon', icon])
    subprocess.run(cmd)

# Usage
notify("Backup Complete", "All files backed up", icon="./backup-icon.png")
```

## Tips

1. **Icon Size**: While any size works, 64x64 or 128x128 pixels provide the best balance
2. **File Format**: PNG with transparency works best for clean-looking notifications
3. **Path Resolution**: Relative paths are resolved from the current working directory
4. **Performance**: Icon loading is fast and doesn't significantly impact notification display time
5. **Reusability**: Define icon constants in your scripts for consistent visual identity

## Testing

Test icon loading:
```bash
# Test all bundled icons
go test -v -run TestLoadIcon

# Test with your own icon
./krankybearnotify -title "Test" -message "Does my icon work?" -icon "./my-icon.png" -timeout 3
```

## Technical Details

### Implementation
- Uses Fyne's `canvas.Image` for rendering
- Images loaded via `storage.NewFileURI()`
- Fill mode: `ImageFillContain` (preserves aspect ratio)
- Minimum size constraint: 64x64 pixels
- Container: Horizontal box layout with vertical separator

### Error Handling
- File existence check before loading
- Nil checks for failed image loads
- Warning logs for debugging
- Automatic fallback to text-only layout

## Upgrade Notes

If upgrading from v1.0.0:
- No breaking changes
- `-icon` flag is optional (backward compatible)
- Window width increased from 400 to 500 pixels
- All existing scripts continue to work without modification

## See Also

- Main README: [README.md](README.md)
- Examples: [examples/README.md](examples/README.md)
- Changelog: [CHANGELOG.md](CHANGELOG.md)
- FyneApp.toml: Configuration file with version 1.1.0

---

**Enjoy your enhanced notifications!** 🎉
