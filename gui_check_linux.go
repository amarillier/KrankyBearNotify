//go:build linux

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// isLinuxGUIAvailable checks if GUI mode is available on Linux
// It checks for:
// 1. DISPLAY environment variable (X11) in current session
// 2. WAYLAND_DISPLAY environment variable (Wayland) in current session
// 3. Active graphical user sessions (even when run as root)
// 4. systemd graphical.target
func isLinuxGUIAvailable() bool {
	// Check for X11 display in current environment
	if display := os.Getenv("DISPLAY"); display != "" {
		return true
	}

	// Check for Wayland display in current environment
	if waylandDisplay := os.Getenv("WAYLAND_DISPLAY"); waylandDisplay != "" {
		return true
	}

	// Check if any user has an active graphical session (important for root)
	if hasActiveGraphicalSession() {
		return true
	}

	// Check systemd graphical.target as fallback
	return isGraphicalTargetActive()
}

// hasActiveGraphicalSession checks if any user has an active graphical session
// This is useful when running as root from SSH but users have GUI sessions
func hasActiveGraphicalSession() bool {
	// Check for active X11 sessions
	if hasX11Session() {
		return true
	}

	// Check for active Wayland sessions
	if hasWaylandSession() {
		return true
	}

	// Check loginctl for active graphical sessions
	return hasLoginctlGraphicalSession()
}

// hasX11Session checks if any X11 server is running
func hasX11Session() bool {
	// Check for X11 processes
	cmd := exec.Command("pgrep", "-x", "X")
	if err := cmd.Run(); err == nil {
		return true
	}

	// Also check for Xorg
	cmd = exec.Command("pgrep", "-x", "Xorg")
	if err := cmd.Run(); err == nil {
		return true
	}

	return false
}

// hasWaylandSession checks if any Wayland compositor is running
func hasWaylandSession() bool {
	// Common Wayland compositors
	compositors := []string{"weston", "sway", "mutter", "kwin_wayland", "gnome-shell"}

	for _, compositor := range compositors {
		cmd := exec.Command("pgrep", "-x", compositor)
		if err := cmd.Run(); err == nil {
			return true
		}
	}

	return false
}

// hasLoginctlGraphicalSession checks for graphical sessions using loginctl
func hasLoginctlGraphicalSession() bool {
	// Run loginctl list-sessions
	cmd := exec.Command("loginctl", "list-sessions", "--no-legend")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}

		sessionID := fields[0]

		// Get session type
		showCmd := exec.Command("loginctl", "show-session", sessionID, "-p", "Type", "--value")
		typeOutput, err := showCmd.Output()
		if err != nil {
			continue
		}

		sessionType := strings.TrimSpace(string(typeOutput))
		if sessionType == "x11" || sessionType == "wayland" {
			return true
		}
	}

	return false
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

// GraphicalSession represents a user's graphical session
type GraphicalSession struct {
	Username    string
	Display     string
	SessionID   string
	SessionType string // "x11" or "wayland"
}

// getGraphicalSessions returns all active graphical sessions
func getGraphicalSessions() []GraphicalSession {
	var sessions []GraphicalSession

	// Run loginctl list-sessions
	cmd := exec.Command("loginctl", "list-sessions", "--no-legend")
	output, err := cmd.Output()
	if err != nil {
		return sessions
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		sessionID := fields[0]
		username := fields[2]

		// Get session type
		typeCmd := exec.Command("loginctl", "show-session", sessionID, "-p", "Type", "--value")
		typeOutput, err := typeCmd.Output()
		if err != nil {
			continue
		}

		sessionType := strings.TrimSpace(string(typeOutput))
		if sessionType != "x11" && sessionType != "wayland" {
			continue
		}

		// Get display for this session
		display := getDisplayForSession(sessionID, username)
		if display == "" {
			continue
		}

		sessions = append(sessions, GraphicalSession{
			Username:    username,
			Display:     display,
			SessionID:   sessionID,
			SessionType: sessionType,
		})
	}

	return sessions
}

