# KrankyBear Notify Examples

This directory contains example scripts demonstrating how to use KrankyBear Notify.

## Scripts

### notify-example.sh (Linux/macOS)

A bash script that demonstrates various notification scenarios.

**Usage:**

```bash
chmod +x notify-example.sh
./notify-example.sh
```

### notify-example.ps1 (Windows)

A PowerShell script that demonstrates various notification scenarios.

**Usage:**

```powershell
.\notify-example.ps1
```

### notify-example.py (Cross-Platform)

A Python script that demonstrates using the unified `notify.py` wrapper. Works on Windows, Linux, and macOS.

**Usage:**

```bash
# Linux/macOS
chmod +x notify-example.py
./notify-example.py

# Windows
python notify-example.py
```

**Note:** This example uses the `notify.py` wrapper script (located in the parent directory), which provides a unified interface for sending notifications across all platforms.

## Examples Included

1. **Simple Notification**: Shows a notification with default settings
2. **Custom Title and Message with Icon**: Demonstrates custom text with a KrankyBear icon
3. **Warning Notification with Icon**: Shows a warning-style notification with hard hat icon
4. **Success Notification with Icon**: Shows a success-style notification with fedora icon
5. **Manual Close**: Shows a notification that requires user interaction

All examples with icons will gracefully fall back to text-only notifications if the icon files are not found.

## Integration Examples

### Cron Job (Linux/macOS)

Send a notification when a backup completes:

```bash
# Add to crontab (crontab -e)
0 2 * * * /usr/local/bin/backup.sh && DISPLAY=:0 /usr/local/bin/notify -title "Backup Complete" -message "Nightly backup finished successfully"
```

### Task Scheduler (Windows)

Create a scheduled task that shows a notification:

```powershell
$action = New-ScheduledTaskAction -Execute "C:\Program Files\KrankyBear\notify.exe" -Argument '-title "Reminder" -message "Time for a break!" -timeout 10'
$trigger = New-ScheduledTaskTrigger -Daily -At 3PM
Register-ScheduledTask -Action $action -Trigger $trigger -TaskName "Break Reminder" -Description "Shows a break reminder"
```

### systemd Service (Linux)

Create a user service that shows notifications:

```ini
# ~/.config/systemd/user/notify@.service
[Unit]
Description=KrankyBear Notification: %i

[Service]
Type=oneshot
Environment="DISPLAY=:0"
ExecStart=/usr/local/bin/notify -title "System Notification" -message "%i"

[Install]
WantedBy=default.target
```

Trigger it with:

```bash
systemctl --user start notify@"Your message here"
```

### Shell Script Integration

```bash
#!/bin/bash

# Function to send notification with optional icon
notify() {
    local title="$1"
    local message="$2"
    local timeout="${3:-10}"
    local icon="${4:-}"
    
    if notify -check-gui; then
        if [ -n "$icon" ]; then
            notify -title "$title" -message "$message" -timeout "$timeout" -icon "$icon"
        else
            notify -title "$title" -message "$message" -timeout "$timeout"
        fi
    else
        echo "[$title] $message"
    fi
}

# Use in your script
if ./run-tests.sh; then
    notify "Tests Passed" "All tests completed successfully!" 5 "./success-icon.png"
else
    notify "Tests Failed" "Some tests failed. Check the logs." 0 "./error-icon.png"
fi
```

### Python Integration (Using notify.py wrapper)

The recommended way to integrate with Python is to use the `notify.py` wrapper script:

```python
#!/usr/bin/env python3
import subprocess
from urllib.parse import quote

def notify(title, message, button="OK", icon="", timeout=5, width=400, height=200, mode=""):
    """Send a notification using notify.py wrapper"""
    # URL-encode parameters
    title_enc = quote(title)
    message_enc = quote(message)
    button_enc = quote(button)
    icon_enc = quote(icon) if icon else ""
    mode_enc = quote(mode) if mode else ""
    
    # Build command
    cmd = [
        'python3',  # or 'python' on Windows
        '../notify.py',  # Adjust path as needed
        title_enc,
        message_enc,
        button_enc,
        icon_enc,
        str(timeout),
        str(width),
        str(height),
        mode_enc
    ]
    
    subprocess.run(cmd)

# Example usage
if __name__ == '__main__':
    notify("Python Script", "Task completed successfully!", icon="./icon.png")
```

See `notify-example.py` for a complete working example.

### Python Integration (Direct call to notify)

Alternatively, call the notify executable directly:

```python
#!/usr/bin/env python3
import subprocess
import sys

def notify(title, message, timeout=10, icon=None):
    """Send a notification using notify"""
    try:
        cmd = [
            'notify',
            '-title', title,
            '-message', message,
            '-timeout', str(timeout)
        ]
        if icon:
            cmd.extend(['-icon', icon])
        subprocess.run(cmd, check=True)
    except subprocess.CalledProcessError:
        print(f"[{title}] {message}", file=sys.stderr)

# Example usage
if __name__ == '__main__':
    notify("Python Script", "Task completed successfully!", 5, "./icon.png")
```

### Node.js Integration

```javascript
const { spawn } = require('child_process');

function notify(title, message, timeout = 10, icon = null) {
    return new Promise((resolve, reject) => {
        const args = [
            '-title', title,
            '-message', message,
            '-timeout', timeout.toString()
        ];
        
        if (icon) {
            args.push('-icon', icon);
        }
        
        const proc = spawn('notify', args);

        proc.on('close', (code) => {
            if (code === 0) {
                resolve();
            } else {
                reject(new Error(`Notification failed with code ${code}`));
            }
        });
    });
}

// Example usage
notify('Node.js App', 'Build completed!', 5, './build-icon.png')
    .then(() => console.log('Notification sent'))
    .catch(err => console.error('Notification failed:', err));
```

## Tips

1. **Always check GUI availability** before showing notifications in automated scripts
2. **Use appropriate timeouts** - shorter for informational messages, longer or 0 for important alerts
3. **Run in background** with `&` (bash) or `Start-Process` (PowerShell) if you don't want to block
4. **Set DISPLAY** environment variable on Linux when running from cron or systemd
5. **Use emojis** in titles/messages for better visual communication (✓, ⚠️, ❌, ℹ️, etc.)
6. **Custom icons** - Use PNG, JPEG, or other image formats for custom icons (64x64 pixels recommended)
7. **Icon paths** - Can be relative or absolute paths; icons are gracefully skipped if not found
