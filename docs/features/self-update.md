# Feature: Self-Update

The `self-update` feature allows users to update their CLI binary to the latest version available on GitHub Releases without needing an external package manager.

## Usage

```bash
journal self-update
```

If installed in a protected directory (like `/usr/local/bin`):
```bash
sudo journal self-update
```

## Implementation Details

Located in: `internal/updater`

1.  **Check**: Queries `https://api.github.com/repos/ops295/journal-cli/releases/latest`.
2.  **Match**: Looks for an asset matching the pattern `journal-{GOOS}-{GOARCH}` (e.g., `journal-darwin-arm64`).
3.  **Download**: Streams the binary to a temporary file in `os.TempDir`.
4.  **Replace**:
    -   Moves the current binary to `{binary}.bak`.
    -   Moves the new binary to the current location.
    -   Restores from backup if the move fails.

## Constraints
-   The binary must have write permissions to its own location.
-   Access to GitHub API is required.
