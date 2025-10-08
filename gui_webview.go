//go:build webview
// +build webview

package main

import (
	"fmt"
	"time"

	webview "github.com/webview/webview_go"
)

// showWebViewNotification shows a notification using HTML/CSS/JavaScript in a webview
// This is a fallback when OpenGL is not available but webview is
func showWebViewNotification(title, message string, timeout int, iconPath string) error {
	// Create webview
	w := webview.New(false)
	defer w.Destroy()

	w.SetTitle(title)
	w.SetSize(500, 250, webview.HintNone)

	// Build HTML content with embedded CSS and JavaScript
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            padding: 20px;
        }
        .notification-card {
            background: white;
            border-radius: 12px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            padding: 30px;
            max-width: 450px;
            width: 100%%;
            animation: slideIn 0.3s ease-out;
        }
        @keyframes slideIn {
            from {
                transform: translateY(-20px);
                opacity: 0;
            }
            to {
                transform: translateY(0);
                opacity: 1;
            }
        }
        .title {
            font-size: 24px;
            font-weight: bold;
            color: #333;
            margin-bottom: 15px;
            display: flex;
            align-items: center;
        }
        .icon {
            width: 32px;
            height: 32px;
            margin-right: 12px;
            font-size: 32px;
        }
        .message {
            font-size: 16px;
            color: #666;
            line-height: 1.6;
            margin-bottom: 20px;
            white-space: pre-wrap;
        }
        .button-container {
            display: flex;
            justify-content: flex-end;
        }
        .ok-button {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            border: none;
            padding: 10px 30px;
            border-radius: 6px;
            font-size: 16px;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        .ok-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        .ok-button:active {
            transform: translateY(0);
        }
        .timer {
            text-align: right;
            color: #999;
            font-size: 12px;
            margin-top: 10px;
        }
    </style>
</head>
<body>
    <div class="notification-card">
        <div class="title">
            <span class="icon">📢</span>
            <span>%s</span>
        </div>
        <div class="message">%s</div>
        <div class="button-container">
            <button class="ok-button" onclick="closeWindow()">OK</button>
        </div>
        <div class="timer" id="timer"></div>
    </div>
    <script>
        let timeLeft = %d;
        
        function closeWindow() {
            window.external.invoke('close');
        }
        
        function updateTimer() {
            if (timeLeft > 0) {
                document.getElementById('timer').textContent = 'Auto-closing in ' + timeLeft + 's';
                timeLeft--;
                setTimeout(updateTimer, 1000);
            } else if (timeLeft === 0) {
                document.getElementById('timer').textContent = 'Closing...';
                closeWindow();
            }
        }
        
        if (timeLeft > 0) {
            updateTimer();
        }
    </script>
</body>
</html>
`, title, message, timeout)

	w.SetHtml(html)

	// Handle close message from JavaScript
	w.Bind("close", func() {
		w.Terminate()
	})

	// Auto-close timer (backup in case JS doesn't work)
	if timeout > 0 {
		go func() {
			time.Sleep(time.Duration(timeout) * time.Second)
			w.Terminate()
		}()
	}

	w.Run()
	return nil
}

// isWebViewAvailable checks if webview can be used
func isWebViewAvailable() bool {
	// Webview is generally available on all platforms
	// Windows needs WebView2 runtime, but it's widely installed
	return true
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
