# Journal CLI

A cross-platform terminal-based daily journaling application written in Go.

## Features
- **Human-first journaling**: Tracks mood, energy, and gratitude.
- **Template-driven**: Customizable templates via YAML.
- **Daily Todos**: Manages daily tasks and automatically carries over unchecked items from the previous day (Backlog).
- **Obsidian-compatible**: Generates Markdown files with frontmatter, ready for your Obsidian vault.
- **Offline & Private**: No database, no cloud, just files on your disk.

## Installation

### Prerequisites
- Go 1.21 or higher

### Build
```bash
go build -o journal cmd/journal/main.go
```

### Run
```bash
./journal
```

## Configuration

The application uses `os.UserConfigDir` for configuration.
- **macOS**: `~/Library/Application Support/journal-cli/`
- **Linux**: `~/.config/journal-cli/`
- **Windows**: `%APPDATA%\journal-cli\`

### config.yaml
Place a `config.yaml` file in the config directory:

```yaml
obsidian_vault: "/Users/username/Documents/ObsidianVault"
journal_dir: "Journal/Daily" # Relative to obsidian_vault
```

### templates
The application looks for YAML template files in the `templates` subdirectory of the config directory:
- **macOS**: `~/Library/Application Support/journal-cli/templates/`
- **Linux**: `~/.config/journal-cli/templates/`
- **Windows**: `%APPDATA%\journal-cli\templates\`

Example `daily-human-dev.yaml`:
```yaml
name: daily-human-dev
description: Human-first daily journal for developers
questions:
  - id: mood
    title: "🧠 How am I feeling today (emotionally)?"
  - id: energy
    title: "⚡ How is my energy level today?"
```

## Usage
1. Run the app.
2. Select a template using Up/Down arrows and Enter.
3. Enter your Mood and Energy.
4. Add Todos for today. Press Enter to add, leave empty and press Enter to finish.
5. Answer the questions. Use Ctrl+S or Ctrl+N to move to the next question.
6. The journal entry will be saved to your configured directory.
