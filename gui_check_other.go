//go:build !linux && !darwin && !windows

package main

import "fmt"

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

// shouldShowToOtherUsers is a stub for unsupported platforms
func shouldShowToOtherUsers() bool {
	return false
}

// isRunningAsSystem is a stub for unsupported platforms
func isRunningAsSystem() bool {
	return false
}

// shouldUseWallBroadcast is a stub for non-Linux platforms
func shouldUseWallBroadcast() bool {
	return false
}

// showNotificationToUsers is a stub for unsupported platforms
func showNotificationToUsers(title, message string, timeout int, iconPath string, width, height int, buttonText string) error {
	return fmt.Errorf("showNotificationToUsers is not supported on this platform")
}

// hideConsoleWindow is a stub for non-Windows platforms
func hideConsoleWindow() {
	// No-op on unsupported platforms
}

// checkLinuxDependencies is a stub for non-Linux platforms
func checkLinuxDependencies() {
	// No-op on other platforms
}

// checkLinuxDependenciesQuiet is a stub for non-Linux platforms
func checkLinuxDependenciesQuiet() {
	// No-op on other platforms
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
