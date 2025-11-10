#!/bin/bash
# Advanced glibc bundle creator with version detection and compatibility checks

set -e

echo "=== Advanced glibc Bundle Creator ==="
echo ""

# Check if notify binary exists
if [ ! -f "bin/notify-linux-amd64" ]; then
    echo "ERROR: bin/notify-linux-amd64 not found. Please build it first."
    exit 1
fi

# Detect current system glibc version
CURRENT_GLIBC=$(ldd --version | head -n1 | grep -o '[0-9]\+\.[0-9]\+')
echo "Current system glibc version: $CURRENT_GLIBC"

# Create bundle directory with version info
BUNDLE_DIR="notify-with-glibc-${CURRENT_GLIBC}"
mkdir -p "$BUNDLE_DIR"

echo "Creating bundle: $BUNDLE_DIR"
echo ""

# Copy the notify binary
cp bin/notify-linux-amd64 "$BUNDLE_DIR/notify"
echo "✓ Copied notify binary"

# Create lib directory
mkdir -p "$BUNDLE_DIR/lib"

# Function to copy library safely
copy_lib() {
    local lib_name="$1"
    local lib_path="$2"
    
    if [ -f "$lib_path" ]; then
        cp "$lib_path" "$BUNDLE_DIR/lib/"
        echo "✓ Copied $lib_name"
        return 0
    else
        echo "⚠️  Warning: $lib_name not found at $lib_path"
        return 1
    fi
}

echo ""
echo "=== Extracting glibc libraries ==="

# Get all required libraries
ldd bin/notify-linux-amd64 > /tmp/ldd_output.txt

# Copy essential glibc libraries
copy_lib "libc.so.6" "$(grep 'libc\.so\.6' /tmp/ldd_output.txt | awk '{print $3}')"
copy_lib "libm.so.6" "$(grep 'libm\.so\.6' /tmp/ldd_output.txt | awk '{print $3}')"
copy_lib "libpthread.so.0" "$(grep 'libpthread\.so\.0' /tmp/ldd_output.txt | awk '{print $3}')"
copy_lib "libdl.so.2" "$(grep 'libdl\.so\.2' /tmp/ldd_output.txt | awk '{print $3}')"

# Copy optional libraries (may not exist on all systems)
copy_lib "libcrypt.so.1" "$(grep 'libcrypt\.so\.1' /tmp/ldd_output.txt | awk '{print $3}')" || true
copy_lib "libnsl.so.1" "$(grep 'libnsl\.so\.1' /tmp/ldd_output.txt | awk '{print $3}')" || true
copy_lib "libresolv.so.2" "$(grep 'libresolv\.so\.2' /tmp/ldd_output.txt | awk '{print $3}')" || true

# Copy X11 libraries (for GUI)
copy_lib "libX11.so.6" "$(grep 'libX11\.so\.6' /tmp/ldd_output.txt | awk '{print $3}')" || true
copy_lib "libXcursor.so.1" "$(grep 'libXcursor\.so\.1' /tmp/ldd_output.txt | awk '{print $3}')" || true
copy_lib "libXrandr.so.2" "$(grep 'libXrandr\.so\.2' /tmp/ldd_output.txt | awk '{print $3}')" || true
copy_lib "libXinerama.so.1" "$(grep 'libXinerama\.so\.1' /tmp/ldd_output.txt | awk '{print $3}')" || true
copy_lib "libXi.so.6" "$(grep 'libXi\.so\.6' /tmp/ldd_output.txt | awk '{print $3}')" || true
copy_lib "libXxf86vm.so.1" "$(grep 'libXxf86vm\.so\.1' /tmp/ldd_output.txt | awk '{print $3}')" || true

# Copy OpenGL libraries
copy_lib "libGL.so.1" "$(grep 'libGL\.so\.1' /tmp/ldd_output.txt | awk '{print $3}')" || true
copy_lib "libGLU.so.1" "$(grep 'libGLU\.so\.1' /tmp/ldd_output.txt | awk '{print $3}')" || true

