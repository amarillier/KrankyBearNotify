//go:build windows

package main

import (
	"syscall"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	getProcessWindowStation = user32.NewProc("GetProcessWindowStation")
)

// isWindowsGUIAvailable checks if GUI mode is available on Windows
// It checks if the process has access to a window station
func isWindowsGUIAvailable() bool {
	ret, _, _ := getProcessWindowStation.Call()
	return ret != 0
}

// isLinuxGUIAvailable is a stub for non-Linux platforms
func isLinuxGUIAvailable() bool {
	return false
}

// isMacGUIAvailable is a stub for non-Mac platforms
func isMacGUIAvailable() bool {
	return false
}

// isGraphicalTargetActive is a stub for non-Linux platforms
func isGraphicalTargetActive() bool {
	return false
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
