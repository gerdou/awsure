package internal

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

//go:embed extension/*
var extensionFS embed.FS

type BrowserType string

const (
	BrowserChrome  BrowserType = "chrome"
	BrowserFirefox BrowserType = "firefox"
	BrowserUnknown BrowserType = "unknown"
)

func DetectDefaultBrowser() (BrowserType, string, error) {
	switch runtime.GOOS {
	case "darwin":
		return detectBrowserMacOS()
	case "linux":
		return detectBrowserLinux()
	case "windows":
		return detectBrowserWindows()
	default:
		return BrowserUnknown, "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func detectBrowserMacOS() (BrowserType, string, error) {
	if _, err := os.Stat("/Applications/Google Chrome.app"); err == nil {
		return BrowserChrome, "Google Chrome", nil
	}
	if _, err := os.Stat("/Applications/Firefox.app"); err == nil {
		return BrowserFirefox, "Firefox", nil
	}
	return BrowserUnknown, "", fmt.Errorf("no supported browser found (Chrome or Firefox required)")
}

func detectBrowserLinux() (BrowserType, string, error) {
	if _, err := exec.LookPath("google-chrome"); err == nil {
		return BrowserChrome, "Chrome", nil
	}
	if _, err := exec.LookPath("chromium"); err == nil {
		return BrowserChrome, "Chromium", nil
	}
	if _, err := exec.LookPath("firefox"); err == nil {
		return BrowserFirefox, "Firefox", nil
	}
	return BrowserUnknown, "", fmt.Errorf("no supported browser found (Chrome/Chromium or Firefox required)")
}

func detectBrowserWindows() (BrowserType, string, error) {
	chromePath := filepath.Join(os.Getenv("ProgramFiles"), "Google", "Chrome", "Application", "chrome.exe")
	if _, err := os.Stat(chromePath); err == nil {
		return BrowserChrome, "Google Chrome", nil
	}
	chromePath = filepath.Join(os.Getenv("ProgramFiles(x86)"), "Google", "Chrome", "Application", "chrome.exe")
	if _, err := os.Stat(chromePath); err == nil {
		return BrowserChrome, "Google Chrome", nil
	}
	firefoxPath := filepath.Join(os.Getenv("ProgramFiles"), "Mozilla Firefox", "firefox.exe")
	if _, err := os.Stat(firefoxPath); err == nil {
		return BrowserFirefox, "Firefox", nil
	}
	firefoxPath = filepath.Join(os.Getenv("ProgramFiles(x86)"), "Mozilla Firefox", "firefox.exe")
	if _, err := os.Stat(firefoxPath); err == nil {
		return BrowserFirefox, "Firefox", nil
	}
	return BrowserUnknown, "", fmt.Errorf("no supported browser found (Chrome or Firefox required)")
}

func GetExtensionPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	extPath := filepath.Join(homeDir, ".config", "awsure", "extension")
	if err := ensureExtensionExtracted(extPath); err != nil {
		return "", fmt.Errorf("failed to setup extension: %w", err)
	}
	return extPath, nil
}

func ensureExtensionExtracted(extPath string) error {
	manifestPath := filepath.Join(extPath, "manifest.json")
	if _, err := os.Stat(manifestPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(extPath, 0755); err != nil {
		return err
	}

	files := []string{
		"manifest.json",
		"background.js",
		"content.js",
		"popup.html",
		"popup.js",
	}

	for _, file := range files {
		data, err := extensionFS.ReadFile("extension/" + file)
		if err != nil {
			return fmt.Errorf("failed to read embedded %s: %w", file, err)
		}
		destPath := filepath.Join(extPath, file)
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", file, err)
		}
	}

	return nil
}

func InstallExtension(browserType BrowserType) error {
	extPath, err := GetExtensionPath()
	if err != nil {
		return fmt.Errorf("failed to get extension path: %w", err)
	}

	switch browserType {
	case BrowserChrome:
		if err := openChromeExtensionsPage(); err != nil {
			return fmt.Errorf("failed to open Chrome: %w", err)
		}
	case BrowserFirefox:
		if err := openFirefoxDebuggingPage(); err != nil {
			return fmt.Errorf("failed to open Firefox: %w", err)
		}
	default:
		return fmt.Errorf("unsupported browser type")
	}

	fmt.Println("Extension folder:")
	fmt.Println("  " + extPath)
	if browserType == BrowserChrome {
		fmt.Println("1. Enable Developer Mode (top-right)")
		fmt.Println("2. Click 'Load unpacked'")
		fmt.Println("3. Select the folder above")
	} else {
		fmt.Println("1. Click 'Load Temporary Add-on...'")
		fmt.Println("2. Select the manifest.json file in the folder above")
	}
	fmt.Print("Press Enter once the extension is loaded... ")
	var input string
	fmt.Scanln(&input)
	time.Sleep(500 * time.Millisecond)
	return nil
}

func openChromeExtensionsPage() error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", "-a", "Google Chrome", "chrome://extensions/").Start()
	case "linux":
		if _, err := exec.LookPath("google-chrome"); err == nil {
			return exec.Command("google-chrome", "chrome://extensions/").Start()
		}
		return exec.Command("chromium", "chrome://extensions/").Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", "chrome", "chrome://extensions/").Start()
	default:
		return fmt.Errorf("unsupported OS")
	}
}

func openFirefoxDebuggingPage() error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", "-a", "Firefox", "about:debugging#/runtime/this-firefox").Start()
	case "linux":
		return exec.Command("firefox", "about:debugging#/runtime/this-firefox").Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", "firefox", "about:debugging#/runtime/this-firefox").Start()
	default:
		return fmt.Errorf("unsupported OS")
	}
}

func IsExtensionPrepared() bool {
	extPath, err := GetExtensionPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(extPath, "manifest.json"))
	return err == nil
}

func OpenURL(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "linux":
		if _, err := exec.LookPath("xdg-open"); err == nil {
			return exec.Command("xdg-open", url).Start()
		}
		if _, err := exec.LookPath("gio"); err == nil {
			return exec.Command("gio", "open", url).Start()
		}
		return fmt.Errorf("no URL opener found")
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		return fmt.Errorf("unsupported OS")
	}
}