# Copy dynamic linker
ld_path=$(grep "ld-linux-x86-64.so.2" /tmp/ldd_output.txt | awk '{print $1}')
if [ -f "$ld_path" ]; then
    cp "$ld_path" "$BUNDLE_DIR/"
    echo "✓ Copied ld-linux-x86-64.so.2"
else
    echo "⚠️  Warning: Could not find ld-linux-x86-64.so.2"
fi

# Create enhanced wrapper script
cat > "$BUNDLE_DIR/notify-wrapper.sh" << 'EOF'
#!/bin/bash

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Set LD_LIBRARY_PATH to include bundled libraries
export LD_LIBRARY_PATH="$SCRIPT_DIR/lib:$LD_LIBRARY_PATH"

# Set other environment variables for X11
export DISPLAY="${DISPLAY:-:0}"

# Check if we have the bundled dynamic linker
if [ -f "$SCRIPT_DIR/ld-linux-x86-64.so.2" ]; then
    echo "Using bundled glibc libraries..."
    exec "$SCRIPT_DIR/ld-linux-x86-64.so.2" "$SCRIPT_DIR/notify" "$@"
else
    echo "Using system glibc libraries..."
    exec "$SCRIPT_DIR/notify" "$@"
fi
EOF

chmod +x "$BUNDLE_DIR/notify-wrapper.sh"

# Create comprehensive test script
cat > "$BUNDLE_DIR/test-bundle.sh" << 'EOF'
#!/bin/bash

echo "=== Testing glibc bundle ==="
echo ""

echo "System information:"
echo "OS: $(uname -a)"
echo "glibc version: $(ldd --version | head -n1)"
echo ""

echo "Bundle contents:"
ls -la "$(dirname "$0")"
echo ""

echo "1. Testing notify version:"
if ./notify-wrapper.sh -version; then
    echo "✓ Version test passed"
else
    echo "✗ Version test failed"
fi

echo ""
echo "2. Testing GUI check:"
if ./notify-wrapper.sh -check-gui; then
    echo "✓ GUI check passed"
else
    echo "✗ GUI check failed (may be normal on headless systems)"
fi

echo ""
echo "3. Testing wall check:"
if ./notify-wrapper.sh -check-wall; then
    echo "✓ Wall check passed"
else
    echo "✗ Wall check failed"
fi

echo ""
echo "4. Testing library dependencies:"
ldd notify 2>/dev/null | head -10

echo ""
echo "=== Bundle test complete ==="
EOF

chmod +x "$BUNDLE_DIR/test-bundle.sh"

# Create deployment instructions
cat > "$BUNDLE_DIR/DEPLOYMENT.md" << EOF
# Glibc Bundle Deployment Instructions

## Bundle Information
- Created on: $(date)
- Source glibc version: $CURRENT_GLIBC
- Target architecture: x86_64

## Deployment Steps

1. **Copy the entire bundle directory** to your target system
2. **Make scripts executable**:
   \`\`\`bash
   chmod +x notify-wrapper.sh
   chmod +x test-bundle.sh
   \`\`\`

3. **Test the bundle**:
   \`\`\`bash
   ./test-bundle.sh
   \`\`\`

4. **Use the wrapper** instead of the direct binary:
   \`\`\`bash
   ./notify-wrapper.sh -title "Test" -message "Hello World"
   \`\`\`

## Troubleshooting

- If you get "No such file or directory" errors, the target system may be missing required libraries
- If GUI doesn't work, ensure X11 is running and DISPLAY is set
- For headless systems, use wall broadcast mode: \`./notify-wrapper.sh -force-wall\`

## Security Note

This bundle includes system libraries. Only use on trusted systems.
EOF

echo ""
echo "=== Bundle creation complete ==="
echo ""
echo "Bundle directory: $BUNDLE_DIR"
echo "Size: $(du -sh "$BUNDLE_DIR" | cut -f1)"
echo ""
echo "Contents:"
ls -la "$BUNDLE_DIR"

echo ""
echo "Next steps:"
echo "1. Test locally: cd $BUNDLE_DIR && ./test-bundle.sh"
echo "2. Copy to target system"
echo "3. Run: ./notify-wrapper.sh instead of ./notify"
echo ""
echo "See $BUNDLE_DIR/DEPLOYMENT.md for detailed instructions"
