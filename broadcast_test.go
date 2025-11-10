package main

import (
	"runtime"
	"testing"
)

// TestWallAvailability tests the wall command availability check
func TestWallAvailability(t *testing.T) {
	result := isWallAvailable()

	if runtime.GOOS == "linux" {
		// On Linux, we expect wall to potentially be available
		t.Logf("Wall available on Linux: %v", result)
		if !result {
			t.Log("Note: wall command not found - this is OK for non-Linux or minimal Linux systems")
		}
	} else {
		// On non-Linux, wall should not be available
		if result {
			t.Errorf("Expected wall to be unavailable on %s, but it was reported as available", runtime.GOOS)
		}
		t.Logf("Wall correctly unavailable on %s", runtime.GOOS)
	}
}

// TestBroadcastWallMessage tests the wall broadcast functionality
func TestBroadcastWallMessage(t *testing.T) {
	// This test just verifies the function can be called
	// It won't actually send messages unless running on Linux with wall available

	err := broadcastWallMessage("Test Title", "Test Message", 0)

	if runtime.GOOS == "linux" && isWallAvailable() {
		// On Linux with wall, it might succeed (if we have permissions)
		// or fail (if we don't have permissions)
		t.Logf("Broadcast attempt on Linux: error=%v", err)
	} else {
		// On non-Linux, should return error
		if err == nil {
			t.Errorf("Expected error on non-Linux platform, got nil")
		}
		t.Logf("Broadcast correctly failed on %s: %v", runtime.GOOS, err)
	}
}

// TestBroadcastFallbackLogic tests the logic for choosing notification method
func TestBroadcastFallbackLogic(t *testing.T) {
	t.Log("Testing notification fallback hierarchy:")

	// Test GUI availability
	guiAvailable := isGUIAvailable()
	t.Logf("1. GUI Available: %v", guiAvailable)

	if !guiAvailable && runtime.GOOS == "linux" {
		// If no GUI on Linux, check wall
		wallAvailable := isWallAvailable()
		t.Logf("2. Wall Broadcast Available (Linux fallback): %v", wallAvailable)

		if !wallAvailable {
			t.Log("3. No notification method available")
		}
	}

	// Test OpenGL (for GUI environments)
	if guiAvailable {
		openglAvailable := isOpenGLAvailable()
		t.Logf("2. OpenGL Available: %v", openglAvailable)
	}
}
