package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	updatechecker "github.com/amarillier/go-update-checker"
)

const (
	appVersion     = "0.1.2"
	appAuthor      = "Allan Marillier"
	defaultTitle   = "Notification"
	defaultMessage = "This is a notification message"
	defaultTimeout = 10  // seconds
	defaultWidth   = 400 // pixels
	defaultHeight  = 250 // pixels
)

var appCopyright = "Copyright (c) Allan Marillier, 2024-" + strconv.Itoa(time.Now().Year())

func init() {
	// Custom usage function for better help output
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `KrankyBear Notify v%s
A cross-platform notification application for Mac, Windows, and Linux

USAGE:
  %s [OPTIONS]

OPTIONS:
`, appVersion, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
EXAMPLES:
  # Show a simple notification
  %s -title "Hello" -message "World!"

  # Notification with custom icon and timeout
  %s -title "Build Complete" -message "Success!" -icon "./icon.png" -timeout 5

  # Check if GUI is available (useful for scripts)
  %s -check-gui

  # Check for updates
  %s -cu

  # Notification that stays until manually closed
  %s -title "Important" -message "Please review" -timeout 0

  # Force basic mode (for VMs where OpenGL detection passes but Fyne fails)
  %s -force-basic -title "VM Alert" -message "Uses MessageBox"

  # Force WebView mode (better UI than MessageBox, requires webview build)
  %s -force-webview -title "Modern Alert" -message "Uses HTML/CSS/JS"

SUPPORTED PLATFORMS:
  • macOS 10.13+
  • Windows 10+
  • Linux (X11/Wayland) - Works on GNOME, KDE, XFCE, Cinnamon, MATE, and more
    - Headless/SSH: Falls back to 'wall' broadcast when no GUI detected

For more information, visit: https://github.com/amarillier/krankybearnotify
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	}
}

