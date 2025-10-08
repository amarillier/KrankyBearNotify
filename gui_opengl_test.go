package main

import (
	"runtime"
	"testing"
)

// TestOpenGLAvailability tests OpenGL detection
func TestOpenGLAvailability(t *testing.T) {
	result := isOpenGLAvailable()

	// On non-Windows platforms, should always return true
	if runtime.GOOS != "windows" {
		if !result {
			t.Error("OpenGL should be available on non-Windows platforms")
		}
	}

	// Log the result for informational purposes
	t.Logf("OpenGL available on %s: %v", runtime.GOOS, result)
}

// TestWindowsMessageBoxStub tests that the function exists
func TestWindowsMessageBoxStub(t *testing.T) {
	// This just verifies the function exists and doesn't panic
	// On non-Windows platforms, it should return nil
	err := showWindowsMessageBox("Test", "Test message", 0)

	if runtime.GOOS != "windows" {
		if err != nil {
			t.Errorf("Expected nil error on non-Windows, got: %v", err)
		}
	}

	t.Logf("Windows MessageBox stub works on %s", runtime.GOOS)
}

// TestOpenGLFallbackLogic tests the fallback decision logic
func TestOpenGLFallbackLogic(t *testing.T) {
	hasOpenGL := isOpenGLAvailable()

	if !hasOpenGL && runtime.GOOS == "windows" {
		t.Log("OpenGL not available on Windows - fallback to MessageBox would be used")
	} else if !hasOpenGL && runtime.GOOS != "windows" {
		t.Log("OpenGL not available on non-Windows platform - would fail")
	} else {
		t.Log("OpenGL available - Fyne GUI would be used")
	}
}
