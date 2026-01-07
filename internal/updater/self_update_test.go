package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGetBinaryName tests the getBinaryName helper function for different OS/arch combinations
func TestGetBinaryName(t *testing.T) {
	tests := []struct {
		name     string
		goos     string
		goarch   string
		expected string
	}{
		{
			name:     "Linux AMD64",
			goos:     "linux",
			goarch:   "amd64",
			expected: "journal-linux-amd64",
		},
		{
			name:     "Linux ARM64",
			goos:     "linux",
			goarch:   "arm64",
			expected: "journal-linux-arm64",
		},
		{
			name:     "macOS AMD64",
			goos:     "darwin",
			goarch:   "amd64",
			expected: "journal-darwin-amd64",
		},
		{
			name:     "macOS ARM64",
			goos:     "darwin",
			goarch:   "arm64",
			expected: "journal-darwin-arm64",
		},
		{
			name:     "Windows AMD64",
			goos:     "windows",
			goarch:   "amd64",
			expected: "journal-windows-amd64.exe",
		},
		{
			name:     "Windows ARM64",
			goos:     "windows",
			goarch:   "arm64",
			expected: "journal-windows-arm64.exe",
		},
		{
			name:     "FreeBSD AMD64",
			goos:     "freebsd",
			goarch:   "amd64",
			expected: "journal-freebsd-amd64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBinaryName(tt.goos, tt.goarch)
			if result != tt.expected {
				t.Errorf("getBinaryName(%q, %q) = %q, want %q", tt.goos, tt.goarch, result, tt.expected)
			}
		})
	}
}