func main() {
	// CRITICAL: Handle version flag BEFORE any other code runs
	// This prevents Fyne GUI initialization which can hang in some environments

	// Quick pre-check for version flag to avoid GUI initialization
	for _, arg := range os.Args[1:] {
		if arg == "-version" || arg == "--version" {
			fmt.Printf("KrankyBear Notify v%s\n", appVersion)
			fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
			fmt.Printf("Copyright: %s\n", appCopyright)
			fmt.Println("License: GNU GPL-3.0")
			fmt.Println("Source: https://github.com/amarillier/krankybearnotify")
			fmt.Println("Documentation: https://github.com/amarillier/krankybearnotify/blob/main/README.md")
			os.Exit(0)
		}
	}

	// Check for help flags - we need to define flags first before showing usage
	// so we check here but display help after flag definitions
	showHelp := false
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "-help" || arg == "--help" || arg == "-?" {
			showHelp = true
			break
		}
	}

	// Show help if no arguments provided at all
	if len(os.Args) == 1 {
		showHelp = true
	}

	// Command-line flags
	title := flag.String("title", defaultTitle, "Notification title")
	message := flag.String("message", defaultMessage, "Notification message")
	timeout := flag.Int("timeout", defaultTimeout, "Timeout in seconds (0 for no timeout)")
	width := flag.Int("width", defaultWidth, "Window width in pixels")
	height := flag.Int("height", defaultHeight, "Window height in pixels")
	checkGUI := flag.Bool("check-gui", false, "Check if GUI mode is available and exit")
	checkOpenGL := flag.Bool("check-opengl", false, "Check if OpenGL is available and exit")
	checkWall := flag.Bool("check-wall", false, "Check if wall broadcast is available (Linux) and exit")
	forceBasic := flag.Bool("force-basic", false, "Force basic GUI mode (skip OpenGL/Fyne, use fallback)")
	forceWebView := flag.Bool("force-webview", false, "Force WebView mode (requires -tags webview build)")
	version := flag.Bool("version", false, "Show version information and exit")

	// Icon flag with alias
	var icon string
	flag.StringVar(&icon, "icon", "", "Path to icon image file (PNG, JPEG, etc.)")
	flag.StringVar(&icon, "image", "", "Path to icon image file (alias for -icon)")

	// Update checker flags (with alias)
	var checkUpdate bool
	flag.BoolVar(&checkUpdate, "checkupdate", false, "Check for updates and exit")
	flag.BoolVar(&checkUpdate, "cu", false, "Check for updates and exit (alias for -checkupdate)")

	// Now show help if requested (flags are defined, so PrintDefaults will work)
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Parse command-line flags (help/version already handled above)
	flag.Parse()

	// Show version if requested
	if *version {
		fmt.Printf("KrankyBear Notify v%s\n", appVersion)
		fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("Copyright: %s\n", appCopyright)
		fmt.Println("License: MIT")
		fmt.Println("Source: https://github.com/amarillier/krankybearnotify")
		fmt.Println("Documentation: https://github.com/amarillier/krankybearnotify/blob/main/README.md")
		os.Exit(0)
	}

	// Check for updates if requested
	if checkUpdate {
		fmt.Printf("Checking for updates...\n")
		fmt.Printf("Current version: %s\n\n", appVersion)

		// Get executable directory for storing update check file
		exePath, err := os.Executable()
		if err != nil {
			log.Printf("Warning: Could not determine executable path: %v", err)
			exePath = "."
		}
		exeDir := filepath.Dir(exePath)
		checkFilePath := filepath.Join(exeDir, "latestcheck.json")

		// Save current directory and change to executable directory
		originalDir, _ := os.Getwd()
		os.Chdir(exeDir)
		defer os.Chdir(originalDir)

		updtmsg, updateAvailable := updateChecker("amarillier", "KrankyBearNotify", "Kranky Bear Notify", "https://github.com/amarillier/KrankyBearNotify/releases/latest")

		if updateAvailable {
			fmt.Println(updtmsg)
			fmt.Printf("\nUpdate check data saved to: %s\n", checkFilePath)
			os.Exit(0)
		} else {
			fmt.Println("You are running the latest version!")
			fmt.Printf("Update check data saved to: %s\n", checkFilePath)
			os.Exit(0)
		}
	}

	// Check GUI mode if requested
	if *checkGUI {
		if isGUIAvailable() {
			fmt.Println("GUI mode is available")
			os.Exit(0)
		} else {
			fmt.Println("GUI mode is not available")
			os.Exit(1)
		}
	}

	// Check OpenGL if requested
	if *checkOpenGL {
		if isOpenGLAvailable() {
			fmt.Println("OpenGL is available")
			fmt.Println("Fyne GUI can be used")
			os.Exit(0)
		} else {
			fmt.Println("OpenGL is not available")
			if runtime.GOOS == "windows" {
				fmt.Println("Will use native Windows MessageBox as fallback")
			}
			os.Exit(1)
		}
	}

	// Check wall broadcast if requested
	if *checkWall {
		if isWallAvailable() {
			fmt.Println("Wall broadcast is available")
			fmt.Println("Can send notifications to all logged-in users")
			os.Exit(0)
		} else {
			if runtime.GOOS != "linux" {
				fmt.Println("Wall broadcast is only available on Linux")
			} else {
				fmt.Println("Wall command not found")
				fmt.Println("Install with: sudo apt install bsdutils (usually pre-installed)")
			}
			os.Exit(1)
		}
	}

	// Force WebView mode if requested (bypass OpenGL check)
	if *forceWebView {
		log.Println("Force-webview mode enabled, skipping OpenGL check")
		if !isWebViewAvailable() {
			log.Fatal("WebView forced but not available. Build with: go build -tags webview")
		}
		log.Println("Using WebView (HTML/CSS/JS) (forced)")
		err := showWebViewNotification(*title, *message, *timeout, icon, *width, *height)
		if err != nil {
			log.Fatalf("Failed to show WebView notification: %v", err)
		}
		os.Exit(0)
	}

	// Force basic GUI mode if requested (bypass OpenGL check)
	if *forceBasic {
		log.Println("Force-basic mode enabled, skipping OpenGL check")
		if runtime.GOOS == "windows" {
			log.Println("Using native Windows MessageBox (forced)")
			err := showWindowsMessageBox(*title, *message, *timeout)
			if err != nil {
				log.Fatalf("Failed to show notification: %v", err)
			}
			os.Exit(0)
		} else {
			log.Println("Force-basic mode only supported on Windows currently")
			// Fall through to normal logic
		}
	}

	// Verify GUI is available before showing notification
	if !isGUIAvailable() {
		// Try wall broadcast on Linux as fallback
		if runtime.GOOS == "linux" && isWallAvailable() {
			log.Println("GUI not available, using wall broadcast")
			err := broadcastWallMessage(*title, *message, *timeout)
			if err != nil {
				log.Fatalf("Failed to broadcast message: %v", err)
			}
			os.Exit(0)
		}
		log.Fatal("GUI mode is not available and no fallback notification method found.")
	}

	// Check OpenGL availability (primarily for Windows)
	openglAvailable := isOpenGLAvailable()
	log.Printf("OpenGL availability check result: %v", openglAvailable)

	if !openglAvailable {
		log.Println("Warning: OpenGL not available, trying alternative GUI")

		// Try WebView first (works on all platforms, better UI)
		if isWebViewAvailable() {
			log.Println("Using WebView (HTML/CSS/JS) for notification")
			err := showWebViewNotification(*title, *message, *timeout, icon, *width, *height)
			if err != nil {
				log.Printf("WebView failed: %v, trying basic fallback", err)
			} else {
				os.Exit(0)
			}
		}

		// Fall back to native OS dialogs as last resort
		if runtime.GOOS == "windows" {
			log.Println("Using native Windows MessageBox")
			err := showWindowsMessageBox(*title, *message, *timeout)
			if err != nil {
				log.Fatalf("Failed to show notification: %v", err)
			}
			os.Exit(0)
		} else {
			log.Fatal("OpenGL not available and no suitable fallback GUI for this platform")
		}
	}

	// Create the notification window with Fyne (when OpenGL is available)
	log.Println("Attempting to create Fyne GUI (OpenGL detected as available)")
	showNotification(*title, *message, *timeout, icon, *width, *height)
}

