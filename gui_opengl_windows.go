//go:build windows

package main

import (
	"log"
	"syscall"
	"unsafe"
)

// Windows constants for OpenGL context creation
const (
	WS_OVERLAPPEDWINDOW = 0x00CF0000
	CW_USEDEFAULT       = 0x80000000

	PFD_DRAW_TO_WINDOW = 0x00000004
	PFD_SUPPORT_OPENGL = 0x00000020
	PFD_DOUBLEBUFFER   = 0x00000001
	PFD_TYPE_RGBA      = 0
	PFD_MAIN_PLANE     = 0
)

// PIXELFORMATDESCRIPTOR structure (simplified)
type PIXELFORMATDESCRIPTOR struct {
	nSize           uint16
	nVersion        uint16
	dwFlags         uint32
	iPixelType      uint8
	cColorBits      uint8
	cRedBits        uint8
	cRedShift       uint8
	cGreenBits      uint8
	cGreenShift     uint8
	cBlueBits       uint8
	cBlueShift      uint8
	cAlphaBits      uint8
	cAlphaShift     uint8
	cAccumBits      uint8
	cAccumRedBits   uint8
	cAccumGreenBits uint8
	cAccumBlueBits  uint8
	cAccumAlphaBits uint8
	cDepthBits      uint8
	cStencilBits    uint8
	cAuxBuffers     uint8
	iLayerType      uint8
	bReserved       uint8
	dwLayerMask     uint32
	dwVisibleMask   uint32
	dwDamageMask    uint32
}

var (
	gdi32       = syscall.NewLazyDLL("gdi32.dll")
	opengl32Dll = syscall.NewLazyDLL("opengl32.dll")

	choosePixelFormat = gdi32.NewProc("ChoosePixelFormat")
	setPixelFormat    = gdi32.NewProc("SetPixelFormat")
	getDC             = user32.NewProc("GetDC")
	releaseDC         = user32.NewProc("ReleaseDC")
	wglCreateContext  = opengl32Dll.NewProc("wglCreateContext")
	wglDeleteContext  = opengl32Dll.NewProc("wglDeleteContext")
	wglMakeCurrent    = opengl32Dll.NewProc("wglMakeCurrent")
)

// isOpenGLAvailable checks if OpenGL is actually functional on Windows
// This is more robust than just checking if the DLL exists
func isOpenGLAvailable() bool {
	// First, basic check: can we load opengl32.dll?
	if err := opengl32Dll.Load(); err != nil {
		log.Printf("OpenGL check: opengl32.dll not found: %v", err)
		return false
	}

	// Check for wglCreateContext - core WGL function
	if err := wglCreateContext.Find(); err != nil {
		log.Printf("OpenGL check: wglCreateContext not found: %v", err)
		return false
	}

	// Try to get a device context from the desktop window
	// This tests if the graphics system is functional
	hdc, _, _ := getDC.Call(0) // 0 = desktop window
	if hdc == 0 {
		log.Println("OpenGL check: Failed to get device context")
		return false
	}
	defer releaseDC.Call(0, hdc)

	// Set up a minimal pixel format descriptor
	pfd := PIXELFORMATDESCRIPTOR{
		nSize:        uint16(unsafe.Sizeof(PIXELFORMATDESCRIPTOR{})),
		nVersion:     1,
		dwFlags:      PFD_DRAW_TO_WINDOW | PFD_SUPPORT_OPENGL | PFD_DOUBLEBUFFER,
		iPixelType:   PFD_TYPE_RGBA,
		cColorBits:   32,
		cDepthBits:   24,
		cStencilBits: 8,
		iLayerType:   PFD_MAIN_PLANE,
	}

	// Try to choose a pixel format
	pixelFormat, _, _ := choosePixelFormat.Call(hdc, uintptr(unsafe.Pointer(&pfd)))
	if pixelFormat == 0 {
		log.Println("OpenGL check: No suitable pixel format found (likely no OpenGL drivers)")
		return false
	}

	// Set the pixel format
	ret, _, _ := setPixelFormat.Call(hdc, pixelFormat, uintptr(unsafe.Pointer(&pfd)))
	if ret == 0 {
		log.Println("OpenGL check: Failed to set pixel format")
		return false
	}

	// NOW THE CRITICAL TEST: Try to actually create an OpenGL context
	hglrc, _, _ := wglCreateContext.Call(hdc)
	if hglrc == 0 {
		log.Println("OpenGL check: Failed to create OpenGL context (this is why Fyne fails in your VM!)")
		return false
	}
	defer wglDeleteContext.Call(hglrc)

	// Try to make the context current - final verification
	ret, _, _ = wglMakeCurrent.Call(hdc, hglrc)
	if ret == 0 {
		log.Println("OpenGL check: Failed to make OpenGL context current")
		return false
	}

	// Clean up - make no context current
	wglMakeCurrent.Call(hdc, 0)

	// If we got here, OpenGL is truly functional!
	log.Println("OpenGL check: OpenGL is fully functional and ready for Fyne")
	return true
}

// showWindowsMessageBox shows a native Windows MessageBox as fallback
func showWindowsMessageBox(title, message string, timeout int) error {
	// Get MessageBoxW from user32.dll (user32 is declared in gui_check_windows.go)
	messageBox := user32.NewProc("MessageBoxW")

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

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
