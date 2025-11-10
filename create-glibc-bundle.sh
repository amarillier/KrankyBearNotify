#!/bin/bash
# Create a glibc bundle for older Linux systems
# This script extracts the required glibc files from a newer system

set -e

echo "=== Creating glibc bundle for older Linux systems ==="
echo ""

# Create bundle directory
BUNDLE_DIR="notify-with-glibc"
mkdir -p "$BUNDLE_DIR"

# Copy the notify binary
if [ -f "bin/notify-linux-amd64" ]; then
    cp bin/notify-linux-amd64 "$BUNDLE_DIR/notify"
    echo "✓ Copied notify binary"
else
    echo "ERROR: bin/notify-linux-amd64 not found. Please build it first."
    exit 1
fi

# Create lib directory for glibc
mkdir -p "$BUNDLE_DIR/lib"

echo ""
echo "=== Extracting required glibc files ==="

# Get the list of required libraries
echo "Analyzing dependencies..."
ldd bin/notify-linux-amd64 | grep -E "(libc\.so|libm\.so|libpthread\.so|libdl\.so|libcrypt\.so|libnsl\.so|libresolv\.so)" > /tmp/required_libs.txt

echo "Required libraries:"
cat /tmp/required_libs.txt

# Copy the required libraries
echo ""
echo "Copying glibc libraries..."

while IFS= read -r line; do
    if [[ $line =~ libc\.so\.6 ]]; then
        lib_path=$(echo "$line" | awk '{print $3}')
        if [ -f "$lib_path" ]; then
            cp "$lib_path" "$BUNDLE_DIR/lib/"
            echo "✓ Copied libc.so.6"
        fi
    elif [[ $line =~ libm\.so\.6 ]]; then
        lib_path=$(echo "$line" | awk '{print $3}')
        if [ -f "$lib_path" ]; then
            cp "$lib_path" "$BUNDLE_DIR/lib/"
            echo "✓ Copied libm.so.6"
        fi
    elif [[ $line =~ libpthread\.so\.0 ]]; then
        lib_path=$(echo "$line" | awk '{print $3}')
        if [ -f "$lib_path" ]; then
            cp "$lib_path" "$BUNDLE_DIR/lib/"
            echo "✓ Copied libpthread.so.0"
        fi
    elif [[ $line =~ libdl\.so\.2 ]]; then
        lib_path=$(echo "$line" | awk '{print $3}')
        if [ -f "$lib_path" ]; then
            cp "$lib_path" "$BUNDLE_DIR/lib/"
            echo "✓ Copied libdl.so.2"
        fi
    elif [[ $line =~ libcrypt\.so\.1 ]]; then
        lib_path=$(echo "$line" | awk '{print $3}')
        if [ -f "$lib_path" ]; then
            cp "$lib_path" "$BUNDLE_DIR/lib/"
            echo "✓ Copied libcrypt.so.1"
        fi
    elif [[ $line =~ libnsl\.so\.1 ]]; then
        lib_path=$(echo "$line" | awk '{print $3}')
        if [ -f "$lib_path" ]; then
            cp "$lib_path" "$BUNDLE_DIR/lib/"
            echo "✓ Copied libnsl.so.1"
        fi
    elif [[ $line =~ libresolv\.so\.2 ]]; then
        lib_path=$(echo "$line" | awk '{print $3}')
        if [ -f "$lib_path" ]; then
            cp "$lib_path" "$BUNDLE_DIR/lib/"
            echo "✓ Copied libresolv.so.2"
        fi
    fi
done < /tmp/required_libs.txt

# Copy the dynamic linker
ld_path=$(ldd bin/notify-linux-amd64 | grep "ld-linux-x86-64.so.2" | awk '{print $1}')
if [ -f "$ld_path" ]; then
    cp "$ld_path" "$BUNDLE_DIR/"
    echo "✓ Copied ld-linux-x86-64.so.2"
else
    echo "WARNING: Could not find ld-linux-x86-64.so.2"
fi

# Create wrapper script
cat > "$BUNDLE_DIR/notify-wrapper.sh" << 'EOF'
#!/bin/bash

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Set LD_LIBRARY_PATH to include bundled glibc
export LD_LIBRARY_PATH="$SCRIPT_DIR/lib:$LD_LIBRARY_PATH"

# Use the bundled dynamic linker if available
if [ -f "$SCRIPT_DIR/ld-linux-x86-64.so.2" ]; then
    exec "$SCRIPT_DIR/ld-linux-x86-64.so.2" "$SCRIPT_DIR/notify" "$@"
else
    exec "$SCRIPT_DIR/notify" "$@"
fi
EOF

chmod +x "$BUNDLE_DIR/notify-wrapper.sh"

# Create a simple test script
cat > "$BUNDLE_DIR/test-bundle.sh" << 'EOF'
#!/bin/bash

echo "=== Testing glibc bundle ==="
echo ""

echo "1. Testing notify version:"
./notify-wrapper.sh -version

echo ""
echo "2. Testing GUI check:"
./notify-wrapper.sh -check-gui

echo ""
echo "3. Testing wall check:"
./notify-wrapper.sh -check-wall

echo ""
echo "=== Bundle test complete ==="
EOF

chmod +x "$BUNDLE_DIR/test-bundle.sh"

echo ""
echo "=== Bundle creation complete ==="
echo ""
echo "Bundle directory: $BUNDLE_DIR"
echo "Contents:"
ls -la "$BUNDLE_DIR"

echo ""
echo "To use the bundle:"
echo "1. Copy the entire '$BUNDLE_DIR' directory to your target system"
echo "2. Run: ./notify-wrapper.sh instead of ./notify"
echo "3. Test with: ./test-bundle.sh"
echo ""
echo "WARNING: This approach may not work on all systems and has security implications."
echo "Consider this a last resort for very old systems."
