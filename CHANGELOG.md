# Changelog

All notable changes to KrankyBear Notify will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-07

### Added
- Initial release of KrankyBear Notify
- Cross-platform GUI notification support (macOS, Windows, Linux)
- Command-line interface with customizable title, message, and timeout
- GUI availability detection for all platforms:
  - Linux: DISPLAY, WAYLAND_DISPLAY, and systemd graphical.target support
  - macOS: WindowServer detection
  - Windows: Window station detection
- `-check-gui` flag to verify GUI availability without showing notifications
- Comprehensive test suite with platform-specific tests
- Auto-dismiss functionality with configurable timeout
- Modern UI built with Fyne v2.6.3
- Example scripts for bash (Linux/macOS) and PowerShell (Windows)
- Integration examples for:
  - Cron jobs
  - systemd services
  - Windows Task Scheduler
  - Python, Node.js, and shell scripts
- Makefile with common build and test targets
- Cross-compilation support for multiple platforms
- Comprehensive documentation and README
- .gitignore for Go projects

### Features
- Display custom notifications with title and message
- Configurable auto-close timeout (0 for manual close)
- Graceful handling of headless environments
- Text wrapping for long messages
- Centered window positioning
- Clean, modern UI with separators and clear buttons

### Testing
- Unit tests for GUI detection on all platforms
- Cross-platform stub function tests
- Environment variable handling tests
- Benchmark tests for performance monitoring
- Platform-specific test skipping for CI/CD compatibility

### Documentation
- Comprehensive README with usage examples
- Platform-specific notes and troubleshooting
- Example scripts directory with integration patterns
- Makefile help system
- Inline code documentation

## [Unreleased]

## [1.1.0] - 2025-10-07

### Added
- **Custom icon support**: New `-icon` flag to display user-specified images in notifications
- Icon loading from PNG, JPEG, and other image formats
- Graceful fallback when icon files are not found or fail to load
- Horizontal layout with icon displayed on the left side of notification content
- Icon sizing (64x64 pixels) with proper aspect ratio handling
- Support for both relative and absolute icon paths

### Changed
- Window width increased from 400 to 500 pixels to accommodate icons
- Updated all example scripts to demonstrate icon usage
- Enhanced documentation with icon usage examples

### Tests
- Added `TestLoadIcon` for icon loading functionality
- Added `TestIconPathHandling` for various path scenarios
- Tests verify graceful handling of missing or invalid icon files

## [1.0.0] - 2025-10-07

### Planned Features
- System tray icon support
- Multiple notification styles (info, warning, error, success)
- Sound notifications
- Notification history
- Configuration file support
- D-Bus integration for Linux
- Native notification APIs (macOS Notification Center, Windows Action Center)
- Multi-monitor support with display selection
- Notification positioning options
- Rich text formatting in messages
- Clickable action buttons
- Progress bar support for long-running tasks

### Known Issues
- None at this time

---

## Version History

- **1.1.0** (2025-10-07): Added custom icon support
- **1.0.0** (2025-10-07): Initial release with core functionality