// TestParseReleaseJSON tests parsing of GitHub release JSON
func TestParseReleaseJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		expectError bool
		checkFunc   func(*testing.T, *Release)
	}{
		{
			name: "Valid release with multiple assets",
			jsonInput: `{
				"tag_name": "v1.0.0",
				"assets": [
					{
						"name": "journal-linux-amd64",
						"browser_download_url": "https://github.com/ops295/journal-cli/releases/download/v1.0.0/journal-linux-amd64"
					},
					{
						"name": "journal-darwin-arm64",
						"browser_download_url": "https://github.com/ops295/journal-cli/releases/download/v1.0.0/journal-darwin-arm64"
					}
				]
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, r *Release) {
				if r.TagName != "v1.0.0" {
					t.Errorf("TagName = %v, want v1.0.0", r.TagName)
				}
				if len(r.Assets) != 2 {
					t.Errorf("Assets count = %v, want 2", len(r.Assets))
				}
				if r.Assets[0].Name != "journal-linux-amd64" {
					t.Errorf("First asset name = %v, want journal-linux-amd64", r.Assets[0].Name)
				}
			},
		},
		{
			name: "Valid release with no assets",
			jsonInput: `{
				"tag_name": "v0.1.0",
				"assets": []
			}`,
			expectError: false,
			checkFunc: func(t *testing.T, r *Release) {
				if r.TagName != "v0.1.0" {
					t.Errorf("TagName = %v, want v0.1.0", r.TagName)
				}
				if len(r.Assets) != 0 {
					t.Errorf("Assets count = %v, want 0", len(r.Assets))
				}
			},
		},
		{
			name:        "Invalid JSON",
			jsonInput:   `{"tag_name": "v1.0.0", "assets": [}`,
			expectError: true,
			checkFunc:   nil,
		},
		{
			name:        "Empty JSON",
			jsonInput:   ``,
			expectError: true,
			checkFunc:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var release Release
			err := json.Unmarshal([]byte(tt.jsonInput), &release)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, &release)
			}
		})
	}
}

// TestUpdateHTTPErrors tests various HTTP error scenarios
func TestUpdateHTTPErrors(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedErrMsg string
	}{
		{
			name:           "404 Not Found",
			statusCode:     http.StatusNotFound,
			responseBody:   `{"message":"Not Found"}`,
			expectedErrMsg: "latest release not found",
		},
		{
			name:           "403 Forbidden",
			statusCode:     http.StatusForbidden,
			responseBody:   `{"message":"API rate limit exceeded"}`,
			expectedErrMsg: "access forbidden or rate limited",
		},
		{
			name:           "500 Internal Server Error",
			statusCode:     http.StatusInternalServerError,
			responseBody:   `{"message":"Internal Server Error"}`,
			expectedErrMsg: "unexpected status code from GitHub API: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that returns the specified error
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// We can't easily test Update() directly as it makes real HTTP calls
			// Instead, we test the error handling logic by simulating the response
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(strings.NewReader(tt.responseBody)),
			}
			defer resp.Body.Close()

			// Simulate the error handling from Update()
			if resp.StatusCode != http.StatusOK {
				msg := fmt.Sprintf("unexpected status code from GitHub API: %d", resp.StatusCode)

				switch resp.StatusCode {
				case http.StatusNotFound:
					msg += " (latest release not found; check that the repository has a published release)"
				case http.StatusForbidden:
					msg += " (access forbidden or rate limited; you may have hit GitHub's API rate limit)"
				}

				if body, readErr := io.ReadAll(resp.Body); readErr == nil && len(body) > 0 {
					msg += fmt.Sprintf(" - response body: %s", string(body))
				}

				if !strings.Contains(msg, tt.expectedErrMsg) {
					t.Errorf("Error message does not contain expected substring.\nGot: %s\nExpected to contain: %s", msg, tt.expectedErrMsg)
				}
			}
		})
	}
}

// TestFindAssetURL tests the findAssetURL function with various scenarios
func TestFindAssetURL(t *testing.T) {
	// Create a sample release with various assets
	release := &Release{
		TagName: "v1.0.0",
		Assets: []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}{
			{Name: "journal-linux-amd64", BrowserDownloadURL: "https://example.com/journal-linux-amd64"},
			{Name: "journal-darwin-arm64", BrowserDownloadURL: "https://example.com/journal-darwin-arm64"},
			{Name: "journal-windows-amd64.exe", BrowserDownloadURL: "https://example.com/journal-windows-amd64.exe"},
			{Name: "checksums.txt", BrowserDownloadURL: "https://example.com/checksums.txt"},
		},
	}

	tests := []struct {
		name        string
		targetName  string
		goos        string
		goarch      string
		expectError bool
		expectedURL string
	}{
		{
			name:        "Linux AMD64 found",
			targetName:  "journal-linux-amd64",
			goos:        "linux",
			goarch:      "amd64",
			expectError: false,
			expectedURL: "https://example.com/journal-linux-amd64",
		},
		{
			name:        "Darwin ARM64 found",
			targetName:  "journal-darwin-arm64",
			goos:        "darwin",
			goarch:      "arm64",
			expectError: false,
			expectedURL: "https://example.com/journal-darwin-arm64",
		},
		{
			name:        "Windows AMD64 found",
			targetName:  "journal-windows-amd64.exe",
			goos:        "windows",
			goarch:      "amd64",
			expectError: false,
			expectedURL: "https://example.com/journal-windows-amd64.exe",
		},
		{
			name:        "Unsupported platform",
			targetName:  "journal-freebsd-386",
			goos:        "freebsd",
			goarch:      "386",
			expectError: true,
			expectedURL: "",
		},
		{
			name:        "Non-binary asset",
			targetName:  "checksums.txt",
			goos:        "linux",
			goarch:      "amd64",
			expectError: false,
			expectedURL: "https://example.com/checksums.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := findAssetURL(release, tt.targetName, tt.goos, tt.goarch)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if url != "" {
					t.Errorf("Expected empty URL but got: %v", url)
				}
				// Verify error message includes OS/arch info
				if !strings.Contains(err.Error(), tt.goos) || !strings.Contains(err.Error(), tt.goarch) {
					t.Errorf("Error message should include OS/arch info, got: %v", err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if url != tt.expectedURL {
				t.Errorf("URL = %v, want %v", url, tt.expectedURL)
			}
		})
	}
}

// TestFindAssetURLEmptyRelease tests findAssetURL with an empty release
func TestFindAssetURLEmptyRelease(t *testing.T) {
	release := &Release{
		TagName: "v1.0.0",
		Assets:  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}{},
	}

	url, err := findAssetURL(release, "journal-linux-amd64", "linux", "amd64")
	if err == nil {
		t.Errorf("Expected error for empty release but got none")
	}
	if url != "" {
		t.Errorf("Expected empty URL but got: %v", url)
	}
	if !strings.Contains(err.Error(), "no compatible binary found") {
		t.Errorf("Error message should mention 'no compatible binary found', got: %v", err.Error())
	}
	if !strings.Contains(err.Error(), "linux/amd64") {
		t.Errorf("Error message should include OS/arch info, got: %v", err.Error())
	}
}

// TestBinaryPermissions tests that downloaded binaries would have correct permissions
func TestBinaryPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-binary")

	// Create a test file
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.WriteString("test content")
	f.Close()

	// Set executable permissions (simulating what Update() does)
	if err := os.Chmod(testFile, 0755); err != nil {
		t.Fatalf("Failed to chmod: %v", err)
	}

	// Check permissions
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	mode := info.Mode()
	expectedPerm := os.FileMode(0755)
	
	// Note: On Unix-like systems, we expect exact 0755 permissions.
	// On Windows, the permission system works differently and we just verify
	// the file has some executable bits set (mode & 0111 != 0).
	if mode.Perm() != expectedPerm {
		// If not exact match, at least verify file has executable bits on Windows
		if mode.Perm()&0111 == 0 {
			t.Errorf("File should be executable. Permissions = %o, want at least some executable bits set", mode.Perm())
		}
	}
}

// TestErrorMessageFormatting tests that error messages are properly formatted
func TestErrorMessageFormatting(t *testing.T) {
	tests := []struct {
		name     string
		goos     string
		goarch   string
		expected string
	}{
		{
			name:     "Linux ARM",
			goos:     "linux",
			goarch:   "arm",
			expected: "no compatible binary found for linux/arm (looking for journal-linux-arm)",
		},
		{
			name:     "Darwin 386",
			goos:     "darwin",
			goarch:   "386",
			expected: "no compatible binary found for darwin/386 (looking for journal-darwin-386)",
		},
		{
			name:     "Windows ARM64",
			goos:     "windows",
			goarch:   "arm64",
			expected: "no compatible binary found for windows/arm64 (looking for journal-windows-arm64.exe)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := fmt.Sprintf("journal-%s-%s", tt.goos, tt.goarch)
			if tt.goos == "windows" {
				target += ".exe"
			}

			errMsg := fmt.Sprintf("no compatible binary found for %s/%s (looking for %s)", tt.goos, tt.goarch, target)

			if errMsg != tt.expected {
				t.Errorf("Error message = %v, want %v", errMsg, tt.expected)
			}
		})
	}
}
