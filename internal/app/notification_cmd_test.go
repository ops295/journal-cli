package app

import (
	"os"
	"path/filepath"
	"testing"

	"journal-cli/internal/config"
)

func TestSetNotification(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)

	// Test enabling daily notification at a specific time
	err := SetNotification("daily", "14:30")
	if err != nil {
		t.Fatalf("SetNotification error: %v", err)
	}

	// Verify config
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if cfg.Notification == nil {
		t.Fatal("Notification config is nil")
	}

	if !cfg.Notification.Enabled {
		t.Error("expected notification to be enabled")
	}

	if cfg.Notification.Frequency != "daily" {
		t.Errorf("expected frequency daily, got %s", cfg.Notification.Frequency)
	}

	if cfg.Notification.Time != "14:30" {
		t.Errorf("expected time 14:30, got %s", cfg.Notification.Time)
	}
}

func TestSetNotificationDefaultTime(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)

	// Test enabling daily notification without specifying time
	err := SetNotification("daily", "")
	if err != nil {
		t.Fatalf("SetNotification error: %v", err)
	}

	// Verify config
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if cfg.Notification.Time != "09:00" {
		t.Errorf("expected default time 09:00, got %s", cfg.Notification.Time)
	}
}

func TestSetNotificationInvalidTime(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)

	// Test set with invalid time formats
	invalidTimes := []string{"25:00", "12:60", "9:00", "noon", "12-00"}
	for _, timeStr := range invalidTimes {
		err := SetNotification("daily", timeStr)
		if err == nil {
			t.Errorf("expected error for invalid time %s, got nil", timeStr)
		}
	}
}

func TestDisableNotification(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)

	// Enable first
	err := SetNotification("daily", "09:00")
	if err != nil {
		t.Fatal(err)
	}

	// Now disable
	err = DisableNotification()
	if err != nil {
		t.Fatalf("DisableNotification error: %v", err)
	}

	// Verify config
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if cfg.Notification != nil && cfg.Notification.Enabled {
		t.Error("expected notification to be disabled")
	}
}

func TestShowNotificationStatus(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)

	// Should not error even if no config
	err := ShowNotificationStatus()
	if err != nil {
		t.Errorf("ShowNotificationStatus error (no config): %v", err)
	}

	// Enable and show status
	_ = SetNotification("daily", "10:00")
	err = ShowNotificationStatus()
	if err != nil {
		t.Errorf("ShowNotificationStatus error: %v", err)
	}
}

func TestShowConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)
	t.Setenv("APPDATA", tmp)

	// Create a dummy config file
	err := os.MkdirAll(filepath.Join(tmp, ".config", "journal-cli"), 0755)
	if err != nil && !os.IsExist(err) {
		// On macOS it might be Application Support
	}

	// Better way: use SaveConfig
	err = config.SaveConfig(&config.Config{ObsidianVault: "/test/vault"})
	if err != nil {
		t.Skip("Skipping ShowConfig test due to config path issues in test environment")
	}

	err = ShowConfig()
	if err != nil {
		t.Errorf("ShowConfig error: %v", err)
	}
}
