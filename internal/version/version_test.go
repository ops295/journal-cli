package version

import (
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("GetVersion() returned empty string")
	}
}

func TestGetBuildInfo(t *testing.T) {
	info := GetBuildInfo()
	
	// Should contain all three components
	if !strings.Contains(info, "Version:") {
		t.Error("GetBuildInfo() missing Version field")
	}
	if !strings.Contains(info, "Build Date:") {
		t.Error("GetBuildInfo() missing Build Date field")
	}
	if !strings.Contains(info, "Git Commit:") {
		t.Error("GetBuildInfo() missing Git Commit field")
	}
}

func TestGetVersionString(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{
			name:     "development version",
			version:  "dev",
			expected: "dev (development build)",
		},
		{
			name:     "release version",
			version:  "1.0.0",
			expected: "v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original version
			originalVersion := Version
			defer func() { Version = originalVersion }()

			// Set test version
			Version = tt.version
			
			result := GetVersionString()
			if result != tt.expected {
				t.Errorf("GetVersionString() = %v, want %v", result, tt.expected)
			}
		})
	}
}
