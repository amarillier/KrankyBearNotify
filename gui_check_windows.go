//go:build windows

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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

// WindowsGUIUser represents a logged-in GUI user on Windows
type WindowsGUIUser struct {
	Username  string
	SessionID string
}

// getWindowsGUIUsers returns all users with active GUI sessions
func getWindowsGUIUsers() []WindowsGUIUser {
	var users []WindowsGUIUser

	// Use query user command (quser/query user)
	// Try quser first (more concise output)
	cmd := exec.Command("quser")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to query user
		cmd = exec.Command("query", "user")
		output, err = cmd.Output()
		if err != nil {
			return users
		}
	}

	lines := strings.Split(string(output), "\n")
	// Skip header line
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		// Parse quser output
		// Format: USERNAME SESSIONNAME ID STATE IDLE TIME LOGON TIME
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		username := fields[0]
		// If username starts with >, it's the current user, strip it
		username = strings.TrimPrefix(username, ">")

		// Session ID is typically field 2 (after username and session name)
		// But if session name is missing (e.g., console), it shifts
		sessionID := ""
		if len(fields) >= 2 {
			// Try to find the numeric session ID
			for _, field := range fields[1:] {
				// Session IDs are typically numeric
				if field != "" && (field[0] >= '0' && field[0] <= '9') {
					sessionID = field
					break
				}
			}
		}

		if sessionID != "" {
			users = append(users, WindowsGUIUser{
				Username:  username,
				SessionID: sessionID,
			})
		}
	}

	return users
}

