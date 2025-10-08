package main

import (
	"os"
	"runtime"
	"testing"
)

// testIsGUIAvailable tests the GUI availability check for the current platform
func testIsGUIAvailable(t *testing.T) {
	// This test will run on the current platform
	result := isGUIAvailable()

	// Log the result for informational purposes
	t.Logf("GUI available on %s: %v", runtime.GOOS, result)

	// We don't assert true/false here because the test might run in CI
	// without a GUI environment. This is more of an integration test.
}

// testLinuxGUIDetection tests Linux-specific GUI detection methods
func testLinuxGUIDetection(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping Linux-specific test on non-Linux platform")
	}

	// Test with DISPLAY environment variable
	t.Run("DISPLAY environment variable", func(t *testing.T) {
		originalDisplay := os.Getenv("DISPLAY")
		defer os.Setenv("DISPLAY", originalDisplay)

		// Test with DISPLAY set
		os.Setenv("DISPLAY", ":0")
		if !isLinuxGUIAvailable() {
			t.Log("GUI not detected with DISPLAY=:0 (might be running in headless environment)")
		}

		// Test with DISPLAY unset
		os.Unsetenv("DISPLAY")
		os.Unsetenv("WAYLAND_DISPLAY")
		result := isLinuxGUIAvailable()
		t.Logf("GUI available without DISPLAY/WAYLAND_DISPLAY: %v", result)
	})

	// Test with WAYLAND_DISPLAY environment variable
	t.Run("WAYLAND_DISPLAY environment variable", func(t *testing.T) {
		originalWayland := os.Getenv("WAYLAND_DISPLAY")
		defer os.Setenv("WAYLAND_DISPLAY", originalWayland)

		os.Setenv("WAYLAND_DISPLAY", "wayland-0")
		if !isLinuxGUIAvailable() {
			t.Log("GUI not detected with WAYLAND_DISPLAY=wayland-0 (might be running in headless environment)")
		}
	})

	// Test graphical.target detection
	t.Run("graphical.target detection", func(t *testing.T) {
		result := isGraphicalTargetActive()
		t.Logf("graphical.target active: %v", result)

		// This is informational - we can't assert because it depends on the system state
		if result {
			t.Log("systemd graphical.target is active")
		} else {
			t.Log("systemd graphical.target is not active or systemctl not available")
		}
	})
}

// TestMacGUIDetection tests macOS-specific GUI detection
func testMacGUIDetection(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping macOS-specific test on non-macOS platform")
	}

	result := isMacGUIAvailable()
	t.Logf("macOS GUI available (WindowServer running): %v", result)

	// On a real Mac with GUI, this should be true
	// In CI or headless environments, it might be false
}

// TestWindowsGUIDetection tests Windows-specific GUI detection
func testWindowsGUIDetection(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping Windows-specific test on non-Windows platform")
	}

	result := isWindowsGUIAvailable()
	t.Logf("Windows GUI available (window station accessible): %v", result)

	// On a real Windows system with GUI, this should be true
	// In headless environments, it might be false
}

// testCrossPlatformStubs tests that stub functions exist for non-native platforms
func testCrossPlatformStubs(t *testing.T) {
	// These should not panic regardless of platform
	switch runtime.GOOS {
	case "linux":
		// On Linux, Mac and Windows stubs should return false
		if isMacGUIAvailable() {
			t.Error("Mac stub should return false on Linux")
		}
		if isWindowsGUIAvailable() {
			t.Error("Windows stub should return false on Linux")
		}
	case "darwin":
		// On Mac, Linux and Windows stubs should return false
		if isLinuxGUIAvailable() {
			t.Error("Linux stub should return false on macOS")
		}
		if isWindowsGUIAvailable() {
			t.Error("Windows stub should return false on macOS")
		}
	case "windows":
		// On Windows, Linux and Mac stubs should return false
		if isLinuxGUIAvailable() {
			t.Error("Linux stub should return false on Windows")
		}
		if isMacGUIAvailable() {
			t.Error("Mac stub should return false on Windows")
		}
	}
}

// benchmarkIsGUIAvailable benchmarks the GUI availability check
func benchmarkIsGUIAvailable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isGUIAvailable()
	}
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
