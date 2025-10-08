//go:build linux

package main

import (
	"os"
	"os/exec"
	"strings"
)

// isLinuxGUIAvailable checks if GUI mode is available on Linux
// It checks for:
// 1. DISPLAY environment variable (X11)
// 2. WAYLAND_DISPLAY environment variable (Wayland)
// 3. systemd graphical.target
func isLinuxGUIAvailable() bool {
	// Check for X11 display
	if display := os.Getenv("DISPLAY"); display != "" {
		return true
	}

	// Check for Wayland display
	if waylandDisplay := os.Getenv("WAYLAND_DISPLAY"); waylandDisplay != "" {
		return true
	}

	// Check systemd graphical.target
	return isGraphicalTargetActive()
}

// isGraphicalTargetActive checks if systemd graphical.target is active
func isGraphicalTargetActive() bool {
	cmd := exec.Command("systemctl", "is-active", "graphical.target")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	status := strings.TrimSpace(string(output))
	return status == "active"
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