// getDisplayForSession gets the DISPLAY value for a specific session
func getDisplayForSession(sessionID, username string) string {
	// Try loginctl show-session to get Display property
	cmd := exec.Command("loginctl", "show-session", sessionID, "-p", "Display", "--value")
	output, err := cmd.Output()
	if err == nil {
		display := strings.TrimSpace(string(output))
		if display != "" {
			return display
		}
	}

	// Fallback: check process environment for X or Wayland compositors
	// Look for processes owned by the user
	pids := findUserGraphicalProcesses(username)
	for _, pid := range pids {
		display := getDisplayFromPID(pid)
		if display != "" {
			return display
		}
	}

	// Last resort: assume :0 if we know they have a graphical session
	return ":0"
}

// findUserGraphicalProcesses finds PIDs of graphical processes for a user
func findUserGraphicalProcesses(username string) []string {
	var pids []string

	// Look for common graphical processes
	processes := []string{"gnome-shell", "kwin_x11", "kwin_wayland", "xfce4-session", "cinnamon", "mate-session"}

	for _, proc := range processes {
		cmd := exec.Command("pgrep", "-u", username, "-x", proc)
		output, err := cmd.Output()
		if err == nil {
			pid := strings.TrimSpace(string(output))
			if pid != "" {
				pids = append(pids, strings.Split(pid, "\n")...)
			}
		}
	}

	return pids
}

// getDisplayFromPID extracts DISPLAY from a process's environment
func getDisplayFromPID(pid string) string {
	// Read /proc/PID/environ
	environFile := "/proc/" + pid + "/environ"
	data, err := os.ReadFile(environFile)
	if err != nil {
		return ""
	}

	// Parse null-separated environment variables
	envVars := strings.Split(string(data), "\x00")
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "DISPLAY=") {
			return strings.TrimPrefix(envVar, "DISPLAY=")
		}
		if strings.HasPrefix(envVar, "WAYLAND_DISPLAY=") {
			return strings.TrimPrefix(envVar, "WAYLAND_DISPLAY=")
		}
	}

	return ""
}

// shouldShowToOtherUsers determines if we should try to show GUI to other logged-in users
// This is true when running as root without our own DISPLAY access
func shouldShowToOtherUsers() bool {
	// Must be running as root
	if os.Geteuid() != 0 {
		return false
	}

	// Must not have our own DISPLAY
	if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" {
		return false
	}

	return true
}

// shouldUseWallBroadcast determines if wall broadcast should be used
// This is now only true if we're on SSH without GUI sessions at all
func shouldUseWallBroadcast() bool {
	// If we have DISPLAY or WAYLAND_DISPLAY, we can use GUI
	if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" {
		return false
	}

	// If running as root and there are GUI sessions available, don't use wall
	// (we'll show GUI to those users instead)
	if os.Geteuid() == 0 {
		sessions := getGraphicalSessions()
		if len(sessions) > 0 {
			return false // Don't use wall, we can show GUI to users
		}
		return true // No GUI sessions, use wall
	}

	// If SSH_CONNECTION is set and no DISPLAY, prefer wall
	if os.Getenv("SSH_CONNECTION") != "" || os.Getenv("SSH_CLIENT") != "" {
		return true
	}

	return false
}

