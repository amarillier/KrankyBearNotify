# PowerShell example script demonstrating KrankyBear Notify usage

# Path to the krankybearnotify binary
$NotifyBin = "..\krankybearnotify.exe"

# Check if the binary exists
if (-not (Test-Path $NotifyBin)) {
    Write-Host "Error: krankybearnotify.exe not found. Please build it first with 'make build'" -ForegroundColor Red
    exit 1
}

# Check if GUI is available
Write-Host "Checking if GUI is available..." -ForegroundColor Cyan
$result = & $NotifyBin -check-gui
if ($LASTEXITCODE -ne 0) {
    Write-Host "GUI is not available. Cannot show notifications." -ForegroundColor Red
    exit 1
}

Write-Host "GUI is available. Showing example notifications..." -ForegroundColor Green
Write-Host ""

# Example 1: Simple notification
Write-Host "Example 1: Simple notification with defaults" -ForegroundColor Yellow
Start-Process -FilePath $NotifyBin -NoNewWindow
Start-Sleep -Seconds 1

# Example 2: Custom title and message with icon
Write-Host "Example 2: Custom title and message with icon" -ForegroundColor Yellow
if (Test-Path "..\KrankyBearBeret.png") {
    Start-Process -FilePath $NotifyBin -ArgumentList "-title", "Build Complete", "-message", "Your project has been built successfully!", "-icon", "..\KrankyBearBeret.png", "-timeout", "5" -NoNewWindow
} else {
    Start-Process -FilePath $NotifyBin -ArgumentList "-title", "Build Complete", "-message", "Your project has been built successfully!", "-timeout", "5" -NoNewWindow
}
Start-Sleep -Seconds 6

# Example 3: Warning notification with hard hat icon
Write-Host "Example 3: Warning notification with hard hat icon" -ForegroundColor Yellow
if (Test-Path "..\KrankyBearHardHat.png") {
    Start-Process -FilePath $NotifyBin -ArgumentList "-title", "⚠️ Warning", "-message", "Disk space is running low. Please free up some space.", "-icon", "..\KrankyBearHardHat.png", "-timeout", "7" -NoNewWindow
} else {
    Start-Process -FilePath $NotifyBin -ArgumentList "-title", "⚠️ Warning", "-message", "Disk space is running low. Please free up some space.", "-timeout", "7" -NoNewWindow
}
Start-Sleep -Seconds 8

# Example 4: Success notification with fedora icon
Write-Host "Example 4: Success notification with fedora icon" -ForegroundColor Yellow
if (Test-Path "..\KrankyBearFedoraRed.png") {
    Start-Process -FilePath $NotifyBin -ArgumentList "-title", "✓ Success", "-message", "All tests passed!", "-icon", "..\KrankyBearFedoraRed.png", "-timeout", "4" -NoNewWindow
} else {
    Start-Process -FilePath $NotifyBin -ArgumentList "-title", "✓ Success", "-message", "All tests passed!", "-timeout", "4" -NoNewWindow
}
Start-Sleep -Seconds 5

# Example 5: Manual close (no timeout)
Write-Host "Example 5: Manual close notification (click OK to close)" -ForegroundColor Yellow
& $NotifyBin -title "Manual Close" -message "This notification will stay until you click OK" -timeout 0

Write-Host ""
Write-Host "All examples completed!" -ForegroundColor Green