// showNotification displays a notification window with the given title, message, timeout, optional icon, and window dimensions
func showNotification(title, message string, timeout int, iconPath string, width, height int) {
	// Add panic recovery in case Fyne initialization fails despite OpenGL check
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Fyne GUI failed to initialize (panic): %v", err)
			log.Println("Falling back to alternative notification method")

			// Try fallbacks
			if runtime.GOOS == "windows" {
				if werr := showWindowsMessageBox(title, message, timeout); werr != nil {
					log.Fatalf("All notification methods failed: %v", werr)
				}
			} else {
				log.Fatalf("Fyne GUI failed and no fallback available for this platform")
			}
		}
	}()

	a := app.New()
	w := a.NewWindow(title)

	// Set the window size BEFORE creating content
	// This ensures the layout managers respect our dimensions
	windowSize := fyne.NewSize(float32(width), float32(height))

	// Create the UI
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle.Bold = true

	messageLabel := widget.NewLabel(message)
	messageLabel.Wrapping = fyne.TextWrapWord // Enable word wrapping

	okButton := widget.NewButton("OK", func() {
		w.Close()
	})

	// Create the main content (title, message, button)
	mainContent := container.NewVBox(
		titleLabel,
		widget.NewSeparator(),
		messageLabel,
		widget.NewSeparator(),
		okButton,
	)

	// Add icon if specified
	var content fyne.CanvasObject
	if iconPath != "" {
		iconImage := loadIcon(iconPath)
		if iconImage != nil {
			// Create horizontal layout with icon on the left
			// Use Border layout to ensure message text gets proper width
			iconContainer := container.NewVBox(iconImage)
			content = container.NewBorder(
				nil,                                // top
				nil,                                // bottom
				container.NewPadded(iconContainer), // left (icon)
				nil,                                // right
				container.NewPadded(mainContent),   // center (content gets remaining space)
			)
		} else {
			// If icon fails to load, just use main content
			content = mainContent
		}
	} else {
		content = mainContent
	}

	// Wrap content in a padded container
	paddedContent := container.NewPadded(content)

	w.SetContent(paddedContent)
	w.Resize(windowSize)
	w.SetFixedSize(false) // Allow manual resizing but start at our size
	w.CenterOnScreen()

	// Set up auto-close if timeout is specified
	if timeout > 0 {
		go func() {
			time.Sleep(time.Duration(timeout) * time.Second)
			fyne.DoAndWait(func() {
				w.Close()
			})
		}()
	}

	// Show the window
	w.Show()

	// Force the window to respect our size after showing
	// This is necessary because Fyne may resize based on content
	w.Resize(windowSize)

	// Run the app
	a.Run()
}

// loadIcon loads an image from the specified file path and returns it as a canvas.Image
func loadIcon(iconPath string) *canvas.Image {
	// Check if file exists
	if _, err := os.Stat(iconPath); os.IsNotExist(err) {
		log.Printf("Warning: Icon file not found: %s", iconPath)
		return nil
	}

	// Load the image using Fyne's storage
	uri := storage.NewFileURI(iconPath)
	img := canvas.NewImageFromURI(uri)

	if img == nil {
		log.Printf("Warning: Failed to load icon: %s", iconPath)
		return nil
	}

	// Set image properties
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(64, 64))

	return img
}

// isGUIAvailable checks if GUI mode is available on the current system
func isGUIAvailable() bool {
	switch runtime.GOOS {
	case "linux":
		return isLinuxGUIAvailable()
	case "darwin":
		return isMacGUIAvailable()
	case "windows":
		return isWindowsGUIAvailable()
	default:
		return false
	}
}

func updateChecker(repoOwner string, repo string, repoName string, repodl string) (string, bool) {
	// Create update checker - it will create latestcheck.json in current directory
	uc := updatechecker.New(repoOwner, repo, repoName, repodl, 0, false)
	uc.CheckForUpdate(appVersion)
	updtmsg := uc.Message
	return updtmsg, uc.UpdateAvailable
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