// showNotificationToUsers shows GUI notifications to all users with active graphical sessions
// This is used when running as root to notify logged-in GUI users
func showNotificationToUsers(title, message string, timeout int, iconPath string, width, height int, buttonText string) error {
	sessions := getGraphicalSessions()
	if len(sessions) == 0 {
		return fmt.Errorf("no graphical sessions found")
	}

	var lastErr error
	successCount := 0

	for _, session := range sessions {
		err := showNotificationAsUser(session, title, message, timeout, iconPath, width, height, buttonText)
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

// showNotificationAsUser shows a notification as a specific user with their display
func showNotificationAsUser(session GraphicalSession, title, message string, timeout int, iconPath string, width, height int, buttonText string) error {
	// Get the path to the current executable
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Check and fix directory permissions in the path
	var restoreDirPerms []func()
	if os.Geteuid() == 0 {
		// Check all parent directories in the path
		exeDir := exePath
		for {
			exeDir = strings.TrimRight(exeDir, "/")
			exeDir = exeDir[:strings.LastIndex(exeDir, "/")+1]
			if exeDir == "" || exeDir == "/" {
				break
			}
			exeDir = strings.TrimRight(exeDir, "/")

			if dirInfo, err := os.Stat(exeDir); err == nil {
				dirMode := dirInfo.Mode()
				// Directories need r-x (0005) for users to traverse and read them
				if (dirMode.Perm() & 0005) != 0005 {
					originalPerm := dirMode.Perm()
					newPerm := originalPerm | 0555

					log.Printf("Note: Temporarily making directory accessible: %s (%o -> %o)\n",
						exeDir, originalPerm, newPerm)

					if err := os.Chmod(exeDir, newPerm); err != nil {
						log.Printf("Warning: Could not change directory permissions: %v\n", err)
					} else {
						// Capture the directory path for the closure
						capturedDir := exeDir
						capturedPerm := originalPerm
						restoreDirPerms = append(restoreDirPerms, func() {
							time.Sleep(time.Duration(timeout+2) * time.Second)
							if err := os.Chmod(capturedDir, capturedPerm); err != nil {
								log.Printf("Warning: Could not restore directory permissions: %v\n", err)
							} else {
								log.Printf("Note: Restored directory permissions: %s (%o)\n", capturedDir, capturedPerm)
							}
						})
					}
				}
			}
		}
	}

	// Check and fix executable file permissions if needed
	var restoreExePerms func()
	if exeInfo, err := os.Stat(exePath); err == nil && os.Geteuid() == 0 {
		exeMode := exeInfo.Mode()
		// Check if readable and executable by others (world r-x: 0005)
		// We need both read (0004) and execute (0001) for the file to be runnable
		needsPermFix := (exeMode.Perm() & 0005) != 0005

		if needsPermFix {
			originalPerm := exeMode.Perm()
			// Make readable and executable by all (add r-x for user, group, and others)
			newPerm := originalPerm | 0555 // Add read+execute for user, group, and others

			log.Printf("Note: Temporarily making executable accessible for user %s: %s (%o -> %o)\n",
				session.Username, exePath, originalPerm, newPerm)

			if err := os.Chmod(exePath, newPerm); err != nil {
				log.Printf("Warning: Could not change executable permissions: %v\n", err)
			} else {
				// Create a function to restore permissions later
				restoreExePerms = func() {
					time.Sleep(time.Duration(timeout+2) * time.Second)
					if err := os.Chmod(exePath, originalPerm); err != nil {
						log.Printf("Warning: Could not restore executable permissions: %v\n", err)
					} else {
						log.Printf("Note: Restored executable permissions: %s (%o)\n", exePath, originalPerm)
					}
				}
			}
		}
	}

	// Handle icon path and permissions
	finalIconPath := ""
	var restoreIconPerms func()

	if iconPath != "" {
		// Make sure the icon path is absolute
		absIconPath := iconPath
		if !strings.HasPrefix(iconPath, "/") {
			// If just a filename, look in the executable's directory
			// This handles cases where the icon is in the same dir as the binary
			exeDir := exePath[:strings.LastIndex(exePath, "/")]
			absIconPath = exeDir + "/" + iconPath
			log.Printf("Converted relative icon path '%s' to '%s'", iconPath, absIconPath)
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
						log.Printf("Found case-insensitive match: '%s' -> '%s'", iconPath, entry.Name())
						fileInfo, err = os.Stat(absIconPath)
						break
					}
				}
			}
		}

		if err == nil {
			// Check if the file is readable by others (or make it so temporarily)
			mode := fileInfo.Mode()
			needsPermFix := (mode.Perm() & 0004) == 0 // Check if world-readable

			if needsPermFix && os.Geteuid() == 0 {
				// We're root, temporarily make it readable
				// Save original permissions
				originalPerm := mode.Perm()

				fmt.Printf("Note: Temporarily making icon readable for user %s: %s (%o -> %o)\n",
					session.Username, absIconPath, originalPerm, originalPerm|0004)

				// Make readable by all (temporarily)
				if err := os.Chmod(absIconPath, mode.Perm()|0004); err != nil {
					fmt.Printf("Warning: Could not change icon permissions: %v\n", err)
				} else {
					// Create a function to restore permissions later
					// We'll call this after the notification timeout
					restoreIconPerms = func() {
						// Wait for notification to finish displaying (add buffer to timeout)
						time.Sleep(time.Duration(timeout+2) * time.Second)
						if err := os.Chmod(absIconPath, originalPerm); err != nil {
							fmt.Printf("Warning: Could not restore icon permissions: %v\n", err)
						} else {
							fmt.Printf("Note: Restored icon permissions: %s (%o)\n", absIconPath, originalPerm)
						}
					}
				}
			}

			finalIconPath = absIconPath
		} else {
			fmt.Printf("Warning: Icon file not accessible: %v\n", err)
		}
	}

	// Build the command arguments (after the environment vars)
	cmdArgs := []string{
		"-title", title,
		"-message", message,
		"-button", buttonText,
		"-timeout", fmt.Sprintf("%d", timeout),
		"-width", fmt.Sprintf("%d", width),
		"-height", fmt.Sprintf("%d", height),
	}

	// Add icon if we have a valid path
	if finalIconPath != "" {
		cmdArgs = append(cmdArgs, "-image", finalIconPath)
	}

	// Build sudo command with proper environment variable handling
	// Use 'env' to set environment variables for the child process
	args := []string{
		"-u", session.Username,
		"env",
		"DISPLAY=" + session.Display,
	}

	// Also set XAUTHORITY if we can find it
	xauth := findXauthorityForUser(session.Username)
	if xauth != "" {
		args = append(args, "XAUTHORITY="+xauth)
	}

	// Add the executable path
	args = append(args, exePath)

	// Add all the command arguments
	args = append(args, cmdArgs...)

	// Execute as the user (non-blocking, notification runs in background)
	cmd := exec.Command("sudo", args...)

	// Let stderr pass through so we can see any errors
	cmd.Stderr = os.Stderr

	err = cmd.Start() // Use Start() instead of Run() to not wait
	if err != nil {
		return fmt.Errorf("failed to run as user %s: %v", session.Username, err)
	}

	// Restore permissions after the notification timeout (in background)
	if restoreExePerms != nil {
		go restoreExePerms()
	}
	for _, restoreDir := range restoreDirPerms {
		go restoreDir()
	}
	if restoreIconPerms != nil {
		go restoreIconPerms()
	}

	return nil
}