// showNotificationToUsers shows notifications to all GUI users on Windows
func showNotificationToUsers(title, message string, timeout int, iconPath string, width, height int, buttonText string) error {
	users := getWindowsGUIUsers()
	if len(users) == 0 {
		return fmt.Errorf("no GUI users found")
	}

	var lastErr error
	successCount := 0

	for _, user := range users {
		err := showNotificationAsWindowsUser(user, title, message, timeout, iconPath, width, height, buttonText)
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

// showNotificationAsWindowsUser shows a notification to a specific Windows user
func showNotificationAsWindowsUser(user WindowsGUIUser, title, message string, timeout int, iconPath string, width, height int, buttonText string) error {
	// Get the path to the current executable
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	// Build the command arguments
	// Let the child process auto-detect the best GUI mode, or pass through forced mode flags
	// (it will run as the target user, so Fyne/WebView should work)
	args := []string{}

	// CRITICAL: Add -target-user flag to prevent infinite loop
	args = append(args, "-target-user")
	log.Println("Adding -target-user flag to prevent re-elevation")

	// Pass through mode flags and feature flags if they were specified
	// Check os.Args to see what the parent was called with
	passedFlags := []string{}
	for _, arg := range os.Args {
		// Pass through mode flags, autosize flag, and debug flag
		if arg == "-win-webview" || arg == "-win-basic" || arg == "-autosize" || arg == "-debug" {
			args = append(args, arg)
			passedFlags = append(passedFlags, arg)
		}
	}
	if len(passedFlags) > 0 {
		log.Printf("Passing flags to child process: %v", passedFlags)
	} else {
		log.Printf("No special flags detected in os.Args: %v", os.Args)
	}

	// Add notification parameters
	args = append(args, "-title", title)
	args = append(args, "-message", message)
	args = append(args, "-button", buttonText)
	args = append(args, "-timeout", fmt.Sprintf("%d", timeout))
	args = append(args, "-width", fmt.Sprintf("%d", width))
	args = append(args, "-height", fmt.Sprintf("%d", height))

	// Add icon if specified
	if iconPath != "" {
		// Ensure absolute path for Windows
		absIconPath := iconPath
		if !strings.Contains(iconPath, ":") && !strings.HasPrefix(iconPath, "\\\\") {
			// Use executable directory as base, not working directory
			// This ensures the icon path is correct when launched as another user
			exeDir := exePath
			if lastSlash := strings.LastIndex(exeDir, "\\"); lastSlash > 0 {
				exeDir = exeDir[:lastSlash]
			}
			absIconPath = exeDir + "\\" + iconPath
		}

		// Verify file exists before passing it
		if _, err := os.Stat(absIconPath); err == nil {
			args = append(args, "-image", absIconPath)
			log.Printf("Including icon in child process args: %s", absIconPath)
		} else {
			log.Printf("Icon file not found, skipping: %s (error: %v)", absIconPath, err)
		}
	}

	// Build command string for PsExec or PowerShell
	cmdStr := fmt.Sprintf("\"%s\"", exePath)
	for _, arg := range args {
		// Escape quotes in arguments
		escapedArg := strings.ReplaceAll(arg, "\"", "\\\"")
		cmdStr += fmt.Sprintf(" \"%s\"", escapedArg)
	}

	// Try PsExec first if available (more reliable)
	// Check if PsExec is available
	psExecPath := ""
	for _, path := range []string{"psexec.exe", "psexec64.exe", "C:\\Windows\\System32\\PsExec64.exe", "C:\\Windows\\SysWOW64\\PsExec.exe"} {
		if _, err := exec.LookPath(path); err == nil {
			psExecPath = path
			break
		}
	}

	if psExecPath != "" {
		log.Printf("Using PsExec to launch notification for user %s in session %s", user.Username, user.SessionID)

		psExecArgs := []string{
			"-accepteula",
			"-nobanner",
			"-i", user.SessionID,
			"-d", // Don't wait for process to terminate
			exePath,
		}
		psExecArgs = append(psExecArgs, args...)

		cmd := exec.Command(psExecPath, psExecArgs...)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000, // CREATE_NO_WINDOW
		}

		output, err := cmd.CombinedOutput()
		if err == nil {
			log.Printf("Successfully launched via PsExec for user %s", user.Username)
			return nil
		}
		log.Printf("PsExec failed: %v (output: %s), falling back to scheduled task", err, string(output))
	}

	// Fallback: Use PowerShell with task scheduler
	taskName := fmt.Sprintf("KrankyBearNotify_%s_%d", user.Username, timeout)

	// Build argument string with proper PowerShell escaping
	// We need to build a single string that will be passed to -Argument parameter
	var argParts []string
	for _, arg := range args {
		// For the -Argument parameter, we need to escape double quotes
		// and wrap each argument in double quotes
		escaped := strings.ReplaceAll(arg, `"`, `\"`)
		argParts = append(argParts, fmt.Sprintf(`"%s"`, escaped))
	}
	// Join with spaces to create the full argument string
	argString := strings.Join(argParts, " ")

	// Escape the argument string for PowerShell single-quoted string
	escapedArgString := strings.ReplaceAll(argString, "'", "''")

	// Escape the executable path for PowerShell
	escapedExePath := strings.ReplaceAll(exePath, "'", "''")

	// Escape the username for PowerShell
	escapedUsername := strings.ReplaceAll(user.Username, "'", "''")

	// PowerShell script with better error handling
	// Use here-strings and proper variable expansion to avoid quoting issues
	psScript := fmt.Sprintf(`
$ErrorActionPreference = 'Stop'
try {
    # Clean up any existing task with same name
    Get-ScheduledTask -TaskName '%s' -ErrorAction SilentlyContinue | Unregister-ScheduledTask -Confirm:$false
    
    # Build the action with the executable and arguments
    $exe = '%s'
    $arguments = '%s'
    $action = New-ScheduledTaskAction -Execute $exe -Argument $arguments
    
    # Settings for immediate execution
    # Multiple settings to ensure task runs without visible console
    $settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -DontStopOnIdleEnd -StartWhenAvailable -ExecutionTimeLimit (New-TimeSpan -Minutes 5)
    
    # Try to prevent console windows - set the task to run hidden
    # This is a best-effort attempt as Task Scheduler has limitations
    $settings.Priority = 4  # Normal priority
    
    # Trigger to run once immediately
    $trigger = New-ScheduledTaskTrigger -Once -At (Get-Date)
    
    # Get the fully qualified username (handles domain vs local users)
    $username = '%s'
    $userPrincipal = $username
    if ($username -notlike '*\*') {
        # If username doesn't contain backslash, it's likely a local user
        # Try to get the computer name and prefix it
        try {
            $computerName = $env:COMPUTERNAME
            $userPrincipal = "$computerName\$username"
        } catch {
            # Fallback to .\username for local users
            $userPrincipal = ".\$username"
        }
    }
    
    # Principal to run as the target user with highest privileges
    # Must use Interactive for GUI access, but we hide console via notify.exe itself
    $principal = New-ScheduledTaskPrincipal -UserId $userPrincipal -LogonType Interactive -RunLevel Highest
    
    # Register the task (suppress output)
    $task = Register-ScheduledTask -TaskName '%s' -Action $action -Settings $settings -Trigger $trigger -Principal $principal -Force | Out-Null
    
    if (-not (Get-ScheduledTask -TaskName '%s' -ErrorAction SilentlyContinue)) {
        Write-Error 'Failed to register task'
        exit 1
    }
    
    # Start the task (suppress output)
    Start-ScheduledTask -TaskName '%s' | Out-Null
    
    # Wait a moment for task to start, then clean up in background
    Start-Sleep -Milliseconds 500
    
    # Clean up scheduled task (suppress output)
    Unregister-ScheduledTask -TaskName '%s' -Confirm:$false -ErrorAction SilentlyContinue | Out-Null
    
    exit 0
} catch {
    Write-Host "ERROR: $_"
    exit 1
}
`, taskName, escapedExePath, escapedArgString, escapedUsername, taskName, taskName, taskName, taskName)

	log.Printf("Attempting scheduled task launch for user %s in session %s", user.Username, user.SessionID)

	// Run PowerShell completely hidden (no window at all)
	cmd := exec.Command("powershell.exe",
		"-WindowStyle", "Hidden",
		"-NoProfile",
		"-NonInteractive",
		"-NoLogo",
		"-ExecutionPolicy", "Bypass",
		"-Command", psScript)

	// Hide the PowerShell window completely
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000 | 0x00000010, // CREATE_NO_WINDOW | CREATE_NEW_CONSOLE (then hide it)
	}

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		log.Printf("PowerShell error for user %s: %v\nOutput: %s", user.Username, err, outputStr)
		return fmt.Errorf("failed to run as user %s: %v (output: %s)", user.Username, err, outputStr)
	}

	// Check for errors in output
	if strings.Contains(outputStr, "ERROR:") {
		log.Printf("Scheduled task creation had errors for user %s: %s", user.Username, outputStr)
		return fmt.Errorf("scheduled task creation failed for user %s: %s", user.Username, outputStr)
	}

	log.Printf("Successfully created and started scheduled task for user %s", user.Username)
	log.Printf("Child process command: %s %v", exePath, args)

	return nil
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

