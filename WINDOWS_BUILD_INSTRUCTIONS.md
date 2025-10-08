# Build WebView version on your Windows computer

## Why you will get an error trying to cross compile on MacOS

The `.exe` files built on macOS **don't include WebView support**. We need to rebuild with the `-tags webview` flag **on a Windows computer**.

## What you need

1. **Go** - Install if you don't have it
2. **Source code** - Either clone the repo or just have the source files
3. **WebView2 Runtime** - Usually pre-installed on Windows 10/11

## Step-by-Step instructions

### Option 1: Using the build script (easiest)

1. **Copy these files to your Windows computer:**
   - All `.go` files
   - `go.mod` and `go.sum`
   - `build-webview-windows.bat` (the build script)

2. **Double-click `build-webview-windows.bat`**
   
   It will:
   - Check if Go is installed
   - Download WebView dependency
   - Build with WebView support
   - Create `krankybearnotify-webview.exe`

3. **Test it:**
   ```cmd
   krankybearnotify-webview.exe -force-webview -title "Test" -message "Success!"
   ```

### Option 2: Manual build

**In PowerShell or Command Prompt:**

```cmd
REM 1. Check Go is installed
go version

REM 2. Get WebView package
go get github.com/webview/webview_go

REM 3. Tidy modules
go mod tidy

REM 4. Build with WebView support (THE KEY!)
go build -tags webview -o krankybearnotify-webview.exe

REM 5. Test it
krankybearnotify-webview.exe -force-webview -title "Test" -message "WebView!"
```

## If you don't have Go installed

### Install Go on Windows:

**Option A: Download Installer**
1. Go to: https://go.dev/dl/
2. Download Windows installer (e.g., `go1.21.x.windows-amd64.msi`)
3. Run installer
4. Restart Command Prompt/PowerShell

**Option B: Using winget**
```cmd
winget install GoLang.Go
```

**Option C: Using Chocolatey**
```cmd
choco install golang
```

**Verify installation:**
```cmd
go version
```

## If you don't have the source code

### Option A: Clone from GitHub
```cmd
git clone https://github.com/amarillier/krankybearnotify.git
cd krankybearnotify
```

### Option B: Copy source files

Copy these files to your Windows computer:
```
Required files:
- *.go (all Go source files)
- go.mod
- go.sum
- build-webview-windows.bat

Optional:
- FyneApp.toml
- Resources/ (if you want icons)
```

## Checking WebView2 runtime

WebView needs Microsoft Edge WebView2 (usually pre-installed):

```cmd
REM Check if installed
reg query "HKLM\SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}"
```

If not installed:
```cmd
REM Install via winget
winget install Microsoft.EdgeWebView2Runtime
```

Or download: https://go.microsoft.com/fwlink/p/?LinkId=2124703

## Expected results

### ✅ Success looks like:
```cmd
PS C:\> go build -tags webview -o krankybearnotify-webview.exe
PS C:\> .\krankybearnotify-webview.exe -force-webview -title "Test" -message "Hello!"

[WebView gradient window appears with "Test" title and "Hello!" message]
[Countdown timer shows "Auto-closing in 10s"]
```

### ❌ Common errors:

**Error: "go: command not found"**
- **Fix:** Install Go (see above)

**Error: "cannot find package"**
- **Fix:** Run `go get github.com/webview/webview_go` then `go mod tidy`

**Error: "WebView2 runtime not found"**
- **Fix:** Install WebView2 Runtime (see above)

**Error: Still says "WebView forced but not available"**
- **Fix:** Make sure you used `-tags webview` in the build command!

## File comparison

After building, you'll have:

| File | Size | Has WebView |
|------|------|-------------|
| `krankybearnotify.exe` (from Mac) | ~22MB | ❌ No |
| `krankybearnotify-debug.exe` (from Mac) | ~22MB | ❌ No |
| `krankybearnotify-webview.exe` (built on Windows) | ~23MB | ✅ Yes |

## Quick test commands

```cmd
REM Test basic mode (works with any build)
krankybearnotify.exe -force-basic -title "Test 1" -message "MessageBox"

REM Test webview mode (only works if built with -tags webview)
krankybearnotify-webview.exe -force-webview -title "Test 2" -message "WebView UI!"

REM Test with auto-close countdown
krankybearnotify-webview.exe -force-webview -title "Alert" -message "Closes in 5s" -timeout 5
```

## What you get with WebView

The WebView build provides a much better UI than the Windows basic MessageBox:

```
┌────────────────────────────────────┐
│  Gradient Background (Purple/Blue) │
│                                    │
│  📢 Test                          │
│                                    │
│  WebView works!                   │
│                                    │
│            [ OK ]                  │
│                                    │
│  Auto-closing in 10s              │
└────────────────────────────────────┘
```

vs MessageBox:
```
┌─────────────────────┐
│ Test                │
├─────────────────────┤
│ This WILL work!     │
│                     │
│      [ OK ]         │
└─────────────────────┘
```

## Summary

1. **Install Go** (if not already installed)
2. **Copy source files to Windows VM**
3. **Run build script** OR manually: `go build -tags webview -o krankybearnotify-webview.exe`
4. **Test:** `krankybearnotify-webview.exe -force-webview -title "Test" -message "Success!"`

The key difference: You MUST use `-tags webview` when building on Windows!

---

**Need help?** The build script (`build-webview-windows.bat`) does all of this automatically!

