//go:build !windows

package main

// isOpenGLAvailable always returns true on non-Windows platforms
// (macOS and Linux handle OpenGL differently and Fyne works well on them)
func isOpenGLAvailable() bool {
	return true
}

// showWindowsMessageBox is not available on non-Windows platforms
func showWindowsMessageBox(title, message string, timeout int) error {
	return nil
}
