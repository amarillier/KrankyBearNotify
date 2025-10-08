#!/bin/bash
# Example script demonstrating KrankyBear Notify usage

# Path to the krankybearnotify binary
NOTIFY_BIN="../krankybearnotify"

# Check if the binary exists
if [ ! -f "$NOTIFY_BIN" ]; then
    echo "Error: krankybearnotify binary not found. Please build it first with 'make build'"
    exit 1
fi

# Check if GUI is available
echo "Checking if GUI is available..."
if ! $NOTIFY_BIN -check-gui; then
    echo "GUI is not available. Cannot show notifications."
    exit 1
fi

echo "GUI is available. Showing example notifications..."
echo ""

# Example 1: Simple notification
echo "Example 1: Simple notification with defaults"
$NOTIFY_BIN &
sleep 1

# Example 2: Custom title and message with icon
echo "Example 2: Custom title and message with icon"
if [ -f "../KrankyBearBeret.png" ]; then
    $NOTIFY_BIN -title "Build Complete" -message "Your project has been built successfully!" -icon "../KrankyBearBeret.png" -timeout 5 &
else
    $NOTIFY_BIN -title "Build Complete" -message "Your project has been built successfully!" -timeout 5 &
fi
sleep 6

# Example 3: Warning notification with hard hat icon
echo "Example 3: Warning notification with hard hat icon"
if [ -f "../KrankyBearHardHat.png" ]; then
    $NOTIFY_BIN -title "⚠️ Warning" -message "Disk space is running low. Please free up some space." -icon "../KrankyBearHardHat.png" -timeout 7 &
else
    $NOTIFY_BIN -title "⚠️ Warning" -message "Disk space is running low. Please free up some space." -timeout 7 &
fi
sleep 8

# Example 4: Success notification with fedora icon
echo "Example 4: Success notification with fedora icon"
if [ -f "../KrankyBearFedoraRed.png" ]; then
    $NOTIFY_BIN -title "✓ Success" -message "All tests passed!" -icon "../KrankyBearFedoraRed.png" -timeout 4 &
else
    $NOTIFY_BIN -title "✓ Success" -message "All tests passed!" -timeout 4 &
fi
sleep 5

# Example 5: Manual close (no timeout)
echo "Example 5: Manual close notification (click OK to close)"
$NOTIFY_BIN -title "Manual Close" -message "This notification will stay until you click OK" -timeout 0

echo ""
echo "All examples completed!"
