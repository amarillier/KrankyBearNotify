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
	appVersion     = "0.1.0"
	appAuthor      = "Allan Marillier"
	defaultTitle   = "Notification"
	defaultMessage = "This is a notification message"
	defaultTimeout = 10 // seconds
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

SUPPORTED PLATFORMS:
  • macOS 10.13+
  • Windows 10+
  • Linux (X11/Wayland) - Works on GNOME, KDE, XFCE, Cinnamon, MATE, and more

For more information, visit: https://github.com/amarillier/krankybearnotify
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	}
}

func main() {
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

	// Command-line flags
	title := flag.String("title", defaultTitle, "Notification title")
	message := flag.String("message", defaultMessage, "Notification message")
	timeout := flag.Int("timeout", defaultTimeout, "Timeout in seconds (0 for no timeout)")
	icon := flag.String("icon", "", "Path to icon image file (PNG, JPEG, etc.)")
	checkGUI := flag.Bool("check-gui", false, "Check if GUI mode is available and exit")
	checkOpenGL := flag.Bool("check-opengl", false, "Check if OpenGL is available and exit")
	version := flag.Bool("version", false, "Show version information and exit")

	// Update checker flags (with alias)
	var checkUpdate bool
	flag.BoolVar(&checkUpdate, "checkupdate", false, "Check for updates and exit")
	flag.BoolVar(&checkUpdate, "cu", false, "Check for updates and exit (alias for -checkupdate)")

	flag.Parse()

	// Show help if no arguments provided
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

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

	// Verify GUI is available before showing notification
	if !isGUIAvailable() {
		log.Fatal("GUI mode is not available. Cannot display notification.")
	}

	// Check OpenGL availability (primarily for Windows)
	if !isOpenGLAvailable() {
		log.Println("Warning: OpenGL not available, using native fallback GUI")
		if runtime.GOOS == "windows" {
			err := showWindowsMessageBox(*title, *message, *timeout)
			if err != nil {
				log.Fatalf("Failed to show notification: %v", err)
			}
			os.Exit(0)
		} else {
			log.Fatal("OpenGL not available and no fallback GUI for this platform")
		}
	}

	// Create the notification window with Fyne
	showNotification(*title, *message, *timeout, *icon)
}

// showNotification displays a notification window with the given title, message, timeout, and optional icon
func showNotification(title, message string, timeout int, iconPath string) {
	a := app.New()
	w := a.NewWindow(title)

	// Create the UI
	titleLabel := widget.NewLabel(title)
	titleLabel.TextStyle.Bold = true

	messageLabel := widget.NewLabel(message)
	messageLabel.Wrapping = 1 // Enable text wrapping

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
			iconContainer := container.NewVBox(iconImage)
			content = container.NewHBox(
				iconContainer,
				widget.NewSeparator(),
				mainContent,
			)
		} else {
			// If icon fails to load, just use main content
			content = mainContent
		}
	} else {
		content = mainContent
	}

	w.SetContent(content)
	w.Resize(fyne.NewSize(500, 200))
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

	// Show the window and run the app
	w.ShowAndRun()
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
