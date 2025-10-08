# Why WebView must be built on Windows

## The problem

When you try to build (cross compile) WebView from macOS for Windows:
```bash
make build-windows-webview
```

You get:
```
fatal error: 'shlobj.h' file not found
```

## Why this happens

### Cross-Compilation toolchain limitations

The `mingw-w64` toolchain on macOS provides:
- ✅ Basic Windows API headers
- ✅ Win32 API basics
- ✅ Standard C/C++ libraries
- ❌ **Advanced Windows SDK headers** (WebView2 needs these!)

### What WebView2 requires

The `webview_go` package wraps Microsoft's WebView2, which needs:
- `shlobj.h` - Shell objects API
- `shlwapi.h` - Shell lightweight API  
- `shobjidl.h` - Shell object interfaces
- Other Windows SDK headers
- COM (Component Object Model) interfaces

These are **only available with the full Windows SDK**, not in cross-compilation toolchains.

## Your WebView2 runtime is fine if:

When you run this on Windows:
```cmd
reg query "HKLM\SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}"
```

And see version `141.0.3357.57` or similar - that's **perfect!** ✅

That's the **runtime** (what runs WebView apps). You have it installed.

The issue is building the app needs the **SDK** (development headers), which mingw-w64 doesn't have.

## The solution: Build on Windows

### What you can build from macOS:
| Build Type | Works from Mac? | Why |
|------------|----------------|-----|
| Standard (Fyne) | ✅ Yes | Uses basic Win32 APIs |
| Debug | ✅ Yes | Same as standard |
| **WebView** | ❌ **No** | **Needs Windows SDK** |

### What you need to do:

**On your Windows computer:**

```cmd
REM 1. Install Go (if not installed)
REM Download from: https://go.dev/dl/

REM 2. Copy source files to Windows
REM (Or clone from Git)

REM 3. Build with WebView
go get github.com/webview/webview_go
go mod tidy
go build -tags webview -o krankybearnotify-webview.exe

REM 4. Test it
krankybearnotify-webview.exe -force-webview -title "Test" -message "Works!"
```

Or just double-click `build-webview-windows.bat` - it does all of this!

## Why not just use Docker?

You might think: "Can I use Docker on Mac with Windows containers?"

**Still doesn't work easily because:**
- Docker Windows containers need Windows host (or complex Hyper-V setup)
- Or: Docker Linux containers with mingw-w64 = same problem (no Windows SDK)
- Would need full Windows SDK in container = huge, complex

**Building directly on Windows is simpler!**

## Comparison: Standard vs WebView build

### Standard Build (from Mac) ✅
```bash
# On Mac - Works!
make build-windows

# Creates:
bin/WindowsAMD64/krankybearnotify.exe

# Features:
- Fyne GUI (if OpenGL works)
- MessageBox fallback
- -force-basic flag
```

### WebView build (needs Windows) ❌→✅
```bash
# On Mac - Doesn't work
make build-windows-webview
# Error: shlobj.h not found

# On Windows - Works!
go build -tags webview -o krankybearnotify-webview.exe

# Features:
- Everything from standard build
- PLUS: HTML/CSS/JS UI fallback
- Beautiful gradient interface
- Auto-close with countdown
```

## What files to copy to Windows

Minimum needed on Windows:
```
Source files:
- *.go (all Go source files)
- go.mod
- go.sum

Optional but helpful:
- build-webview-windows.bat (build script)
- WINDOWS_VM_BUILD_INSTRUCTIONS.md
- FyneApp.toml
- Resources/ (for icons, if you want)
```

## Your options

### Option 1: Use what you have ✅
```cmd
REM The builds from Mac work fine with force-basic
krankybearnotify.exe -force-basic -title "Alert" -message "Message"
```
Uses MessageBox - simple but works everywhere.

### Option 2: Build WebView on Windows (Better!) ⭐
```cmd
REM Build once on Windows with webview support
go build -tags webview -o krankybearnotify-webview.exe

REM Then use beautiful HTML UI
krankybearnotify-webview.exe -force-webview -title "Alert" -message "Message"
```
Uses HTML/CSS/JS - much better UI with animations and countdown timer!

## Summary

| Question | Answer |
|----------|--------|
| Can I build standard version from Mac? | ✅ Yes |
| Can I build WebView version from Mac? | ❌ No |
| Why not? | Needs Windows SDK headers |
| Do I have WebView2 runtime? | ✅ Yes (v141.0.3357.57) |
| What do I need to build WebView? | Build on Windows with Go |
| Is it hard? | No - 3 commands (or use build script) |
| Is it worth it? | Yes! Much better UI than MessageBox |

## Next Steps

1. ✅ **Use `-force-basic` for now** (works with Mac-built exe)
2. 🔨 **Install Go on your Windows VM** (https://go.dev/dl/)
3. 📦 **Copy source files to Windows**
4. 🎨 **Build with WebView** (3 commands or run script)
5. 🎉 **Enjoy beautiful HTML/CSS UI!**

---

**TL;DR:** WebView needs Windows SDK headers that mingw-w64 doesn't have. Build on Windows - it's easy!

See: `WINDOWS_VM_BUILD_INSTRUCTIONS.md` for step-by-step guide.