// findXauthorityForUser tries to find the .Xauthority file for a user
func findXauthorityForUser(username string) string {
	// Try to get user's UID to check /run/user/<uid>
	var uid string
	cmd := exec.Command("id", "-u", username)
	if output, err := cmd.Output(); err == nil {
		uid = strings.TrimSpace(string(output))
	}

	// Try common locations
	possiblePaths := []string{
		"/home/" + username + "/.Xauthority",
	}

	// Add UID-specific path if we found it
	if uid != "" {
		possiblePaths = append(possiblePaths, "/run/user/"+uid+"/.Xauthority")
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// isMacGUIAvailable is a stub for non-Mac platforms
func isMacGUIAvailable() bool {
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

// hideConsoleWindow is a stub for non-Windows platforms
func hideConsoleWindow() {
	// No-op on Linux (no console window to hide)
}

// LinuxDistro represents a detected Linux distribution
type LinuxDistro struct {
	Name           string // "ubuntu", "debian", "fedora", "rhel", "centos", "arch", "opensuse", etc.
	Version        string
	PrettyName     string
	PackageManager string // "apt", "dnf", "yum", "pacman", "zypper"
}

// detectLinuxDistro detects the current Linux distribution
func detectLinuxDistro() LinuxDistro {
	distro := LinuxDistro{
		Name:           "unknown",
		PackageManager: "apt", // default fallback
	}

	// Read /etc/os-release (standard on systemd-based systems)
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		// Fallback to /etc/lsb-release
		data, err = os.ReadFile("/etc/lsb-release")
	}

	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "ID=") {
				distro.Name = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
			} else if strings.HasPrefix(line, "VERSION_ID=") {
				distro.Version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
			} else if strings.HasPrefix(line, "PRETTY_NAME=") {
				distro.PrettyName = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
			}
		}
	}

	// Determine package manager based on distro
	switch distro.Name {
	case "ubuntu", "debian", "linuxmint", "pop", "elementary":
		distro.PackageManager = "apt"
	case "fedora":
		distro.PackageManager = "dnf"
	case "rhel", "centos", "rocky", "almalinux":
		if distro.Version != "" {
			versionNum := 0
			fmt.Sscanf(distro.Version, "%d", &versionNum)
			if versionNum >= 8 {
				distro.PackageManager = "dnf"
			} else {
				distro.PackageManager = "yum"
			}
		} else {
			distro.PackageManager = "dnf" // assume newer
		}
	case "arch", "manjaro":
		distro.PackageManager = "pacman"
	case "opensuse", "opensuse-leap", "opensuse-tumbleweed", "sles":
		distro.PackageManager = "zypper"
	default:
		// Try to detect by available commands
		if _, err := exec.LookPath("apt-get"); err == nil {
			distro.PackageManager = "apt"
		} else if _, err := exec.LookPath("dnf"); err == nil {
			distro.PackageManager = "dnf"
		} else if _, err := exec.LookPath("yum"); err == nil {
			distro.PackageManager = "yum"
		} else if _, err := exec.LookPath("pacman"); err == nil {
			distro.PackageManager = "pacman"
		} else if _, err := exec.LookPath("zypper"); err == nil {
			distro.PackageManager = "zypper"
		}
	}

	return distro
}

