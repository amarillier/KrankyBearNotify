//go:build !linux && !darwin && !windows

package main

// isLinuxGUIAvailable is a stub for non-Linux platforms
func isLinuxGUIAvailable() bool {
	return false
}

// isMacGUIAvailable is a stub for non-Mac platforms
func isMacGUIAvailable() bool {
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