// isRunningAsSystem checks if we're running as SYSTEM account on Windows
func isRunningAsSystem() bool {
	cmd := exec.Command("whoami")
	output, err := cmd.Output()
	if err == nil {
		username := strings.TrimSpace(strings.ToLower(string(output)))
		return strings.Contains(username, "system")
	}
	return false
}

// shouldShowToOtherUsers determines if we should try to show GUI to other logged-in users
// On Windows, check if we're running as SYSTEM or elevated Administrator
func shouldShowToOtherUsers() bool {
	// CRITICAL: If we were launched as a target user from an elevated parent,
	// DO NOT try to elevate again (prevents infinite loop)
	// Check both environment variable (old method) and command-line flag (new method)
	if os.Getenv("NOTIFY_TARGET_USER") == "1" {
		log.Println("Running as target user (NOTIFY_TARGET_USER=1), will not elevate")
		return false
	}

	// Check for -target-user flag in os.Args (more reliable than env var)
	for _, arg := range os.Args {
		if arg == "-target-user" {
			log.Println("Running as target user (-target-user flag), will not elevate")
			return false
		}
	}

	// Check if running as SYSTEM
	if isRunningAsSystem() {
		return true
	}

	// Check if we're running elevated (as Administrator)
	// Try to open a privileged registry key
	cmd := exec.Command("net", "session")
	err := cmd.Run()
	// If this succeeds, we're running elevated
	return err == nil
}

// shouldUseWallBroadcast is a stub for non-Linux platforms
func shouldUseWallBroadcast() bool {
	return false
}

// hideConsoleWindow hides the console window using Windows API
// This is called when running as target user to prevent console from showing
func hideConsoleWindow() {
	// Get handle to kernel32.dll
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")

	// Get handle to user32.dll
	user32 := syscall.NewLazyDLL("user32.dll")
	showWindow := user32.NewProc("ShowWindow")

	// Get console window handle
	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return // No console window
	}

	// SW_HIDE = 0
	const SW_HIDE = 0

	// Hide the console window
	showWindow.Call(hwnd, SW_HIDE)
	log.Println("Console window hidden via Windows API")
}

// checkLinuxDependencies is a stub for non-Linux platforms
func checkLinuxDependencies() {
	// No-op on Windows
}

// checkLinuxDependenciesQuiet is a stub for non-Linux platforms
func checkLinuxDependenciesQuiet() {
	// No-op on Windows
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
