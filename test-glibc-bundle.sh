#!/bin/bash
# Test the glibc bundle approach

set -e

echo "=== Testing glibc bundle approach ==="
echo ""

# First, ensure we have a Linux binary
if [ ! -f "bin/notify-linux-amd64" ]; then
    echo "Building Linux binary first..."
    ./compile-linux.sh
fi

echo "Current system glibc version:"
ldd --version | head -n1

echo ""
echo "Dependencies of notify binary:"
ldd bin/notify-linux-amd64 | head -10

echo ""
echo "Creating glibc bundle..."
./create-glibc-bundle-advanced.sh

echo ""
echo "Testing the bundle..."
cd notify-with-glibc-*
./test-bundle.sh

echo ""
echo "=== Bundle test complete ==="
echo "If this works, you can deploy the bundle to older systems"
