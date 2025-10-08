# Windows computer Troubleshooting - Still Getting Fyne Error

## Issue
After recompiling, still getting error on Windows computer:
```
Fyne error: window creation error
Cause APIUnavailable: WGL: the driver does not appear to support opengl
```

## Diagnosis Steps

### Step 1: Use the DEBUG Build

The regular build uses `-H windowsgui` which hides console output. You need the DEBUG version to see what's happening.

**I've created a debug build for you:**
```
bin/WindowsAMD64/krankybearnotify-debug.exe
```

**Copy this to your computer and run from Command Prompt:**

```cmd
cd C:\path\to\executable
krankybearnotify-debug.exe -check-opengl
```

### Step 2: Check What the Logs Say

Look for these specific log messages:

**If OpenGL detection is working:**
```
OpenGL check: No suitable pixel format found (likely no OpenGL drivers)
OpenGL is not available
Will use native Windows MessageBox as fallback
Exit code: 1
```

**If OpenGL detection is falsely passing:**
```
OpenGL check: OpenGL appears to be available and functional  
OpenGL is available
Fyne GUI can be used
Exit code: 0
```

If you see the second output, the detection is WRONG for your computer.

### Step 3: Test Actual Notification with Debug Build

```cmd
krankybearnotify-debug.exe -title "Test" -message "Check console output"
```

**Look for these log lines:**
```
OpenGL availability check result: false   ← Should be FALSE
Warning: OpenGL not available, trying alternative GUI
Using native Windows MessageBox
```

**If instead you see:**
```
OpenGL availability check result: true    ← Wrong!
Attempting to create Fyne GUI (OpenGL detected as available)
[Then Fyne crashes]
```

This means the OpenGL detection isn't working correctly for your specific computer configuration.

## Possible Causes

### Cause 1: Wrong Executable
You might be running the old version.

**Solution:**
```cmd
REM Check version and date
krankybearnotify-debug.exe -version

REM Delete old versions
del krankybearnotify.exe
del krankybearnotify-old.exe

REM Copy new debug build
copy krankybearnotify-debug.exe krankybearnotify.exe
```

### Cause 2: Your computer Has Partial OpenGL Support
Some computers have just enough OpenGL to pass the pixel format test but not enough for Fyne.

**Check your computer's OpenGL:**
```cmd
REM Install OpenGL Extensions Viewer (free tool)
REM Or use DirectX Diagnostic Tool
dxdiag

REM Check Display tab for:
REM - Driver Model: WDDM or XDDM?
REM - Feature Levels: What version?
```

**If you see:**
- "Microsoft Basic Display Adapter" → No real OpenGL
- "Standard VGA" → No real OpenGL  
- "QXL" (QEMU/KVM/Proxmox) → Limited or no OpenGL
- Feature Level 9.x or lower → Old drivers

### Cause 3: Proxmox-Specific Issue
Proxmox VMs with QXL graphics adapter often have this issue.

**Check your VM config:**
```bash
# On Proxmox host
qm config <VMID> | grep vga

# Common outputs:
vga: qxl     ← No real OpenGL support
vga: std     ← Basic VGA, no OpenGL
vga: vmware  ← Might have basic OpenGL
```

## Solutions

### Solution 1: Force Fallback Mode (Temporary Fix)

Modify the OpenGL check to always return false for testing:

**Edit `gui_opengl_windows.go`:**
```go
func isOpenGLAvailable() bool {
    // TEMPORARY: Force false for VM testing
    return false
    
    // ... rest of function commented out ...
}
```

Rebuild and test:
```bash
make build-windows
```

If this works, it confirms the detection is the problem.

### Solution 2: More Aggressive OpenGL Detection

Try actually creating a GL context, not just checking pixel format:

```go
// In gui_opengl_windows.go
func isOpenGLAvailable() bool {
    // ... existing checks ...
    
    // Actually try to create an OpenGL context
    wglCreateContext := opengl32.NewProc("wglCreateContext")
    hglrc, _, _ := wglCreateContext.Call(hdc)
    if hglrc == 0 {
        log.Println("OpenGL check: Failed to create GL context")
        return false
    }
    defer wglDeleteContext.Call(hglrc)
    
    // Try to make it current
    wglMakeCurrent := opengl32.NewProc("wglMakeCurrent")  
    ret, _, _ := wglMakeCurrent.Call(hdc, hglrc)
    if ret == 0 {
        log.Println("OpenGL check: Failed to make GL context current")
        return false
    }
    
    return true
}
```

### Solution 3: Skip Fyne Entirely in VMs

Add a command-line flag to force MessageBox mode:

```go
// main.go
forceBasicGUI := flag.Bool("force-basic", false, "Force basic GUI (skip OpenGL/Fyne)")

if *forceBasicGUI {
    if runtime.GOOS == "windows" {
        showWindowsMessageBox(*title, *message, *timeout)
        os.Exit(0)
    }
}
```

Then use:
```cmd
krankybearnotify.exe -force-basic -title "Test" -message "Hello"
```

### Solution 4: Use WebView Build (Recommended)

WebView doesn't need OpenGL:

```bash
# On Mac build host
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" \
  go build -tags webview -ldflags="-w -s" -o krankybearnotify-webview-debug.exe
```

Test on VM:
```cmd
krankybearnotify-webview-debug.exe -title "Test" -message "WebView test"
```

WebView should work even without OpenGL!

## Diagnostic Commands to Run

Run these in your VM and send me the output:

```cmd
REM 1. Check version
krankybearnotify-debug.exe -version

REM 2. Check OpenGL detection
krankybearnotify-debug.exe -check-opengl

REM 3. Try notification with full output
krankybearnotify-debug.exe -title "Diagnostic Test" -message "Testing" 2>&1

REM 4. Check Windows version
ver

REM 5. Check display adapter
wmic path win32_VideoController get name

REM 6. Check if OpenGL DLL exists
dir C:\Windows\System32\opengl32.dll
```

## What I Need to Know

To help fix this, please provide:

1. **Output from diagnostic commands above**

2. **Your Proxmox VM config:**
   ```bash
   # On Proxmox host:
   qm config <VMID>
   ```

3. **What you see when running debug build:**
   - Does it say "OpenGL availability check result: true" or "false"?
   - What happens next?

4. **VM specs:**
   - Proxmox version?
   - Windows 11 version?
   - Display adapter in Windows Device Manager?

## Quick Test: Force MessageBox

While we diagnose, you can force MessageBox mode by modifying the code:

**In `main.go`, line ~210, change:**
```go
openglAvailable := isOpenGLAvailable()
```

**To:**
```go
openglAvailable := false  // FORCE FALSE FOR TESTING
```

Rebuild and test. If this works, we know the issue is the detection.

## Expected vs Actual

### What SHOULD Happen
```
1. Run krankybearnotify.exe
2. Check GUI available → TRUE (Windows)
3. Check OpenGL available → FALSE (VM)
4. Log: "OpenGL not available, trying alternative GUI"
5. Use MessageBox
6. Success!
```

### What's PROBABLY Happening
```
1. Run krankybearnotify.exe
2. Check GUI available → TRUE
3. Check OpenGL available → TRUE (WRONG!)
4. Log: "Attempting to create Fyne GUI"
5. Fyne tries to create window
6. WGL error → CRASH
```

The key is figuring out why step 3 returns TRUE when it should return FALSE for your VM.

---

**Next Steps:**
1. Use the debug build (`krankybearnotify-debug.exe`)
2. Run the diagnostic commands
3. Share the output with me
4. We'll fix the detection based on your specific VM configuration

Let me know what you find!

