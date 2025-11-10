//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// broadcastWallMessage sends a message to all logged-in users via wall command
// This is used when GUI is not available (headless, SSH, etc.)
func broadcastWallMessage(title, message string, timeout int) error {
	// Check if wall command is available
	_, err := exec.LookPath("wall")
	if err != nil {
		return fmt.Errorf("wall command not found: %v", err)
	}

	// Build the broadcast message
	var sb strings.Builder
	sb.WriteString("=" + strings.Repeat("=", 60) + "=\n")
	sb.WriteString(fmt.Sprintf("  %s\n", strings.ToUpper(title)))
	sb.WriteString("=" + strings.Repeat("=", 60) + "=\n\n")
	sb.WriteString(message)
	sb.WriteString("\n\n")
	if timeout > 0 {
		sb.WriteString(fmt.Sprintf("[This notification will be displayed for %d seconds]\n", timeout))
	}
	sb.WriteString("=" + strings.Repeat("=", 60) + "=\n")
	sb.WriteString(fmt.Sprintf("Sent: %s\n", time.Now().Format("2006-01-02 15:04:05")))

	broadcastMsg := sb.String()

	// Send the message via wall
	cmd := exec.Command("wall")
	cmd.Stdin = strings.NewReader(broadcastMsg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send wall broadcast: %v", err)
	}

	// If timeout is specified, wait and send a "notification expired" message
	if timeout > 0 {
		time.Sleep(time.Duration(timeout) * time.Second)

		expiryCmd := exec.Command("wall")
		expiryMsg := fmt.Sprintf("\n[Notification '%s' has expired]\n", title)
		expiryCmd.Stdin = strings.NewReader(expiryMsg)
		expiryCmd.Run() // Ignore errors on expiry message
	}

	return nil
}

// isWallAvailable checks if the wall command is available on this system
func isWallAvailable() bool {
	_, err := exec.LookPath("wall")
	return err == nil
}
