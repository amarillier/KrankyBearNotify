@echo off
REM Build script for Windows with WebView support
REM Run this on your Windows VM

echo ============================================================
echo KrankyBear Notify - WebView Build Script
echo ============================================================
echo.

REM Check if Go is installed
where go >nul 2>&1
if errorlevel 1 (
    echo ERROR: Go is not installed!
    echo Please install Go from: https://go.dev/dl/
    pause
    exit /b 1
)

echo [1/4] Go is installed
go version
echo.

REM Get WebView dependency
echo [2/4] Getting WebView dependency...
go get github.com/webview/webview_go
if errorlevel 1 (
    echo ERROR: Failed to get WebView package
    pause
    exit /b 1
)
echo.

REM Tidy modules
echo [3/4] Tidying modules...
go mod tidy
echo.

REM Build with WebView support
echo [4/4] Building with WebView support...
go build -tags webview -o krankybearnotify-webview.exe
if errorlevel 1 (
    echo ERROR: Build failed
    pause
    exit /b 1
)
echo.

echo ============================================================
echo SUCCESS! Built: krankybearnotify-webview.exe
echo ============================================================
echo.
echo Test it with:
echo krankybearnotify-webview.exe -force-webview -title "Test" -message "WebView works!"
echo.
echo Or with countdown:
echo krankybearnotify-webview.exe -force-webview -title "Alert" -message "Auto-closes in 10s" -timeout 10
echo.
pause