// RequiredLibrary represents a shared library dependency
type RequiredLibrary struct {
	SoName      string // e.g., "libGL.so.1"
	DebPackage  string // apt package name
	RpmPackage  string // dnf/yum package name
	ArchPackage string // pacman package name
	SusePackage string // zypper package name
	Description string
}

// getRequiredLibraries returns the list of runtime libraries needed by notify
func getRequiredLibraries() []RequiredLibrary {
	return []RequiredLibrary{
		{
			SoName:      "libGL.so.1",
			DebPackage:  "libgl1",
			RpmPackage:  "mesa-libGL",
			ArchPackage: "mesa",
			SusePackage: "Mesa-libGL1",
			Description: "OpenGL library (required for GUI)",
		},
		{
			SoName:      "libXcursor.so.1",
			DebPackage:  "libxcursor1",
			RpmPackage:  "libXcursor",
			ArchPackage: "libxcursor",
			SusePackage: "libXcursor1",
			Description: "X11 cursor management",
		},
		{
			SoName:      "libXrandr.so.2",
			DebPackage:  "libxrandr2",
			RpmPackage:  "libXrandr",
			ArchPackage: "libxrandr",
			SusePackage: "libXrandr2",
			Description: "X11 screen resolution",
		},
		{
			SoName:      "libXinerama.so.1",
			DebPackage:  "libxinerama1",
			RpmPackage:  "libXinerama",
			ArchPackage: "libxinerama",
			SusePackage: "libXinerama1",
			Description: "X11 multi-screen support",
		},
		{
			SoName:      "libXi.so.6",
			DebPackage:  "libxi6",
			RpmPackage:  "libXi",
			ArchPackage: "libxi",
			SusePackage: "libXi6",
			Description: "X11 input extension",
		},
		{
			SoName:      "libXxf86vm.so.1",
			DebPackage:  "libxxf86vm1",
			RpmPackage:  "libXxf86vm",
			ArchPackage: "libxxf86vm",
			SusePackage: "libXxf86vm1",
			Description: "X11 video mode extension",
		},
	}
}

