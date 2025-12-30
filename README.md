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
make build
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
4. Add Todos for today.
  - Type a todo and press Enter to add it.
  - Use `Tab` / `Shift+Tab` to switch focus between the todo input and the backlog/added list.
  - Use `Up`/`Down` (or `k`/`j`) to navigate the list when it has focus.
  - Press `Space` to toggle backlog selection.
  - Press `Enter` on an added todo (when the list has focus) to load it into the input for editing.
  - Leave the input empty and press Enter to finish todos.
5. Answer the questions.
  - `Enter` saves the current answer and advances to the next question.
  - `Shift+Right` moves to the next question. `Shift+Left` moves to the previous question.
  - `Ctrl+S` or `Ctrl+N` still advance.
6. The journal entry will be saved to your configured directory.

CLI: Update todos from terminal

You can update todos for today (or a specific date) without launching the full TUI using the `--todos` flag.

- Update today (recommended):

```bash
./journal --todos ""
# or simply
./journal --todo
```

- Update a specific date (YYYY-MM-DD):

```bash
./journal --todos 2025-12-30
```

Behavior of `--todos` mode:
- The program loads the specified day's markdown file and prompts for each todo, one by one.
- Available responses:
  - `c` or `complete` — mark todo complete.
  - `p` or `partial` — mark as partially completed (appends `(partial)` to the todo text).
  - `n` or `not` — move todo to backlog (it will be carried forward to the next day).
  - any other input — leave the todo unchanged.

Note: shells treat flags-without-values differently. Using `--todos ""` explicitly is reliable across shells to mean "today." If you prefer, I can add a separate boolean flag `--todo-mode` that always updates today's todos.

## Keywords

- journaling
- daily-journal
- cli
- markdown
- obsidian
- productivity
- templates
- go

## Contributing

We welcome contributions! See `CONTRIBUTING.md` for guidelines on filing issues and submitting pull requests. Maintainers will review incoming PRs — the repository uses a review-first workflow and code owners to ensure one or more reviews are required before merging.

If you'd like to help, open an issue or submit a draft PR and we will guide you through the process.

