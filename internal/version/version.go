package version

import "fmt"

// Variables injected via -ldflags
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func Info() string {
	return fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, BuildDate)
}

// GetVersion returns the raw Version value
func GetVersion() string { return Version }

// GetBuildInfo returns a multi-line string with build metadata
func GetBuildInfo() string {
	return fmt.Sprintf("Version: %s\nBuild Date: %s\nGit Commit: %s", Version, BuildDate, Commit)
}

// GetVersionString returns a human-friendly version string
func GetVersionString() string {
	if Version == "dev" {
		return fmt.Sprintf("%s (development build)", Version)
	}
	return fmt.Sprintf("v%s", Version)
}

