# Building with WebView on Linux

## Prerequisites

WebView on Linux requires `webkit2gtk`. Install it first:

### Ubuntu/Debian
```bash
sudo apt update
sudo apt install -y libwebkit2gtk-4.0-dev pkg-config
```

### Fedora/RHEL
```bash
sudo dnf install -y webkit2gtk3-devel pkg-config
```

### Arch Linux
```bash
sudo pacman -S webkit2gtk pkg-config
```

### openSUSE
```bash
sudo zypper install webkit2gtk3-devel pkg-config
```

## Building Steps

1. **Install system dependencies** (see above)

2. **Ensure CGO is enabled:**
```bash
export CGO_ENABLED=1
```

3. **Get the webview dependency:**
```bash
go get github.com/webview/webview_go
go mod tidy
```

4. **Build with webview tag:**
```bash
go build -tags webview
```

## Alternative: Standard Build

If you don't need WebView fallback, just build normally:
```bash
go build
```

This gives you Fyne GUI (if OpenGL available) with native MessageBox fallback.

## Troubleshooting

### Error: "no required module provides package github.com/webview/webview_go"

**Solution:**
```bash
go get github.com/webview/webview_go
go mod tidy
go build -tags webview
```

### Error: "Package webkit2gtk-4.0 was not found"

**Solution:** Install webkit2gtk development libraries (see Prerequisites above)

### Error: "CGO is not enabled"

**Solution:**
```bash
export CGO_ENABLED=1
go build -tags webview
```

### WebView crashes at runtime

**Solution:** Install webkit2gtk runtime (usually installed with -dev package):
```bash
# Ubuntu/Debian
sudo apt install libwebkit2gtk-4.0-37

# Fedora
sudo dnf install webkit2gtk3
```

## Quick Start (One-liner)

### Ubuntu/Debian
```bash
sudo apt install -y libwebkit2gtk-4.0-dev pkg-config && export CGO_ENABLED=1 && go get github.com/webview/webview_go && go mod tidy && go build -tags webview
```

### Fedora
```bash
sudo dnf install -y webkit2gtk3-devel pkg-config && export CGO_ENABLED=1 && go get github.com/webview/webview_go && go mod tidy && go build -tags webview
```

## Verify Build

Test that webview is working:
```bash
# This should not show any webview-related errors
./krankybearnotify -version

# Test a notification
./krankybearnotify -title "Test" -message "WebView build works!" -timeout 5
```

## Docker Build (Optional)

If you want to build in a container:

```dockerfile
FROM golang:1.21-bullseye

RUN apt-get update && apt-get install -y \
    libwebkit2gtk-4.0-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build
COPY . .

RUN go get github.com/webview/webview_go && \
    go mod tidy && \
    CGO_ENABLED=1 go build -tags webview -o krankybearnotify

CMD ["./krankybearnotify"]
```

Build:
```bash
docker build -t krankybearnotify-webview .
```

## Notes

- WebView adds ~5-10MB to binary size
- Requires webkit2gtk runtime on target systems
- Falls back to native dialogs if WebView2GTK not available
- Standard build (without webview) is recommended for most Linux users since OpenGL is widely available

