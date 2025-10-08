//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32    = syscall.NewLazyDLL("kernel32.dll")
	loadLibrary = kernel32.NewProc("LoadLibraryW")
	freeLibrary = kernel32.NewProc("FreeLibrary")

	user32     = syscall.NewLazyDLL("user32.dll")
	messageBox = user32.NewProc("MessageBoxW")
)

// isOpenGLAvailable checks if OpenGL is available on Windows
func isOpenGLAvailable() bool {
	// Try to load opengl32.dll
	dll, _, _ := syscall.NewLazyDLL("opengl32.dll").NewProc("wglGetProcAddress").Find()
	return dll != 0
}

// showWindowsMessageBox shows a native Windows MessageBox as fallback
func showWindowsMessageBox(title, message string, timeout int) error {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	messagePtr, _ := syscall.UTF16PtrFromString(message)

	// MB_OK | MB_ICONINFORMATION | MB_TOPMOST
	const MB_OK = 0x00000000
	const MB_ICONINFORMATION = 0x00000040
	const MB_TOPMOST = 0x00040000

	flags := MB_OK | MB_ICONINFORMATION | MB_TOPMOST

	if timeout > 0 {
		// For timeout, we'd need to use a timer and close the window
		// For simplicity, we'll just show the message
		messageWithTimeout, _ := syscall.UTF16PtrFromString(message + "\n\n(Auto-close not supported in fallback mode)")
		messageBox.Call(
			0,
			uintptr(unsafe.Pointer(messageWithTimeout)),
			uintptr(unsafe.Pointer(titlePtr)),
			uintptr(flags),
		)
	} else {
		messageBox.Call(
			0,
			uintptr(unsafe.Pointer(messagePtr)),
			uintptr(unsafe.Pointer(titlePtr)),
			uintptr(flags),
		)
	}

	return nil
}
