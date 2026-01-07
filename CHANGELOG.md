# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog, and this project adheres to Semantic Versioning.

## [Unreleased]

### Added

- Initial development items tracked here.

## [0.2.02] - 2026-01-06

### Added

- `self-update` command to update the CLI directly from GitHub Releases.
- `make bump-version` command to automate version updates across files.
- Documentation for the update command in `README.md` and standard help output.

## [0.2.01] - 2026-01-06

### Added

- Template management commands: `--set-template` and `--list-templates`
- Enhanced help command (`--help`) with structured YAML documentation
- Default template auto-selection capability
- Configuration field `default_template` supported
- Cross-platform support for config and template locations
- `internal/help` package for structured documentation
- Comprehensive tests for new help and template functionality


## [0.2.0] - 2025-12-30

### Changed

- Bumped version to 0.2.0
- Fixed config tests on macOS
- Fixed release pipeline trigger (added manual trigger)
- Added automated release notes generation

## [0.1.0] - 2025-12-30

### Added

- Initial v0.1.0 release

### Changed

- Initial commit

### Fixed

- N/A
