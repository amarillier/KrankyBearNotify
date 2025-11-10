package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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
		fmt.Fprintf(os.Stderr, `Notify v%s
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

  # Custom button text
  %s -title "Confirm" -message "Please review" -button "Got it!" -timeout 0

  # URL-encoded parameters (automatically decoded)
  %s -title "Build%%20Complete" -message "Path:%%20/home/user%%2Fproject" -icon "%%2Fpath%%2Fto%%2Ficon.png"

  # Check if GUI is available (useful for scripts)
  %s -check-gui

  # Check for missing runtime dependencies (Linux)
  %s -check-deps

  # Check for updates
  %s -cu

  # Notification that stays until manually closed
  %s -title "Important" -message "Please review" -timeout 0

  # Windows: Force MessageBox mode (for VMs where OpenGL fails)
  %s -win-basic -title "VM Alert" -message "Uses Windows MessageBox"

  # Windows: Force WebView mode (better UI, requires webview build)
  %s -win-webview -title "Modern Alert" -message "Uses HTML/CSS/JS"

  # Linux: Send to GUI users only (no wall broadcast)
  %s -gui-only -title "GUI Alert" -message "Only GUI users see this"

  # Linux: Force wall broadcast only (no GUI)
  %s -force-wall -title "Terminal Alert" -message "Only terminal users see this"

SUPPORTED PLATFORMS:
  • macOS 10.13+
  • Windows 10+
  • Linux (X11/Wayland) - Works on GNOME, KDE, XFCE, Cinnamon, MATE, and more
    - Headless/SSH: Falls back to 'wall' broadcast when no GUI detected

For more information, visit: https://github.com/amarillier/krankybearnotify
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	}
}

func main() {
	// CRITICAL: Handle version flag BEFORE any other code runs
	// This prevents Fyne GUI initialization which can hang in some environments

	// Windows: Hide console window IMMEDIATELY if running as target user
	// Check os.Args directly before flag parsing for fastest hide
	if runtime.GOOS == "windows" {
		hasTargetUser := false
		hasDebug := false
		for _, arg := range os.Args[1:] {
			if arg == "-target-user" {
				hasTargetUser = true
			}
			if arg == "-debug" {
				hasDebug = true
			}
		}
		if hasTargetUser && !hasDebug {
			hideConsoleWindow()
		}
	}

	// Windows 7 compatibility check - must be early to prevent crashes
	if runtime.GOOS == "windows" {
		if isWindows7() {
			fmt.Fprintf(os.Stderr, "Error: Not supported on Windows 7\n")
			fmt.Fprintf(os.Stderr, "This application requires Windows 10 or later.\n")
			fmt.Fprintf(os.Stderr, "Please upgrade your operating system.\n")
			os.Exit(1)
		}
	}

	// Quick pre-check for version flag to avoid GUI initialization
	for _, arg := range os.Args[1:] {
		if arg == "-version" || arg == "--version" {
			// Declare glibcver outside the Linux-specific blocks so it's in scope for both
			glibcver := ""
			if runtime.GOOS == "linux" {
				glibcVer, glibcErr := getGlibcVersion()
				if glibcErr != nil {
					glibcver = "(glibc version undetected)"
				} else {
					glibcver = glibcVer
				}
			}
			fmt.Printf("Notify: v%s\n", appVersion)
			if runtime.GOOS == "linux" {
				fmt.Printf("Platform: %s/%s (glibc %s)\n", runtime.GOOS, runtime.GOARCH, glibcver)
			} else {
				fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
			}
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
	title := flag.String("title", defaultTitle, "Notification title (URL/percent-encoded characters will be decoded)")
	message := flag.String("message", defaultMessage, "Notification message (URL/percent-encoded characters will be decoded)")
	buttonText := flag.String("button", "OK", "Button text (URL/percent-encoded characters will be decoded)")
	timeout := flag.Int("timeout", defaultTimeout, "Timeout in seconds (0 for no timeout)")
	width := flag.Int("width", defaultWidth, "Window width in pixels")
	height := flag.Int("height", defaultHeight, "Window height in pixels")
	autosize := flag.Bool("autosize", false, "Auto-size window based on message length (max 600x400)")
	checkGUI := flag.Bool("check-gui", false, "Check if GUI mode is available and exit")
	checkOpenGL := flag.Bool("check-opengl", false, "Check if OpenGL is available and exit")
	checkWall := flag.Bool("check-wall", false, "Check if wall broadcast is available (Linux) and exit")
	checkDeps := flag.Bool("check-deps", false, "Check for missing runtime dependencies (Linux) and exit")
	winBasic := flag.Bool("win-basic", false, "Windows: Force basic mode (MessageBox instead of Fyne)")
	winWebView := flag.Bool("win-webview", false, "Windows: Force WebView mode (requires -tags webview build)")
	guiOnly := flag.Bool("gui-only", false, "Linux: Send to GUI users only (no wall broadcast)")
	forceWall := flag.Bool("force-wall", false, "Linux: Force wall broadcast only (no GUI)")
	targetUser := flag.Bool("target-user", false, "Internal: Marks process as already running as target user (prevents re-elevation)")
	debug := flag.Bool("debug", false, "Enable debug output (shows log messages)")
	version := flag.Bool("version", false, "Show version information and exit")

	// Icon flag with alias
	var icon string
	flag.StringVar(&icon, "icon", "", "Path to icon image file (PNG, JPEG, etc.) (URL/percent-encoded characters will be decoded)")
	flag.StringVar(&icon, "image", "", "Path to icon image file (alias for -icon) (URL/percent-encoded characters will be decoded)")

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

	// Suppress unused variable warning for targetUser
	// This flag is checked in shouldShowToOtherUsers() via os.Args
	_ = targetUser

	// Configure logging based on debug flag
	// When running via scheduled task (target-user), default to quiet unless debug is enabled
	if !*debug {
		// When running via scheduled task with -target-user, log to file for debugging
		if *targetUser && runtime.GOOS == "windows" {
			logFile, err := os.OpenFile("C:\\Temp\\notify-debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				log.SetOutput(logFile)
				defer logFile.Close()
				log.Printf("=== Notify started via scheduled task ===")
				log.Printf("Args: %v", os.Args)
			} else {
				// Fallback to discard if can't open log file
				log.SetOutput(io.Discard)
			}
		} else {
			// Disable all log output (no console spam)
			log.SetOutput(io.Discard)
		}
	}

	// URL decode title, message, button text, and icon parameters
	// This handles percent-encoded characters like %2d (-), %2f (/), %20 (space), etc.
	if decodedTitle, err := url.QueryUnescape(*title); err == nil {
		*title = decodedTitle
	} else {
		log.Printf("Warning: Failed to URL decode title: %v", err)
	}
	if decodedMessage, err := url.QueryUnescape(*message); err == nil {
		*message = decodedMessage
	} else {
		log.Printf("Warning: Failed to URL decode message: %v", err)
	}
	if decodedButtonText, err := url.QueryUnescape(*buttonText); err == nil {
		*buttonText = decodedButtonText
	} else {
		log.Printf("Warning: Failed to URL decode button text: %v", err)
	}
	if icon != "" {
		if decodedIcon, err := url.QueryUnescape(icon); err == nil {
			icon = decodedIcon
		} else {
			log.Printf("Warning: Failed to URL decode icon path: %v", err)
		}

		// Add .png extension if no extension provided
		// This ensures all modes (Fyne, WebView, MessageBox) get the same icon path processing
		ext := filepath.Ext(icon)
		if ext == "" {
			icon = icon + ".png"
			log.Printf("No extension provided, added .png: %s", icon)
		}
	}

	// Show version if requested
	if *version {
		fmt.Printf("Notify v%s\n", appVersion)
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

	// Check dependencies if requested (Linux only)
	if *checkDeps {
		if runtime.GOOS == "linux" {
			checkLinuxDependencies()
		} else {
			fmt.Println("Dependency check is only available on Linux")
			os.Exit(1)
		}
	}

	// Check GUI mode if requested
	if *checkGUI {
		if isGUIAvailable() {
			fmt.Println("GUI mode is available")
			// On Linux, also check for missing libraries
			if runtime.GOOS == "linux" {
				checkLinuxDependenciesQuiet()
			}
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

	// Force wall broadcast mode if requested (Linux only)
	if *forceWall {
		if runtime.GOOS != "linux" {
			log.Fatal("Force-wall mode is only available on Linux")
		}
		if !isWallAvailable() {
			log.Fatal("Wall command not found. Install with: sudo apt install bsdutils")
		}
		log.Println("Force-wall mode enabled, using wall broadcast")
		err := broadcastWallMessage(*title, *message, *timeout)
		if err != nil {
			log.Fatalf("Failed to send wall broadcast: %v", err)
		}
		os.Exit(0)
	}

	// Windows: Force WebView mode if requested (bypass OpenGL check)
	// BUT skip if running as SYSTEM with other users (will be handled by elevated notification logic)
	if *winWebView {
		if runtime.GOOS != "windows" {
			log.Fatal("-win-webview flag is only supported on Windows")
		}

		// If running as SYSTEM with logged-in users, defer to the elevated notification handler
		if shouldShowToOtherUsers() {
			log.Println("-win-webview flag detected, but running as SYSTEM with logged-in users")
			log.Println("Will launch as target user (flag will be passed to child process)")
			// Continue to the elevated notification logic below
		} else {
			log.Println("Windows WebView mode enabled, skipping OpenGL check")
			if !isWebViewAvailable() {
				log.Fatal("WebView not available. Build with: go build -tags webview")
			}
			log.Println("Using WebView (HTML/CSS/JS)")
			err := showWebViewNotification(*title, *message, *timeout, icon, *width, *height, *buttonText)
			if err != nil {
				log.Fatalf("Failed to show WebView notification: %v", err)
			}
			os.Exit(0)
		}
	}

	// Windows: Force basic mode if requested (bypass OpenGL check)
	// BUT skip if running as SYSTEM with other users (will be handled by elevated notification logic)
	if *winBasic {
		if runtime.GOOS != "windows" {
			log.Fatal("-win-basic flag is only supported on Windows")
		}

		// If running as SYSTEM with logged-in users, defer to the elevated notification handler
		if shouldShowToOtherUsers() {
			log.Println("-win-basic flag detected, but running as SYSTEM with logged-in users")
			log.Println("Will launch as target user (flag will be passed to child process)")
			// Continue to the elevated notification logic below
		} else {
			log.Println("Windows basic mode enabled, using MessageBox")
			err := showWindowsMessageBox(*title, *message, *timeout)
			if err != nil {
				log.Fatalf("Failed to show notification: %v", err)
			}
			os.Exit(0)
		}
	}

	// Special handling when running as root/SYSTEM/Administrator
	// Show to BOTH GUI users and terminal users (Linux only has wall)
	if shouldShowToOtherUsers() {
		log.Println("Running with elevated privileges, notifying logged-in users")

		guiSuccess := false
		wallSuccess := false

		// Try to show GUI to logged-in GUI users (unless force-wall is set)
		if !*forceWall {
			if err := showNotificationToUsers(*title, *message, *timeout, icon, *width, *height, *buttonText); err == nil {
				log.Println("✓ Notification shown to GUI user(s)")
				guiSuccess = true
			} else {
				log.Printf("✗ Could not show GUI to users: %v", err)
			}
		}

		// Linux-specific: Send wall broadcast to terminal sessions
		// Skip if -gui-only flag is set
		if runtime.GOOS == "linux" && !*guiOnly && isWallAvailable() {
			if *forceWall {
				log.Println("Sending wall broadcast only (force-wall mode)")
			} else {
				log.Println("Also sending wall broadcast to terminal sessions")
			}
			err := broadcastWallMessage(*title, *message, *timeout)
			if err != nil {
				log.Printf("✗ Wall broadcast failed: %v", err)
			} else {
				log.Println("✓ Wall broadcast sent to terminal users")
				wallSuccess = true
			}
		}

		// Exit if at least one method succeeded
		if guiSuccess || wallSuccess {
			os.Exit(0)
		}

		// If both failed, check if we're running as SYSTEM on Windows
		// SYSTEM doesn't have a desktop, so don't try to show GUI to SYSTEM itself
		if runtime.GOOS == "windows" && isRunningAsSystem() {
			log.Println("ERROR: Running as SYSTEM but could not notify any users via scheduled task")
			log.Println("SYSTEM account has no desktop/display - cannot show GUI directly")
			log.Fatal("Notification failed: No logged-in users found or scheduled task creation failed")
		}

		// If both failed, log and continue to try normal GUI
		log.Println("Warning: Could not notify via GUI or wall, trying normal GUI mode")
	}

	// Auto-size window if requested
	if *autosize {
		calculatedWidth, calculatedHeight := calculateWindowSize(*title, *message, *buttonText, icon != "")
		// Use calculated size but respect user-provided maximums
		if *width == defaultWidth {
			*width = calculatedWidth
		}
		if *height == defaultHeight {
			*height = calculatedHeight
		}
		log.Printf("Auto-sizing enabled: calculated %dx%d, using %dx%d", calculatedWidth, calculatedHeight, *width, *height)
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

		// On Windows, when running as SYSTEM, skip WebView (it has permission issues)
		// and go directly to MessageBox which is more reliable for services
		skipWebView := false
		if runtime.GOOS == "windows" && isRunningAsSystem() {
			log.Println("Running as SYSTEM on Windows, skipping WebView (using MessageBox for reliability)")
			skipWebView = true
		}

		// Try WebView first (works on all platforms, better UI) unless skipped
		if !skipWebView && isWebViewAvailable() {
			log.Println("Using WebView (HTML/CSS/JS) for notification")
			err := showWebViewNotification(*title, *message, *timeout, icon, *width, *height, *buttonText)
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
	showNotification(*title, *message, *timeout, icon, *width, *height, *buttonText)
}

// showNotification displays a notification window with the given title, message, timeout, optional icon, window dimensions, and button text
func showNotification(title, message string, timeout int, iconPath string, width, height int, buttonText string) {
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
	w.SetIcon(resourceKrankyBearBeretPng)

	// Windows-specific: Add zombie process prevention timeout
	// In VMs without proper OpenGL, Fyne may hang invisibly without crashing
	if runtime.GOOS == "windows" {
		// Calculate a reasonable zombie prevention timeout
		// Use the larger of: (user timeout + 15 seconds) or 30 seconds minimum
		zombieTimeout := timeout + 15
		if zombieTimeout < 30 {
			zombieTimeout = 30
		}

		go func() {
			time.Sleep(time.Duration(zombieTimeout) * time.Second)
			log.Printf("Warning: Zombie prevention timeout reached (%d seconds), forcing exit", zombieTimeout)

			// Try graceful quit using DoAndWait (proper Fyne thread-safe call)
			go func() {
				defer func() {
					// Catch any panic from Quit()
					if r := recover(); r != nil {
						log.Printf("Panic during graceful quit (expected if hung): %v", r)
					}
				}()
				fyne.DoAndWait(func() {
					a.Quit()
				})
			}()

			// Force exit immediately - if we've reached timeout, Fyne is hung anyway
			time.Sleep(100 * time.Millisecond) // Brief moment for quit attempt
			log.Printf("Forcing process termination")
			os.Exit(0)
		}()
		log.Printf("Zombie prevention timeout set: %d seconds", zombieTimeout)
	}

	// Set the window size BEFORE creating content
	// This ensures the layout managers respect our dimensions
	windowSize := fyne.NewSize(float32(width), float32(height))

	// Create the UI
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle.Bold = true

	messageLabel := widget.NewLabel(message)
	messageLabel.Wrapping = fyne.TextWrapWord // Enable word wrapping

	okButton := widget.NewButton(buttonText, func() {
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

// calculateWindowSize calculates optimal window dimensions based on content
// Returns width and height capped at reasonable maximums
func calculateWindowSize(title, message, buttonText string, hasIcon bool) (int, int) {
	// Base dimensions
	minWidth := 300
	minHeight := 150
	maxWidth := 600
	maxHeight := 400

	// Estimate based on text length
	// Average character width: ~7 pixels for normal text
	// Average line height: ~20 pixels

	// Calculate width based on longest line in message
	messageWidth := estimateTextWidth(message)
	titleWidth := estimateTextWidth(title)
	buttonWidth := 100 + len(buttonText)*7 // Button has padding

	// Use the longest element
	contentWidth := messageWidth
	if titleWidth > contentWidth {
		contentWidth = titleWidth
	}
	if buttonWidth > contentWidth {
		contentWidth = buttonWidth
	}

	// Add padding and icon space
	width := contentWidth + 60 // 30px padding on each side
	if hasIcon {
		width += 80 // Space for icon
	}

	// Apply width constraints BEFORE calculating line count
	// This ensures line count is based on the actual available width
	if width < minWidth {
		width = minWidth
	}
	if width > maxWidth {
		width = maxWidth
	}

	// Calculate height based on message lines (using constrained width)
	messageLines := estimateLineCount(message, width-60)
	titleLines := 1
	if len(title) > 50 {
		titleLines = 2
	}

	// Calculate total height
	height := 40 + // Top padding
		(titleLines * 30) + // Title
		(messageLines * 25) + // Message lines
		50 + // Button
		30 // Bottom padding

	// Apply height constraints
	if height < minHeight {
		height = minHeight
	}
	if height > maxHeight {
		height = maxHeight
	}

	return width, height
}

// estimateTextWidth estimates the pixel width of text
func estimateTextWidth(text string) int {
	const avgCharWidth = 7
	maxLineLength := 0

	// Split by newlines and find longest line
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if len(line) > maxLineLength {
			maxLineLength = len(line)
		}
	}

	return maxLineLength * avgCharWidth
}

// estimateLineCount estimates how many lines the text will take with word wrapping
func estimateLineCount(text string, availableWidth int) int {
	if text == "" {
		return 1
	}

	const avgCharWidth = 7
	charsPerLine := availableWidth / avgCharWidth

	if charsPerLine <= 0 {
		charsPerLine = 40 // Fallback
	}

	// Count lines considering word wrap
	words := strings.Fields(text)
	lines := 1
	currentLineLength := 0

	for _, word := range words {
		wordLength := len(word) + 1 // +1 for space
		if currentLineLength+wordLength > charsPerLine {
			lines++
			currentLineLength = wordLength
		} else {
			currentLineLength += wordLength
		}
	}

	// Add explicit newlines
	lines += strings.Count(text, "\n")

	return lines
}

// resolveIconPath resolves an icon path by looking in the executable's directory if it's just a filename
// Returns the resolved path that should be used to load the icon
func resolveIconPath(iconPath string) string {
	if iconPath == "" {
		return ""
	}

	// Determine the actual path to use
	actualPath := iconPath

	// Check if the path is just a filename (no directory separators)
	if filepath.Base(iconPath) == iconPath {
		// It's just a filename, look in the executable's directory
		exePath, err := os.Executable()
		if err != nil {
			log.Printf("Warning: Could not determine executable path: %v", err)
		} else {
			exeDir := filepath.Dir(exePath)
			exeDirIconPath := filepath.Join(exeDir, iconPath)

			// Check if the file exists in the executable's directory
			if _, err := os.Stat(exeDirIconPath); err == nil {
				log.Printf("Found icon in executable directory: %s", exeDirIconPath)
				actualPath = exeDirIconPath
			} else {
				log.Printf("Icon not found in executable directory (%s), trying current directory", exeDirIconPath)
			}
		}
	}

	return actualPath
}

// loadIcon loads an image from the specified file path and returns it as a canvas.Image
// If only a filename is provided (no directory separators), it will look for the file
// in the executable's directory first, then fall back to the current directory
// Note: .png extension is added earlier in main() if needed, so iconPath should already have an extension
func loadIcon(iconPath string) *canvas.Image {
	if iconPath == "" {
		return nil
	}

	// Resolve the icon path (look in exe directory if needed)
	actualPath := resolveIconPath(iconPath)

	// Check if file exists at the determined path
	if _, err := os.Stat(actualPath); os.IsNotExist(err) {
		log.Printf("Warning: Icon file not found: %s", actualPath)
		return nil
	}

	log.Printf("Loading icon from: %s", actualPath)

	// Convert to absolute path to ensure Fyne can find it
	absPath, err := filepath.Abs(actualPath)
	if err != nil {
		log.Printf("Warning: Could not get absolute path for icon: %v", err)
		absPath = actualPath
	} else {
		log.Printf("Absolute icon path: %s", absPath)
	}

	// Load the image using Fyne's storage
	// Note: NewFileURI handles Windows paths correctly, including spaces
	uri := storage.NewFileURI(absPath)
	log.Printf("Icon URI: %s", uri.String())

	img := canvas.NewImageFromURI(uri)

	if img == nil {
		log.Printf("Warning: Failed to load icon from URI: %s (path: %s)", uri.String(), absPath)
		return nil
	}

	log.Printf("Successfully loaded icon: %s", absPath)

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

// isWindows7 checks if the current system is running Windows 7
func isWindows7() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	// Use the 'ver' command to get Windows version
	cmd := exec.Command("ver")
	output, err := cmd.Output()
	if err != nil {
		// If we can't determine version, assume it's not Windows 7
		// This prevents false positives on newer systems
		return false
	}

	versionStr := strings.ToLower(string(output))

	// Windows 7 version strings typically contain "6.1"
	// Examples: "Microsoft Windows [Version 6.1.7601]" or "Microsoft Windows [Version 6.1.7600]"
	return strings.Contains(versionStr, "6.1")
}

func updateChecker(repoOwner string, repo string, repoName string, repodl string) (string, bool) {
	// Create update checker - it will create latestcheck.json in current directory
	uc := updatechecker.New(repoOwner, repo, repoName, repodl, 0, false)
	uc.CheckForUpdate(appVersion)
	updtmsg := uc.Message
	return updtmsg, uc.UpdateAvailable
}

func getGlibcVersion() (string, error) {
	glibcver, glibcerr := exec.Command("getconf", "GNU_LIBC_VERSION").Output()
	if glibcerr == nil {
		// Trim whitespace and newlines from the output
		return strings.TrimSpace(string(glibcver)), nil
	}
	return "", glibcerr
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
