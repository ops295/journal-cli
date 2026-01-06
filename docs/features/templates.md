# Feature: Templates

Templates allow users to define custom prompts for their daily journal entries.

## Configuration

Templates are YAML files stored in the `templates/` subdirectory of the config folder.

### Location
-   **macOS**: `~/Library/Application Support/journal-cli/templates/`
-   **Linux**: `~/.config/journal-cli/templates/`
-   **Windows**: `%APPDATA%\journal-cli\templates\`

## File Format

Example `daily-reflection.yaml`:

```yaml
name: daily-reflection
description: A simple daily reflection
questions:
  - id: gratitude
    title: "What are you grateful for?"
  - id: improvement
    title: "What could have gone better?"
```

## Management Commands

-   **List**: `journal --list-templates`
-   **Set Default**: `journal --set-template <name>`
    -   Example: `journal --set-template daily-reflection`
    -   Logic: Sets `default_template` in `config.yaml`. The app will skip the template selection screen and auto-load this template.

## Implementation Details

Located in: `internal/template` (loading) and `internal/app` (management).

-   `template.LoadTemplates()`: Scans the directory and parses valid YAML files.
-   `app.SetDefaultTemplate()`: Updates the main `config.yaml` file.
