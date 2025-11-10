//go:build !linux

package main

import "fmt"

// broadcastWallMessage is a stub for non-Linux platforms
func broadcastWallMessage(title, message string, timeout int) error {
	return fmt.Errorf("wall broadcast is only available on Linux")
}

// isWallAvailable is a stub for non-Linux platforms
func isWallAvailable() bool {
	return false
}
