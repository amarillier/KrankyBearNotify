.PHONY: all build test clean install run help

# Binary name
BINARY_NAME=krankybearnotify

# Build directory
BUILD_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
LDFLAGS=-w -s

# Default target
all: test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) -v

# Build for all platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "==================== Linux Build Instructions ===================="
	@echo "Cross-compiling Fyne apps from macOS to Linux requires fyne-cross"
	@echo ""
	@echo "Option 1 - Use fyne-cross (Docker-based, works from macOS):"
	@echo "  go install github.com/fyne-io/fyne-cross@latest"
	@echo "  fyne-cross linux -arch=amd64,arm64"
	@echo ""
	@echo "Option 2 - Build directly on Linux:"
	@echo "  go build -ldflags=\"-w -s\" -o krankybearnotify"
	@echo "=================================================================="
	@false

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	# GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 -v
	# GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 -v
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 $(GOBUILD) -ldflags="-w -s" -o bin/MacOSARM64/
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 $(GOBUILD) -ldflags="-w -s" -o bin/MacOSAMD64/
	# set executable icon
	./setIcon.sh Resources/Images/KrankyBearBeret.png bin/MacOSARM64/krankybearnotify
	./setIcon.sh Resources/Images/KrankyBearBeret.png bin/MacOSAMD64/krankybearnotify

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@echo "Note: Requires mingw-w64 (brew install mingw-w64 on macOS)"
	# GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" $(GOBUILD) -ldflags="-w -s -H windowsgui" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe -v
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" $(GOBUILD) -ldflags="-w -s -H windowsgui" -o bin/WindowsAMD64/krankybearnotify.exe -v
	# GOOS=windows GOARCH=arm64 CGO_ENABLED=1 CC="x86_64-w64-mingw32-gcc" $(GOBUILD) -ldflags="-w -s -H windowsgui" -o bin/WindowsARM64/krankybearnotify.exe -v
	./setIcon.sh Resources/Images/KrankyBearBeret.png bin/WindowsAMD64/krankybearnotify.exe

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -cover -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BINARY_NAME) /usr/local/bin/

# Run the application with default settings
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Check if GUI is available
check-gui: build
	@echo "Checking GUI availability..."
	./$(BINARY_NAME) -check-gui

# Show application help
app-help: build
	@echo "Showing application help..."
	./$(BINARY_NAME) -help

# Show version
version: build
	@echo "Showing version information..."
	./$(BINARY_NAME) -version

# Show help
help:
	@echo "KrankyBear Notify - Makefile commands:"
	@echo ""
	@echo "  make build          - Build the application for current platform"
	@echo "  make build-all      - Build for all platforms (Linux, macOS, Windows)"
	@echo "  make build-linux    - Build for Linux"
	@echo "  make build-darwin   - Build for macOS (Intel and ARM)"
	@echo "  make build-windows  - Build for Windows"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make bench          - Run benchmarks"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make deps           - Install/update dependencies"
	@echo "  make install        - Install binary to /usr/local/bin"
	@echo "  make run            - Build and run the application"
	@echo "  make check-gui      - Check if GUI is available"
	@echo "  make app-help       - Show application help (flags and options)"
	@echo "  make version        - Show application version"
	@echo "  make help           - Show this help message"
	@echo ""
