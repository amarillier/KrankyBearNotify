//go:build darwin

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// isMacGUIAvailable checks if GUI mode is available on macOS
// On macOS, we check if the WindowServer is running
func isMacGUIAvailable() bool {
	// Check if WindowServer is running
	cmd := exec.Command("pgrep", "-x", "WindowServer")
	err := cmd.Run()
	return err == nil
}

// MacGUIUser represents a logged-in GUI user on macOS
type MacGUIUser struct {
	Username string
	UID      string
}

// getMacGUIUsers returns all users logged into the GUI
func getMacGUIUsers() []MacGUIUser {
	var users []MacGUIUser

	// Get console user (the one at the login screen/desktop)
	cmd := exec.Command("stat", "-f", "%Su", "/dev/console")
	output, err := cmd.Output()
	if err == nil {
		username := strings.TrimSpace(string(output))
		if username != "" && username != "root" {
			uid := getUIDForUser(username)
			users = append(users, MacGUIUser{
				Username: username,
				UID:      uid,
			})
		}
	}

	return users
}

// getUIDForUser gets the UID for a username
func getUIDForUser(username string) string {
	cmd := exec.Command("id", "-u", username)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// showNotificationToUsers shows notifications to all GUI users on macOS
func showNotificationToUsers(title, message string, timeout int, iconPath string, width, height int, buttonText string) error {
	users := getMacGUIUsers()
	if len(users) == 0 {
		return fmt.Errorf("no GUI users found")
	}

	var lastErr error
	successCount := 0

	for _, user := range users {
		err := showNotificationAsMacUser(user, title, message, timeout, iconPath, width, height, buttonText)
		if err != nil {
			lastErr = err
		} else {
			successCount++
		}
	}

	if successCount == 0 && lastErr != nil {
		return fmt.Errorf("failed to show notification to any user: %v", lastErr)
	}

	return nil
}

// showNotificationAsMacUser shows a notification as a specific macOS user
func showNotificationAsMacUser(user MacGUIUser, title, message string, timeout int, iconPath string, width, height int, buttonText string) error {
	// Get the path to the current executable
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Build the command to run as the user using launchctl asuser
	args := []string{
		"asuser",
		user.UID,
		exePath,
		"-title", title,
		"-message", message,
		"-button", buttonText,
		"-timeout", fmt.Sprintf("%d", timeout),
		"-width", fmt.Sprintf("%d", width),
		"-height", fmt.Sprintf("%d", height),
	}

	// Add icon if specified
	if iconPath != "" {
		// Make sure the icon path is absolute
		absIconPath := iconPath
		if !strings.HasPrefix(iconPath, "/") {
			// If just a filename, look in the executable's directory
			exeDir := exePath[:strings.LastIndex(exePath, "/")]
			absIconPath = exeDir + "/" + iconPath
		}

		// Check if the file exists and is readable
		fileInfo, err := os.Stat(absIconPath)
		if err != nil && !strings.HasPrefix(iconPath, "/") {
			// File not found and it's a relative path - try case-insensitive search
			exeDir := exePath[:strings.LastIndex(exePath, "/")]
			entries, readErr := os.ReadDir(exeDir)
			if readErr == nil {
				lowerIconPath := strings.ToLower(iconPath)
				for _, entry := range entries {
					if strings.ToLower(entry.Name()) == lowerIconPath {
						// Found a case-insensitive match
						absIconPath = exeDir + "/" + entry.Name()
						fileInfo, err = os.Stat(absIconPath)
						break
					}
				}
			}
		}

		if err == nil {
			mode := fileInfo.Mode()
			needsPermFix := (mode.Perm() & 0004) == 0

			if needsPermFix && os.Geteuid() == 0 {
				originalPerm := mode.Perm()
				os.Chmod(absIconPath, mode.Perm()|0004)
				defer os.Chmod(absIconPath, originalPerm)
			}
			args = append(args, "-image", absIconPath)
		}
	}

	// Execute using launchctl
	cmd := exec.Command("launchctl", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run as user %s: %v (output: %s)", user.Username, err, string(output))
	}

	return nil
}

// isLinuxGUIAvailable is a stub for non-Linux platforms
func isLinuxGUIAvailable() bool {
	return false
}

// isWindowsGUIAvailable is a stub for non-Windows platforms
func isWindowsGUIAvailable() bool {
	return false
}

// isRunningAsSystem is a stub for non-Windows platforms
func isRunningAsSystem() bool {
	return false
}

// isGraphicalTargetActive is a stub for non-Linux platforms
func isGraphicalTargetActive() bool {
	return false
}

// shouldShowToOtherUsers determines if we should try to show GUI to other logged-in users
// On macOS, this is true when running as root
func shouldShowToOtherUsers() bool {
	// Check if we're running as root (UID 0)
	return os.Geteuid() == 0
}

// shouldUseWallBroadcast is a stub for non-Linux platforms
func shouldUseWallBroadcast() bool {
	return false
}

// hideConsoleWindow is a stub for non-Windows platforms
func hideConsoleWindow() {
	// No-op on macOS (no console window to hide)
}

// checkLinuxDependencies is a stub for non-Linux platforms
func checkLinuxDependencies() {
	// No-op on macOS
}

// checkLinuxDependenciesQuiet is a stub for non-Linux platforms
func checkLinuxDependenciesQuiet() {
	// No-op on macOS
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
