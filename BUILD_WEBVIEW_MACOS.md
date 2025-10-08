# Building with WebView on macOS

## Quick Start

Building with WebView on macOS is straightforward - no external dependencies required! macOS includes the WebKit framework natively.

```bash
# Get the Go webview package
go get github.com/webview/webview_go

# Tidy dependencies
go mod tidy

# Build with webview tag
go build -tags webview -o krankybearnotify
```

Or use the Makefile:
```bash
make build-webview
```

## Requirements

**None!** macOS includes WKWebView (WebKit) by default.

- ✅ macOS 10.13 (High Sierra) or later
- ✅ No additional packages needed
- ✅ No system dependencies to install

## What You Get

With WebView support enabled, your notification system has three levels of fallback:

1. **Fyne GUI** (preferred) - Full-featured OpenGL UI
2. **WebView** (fallback) - Modern HTML/CSS/JavaScript UI using WKWebView
3. **Native Dialogs** (last resort) - Basic system alerts

## Testing

After building, test that it works:

```bash
# Check version
./krankybearnotify -version

# Show a test notification
./krankybearnotify -title "Test" -message "WebView build works!" -timeout 5

# Check OpenGL availability
./krankybearnotify -check-opengl
```

## Standard Build

If you don't need WebView fallback, just build normally:

```bash
go build -o krankybearnotify
```

This is recommended for most macOS users since OpenGL is typically available.

## Technical Details

- Uses Cocoa's `WKWebView` framework
- No CGO dependencies beyond standard Fyne
- Binary size increase: ~1-2MB
- Zero runtime dependencies

## Troubleshooting

### Error: "no required module provides package github.com/webview/webview_go"

**Solution:**
```bash
go get github.com/webview/webview_go
go mod tidy
go build -tags webview
```

### Linker warning: "ignoring duplicate libraries: '-lobjc'"

**Status:** This is harmless. Both Fyne and webview link against Objective-C libraries. The linker handles this correctly by ignoring duplicates.

### Runtime error: "WebView not available"

**Solution:** This typically means the webview wasn't compiled in. Rebuild with:
```bash
go build -tags webview
```

## Cross-Compilation

To build for macOS from macOS (both architectures):

```bash
# ARM64 (M1/M2/M3 Macs)
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -tags webview -o krankybearnotify-arm64

# AMD64 (Intel Macs)
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -tags webview -o krankybearnotify-amd64
```

## Comparison: With vs Without WebView

| Feature | Standard Build | WebView Build |
|---------|---------------|---------------|
| Fyne GUI | ✅ | ✅ |
| OpenGL Required | ✅ | No (has fallback) |
| HTML/CSS UI Fallback | ❌ | ✅ |
| Binary Size | ~50MB | ~51-52MB |
| macOS Dependencies | Standard | Standard (WKWebView included) |
| Best For | Normal use | Maximum compatibility |

## When to Use WebView Build

Use the webview-enabled build if:

- You need guaranteed GUI even without OpenGL
- You're running in virtual machines or remote desktop
- You want a more visually appealing fallback than system dialogs
- You're distributing to users who might not have GPU acceleration

For most macOS users, the standard build is sufficient since OpenGL is widely available.