// checkLibraryAvailable checks if a shared library can be loaded
func checkLibraryAvailable(soName string) bool {
	// Try using ldconfig to check if library is available
	cmd := exec.Command("ldconfig", "-p")
	output, err := cmd.Output()
	if err == nil {
		return strings.Contains(string(output), soName)
	}

	// Fallback: try using find on common library directories
	commonPaths := []string{
		"/lib",
		"/lib64",
		"/usr/lib",
		"/usr/lib64",
		"/usr/lib/x86_64-linux-gnu",
		"/usr/lib/i386-linux-gnu",
	}

	for _, path := range commonPaths {
		testPath := path + "/" + soName
		if _, err := os.Stat(testPath); err == nil {
			return true
		}
		// Also check for symlinks with version numbers
		testPathStar := path + "/" + strings.Split(soName, ".so")[0] + ".so*"
		cmd := exec.Command("sh", "-c", "ls "+testPathStar+" 2>/dev/null")
		if output, err := cmd.Output(); err == nil && len(output) > 0 {
			return true
		}
	}

	return false
}

// checkDependencies checks for missing libraries and returns helpful info
func checkDependencies() (bool, []RequiredLibrary, LinuxDistro) {
	distro := detectLinuxDistro()
	required := getRequiredLibraries()
	var missing []RequiredLibrary

	for _, lib := range required {
		if !checkLibraryAvailable(lib.SoName) {
			missing = append(missing, lib)
		}
	}

	return len(missing) == 0, missing, distro
}

// getInstallCommand generates the appropriate install command for missing libraries
func getInstallCommand(missing []RequiredLibrary, distro LinuxDistro) string {
	if len(missing) == 0 {
		return ""
	}

	var packages []string
	var cmd string

	switch distro.PackageManager {
	case "apt":
		for _, lib := range missing {
			packages = append(packages, lib.DebPackage)
		}
		cmd = "sudo apt install -y " + strings.Join(packages, " ")

	case "dnf":
		for _, lib := range missing {
			packages = append(packages, lib.RpmPackage)
		}
		cmd = "sudo dnf install -y " + strings.Join(packages, " ")

	case "yum":
		for _, lib := range missing {
			packages = append(packages, lib.RpmPackage)
		}
		cmd = "sudo yum install -y " + strings.Join(packages, " ")

	case "pacman":
		for _, lib := range missing {
			packages = append(packages, lib.ArchPackage)
		}
		cmd = "sudo pacman -S --needed " + strings.Join(packages, " ")

	case "zypper":
		for _, lib := range missing {
			packages = append(packages, lib.SusePackage)
		}
		cmd = "sudo zypper install -y " + strings.Join(packages, " ")

	default:
		return "# Unknown package manager - please install the required libraries manually"
	}

	return cmd
}

// printDependencyReport prints a detailed dependency report
func printDependencyReport() {
	allOk, missing, distro := checkDependencies()

	fmt.Println("=== Dependency Check ===")
	fmt.Printf("Distribution: %s\n", distro.PrettyName)
	fmt.Printf("Package Manager: %s\n", distro.PackageManager)
	fmt.Println()

	if allOk {
		fmt.Println("- All required libraries are installed")
		fmt.Println("- GUI notifications should work properly")
		return
	}

	fmt.Println("- Missing required libraries:")
	fmt.Println()
	for _, lib := range missing {
		fmt.Printf("  ✗ %s - %s\n", lib.SoName, lib.Description)
	}
	fmt.Println()

	installCmd := getInstallCommand(missing, distro)
	fmt.Println("To install missing dependencies, run:")
	fmt.Println()
	fmt.Printf("  %s\n", installCmd)
	fmt.Println()
}

// checkLinuxDependencies runs a full dependency check and exits
func checkLinuxDependencies() {
	printDependencyReport()
	allOk, _, _ := checkDependencies()
	if allOk {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// checkLinuxDependenciesQuiet checks dependencies and prints a warning if any are missing
// This is called during -check-gui to provide helpful feedback
func checkLinuxDependenciesQuiet() {
	allOk, missing, distro := checkDependencies()
	if !allOk {
		fmt.Println()
		fmt.Println("Warning: Some runtime libraries are missing")
		fmt.Printf("Missing: ")
		libNames := []string{}
		for _, lib := range missing {
			libNames = append(libNames, lib.SoName)
		}
		fmt.Println(strings.Join(libNames, ", "))
		fmt.Println()
		installCmd := getInstallCommand(missing, distro)
		fmt.Printf("To fix: %s\n", installCmd)
		fmt.Println()
		fmt.Println("Run './notify -check-deps' for detailed information")
	}
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
