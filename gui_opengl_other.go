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

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
