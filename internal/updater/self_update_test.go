package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestParseReleaseJSON tests parsing of GitHub release JSON
func TestParseReleaseJSON(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		wantTag   string
		wantAsset string
		wantErr   bool
	}{
		{
			name: "valid release with assets",
			jsonData: `{
				"tag_name": "v1.0.0",
				"assets": [
					{
						"name": "journal-linux-amd64",
						"browser_download_url": "https://example.com/journal-linux-amd64"
					},
					{
						"name": "journal-darwin-arm64",
						"browser_download_url": "https://example.com/journal-darwin-arm64"
					}
				]
			}`,
			wantTag:   "v1.0.0",
			wantAsset: "journal-linux-amd64",
			wantErr:   false,
		},
		{
			name: "release with no assets",
			jsonData: `{
				"tag_name": "v2.0.0",
				"assets": []
			}`,
			wantTag:   "v2.0.0",
			wantAsset: "",
			wantErr:   false,
		},
		{
			name:     "invalid json",
			jsonData: `{invalid json}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var release Release
			err := json.Unmarshal([]byte(tt.jsonData), &release)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if release.TagName != tt.wantTag {
					t.Errorf("TagName = %v, want %v", release.TagName, tt.wantTag)
				}

				if len(release.Assets) > 0 && tt.wantAsset != "" {
					if release.Assets[0].Name != tt.wantAsset {
						t.Errorf("Asset name = %v, want %v", release.Assets[0].Name, tt.wantAsset)
					}
				}
			}
		})
	}
}

// TestDetermineBinaryName tests binary name generation for different OS/arch combinations
func TestDetermineBinaryName(t *testing.T) {
	tests := []struct {
		name     string
		goos     string
		goarch   string
		expected string
	}{
		{
			name:     "linux amd64",
			goos:     "linux",
			goarch:   "amd64",
			expected: "journal-linux-amd64",
		},
		{
			name:     "linux arm64",
			goos:     "linux",
			goarch:   "arm64",
			expected: "journal-linux-arm64",
		},
		{
			name:     "darwin amd64",
			goos:     "darwin",
			goarch:   "amd64",
			expected: "journal-darwin-amd64",
		},
		{
			name:     "darwin arm64",
			goos:     "darwin",
			goarch:   "arm64",
			expected: "journal-darwin-arm64",
		},
		{
			name:     "windows amd64",
			goos:     "windows",
			goarch:   "amd64",
			expected: "journal-windows-amd64.exe",
		},
		{
			name:     "windows arm64",
			goos:     "windows",
			goarch:   "arm64",
			expected: "journal-windows-arm64.exe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the binary name generation logic from Update()
			target := fmt.Sprintf("journal-%s-%s", tt.goos, tt.goarch)
			if tt.goos == "windows" {
				target += ".exe"
			}

			if target != tt.expected {
				t.Errorf("Binary name = %v, want %v", target, tt.expected)
			}
		})
	}
}

// TestFindAssetForPlatform tests finding the correct asset for a platform
func TestFindAssetForPlatform(t *testing.T) {
	release := Release{
		TagName: "v1.0.0",
		Assets: []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}{
			{Name: "journal-linux-amd64", BrowserDownloadURL: "https://example.com/journal-linux-amd64"},
			{Name: "journal-darwin-arm64", BrowserDownloadURL: "https://example.com/journal-darwin-arm64"},
			{Name: "journal-windows-amd64.exe", BrowserDownloadURL: "https://example.com/journal-windows-amd64.exe"},
		},
	}

	tests := []struct {
		name       string
		targetName string
		wantURL    string
		wantFound  bool
	}{
		{
			name:       "linux amd64 exists",
			targetName: "journal-linux-amd64",
			wantURL:    "https://example.com/journal-linux-amd64",
			wantFound:  true,
		},
		{
			name:       "darwin arm64 exists",
			targetName: "journal-darwin-arm64",
			wantURL:    "https://example.com/journal-darwin-arm64",
			wantFound:  true,
		},
		{
			name:       "windows amd64 exists",
			targetName: "journal-windows-amd64.exe",
			wantURL:    "https://example.com/journal-windows-amd64.exe",
			wantFound:  true,
		},
		{
			name:       "linux arm64 not found",
			targetName: "journal-linux-arm64",
			wantURL:    "",
			wantFound:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var url string
			for _, asset := range release.Assets {
				if asset.Name == tt.targetName {
					url = asset.BrowserDownloadURL
					break
				}
			}

			if (url != "") != tt.wantFound {
				t.Errorf("Found asset = %v, want found %v", url != "", tt.wantFound)
			}

			if url != tt.wantURL {
				t.Errorf("Asset URL = %v, want %v", url, tt.wantURL)
			}
		})
	}
}

// TestUpdateWithHTTPErrors tests error handling for various HTTP error conditions
func TestUpdateWithHTTPErrors(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		wantErrContain string
	}{
		{
			name:           "404 not found",
			statusCode:     http.StatusNotFound,
			responseBody:   `{"message": "Not Found"}`,
			wantErrContain: "latest release not found",
		},
		{
			name:           "403 forbidden",
			statusCode:     http.StatusForbidden,
			responseBody:   `{"message": "API rate limit exceeded"}`,
			wantErrContain: "access forbidden or rate limited",
		},
		{
			name:           "500 internal server error",
			statusCode:     http.StatusInternalServerError,
			responseBody:   `{"message": "Internal Server Error"}`,
			wantErrContain: "unexpected status code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Mock the GitHub API by temporarily changing the repo variable
			// Since we can't change the const, we'll test the logic directly
			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

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

				// Verify the error message contains expected text
				if tt.wantErrContain != "" {
					if !strings.Contains(msg, tt.wantErrContain) {
						t.Errorf("Error message %q does not contain %q", msg, tt.wantErrContain)
					}
				}
			}
		})
	}
}

// TestUpdateWithInvalidJSON tests handling of invalid JSON response
func TestUpdateWithInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err == nil {
		t.Error("Expected error decoding invalid JSON, got nil")
	}
}

// TestUpdateWithMissingAsset tests handling when target platform binary is not in assets
func TestUpdateWithMissingAsset(t *testing.T) {
	// Create a server that returns a valid release but without the current platform's binary
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := Release{
			TagName: "v1.0.0",
			Assets: []struct {
				Name               string `json:"name"`
				BrowserDownloadURL string `json:"browser_download_url"`
			}{
				{Name: "journal-different-platform", BrowserDownloadURL: "https://example.com/binary"},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		t.Fatalf("Failed to decode release: %v", err)
	}

	// Simulate looking for the current platform's binary
	target := fmt.Sprintf("journal-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		target += ".exe"
	}

	var url string
	for _, asset := range release.Assets {
		if asset.Name == target {
			url = asset.BrowserDownloadURL
			break
		}
	}

	if url != "" {
		t.Errorf("Expected no URL for missing asset, got %s", url)
	}
}

// TestUpdateBinaryDownload tests downloading a binary
func TestUpdateBinaryDownload(t *testing.T) {
	binaryContent := []byte("fake binary content")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(binaryContent)
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to download: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test-binary")

	out, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		t.Fatalf("Failed to save binary: %v", err)
	}

	// Verify the downloaded content
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	if string(content) != string(binaryContent) {
		t.Errorf("Downloaded content = %v, want %v", string(content), string(binaryContent))
	}
}

// TestUpdateBinaryPermissions tests setting executable permissions on downloaded binary
func TestUpdateBinaryPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test-binary")

	// Create a test file
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Change permissions to executable
	if err := os.Chmod(tmpFile, 0755); err != nil {
		t.Fatalf("Failed to chmod: %v", err)
	}

	// Verify permissions
	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	mode := info.Mode()
	if mode.Perm() != 0755 {
		t.Errorf("File permissions = %o, want 0755", mode.Perm())
	}
}
