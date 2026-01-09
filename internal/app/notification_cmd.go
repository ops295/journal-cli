package app

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"journal-cli/internal/config"
)

// SetNotification enables daily notifications at the specified time
func SetNotification(frequency, timeStr string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate frequency
	if frequency != "daily" && frequency != "" {
		return fmt.Errorf("invalid frequency: %s (only 'daily' is supported)", frequency)
	}

	// Default to daily if not specified
	if frequency == "" {
		frequency = "daily"
	}

	// Default time to 09:00 if not specified
	if timeStr == "" {
		timeStr = "09:00"
	}

	// Validate time format (HH:MM)
	timeRegex := regexp.MustCompile(`^([0-1][0-9]|2[0-3]):([0-5][0-9])$`)
	if !timeRegex.MatchString(timeStr) {
		return fmt.Errorf("invalid time format: %s (expected HH:MM in 24-hour format)", timeStr)
	}

	// Initialize notification if nil
	if cfg.Notification == nil {
		cfg.Notification = &config.Notification{}
	}

	// Update notification settings
	cfg.Notification.Enabled = true
	cfg.Notification.Frequency = frequency
	cfg.Notification.Time = timeStr

	// Save config
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✓ Daily notification enabled at %s\n", timeStr)
	fmt.Printf("  Frequency: %s\n", frequency)
	fmt.Printf("\nNote: You'll need to set up your system scheduler to run 'journal' at the specified time.\n")
	fmt.Printf("See documentation for platform-specific instructions.\n")

	return nil
}

// DisableNotification disables notifications
func DisableNotification() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Notification == nil {
		fmt.Println("Notifications are not configured.")
		return nil
	}

	cfg.Notification.Enabled = false

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("✓ Notifications disabled")
	return nil
}

// ShowNotificationStatus displays current notification settings
func ShowNotificationStatus() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Notification Settings")
	fmt.Println(strings.Repeat("─", 40))

	if cfg.Notification == nil || !cfg.Notification.Enabled {
		fmt.Println("Status: Disabled")
		return nil
	}

	fmt.Printf("Status:    Enabled\n")
	fmt.Printf("Frequency: %s\n", cfg.Notification.Frequency)
	fmt.Printf("Time:      %s\n", cfg.Notification.Time)

	return nil
}

// TriggerNotification checks if it's time to notify and shows an OS notification
func TriggerNotification() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Notification == nil || !cfg.Notification.Enabled {
		fmt.Println("Notifications are not enabled.")
		return nil
	}

	title := "Journal Reminder"
	message := fmt.Sprintf("Time to write in your journal! (Current setting: %s)", cfg.Notification.Time)

	// Trigger OS-specific notification
	switch runtime.GOOS {
	case "darwin":
		err = showMacOSNotification(title, message)
	case "linux":
		err = showLinuxNotification(title, message)
	case "windows":
		err = showWindowsNotification(title, message)
	default:
		err = fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if err != nil {
		return fmt.Errorf("failed to show notification: %w", err)
	}

	fmt.Printf("✓ %s notification triggered successfully\n", strings.Title(runtime.GOOS))
	return nil
}

// showMacOSNotification displays a notification on macOS using osascript
func showMacOSNotification(title, message string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "default"`, message, title)
	return executeCommand("osascript", "-e", script)
}

// showLinuxNotification displays a notification on Linux using notify-send
func showLinuxNotification(title, message string) error {
	return executeCommand("notify-send", title, message)
}

// showWindowsNotification displays a notification on Windows using PowerShell
func showWindowsNotification(title, message string) error {
	// A simple PowerShell command to show a balloon tip notification
	psCommand := fmt.Sprintf(
		`$wshell = New-Object -ComObject WScript.Shell; $wshell.Popup("%s", 0, "%s", 64)`,
		message, title,
	)
	return executeCommand("powershell", "-Command", psCommand)
}

// ShowConfig displays the configuration path and content
func ShowConfig() error {
	configPath, err := config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	fmt.Printf("Configuration file: %s\n", configPath)

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		fmt.Println("\nConfiguration file does not exist (using defaults).")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	fmt.Println("\nContent:")
	fmt.Println(strings.Repeat("─", 40))
	fmt.Println(string(data))
	fmt.Println(strings.Repeat("─", 40))

	return nil
}

// executeCommand runs a command with arguments
func executeCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}
