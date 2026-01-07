package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const repo = "ops295/journal-cli"

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// Update handles the self-update process
func Update() error {
	fmt.Println("🔍 Checking latest version...")

	// 1. Fetch latest release info
	resp, err := http.Get("https://api.github.com/repos/" + repo + "/releases/latest")
	if err != nil {
		return fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Provide clearer guidance about common GitHub API errors and include the response body.
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

		return fmt.Errorf("%s", msg)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to decode release info: %w", err)
	}

	// 2. Determine target binary name
	// Naming convention assumed: journal-darwin-arm64, journal-linux-amd64, etc.
	target := fmt.Sprintf(
		"journal-%s-%s",
		runtime.GOOS,
		runtime.GOARCH,
	)
	
	// Windows binaries usually have .exe extension
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

	if url == "" {
		return fmt.Errorf("no compatible binary found for %s/%s (looking for %s)", runtime.GOOS, runtime.GOARCH, target)
	}

	fmt.Printf("⬇️  Downloading %s (%s)...\n", release.TagName, target)

	// 3. Download the new binary to a temporary file
	tmpFile := filepath.Join(os.TempDir(), "journal-new")
	// Ensure we don't conflict if multiple runs or stale files
	_ = os.Remove(tmpFile)
	
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	respBin, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer respBin.Body.Close()

	if respBin.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code during download: %d", respBin.StatusCode)
	}

	if _, err := io.Copy(out, respBin.Body); err != nil {
		return fmt.Errorf("failed to save binary: %w", err)
	}

	// Make the temp file executable
	if err := out.Chmod(0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}
	
	// Close the file explicitly before renaming to ensure all writes are flushed
	out.Close()

	// 4. Replace the current binary
	current, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to locate current executable: %w", err)
	}
	
	// Resolve symlinks if any (common in some installs), though os.Executable usually handles this.
	// We'll stick to what os.Executable returns for now.

	backup := current + ".bak"

	fmt.Println("🔄 Replacing binary...")
	
	// First move the current binary to .bak
	if err := os.Rename(current, backup); err != nil {
		// If permission denied, give a helpful hint
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied: try running with sudo/admin privileges")
		}
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Then move the new binary to the current location
	if err := os.Rename(tmpFile, current); err != nil {
		// Rollback: try to restore the backup
		_ = os.Rename(backup, current)
		return fmt.Errorf("failed to install new binary: %w", err)
	}
	
	// Cleanup backup
	_ = os.Remove(backup)

	fmt.Println("✅ Update successful! Please run 'journal --version' to verify.")
	return nil
}
