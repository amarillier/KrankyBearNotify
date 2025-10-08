# Testing KrankyBear Notify on Windows VM (Proxmox)

## Quick Test Instructions

After rebuilding with the improved OpenGL detection, test on your Windows 11 Proxmox VM:

### 1. Copy the New Executable to Your VM

```bash
# On your Mac (build host)
make build-windows

# Copy to VM (adjust paths/method as needed)
scp bin/WindowsAMD64/krankybearnotify.exe user@your-vm-ip:C:/path/
# OR use shared folder, RDP copy, etc.
```

### 2. Test OpenGL Detection

```cmd
REM On Windows VM command prompt
cd C:\path\to\executable

REM Test OpenGL detection
krankybearnotify.exe -check-opengl
```

**Expected output in VM without GPU:**
```
OpenGL check: No suitable pixel format found (likely no OpenGL drivers)
OpenGL is not available
Will use native Windows MessageBox as fallback
Exit code: 1
```

**What this means:**
✅ Detection now correctly identifies that OpenGL won't work  
✅ App will automatically use MessageBox instead of crashing

### 3. Test Actual Notification (MessageBox Fallback)

```cmd
REM Simple notification
krankybearnotify.exe -title "Test Notification" -message "This should work in your VM!"
```

**Expected behavior:**
- Should display a Windows MessageBox (native dialog)
- Should NOT crash with "window creation error"
- Should NOT try to initialize Fyne/OpenGL

**If you see:**
- ✅ Windows MessageBox appears → SUCCESS!
- ❌ "WGL: the driver does not appear to support opengl" → Old version, rebuild needed
- ❌ App crashes → Check if you're using the new executable

### 4. Test with Custom Parameters

```cmd
REM Notification with title and timeout
krankybearnotify.exe -title "Build Complete" -message "Your deployment succeeded!" -timeout 10

REM Check version
krankybearnotify.exe -version
```

### 5. (Optional) Test with WebView

If you want a nicer UI than MessageBox, build with WebView support:

```bash
# On Mac build host
go get github.com/webview/webview_go
go mod tidy
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" go build -tags webview -ldflags="-w -s -H windowsgui" -o krankybearnotify-webview.exe

# Copy to VM and test
```

On Windows VM:
```cmd
REM Test with WebView
krankybearnotify-webview.exe -title "Test webview Notification" -message "This uses HTML/CSS/JS!"
```

**Expected:** HTML-based notification window with gradient background (uses Edge WebView2)

## Troubleshooting

### Issue: Still getting "WGL: the driver does not appear to support opengl"

**Solution:** You're still using the old executable. Make sure you:
1. Rebuilt with `make build-windows`
2. Copied the NEW executable from `bin/WindowsAMD64/krankybearnotify.exe`
3. Overwrote the old one on the VM
4. Running the new version (check with `-version`)

### Issue: MessageBox appears but says "Auto-close not supported"

**This is expected!** MessageBox doesn't support auto-close timeout. The message will say:
```
Your notification message

(Auto-close not supported in fallback mode)
```

This is normal for the MessageBox fallback.

### Issue: Want auto-close in VM

**Solution:** Use WebView build instead:
```bash
go build -tags webview ...
```

WebView supports:
- ✅ Auto-close with countdown timer
- ✅ Webview gradient UI
- ✅ Custom styling
- ✅ No OpenGL required

## Logging

To see detailed logging of what's happening:

```cmd
REM Run from command prompt (not double-click)
krankybearnotify.exe -title "Test" -message "Check logs"
```

You should see output like:
```
OpenGL check: No suitable pixel format found (likely no OpenGL drivers)
Warning: OpenGL not available, trying alternative GUI
Using native Windows MessageBox
```

## Expected Results Summary

| Test | Old Version | New Version |
|------|-------------|-------------|
| `-check-opengl` | Returns true (wrong) | Returns false (correct) |
| Show notification | CRASH | MessageBox appears ✅ |
| With `-timeout 10` | CRASH | MessageBox with note ✅ |
| Exit behavior | Error exit | Clean exit ✅ |

## What Fixed It

The old detection only checked if `opengl32.dll` existed:
```go
// Old (broken)
dll := syscall.NewLazyDLL("opengl32.dll")
return err == nil  // DLL exists in VM, returns true
```

The new detection actually tests if OpenGL works:
```go
// New (working)
pixelFormat := choosePixelFormat.Call(hdc, &pfd)
if pixelFormat == 0 {
    // No suitable format = No OpenGL support
    return false  // Returns false in VM!
}
```

The key is `ChoosePixelFormat()` - this Windows API call fails in VMs without GPU drivers, even though the DLL exists.

## For Production Use

Once you've verified it works:

1. **Standard build** (MessageBox fallback):
   ```bash
   make build-windows
   ```
   - Smaller binary (~50MB)
   - Always works in VMs
   - Basic MessageBox UI

2. **WebView build** (Better UI):
   ```bash
   go build -tags webview -o krankybearnotify.exe
   ```
   - Slightly larger (~51-52MB)
   - Modern HTML/CSS UI
   - Requires Edge WebView2 (usually pre-installed on Windows 10/11)

## Questions?

If you still have issues:
1. Check the logs (run from command prompt, not double-click)
2. Verify you're using the newly built executable
3. Try `-check-opengl` to see what detection returns
4. Check `krankybearnotify.exe -version` to confirm version

The fix is specifically for your reported issue:
> "Fyne error: window creation error Cause APIUnavailable: WGL: the driver does not appear to support opengl"

This should no longer happen - the app will detect the issue and use MessageBox instead! 🎉

---

**Test Date:** October 8, 2025  
**Target:** Windows 11 VM on Proxmox  
**Expected:** MessageBox fallback works without crash  
**Build:** v0.1.0 with improved OpenGL detection

