#!/bin/bash
# Bundle glibc with the notify application for older Linux systems
# WARNING: This is complex and may not work on all systems

set -e

echo "=== Bundling glibc with notify application ==="
echo "WARNING: This approach has limitations and may not work on all systems"
echo ""

# Create bundle directory
BUNDLE_DIR="notify-with-glibc"
mkdir -p "$BUNDLE_DIR"

# Copy the notify binary
cp bin/notify-linux-amd64 "$BUNDLE_DIR/notify"

# Create a wrapper script that sets LD_LIBRARY_PATH
cat > "$BUNDLE_DIR/notify-wrapper.sh" << 'EOF'
#!/bin/bash

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Set LD_LIBRARY_PATH to include bundled glibc
export LD_LIBRARY_PATH="$SCRIPT_DIR/lib:$LD_LIBRARY_PATH"

# Run the actual notify binary
exec "$SCRIPT_DIR/notify" "$@"
EOF

chmod +x "$BUNDLE_DIR/notify-wrapper.sh"

# Create lib directory for glibc
mkdir -p "$BUNDLE_DIR/lib"

echo "To complete the glibc bundling, you need to:"
echo "1. Copy the required glibc libraries to $BUNDLE_DIR/lib/"
echo "2. Copy the ld-linux-x86-64.so.2 loader to $BUNDLE_DIR/"
echo ""
echo "Required files (from a newer Linux system):"
echo "  - libc.so.6"
echo "  - libm.so.6" 
echo "  - libpthread.so.0"
echo "  - libdl.so.2"
echo "  - libcrypt.so.1"
echo "  - libnsl.so.1"
echo "  - libresolv.so.2"
echo "  - ld-linux-x86-64.so.2"
echo ""
echo "You can find these with: ldd bin/notify-linux-amd64"
echo ""
echo "Then users would run: ./notify-wrapper.sh instead of ./notify"
