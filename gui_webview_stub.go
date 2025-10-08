//go:build !webview
// +build !webview

package main

import "fmt"

// showWebViewNotification stub when webview is not available
func showWebViewNotification(title, message string, timeout int, iconPath string) error {
	return fmt.Errorf("webview support not compiled in (use build tag: -tags webview)")
}

// isWebViewAvailable always returns false when webview is not compiled
func isWebViewAvailable() bool {
	return false
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
