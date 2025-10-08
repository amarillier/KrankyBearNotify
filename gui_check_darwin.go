//go:build darwin

package main

import (
	"os/exec"
)

// isMacGUIAvailable checks if GUI mode is available on macOS
// On macOS, we check if the WindowServer is running
func isMacGUIAvailable() bool {
	// Check if WindowServer is running
	cmd := exec.Command("pgrep", "-x", "WindowServer")
	err := cmd.Run()
	return err == nil
}

// isLinuxGUIAvailable is a stub for non-Linux platforms
func isLinuxGUIAvailable() bool {
	return false
}

// isWindowsGUIAvailable is a stub for non-Windows platforms
func isWindowsGUIAvailable() bool {
	return false
}

// isGraphicalTargetActive is a stub for non-Linux platforms
func isGraphicalTargetActive() bool {
	return false
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
