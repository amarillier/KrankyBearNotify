# Building WebView Version on Windows

## Why Build on Windows

The WebView version with `-tags webview` requires Windows SDK headers that are not easily available in cross-compilation from macOS/Linux. The easiest way is to build directly on your Windows machine.

## Prerequisites (Windows)

1. **Go** - Install from https://go.dev/dl/
2. **Git** (optional) - For cloning the repo
3. **WebView2 Runtime** - Usually pre-installed on Windows 10/11

## Build Instructions

### On Your Windows VM

```cmd
REM 1. Navigate to your project directory
cd C:\Users\Allan\Desktop\KrankyBearNotify

REM 2. Get the WebView dependency
go get github.com/webview/webview_go
go mod tidy

REM 3. Build with WebView support
go build -tags webview -o krankybearnotify-webview.exe

REM 4. Test it
krankybearnotify-webview.exe -force-webview -title "Test" -message "WebView works!"
```

### Build Both Versions

```cmd
REM Standard version (MessageBox fallback)
go build -o krankybearnotify.exe

REM WebView version (HTML/CSS/JS UI)
go build -tags webview -o krankybearnotify-webview.exe
```

## Usage

### With Force-WebView Flag

```cmd
REM Force WebView mode (best UI)
krankybearnotify-webview.exe -force-webview -title "Alert" -message "Modern UI!"
```

### Auto-Detection Mode

```cmd
REM Let it auto-detect (will use WebView if OpenGL fails)
krankybearnotify-webview.exe -title "Alert" -message "Auto-select UI"
```

### For Your VM (Recommended)

Since your VM fails OpenGL but passes detection, use `-force-webview`:

```cmd
krankybearnotify-webview.exe -force-webview -title "VM Alert" -message "Better than MessageBox!"
```

## Comparison

| Build | OpenGL | WebView UI | MessageBox | Auto-Close | Icons |
|-------|--------|------------|------------|------------|-------|
| Standard | If available | ❌ | ✅ | ❌ | ❌ |
| WebView | If available | ✅ | ✅ | ✅ | ❌* |
| Fyne (works) | ✅ Required | ❌ | ❌ | ✅ | ✅ |

*Icons not displayed in WebView fallback mode

## WebView Features

The WebView UI provides:
- ✅ **Gradient background** (purple/blue)
- ✅ **Smooth animations** (slide-in effect)
- ✅ **Countdown timer** (shows time remaining)
- ✅ **Auto-close** (timeout support)
- ✅ **Modern styling** (rounded corners, shadows)
- ✅ **Responsive layout**

## For Your Specific VM

Your Proxmox VM has partial OpenGL that passes detection but fails Fyne. Solutions:

### Option 1: Force Basic (MessageBox)
```cmd
krankybearnotify.exe -force-basic -title "Alert" -message "Simple"
```
- ✅ Always works
- ❌ Basic UI (standard Windows dialog)
- ❌ No auto-close
- ❌ No styling

### Option 2: Force WebView (Better!) ⭐
```cmd
krankybearnotify-webview.exe -force-webview -title "Alert" -message "Beautiful"
```
- ✅ Always works (no OpenGL needed)
- ✅ Modern UI (HTML/CSS)
- ✅ Auto-close with countdown
- ✅ Professional appearance

## Create Wrapper Script

### webview-notify.bat
```batch
@echo off
REM WebView notification wrapper
krankybearnotify-webview.exe -force-webview -title %1 -message %2 -timeout %3
```

Usage:
```cmd
webview-notify.bat "Build Complete" "Deployment successful" 30
```

## Troubleshooting

### Error: "WebView2 runtime not found"

**Install WebView2:**
Download from: https://developer.microsoft.com/microsoft-edge/webview2/

Or use the installer:
```cmd
winget install Microsoft.EdgeWebView2Runtime
```

### Error: "WebView forced but not available"

**Cause:** Built without `-tags webview`

**Solution:**
```cmd
go build -tags webview -o krankybearnotify-webview.exe
```

### Check if WebView is available

```cmd
REM If you have the WebView build, try this:
krankybearnotify-webview.exe -force-webview -title "Test" -message "Testing WebView"
```

If it shows an error about WebView not being available, the build doesn't have WebView support.

## File Sizes

| Build | Size | Features |
|-------|------|----------|
| Standard | ~50MB | Fyne + MessageBox |
| WebView | ~52MB | Fyne + WebView + MessageBox |

The WebView build is only slightly larger but provides much better fallback UI!

## Recommended Setup for Your VM

1. Build WebView version on Windows:
   ```cmd
   go build -tags webview -o krankybearnotify.exe
   ```

2. Always use `-force-webview`:
   ```cmd
   krankybearnotify.exe -force-webview -title "Alert" -message "Message"
   ```

3. Or create an alias/wrapper script to make it automatic

## Summary

- ❌ **Cross-compile WebView from Mac**: Complex, needs Windows SDK headers
- ✅ **Build on Windows**: Simple, just add `-tags webview`
- ✅ **Use `-force-webview`**: Best UI for your VM
- ✅ **Much better than MessageBox**: Professional appearance with animations

---

**TL;DR:**  
On your Windows VM: `go build -tags webview -o krankybearnotify.exe`  
Then use: `krankybearnotify.exe -force-webview -title "Alert" -message "Message"`

